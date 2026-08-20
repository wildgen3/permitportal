---
title: Data source catalog
status: scaffolded
last_reviewed: 2026-08-19
adr: [0009, 0017]
---

# Data source catalog

> **Status: scaffolded.** Sources are registered in `03-source-register.md`; this
> document adds the per-source **ingestion contract**. **Minimum viable version:** the
> federal text backbone (bulk XML primary, versioner for deltas) with its freshness SLO
> and fallback, since that is the only ingestion path the vertical slice needs.

## Per-source contract template

Each source gets: endpoint shapes, authentication, rate posture and robots policy,
keys and join paths, known defects, freshness SLO, refresh cadence, licence, and the
documented fallback when it is unavailable.

## To specify

- Federal text: bulk XML as the primary ingest path, the versioner for deltas, a polite
  client with a contact user-agent and backoff. **The point-in-time source actively
  blocks automated traffic**, and the retrieval floor of 2017-01-01 is surfaced in the UI
  as a bound rather than silently misrepresented.
- Rulemaking: the dockets API has no regulatory-citation or industry filters and must be
  joined through the Federal Register document number.
- Classification and concordance files, including revision status.
- The long official crosswalk document that is an artifact to ingest, not to reproduce.
- Size standards, as the starting point for the self-owned code-to-regulation table.
- The credential vocabulary registry.
- Pilot-jurisdiction state and municipal sources.

## Already settled

- Regulatory-relevance scoring orders a curation queue and never gates inclusion.
  (ADR-0017)
- Ingest availability is a non-functional requirement with a stated fallback, not an
  assumption.
