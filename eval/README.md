---
title: Evaluation
status: scaffolded
last_reviewed: 2026-08-19
---

# Evaluation

Top-level rather than nested under a service, because there are two evaluable subjects and one
cross-cutting concern:

| Path | Measures |
| --- | --- |
| `classification/` | Recall@k and top-1 accuracy at 6 digits, **per sector**, against a named vintage |
| `obligations/` | Obligation-set correctness against hand-adjudicated synthetic profiles |
| `harness/` | The runner, metrics, cost accounting, and report emitter |
| `fixtures/` | Synthetic business descriptions from a committed seeded generator |
| `baselines/` | Committed baselines. CI fails on regression against these. |

## The headline metric is false negatives

For obligation resolution, a missed obligation becomes a fine. False-negative rate on must-apply
fixtures is gated at **zero**; accuracy is secondary.

## On the published baseline

The Census BEACON figures (0.901 code accuracy, 0.755 end-to-end) are reported for comparison,
decomposed the same way (`p_success = p_return × p_select|return`) so the numbers are actually
comparable rather than rhetorically adjacent. Three caveats travel with them wherever they appear: the
denominator excludes the "not listed" cohort, the end-to-end figure comes from a 2021 pre-test
instrument, and the population differs from ours.

**Gating is against our own committed baseline.** Claiming to beat 75.5% on a different instrument
would be exactly the overclaim this project is built to avoid.

## Cost

`eval.yml` runs against recorded fixtures and calls no paid model. `eval-live.yml` is
`workflow_dispatch` only and is the sole workflow in this repository that can spend money.
