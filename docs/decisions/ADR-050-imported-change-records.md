---
id: ADR-050
type: decision
status: verified
links: [ADR-034, ADR-048, ADR-049, P-011, CAP-003, C-013]
title: In-flight source work becomes a durable imported-change record, never a transient workspace
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-050 — Imported change records

## Context and problem statement

Source extraction can preserve criteria while losing the source change that explained their intent and proof mapping. That work needs a durable native record before source material is deleted, with a lifecycle that expresses completion without inventing a retirement state.

## Decision outcome

**`docs/imported-changes/` holds one durable `imported-change` record per preserved source change.** IDs use the `IC` prefix from the ledger; records use `in-progress → complete`, carry `source-revision` and `source-location`, link their dependencies, and contain Origin, Intent, Design rationale, and a Proof links table. A complete record must link only declared, non-draft, non-retired criteria; an in-progress record may name unfinished proof.

The folder has its own index and corpus documentation, and the status vocabulary omits `retired` because the record is durable. `clue` does not read the source repository or decide whether source work may be deleted; the extraction skill keeps that human-reviewed rule, while the machine check verifies the record's completion claim.

**Carrier:** imported-change corpus validation, the ID ledger, the extraction rehearsal and skill, and the folder/index documentation.
