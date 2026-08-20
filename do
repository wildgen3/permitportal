#!/usr/bin/env bash
# PermitPortal task runner.
#
# A shell script rather than a Makefile: `make` is not installed on the development
# machine and installing it requires root, which the rest of this toolchain does not.
# Every dependency here is already present (python3, node, uv-installed linkml).
#
#   ./do spec     regenerate spec/generated/ from the LinkML model
#   ./do check    run the same gates CI runs
#   ./do fmt      format what can be formatted
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$REPO"
export PATH="$HOME/.local/bin:$PATH"

MODEL="spec/model/core.yaml"
GEN="spec/generated"

die() { echo "error: $*" >&2; exit 1; }
have() { command -v "$1" >/dev/null 2>&1 || die "$1 not found on PATH. See AGENTS.md."; }

# LinkML's TypeScript generator emits one WARNING per decimal/datetime slot -- a known
# limitation recorded in ADR-0010, not an error. Filter those and keep everything else.
#
# Deliberately NOT `2> >(grep ...)`. Process substitution runs the filter
# asynchronously: when it exits first the generator takes an EPIPE on stderr and dies
# mid-write, leaving a TRUNCATED stdout file. That made the contracts drift gate flap --
# it emptied core.shacl.ttl on roughly two runs in three and blocked a publish for a
# reason that did not exist. Buffer stderr to a file and filter after the command has
# fully exited; the exit status is then the generator's own, not the filter's.
quiet() {
  local err rc=0
  err="$(mktemp)"
  "$@" 2>"$err" || rc=$?
  grep -v "^WARNING:linkml" "$err" >&2 || true
  rm -f "$err"
  return $rc
}

# Normalisation needs rdflib, which ships with LinkML. In CI that is the same interpreter
# that runs the generators; locally LinkML lives in its own uv tool venv.
resolve_python() {
  if python3 -c "import rdflib" 2>/dev/null; then
    command -v python3
    return
  fi
  local venv="$HOME/.local/share/uv/tools/linkml/bin/python"
  if [[ -x "$venv" ]] && "$venv" -c "import rdflib" 2>/dev/null; then
    echo "$venv"
    return
  fi
  die "no interpreter with rdflib found. Install LinkML: uv tool install linkml"
}

cmd_spec() {
  have gen-json-schema
  mkdir -p "$GEN"
  echo "spec: regenerating from $MODEL"
  quiet gen-json-schema     "$MODEL" > "$GEN/core.schema.json"
  quiet gen-shacl           "$MODEL" > "$GEN/core.shacl.ttl"
  quiet gen-sqlddl          "$MODEL" > "$GEN/core.sql"
  quiet gen-pydantic        "$MODEL" > "$GEN/core_pydantic.py"
  quiet gen-typescript      "$MODEL" > "$GEN/core.d.ts"
  quiet gen-jsonld-context  "$MODEL" > "$GEN/core.context.jsonld"
  quiet gen-erdiagram       "$MODEL" > "$GEN/core.er.mmd"
  "$(resolve_python)" scripts/normalize-generated.py
  echo "spec: $(ls -1 "$GEN" | wc -l) artifacts regenerated"
}

cmd_check() {
  local failed=0
  echo "== docs =="
  node scripts/check-docs.mjs || failed=1
  node scripts/gen-adr-index.mjs --check || failed=1

  echo "== status =="
  node scripts/check-status-table.mjs || failed=1

  echo "== clean-room =="
  CLEAN_ROOM_DENYLIST_FILE="${CLEAN_ROOM_DENYLIST_FILE:-$HOME/.config/permitportal/denylist.txt}" \
    python3 scripts/clean-room-check.py || failed=1

  echo "== spec drift =="
  cmd_spec >/dev/null
  if git diff --quiet -- "$GEN"; then
    echo "spec: generated artifacts are in sync."
  else
    echo "spec: FAILED — generated artifacts are stale. Commit the regenerated files:"
    git diff --stat -- "$GEN"
    failed=1
  fi

  echo "== rules =="
  python3 scripts/check-rules.py || failed=1

  echo "== engine purity =="
  python3 scripts/check-engine-purity.py || failed=1

  if [[ $failed -ne 0 ]]; then
    echo
    echo "check: FAILED"
    exit 1
  fi
  echo
  echo "check: all local gates passed"
}

cmd_fmt() {
  if command -v terraform >/dev/null 2>&1; then
    terraform fmt -recursive infra/ && echo "fmt: terraform formatted"
  fi
}

case "${1:-help}" in
  spec)  cmd_spec ;;
  check) cmd_check ;;
  fmt)   cmd_fmt ;;
  *)     sed -n '3,12p' "$0" | sed 's/^# \{0,1\}//' ;;
esac
