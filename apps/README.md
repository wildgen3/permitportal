---
title: Applications
status: scaffolded
last_reviewed: 2026-08-19
---

# Applications

| Path | Contents |
| --- | --- |
| `web/` | Reference UI: intake, ranked candidate confirmation, determination and evidence display |

The UI is not decoration. Two of its screens are load-bearing architecture:

- **Ranked candidate confirmation.** Published benchmarking found a classification tool returned the
  correct code 90.1% of the time while end-to-end success was only 75.5% — the gap is presentation,
  not retrieval. "None of these" is a first-class path, not an error state.
- **The evidence tree.** Every determination renders as the predicate tree that produced it, annotated
  with truth values and citations. A determination a user cannot interrogate is a determination they
  cannot appeal.
