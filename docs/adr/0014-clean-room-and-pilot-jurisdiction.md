---
title: Clean room — public sources only, synthetic fixtures, and a pilot jurisdiction chosen for non-overlap
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0014
---

# ADR-0014: Clean room, on the record

## Context and problem statement

Work in this problem class is routinely performed under confidentiality obligations, and
this repository is public. Whatever a contributor knows from non-public sources must stay
out of it. How is that boundary maintained in a way that is verifiable rather than merely
asserted?

## Decision drivers

- A promise is not a mechanism. Anything that depends on remembering to be careful will
  eventually fail.
- An artifact fully traceable to public sources is defensible in a way that "trust me,
  this is generalized" is not.
- The problem class is public: many states publish partial one-stop portals, and the
  regulatory sources are federal and state law.

## Considered options

1. **A disclaimer in the README.**
2. **A review checklist on pull requests.**
3. **Mechanical enforcement: citation lint, a scanner, and a pilot jurisdiction chosen
   by a published rubric.**

## Decision outcome

Chosen: **option 3**, with options 1 and 2 retained as supporting layers.

1. **Citation lint.** Every rule requires a `source_url`, and every external statistic in
   `docs/` requires a `[SRC-nn]` that resolves to the source register — both fail the
   build. The statistics gate is narrow on purpose: "normative claim" is not detectable by
   machine, so the gate covers the class where overclaiming does damage, and the
   documentation says so rather than implying more.
2. **A denylist scanner** over the full tree on every pull request and weekly. The term
   list lives in a repository **secret**, never in a file — a plaintext denylist in a
   public repo discloses exactly what it exists to protect. Findings never print the
   matched term, and in CI they do not print the location either: on a public repository
   a location is a pointer to a public line, and reading it recovers the term.
3. **The pilot jurisdiction is chosen by a published rubric** — bulk-retrievable state
   and administrative code, a machine-readable licence catalogue, a stable citable
   municipal source, and fees published with effective dates. See
   `../18-jurisdiction-onboarding.md`. A decision, recorded here, not an accident of
   convenience.

All fixture data is synthetic, produced by a committed seeded generator.

### Consequences

- Good: the boundary is checkable by anyone, including by the author's own CI.
- Bad: the scanner cannot be fully exercised by an outside contributor, because they do
  not have the secret. The generic list is public so that the mechanism is still visible.
- Enforced by: `scripts/clean-room-check.py`, the `clean-room` workflow, and the
  `source_url` requirement in `scripts/check-rules.py`.

Changing the pilot jurisdiction requires a superseding ADR.
