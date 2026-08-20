---
title: Roadmap and cost
status: scaffolded
last_reviewed: 2026-08-19
adr: [0014]
---

# Roadmap and cost

> **Status: scaffolded.** **Minimum viable version:** the phase table and the cost
> honesty statement.

## Phases and gates

| Phase | Content | Exit gate |
| --- | --- | --- |
| 0 — Source verification | Live probes of every source in the catalog | Every access assumption verified or the catalog amended; pilot jurisdiction locked |
| 1 — Specification | The document set and the executable spec | ADRs accepted; the model generates cleanly; rules lint-clean; fixtures adjudicated with published agreement |
| 2 — Vertical slice | `19-vertical-slice.md` | Demo runs end to end; zero false negatives on the gold set; baselines committed |
| 3 — Multi-regime federal | Ingestion, the self-owned code-to-regulation table, per-sector completeness checklists | Per-regime zero false negatives; **a human signs off the completeness checklist per sector** |
| 4 — Second jurisdiction | Onboard using only the playbook | The playbook, not the concept, is what is tested |
| 5 — Depth and productization | Municipal depth in one metro, accounts, accessibility audit | Legal review of the non-advice posture; measured findability; WCAG 2.2 AA |

## Cost honesty

**Compute is not the constraint.** Through phase 4, total cloud and model spend plausibly
stays in the low hundreds of dollars a month: Cloud Run scales to zero, the eval harness
runs against recordings, and only a manually dispatched workflow can call a paid model.

The dominant cost is **per-jurisdiction curation labour**. That number should be
*measured* in phases 2 and 4 and stated as a range with its basis — not guessed now.

Phase 5 is the genuinely expensive one: counsel, possibly a compliance audit for
public-sector sales, data licensing, and staffed curation.
