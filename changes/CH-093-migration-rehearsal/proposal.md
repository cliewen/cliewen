---
id: CH-093
type: change
status: open
links: [P-009, M-042, PDR-004, PDR-020, PDR-021, AN-014, AN-013, CAP-003]
title: Rehearse a representative brownfield migration without mutation
---

# CH-093 — Rehearse a representative brownfield migration without mutation

## What

This full change serves P-009 milestone M-042 by running the mandatory report-only `clue-extract` rehearsal against the clean local Robocode Tank Royale checkout, the repository the human named as this milestone's representative target. Two revisions are pinned: the head pin `384d27d55176a2d2ad4668ac381852e629e4540a`, the public adopter revision already used by AN-013, and the source pin `4e579878fd6667fab3b75515b6e68135a935c8df`, the last state carrying the parallel `openspec/` corpus before adoption removed it.

Because that target is already adopted, the rehearsal is retrospective: it reconstructs a conversion that happened between the two pins rather than previewing a pending one. The exit items this cannot exercise are recorded as narrower boundaries rather than claimed.

The rehearsal will inventory the source formats and entry points, artifact mappings, preserved and minted criterion IDs, confidence and reversal cost, test-purpose work, instruction conflicts, governance effects, planned deletions, CI and installation constraints, current-to-target corpus migration, merge-mode compatibility, and named follow-up doors. It will record Cliewen's deterministic validation separately from the target repository's own test and build results, including whether each result came from a clean disposable or prepared environment.

The target checkout, tests, routing, and hosted state remain unchanged. The rehearsal report was drafted as transient evidence under this change workspace; on 2026-08-01 the human answered OQ-003 by directing that it survive the merge, so this change's one durable corpus addition is the rehearsal analysis artifact AN-015, landed under explicit human direction rather than under PDR-020's mutate phase, which this change still does not authorize. Any unresolved conflict becomes an open question and stops before mutation.

## Why

M-042 is the final migration-readiness gate before P-009's distributed-work contracts. AN-014 records the evidence boundary that the confidential assessment lacked, while PDR-020 requires a pinned, reviewable rehearsal before any brownfield conversion can delete a parallel source corpus or change governance.

## Decision boundary

This change assesses migration truth and operational viability only. It does not mutate Tank Royale, perform an extraction, add a Cliewen configuration interface, import foreign test results as local proof, resolve external references over the network, or widen the support boundary without a recorded decision.
