---
title: Diagram sources
status: scaffolded
last_reviewed: 2026-08-19
adr: [0016]
---

# Diagram sources

Standalone `.mmd` sources for diagrams reused across more than one document. Diagrams
used in exactly one place live inline in that document, because a diagram in a separate
file is a diagram that drifts from its prose.

Mermaid, per ADR-0016 — it renders natively both on the repository host and in the
note-taking application the docs mirror into, so there is no build step and no committed
image to go stale.

Conventions, since Mermaid has no C4 model and the discipline is therefore ours:

- One diagram per C4 level. Context and containers only; there is no level 4, because
  there is no code yet.
- A fixed legend across diagrams at the same level.
- **No diagram exceeds twelve nodes.** A diagram that needs more is two diagrams.

The domain-model ER diagram is not here: it is *generated* from the canonical model to
`spec/generated/core.er.mmd` and diff-gated, so it cannot drift from the schema.
