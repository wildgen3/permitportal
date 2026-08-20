---
title: Terraform bootstrap
status: scaffolded
last_reviewed: 2026-08-19
---

# Bootstrap

Run **once**, with local state, then never again. Everything this creates is a
precondition for the environments in `../envs/`, including the bucket their remote state
lives in — which is why this step cannot itself use remote state.

## What it creates

| Resource | Why |
| --- | --- |
| Terraform state bucket | Versioned, uniform bucket-level access, public access prevented |
| Enabled APIs | Cloud Run, Cloud SQL, Artifact Registry, Cloud Build, Secret Manager, Monitoring |
| Deploy service account | Least privilege: `run.developer`, `artifactregistry.writer`, `secretmanager.secretAccessor` |
| Workload Identity pool and provider | So GitHub Actions authenticates **without a service account key** |
| Billing budget | Created here rather than later, so no environment can exist without one |

## Sequence

First, create `terraform.tfvars` (gitignored -- it carries your project id, not this
repository's):

```hcl
project_id      = "your-gcp-project-id"
github_repo     = "your-org/permitgraph"
billing_account = "000000-000000-000000"
```

Then:

```bash
gcloud auth application-default login
gcloud config set project "$(grep project_id infra/bootstrap/terraform.tfvars | cut -d'"' -f2)"

terraform -chdir=infra/bootstrap init
terraform -chdir=infra/bootstrap plan      # read this before applying
terraform -chdir=infra/bootstrap apply

# Copy the emitted bucket name into ../envs/dev/backend.tf, then:
terraform -chdir=infra/envs/dev init
terraform -chdir=infra/envs/dev apply
```

## Rules

- **No service account keys.** If a key would solve the problem, Workload Identity
  Federation solves it without a credential to leak. The bootstrap creates no keys and
  the environments accept none.
- **The budget is not optional.** A cloud project without a budget alert is a project
  that surprises you.
- The target project is whatever `terraform.tfvars` names. It is deliberately not committed.

## Teardown

`terraform -chdir=infra/envs/dev destroy` removes the environment. The bootstrap is left
in place deliberately: destroying the state bucket destroys the record of everything
else.
