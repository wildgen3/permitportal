---
title: The canonical classification is an internal activity taxonomy; published schemes are derived projections
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0003
---

# ADR-0003: Published schemes are derived projections, not the primary key

## Context and problem statement

Four published industry schemes are in scope, and none is a superset of the others.
Which one is the compliance engine keyed to?

## Decision drivers

- The central finding of the research this project is built on: *do not make the
  industry code the primary key of the compliance engine; make it a derived,
  vintage-stamped projection over a canonical model of the business, and key obligations
  to the attributes that regulations actually test.*
- Keying on one scheme makes every other scheme lossy and makes each revision a rewrite.
- Regulations mostly do not test industry codes. They test chemical quantities,
  headcount, equipment, jurisdiction, and floor area — with an industry code sometimes
  as one qualifier among several.

## Considered options

1. **Key on the dominant national scheme** — one code per establishment, everything else derived.
2. **Key on an internal activity taxonomy**, projecting outward to every published scheme.
3. **Store all schemes side by side with no canonical form.**

## Decision outcome

Chosen: **an internal activity taxonomy as canonical**, with published schemes as
derived, vintage-stamped projections.

Option 3 has no answer to "what is this business" and pushes reconciliation into every
consumer. Option 1 fails the moment a second scheme or a revision arrives.

### Consequences

- Good: a scheme revision is an added projection, not a migration.
- Good: obligations key to attributes, so a rule survives reclassification.
- Bad: an internal taxonomy is a thing to maintain, and it has no external authority.
- Enforced by: obligations reference the attribute registry; a rule whose predicates are
  *all* code predicates must justify itself explicitly (linter L-02).
