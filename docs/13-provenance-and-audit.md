---
title: Provenance and audit
status: scaffolded
last_reviewed: 2026-08-19
adr: [0005, 0018]
---

# Provenance and audit

> **Status: scaffolded.** The model is implemented; ADR-0018 settles the approach.
> **Minimum viable version:** the replay contract and what a determination must persist.

## To specify

- The evidence tree as a stored artifact: shape, versioning, and how it is rendered.
- Determination replay: the exact tuple that must be pinned, and the property test that
  asserts byte-identical reproduction.
- Retention: rule versions are never deleted, and why that is a requirement rather than a
  storage decision.
- The appeal path, end to end, from a user's "this is wrong" to a staff view of the
  original evidence tree with the alternatives that were shown.
- Provenance on model-assisted steps: model, prompt version, index version, rank at
  selection.

## Already settled

- Every determination persists `rule_version_id`, `engine_version`, `as_of_law`, and
  `input_snapshot_hash`. (ADR-0018)
- `alternatives_shown[]` is persisted, because an appeal needs what was offered.
  (ADR-0005)
