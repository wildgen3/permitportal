---
title: Legacy integration and the anti-corruption layer
status: scaffolded
last_reviewed: 2026-08-19
adr: [0003, 0013]
---

# Legacy integration and the anti-corruption layer

> **Status: scaffolded.** **Minimum viable version:** the adapter contract and the
> no-inbound-trust rule.

## The premise

Agency systems of record stay where they are. Nothing here proposes that a state replace
its licensing system, because no state is going to, and a design that requires it is a
design that never ships.

## To specify

- The adapter contract per agency system: what is read, what is never written, and how
  identity is resolved without a shared business identifier.
- **The no-inbound-trust rule.** A code supplied by an external system enters as an
  *unconfirmed candidate* and must be confirmed like any other. There is no "import your
  existing classification and trust it" path — the error rate of self- and
  agency-assigned codes in the wild is unmeasured, and treating them as ground truth
  would import that unknown directly into determinations.
- Identity resolution and record linkage across systems that share no key.
- The strangler-fig sequence: which capability moves first, and what has to be true
  before the next one moves.
- Failure and reconciliation: what happens when the upstream and our model disagree.
