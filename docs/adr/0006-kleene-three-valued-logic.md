---
title: Kleene three-valued logic; UNKNOWN is never coerced to false
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0006
---

# ADR-0006: Missing input is UNKNOWN, never false

## Context and problem statement

Intake is incremental. What does a rule return when the facts it tests have not been
collected yet?

## Decision drivers

- In a compliance product, a **false negative is the harm.** A missed obligation becomes
  a fine; a spurious question is an inconvenience.
- Chemical-process rules cannot be evaluated at all without a chemical inventory. If
  absence evaluates to false, every such rule silently returns "does not apply" to every
  business that has not yet been asked.
- The set of unresolved facts is exactly the set of questions worth asking next. Losing
  it means the system cannot drive its own intake.

## Considered options

1. **Two-valued logic, absence is false.**
2. **Two-valued logic, absence is an error.**
3. **Kleene three-valued logic with unknown propagation.**

## Decision outcome

Chosen: **Kleene three-valued logic**, surfaced as four results: `TRUE`, `FALSE`,
`UNKNOWN_MISSING_INPUT` (ask a question), and `UNKNOWN_NOT_SELF_DETERMINABLE` (route to
the issuing authority — some determinations are assignable only by the regulator).

`any` returns TRUE if any child is TRUE regardless of unknown siblings; `all` returns
FALSE if any child is FALSE; otherwise unknown wins.

### Consequences

- Good: an unevaluable rule returns INDETERMINATE with `missing_attributes[]`, which
  drives the next question. It never returns "does not apply".
- Bad: every consumer must handle three outcomes, not two.
- Enforced by: the evaluator uses partial activation (ADR-0007); `TruthValue` in the
  canonical model has no two-valued representation to fall back to.
