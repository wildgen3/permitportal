#!/usr/bin/env python3
"""Rules integrity gate.

Validates the executable rule data in spec/. These are domain checks, not schema
checks -- none of them is expressible in JSON Schema, which is why this is a script
rather than an ajv invocation.

Structural integrity:
  - every rule declares scope, effective dates, a source_url, and citations
  - every `cites:` reference resolves to a declared citation on that rule
  - every `list_ref:` resolves to a list that exists
  - the credential dependency graph is acyclic and every target credential exists

Linter rules (each fails the build):
  L-01  negation over an illustrative_open list -- a code miss proves nothing, so
        `not(code_in(open_list))` asserts something the source text does not support
  L-02  a rule whose predicates are ALL code predicates must carry
        sole_key_justification quoting the text that keys to codes
  L-03  every code predicate declares list_vintage and a list_ref pinned to an edition
  L-04  every attribute resolves in the registry AND the predicate's declared scope
        matches the registry's scope_unit for it
  L-05  every rule has at least one positive and one negative fixture
  L-06  a rule referencing a list marked is_complete: false declares fixture_only

Exit 0 clean, 1 on violations, 2 on configuration error.
"""

from __future__ import annotations

import sys
from pathlib import Path

try:
    import yaml
except ImportError:  # pragma: no cover
    print("check-rules: PyYAML is required (python3 -m pip install --user pyyaml)", file=sys.stderr)
    raise SystemExit(2)

REPO = Path(__file__).resolve().parent.parent
SPEC = REPO / "spec"

# Predicates whose truth depends on membership in a versioned list. L-01, L-03 and
# L-06 all turn on list polarity and list completeness, so they must cover every
# predicate that consults a list -- not only the ones that consult an *industry* code
# list. They originally covered only the latter, which left `activity_in` and `any_of`
# entirely unchecked; see the "What failed" log in AGENTS.md.
LIST_PREDICATES = {"code_in", "code_not_in", "activity_in", "any_of"}

# The subset that keys to an industry classification. L-02 is about industry codes
# specifically: a rule that decides applicability on nothing but a NAICS code has to
# quote the text that keys to codes.
CODE_PREDICATES = {"code_in", "code_not_in"}
BOOLEAN_NODES = ("all", "any", "none", "not")

errors: list[str] = []


def err(where: str, code: str, message: str) -> None:
    errors.append(f"{where}: [{code}] {message}")


def load_yaml(path: Path):
    try:
        return yaml.safe_load(path.read_text(encoding="utf-8"))
    except yaml.YAMLError as exc:
        err(path.relative_to(REPO).as_posix(), "PARSE", f"invalid YAML: {exc}")
        return None


# --- load the registry and the lists ----------------------------------------
registry: dict[str, dict] = {}
registry_path = SPEC / "registry" / "attributes.yaml"
if not registry_path.is_file():
    print("check-rules: spec/registry/attributes.yaml is missing", file=sys.stderr)
    raise SystemExit(2)
reg_doc = load_yaml(registry_path) or {}
for entry in reg_doc.get("attributes", []):
    registry[entry["uri"]] = entry

lists: dict[str, dict] = {}
for path in sorted(SPEC.glob("lists/*.yaml")):
    doc = load_yaml(path)
    if not doc:
        continue
    if "id" not in doc:
        err(path.relative_to(REPO).as_posix(), "LIST", "list has no `id`")
        continue
    for field in ("list_semantics", "source_url", "citation", "edition"):
        if field not in doc:
            err(path.relative_to(REPO).as_posix(), "LIST", f"list is missing `{field}`")
    lists[doc["id"]] = doc


# --- walk a rule's expression tree ------------------------------------------
def walk(node, rule_name: str, negated: bool, state: dict) -> None:
    """Recurse the predicate tree. `negated` tracks whether we are underneath a
    `not:` or a `none:`, which is what L-01 turns on."""
    if node is None:
        return
    if isinstance(node, list):
        for child in node:
            walk(child, rule_name, negated, state)
        return
    if not isinstance(node, dict):
        return

    for key in BOOLEAN_NODES:
        if key in node:
            flips = key in ("not", "none")
            walk(node[key], rule_name, negated != flips, state)
            return

    if "ref" in node:  # subtree reference to another rule
        state["refs"].add(node["ref"])
        return

    predicate = node.get("predicate")
    if predicate is None:
        return
    state["predicates"] += 1

    # --- citations resolve
    for cite_id in node.get("cites", []) or []:
        if cite_id not in state["citation_ids"]:
            err(rule_name, "CITE", f"predicate cites `{cite_id}`, which is not declared on the rule")

    # --- L-04: attribute registered, at the declared scope
    attribute = node.get("attribute")
    scope = node.get("scope")
    if attribute:
        entry = registry.get(attribute)
        if entry is None:
            err(rule_name, "L-04", f"attribute `{attribute}` is not in the attribute registry")
        elif scope is None:
            err(rule_name, "L-04", f"predicate on `{attribute}` declares no scope")
        elif entry["scope_unit"] != scope:
            err(
                rule_name,
                "L-04",
                f"scope mismatch on `{attribute}`: predicate says `{scope}`, "
                f"registry says `{entry['scope_unit']}`",
            )

    if predicate not in LIST_PREDICATES:
        state["non_code_predicates"] += 1
        return

    if predicate in CODE_PREDICATES:
        state["code_predicates"] += 1
    else:
        state["non_code_predicates"] += 1

    # --- L-03: vintage and a resolvable, pinned list
    list_ref = node.get("list_ref")
    if not node.get("list_vintage"):
        err(rule_name, "L-03", f"list predicate on `{attribute}` declares no list_vintage")
    if not list_ref:
        err(rule_name, "L-03", f"list predicate on `{attribute}` declares no list_ref")
        return
    referenced = lists.get(list_ref)
    if referenced is None:
        err(rule_name, "L-03", f"list_ref `{list_ref}` does not resolve to any list in spec/lists/")
        return
    state["lists_used"].add(list_ref)

    declared = node.get("list_semantics")
    actual = referenced.get("list_semantics")
    if declared and actual and declared != actual:
        err(
            rule_name,
            "L-03",
            f"predicate declares list_semantics `{declared}` but `{list_ref}` is `{actual}`",
        )

    # --- L-01: negation over an open list
    semantics = actual or declared
    if negated and semantics == "illustrative_open":
        err(
            rule_name,
            "L-01",
            f"negation over `{list_ref}`, which is illustrative_open. A code miss proves "
            f"nothing; this predicate would assert exclusion the source text does not support.",
        )


# --- validate each rule ------------------------------------------------------
rule_files = sorted(SPEC.glob("rules/**/*.rule.yaml"))
if not rule_files:
    print("check-rules: no rule files found under spec/rules/", file=sys.stderr)
    raise SystemExit(2)

for path in rule_files:
    rel = path.relative_to(REPO).as_posix()
    doc = load_yaml(path)
    if not doc:
        continue
    name = doc.get("rule", rel)

    for field in ("rule", "version", "scope", "source_url", "expr"):
        if field not in doc:
            err(name, "STRUCT", f"missing required field `{field}`")
    if "effective" not in doc or "law_from" not in (doc.get("effective") or {}):
        err(name, "STRUCT", "missing `effective.law_from`")

    citations = doc.get("citations") or []
    if not citations:
        err(name, "STRUCT", "rule declares no citations")
    citation_ids = {c.get("id") for c in citations if isinstance(c, dict)}

    state = {
        "citation_ids": citation_ids,
        "predicates": 0,
        "code_predicates": 0,
        "non_code_predicates": 0,
        "lists_used": set(),
        "refs": set(),
    }
    walk(doc.get("expr"), name, False, state)

    # --- L-02: code predicates alone need an explicit justification
    if state["code_predicates"] > 0 and state["non_code_predicates"] == 0:
        if not doc.get("sole_key_justification"):
            err(
                name,
                "L-02",
                "every predicate is a code predicate, but the rule carries no "
                "sole_key_justification quoting the text that keys to codes",
            )

    # --- L-05: fixtures
    fixtures = doc.get("fixtures") or {}
    if not fixtures.get("positive"):
        err(name, "L-05", "no positive fixture")
    if not fixtures.get("negative"):
        err(name, "L-05", "no negative fixture")

    # --- L-06: incomplete lists must be declared as fixtures
    for list_id in state["lists_used"]:
        if lists[list_id].get("is_complete") is False and not doc.get("fixture_only"):
            err(
                name,
                "L-06",
                f"references `{list_id}`, which is marked is_complete: false, but the rule "
                f"does not declare `fixture_only: true`. An incomplete closed list cannot "
                f"support a production FALSE.",
            )


# --- credential graph --------------------------------------------------------
credentials: dict[str, dict] = {}
edges: list[tuple[str, str]] = []

for path in sorted(SPEC.glob("credentials/*.yaml")):
    rel = path.relative_to(REPO).as_posix()
    doc = load_yaml(path)
    if not doc:
        continue
    for cred in doc.get("credentials", []):
        cid = cred.get("id")
        if not cid:
            err(rel, "CRED", "credential has no `id`")
            continue
        if cid in credentials:
            err(rel, "CRED", f"duplicate credential id `{cid}`")
        credentials[cid] = cred
        for field in ("type", "label", "issuing_authority", "source_url", "citation"):
            if field not in cred:
                err(cid, "CRED", f"missing required field `{field}`")

        def visit_requirements(nodes, owner):
            for node in nodes or []:
                if "legal_source" not in node:
                    err(owner, "CRED", f"requirement `{node.get('id')}` has no legal_source")
                target = node.get("target_credential")
                if target:
                    edges.append((owner, target))
                visit_requirements(node.get("children"), owner)

        visit_requirements(cred.get("requirements"), cid)

for owner, target in edges:
    if target not in credentials:
        err(owner, "CRED", f"requires `{target}`, which is not declared")

# Cycle detection over the prerequisite graph. A cycle is a transcription error in
# the source law, not a runtime condition -- it fails the build.
colour: dict[str, int] = {}
adjacency: dict[str, list[str]] = {}
for owner, target in edges:
    adjacency.setdefault(owner, []).append(target)


def find_cycle(node: str, stack: list[str]) -> None:
    colour[node] = 1
    for nxt in adjacency.get(node, []):
        if nxt not in credentials:
            continue
        if colour.get(nxt) == 1:
            loop = " -> ".join(stack[stack.index(nxt):] + [nxt]) if nxt in stack else f"{node} -> {nxt}"
            err("credential-graph", "CYCLE", f"dependency cycle: {loop}")
        elif colour.get(nxt, 0) == 0:
            find_cycle(nxt, stack + [nxt])
    colour[node] = 2


for cid in credentials:
    if colour.get(cid, 0) == 0:
        find_cycle(cid, [cid])


# --- report ------------------------------------------------------------------
if errors:
    print("rules: FAILED\n")
    for e in errors:
        print(f"  - {e}")
    print(f"\n{len(errors)} violation(s). See docs/06-rules-dsl.md and AGENTS.md.")
    raise SystemExit(1)

print(
    f"rules: clean ({len(rule_files)} rules, {len(lists)} lists, "
    f"{len(registry)} registered attributes, {len(credentials)} credentials, "
    f"{len(edges)} dependency edges, no cycles)."
)
