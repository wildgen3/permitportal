---
title: Regulatory-relevance scoring is candidate generation only, never an include/exclude gate
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0017
---

# ADR-0017: Relevance scoring never gates inclusion

## Context and problem statement

There is published work mapping regulatory text to industry codes at scale. Can it be
used to decide which regulations apply to a business?

## Decision drivers

- The best published classifier for this mapping posts **mean F1 around 56% with mean
  recall around 50%** [SRC-14]. Half of the relevant text is missed.
- Using it as a filter means the missed half is never surfaced, never reviewed, and
  never known to be missing.
- Using it as an ordering signal for a human curation queue costs nothing when it is
  wrong.

## Considered options

1. **Use relevance scores to select applicable regulations.**
2. **Use them to rank a curation queue.**
3. **Do not use them.**

## Decision outcome

Chosen: **option 2.** Never use industry-relevance as an include/exclude gate.

### Consequences

- Good: recall is protected, which is the quantity that matters when a false negative is
  a fine.
- Bad: per-jurisdiction, per-regime curation remains unavoidable human work — the
  largest and least glamorous line item in the build. This ADR does not reduce it; it
  refuses to pretend it has been reduced.
- Enforced by: relevance scores are not an input to any rule predicate. A rule that
  attempted to reference one would fail linter L-04, since no such attribute is
  registered.

A related discipline: the normalized metrics published alongside that dataset are
internally inconsistent with its own methodology guide and are never cited in this
repository. See `docs/03-source-register.md`.
