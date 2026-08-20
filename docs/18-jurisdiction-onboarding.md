---
title: Jurisdiction onboarding
status: scaffolded
last_reviewed: 2026-08-19
adr: [0013, 0014]
---

# Jurisdiction onboarding

> **Status: scaffolded.** **Minimum viable version:** the source-availability rubric, which
> gates the pilot jurisdiction choice before any rules are authored.

## Why this document is first-class

Per-jurisdiction ingestion is **unavoidable scope, not an optimization**, and it is the
largest and least glamorous line item in any real build of this system. Specifications
that treat it as an implementation detail are the ones that die in month three.

## The source-availability rubric

Applied *before* committing to a jurisdiction, verified live:

1. Are the state code and administrative code retrievable in bulk?
2. Is the licence and endorsement catalogue machine-readable?
3. Does the municipal code have a stable, citable source?
4. Are fees and bond amounts published with effective dates?

Two or more failures means choose a different jurisdiction. Selecting one on *assumed*
data availability is the characteristic way this kind of project fails.

## To specify

- The onboarding playbook, step by step.
- A **measured** person-days budget, taken from the pilot rather than guessed, and a
  fifty-state extrapolation stated as a range with its measurement basis.
- The coverage matrix, which ships **in the product UI** and not only in these docs.
  "We have onboarded N of 50" is a user-visible fact, because a user in an unsupported
  state deserves to know that immediately.

## The gate that matters

The second jurisdiction tests the **playbook**, not the concept. If jurisdiction two
costs more than twice jurisdiction one in person-days, stop and fix the playbook before
jurisdiction three.
