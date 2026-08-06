---
id: PDR-025
type: decision
status: inferred
links: [P-010, P-011, PDR-019, PDR-026, C-013]
title: Truthful brownfield migration precedes simplification
author: agent
accepted-by: []
---

# PDR-025 — Truthful brownfield migration precedes simplification

## Context and problem statement

P-010 defers simplification to P-011. A brownfield migration assessment instead established that the current extraction contract can lose identity reservations, evidence state, active-work traceability, and operational carriers. These are failures of the verifiable thread and deterministic judge, not a local source-format inconvenience. Simplification cannot be allowed to delay a contract that would otherwise let an extraction look green while changing what the source proves.

## Decision outcome

**P-011 is the migration-only campaign, and simplification moves to P-012.** [PDR-026](PDR-026-campaigns-close-on-re-derived-evidence.md) supersedes the second half of that sentence and nothing else: P-012 finishes the migration gap on re-derived evidence and P-013 is the simplification campaign. P-011 remains `draft` until P-010 completes. Its work is ordered from identity preservation through evidence and carrier parity to disposable end-to-end fixtures; it does not start a second active campaign or bypass P-010's remaining milestone.

**Brownfield identity is preserved, not inferred from a visible sequence.** A later P-011 decision will define a permanent, machine-checked identity ledger and the exact opaque-token grammar. New corpora retain the current numeric allocation convention. The planned contract changes what `clue validate` asserts about criterion identity and evidence, so each semantic milestone is a full change with its own decision record and human acceptance under C-013.

## Rejected: retain P-011 for simplification

The present migration boundary is a higher-risk defect in an existing promise. Keeping its resolution behind an unrelated simplification campaign would make a known unsafe extraction rule the published default.

## Rejected: implement migration work inside P-010

P-010 exists to make an already-defined upgrade path discoverable and to close its own measured backlogs. Adding a new core migration contract would conflate the campaigns and defer an explicit plan decision until implementation is already under way.

## Carrier

P-010's forward-plan boundary, P-011's scope and milestones, the plans index, and this record carry the sequencing decision. P-011's future contract decisions will update every live methodology carrier they affect under PDR-019.
