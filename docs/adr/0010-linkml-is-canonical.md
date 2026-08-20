---
title: LinkML is the canonical model language; all other schemas are generated
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0010
---

# ADR-0010: LinkML is the canonical model language

## Context and problem statement

The domain model must be expressed once and consumed as JSON Schema, SQL DDL, Python
types, TypeScript types, and RDF shapes. What is the source?

## Decision drivers

- The credential vocabulary this project adopts is JSON-LD, and the classification
  vocabularies are RDF. **That vocabulary ships no OWL axioms or SHACL shapes enforcing
  its own ordering and boolean semantics**, so those shapes are ours to write.
- The house pattern in adjacent projects is TypeScript-first schema definition with
  JSON Schema emitted from it.
- Hand-maintaining five representations of one model guarantees drift.

## Considered options

1. **TypeScript-first (the house pattern), emitting JSON Schema.**
2. **JSON Schema by hand as the source.**
3. **LinkML as the source, generating everything.**

## Decision outcome

Chosen: **LinkML.**

The house pattern cannot reach RDF, and RDF is not optional here — writing the SHACL
shapes the credential vocabulary omits is a load-bearing requirement, not a nicety.
House ergonomics are preserved by generating the TypeScript validator from the emitted
JSON Schema, so there is still exactly one source of truth.

### Consequences

- Good: one model, seven generated artifacts, no drift.
- Good: SHACL for free, which is the whole reason.
- Bad: a Python toolchain dependency (`uv tool install linkml`) that the rest of the
  frontend stack does not need.
- Bad: the TypeScript generator has no native decimal or datetime, and emits warnings.
  Downstream types for those fields need review.
- Bad: two generators are not byte-reproducible — one embeds a wall-clock timestamp, and
  the SHACL generator varies blank-node ordering between processes. `./do spec` therefore
  runs `scripts/normalize-generated.py`, which pins the timestamp and re-serialises the
  graph canonically. Without it the diff gate would fail on every run and be worthless.
- Enforced by: `./do spec` regenerates, normalises, and CI runs
  `git diff --exit-code spec/generated`.

Generated DDL is the **target**, not a migration tool. Migrations are hand-authored and
numbered; CI asserts the migrated schema equals the generated target.
