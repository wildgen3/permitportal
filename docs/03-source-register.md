---
title: Source register
status: scaffolded
last_reviewed: 2026-08-19
adr: [0014, 0017]
---

# Source register

Every normative claim in this repository cites an entry here. This document also carries
two things most specifications omit: the claims that were **investigated and killed**,
and the places where **published sources are themselves wrong**.

> **Status: scaffolded.** The entries below are the sources actually relied on by the
> current content. The full register is reconciled against the adversarial research pass
> that produced this architecture, and entries are added as content cites them. An
> uncited external statistic fails the build, so this register grows with the documents
> rather than ahead of them.

## Primary sources

| ID | Source | Used for | Access notes |
| --- | --- | --- | --- |
| SRC-01 | eCFR — Electronic Code of Federal Regulations | Federal regulatory text, point-in-time | No API key. Versioned retrieval with per-section amendment metadata, which gives change alerts for free. **Point-in-time floor is 2017-01-01.** Actively blocks automated traffic — see the ingestion note below. |
| SRC-02 | GPO Bulk Data (CFR XML) | Primary bulk ingest path | Bulk XML is the ingest backbone; the eCFR versioner is used for deltas. |
| SRC-03 | Federal Register API | CFR-keyed entry point to rulemaking | Joins to the public comment system on the Federal Register document number. |
| SRC-04 | Regulations.gov API v4 | Rulemaking dockets and comments | **No CFR or industry-code filters.** Must be joined via SRC-03. |
| SRC-05 | 29 CFR 1904 (OSHA recordkeeping) | Worked example: two keys at two scopes | Appendix A's NAICS-based list was adopted 2014-09-18 (79 FR 56186). |
| SRC-06 | 40 CFR 68 (EPA Risk Management Program) | Worked example: process-scope thresholds | Applicability is at **68.10(l)**, not 68.10(d) — see corrections. |
| SRC-07 | 29 CFR 1910.119 (Process Safety Management) | Worked example: chemical thresholds, categorical exclusions | |
| SRC-08 | 40 CFR 122.26(b)(14) | Evidence for open-list semantics | Only 5 of 11 industrial-activity categories reference an industry code at all. |
| SRC-09 | Census industry classification and concordance files | Scheme versions, crosswalks, revision status | Source for the reused-code list in ADR-0002. |
| SRC-10 | Census BEACON published accuracy results | The baseline this project reports against | 0.901 code accuracy; 0.755 end-to-end. Caveats in `12-eval-and-quality.md`. |
| SRC-11 | CTDL — Credential Transparency Description Language | Credential vocabulary | JSON-LD. **Ships no OWL axioms or SHACL shapes** enforcing its own ordering and boolean semantics. Two adopted terms are marked unstable. |
| SRC-12 | SKOS | Mapping relations between concepts | `exactMatch` composes; `closeMatch` does not. |
| SRC-13 | XKOS 1.2 | Correspondence as a first-class object | Defines neither apportionment ratio nor match strength; both are local extensions. |
| SRC-14 | RegData / QuantGov | Curation-queue ordering only | See ADR-0017 and the do-not-cite note below. |
| SRC-15 | 13 CFR 121.201 (SBA size standards) | Size standards by industry | Starting point for the self-owned code-to-regulation table. |
| SRC-16 | EPA MSGP Appendix N | Published SIC-to-NAICS crosswalk | A long official crosswalk document; an artifact to ingest rather than reproduce. |
| SRC-17 | GAO findings on classification consistency | Why `alternatives_shown[]` is persisted | Competent reviewers disagree on the same description. |
| SRC-18 | RCW 19.28 (Washington — electrical) | Pilot jurisdiction: trade credential chain | |
| SRC-19 | RCW 18.27 (Washington — contractor registration) | Pilot jurisdiction: bond and insurance requirements | |
| SRC-20 | RCW 19.02 (Washington — business licensing) | Pilot jurisdiction: state business licence | |
| SRC-21 | Seattle Municipal Code 5.55 | Pilot jurisdiction: municipal layer | |
| SRC-22 | ANSI/NISO Z39.19 | Facet and polyhierarchy guidance | Primary standard. Cite this rather than practitioner commentary. |
| SRC-23 | ISO 25964 | Thesaurus structure and interoperability | Primary standard. |

## Corrections — do not cite

Places where a published source is itself wrong, or where a plausible secondary claim is
false. Each of these was found by checking a primary source against a secondary one.

| Claim in circulation | Correct position |
| --- | --- |
| RMP applicability is at 40 CFR 68.10(d) | It is at **68.10(l)**. The provision was redesignated by 82 FR 4696, 84 FR 69913, and 89 FR 17686. The agency's own web page has been observed stale on this. |
| The OSHA partial-exemption industry list has been unchanged since 2001 | The **NAICS-based** list was adopted **2014-09-18** (79 FR 56186). Before that date the appendix keyed to SIC, not NAICS. |
| F1 is the geometric mean of precision and recall | It is the **harmonic** mean. This error appears in secondary write-ups of classifier results. |
| RegData's published normalized metrics can be cited as-is | They are internally inconsistent with that project's own methodology guide. **These metrics are never cited in this repository.** Use the dataset for curation-queue ordering only (ADR-0017). |
| Citations to a CFR *chapter* where the source means a *part* | Part and chapter are different levels. Check the level before transcribing a citation. |

## Killed claims

Claims that were investigated and did not survive verification. Recorded so they are not
reintroduced by a later contributor who finds the same secondary source.

Eleven claims were refuted in the research pass that produced this architecture. The
categories, so that a familiar-sounding number triggers a check rather than a citation:

- **Specific per-sector accuracy figures** attributed to the published benchmark that do
  not appear in it. Only the aggregate 0.901 / 0.755 and the sector *ordering* are
  supportable.
- **A ceiling figure for two-digit classification accuracy** that traces to no primary
  source.
- **Reported LLM classification accuracies in the high nineties**, which do not survive
  a look at their evaluation instruments.
- **A tax-agency misclassification rate** with no published basis.

> Rule: if a number in this domain is memorable, round, and uncited, treat it as false
> until a primary source is in hand. Add it here if it fails.

## Open uncertainty

`21-open-questions.md` carries what is genuinely unknown, including the one measurement
that would most cleanly size the rules layer and has never been published.

## Ingestion note

The point-in-time regulatory source actively blocks automated traffic. The ingestion
design therefore uses bulk XML (SRC-02) as the primary path with the versioner used for
deltas, a polite client with a contact user-agent and backoff, and a stated fallback.
Ingest availability is a non-functional requirement with a freshness SLO, not an
assumption. See `08-data-source-catalog.md`.
