---
title: Packages
status: scaffolded
last_reviewed: 2026-08-19
---

# Packages

| Path | Contents | Status |
| --- | --- | --- |
| `rules/engine/` | The evaluator's contract. The implementation is `services/resolver/`; see that README for why. | `scaffolded` |
| `rules/tests/golden/` | Profile → expected obligation set. The correctness specification. | `scaffolded` |

## The constraint that matters here

`rules/engine` declares **zero dependencies on any model client**, and CI asserts it against the
import graph. This is the mechanical enforcement of ADR-0001: the deterministic engine is the system
of record, and if the AI layer ever becomes load-bearing the build breaks.

The golden tests are the correctness spec, not a regression suite. A rule change that alters a golden
expectation requires the pull request to say why the law changed — or to admit the rule was wrong.
