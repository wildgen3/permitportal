## What and why

<!-- One paragraph. Link the issue this closes. -->

Closes #

## Decision reference

<!-- The ADR this implements, or an explicit statement that nothing architectural changed. -->

- [ ] Implements ADR-____
- [ ] No architectural change

## Clean-room checklist

The scanner runs regardless of what is ticked here. These are the things it cannot check.

- [ ] No third-party names appear anywhere in this change
- [ ] Every new or modified rule entry carries a `source_url` to public authority
- [ ] Every new normative claim in `docs/` carries a source reference
- [ ] No non-synthetic data was added
- [ ] Nothing here originates from a non-public source, in substance or in structure

## Eval impact

<!-- Baseline delta table, or "n/a". A moved baseline needs a reason, not just a number. -->

n/a

## Status discipline

- [ ] Directory `status:` values updated if a phase changed, and the root README table agrees
- [ ] `./do check` passes locally
