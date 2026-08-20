# Resolver service

**Status: not yet implemented.** Go.

Evaluates rules against a subject's facts and resolves the credential dependency graph
into a topological order.

**Go rather than Python**, deliberately: this service evaluates CEL with
`PartialActivation`, and cel-go is the mature implementation of partial evaluation. That
is not a performance preference — unknown propagation is the mechanism by which an
unevaluable rule returns INDETERMINATE instead of "does not apply" (ADR-0006, ADR-0007).

Topological ordering and cycle detection are implemented here, because the credential
vocabulary specifies neither (ADR-0012). A cycle in committed rule data fails CI before
it ever reaches this service.
