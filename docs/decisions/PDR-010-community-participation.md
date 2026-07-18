---
id: PDR-010
type: decision
status: verified
links: [G-003, P-003, M-009, PDR-007]
title: Community participation enters through structured intake, private safety channels, and human review
author: agent
accepted-by: Flemming N. Larsen (2026-07-19, planning conversation — conduct standard, reporting channels, response targets, and issue intake chosen explicitly)
---

# PDR-010 — Community participation has a structured and safe front door

## Context and problem statement

A public repository needs to welcome bugs, new direction, contributions, conduct reports, and vulnerability disclosures without turning issue activity into accepted plans, exposing sensitive reports, or weakening the human review boundary. Which entry points let strangers participate while preserving Cliewen's system-of-record and decision mechanics?

## Decision outcome

**Cliewen's public front door separates demand, proposed changes, and private safety reports, and every code or corpus change still ends at human review.**

- **Issues distinguish defects from direction.** A bug form asks for reproducible evidence. A proposed-goal form asks who wants an outcome and why; submission records demand, not acceptance or plan authority. Blank issues stay disabled so security reports and unsupported free-form requests are not invited into public discussion.
- **Pull requests carry the change loop.** The template requires a CH identity, plan item or explicit plan-less declaration, tier, traceability and decision impact, verification evidence, digest state, and acknowledgment that a human maintainer merges. [PDR-007](PDR-007-review-boundary.md) remains unchanged.
- **Contributor Covenant 3.0 governs participation.** Conduct reports go privately to `flemming.n.larsen+cliewen-conduct@gmail.com` with a `[Cliewen Conduct]` subject. Access to that channel stays narrower than general project discussion even if its backing mailbox or team routing changes.
- **Vulnerabilities use coordinated disclosure.** Reports go privately to `flemming.n.larsen+cliewen-security@gmail.com` with a `[Cliewen Security]` subject; reporters are told not to open public issues. Cliewen aims to acknowledge within seven calendar days and provide an initial status within fourteen, without promising a fixed remediation date. Security support covers `main` and the latest release; older releases must upgrade.
- **Forge-native security intake is additive.** Private vulnerability reporting may supplement email where the forge supports it, but the published email remains stable across hosting changes.

**Carrier:** `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, and the repository's issue and pull-request templates. Repo-local — none of these files ships to adopting repositories through `clue init`.

### Rejected: unstructured public intake

Blank issues make the exceptional path the easiest path and invite vulnerability disclosure, support requests without a support strategy, and feature activity that looks more authoritative than it is.

### Rejected: public security issues

Public disclosure before triage can expose users before a fix exists. A coordinated private channel lets the reporter and maintainer establish impact, remediation, and disclosure timing first.

### Rejected: one safety mailbox without routing

Security and conduct reports need different handling and potentially different future access. Separate stable aliases preserve that boundary while still allowing a solo maintainer to route both into one underlying mailbox today.
