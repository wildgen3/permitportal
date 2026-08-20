---
title: Every determination is reproducible by replay
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0018
---

# ADR-0018: Determinations are reproducible by replay

## Context and problem statement

An applicant is told an obligation applies. Six months later they dispute it. What does
the system owe them?

## Decision drivers

- Government determinations get appealed. An answer that cannot be reconstructed cannot
  be defended, and cannot be shown to be wrong either.
- Rule data, engine code, and the applicant's own facts all change independently.
- A determination whose explanation is regenerated at read time can drift from the
  determination itself.

## Considered options

1. **Store the outcome.**
2. **Store the outcome and a rendered explanation.**
3. **Store the outcome plus everything needed to recompute it, and replay on demand.**

## Decision outcome

Chosen: **option 3.** Every `Determination` persists `rule_version_id`,
`engine_version`, `as_of_law`, and `input_snapshot_hash`, and the evidence tree is
stored rather than regenerated.

### Consequences

- Good: the appeal path is a replay, not an argument.
- Good: an engine change that alters historical outcomes is detectable rather than
  silent.
- Bad: storage grows with determinations, and rule versions can never be deleted.
- Enforced today by: the canonical model, which makes `rule_version_id`,
  `engine_version`, `as_of_law` and `input_snapshot_hash` required on `Determination`,
  and by `scripts/check-engine-purity.py`, which keeps sampled model output out of the
  decision plane (ADR-0001).
- **Will be enforced by** a property test asserting that identical inputs, rule version
  and `as_of_law` produce a **byte-identical** evidence tree. That test needs an
  evaluator, which does not exist yet.
