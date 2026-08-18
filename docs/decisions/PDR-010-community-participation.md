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

A public repository needs to receive defects, proposed direction, contributions, conduct reports, and vulnerability disclosures without turning demand into accepted plans, exposing sensitive reports, or weakening human review.

## Decision outcome

**The public front door separates demand, proposed changes, and private safety reports, and every code or corpus change still ends at human review.** Bug and proposed-goal forms distinguish reproducible defects from desired outcomes; submissions are demand, not acceptance or plan authority; and blank issues stay disabled.

Pull requests carry the change-loop declarations, evidence, and human-merge acknowledgment. Contributor reports use a dedicated `[Cliewen Conduct]` channel, and vulnerabilities use a separate `[Cliewen Security]` channel with coordinated disclosure, a seven-calendar-day acknowledgement target, a fourteen-day initial-status target, and support for `main` and the latest release. Forge-native private security intake may supplement but not replace the published security channel.

`CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, and the issue and pull-request templates carry this repository-local front door. Separate channels reject both unstructured public intake and a shared safety mailbox whose routing would blur conduct and security handling.
