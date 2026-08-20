# Security

## Scope

PermitPortal is a specification and reference implementation. There is no hosted
production system, no user accounts, and no personal data. The repository
contains synthetic fixtures only.

## Reporting

Report suspected vulnerabilities through GitHub's
[private security advisory](https://github.com/wildgen3/permitportal/security/advisories/new)
flow rather than a public issue.

## What counts as a security issue here

- A dependency with a known exploitable vulnerability reachable from committed code.
- A CI workflow configuration that could leak repository secrets.
- Anything that would cause the clean-room scanner to disclose the terms it screens for.

## Data classification

`spec/model/` defines a `restricted` data class covering chemical inventory
attributes. Attributes in that class carry `llm_egress_allowed: false` and must
never be transmitted to a third-party model provider or analytics endpoint. A
change that widens LLM egress for a restricted attribute is a security change and
requires an ADR.
