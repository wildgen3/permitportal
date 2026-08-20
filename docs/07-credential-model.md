---
title: Credential model
status: scaffolded
last_reviewed: 2026-08-19
adr: [0012, 0013]
---

# Credential model

> **Status: scaffolded.** The application profile is decided in ADR-0012 and the model is
> implemented in `spec/model/core.yaml`; a worked chain exists in
> `spec/credentials/us-wa-electrical.yaml`. **Minimum viable version:** the adopt /
> extend / refuse table from ADR-0012, plus the resolver semantics below.

## To specify

- The full application profile, term by term: adopted verbatim, extended locally,
  refused, with the reason for each.
- The SHACL shapes we author because the vocabulary ships none — ordering and boolean
  evaluation are ours to implement and validate.
- Resolver semantics: topological ordering, cycle handling (a cycle is a data error that
  fails CI, not a runtime condition), and how OR alternatives are presented to a user
  who must pick a path.
- Industry-varying requirements: the vocabulary cannot scope industry per requirement, so
  these are modeled as separate alternative branches. Export is lossy and must say so.
- Holder state and credential instances, which the vocabulary does not model at all.
- The two adopted terms marked unstable, their pinned versions, and the review date.

## Already settled

- Requirements are reified, recursive, and carry jurisdiction, residency, experience,
  cost, and dates. (ADR-0012)
- Jurisdiction is a profile with inclusion **and** exception lists. (ADR-0013)
- Ordering is computed by us. The vocabulary specifies none.
