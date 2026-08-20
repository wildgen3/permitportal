---
title: Apache-2.0 for code, CC-BY-4.0 for documentation and reference data
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0015
---

# ADR-0015: Apache-2.0 for code, CC-BY-4.0 for docs

## Context and problem statement

This repository is mostly specification, partly code, and is intended to be read by
people at companies whose legal teams filter what their engineers may look at. What
licence?

## Decision drivers

- An explicit patent grant matters the moment this becomes a product or is evaluated by
  a company's counsel.
- A strong copyleft licence on a reference architecture gets it filtered out by exactly
  the audience it is written for.
- Documentation and data are not software, and licensing them as software is a category
  error a sharp reader will notice.

## Considered options

1. **GPL-3.0 throughout.**
2. **MIT throughout.**
3. **Apache-2.0 for code, CC-BY-4.0 for documentation and data.**

## Decision outcome

Chosen: **option 3.**

MIT is fine but has no patent grant, which is the one thing Apache-2.0 adds that matters
here. GPL-3.0 would be actively counterproductive.

### Consequences

- Good: readable by the intended audience without a licence review.
- Bad: two licence files and a `NOTICE`, which must be kept accurate.
- Enforced by: `LICENSE`, `LICENSE-DOCS`, and a licensing section in the README.
