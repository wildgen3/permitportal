---
title: Classification and schemes
status: scaffolded
last_reviewed: 2026-08-19
adr: [0002, 0003]
---

# Classification and schemes

> **Status: scaffolded.** The decisions are made and recorded in ADR-0002 and ADR-0003;
> this document expands them into a full specification. **Minimum viable version:** the
> concept-identity rules and the crosswalk composition rule, both of which are already
> enforced by `scripts/check-rules.py` and the canonical model.

## To specify

- Concept scheme registration: what it takes to add a scheme, and what must be true of
  its versions.
- Vintage handling end to end: how a current code of record is compared against a list
  frozen at an older vintage, and where the translation is recorded.
- Crosswalk composition: exact matches compose; close matches do not. The hop path and
  match-type chain are retained, and `is_composable` gates automatic use.
- Reused codes: the seven known cases, and how a query spanning a revision boundary must
  be written.
- Roll-up: how `RollupRule` rows are authored, and the regimes that forbid roll-up
  entirely.
- Code spaces that resemble industry codes but are not, and must never be crosswalked as
  though they were.

## Already settled

- Identity is `(scheme, vintage, code)`. Never a bare code. (ADR-0002)
- Published schemes are derived projections, not the primary key. (ADR-0003)
- Crosswalks are addressable, versioned objects, not lookup tables.
