---
id: ADR-057
type: decision
status: inferred
links: [P-015, CAP-002, CAP-007, PDR-029, ADR-056, C-022, ADR-058]
title: Read cost is reported as a structural backlog, never scored by size
author: agent
accepted-by: []
---

# ADR-057 — Read-cost measurements

## Context and problem statement

Read cost needs observable signals without pretending that bytes, words, or artifact count are a quality score. Two structural shapes reveal that one durable file serves multiple documents or that a default context entry point is wider than its focused budget.

## Decision outcome

**`clue validate` reports two non-blocking populations: live `docs/` artifacts with more than one rendered H1, and identities whose default context slice prints more than eight artifacts.** Setext and ATX H1 forms count; completed plans are excluded because C-008 makes their historical shape immutable. Identity rows count entry paths, including criteria and milestones, while the measurement examines only the bounded slice, not its frontier or wider closure.

The counts appear on successful validation and are listable with `--read-cost`; they are derived on every run, write no registry, and never change the exit code. They measure navigation structure only; PDR-029 remains the judgment for whether a document or link belongs together.

**Carrier:** the read-cost checks and flag, the focused-context slice, P-015, C-022, and the structural guidance.
