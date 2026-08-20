---
title: Specification fixtures evaluate in an explicit fixture mode, never by relaxing an invariant
status: accepted
date: 2026-08-20
deciders: [rome]
adr: 0019
---

# ADR-0019: Fixture mode is explicit and narrow

## Context and problem statement

Two invariants block this repository's own worked examples.

**Completeness.** `lists/osha_1904_app_a` is semantically closed — the regulation's
appendix really is exhaustive — but the copy here is a representative subset, declared
`is_complete: false`. Linter L-06 exists because a subset cannot support a confident
FALSE. Yet the golden case `large-company-nonlisted-industry` expects exactly that
FALSE, and it is the case that demonstrates the two-keys-two-granularities bug.

**Vintage.** `spec/rules/us-federal/` keys to the 2007 NAICS edition the appendix was
adopted against; the fixture subject carries a 2022 code of record. ADR-0002 forbids
comparing them without a crosswalk, and `data/crosswalks/` is empty — nothing has been
retrieved, and fabricating a concordance would violate the clean-room rule outright.

So the engine cannot evaluate its own demonstration rules without doing something the
specification forbids. The question is what.

## Decision drivers

- The invariants are the product. Weakening either to make a demo run inverts the
  argument this repository exists to make.
- A determination must never be mistakable for something it is not.
- The failure mode to avoid is silent: an engine that quietly treats a subset as
  exhaustive produces confident, wrong, unfalsifiable answers.

## Considered options

1. **Relax the invariants** — let a closed-but-incomplete list return FALSE, and compare
   codes across vintages when no crosswalk exists.
2. **Ship crosswalk data and a complete appendix**, then evaluate normally.
3. **Author fixture-only rules against synthetic lists**, leaving the real rules
   unevaluable.
4. **An explicit evaluation mode** that suspends both invariants, refuses any rule that
   has not declared itself a fixture, and stamps itself on the output.

## Decision outcome

Chosen: **option 4.**

`Mode` is `production` or `fixture`, and it is a required constructor argument — there
is no default, because a default is how a fixture escapes.

- A rule reachable in fixture mode must declare `fixture_only: true`. Production mode
  refuses those rules outright, with an error naming L-06, rather than answering.
- The crosswalk is selected by mode. Production gets `StrictCrosswalk`: identity within
  a vintage, **not composable** across one, which propagates as UNKNOWN. Fixture gets
  `FixtureVintageStableCrosswalk`, which asserts something false in general and says so
  in its name and in every evidence node it touches.
- Every determination carries `"mode"`, and a waived completeness check is recorded on
  the evidence node as `list_completeness_waived`.

Option 1 is the defect this repository documents in others. Option 2 is correct and is
what phase 4 does, but it makes retrieving and licensing a full appendix a prerequisite
for the engine existing at all. Option 3 keeps the invariants but moves the
demonstration away from the rules the documentation discusses, so the worked example
would no longer be the thing that runs.

### Consequences

- Good: the invariants are never relaxed, only visibly suspended, and the suspension is
  scoped to rules that have opted in and stamped on the output.
- Good: it converts an inconvenience into a demonstration — `production` mode currently
  refuses **both** rules in the corpus, which is the honest state of a system with no
  crosswalk data and a partial appendix, and the test suite asserts it.
- Bad: two evaluation paths exist, and a reviewer must check which one produced a given
  determination. Mitigated by making mode a required argument and a required field.
- Enforced today by: `TestProductionRefusesAFixtureOnlyRule` and the golden case
  `production-mode-refuses-a-fixture-rule`, plus linter L-06 on the data side.
- Retired when: `data/crosswalks/` carries a real concordance with its provenance row
  and the appendix is reproduced in full. Fixture mode does not disappear then — it
  stops being the only way to run the worked examples.
