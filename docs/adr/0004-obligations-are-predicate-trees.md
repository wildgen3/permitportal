---
title: Obligations are boolean expression trees over typed predicates
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0004
---

# ADR-0004: Obligations are boolean expression trees over typed predicates

## Context and problem statement

How is "does this obligation apply to this business" represented so that it is both
evaluable by a machine and auditable by a person who is not an engineer?

## Decision drivers

- Rules-as-prose is unevaluable. Rules-as-code is untestable by the domain experts who
  must verify it, and makes every legal change a deployment.
- Citations, effective dates, list vintages, and scope hops attach **per node**, not per
  rule. A rule stored as an expression string cannot carry them.
- The explanation shown to a user is the same tree annotated with truth values. If the
  tree is not the representation, the explanation is a separate thing that can drift
  from the decision.

## Considered options

1. **Prose plus a hand-written function per regime.**
2. **A typed predicate tree, authored as data.**
3. **A general-purpose expression string per rule.**

## Decision outcome

Chosen: **a typed predicate tree, authored as data.**

### Consequences

- Good: the evidence tree and the decision are the same object.
- Good: a domain expert can author and review a rule without touching the evaluator.
- Bad: the DSL is a thing to design, document, and version.
- Enforced by: `scripts/check-rules.py` walks the tree and enforces L-01 through L-06.

For multi-activity establishments, obligations are a **set union**, not a
winner-take-all classification.
