---
title: Fixtures
status: scaffolded
last_reviewed: 2026-08-19
adr: [0014]
---

# Fixtures

> **Status: scaffolded.** **Minimum viable version:** the synthesis rule and the
> adjudication protocol — both are clean-room controls, not test tooling.

All fixture data is **synthetic**, produced by a committed, seeded generator. There is no
ingestion path by which real business data enters this repository. Costs, hour counts,
and fee amounts appearing in specification fixtures are illustrative and are not
represented as current published schedules.

## To specify

- The generator: seed handling, the distribution it samples from, and why it is
  committed (a fixture set nobody can reproduce is not evidence of anything).
- Coverage design: every rule needs at least one positive and one negative fixture
  (linter L-05), and every regime needs enough profiles to gate on — roughly 150–300.
- The adjudication protocol: who labels, against what, and how disagreements resolve.
- **Inter-rater agreement, published.** Independent reviewers disagree on the same
  business description, so a gold set built by one person and presented as ground truth
  overclaims. Publishing the agreement statistic bounds how much residual error is
  measurement rather than model.
- Fixture versioning, so a baseline is always attributable to a specific fixture set.
