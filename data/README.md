---
title: Reference data
status: scaffolded
last_reviewed: 2026-08-19
---

# Reference data

**Public reference data only.** No business data, real or otherwise, enters this directory —
synthetic fixtures live in [`../eval/fixtures/`](../eval/fixtures/).

## Provenance table

Every file here must have a row. A file without one fails CI.

| File | Source | Publisher | Licence | Retrieved | SHA-256 |
| --- | --- | --- | --- | --- | --- |
| _(none yet)_ | | | | | |

## Why this directory exists at all

The code lists are the substrate of the entire system, and a dedicated directory with a mandatory
provenance table is what makes "no client data" *checkable* rather than remembered. Every crosswalk is
an addressable, versioned object — see ADR-0002 and ADR-0003.
