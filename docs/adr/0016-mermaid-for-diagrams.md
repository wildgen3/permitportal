---
title: Diagrams-as-code in Mermaid, parse-checked in CI
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0016
---

# ADR-0016: Mermaid for diagrams

## Context and problem statement

This is a documentation-first repository whose diagrams are a primary artifact. What
draws them?

## Decision drivers

- The diagrams must render where the documents are read: on the repository host, and in
  the author's note-taking application.
- Committed image files drift from their source and nobody notices.
- A build step for diagrams is a build step that will eventually be skipped.

## Considered options

1. **Structurizr DSL** — the correct C4 tool.
2. **PlantUML.**
3. **D2.**
4. **Mermaid.**

## Decision outcome

Chosen: **Mermaid.**

It renders natively both on the repository host and in the note-taking application — no
build step, no committed images, no drift. Structurizr is the more *correct* C4 tool but
requires a Java toolchain, an export step, and a hosted workspace, and produces images
that go stale relative to the DSL. D2 renders in neither target.

### Consequences

- Good: a diagram is always in sync with its source, because it *is* its source.
- Bad: **Mermaid has no C4 model.** C4 discipline is therefore convention, not tooling:
  one diagram per level, a fixed legend, and no diagram exceeding 12 nodes.
- Bad: a *full* parse needs a headless browser and would be the most flake-prone job in
  CI, so it is deliberately not run.
- Enforced by: the `docs` gate extracts every fenced Mermaid block and checks the two
  things that actually break in practice — an unrecognised diagram type, and unbalanced
  brackets in node labels. This is a structural check, not a parse, and claiming
  otherwise would be the same defect this repository has already logged twice. A full
  render check can be added later if structural checking proves insufficient.

Use `flowchart` for context and container levels, `sequenceDiagram` for intake through
determination, `stateDiagram-v2` for credential lifecycle, `erDiagram` for the domain
model (generated from the canonical model — see ADR-0010).
