# Environment: dev

**Status: not yet written.** Created in phase 4, after `../../bootstrap` has run.

Intended shape: `min_instances = 0` on Cloud Run, the smallest database tier, and
`budget.tf` in the same pull request as the first billable resource.

`backend.tf` takes the bucket name emitted by `terraform -chdir=infra/bootstrap output
tfstate_bucket`.
