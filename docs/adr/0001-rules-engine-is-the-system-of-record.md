---
title: The deterministic rules engine is the system of record; the AI layer never decides
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0001
---

# ADR-0001: The deterministic rules engine is the system of record

## Context and problem statement

This system tells businesses what the law requires of them. Being wrong means fines and
legal exposure. Where, if anywhere, may a language model sit in the path of a
legally consequential answer?

## Decision drivers

- Published guidance on high-impact applications declines to recommend generative AI as
  the decision-maker in exactly this quadrant.
- Determinations get appealed. "Why did it say that in March" must have an answer, which
  requires reproducibility a sampled model cannot provide.
- Retrieval and language understanding genuinely outperform hand-written heuristics at
  turning "I do residential electrical work" into candidate classifications.

## Considered options

1. **Model decides, human reviews** — fastest to build, familiar pattern.
2. **Model proposes, deterministic engine decides** — a hard seam between the two.
3. **No model at all** — forms and decision trees only.

## Decision outcome

Chosen: **Model proposes, deterministic engine decides.**

Recall-weighted retrieval and model judgment for candidate generation, with
deterministic, citable rules for statutory triggers, stored provenance on every
surfaced obligation, and an independent completeness check.

Option 3 is defensible but throws away the one thing that makes this product better
than a PDF: a person who cannot classify their own business still gets an answer.
Option 1 is indefensible — a sampled decision is not reproducible, and an
unreproducible determination cannot be appealed.

### Consequences

- Good: every determination is replayable and citable.
- Good: model quality changes cannot silently change legal outcomes.
- Bad: two systems to build and keep in contract with each other.
- Enforced by: `scripts/check-engine-purity.py`, which fails the build if anything in
  `packages/rules` or `services/resolver` imports or declares a dependency on a model
  client. If the AI layer ever becomes load-bearing, the build breaks. The evaluator is
  not written yet, so the gate is armed rather than satisfied, and it reports that
  distinction rather than printing a bare pass. A database `CHECK` additionally prevents
  an unconfirmed classification from reaching a determination.

See ADR-0005, ADR-0009, and `docs/11-ai-layer.md`.
