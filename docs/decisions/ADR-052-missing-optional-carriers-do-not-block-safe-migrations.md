---
id: ADR-052
type: decision
status: verified
links: [ADR-038, ADR-039, CAP-001, P-012]
title: A missing optional generated carrier does not block an independent safe migration
author: agent
accepted-by: Flemming N. Larsen (2026-08-06, conversation)
---

# ADR-052 — A missing optional generated carrier does not block an independent safe migration

## Context and problem statement

`clue migrate` plans all safe changes as one atomic write. It also reports a missing thin CI caller, because `clue init` is the only authority that may materialize it. Treating that absence as a blocking finding made an independent ledger backfill impossible in a repository whose existing CI does not use the caller.

## Decision outcome

**A missing thin CI caller is a notice, not a blocking migration finding.** The notice names `clue init` as the materialization route, and migration never creates or rewrites the caller. An independent safe change in the same plan may apply. A present caller whose content cannot be safely recognized remains a blocking finding: migration has no authority to overwrite its semantics.

**Carrier:** `planCaller` classifies absence as a notice; AC-124 proves that the ledger backfill applies while the caller remains absent. CAP-001's design states the distinction.

### Rejected: create the caller during migration

The caller is an optional generated entry point with adopter-owned runner, source, and installation choices. Materializing it merely to unblock another migration would give migration authority that `init` deliberately owns.

### Rejected: let every finding permit partial writes

Ambiguous corpus meaning and locally modified managed carriers are unsafe preflight states. Reclassifying this one absence does not weaken atomicity for those states.
