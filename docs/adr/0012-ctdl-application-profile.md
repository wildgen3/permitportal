---
title: CTDL is the credential vocabulary, adopted via a pinned local profile
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0012
---

# ADR-0012: CTDL via a pinned local application profile

## Context and problem statement

Credentials, and the conditions attached to them, need a vocabulary. Invent one, or
adopt the published standard?

## Decision drivers

- An existing standard means interoperability with published credential registries.
- The standard has real gaps for this use case, and some of its terms are marked
  unstable.
- Requirements are genuinely recursive: "four years of experience plus an exam" OR
  "a degree plus two years plus a different exam" is a real licensing path.

## Considered options

1. **A bespoke credential model.**
2. **CTDL adopted wholesale.**
3. **CTDL via a documented local application profile — adopt, extend, and refuse
   explicitly.**

## Decision outcome

Chosen: **option 3.**

Adopted: the credential classes, the condition profile as a reified `Requirement`, the
alternative-condition structure for OR groups, the jurisdiction profile, and renewal.

Extended, with the extension documented: a vintage field alongside the industry-code
field, because the standard types that field as an untyped string with no vintage
validation, so codes from different revisions are indistinguishable. Our own SHACL
shapes, because the standard ships none. A native exclusion model, because the
standard's exclusion term is URI-ranged and cannot express "must not hold X".

Refused: the course-scoped prerequisite term (its domain and range are both courses, so
it cannot express licence-before-licence) and the entry-condition term (its domain
excludes licences).

### Consequences

- Good: interoperable where it matters, honest where it does not fit.
- Bad: export to the standard is lossy for exclusions, and that must be documented
  wherever export is offered.
- Bad: two adopted terms are marked unstable; the profile pins them and carries a
  scheduled review date.
- Enforced by: `spec/model/core.yaml` carries `exact_mappings` to the standard's terms;
  requirement trees are validated by `scripts/check-rules.py`.

Ordering is computed by us, not by the vocabulary — see ADR-0011.
