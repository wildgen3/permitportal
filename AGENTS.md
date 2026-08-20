# AGENTS.md — PermitPortal

Conventions and hard constraints for humans and AI agents working in this repo. Read this before
writing anything. Every rule here exists because something specific went wrong — either in the
adversarial research pass that produced this architecture, or in the published systems this project
studies.

## What this project is

A reference architecture for **state business one-stop portals**: the class of government system that
is supposed to tell a business owner what they must do to operate legally, and in what order.

The five stages: a business owner describes their business in free text → the description is
classified into a formal industry taxonomy → the business is cataloged and compared against peers →
applicable regulations are surfaced → the dependency chain of credentials required to satisfy them is
resolved into an ordered graph.

**Stage 3 is out of scope for v1** and nothing in `spec/` implements it. It stays in the statement of
the product because it is part of the design, not because it is built — see rule 8 below.

## What this project deliberately is not

- **Not legal advice.** Every determination is informational, carries a citation, and instructs the
  reader to confirm against the current text of the rule for their jurisdiction.
- **Not an auto-filer.** v1 never submits anything to an agency on a user's behalf.
- **Not a national municipal dataset.** Municipal law is modeled as a layer and implemented to depth
  in exactly one metro. No national municipal coverage is claimed anywhere.
- **Not a classifier benchmark.** Classification accuracy is a means; obligation-set correctness is
  the product.

---

## Hard constraints

### 1. Clean room — this rule outranks every other rule here

Everything in this repository must be derivable from **public, citable sources**. Not from
recollection, not from a prior employer's or client's systems, not from any non-public document.

- Every external statistic in `docs/` carries a `[SRC-nn]` reference into
  `docs/03-source-register.md`, and every `[SRC-nn]` must resolve to a row there. Our own targets and
  gates carry `[internal]` instead, so a number about the world and a number about us are never
  confusable. Both are enforced by `scripts/check-docs.mjs`.
  **The gate is deliberately narrow**: "normative claim" is not machine-detectable in general, so it
  covers the class where overclaiming actually does damage — rates and percentages. Prose claims still
  need citations; that part rests on review, and this file does not pretend otherwise.
- Every rule entry in `spec/rules/` requires a `source_url` pointing at public authority. A rule
  without one cannot even be filed as an issue — the issue form rejects it.
- All fixture data is synthetic, produced by the committed seeded generator in `eval/fixtures/`.
  There is no path by which real business data enters this repository.
- No client, employer, agency, program, vendor, or product name appears anywhere.
  `scripts/clean-room-check.py` runs on every pull request and on a weekly schedule.
- The exemplar jurisdiction was chosen against the published source-availability rubric in
  `docs/18-jurisdiction-onboarding.md`. See ADR-0014. Do not change it without a superseding ADR.

**If you are unsure whether something is client-derived, it is.** Find a public source that supports
it, cite that, and write it from the source. If no public source supports it, it does not go in.

The scanner's terms live in the repository secret `CLEAN_ROOM_DENYLIST`, never in a file. A plaintext
denylist committed to a public repo discloses exactly what it exists to protect. Locally the
pre-commit hook reads an untracked file outside the repository.

Findings never print the matched term, and in CI they do not print the location either — on a public
repository, "line 87 of this file matched" plus a public line 87 recovers the term by inspection.
Locations appear only under `--local`, where the tree is not yet public.

### 2. The AI layer never decides

The deterministic rules engine is the system of record for every determination. Models generate
candidates and explanations. That boundary is architectural, not a matter of care:

- `packages/rules/engine` declares **zero dependencies on any model client**, and
  `scripts/check-engine-purity.py` asserts it in CI over both source imports and dependency
  manifests. If the AI layer ever becomes load-bearing, the build breaks. The evaluator is not
  written yet, so the gate is currently **armed rather than satisfied** — it says so when it runs.
- A database `CHECK` constraint prevents a `Determination` from referencing a
  `ClassificationAssignment` whose `confirmation_state` is `unconfirmed`. An unconfirmed code cannot
  reach a compliance determination even if application code tries.

Permitted model surfaces: structured extraction from free text, candidate generation and reranking,
question ordering and phrasing, explanation rendering constrained to nodes present in the evidence
tree, rule drafting into a human curation queue, plain-language rewriting, and form pre-fill into an
editable unconfirmed state.

### 3. Never generate a citation, code, threshold, fee, deadline, or identifier with a model

These are retrieved and passed by reference, always. A model that emits `29 CFR 1904.1(b)(1)` from
memory will eventually emit a subsection that does not exist, and it will look exactly as confident as
the correct one. The published record already contains this failure: EPA's own web page cites RMP
applicability at 40 CFR 68.10(d) after the provision was redesignated to 68.10(l) by 82 FR 4696,
84 FR 69913, and 89 FR 17686. If the government's own page is stale, a language model's recollection
is worthless.

This generalizes the predecessor rule "never generate a primary key with a model" to every
externally-owned string.

### 4. A code hit proves inclusion; a code miss proves nothing

The single most consequential rule in the domain, and it is enforced by the type system rather than by
review. Every code predicate declares `list_semantics`:

- `enumerative_closed` — the list is exhaustive. May return TRUE or FALSE.
- `illustrative_open` — the list is examples. May return TRUE or **UNKNOWN, never FALSE**.

Only 5 of the 11 industrial-activity categories at 40 CFR 122.26(b)(14) reference an industry code at
all. Negation over an open list is linter error L-01 and fails the build.

### 5. Codes are `(scheme, vintage, code)` triples — never bare codes

NAICS 2017 `454110` is not NAICS 2022 `454110`. Seven codes were **reused for different concepts**
across revisions (325199, 332995, 336310, 336998, 4881, 48811, 488119), so a bare code silently
corrupts every join and every historical query. Cheap on day one, extremely expensive to retrofit.

Crosswalks are addressable, versioned objects. Close-match chains are **never** auto-composed —
`is_composable` is true only when every hop is an exactMatch.

### 6. Scope is declared, and the linter checks it

Every predicate declares a `scope` from {business, establishment, process, activity}, and every
attribute is registered in the Attribute Registry with its own `scope_unit`. A mismatch is linter
error L-04 and fails the build.

This exists because of a real, published trap: OSHA's recordkeeping partial exemption uses two
independent keys at different granularity. The industry key is **establishment**-level; the size key
is **company**-level (29 CFR 1904.1(b)(1)). Conflating them produces wrong determinations in both
directions. With scope declared and checked, that bug is unwriteable.

### 7. Missing input is UNKNOWN, never false

Three-valued Kleene logic, surfaced as four results: `TRUE`, `FALSE`, `UNKNOWN_MISSING_INPUT` (ask a
question), `UNKNOWN_NOT_SELF_DETERMINABLE` (route to the issuing authority).

A rule that cannot be evaluated returns INDETERMINATE with a `missing_attributes[]` list that drives
the next question. It never returns "does not apply." In a compliance product a false negative is the
harm — a missed obligation becomes a fine.

This is why the DSL compiles to CEL and evaluates with `PartialActivation`, and why OPA/Rego was
rejected: negation-as-failure succeeding on absent data is precisely the bug being engineered against.

### 8. Never claim work that didn't happen

Every top-level directory has a `README.md` with front-matter `status: specified | scaffolded |
implemented`. `scripts/check-status-table.mjs` fails CI if a directory lacks a README, carries an
invalid status, or disagrees with the root README's status table.

The repository cannot overstate its own completeness. An empty directory with an honest status is
legitimate; an empty directory implied to be finished is a defect.

### 9. Generated artifacts are committed and diff-gated

`spec/generated/`, `docs/adr/README.md`, and `docs/eval/results.md` are generated. They are committed
so that GitHub renders them, and CI regenerates them and runs `git diff --exit-code`. **Regenerate;
never hand-edit.** A hand-edit that CI then overwrites is worse than no generation at all.

`spec/model/*.yaml` (LinkML) is the single source of truth. JSON Schema, SHACL, Postgres DDL,
Pydantic, TypeScript, and Zod are all downstream of it.

### 10. Rules are data, not code

Adding an obligation means adding YAML with a citation to `spec/rules/`. It never means adding a
branch to the engine. If a rule cannot be expressed in the DSL, that is a finding about the DSL —
open a `spec-change` issue, do not special-case it in the evaluator.

### 11. Pin versions; verify against the registry, not against blog posts

Agents read blog posts; blog posts are wrong. Check npm and PyPI directly.

- LinkML **1.11.x**. `gen-shacl` and `gen-sqlddl` are the load-bearing generators.
- Go for `services/resolver` — cel-go is the mature partial-evaluation implementation. Python's CEL
  bindings are not.
- Python for `services/classifier`, managed with **uv**, not poetry or pipx.
- **npm**, not pnpm/yarn/bun. **podman**, not docker — images build in Cloud Build, not locally.

---

## Repository layout

| Path | Contents |
| --- | --- |
| `docs/` | Prose: decided, not yet executable. Numbered, front-matter validated. |
| `docs/adr/` | Architecture decision records. MADR 4.0. One ADR per pull request. |
| `spec/` | **Executable specification.** The source of truth. |
| `spec/model/` | LinkML canonical model. Everything else is generated from it. |
| `spec/rules/` | Obligation rules as validated, effective-dated, cited data. |
| `spec/lists/` | Versioned code lists, pinned to a named edition. |
| `spec/generated/` | Generated + committed + diff-gated. Never hand-edit. |
| `packages/rules/` | The evaluator and the golden correctness tests. |
| `services/classifier/` | Python: free text → ranked candidates. |
| `services/resolver/` | Go: rules evaluation + credential DAG resolution. |
| `apps/web/` | Reference UI. |
| `eval/` | Harness, synthetic fixtures, labeled sets, committed baselines. |
| `data/` | Public reference data only, with a mandatory provenance table. |
| `infra/` | Terraform. **Every** cloud resource. |
| `scripts/` | Repository tooling and one-off spikes. |

Code goes in exactly one of these. If you are unsure where something belongs, ask rather than
inventing a new top-level directory.

Two deviations from the monorepo convention I use elsewhere, both deliberate:

- **`eval/` is top-level** rather than nested under a service, because there are two evaluable subjects
  (classification accuracy, obligation-set correctness) plus a cross-cutting cost metric.
- **`data/` is top-level**, because the code lists are the substrate of the whole system and a
  dedicated directory with a mandatory provenance table makes "no client data" checkable.

---

## CI gates

These **fail** the build. They do not warn.

Running today:

| Gate | Enforces |
| --- | --- |
| `docs` | Front-matter schema, ADR back-references resolve, internal links resolve, every external statistic carries a resolving `[SRC-nn]`, referenced `scripts/` paths exist, Mermaid blocks are structurally sane, ADR index in sync, `terraform fmt` |
| `contracts` | LinkML regenerates byte-identically and the committed artifacts match |
| `rules` | Every ruleset validates; no cycles in the credential graph; every credential declared; every predicate registered at the declared scope; every entry cited and effective-dated; linter L-01…L-06 |
| `engine-purity` | Nothing in the decision plane imports or declares a model client |
| `clean-room` | Denylist scan over the full tree |
| `status` | Every directory has a README with a valid status matching the root table |

Arriving with the code they check, and **not running today**:

| Gate | Waiting on |
| --- | --- |
| TypeScript typecheck, OpenAPI lint | `packages/` code and `spec/api/openapi.yaml` |
| `eval` | The harness and a real baseline |
| Determinism property test | The evaluator |
| Migration-equals-target assertion | The first migration |

The distinction is deliberate. Listing a gate that does not run as though it does is the
specific defect logged twice below, and it is the one this repository can least afford.

`eval-live` is the only workflow that can spend money, and it runs on `workflow_dispatch` only.

---

## Working environment

Operator-local configuration — absolute paths, cloud project, mirror targets — lives in an untracked
`AGENTS.local.md`, not here. A public constraints file is not the place for one machine's layout.

The one rule that is general: **`budget.tf` ships in the same pull request as the first billable
resource.** Not the next one.

---

## What failed

A running log of failure modes actually encountered on this repository. Add to it; do not prune it.
This section is the point of the file.

- **`make` is not installed and cannot be installed.** The first task runner here was a
  Makefile, on the reasonable assumption that a Linux development machine has `make`. This one does
  not (`rpm -q make` → not installed), and installing it requires root, which nothing else in this
  toolchain does. Replaced with `./do`, a POSIX shell task runner depending only on tools that are
  actually present. **Rule: verify a tool exists before building a workflow on it.** `command -v` is
  one line and would have caught this before the file was written.
- **The clean-room scanner's own privacy guarantee was void on a public repository.** It
  printed `file:line` for every finding while withholding the term — which is correct
  pre-publication and self-defeating after it, because Actions logs on a public repo are
  world-readable and the location is a pointer to a public line. Anyone reading a failing
  run could open that line and recover the protected term by inspection. Locations are now
  printed only under `--local`, which the pre-commit hook passes. **Rule: a control's
  threat model changes the moment the artifact becomes public. Re-derive it for the
  published state, not the state you developed in.**
- **This file and the README both claimed a CI gate that did not exist.** Both asserted
  that "uncited normative claims fail CI." No such check was implemented — and the
  README's own opening sentence was an uncited statistical claim, exactly what the missing
  gate would have caught. Now implemented in `scripts/check-docs.mjs`, narrowly and
  honestly: external statistics require a resolving `[SRC-nn]`, our own targets require
  `[internal]`, and the documentation says plainly what the gate does not cover. **Rule:
  in a repository whose entire argument is "controls, not promises", an unimplemented
  control is the most damaging possible defect. Grep for every claim of enforcement and
  verify each one resolves to code.**
- **The corrected scanner blocked its own fix commit, and was right to.** The rename
  updated `scripts/install-hooks.sh`, but the *already-installed* `.git/hooks/pre-commit`
  had the old denylist path baked in at install time. Updating a generator does not update
  what it previously generated. Under the old fail-open behaviour this would have passed
  silently against the wrong list; under the fix it exited 2 and aborted the commit.
  **Rule: after renaming anything a script writes into place, re-run the installer.
  `scripts/install-hooks.sh` is idempotent precisely so this is cheap.**

- **The clean-room scanner had two silent-degradation holes, found by probing it rather
  than reading it.** Pointing `CLEAN_ROOM_DENYLIST_FILE` at a non-existent path made it
  fall back to the default list and report *clean* — a pass against a list nobody asked
  for. And an empty sensitive list was only a warning, which matters because in CI
  `CLEAN_ROOM_DENYLIST: ${{ secrets.X }}` expands to an empty string when the secret is
  missing or deleted: the gate would have degraded to generic-terms-only and stayed green
  forever. Both are now exit 2. **Rule: test a control by feeding it a broken
  configuration, not by reading the code. Both of these read as correct and neither
  behaved correctly, and a control that fails open is indistinguishable from a control
  that works right up until the day it matters.**

- **Three GitHub Action SHAs were written from memory, and all three were wrong.** Setting
  up the Pages workflow, the pins for `configure-pages`, `upload-pages-artifact` and
  `deploy-pages` were confidently wrong in both the SHA *and* the major version — the real
  releases were a full major ahead. The workflow would have failed to resolve the action.
  This is the rule already written at the top of this file ("pin exact versions and verify
  against the registry, not against training data") violated while writing the file that
  contains it. Every pin in both repositories is now verified against the API. **Rule:
  a SHA is never recalled, only looked up. `gh api repos/OWNER/REPO/tags` costs one call.**

- **A second claimed gate did not exist — the one behind the central decision.** Five
  places (this file, ADR-0001, the AI-layer spec, the architecture doc, and the eval gate
  table) stated that CI asserts the rules engine has no model-client dependency. Nothing
  did. This is the enforcement mechanism for ADR-0001, which is the decision every other
  decision here derives from, so it was the single most load-bearing unimplemented claim
  in the repository. Now `scripts/check-engine-purity.py`, checking both source imports
  and dependency manifests across `packages/rules` and `services/resolver`, verified by
  planting both. **Rule: after writing a claim of enforcement, grep for it and follow it
  to executable code. Two of these were found in one review; assume there are more until
  each has been traced.**

- **The repository advertised the taxonomy of what it was hiding.** A setup script's
  template listed the categories of protected term, and the failure log named the *kind*
  of thing that had leaked. Individually harmless; together they convert "nothing to see
  here" into "there is a specific named set of things, and here is the shape of it."
  **Rule: never publish the schema of a secret. Describing what a control protects is a
  disclosure about the thing protected.**

- **Two generated artifacts were not byte-reproducible, so the diff gate would have failed
  on every single CI run.** `core.context.jsonld` embeds a wall-clock timestamp, and
  `gen-shacl` emits blank-node property shapes in an order that varies between processes.
  `PYTHONHASHSEED=0` does **not** fix the second — tested. `scripts/normalize-generated.py`
  now pins the timestamp and re-serialises the graph canonically via rdflib, which relabels
  blank nodes deterministically. **Rule: a "regenerate and diff" gate is only a gate if the
  generator is deterministic. Verify that by generating twice and comparing, before relying
  on the gate.** A check that fails when nothing is wrong is the same defect as a check that
  passes when something is — both teach the reader to ignore it.

- **The scanner caught a denylisted term hard-coded in a five-line shell script**, on its
  first run with a real term list — in a local sync path, not in prose. The same run also
  flagged the scanner's own docstring, which had used a denylisted term as its worked
  example. Both fixed. **This is the finding that justifies the whole control**: the leak
  was not anywhere a reviewer would look carefully.
- **Two of the first denylist entries were wrong, in the opposite direction.** "business
  one stop" was listed as sensitive, but it is the *public name of the problem class* —
  several states publish portals under that exact phrase — so it flagged the README, this
  file, and the vision document for correctly describing what this project is. **Rule:
  before adding a term, ask whether it names something confidential or something public.
  A category name is not a secret, and listing one guarantees false positives on your own
  best documents.**

- **The clean-room scanner's own generic list produced false positives on the first full
  run.** The list carried bare single words — `confidential`, `proprietary` — and this
  model has a `business_confidential` data class, so the scanner flagged eleven lines of
  entirely correct content, including `docs/clean-room.md` itself. This is the exact
  failure mode that document had *already argued against* for state names, reproduced one
  file away. Generic markers are now **phrases only**, plus a `clean-room-allow` line
  directive for prose that legitimately discusses a marker. **Rule: a check that fires on
  correct content trains its reader to skip it, which is worse than not having the check.
  Tune for zero false positives on a known-good tree before trusting a single finding.**

- **YAML parsed `TRUE:` and `FALSE:` as booleans.** The `TruthValue` enum declared permissible values
  `TRUE`, `FALSE`, `UNKNOWN_MISSING_INPUT`, `UNKNOWN_NOT_SELF_DETERMINABLE`. YAML 1.1 coerces the
  first two to booleans, so they silently became `True` and `False` in *every* generated artifact —
  JSON Schema, SHACL, the SQL DDL, the TypeScript types. Nothing errored. Caught only by asserting on
  the generated output rather than on the generator's exit code. **Rule: quote enum values that
  collide with YAML scalars, and assert on generated content, not on exit status.**
