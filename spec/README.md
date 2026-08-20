---
title: Executable specification
status: specified
last_reviewed: 2026-08-19
---

# Executable specification

The source of truth. Prose in [`../docs/`](../docs/) explains and justifies what is defined here, but
never overrides it.

| Path | Contents |
| --- | --- |
| `model/` | LinkML canonical model. **Everything else is generated from this.** |
| `rules/` | Obligation rules as validated, cited, effective-dated data |
| `lists/` | Versioned code lists, each pinned to a named edition |
| `api/` | OpenAPI 3.1, `$ref`-ing into generated JSON Schema |
| `generated/` | Generated, committed, and diff-gated. **Never hand-edit.** |

## Regenerating

```bash
./do spec                        # regenerate everything under generated/
git diff --exit-code spec/generated   # CI runs exactly this
```

Generated artifacts are committed so GitHub renders them and so a reviewer can read the JSON Schema
without a toolchain. They are diff-gated so they can never drift from the model.

## Why LinkML rather than a TypeScript-first schema

This model has to speak RDF. CTDL is JSON-LD; SKOS and XKOS are RDF vocabularies. CTDL ships no OWL
axioms or SHACL shapes enforcing its own ordering and boolean semantics, so those shapes are ours to
write — and LinkML emits SHACL. Zod cannot reach that side of the problem.

The ergonomics are preserved by generating Zod from the emitted JSON Schema, so there is still
exactly one source of truth. See ADR-0010.
