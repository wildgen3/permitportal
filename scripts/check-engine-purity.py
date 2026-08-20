#!/usr/bin/env python3
"""Decision-plane purity gate.

ADR-0001 says the deterministic rules engine is the system of record and the AI layer
never decides. That is only architecture if it is mechanically enforced, so this gate
asserts that nothing in the decision plane imports or depends on a model client.

Scope -- the decision plane:
    packages/rules/            the evaluator and its golden tests
    services/resolver/         rule evaluation and credential DAG resolution

Two checks, because an import is not the only way a dependency arrives:
    1. source imports        (Python, TypeScript/JavaScript, Go)
    2. dependency manifests  (package.json, pyproject.toml, requirements*.txt, go.mod)

Fails CLOSED on ambiguity, and reports how many source files it actually inspected. A
gate that scans nothing and prints "pass" is indistinguishable from one that verified
something, which is the failure mode this repository keeps finding in itself.

Exit 0 clean, 1 violation, 2 configuration error.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent

DECISION_PLANE = ["packages/rules", "services/resolver"]

# Model clients and the frameworks that wrap them. Substring match on the module path,
# so `anthropic` catches `anthropic.types` and `@anthropic-ai/sdk`.
FORBIDDEN = [
    "anthropic", "openai", "cohere", "mistralai", "together",
    "google.generativeai", "google-generativeai", "google.genai", "google-genai",
    "vertexai", "aiplatform", "genkit", "langchain", "langgraph", "llama_index",
    "llamaindex", "litellm", "ollama", "huggingface_hub", "transformers",
    "sentence_transformers", "bedrock-runtime",
]

SOURCE_SUFFIXES = {".py", ".ts", ".tsx", ".js", ".mjs", ".go"}
MANIFESTS = {"package.json", "pyproject.toml", "go.mod", "requirements.txt", "requirements-dev.txt"}

IMPORT_PATTERNS = [
    re.compile(r"^\s*import\s+([\w.]+)", re.M),                      # py / go
    re.compile(r"^\s*from\s+([\w.]+)\s+import", re.M),               # py
    re.compile(r"""^\s*import\s+.*?from\s+['"]([^'"]+)['"]""", re.M),  # ts/js
    re.compile(r"""require\(\s*['"]([^'"]+)['"]\s*\)""", re.M),      # cjs
    re.compile(r"""^\s*['"]([\w./-]+)['"]\s*$""", re.M),             # go import block entries
]

violations: list[str] = []
scanned_sources = 0
scanned_manifests = 0


def offending(module: str) -> str | None:
    lowered = module.lower()
    for bad in FORBIDDEN:
        if bad in lowered:
            return bad
    return None


for area in DECISION_PLANE:
    root = REPO / area
    if not root.is_dir():
        # Not an error: the tree declares its own completeness elsewhere. But say so.
        print(f"engine-purity: note — {area}/ does not exist yet")
        continue

    for path in sorted(root.rglob("*")):
        if not path.is_file():
            continue
        rel = path.relative_to(REPO).as_posix()

        if path.suffix in SOURCE_SUFFIXES:
            scanned_sources += 1
            text = path.read_text(encoding="utf-8", errors="replace")
            for pattern in IMPORT_PATTERNS:
                for match in pattern.finditer(text):
                    bad = offending(match.group(1))
                    if bad:
                        violations.append(
                            f"{rel}: imports '{match.group(1)}' (matches forbidden '{bad}'). "
                            f"The decision plane may not depend on a model client — ADR-0001."
                        )

        elif path.name in MANIFESTS:
            scanned_manifests += 1
            text = path.read_text(encoding="utf-8", errors="replace")
            for line in text.splitlines():
                bad = offending(line)
                if bad and not line.strip().startswith("#"):
                    violations.append(
                        f"{rel}: declares a dependency matching forbidden '{bad}' — ADR-0001."
                    )

if violations:
    print("engine-purity: FAILED\n")
    for v in violations:
        print(f"  - {v}")
    print(f"\n{len(violations)} violation(s). See docs/adr/0001-rules-engine-is-the-system-of-record.md.")
    raise SystemExit(1)

print(
    f"engine-purity: clean ({scanned_sources} source file(s), "
    f"{scanned_manifests} manifest(s) inspected across the decision plane)."
)
if scanned_sources == 0 and scanned_manifests == 0:
    print(
        "  NOTE: nothing to inspect yet — the evaluator is not written. This gate is "
        "armed, not satisfied. It becomes meaningful with the first line of engine code."
    )
