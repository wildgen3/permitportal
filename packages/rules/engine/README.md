# Rules engine

**Status: contract defined here; implemented in [`services/resolver/`](../../../services/resolver/).**

The implementation lives with the service that runs it rather than in this package,
because the choice of language is forced by the semantics: evaluation needs cel-go's
`PartialActivation`, and the resolver is Go for exactly that reason (ADR-0007). Splitting
the evaluator from the service that is the only consumer of it would buy nothing and
would put the contract and its one implementation in two languages.

What stays here is the contract and the golden correctness cases in
[`../tests/golden/`](../tests/golden/), which are data: a profile in, an expected
determination out. Adding a case never means adding a branch to the evaluator.

## The constraint that defines this package

**Zero dependencies on any model client.** CI asserts it against the import graph and
against `go.mod`. This is the mechanical enforcement of ADR-0001 — if the AI layer ever
becomes load-bearing in a determination, the build breaks rather than the determination
becoming unreproducible.

## Contract

```
evaluate(rule_version, subject_facts, as_of_law) -> Determination
```

Returns one of four outcomes (`docs/06-rules-dsl.md`):

| Outcome | Carries |
| --- | --- |
| `TRUE` | Evidence tree with the citation for every contributing node |
| `FALSE` | Evidence tree showing which predicate excluded it |
| `UNKNOWN_MISSING_INPUT` | `missing_attributes[]` — the next questions to ask |
| `UNKNOWN_NOT_SELF_DETERMINABLE` | The reasons to route to the authority |

The split between the two UNKNOWNs is not cosmetic. Only a missing *attribute* can be
resolved by asking the applicant something; an incomplete code list, an unavailable
crosswalk, or an unloaded jurisdiction rule are gaps in the system's own knowledge, and
offering the subject a form field for one of those is a dead end. Where both appear in
the same determination, route-to-authority wins.

Authored YAML compiles leaf-by-leaf to CEL; evaluation uses `PartialActivation` so
unknown attributes enter the unknown set rather than defaulting to false (ADR-0006,
ADR-0007). The Kleene combinators are in Go, over the tree — which is also the evidence
tree.

Determinations are byte-identically reproducible from
`(rule_version_id, engine_version, as_of_law, input_snapshot_hash)`. A property test
asserts this, and it is the reason nothing sampled may enter this path.
