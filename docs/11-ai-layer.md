---
title: AI layer
status: specified
last_reviewed: 2026-08-19
adr: [0001, 0005, 0017]
---

# AI layer

A US business-compliance portal, where being wrong means fines and legal exposure, sits
squarely in the high-impact quadrant where published guidance declines to recommend
generative models as decision-makers. It does not follow that models have no place —
only that the place must be specified precisely and enforced mechanically.

The defensible architecture: **recall-weighted retrieval and model judgment for
candidate generation, with deterministic, citable rules for statutory triggers, stored
provenance on every surfaced obligation, and an independent completeness check.**

## Permitted surfaces

1. **Structured extraction** from free text into proposed facts. Always
   `source: llm_proposed`, never confirmed.
2. **Candidate generation.** Hybrid lexical and dense retrieval over official titles,
   definitions, and index entries; top-50 recall-oriented; model rerank to top-k. The
   query is enriched with the whole profile — activity, location, size, equipment — not
   just the bare description.
3. **Question selection.** The resolver emits `missing_attributes[]`; the model chooses
   order and phrasing only. **The engine owns *what* is asked; the model owns *how*.**
4. **Explanation rendering** of an already-computed evidence tree, constrained to nodes
   and citations present in that tree.
5. **Rule drafting into a human curation queue.** Proposed rules and proposed
   regulation matches go to a human author. Nothing reaches production unsigned.
6. **Plain-language rewriting, translation, reading-level adaptation** of curated text.
7. **Form pre-fill**, always into an editable, unconfirmed state.

## Forbidden

1. Emitting or modifying a determination.
2. Setting `is_code_of_record` or `confirmation_state`.
3. **Generating a citation, code string, threshold, fee, deadline, or identifier.**
   These are retrieved and passed by reference. A model that emits a subsection from
   memory will eventually emit one that does not exist, and it will look exactly as
   confident as a correct one. The published record already contains a case of a
   government agency's own page citing a redesignated provision — if the authoritative
   page is stale, recollection is worthless.
4. Pruning candidates below the confirmation UI. Recall is the protected quantity;
   review cannot catch what was never retrieved.
5. Deciding a jurisdiction.
6. Sitting in the gating path of a legally consequential answer. The deterministic result
   exists first; the model decorates it.
7. Reading any attribute whose registry `data_class` is `restricted` — chemical
   inventories. No model egress, no third-party analytics.

## Enforcement, not discipline

| Rule | Mechanism |
| --- | --- |
| The engine never depends on a model | `packages/rules/engine` declares zero model-client dependencies; CI asserts it against the import graph |
| Unconfirmed codes never reach determinations | Database `CHECK` constraint on `Determination` |
| Restricted attributes never leave | `llm_egress_allowed: false` in the registry; a change to that field is a security change requiring an ADR |
| Explanations cannot invent citations | Faithfulness gate: every citation in generated text must be present in the evidence tree. 100% [internal], not a metric — a gate |

## The flow

```
free text
  → extractor (schema-validated structured output) → Fact[source=llm_proposed]
  → retrieval (top-50, recall-weighted) → rerank    → Candidate[1..k]
  → RANKED PRESENTATION
       plain language, per-sector k,
       "none of these" as a first-class path
  → HUMAN CONFIRMATION
  → ClassificationAssignment{
        confirmation_state: owner_confirmed,
        is_code_of_record: true,
        provenance: {model_id, prompt_version, index_version, rank_at_selection,
                     alternatives_shown[], confirmed_by, confirmed_at, ui_version}}
  → RESOLVER (deterministic, no model dependency) → Determination + evidence tree
  → EXPLAINER (citation-constrained) → prose
```

`alternatives_shown[]` is persisted because competent reviewers disagree on the same
business description. An appeal needs to know what the applicant was **offered**, not
only what they picked.

## "None of these" is designed, not exceptional

In the published benchmark a substantial minority of users found no acceptable candidate.
That cohort is a designed route — into a clarifying question, a broader candidate set,
or a staff referral — not an error state and not a dead end.

## Per-sector policy, not a global threshold

Accuracy varies by roughly a factor of two across sectors. A single confidence cut is
therefore wrong nearly everywhere. Policy is a config table:

| Field | Purpose |
| --- | --- |
| `top_k_shown` | Wider in low-accuracy sectors |
| `min_margin` | Minimum gap between first and second candidate |
| `escalate_below` | Threshold for staff review |
| `require_clarifying_question` | Mandatory in the weakest sectors |
| `allow_silent_assignment` | **`false` everywhere in v1** |

`allow_silent_assignment` exists only so that a future high-confidence, low-consequence
carve-out has somewhere to live. It is not enabled.

## Cost and evaluation

`eval.yml` runs against recorded fixtures and calls no paid model. `eval-live.yml` is
`workflow_dispatch` only. Cost per 1,000 classifications is a tracked metric with a
committed ceiling — see `12-eval-and-quality.md`.
