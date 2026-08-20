---
title: Record architecture decisions using MADR 4.0
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0000
---

# ADR-0000: Record architecture decisions using MADR 4.0

## Context and problem statement

This project's value is largely in its decisions rather than its code — it is a
specification first. Decisions therefore need a durable, reviewable, machine-checkable
format. Which one?

## Decision drivers

- The same front-matter validation stack used for rules and docs should validate ADRs.
- A reviewer should be able to see what was rejected, not only what was chosen.
- One decision per pull request is how this repository manufactures a visible
  engineering-process history out of writing that would happen anyway.

## Considered options

1. **MADR 4.0** — YAML front-matter plus structured "considered options" and "pros and cons".
2. **Nygard-style ADRs** — five prose headings, no metadata.
3. **A decision log in one file** — a single running table.

## Decision outcome

Chosen: **MADR 4.0**.

Nygard's five prose headings give CI nothing to hold onto: no status enum to validate,
no numbering to check, no supersession links to verify. MADR's front-matter is
machine-checkable, and its explicit "considered options" section is precisely the
artifact a senior reviewer reads for.

### Consequences

- Good: `docs/adr/README.md` can be generated and diff-gated; supersession is checkable.
- Bad: more ceremony per ADR than Nygard.
- Enforced by: the `docs` gate validates ADR front-matter and regenerates the index.

A standing rule: **a finding the source research marks 2–1 in adversarial verification
may inform design but may never be the sole justification in an ADR.** Anchor on the
primary standard instead.
