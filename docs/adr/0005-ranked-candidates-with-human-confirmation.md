---
title: Classification output is a ranked candidate set with human confirmation
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0005
---

# ADR-0005: Ranked candidates with human confirmation

## Context and problem statement

Free-text classification is imperfect. What does the system do with an imperfect
classification?

## Decision drivers

- The published benchmark for a government classification tool: **90.1% correct code,
  but 75.5% end-to-end success** [SRC-10]. The gap is presentation, not retrieval. Recall is not
  the problem; candidate presentation is.
- Accuracy varies enormously by sector — roughly 0.55 in the worst sectors against
  above 0.9 in the best [SRC-10]. A single global confidence threshold is wrong everywhere.
- Independent reviewers disagree on the same business description. There is no single
  correct answer to defer to.

## Considered options

1. **Auto-assign the top code.**
2. **Auto-assign above a confidence threshold, ask otherwise.**
3. **Always present a ranked candidate set for confirmation.**

## Decision outcome

Chosen: **always present a ranked candidate set for confirmation**, with "none of
these" as a first-class path rather than an error state.

Never let an unconfirmed code flow into a compliance determination.

### Consequences

- Good: the applicant owns the classification, and the appeal record shows what they
  were offered.
- Bad: an extra interaction on every intake.
- Enforced by: a database `CHECK` — a `Determination` may only reference a
  `ClassificationAssignment` whose `confirmation_state` is not `unconfirmed`.

`alternatives_shown[]` is persisted, because an appeal needs to know what the applicant
was offered, not only what they picked. Per-sector policy (`top_k`, margin, mandatory
clarifying question) is a config table, not a constant.
