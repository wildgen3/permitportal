---
title: Services
status: scaffolded
last_reviewed: 2026-08-19
---

# Services

| Path | Language | Responsibility |
| --- | --- | --- |
| `classifier/` | Python 3.14 + uv | Free text → ranked candidate codes with confidence and an explicit abstain |
| `resolver/` | Go | Rule evaluation and credential DAG resolution |

## Why two languages

`resolver` is Go because it evaluates CEL with `PartialActivation`, and cel-go is the mature
implementation of partial evaluation. Unknown-propagation is not an optimization here — it is the
mechanism by which a rule that cannot be evaluated returns INDETERMINATE instead of "does not apply."
Python's CEL bindings do not offer it.

`classifier` is Python because retrieval, embeddings, and evaluation tooling live there.

See ADR-0007.
