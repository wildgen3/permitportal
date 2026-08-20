---
title: Codes are (scheme, vintage, code) triples, never bare codes
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0002
---

# ADR-0002: Codes are `(scheme, vintage, code)` triples

## Context and problem statement

Industry classification schemes are revised on a multi-year cycle. Is a code string
sufficient to identify a concept?

## Decision drivers

- **Codes have been reused for different concepts across revisions.** At least seven are
  known: 325199, 332995, 336310, 336998, 4881, 48811, 488119. A bare code is therefore
  ambiguous, not merely imprecise.
- Regulatory code lists are frozen at the vintage in force when the rule was adopted,
  while a business's code of record is current. Every applicability test is implicitly a
  translation.
- Historical determinations must remain reproducible after a scheme revision.

## Considered options

1. **Bare code string** — simplest.
2. **`(scheme, code)`** — disambiguates schemes but not revisions.
3. **`(scheme, vintage, code)`** — full identity.

## Decision outcome

Chosen: **`(scheme, vintage, code)`**, as the natural key of `Concept`, with a synthetic
surrogate for foreign keys.

This is cheap on day one and extremely expensive to retrofit — every join, every stored
assignment, and every historical query would need rewriting.

### Consequences

- Good: reused codes cannot silently corrupt a join.
- Good: point-in-time queries are possible without heroics.
- Bad: every code reference is three fields, and every comparison is a translation.
- Enforced by: `unique_keys.natural_key` on `Concept` in `spec/model/core.yaml`, and
  linter L-03 rejecting any code predicate without a `list_vintage`.
