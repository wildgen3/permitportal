---
title: Documentation
status: specified
last_reviewed: 2026-08-19
---

# Documentation

Prose: decided, but not yet executable. The executable specification lives in [`../spec/`](../spec/).

## Reading order

| # | Document | What it settles |
| --- | --- | --- |
| 01 | [Vision and problem](01-vision-and-problem.md) | The unsiloing thesis, the five-stage flow, non-goals |
| 02 | [Glossary](02-glossary.md) | Ubiquitous language. Read before the ADRs — the terms are load-bearing |
| 03 | [Source register](03-source-register.md) | The public sources, the killed claims, and the "do not cite" list |
| 04 | [Domain model](04-domain-model.md) | Entity graph, scope units, roll-up, the attribute registry |
| 05 | [Classification and schemes](05-classification-and-schemes.md) | Vintages, crosswalks, projections |
| 06 | [Rules DSL](06-rules-dsl.md) | Grammar, predicates, three-valued semantics, the linter |
| 07 | [Credential model](07-credential-model.md) | The CTDL application profile |
| 08 | [Data source catalog](08-data-source-catalog.md) | Ingestion contracts and known defects |
| 09 | [Architecture (C4)](09-architecture-c4.md) | Context and containers |
| 10 | [API contracts](10-api-contracts.md) | Companion to `spec/api/openapi.yaml` |
| 11 | [AI layer](11-ai-layer.md) | Permitted and forbidden surfaces, the confirmation gate |
| 12 | [Eval and quality](12-eval-and-quality.md) | Metrics, baselines, CI gates |
| 13 | [Provenance and audit](13-provenance-and-audit.md) | The evidence tree and determination replay |
| 14 | [Non-functional requirements](14-nfr.md) | Determinism, freshness, accessibility, data classes |
| 15 | [Threat and risk](15-threat-and-risk.md) | STRIDE, legal exposure, classification gaming |
| 16 | [Legacy integration](16-legacy-integration-acl.md) | Anti-corruption layer for agency systems of record |
| 18 | [Jurisdiction onboarding](18-jurisdiction-onboarding.md) | The playbook for the largest line item |
| 19 | [Vertical slice](19-vertical-slice.md) | The first thing to build |
| 20 | [Roadmap and cost](20-roadmap-and-cost.md) | Phases, gates, what money is spent where |
| 21 | [Open questions](21-open-questions.md) | What is not known. A feature, not an appendix |
| 22 | [Fixtures](22-fixtures.md) | Synthetic data: the generator, coverage, adjudication |

Plus [`clean-room.md`](clean-room.md) (the provenance rule, in force) and [`adr/`](adr/) (every
decision).

## Minimum viable version

Documents 01, 02, 03, 04, 06, 11, and 19 are the spine. If the set is ever cut short, those seven
stand alone and remain coherent. Every other document defines a minimum viable version in its own
front-matter so the set degrades gracefully rather than ending half-written.

## Conventions

Every document opens with front-matter validated in CI:

```yaml
---
title: Obligation graph
status: specified          # specified | scaffolded | implemented
last_reviewed: 2026-08-19
supersedes: null
adr: [0004, 0005, 0011]    # the decisions this document implements
---
```

The `adr:` back-reference is checked: a document referencing an ADR number that does not exist fails
the build. Normative claims carry a `[SRC-nn]` citation into [`03-source-register.md`](03-source-register.md).
