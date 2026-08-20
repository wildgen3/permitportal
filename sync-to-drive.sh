#!/usr/bin/env bash
# Mirror this repository to a local note-vault path.
#
# The destination is NOT hard-coded: it comes from $PERMITPORTAL_VAULT. An earlier version
# of this script carried an absolute path, which is how a denylisted term reached a
# committed file. Operator-local paths do not belong in a public repository.
#
#   PERMITPORTAL_VAULT=/path/to/vault ./sync-to-drive.sh
set -euo pipefail
SRC="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/"
DST="${PERMITPORTAL_VAULT:?set PERMITPORTAL_VAULT to the mirror destination}"

# -L: resolve symlinks to real files. A FUSE mount may reject symlinks outright
# (Input/output error), so CLAUDE.md -> AGENTS.md must be flattened.
rsync -aL --delete \
  --exclude '.git/' \
  --exclude 'node_modules/' --exclude '.next/' --exclude '.venv/' \
  --exclude '__pycache__/' --exclude '.terraform/' --exclude '*.tfstate*' \
  --exclude '*.tfvars' \
  --exclude '*.tsbuildinfo' --exclude '.env*' --exclude 'site/' \
  --exclude '*debug.log' --exclude '.pytest_cache/' --exclude '.ruff_cache/' \
  "$SRC" "$DST"
echo "Synced -> $DST"
