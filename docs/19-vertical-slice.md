---
title: Vertical slice
status: scaffolded
last_reviewed: 2026-08-19
adr: [0005, 0006, 0008, 0012]
---

# Vertical slice

> **Status: scaffolded.** The slice is chosen and its rules and credential chain are
> authored in `spec/`. **Minimum viable version:** this page — it is an implementation
> brief, and the specification it depends on already exists.

## The persona

> *"I do residential electrical work and small remodels out of a shop in Seattle. Eight
> employees. I keep a 55-gallon drum of solvent in the back."*

One persona, three regimes, two jurisdiction layers, one credential chain.

## What it demonstrates

**1. Recordkeeping — two keys at two scopes.** Exempt by the *size* key (8 ≤ 10, company
scope), **not** by the industry key, with severe-injury reporting surfaced anyway. One
screen showing the scope bug avoided, vintage translation from the current code of record
down to the frozen list, and a non-waivable floor.

**2. Chemical-process applicability — at process scope.** Initially INDETERMINATE because
chemical attributes have not been collected. Three targeted questions. Resolves to
not-applicable *with the reason and citation shown*. This is the "a code miss proves
nothing" demonstration and the unknown-propagation demonstration in one screen.

**3. A state trade credential chain.** State business licence → contractor registration
(bond and insurance) → electrical contractor licence → electrician certification (with a
genuine OR: an experience-hours path or an education-plus-reduced-hours path) →
municipal business licence. Topologically ordered, with jurisdiction as a reified node.

## Size

~12 rules, ~40 registered attributes, ~250 synthetic fixtures, a five-node credential
chain, a CLI plus a single-page demo, Postgres in a container, no authentication, no
cloud spend. **Four minutes to demo.** It exercises all five product stages and every
load-bearing architectural claim.

## Gate before locking the jurisdiction

Run the rubric in `18-jurisdiction-onboarding.md` live. Two or more failures means fall
back to another state on the same rubric.
