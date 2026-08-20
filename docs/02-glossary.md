---
title: Glossary
status: specified
last_reviewed: 2026-08-19
adr: [0002, 0004, 0013]
---

# Glossary

Ubiquitous language. These terms appear in the canonical model, the rules DSL, and the
ADRs with exactly these meanings. Several pairs below are routinely conflated in
practice, and each conflation produces a specific class of wrong determination.

## Units of analysis

The distinctions in this section are the ones that matter most.

**Business entity** — the legal entity. The *company* scope. Company-wide headcount and
receipts attach here. Size-based exemptions almost always count at this scope.

**Establishment** — a single physical site. Most industry-code keying happens here. An
entity may have many; each has its own jurisdiction path, its own employment count, and
its own classification.

**Process** — an activity unit *within* an establishment. Chemical-process regimes
attach codes and quantity thresholds at this level and forbid inheriting the site's
primary code as a default. A model without this node cannot express those regimes at
all.

**Activity** — a line of business carried on at an establishment. Multiple activities
may be primary simultaneously; see `primacy_basis`.

> **The conflation to avoid:** an exemption keyed to *industry at the establishment* and
> an exemption keyed to *size at the company* are two independent tests at two different
> scopes. Applying either at the wrong scope produces wrong answers in both directions.
> Every predicate declares its scope, and the linter checks it against the registry.

## Classification

**Scheme** — a classification system.

**Vintage** — a dated edition of a scheme. Part of concept identity, never optional.

**Concept** — one code within one scheme version. Identity is the triple
`(scheme, vintage, code)`. Never a bare code: codes have been reused for different
concepts across revisions.

**Crosswalk / correspondence** — a published mapping between two scheme versions. An
addressable, versioned object, not a lookup table.

**Projection** — a published-scheme code derived from the canonical activity model. The
output of classification, not its key.

**Code of record** — the confirmed classification a determination may rely on. A
candidate is not a code of record until a person confirms it.

## Obligations and credentials

These four are distinct and are routinely used interchangeably in the wild.

**Regime** — a regulatory programme that generates obligations.

**Obligation** — a thing a business must do: keep records, report, obtain, notify.

**Determination** — the *result* of evaluating whether an obligation applies to a
specific subject, at a specific point in law time, given specific facts. Reproducible by
replay.

**Credential** — a licence, certification, registration, or permit: an artifact a
business or person holds.

**Requirement** — a reified condition between a credential and what it demands. A table,
not a foreign key, because the edge carries jurisdiction, residency, experience, cost,
and dates. Requirements nest: AND groups contain OR alternative sets contain further
groups.

> An obligation may be *satisfied by* obtaining a credential. They are not the same
> thing, and the credential graph is not the obligation set.

## Evaluation

**Predicate** — a typed leaf test in a rule's expression tree.

**Scope** — the unit of analysis a predicate evaluates against. Declared explicitly.

**Scope hop** — evaluating a predicate at one scope against a fact held at another
(process to establishment, establishment to business). Always explicit, always logged.

**List semantics** — whether a code list is exhaustive (`enumerative_closed`) or
illustrative (`illustrative_open`). Determines whether a miss can produce FALSE.

**Evidence tree** — the predicate tree annotated with truth values and citations. The
determination and its explanation are the same object.

**INDETERMINATE** — a determination that could not be reached because facts are missing.
Carries `missing_attributes[]`, which drives the next question. Distinct from "does not
apply".

## Time

**Valid time / law time** — when a rule is in force. `law_from`, `law_to`.

**Transaction time / knowledge time** — when the system learned something.

**As-of query** — a read pinned to a law date, a knowledge date, or both.

## Provenance

**Provenance record** — what produced a surfaced item: model, prompt version, index
version, rank at selection, alternatives shown, who confirmed and when.

**Alternatives shown** — the candidate set the applicant was offered. Persisted, because
an appeal needs to know what they were offered, not only what they chose.

**Citation** — a structured reference to public authority, with a retrieval date and a
text hash.
