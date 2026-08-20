# Rules engine

**Status: interface defined, implementation lands with the vertical slice.**

## The constraint that defines this package

**Zero dependencies on any model client.** CI asserts it against the import graph. This
is the mechanical enforcement of ADR-0001 — if the AI layer ever becomes load-bearing in
a determination, the build breaks rather than the determination becoming unreproducible.

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
| `UNKNOWN_NOT_SELF_DETERMINABLE` | The authority to route to |

Authored YAML compiles to CEL; evaluation uses cel-go's `PartialActivation` so unknown
attributes enter the unknown set rather than defaulting to false (ADR-0006, ADR-0007).

Determinations are byte-identically reproducible from
`(rule_version_id, engine_version, as_of_law, input_snapshot_hash)`. A property test
asserts this, and it is the reason nothing sampled may enter this path.
