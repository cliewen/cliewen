---
id: CH-152
type: change
status: open
links: [G-005, CAP-002, CAP-005, C-016]
title: Constraint index badges state enforcement
---

# CH-152 — Constraint index badges state enforcement

## Proposal

This change is plan-less: it serves G-005's targeted repair to generated constraint index badges rather than a planned campaign item.

New constraint rows in a taxonomy README will show their `enforcement:` class in the badge instead of the generic artifact status, and `clue validate --index-rows` will name rows whose constraint badge disagrees with that field. The change will retire the status-only index-row criterion and add AC-138 with focused positive and negative Unit evidence for generation and reporting.

## Scope boundary

The change does not rewrite existing curated rows, normalize old indexes, or alter badges for artifact types other than constraints.
