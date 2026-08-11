---
id: CH-147-proposal
type: change
status: open
links: [P-015, M-072, C-008, C-022, ADR-057, ADR-058]
title: Repair the read-cost backlog
---

# CH-147 — Repair the read-cost backlog

## What

Inspect every live artifact reported by `clue validate --read-cost`. For each multi-document artifact, split content only when it has distinct primary consumers or record why its current form remains appropriate. For each over-budget context slice, remove only links that are genuinely redundant for that identity's reader or record the accepted path and its reason.

## Why

P-015/M-072 makes the read-cost report a repair backlog rather than a target to optimize. Completing the inspection makes its report actionable without weakening the meaning-bearing graph merely to lower a count.

## Boundary

This change does not alter the read-cost measure, its budget, or the underlying capability and policy contracts. It works only the artifacts reported by the current measurement.
