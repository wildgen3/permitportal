---
title: Services
status: scaffolded
last_reviewed: 2026-08-20
---

# Services

| Path | Language | Status | Responsibility |
| --- | --- | --- | --- |
| `resolver/` | Go | **implemented** | Rule evaluation and credential DAG resolution |
| `classifier/` | Python 3.14 + uv | not started | Free text → ranked candidate codes with confidence and an explicit abstain |

This directory's own status is `scaffolded`, which understates the resolver and
overstates the classifier. The status vocabulary is one word per top-level directory
(AGENTS.md rule 8), and rounding **down** is the only safe direction: a directory that
claimed `implemented` while half of it does not exist would be the exact defect that
rule was written to prevent.

## Why two languages

`resolver` is Go because it evaluates CEL with `PartialActivation`, and cel-go is the
mature implementation of partial evaluation. Unknown-propagation is not an optimization
here — it is the mechanism by which a rule that cannot be evaluated returns
INDETERMINATE instead of "does not apply." Python's CEL bindings do not offer it.

`classifier` is Python because retrieval, embeddings, and evaluation tooling live there.

See ADR-0007.
