# Harness

**Status: not yet implemented.**

Runs the fixture sets, scores against `../baselines/`, accounts for cost, and emits the
report that `docs/eval/results.md` is generated from.

Two modes, and the split is a cost control:

- **Recorded** — replays stored model responses. Calls nothing paid. This is what runs in
  CI on every relevant pull request.
- **Live** — real model calls, refreshes recordings and cost-per-1k. Runs only via
  `workflow_dispatch`. It is the only thing in this repository that can spend money.
