#!/usr/bin/env node
/**
 * Documentation gate.
 *
 *   1. Every doc has front-matter with title, status (valid enum), and last_reviewed.
 *   2. Every `adr:` back-reference resolves to an ADR that exists. A document that
 *      claims to implement decision 0042 when there is no ADR-0042 is lying about its
 *      own provenance.
 *   3. Every relative markdown link resolves to a file that exists. Dead links in a
 *      specification repository are the fastest credibility leak available.
 *
 * No dependencies. Exit 0 clean, 1 on violations.
 */

import { readFileSync, readdirSync, existsSync, statSync } from "node:fs";
import { join, dirname, resolve, relative } from "node:path";
import { fileURLToPath } from "node:url";

const REPO = join(dirname(fileURLToPath(import.meta.url)), "..");
const VALID_STATUS = new Set(["specified", "scaffolded", "implemented"]);
const SKIP_DIRS = new Set([".git", "node_modules", ".venv", "spec/generated"]);
// Generated or exemplar files: the index is produced by gen-adr-index.mjs and validated
// there; the template carries deliberately-invalid placeholder values.
const SKIP_FILES = new Set(["docs/adr/README.md", "docs/adr/template.md"]);

const errors = [];

function walk(dir, out = []) {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const rel = relative(REPO, full);
    if (SKIP_DIRS.has(entry) || SKIP_DIRS.has(rel) || entry.startsWith(".")) continue;
    if (statSync(full).isDirectory()) walk(full, out);
    else if (entry.endsWith(".md")) out.push(full);
  }
  return out;
}

function frontMatter(text) {
  const m = text.match(/^---\r?\n([\s\S]*?)\r?\n---/);
  if (!m) return null;
  const out = {};
  for (const line of m[1].split(/\r?\n/)) {
    const kv = line.match(/^([A-Za-z_][A-Za-z0-9_]*):\s*(.*)$/);
    if (kv) out[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, "");
  }
  return out;
}

// Which ADR numbers actually exist.
const adrNumbers = new Set();
const adrDir = join(REPO, "docs", "adr");
for (const f of readdirSync(adrDir)) {
  const m = f.match(/^(\d{4})-.*\.md$/);
  if (m) adrNumbers.add(m[1]);
}

const files = walk(REPO);

for (const file of files) {
  const rel = relative(REPO, file);
  if (SKIP_FILES.has(rel)) continue;
  const text = readFileSync(file, "utf8");

  // --- 1 & 2: front-matter, for docs/ only (root-level prose is narrative)
  if (rel.startsWith("docs/") || /^[a-z]+\/README\.md$/.test(rel)) {
    const fm = frontMatter(text);
    if (!fm) {
      errors.push(`${rel}: no front-matter block`);
    } else {
      if (!fm.title) errors.push(`${rel}: front-matter missing 'title'`);
      if (!fm.status) errors.push(`${rel}: front-matter missing 'status'`);
      else if (!VALID_STATUS.has(fm.status) && !/^(proposed|accepted|rejected|deprecated|superseded by ADR-\d{4})$/.test(fm.status))
        errors.push(`${rel}: invalid status '${fm.status}'`);
      if (!fm.last_reviewed && !fm.date) errors.push(`${rel}: front-matter missing 'last_reviewed'`);

      if (fm.adr) {
        const nums = [...fm.adr.matchAll(/\d{4}/g)].map((m) => m[0]);
        for (const n of nums) {
          if (!adrNumbers.has(n)) errors.push(`${rel}: references ADR-${n}, which does not exist`);
        }
      }
    }
  }

  // --- 3: relative links resolve
  const body = text.replace(/^---\r?\n[\s\S]*?\r?\n---/, "");
  for (const m of body.matchAll(/\[[^\]]*\]\(([^)\s]+)\)/g)) {
    let target = m[1];
    if (/^(https?:|mailto:|#)/.test(target)) continue;
    target = target.split("#")[0];
    if (!target) continue;
    const resolved = resolve(dirname(file), target);
    if (!existsSync(resolved)) {
      errors.push(`${rel}: dead link -> ${m[1]}`);
    }
  }
}

// --- 4: citation discipline on statistical claims ---------------------------
// AGENTS.md and README.md both promise that uncited normative claims fail CI. This is
// that gate. It is deliberately narrow -- "normative claim" is not machine-detectable in
// general, but the class where overclaiming actually does damage is external statistics:
// accuracy rates, recall figures, percentages. Every one of those must carry either a
// [SRC-nn] reference into the source register, or an explicit [internal] marker meaning
// "this is our own target or gate, not a claim about the world".
const STAT = /(?<![\w.])(?:\d{1,3}(?:\.\d+)?\s?%|0\.\d{2,3})(?![\w])/;
const CITED = /\[SRC-\d{2}\]|\[internal\]/;

const registerText = readFileSync(join(REPO, "docs", "03-source-register.md"), "utf8");
const knownSrc = new Set([...registerText.matchAll(/SRC-(\d{2})/g)].map((m) => m[1]));

for (const file of files) {
  const rel = relative(REPO, file);
  if (SKIP_FILES.has(rel)) continue;
  if (!(rel.startsWith("docs/") || rel === "README.md")) continue;
  if (rel === "docs/03-source-register.md") continue;

  const text = readFileSync(file, "utf8").replace(/^---\r?\n[\s\S]*?\r?\n---/, "");
  const lines = text.split(/\r?\n/);
  let inFence = false;

  lines.forEach((line, i) => {
    if (/^\s*```/.test(line)) { inFence = !inFence; return; }
    if (inFence) return;
    if (!STAT.test(line)) return;
    if (CITED.test(line)) return;
    // Allow the citation to sit on an adjacent line of the same paragraph.
    const neighbours = [lines[i - 1], lines[i + 1], lines[i + 2]].filter(Boolean).join(" ");
    if (CITED.test(neighbours)) return;
    errors.push(
      `${rel}:${i + 1}: statistical claim without [SRC-nn] or [internal] — ` +
        `"${line.trim().slice(0, 68)}"`
    );
  });

  for (const m of text.matchAll(/\[SRC-(\d{2})\]/g)) {
    if (!knownSrc.has(m[1])) errors.push(`${rel}: cites SRC-${m[1]}, which is not in the source register`);
  }
}

// --- 5: referenced tooling must exist ---------------------------------------
// Cheap guard against the failure this repository has already hit twice: prose that
// names an enforcement mechanism which does not resolve to anything executable.
for (const file of files) {
  const rel = relative(REPO, file);
  if (SKIP_FILES.has(rel)) continue;
  const text = readFileSync(file, "utf8");
  for (const m of text.matchAll(/scripts\/[a-z0-9-]+\.(?:py|mjs|sh)/g)) {
    if (!existsSync(join(REPO, m[0]))) {
      errors.push(`${rel}: references ${m[0]}, which does not exist`);
    }
  }
}

// --- 6: Mermaid blocks are structurally sane --------------------------------
// Not a full parse -- that needs a headless browser and is the flakiest job available.
// This checks what actually breaks in practice: an unrecognised diagram type, and
// unbalanced brackets in node labels. ADR-0016 describes exactly this and no more.
const MERMAID_TYPES = [
  "flowchart", "graph", "sequenceDiagram", "stateDiagram", "stateDiagram-v2",
  "erDiagram", "classDiagram", "journey", "gantt", "pie", "mindmap", "timeline",
];
for (const file of files) {
  const rel = relative(REPO, file);
  if (SKIP_FILES.has(rel)) continue;
  const text = readFileSync(file, "utf8");
  const blocks = [...text.matchAll(/```mermaid\r?\n([\s\S]*?)```/g)];
  for (const [, block] of blocks) {
    const first = block.split(/\r?\n/).map((l) => l.trim()).filter(Boolean)[0] || "";
    const type = first.split(/[\s;]/)[0];
    if (!MERMAID_TYPES.includes(type)) {
      errors.push(`${rel}: mermaid block starts with '${type}', which is not a known diagram type`);
    }
    for (const [open, close] of [["[", "]"], ["(", ")"], ["{", "}"]]) {
      const a = (block.match(new RegExp("\\" + open, "g")) || []).length;
      const b = (block.match(new RegExp("\\" + close, "g")) || []).length;
      if (a !== b) errors.push(`${rel}: mermaid block has unbalanced '${open}${close}' (${a} vs ${b})`);
    }
  }
}

// --- 7: workflow YAML must parse -------------------------------------------
// A broken workflow file is invisible to every other gate here and only surfaces as a
// failed run AFTER the push. This caught a real one: `run: echo "Skipped: ..."` -- a
// colon-space inside an unquoted scalar, which YAML reads as a nested mapping.
const wfDir = join(REPO, ".github", "workflows");
if (existsSync(wfDir)) {
  for (const f of readdirSync(wfDir).filter((n) => /\.ya?ml$/.test(n))) {
    const text = readFileSync(join(wfDir, f), "utf8");
    // Dependency-free structural check: every `key: value` line whose value is
    // unquoted must not contain a further ": " outside quotes.
    text.split(/\r?\n/).forEach((line, i) => {
      const m = line.match(/^(\s*)(-?\s*[A-Za-z_][\w-]*):\s+(.*)$/);
      if (!m) return;
      let value = m[3];
      if (/^['"]/.test(value) || value.startsWith("|") || value.startsWith(">")) return;
      if (/:\s/.test(value)) {
        errors.push(
          `.github/workflows/${f}:${i + 1}: unquoted value contains ": " — YAML will read ` +
            `this as a nested mapping. Quote it. (${value.slice(0, 50)})`
        );
      }
    });
  }
}

// --- 8: every action must be SHA-pinned ------------------------------------
// A version tag is mutable: whoever controls the tag controls what runs with this
// repository's secrets. AGENTS.md states the pinning rule; this enforces it. Caught a
// real one -- hashicorp/setup-terraform@v3 sat unpinned among five pinned actions.
if (existsSync(wfDir)) {
  for (const f of readdirSync(wfDir).filter((n) => /\.ya?ml$/.test(n))) {
    const text = readFileSync(join(wfDir, f), "utf8");
    text.split(/\r?\n/).forEach((line, i) => {
      const m = line.match(/^\s*-?\s*uses:\s*([^\s#]+)/);
      if (!m) return;
      const ref = m[1];
      if (ref.startsWith("./") || ref.startsWith("docker://")) return;
      if (!/@[a-f0-9]{40}$/.test(ref)) {
        errors.push(
          `.github/workflows/${f}:${i + 1}: '${ref}' is not SHA-pinned. A tag is mutable; ` +
            `whoever controls it controls what runs with this repository's secrets.`
        );
      }
    });
  }
}

if (errors.length) {
  console.error("docs: FAILED\n");
  for (const e of errors) console.error(`  - ${e}`);
  console.error(`\n${errors.length} violation(s).`);
  process.exit(1);
}
console.log(`docs: clean (${files.length} markdown files, ${adrNumbers.size} ADRs, links and ADR references resolve).`);
