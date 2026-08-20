# Bootstrap. Run once, with LOCAL state, then never again.
#
# This creates the bucket that every other environment's remote state lives in, which
# is why it cannot itself use remote state.

terraform {
  required_version = ">= 1.15"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
  # No backend block: local state, deliberately. See README.md.
}

variable "project_id" {
  description = "Target GCP project. Required -- supply via terraform.tfvars, which is gitignored."
  type        = string
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "github_repo" {
  description = "owner/repo that may assume the deploy identity."
  type        = string
}

variable "billing_account" {
  description = "Billing account ID for the budget. Empty disables the budget resource."
  type        = string
  default     = ""
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# --- APIs -------------------------------------------------------------------
resource "google_project_service" "required" {
  for_each = toset([
    "run.googleapis.com",
    "sqladmin.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "secretmanager.googleapis.com",
    "monitoring.googleapis.com",
    "iamcredentials.googleapis.com",
  ])
  service            = each.value
  disable_on_destroy = false
}

# --- Terraform state --------------------------------------------------------
resource "google_storage_bucket" "tfstate" {
  name                        = "${var.project_id}-permitportal-tfstate"
  location                    = var.region
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  force_destroy               = false

  versioning {
    enabled = true
  }

  lifecycle {
    # Destroying this destroys the record of everything else.
    prevent_destroy = true
  }
}

# --- Deploy identity, with NO service account key ---------------------------
resource "google_service_account" "deployer" {
  account_id   = "permitportal-deployer"
  display_name = "PermitPortal CI deployer"
  description  = "Assumed by GitHub Actions via Workload Identity Federation. No keys."
}

resource "google_project_iam_member" "deployer" {
  for_each = toset([
    "roles/run.developer",
    "roles/artifactregistry.writer",
    "roles/secretmanager.secretAccessor",
  ])
  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.deployer.email}"
}

resource "google_iam_workload_identity_pool" "github" {
  workload_identity_pool_id = "permitportal-github"
  display_name              = "PermitPortal GitHub Actions"
}

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-oidc"

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
  }

  # Scoped to this repository. Without this condition any GitHub repository on the
  # internet could assume the identity.
  attribute_condition = "assertion.repository == '${var.github_repo}'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account_iam_member" "github_impersonation" {
  service_account_id = google_service_account.deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_repo}"
}

# --- Budget, created here so no environment can exist without one -----------
module "budget" {
  source          = "../modules/budget"
  count           = var.billing_account == "" ? 0 : 1
  project_id      = var.project_id
  billing_account = var.billing_account
  amount_usd      = 50
}

output "tfstate_bucket" {
  description = "Copy this into ../envs/dev/backend.tf."
  value       = google_storage_bucket.tfstate.name
}

output "deployer_service_account" {
  value = google_service_account.deployer.email
}

output "workload_identity_provider" {
  description = "Pass to google-github-actions/auth as workload_identity_provider."
  value       = google_iam_workload_identity_pool_provider.github.name
}
