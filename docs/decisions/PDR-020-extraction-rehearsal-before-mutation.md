---
id: PDR-020
type: decision
status: verified
links: [P-008, CAP-003, AN-003, AN-004, AN-005, ADR-008, PDR-019, C-013]
title: Extraction rehearses before it mutates
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-020 — Extraction rehearses before it mutates

## Context and problem statement

Brownfield extraction converts uncertain source meaning and deletes the parallel source corpus in one pull request. Source disagreement, prepared build environments, and governance changes must be reviewable before that mutation becomes difficult to reverse.

## Decision outcome

**Every extraction begins with a mandatory report-only rehearsal after its full change is proposed and before the target mutates.** The transient report inventories source formats and entry points, artifact mappings, preserved and minted IDs, confidence and reversal cost, test-purpose work, instruction conflicts, planned deletions, and plan doors. An unresolved conflict becomes an open question and stops mutation.

Only explicit human direction begins the mutate phase. The full change then digests the rehearsal into the durable analysis report, performs the accepted conversion, and deletes the transient workspace and parallel source corpus only when the resulting change is ready for review. The rehearsal is branch-local evidence, not a second system of record or separate change.

`clue-extract`, its generated and scaffolded copies, CAP-003, AC-056, guide and README explanations, and generator tests carry this checkpoint.
