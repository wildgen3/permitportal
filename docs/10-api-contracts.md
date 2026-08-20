---
title: API contracts
status: scaffolded
last_reviewed: 2026-08-19
adr: [0005, 0009]
---

# API contracts

> **Status: scaffolded.** `spec/api/openapi.yaml` is not yet written; this document is its
> narrative companion. **Minimum viable version:** the resource model and the two
> behaviours below, both of which are architectural rather than cosmetic.

## Two behaviours that are not negotiable

- **A determination request against an unconfirmed classification returns 409**, not a
  determination. The database constraint would reject it anyway; the API says so
  explicitly rather than surfacing a storage error.
- **Every read accepts `as_of_law` and `as_of_knowledge`.** Bitemporal queries are the
  default shape, not an advanced feature. (ADR-0009)

## To specify

- Resource model: subjects, classifications, determinations, obligations, credential
  chains.
- The intake contract: free text in, ranked candidates out, with confidence and an
  explicit abstain.
- The INDETERMINATE response shape, carrying `missing_attributes[]` — this is what drives
  the next question, so it is part of the contract, not an error payload.
- Evidence-tree serialization.
- Pagination, idempotency, and versioning.

OpenAPI 3.1, `$ref`-ing into `spec/generated/core.schema.json` so the API cannot drift
from the model. Linted in the `contracts` gate.
