---
title: Non-functional requirements
status: scaffolded
last_reviewed: 2026-08-19
adr: [0009, 0018]
---

# Non-functional requirements

> **Status: scaffolded.** **Minimum viable version:** the determinism requirement and the
> data-classification table, both of which constrain the design today.

## To specify

- **Determinism and replayability.** The strongest non-functional requirement here, and
  the reason no sampled output may enter the determination path.
- **Freshness SLOs per source**, with a stale-list report and an alert when a pinned code
  list passes its review date.
- Latency budgets: determination is synchronous and sub-second; classification retrieval
  has its own budget.
- **Accessibility: WCAG 2.2 AA and Section 508**, plus a plain-language reading-level
  target. This is a government-facing product; accessibility is a requirement, not a
  polish item.
- Availability and degradation posture: what the system does when a source is
  unreachable. It says so, rather than answering from stale data silently.
- **Data classification**, including the `restricted` class covering chemical inventory,
  which carries `llm_egress_allowed: false` and is excluded from third-party analytics.
