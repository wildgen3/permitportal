---
title: Author rules in YAML, compile to CEL, evaluate with partial activation
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0007
---

# ADR-0007: Author in YAML, compile to CEL

## Context and problem statement

Given a typed predicate tree (ADR-0004) with three-valued semantics (ADR-0006), what
actually evaluates it?

## Decision drivers

- Unknown propagation is mandatory, not an optimization.
- Evaluation must be sandboxed, typed, deterministic, and fast enough to run
  interactively.
- Per-leaf attribution must survive evaluation, because the evidence tree is the
  explanation.

## Considered options

1. **JSON Logic** — small, ubiquitous.
2. **OPA / Rego** — a real policy engine with tooling.
3. **Drools or another RETE engine.**
4. **A bespoke evaluator.**
5. **Author in YAML, compile leaves to CEL, evaluate on cel-go with `PartialActivation`.**

## Decision outcome

Chosen: **option 5.**

**Rego is rejected specifically because negation-as-failure is the bug being engineered
against** — `not p` succeeding on absent data is precisely the failure mode ADR-0006
exists to prevent. JSON Logic is untyped with no unknown semantics and no per-node
provenance. Drools' win is many-rules-against-many-facts throughput, which this workload
does not have, and its cost is a JVM dependency plus an agenda that is hard to render as
a citable explanation. A bespoke evaluator means writing an expression engine from
scratch — unpaid risk, when the semantics are the novel part and the evaluator is not.

### Consequences

- Good: typed, sandboxed, sub-millisecond, with mature partial evaluation.
- Bad: **`services/resolver` must be Go**, because cel-go is the mature implementation
  of partial evaluation and the Python bindings are not. This splits the service layer
  across two languages.
- Enforced today by: linter L-01 and L-03 in `scripts/check-rules.py`, which reject
  negation over an open list and any code predicate without a pinned, resolving list.
- **Will be enforced by** the compiler, which must emit `code_in(...) ? true : unknown`
  for open lists so that list polarity (ADR-0008) survives compilation. No compiler
  exists yet.
