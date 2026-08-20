---
title: Rules DSL
status: specified
last_reviewed: 2026-08-19
adr: [0004, 0006, 0007, 0008, 0017]
---

# Rules DSL

Rules are **data**: effective-dated, cited YAML in `spec/rules/`, validated by
`scripts/check-rules.py`. Adding an obligation never means adding a branch to the
evaluator. If a rule cannot be expressed here, that is a finding about the DSL.

## Why a tree, not an expression string

Citations, effective dates, list vintages, and scope hops attach **per node**. A rule
stored as `employees <= 10 || naics in [...]` has nowhere to put them, and — decisively
— the evidence tree shown to the user *is* the rule tree annotated with truth values. If
the tree is not the representation, the explanation is a separate artifact that drifts
from the decision it purports to explain.

## Evaluation semantics

Kleene three-valued logic (ADR-0006), surfaced as four outcomes:

| Result | Meaning | What the system does |
| --- | --- | --- |
| `TRUE` | Applies | Surface the obligation with its citation |
| `FALSE` | Does not apply | Show why, with the citation |
| `UNKNOWN_MISSING_INPUT` | Not yet evaluable | Ask for `missing_attributes[]` |
| `UNKNOWN_NOT_SELF_DETERMINABLE` | Not the applicant's to decide | Route to the issuing authority |

- `any` returns TRUE if any child is TRUE, regardless of unknown siblings.
- `all` returns FALSE if any child is FALSE, regardless of unknown siblings.
- Otherwise, unknown wins.

**Missing input never becomes false.** In a compliance product a false negative is the
harm, and a rule that silently returns "does not apply" because nobody asked about
chemicals is the worst failure this system can produce.

## Node types

```yaml
all:  [ ... ]     # AND
any:  [ ... ]     # OR
none: [ ... ]     # NOR — flips negation context
not:  { ... }     # NOT — flips negation context
ref:  rule.id     # subtree reference to another rule, never a copy
```

Leaves are predicates. Every leaf declares `scope` and `attribute`; code predicates
additionally declare `scheme`, `list_ref`, `list_vintage`, and `list_semantics`.

## List polarity — the rule that matters most

Every code list declares `list_semantics`:

- **`enumerative_closed`** — exhaustive. A predicate may return TRUE or FALSE.
- **`illustrative_open`** — examples. A predicate may return TRUE or **UNKNOWN, never
  FALSE**. Compiles to `code_in(...) ? true : unknown`.

A code hit is dispositive for inclusion. **A code miss proves nothing** and must not
prune the tree. Only 5 of the 11 industrial-activity categories at 40 CFR 122.26(b)(14)
reference an industry code at all; the rest are narrative, and the codes that do appear
are cited as examples.

For multi-activity establishments, obligations are a **set union**, not a
winner-take-all classification.

## Linter rules

All of these fail the build. Implemented in `scripts/check-rules.py`.

| Rule | Rejects |
| --- | --- |
| **L-01** | Negation over an `illustrative_open` list. Asserts an exclusion the source text does not support. |
| **L-02** | A rule whose predicates are *all* code predicates, without a `sole_key_justification` quoting the text that keys to codes. |
| **L-03** | A list predicate without `list_vintage`, without a `list_ref`, with a `list_ref` that does not resolve, or whose declared semantics disagree with the list's. |
| **L-04** | An attribute not in the registry, or a predicate whose declared `scope` disagrees with the registry's `scope_unit`. |
| **L-05** | A rule without at least one positive and one negative fixture. |
| **L-06** | A rule referencing a list marked `is_complete: false` without declaring `fixture_only: true`. An incomplete closed list cannot support a production FALSE. |

Plus structural integrity: required fields, `cites:` references resolving to declared
citations, credential targets existing, and **no cycles in the credential dependency
graph**.

L-06 deserves a note. A list can be *semantically* closed — the regulation's list really
is exhaustive — while the copy in this repository is a representative subset. Those are
different facts, and conflating them would let a specification fixture produce a
confident FALSE. `is_complete` records the second; `list_semantics` records the first.

## Worked example: two keys at two scopes

The canonical demonstration. OSHA's recordkeeping partial exemption uses two independent
keys at **different granularity**: the industry key is establishment-level, the size key
is company-level. Conflating them produces wrong determinations in both directions, and
it is the most common transcription bug in this domain.

```yaml
rule: osha.1904.partial_exemption
scope: establishment
expr:
  any:
    - predicate: attr_lte
      scope: business                                # <-- COMPANY scope
      attribute: company.peak_employment_prior_cy
      value: 10
      cites: [c1]
    - all:
        - predicate: code_in
          scope: establishment                       # <-- SITE scope
          attribute: establishment.classification
          list_ref: lists/osha_1904_app_a
          list_vintage: "2007"
          list_semantics: enumerative_closed
          translation: { from_code_of_record: true, require_composable: true }
          cites: [c3]
        - predicate: is_false
          scope: establishment
          attribute: establishment.osha_bls_state_written_exemption_override
          cites: [c2]
```

Swap either `scope` value and linter L-04 fails the build, naming the attribute and both
scopes. The full rule is in
[`spec/rules/us-federal/osha-1904-partial-exemption.rule.yaml`](../spec/rules/us-federal/osha-1904-partial-exemption.rule.yaml).

Two further features appear there:

**Overlays.** The federal appendix is a *floor*. A state plan may require records from
establishments the federal rule exempts, so the rule carries
`effect: MAY_NARROW, resolution: REQUIRE_JURISDICTION_RULE` — absent a loaded state
rule, the result is INDETERMINATE, never "exempt".

**Non-waivable obligations.** Severe-injury reporting survives *both* keys. An "exempt"
answer that does not also surface the reporting obligation is a harmful answer, so the
non-waivable obligation attaches to the rule rather than living somewhere a reader has
to know to look.

## Worked example: negation, legally

```yaml
- not:
    predicate: any_of
    scope: establishment
    attribute: establishment.exemption_claims
    list_ref: lists/rcw_19_28_exemptions
    list_vintage: "current"
    list_semantics: enumerative_closed
```

Legal, because the exemption list is closed: the statute enumerates its exemptions
exhaustively, so *not being on the list* is a supportable conclusion. The same construct
over an open list is L-01.

Closed is necessary but not sufficient. The copy of that list here is a representative
subset (`is_complete: false`), so the rule that contains this node must also declare
`fixture_only: true` — L-06. Closed-in-law and complete-in-this-repository are different
facts, and only the second licenses a FALSE a user could rely on.

**L-01, L-03 and L-06 apply to every list predicate**, not only to the industry-code
ones: `code_in`, `code_not_in`, `activity_in`, `any_of`. They originally covered the
first two, which left the negation above unchecked by the very rule it illustrates.

## Vintage translation

A code predicate compares a **current** code of record against a list frozen at the
vintage in force when the rule was adopted. The `translation` block controls this:

```yaml
translation:
  from_code_of_record: true      # translate the confirmed code, not a candidate
  require_composable: true       # refuse close-match chains
  on_close_match: REVIEW         # flag rather than guess
```

Retired codes with successors are declared on the list entry itself, so a rule does not
have to encode revision history.

## Compilation

**Leaves** compile to CEL and evaluate on cel-go with `PartialActivation` (ADR-0007).
An attribute absent from the fact set is registered as an unknown attribute pattern
rather than bound to a zero value, so cel-go returns an unknown and the leaf becomes
UNKNOWN instead of FALSE. That is what carries the semantics above through compilation.

The Kleene combinators stay in the host language rather than compiling into the same
expression, for two reasons that only became concrete once the evaluator was written:
CEL has no three-valued logic of its own, and a rule collapsed into a single expression
has nowhere to hang a per-node citation — which would mean the evidence tree was a
separate artifact from the decision, the thing this document opens by rejecting.

List polarity survives compilation the same way. A membership predicate returns a
tri-state — hit, miss, or *not comparable* — and the polarity of the list decides what a
miss means: FALSE from a closed and complete list, UNKNOWN from an illustrative one,
UNKNOWN from a closed list this repository carries only a subset of. "Not on the list"
and "could not be compared to the list" are different facts, and only the first of them
is ever allowed to become FALSE.

Rego was rejected specifically because negation-as-failure — `not p` succeeding on
absent data — is the exact bug this design exists to prevent. The golden case
`exemption-claims-not-asked` pins that down: under negation-as-failure the `not:` node
above would succeed on an unanswered question and produce a confident obligation.

Implemented in `services/resolver/` (`internal/kleene`, `internal/engine`).

## Provenance

Every predicate may carry `cites: [c1, c2]` referencing citations declared on the rule;
unresolvable references fail the build. Every rule carries a `source_url` and
`effective.law_from`. The evidence tree returned with a determination carries the
citation for every node that contributed to the outcome — which is what makes a
determination arguable rather than merely announced.
