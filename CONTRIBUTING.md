# Contributing

## The short version

1. Open an issue from a template. Blank issues are disabled — the forms exist because they collect the
   fields the CI gates require.
2. Branch: `docs/*`, `feat/*`, `fix/*`, `chore/*`, or `adr/NNNN-*`.
3. One decision per ADR. One ADR per pull request.
4. Pull request titles follow [Conventional Commits](https://www.conventionalcommits.org/). Individual
   commits are not linted.
5. Squash merge. `main` keeps a linear history.

## Single maintainer, real gates

This is a one-person repository. Required approving reviews are set to **zero**, because a solo owner
cannot approve their own pull request and requiring one would deadlock the repo permanently.

Everything else is enforced. Six status checks block merges, force pushes are blocked, conversations
must be resolved, and there are **no bypass actors** — including the owner. `CODEOWNERS` marks the
paths that would require a domain reviewer on a team; it is documentation of that intent, not an
active gate.

## Before you open a pull request

```bash
# Regenerate everything downstream of the LinkML model
./do spec

# Run the gates locally — same commands CI runs
./do check
```

If `./do spec` produces a diff, commit it. Generated artifacts are committed *and* diff-gated:
CI regenerates them and runs `git diff --exit-code`.

## Adding an obligation rule

Rules are data. Adding one means adding YAML to `spec/rules/`, never adding a branch to the evaluator.

Required on every rule:

- `source_url` and a structured citation to public authority
- `effective` dates (`law_from`, optionally `law_to`)
- every code predicate declaring `list_semantics`, `list_vintage`, and a pinned `list_ref`
- every attribute registered in the attribute registry at the declared `scope`
- at least one positive and one negative fixture

The linter (L-01 … L-05) checks all of these and fails the build. If a rule cannot be expressed in the
DSL, that is a finding about the DSL — open a `spec-change` issue rather than special-casing it.

## Writing an ADR

Copy `docs/adr/template.md`, take the next number (numbers are never reused), and fill in the MADR
sections. Superseding an ADR edits both files: the old one gets `status: superseded by ADR-NNNN` and a
forward link. Nothing is ever deleted.

`docs/adr/README.md` is generated. Do not hand-edit it.

A standing rule: **a finding the source research marks 2–1 may inform design but may never be the sole
justification in an ADR.** Anchor on the primary standard instead.

## The clean-room rule

Read the first section of [`AGENTS.md`](AGENTS.md) before your first contribution. Everything here
must be derivable from public, citable sources. The pull request template has a checklist; the scanner
runs regardless of what you check.

To run the scanner locally, create `~/.config/permitportal/denylist.txt` (untracked, never committed)
and install the pre-commit hook:

```bash
scripts/install-hooks.sh
```
