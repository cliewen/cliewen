---
id: PDR-026
type: decision
status: verified
links: [P-011, P-012, P-013, PDR-025, C-013, G-001]
title: A campaign closes on re-derived gate status, not on its own evidence column
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-026 — A campaign closes on re-derived gate status, not on its own evidence column

## Context and problem statement

A milestone's evidence cell is written by the implementing change and can report a true form while the external assessment's underlying judgment remains open, or call a deliberately declined request satisfied.

## Decision outcome

**A campaign closes on gate status re-derived from the corpus and tool, not on its evidence column.** Before a plan becomes `completed` against an external assessment, every gap is checked from the current artifacts and commands and classified as closed with a named mechanism and failure-path evidence, closed as a declined request with its cost stated, or open as work for a successor campaign. A `done` row starts that check; it is never the result.

P-012 closes the brownfield migration gap on re-derived evidence and simplification moves to P-013. This supersedes only PDR-025's clause naming P-012 as simplification; its migration-first reasoning and all other clauses remain. A milestone that changes what `clue validate` or `clue parity` asserts remains core-adjacent under C-013 and carries its own decision.

## Rejected: treat completed rows as closure

That makes the campaign its own judge and cannot detect an unclosed judgment or a declined request recorded as satisfied.

## Rejected: reopen P-011 or fold the work into simplification

Completed plans are immutable, and the remaining migration work has a different scope; combining it with simplification would let either half postpone the other.

## Carrier

P-011's terminal state, P-012 and P-013's scope and milestones, PDR-025's superseded clause, the plans index, and this record carry the gate.
