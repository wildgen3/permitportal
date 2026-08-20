---
title: Code predicates are polarity-typed; negation over an open list is a lint error
status: accepted
date: 2026-08-19
deciders: [rome]
adr: 0008
---

# ADR-0008: A code hit proves inclusion; a code miss proves nothing

## Context and problem statement

Regulatory code lists come in two kinds that look identical in the text: exhaustive
lists, and lists of examples. Treating the second as the first produces confident wrong
exclusions. How is the difference represented?

## Decision drivers

- **Only 5 of the 11 industrial-activity categories at 40 CFR 122.26(b)(14) reference an
  industry code at all** — the remainder are narrative, and the codes that do appear are
  cited as non-exhaustive examples.
- The best published mapping between regulatory text and industry codes posts mean F1
  around 56% with recall around 50% [SRC-14]. Industry relevance is a triage signal, never a gate.
- A wrong exclusion is invisible: the applicant is never told about the obligation they
  were wrongly filtered out of.

## Considered options

1. **Convention and code review** — document that some lists are open.
2. **A boolean flag on the list, checked at runtime.**
3. **A typed property on the list, checked by a linter at author time, and preserved
   through compilation.**

## Decision outcome

Chosen: **option 3.** Every code list declares `list_semantics`:

- `enumerative_closed` — exhaustive. A predicate over it may return TRUE or FALSE.
- `illustrative_open` — examples. A predicate over it may return TRUE or **UNKNOWN,
  never FALSE**.

A NAICS hit is dispositive for inclusion. A NAICS miss proves nothing and must not prune
the tree.

### Consequences

- Good: the most dangerous class of error in the domain becomes unwriteable.
- Bad: authors must classify every list, and getting it wrong in the permissive
  direction produces more questions than necessary.
- Enforced by: linter **L-01** rejects negation over an `illustrative_open` list;
  **L-06** additionally rejects a production FALSE from a closed list that is marked
  `is_complete: false`.
