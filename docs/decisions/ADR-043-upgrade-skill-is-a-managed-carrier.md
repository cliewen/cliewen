---
id: ADR-043
type: decision
status: verified
links: [P-010, CAP-001, CAP-004, ADR-021, ADR-022, ADR-031, ADR-039]
title: The managed skill set includes a human-authorized upgrade entry point
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-043 — The upgrade skill is a managed carrier

## Context and problem statement

Once `clue latest` reports drift, adopters need a named route that can explain the release, ask the human whether to act, and carry an explicit approval through the normal reviewed repository change without inventing a second installation mechanism.

## Decision outcome

**The six generated, repository-scoped skills are `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, `clue-upgrade`, and `clue-verify`; `clue-upgrade` is the human-authorized upgrade entry point.** It checks the release, reads migration guidance, asks whether to act now or later, and carries a yes through `clue migrate`; it contains no installation command because `clue latest` already selects the machine-specific route.

The sixth skill has the same ownership marker, version stamp, canonical source, generated tree, embedded scaffold template, and reserved managed slot as the other five. MIG-003 may add it only from a complete recognized prior release set; missing, partial, or modified carriers block the write, and an optional Claude mirror is never materialized. The Claude plugin remains a bootstrap and ships none of the managed skills.

**Carrier:** the canonical and generated skill trees, scaffold templates, migration manifest, and upgrade guidance.
