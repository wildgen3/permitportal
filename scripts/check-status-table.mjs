#!/usr/bin/env node
/**
 * Status discipline check.
 *
 * Enforces AGENTS.md rule 8, "never claim work that didn't happen":
 *
 *   1. Every tracked top-level directory has a README.md.
 *   2. That README has front-matter with `status` from the allowed enum.
 *   3. The status table in the root README.md agrees with every one of them.
 *
 * The repository cannot overstate its own completeness. An empty directory with an
 * honest status is legitimate; an empty directory implied to be finished is a defect.
 *
 * No dependencies — front-matter is parsed directly. Exit 0 clean, 1 on violations.
 */

import { readFileSync, existsSync, readdirSync, statSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const REPO = join(dirname(fileURLToPath(import.meta.url)), "..");
const ALLOWED = ["specified", "scaffolded", "implemented"];
const IGNORED = new Set([".git", ".github", "node_modules", "scripts", "spec"]);

/** Directories that must appear in the root status table. `scripts` and `spec` are
 * handled explicitly below so the ignore list above stays about tooling noise. */
const TRACKED = ["docs", "spec", "packages", "services", "apps", "eval", "data", "infra", "scripts"];

const errors = [];

function parseFrontMatter(path) {
  const text = readFileSync(path, "utf8");
  const match = text.match(/^---\r?\n([\s\S]*?)\r?\n---/);
  if (!match) return null;
  const fields = {};
  for (const line of match[1].split(/\r?\n/)) {
    const kv = line.match(/^([A-Za-z_][A-Za-z0-9_]*):\s*(.*)$/);
    if (kv) fields[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, "");
  }
  return fields;
}

// --- 1 & 2: every tracked directory declares an honest status -----------------
const declared = new Map();
for (const dir of TRACKED) {
  const dirPath = join(REPO, dir);
  if (!existsSync(dirPath) || !statSync(dirPath).isDirectory()) {
    errors.push(`${dir}/ is in the tracked list but does not exist`);
    continue;
  }
  const readme = join(dirPath, "README.md");
  if (!existsSync(readme)) {
    errors.push(`${dir}/README.md is missing`);
    continue;
  }
  const fm = parseFrontMatter(readme);
  if (!fm) {
    errors.push(`${dir}/README.md has no front-matter block`);
    continue;
  }
  if (!fm.status) {
    errors.push(`${dir}/README.md front-matter has no 'status' field`);
    continue;
  }
  if (!ALLOWED.includes(fm.status)) {
    errors.push(`${dir}/README.md status '${fm.status}' is not one of: ${ALLOWED.join(", ")}`);
    continue;
  }
  declared.set(dir, fm.status);
}

// --- 3: the root table agrees -------------------------------------------------
const rootReadme = readFileSync(join(REPO, "README.md"), "utf8");
const tabled = new Map();
// Matches:  | [`docs/`](docs/) | `specified` | ... |
const rowPattern = /^\|\s*\[?`([a-z]+)\/`\]?[^|]*\|\s*`([a-z]+)`\s*\|/gm;
let row;
while ((row = rowPattern.exec(rootReadme)) !== null) {
  tabled.set(row[1], row[2]);
}

for (const [dir, status] of declared) {
  if (!tabled.has(dir)) {
    errors.push(`root README status table has no row for ${dir}/`);
  } else if (tabled.get(dir) !== status) {
    errors.push(
      `status mismatch for ${dir}/: README table says '${tabled.get(dir)}', ` +
        `${dir}/README.md says '${status}'`
    );
  }
}
for (const dir of tabled.keys()) {
  // Tested against TRACKED, not `declared` — a directory whose README failed to parse
  // has already produced a precise error above and must not also produce this one.
  if (!TRACKED.includes(dir)) {
    errors.push(`root README status table lists ${dir}/, which is not a tracked directory`);
  }
}

// --- 4: no undeclared top-level directories -----------------------------------
for (const entry of readdirSync(REPO)) {
  if (IGNORED.has(entry) || entry.startsWith(".")) continue;
  if (!statSync(join(REPO, entry)).isDirectory()) continue;
  if (!TRACKED.includes(entry)) {
    errors.push(
      `top-level directory '${entry}/' is undeclared. Add it to TRACKED and to the ` +
        `root status table, or move its contents into an existing directory (AGENTS.md).`
    );
  }
}

if (errors.length) {
  console.error("status: FAILED\n");
  for (const e of errors) console.error(`  - ${e}`);
  console.error(`\n${errors.length} violation(s). See AGENTS.md rule 8.`);
  process.exit(1);
}
console.log(`status: clean (${declared.size} directories, all statuses agree with the root table).`);
