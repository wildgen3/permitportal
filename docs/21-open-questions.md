---
title: Open questions
status: scaffolded
last_reviewed: 2026-08-19
adr: [0017]
---

# Open questions

> **Status: scaffolded, and permanently incomplete by design.** This document is a
> feature, not an appendix. The credibility of this specification rests on the same
> property as the research it is built on: it says what it does not know.

## The measurement that would most cleanly size the rules layer

**What fraction of federal obligations is resolvable from an industry code at all?**
Nothing published answers this. It is the single number that would most cleanly size the
rules layer, and its absence is why this specification refuses to extrapolate a
fifty-state or all-regime estimate.

The response is not to guess. Phase 3's exit gate is a *measurement*: publish the
measured fraction for the regimes actually implemented. That turns the largest unknown
into this project's first original contribution rather than a hole in the estimate.

## Also unknown

- **How many states actually run a one-stop portal.** There is no authoritative published
  count, and the boundary is genuinely fuzzy: several states run a unified *filing* page
  without unified *guidance*, which is the capability that matters here. This document
  therefore makes no numeric claim, and neither does anything else in the repository. A
  survey with published inclusion criteria would be a genuine contribution.

- **The error rate of self-assigned and agency-assigned codes in the wild.** Unmeasured.
  This is why there is no "import your existing code and trust it" path, and why sizing a
  reclassification workload is explicitly deferred.
- **Municipal coverage.** No authoritative national dataset exists, and commercial
  aggregator coverage is unestablished. Scope is therefore one metro in depth, and no
  national municipal claim appears anywhere in this repository.
- **Peer comparison (stage 3).** The five-stage product includes cataloging a business
  against its peers, and nothing in this repository specifies it. It needs a peer corpus
  the other four stages do not require, and the obvious sources are either proprietary or
  aggregated to a level that defeats the purpose. Listed here rather than quietly dropped:
  it is a real part of the design that has no credible public data path yet.
- **Real-user findability.** Whether an owner can actually navigate a credential chain is
  an empirical question requiring card sorting and tree testing, not an architectural
  one. Deferred to phase 5, and named rather than assumed.
- **Standard stability.** Two adopted credential-vocabulary terms are marked unstable.
  Pinned, with a scheduled review date.

## Standing rule

A finding that survived adversarial verification only 2–1 may inform design but may never
be the sole justification in an ADR. Anchor on the primary standard instead.
