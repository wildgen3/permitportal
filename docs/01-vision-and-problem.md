---
title: Vision and problem
status: specified
last_reviewed: 2026-08-19
adr: [0001, 0003, 0013]
---

# Vision and problem

## The problem class

A person wants to open a business. Many US states offer some form of "business one-stop"
portal; many others offer a registration page and nothing else. No authoritative count of
either group is published — see `21-open-questions.md`. Even the good ones share a failure
mode:

**Each agency owns its own form. Nobody owns the order.**

An applicant can find the state business licence, the trade licence, the local
tax certificate, and the health permit — each on its own site, each with its own account
and identifier — and still not know that three of them require a fourth to exist first,
or that one of them is issued only after an inspection that takes six weeks to schedule.

The information is public. The *sequence* is not written down anywhere.

## What "unsiloed" means here

Not a single sign-on across agencies, though that helps. Not a shared database, which no
state is going to build. The thing that is actually missing is a **model**: one
representation of a business, against which every agency's rules can be evaluated, that
produces an ordered graph of what that specific business must obtain.

The agencies stay where they are. What gets unified is the reasoning.

## The five stages

1. A business owner describes their business in **free text**.
2. The description is **classified** into formal industry taxonomies.
3. The business is **cataloged and compared** against peers. *(Product scope; deliberately
   out of scope for v1 — it needs a peer corpus the other four stages do not, and nothing
   downstream of it is specified yet. See `21-open-questions.md`.)*
4. Applicable **regulations are correlated and surfaced**.
5. The **dependency chain of credentials** is resolved — what must be obtained, and in
   what order.

## The architectural thesis

> Do not make the industry code the primary key of the compliance engine. Make it a
> derived, vintage-stamped projection over a canonical model of the business, and key
> obligations to the attributes that regulations actually test.

This is not a modeling preference. It follows from what the regulations do:

- Most rules test **attributes**, not industry: chemical quantities, headcount,
  equipment, floor area, jurisdiction. An industry code is one qualifier among several.
- Where rules do cite codes, they are frequently citing **examples**, not an exhaustive
  list. Treating an illustrative list as exhaustive produces confident wrong exclusions
  that nobody ever sees. See ADR-0008.
- Rules key to different **units of analysis** — the company, the site, a process within
  a site — and a model with one unit gets a whole class of determinations wrong in both
  directions. See ADR-0013 and the OSHA example in `06-rules-dsl.md`.

## Non-goals

- **Not legal advice.** Determinations are informational, carry citations, and instruct
  the reader to confirm against current text for their jurisdiction.
- **Not an auto-filer.** v1 submits nothing to any agency.
- **Not a national municipal dataset.** No authoritative one exists. Municipal law is
  modeled as a layer and implemented to depth in exactly one metro.
- **Not a classifier benchmark.** Classification accuracy is a means. Obligation-set
  correctness is the product, and a missed obligation is the harm.

## Why this is hard, honestly

Per-jurisdiction ingestion — starting from state registries and working down to
municipal codes — is **unavoidable scope, not an optimization**. It is the largest and
least glamorous line item in any real build of this, and no amount of retrieval quality
removes it. `18-jurisdiction-onboarding.md` treats it as a first-class problem with a
measured budget rather than as an implementation detail.

The honest position on scale is in `21-open-questions.md`: nobody has published what
fraction of federal obligations is even resolvable from an industry code, and that is
the single number that would most cleanly size the rules layer. This project measures it
for the regimes it implements rather than extrapolating from a guess.
