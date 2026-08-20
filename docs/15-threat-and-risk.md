---
title: Threat and risk
status: scaffolded
last_reviewed: 2026-08-19
adr: [0001, 0014]
---

# Threat and risk

> **Status: scaffolded.** **Minimum viable version:** the legal-exposure posture and the
> hazard-data sensitivity note, both of which already constrain the model.

## To specify

- STRIDE over the system, with the intake and determination paths as the primary
  surfaces.
- **Legal exposure and unauthorized practice.** The posture is settled: informational
  determinations, mandatory citation on every surfaced item, no auto-filing in v1,
  documented escalation to the issuing authority, and an explicit instruction to confirm
  against current text for the applicant's jurisdiction.
- **Adversarial classification gaming.** A business that selects a favourable code to
  avoid an obligation. Mitigations: the confirmed code is recorded with provenance and
  the alternatives shown; obligations key to attributes rather than to codes, so gaming
  the code alone moves fewer determinations than it appears to.
- **The security sensitivity of hazard intake.** An aggregated database of who stores
  what chemicals, in what quantity, and where, is a target. This is why chemical
  attributes are `restricted`, why they never reach a model provider, and why access
  control on them is not the same as on the rest of the profile.
- Supply chain: pinned actions, dependency review, and the secret handling that the
  clean-room scanner depends on.
