---
title: Jurisdiction is a reified node with inclusion and exception lists
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0013
---

# ADR-0013: Jurisdiction is a set, not a scalar

## Context and problem statement

Obligations attach at federal, state, county, and municipal levels simultaneously. How
is "where does this apply" represented?

## Decision drivers

- **This is the reason existing portals are siloed.** A state-only model cannot express
  a food truck's city permit, so the city requirement lives in a different system, or
  nowhere.
- "Everywhere except one state" is a common real pattern that a flat state column cannot
  represent.
- Municipal boundaries and incorporation status change, and unincorporated areas need a
  distinguishable representation.

## Considered options

1. **A state column.**
2. **A jurisdiction column plus a level enum.**
3. **A reified jurisdiction node with a parent hierarchy, plus a profile carrying
   inclusion and exception lists.**

## Decision outcome

Chosen: **option 3.** An establishment carries an ordered `jurisdiction_path` from
nation to place, with an `is_unincorporated` flag; credentials and requirements carry a
`JurisdictionProfile` with `main_jurisdiction[]` and `jurisdiction_exception[]`.

### Consequences

- Good: the municipal layer is expressible from day one, even where it is not populated.
- Good: exception patterns are first-class rather than encoded in prose.
- Bad: a two-hop join on every jurisdiction test.
- Scope honesty: municipal law is modeled as a layer and implemented to depth in exactly
  **one metro**. No national municipal coverage is claimed. There is no authoritative
  national municipal dataset, and pretending otherwise would be the most easily falsified
  claim this project could make.
