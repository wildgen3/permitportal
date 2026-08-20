---
title: Clean room
status: specified
last_reviewed: 2026-08-19
adr: [0014]
---

# Clean room

Everything in this repository is derivable from public, citable sources.

That is a claim about provenance, and claims about provenance are worthless unless they
are checkable. This document describes the mechanisms that make it checkable, and is
honest about what they do and do not catch.

## The rule

> No content in this repository originates from any non-public source. Every normative
> claim carries a citation to public authority. All data is synthetic.

Applies to prose, rules, code lists, fixtures, diagrams, commit messages, issue text,
and pull request descriptions.

## Three mechanisms

### 1. Citation is required, not encouraged

Every rule in `spec/rules/` must carry a `source_url` and structured citations, and
every predicate that cites must reference a citation declared on the rule.
`scripts/check-rules.py` enforces this and fails the build.

Every external statistic in `docs/` must carry a `[SRC-nn]` reference that resolves to a
row in `03-source-register.md`; our own targets carry `[internal]` instead, so a number
about the world is never confusable with a number about us. `scripts/check-docs.mjs`
enforces both.

**That gate is narrow, deliberately, and it is worth being precise about what it does not
do.** "Normative claim" is not detectable by machine in the general case. What is
detectable is the class where overclaiming does real damage — rates and percentages — so
that is what is gated. Prose assertions still require citation, and that part rests on
review. An earlier version of this document claimed the broad version while implementing
nothing at all; that was the wrong kind of error for this document in particular.

Every code list in `spec/lists/` must carry `source_url`, `citation`, `edition`, and
`retrieved_at`.

The effect is stronger than a provenance guarantee. A specification in which every
substantive claim traces to a public source **cannot** be client-derived, because
anything client-derived would have no citation to offer.

### 2. A scanner, whose term list is not in this repository

`scripts/clean-room-check.py` runs on every pull request and weekly on a schedule.

The sensitive terms live in the repository secret `CLEAN_ROOM_DENYLIST`. They are
deliberately **not** committed: a plaintext denylist in a public repository discloses
exactly what it exists to protect, and would be the single most informative file here
for anyone trying to work out what was being kept out.

Findings **never print the matched term**. In CI they do not print the location either.

That second restriction was not in the first version of this document, and its absence
was a real defect. On a **public** repository, Actions logs are world-readable, so a
finding of the form "line 87 of `docs/01-vision-and-problem.md` matched" is itself the
disclosure: a reader opens the public line 87 and recovers the protected term by
inspection. The control was sound before publication and self-defeating after it.

Locations are printed only under `--local`, which the pre-commit hook passes, because at
that point the tree is not yet public.

The scanner normalizes case, accents, and punctuation before matching, so `Foo-Bar`,
`foo bar`, and `FOO_BAR` are one term. Acronyms match as whole words only.

A public generic list ships in `.github/clean-room-generic.txt` covering
confidentiality markers and personal-data markers. It discloses nothing and exists so
that outside readers can see the mechanism working.

**On state names.** An earlier design put all fifty state names on the generic list, on
the theory that a state name outside a survey table is worth a manual look. That was
dropped: this repository's pilot jurisdiction is a US state, so the scanner would fire
on hundreds of legitimate lines. A check that cries wolf on every correct document
trains its reader to skip it, which is worse than not having it. Jurisdiction discipline
is handled instead by ADR-0014 and by the requirement that every jurisdiction-specific
rule cite public law.

**The escape hatch.** A line containing `clean-room-allow` is skipped. It exists for
prose that legitimately discusses a marker phrase — this document, and the generic list
itself. It is deliberately verbose so that it is obvious in a diff, and it is deliberately
**not** available for the sensitive list: a contributor cannot wave through a term they
should not have written.

**The scanner exits 2 when no terms are configured at all.** It refuses to report a pass
without having checked anything. "Nothing was checked" and "nothing was found" produce
identical output otherwise, and the first silently reads as the second in exactly the
situation where that is most dangerous.

### 3. The pilot jurisdiction was chosen by a rubric, not defaulted

The exemplar jurisdiction is US-federal plus one state and one city, selected against the
published source-availability rubric in `18-jurisdiction-onboarding.md`: bulk-retrievable
state and administrative code, a machine-readable licence catalogue, a stable citable
municipal source, and fees published with effective dates. Recorded as ADR-0014; changing
it requires a superseding ADR.

## Local use

The scanner should run before a commit exists, not after it is public:

```bash
scripts/install-hooks.sh          # creates ~/.config/permitgraph/denylist.txt, installs the hook
```

The local denylist is untracked, mode 600, and excluded from the vault sync.

## What these mechanisms do not catch

Stated plainly, because a control whose limits are undocumented gets trusted past them.

- **Paraphrase.** A scanner matches strings. Architecture described in original words
  carries no string to match. The citation requirement is the control that addresses
  this: content that cannot cite a public source does not belong here regardless of how
  it is worded.
- **Structure.** A distinctive system decomposition is not a term. This is why the
  design here is derived from published sources and reasoned from first principles in
  the ADRs, where the reasoning is visible and can be checked against the cited
  material.
- **Judgment.** Knowing which of several public options is the right one is experience,
  and experience is not confidential. That distinction is the entire basis of this
  project: the sources are public, and the synthesis over them is original work.

## Fixtures

All fixture data in `eval/fixtures/` is synthetic and produced by a committed, seeded
generator. There is no ingestion path by which real business data enters this
repository. Costs, hour counts, and fee amounts appearing in specification fixtures are
illustrative and are not represented as current published schedules.
