#!/usr/bin/env bash
# Install the local pre-commit hook so the clean-room scanner runs before a commit
# exists, rather than after it is already public.
set -euo pipefail
REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DENYLIST="$HOME/.config/permitgraph/denylist.txt"

if [[ ! -f "$DENYLIST" ]]; then
  mkdir -p "$(dirname "$DENYLIST")"
  cat > "$DENYLIST" <<'TEMPLATE'
# Local clean-room denylist. Never committed, never synced.
# One term per line, as a phrase. Lines starting with # are ignored.
TEMPLATE
  chmod 600 "$DENYLIST"
  echo "Created $DENYLIST — add your terms to it."
fi

cat > "$REPO/.git/hooks/pre-commit" <<'HOOK'
#!/usr/bin/env bash
set -euo pipefail
REPO="$(git rev-parse --show-toplevel)"
export CLEAN_ROOM_DENYLIST_FILE="$HOME/.config/permitgraph/denylist.txt"
python3 "$REPO/scripts/clean-room-check.py" --local || {
  echo
  echo "pre-commit: clean-room check failed. Commit aborted."
  echo "Override only if you are certain: git commit --no-verify"
  exit 1
}
HOOK
chmod +x "$REPO/.git/hooks/pre-commit"
echo "Installed pre-commit hook -> $REPO/.git/hooks/pre-commit"
