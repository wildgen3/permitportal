---
title: Bitemporal storage — valid (law) time and transaction (knowledge) time
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0009
---

# ADR-0009: Bitemporal storage

## Context and problem statement

Law changes, and our knowledge of law changes — independently. A rule may be amended
effective next January, and we may learn about a 2019 amendment tomorrow. How is time
modeled?

## Decision drivers

- Determinations get appealed months later. What was told to the applicant, on the
  basis of what law, must be recoverable exactly.
- The authoritative point-in-time source has a floor: reliable retrieval begins at
  2017-01-01. Earlier dates cannot be honestly served.
- Retroactively rewriting what a business was told is not acceptable.

## Considered options

1. **Current state only.**
2. **Valid time only** (effective dating on rules).
3. **Bitemporal** — valid time and transaction time.

## Decision outcome

Chosen: **bitemporal.** Rules carry `law_from`/`law_to`; the store carries transaction
time; published rulesets are immutable and every read is an as-of query.

### Consequences

- Good: "why did it say that in March" is answerable by replay, not by reasoning.
- Bad: every query carries two dates, and the storage layer is more complex.
- Enforced by: `Determination` requires `rule_version_id`, `engine_version`,
  `as_of_law`, and `input_snapshot_hash`; a property test asserts byte-identical replay.

The 2017-01-01 retrieval floor is surfaced in the UI as a bound, not silently
misrepresented.
