# Classifier service

**Status: not yet implemented.** Python 3.14, managed with `uv`.

Free text in; **ranked candidates** out, with confidence and an explicit abstain. It
never assigns a code of record (ADR-0005), and it never emits a citation, threshold, or
identifier (`docs/11-ai-layer.md`).

Retrieval is hybrid lexical plus dense over official titles, definitions, and index
entries — top-50, recall-oriented, then reranked to top-k. The query is enriched with the
whole profile (activity, location, size, equipment), not the bare description.

Measured per sector, because accuracy varies by roughly a factor of two across them. See
`docs/12-eval-and-quality.md`.
