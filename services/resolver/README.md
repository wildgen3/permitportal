# Resolver

**Status: implemented.** Go.

Evaluates rules against a subject's facts and resolves the credential dependency graph
into a topological order.

```
go test ./...                              # 11 golden cases, property tests, negative controls
go run ./cmd/resolve --list-rules
go run ./cmd/resolve --mode fixture \
  --rule us-wa.electrical_contractor.required \
  --profile testdata/seattle-electrical.yaml \
  --as-of-law 2026-08-19
```

## What is here

| Package | Responsibility |
| --- | --- |
| `internal/kleene` | Three-valued logic. Truth tables written out, not derived. |
| `internal/spec` | Loads `spec/` — rules, code lists, the attribute registry, credentials. |
| `internal/engine` | Compiles leaves to CEL, walks the tree, produces the evidence tree. |
| `internal/credential` | Requirement trees → topological order, with runtime cycle detection. |
| `cmd/resolve` | CLI. A determination is printed as JSON with its evidence tree. |

## Go rather than Python

This service evaluates CEL with `PartialActivation`, and cel-go is the mature
implementation of partial evaluation. That is not a performance preference — unknown
propagation is the mechanism by which an unevaluable rule returns INDETERMINATE
instead of "does not apply" (ADR-0006, ADR-0007).

Only the **leaves** compile to CEL. The Kleene combinators stay in Go because CEL has
no three-valued semantics of its own and no place to hang a per-node citation, and
because collapsing the tree into one expression would destroy the evidence tree — which
is the explanation, not a rendering of it.

Topological ordering and cycle detection are implemented here, because the credential
vocabulary specifies neither (ADR-0012). A cycle in committed rule data fails CI before
it reaches this service; `credential.Resolve` checks again anyway, because it runs
against whatever was loaded rather than against whatever was linted.

## Modes

`production` and `fixture` (ADR-0019). Two invariants — an incomplete code list cannot
support a FALSE, and a code cannot cross vintages without a crosswalk — would otherwise
block the worked examples that demonstrate the engine. Rather than weakening them,
Fixture mode suspends them explicitly, refuses any rule that has not declared
`fixture_only: true`, and stamps `"mode": "fixture"` on every determination it produces.

`data/crosswalks/` is empty, so **production mode currently refuses both rules in the
corpus.** That is the correct behaviour and the test suite asserts it.

## What is not here

No HTTP surface. `spec/api/openapi.yaml` is not written, and a service whose contract
is not specified would be inventing one. The determination shape in
`internal/engine/engine.go` is what that contract will `$ref`.

No model client, ever. `scripts/check-engine-purity.py` asserts it against both the
import graph and `go.mod` on every pull request (ADR-0001).
