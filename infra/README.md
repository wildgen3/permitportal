---
title: Infrastructure
status: scaffolded
last_reviewed: 2026-08-19
---

# Infrastructure

Terraform. **Every** cloud resource. Nothing is created by hand in the console.

```
bootstrap/   Run once, with local state. Creates the state bucket, enables APIs,
             creates the deploy service account and the Workload Identity pool.
             Takes project_id from an untracked terraform.tfvars -- no project id
             is committed to this repository.
modules/     Reusable, single-purpose modules.
envs/dev/    Development environment. min_instances = 0, smallest database tier.
envs/prod/   Not created until there is something to run in it.
```

## Rules

- **`budget.tf` ships in the same pull request as the first billable resource.** Not the next one.
- **No service account keys, ever.** GitHub Actions authenticates through Workload Identity
  Federation to a least-privilege deploy service account.
- **Cloud Run scales to zero.** In dev, idle cost is approximately nothing; the standing cost is the
  database, so dev uses the smallest tier and is destroyable.
- Images build in Cloud Build. There is no local Docker daemon on the development machine — podman
  runs containers, it does not build the deployment images.

`terraform fmt -check` and `terraform validate` run in CI from day one, before any resource exists.
`terraform plan` on pull requests arrives with the first real environment, behind Workload Identity.

## Bootstrap

See [`bootstrap/README.md`](bootstrap/README.md) for the exact sequence.
