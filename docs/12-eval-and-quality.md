---
title: Evaluation and quality
status: specified
last_reviewed: 2026-08-19
adr: [0005, 0017, 0018]
---

# Evaluation and quality

## The headline metric is false negatives

For obligation resolution, a missed obligation becomes a fine and a spurious question is
an inconvenience. **False-negative rate on must-apply fixtures is gated at zero.**
Accuracy is secondary and reported for information.

## What is measured

### Classification retrieval

- `recall@{1,3,5,10}` at six digits against a **named vintage**, reported **per sector**.
- Top-1 six-digit accuracy.
- Abstention rate.

Five qualifiers travel with every classification number in this repository — the digit
level, the scheme vintage, the sector breakdown, whether the denominator includes
abstentions, and the instrument. A figure without them is not comparable to anything.

### End-to-end

Decomposed as `p_success = p_return × p_select|return`, so the numbers are actually
comparable to the published baseline rather than rhetorically adjacent to it.

### The published baseline, with its caveats

The government's own tool for this task reports **0.901** code accuracy and **0.755**
end-to-end success [SRC-10]. Three caveats travel with those figures wherever they appear:

1. The denominator excludes the "not listed" cohort and back-outs.
2. The end-to-end figure comes from a pre-test instrument, not the production tool.
3. The population differs from ours.

**Gating is against our own committed baseline in `eval/baselines/`, never against
theirs.** Claiming to beat 0.755 [SRC-10] on a different instrument with a different population
would be precisely the overclaim this project is built to avoid, and it is the first
thing a competent reviewer would catch.

Per-sector spread is the reason for per-sector policy: the weakest sectors sit near 0.55
while the strongest exceed 0.9 [SRC-10]. Only the aggregate figures and the sector *ordering* are
supportable — specific per-sector accuracy figures attributed to that benchmark did not
survive verification and are listed in `03-source-register.md`.

### Rules correctness

- **False-negative rate on must-apply fixtures — gated at zero.**
- Obligation-set exact match against hand-adjudicated synthetic profiles.
- INDETERMINATE rate, and mean questions-to-resolution.

### Calibration

Per-sector expected calibration error, reliability diagrams, and a risk-coverage curve.
The operating point is chosen by cost-of-being-wrong, not by accuracy.

### Explanation faithfulness

Fraction of citations in generated text that are present in the evidence tree. Gated at
100% [internal]. This is not a metric with a target; it is a correctness property.

### Cost

Dollars per 1,000 classifications, tokens per session, p50 and p95 latency.

### Inter-rater agreement on the gold set itself

Published. Independent reviewers demonstrably disagree on the same business description,
so a gold set built by one person and presented as ground truth is overclaiming.
Publishing our own agreement statistic is a credibility asset, not an embarrassment —
and it bounds how much of the residual error is measurement rather than model.

## CI gates

These fail the build:

| # | Gate |
| --- | --- |
| 1 | `recall@5` (six-digit, per sector) at or above the committed baseline minus 1.0 point, with no single sector down more than 3 points |
| 2 | Rules gold set: **zero false negatives** on must-apply fixtures |
| 3 | Citation faithfulness exactly 100% [internal] |
| 4 | Rule linter L-01 … L-06: zero errors |
| 5 | Schema drift: regenerated artifacts identical to committed |
| 6 | Determinism: identical inputs, rule version, and `as_of_law` produce a **byte-identical** evidence tree |
| 7 | Cost per 1,000 at or below the committed ceiling times 1.25 |
| 8 | `packages/rules/engine` import graph contains no model client |

Gates 4, 5, and 8 run today. Gates 1, 2, 3, 6, and 7 arrive with the vertical slice.
`eval/baselines/` currently holds honest placeholders, clearly marked as such — a
baseline file containing invented numbers would be worse than an empty one.

## Fixtures

Synthetic, from a committed seeded generator. Adjudication protocol, sample sizes, and
the labeling procedure are in `22-fixtures.md`. Roughly 150–300 profiles per regime is
the target for a gold set to be worth gating on.

## Cost discipline

`eval.yml` runs against **recorded** fixtures and calls no paid model — it is free and
runs on every relevant pull request. `eval-live.yml` is `workflow_dispatch` only and is
the sole workflow in this repository that can spend money.
