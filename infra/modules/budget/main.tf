# Billing budget with threshold alerts.
#
# This module exists so that no environment can be created without one. A cloud project
# without a budget alert is a project that surprises you.

terraform {
  required_version = ">= 1.15"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
}

variable "project_id" {
  type = string
}

variable "billing_account" {
  type = string
}

variable "amount_usd" {
  description = "Monthly budget in USD."
  type        = number
  default     = 50
}

variable "thresholds" {
  description = "Fractions of the budget at which to alert."
  type        = list(number)
  default     = [0.5, 0.9, 1.0]
}

resource "google_billing_budget" "this" {
  billing_account = var.billing_account
  display_name    = "permitportal-${var.project_id}"

  budget_filter {
    projects = ["projects/${var.project_id}"]
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units         = tostring(var.amount_usd)
    }
  }

  dynamic "threshold_rules" {
    for_each = var.thresholds
    content {
      threshold_percent = threshold_rules.value
    }
  }
}
