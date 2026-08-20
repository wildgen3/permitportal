#!/usr/bin/env python3
"""Make generated artifacts byte-reproducible.

Two LinkML generators are non-deterministic, which would make the `contracts` diff gate
fail on every CI run regardless of whether anything changed -- a gate that cries wolf is
a gate people learn to ignore.

  core.context.jsonld   embeds a wall-clock `generation_date`
  core.shacl.ttl        emits blank-node property shapes in an order that varies between
                        processes (PYTHONHASHSEED=0 does NOT fix it -- verified)

The fix for the timestamp is to pin it. The fix for the ordering is to parse the graph
and re-serialise it canonically: rdflib's `to_canonical_graph` relabels blank nodes
deterministically, so two independent generations normalise to identical bytes.

Requires rdflib, which ships with LinkML. Exits 2 rather than silently skipping if it is
absent -- a normaliser that quietly does nothing produces exactly the flapping gate it
exists to prevent.
"""

from __future__ import annotations

import json
import sys
from pathlib import Path

GEN = Path(__file__).resolve().parent.parent / "spec" / "generated"

# Pinned rather than removed: the field documents which generator produced the file, and
# a fixed value keeps the artifact self-describing while making it reproducible.
PINNED_DATE = "1970-01-01T00:00:00"


def normalize_jsonld(path: Path) -> bool:
    doc = json.loads(path.read_text(encoding="utf-8"))
    comments = doc.get("comments")
    if not isinstance(comments, dict) or "generation_date" not in comments:
        return False
    if comments["generation_date"] == PINNED_DATE:
        return False
    comments["generation_date"] = PINNED_DATE
    path.write_text(json.dumps(doc, indent=3, ensure_ascii=False) + "\n", encoding="utf-8")
    return True


def normalize_turtle(path: Path) -> bool:
    try:
        import rdflib
        from rdflib.compare import to_canonical_graph
    except ImportError:
        print(
            "normalize: rdflib is not importable by this interpreter.\n"
            "  It ships with LinkML. Refusing to skip normalisation silently, because an\n"
            "  un-normalised artifact makes the contracts gate fail on every run.",
            file=sys.stderr,
        )
        raise SystemExit(2)

    graph = rdflib.Graph()
    graph.parse(path, format="turtle")
    canonical = to_canonical_graph(graph).serialize(format="longturtle")
    if not canonical.endswith("\n"):
        canonical += "\n"
    if path.read_text(encoding="utf-8") == canonical:
        return False
    path.write_text(canonical, encoding="utf-8")
    return True


def main() -> int:
    changed = []
    jsonld = GEN / "core.context.jsonld"
    if jsonld.is_file() and normalize_jsonld(jsonld):
        changed.append(jsonld.name)

    shacl = GEN / "core.shacl.ttl"
    if shacl.is_file() and normalize_turtle(shacl):
        changed.append(shacl.name)

    print(f"normalize: {len(changed)} artifact(s) normalised" + (f" ({', '.join(changed)})" if changed else ""))
    return 0


if __name__ == "__main__":
    sys.exit(main())
