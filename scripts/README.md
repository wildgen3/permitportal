---
title: Scripts
status: implemented
last_reviewed: 2026-08-19
---

# Scripts

Repository tooling. These run in CI and locally, and they work today.

| Script | What it does | Exit codes |
| --- | --- | --- |
| `clean-room-check.py` | Scans the tree for terms that must never appear. Never prints the matched term. | 0 clean, 1 findings, 2 misconfigured |
| `check-status-table.mjs` | Every directory declares an honest status, and the root table agrees. | 0 clean, 1 violations |
| `install-hooks.sh` | Installs the pre-commit hook so the scanner runs before a commit exists. | 0 |
| `spike/` | One-off experiments. Nothing here is imported by anything else. | — |

## On the scanner returning 2

`clean-room-check.py` exits **2** when no terms are configured at all. It refuses to report a pass
without having checked anything — a scanner that silently succeeds when misconfigured is worse than no
scanner, because it manufactures false confidence in exactly the situation where confidence is
unwarranted.
