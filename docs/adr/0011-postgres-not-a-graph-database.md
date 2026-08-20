---
title: Postgres is the system of record; no graph database
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0011
---

# ADR-0011: Postgres, not a graph database

## Context and problem statement

The product's headline output is a dependency graph. Does the storage layer need to be
a graph database?

## Decision drivers

- The graph is small: a credential chain is single-digit to low-double-digit nodes per
  jurisdiction. This is not a graph-scale problem; it is a graph-shaped problem.
- The hard requirements are bitemporal correctness, check constraints, and transactional
  integrity across the entity, classification, and obligation tables.
- One database is one thing to operate, back up, and reason about.

## Considered options

1. **A dedicated graph database.**
2. **Postgres with recursive CTEs for traversal.**
3. **Both — Postgres of record, graph projection for queries.**

## Decision outcome

Chosen: **Postgres**, with recursive CTEs for traversal and topological ordering
computed in the resolver.

Option 3 is the shape that gets adopted under real load, but adopting it before there is
load buys a synchronization problem and a second operational surface in exchange for
nothing measurable.

### Consequences

- Good: the confirmation-state `CHECK` constraint from ADR-0005 is expressible and
  enforced by the database.
- Good: bitemporal queries are ordinary SQL.
- Bad: cycle detection and topological sort are ours to implement — which is required
  regardless, since the credential vocabulary specifies neither.
- Revisit when: a single jurisdiction's credential graph exceeds a few thousand edges,
  or multi-hop path queries become interactive-latency bound.
