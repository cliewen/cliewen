---
id: ADR-048
type: decision
status: verified
links: [ADR-009, ADR-034, ADR-007, P-011, C-013]
title: A persisted ledger replaces scan-and-max allocation for every native ID prefix
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-048 — A corpus-wide ID ledger

## Context and problem statement

Scanning the corpus to allocate the next ID cannot remember reserved or deleted identities and becomes increasingly expensive. Brownfield migrations also need to preserve opaque source IDs and their provenance without renumbering live work.

## Decision outcome

**`.clue/id-ledger.yaml` becomes the registry for every native ID prefix, and `clue validate` cross-checks the corpus against it rather than allocating from a scan.** Entries carry exact canonical IDs, `numeric | opaque` kind, `reserved | live | retired` state, numeric prefix/component where applicable, and per-prefix counters; imported entries may carry `source-revision` and `source-location`. Reserved IDs become live when the artifact appears, retired IDs are never removed or reissued, and opaque allocation remains source-owned.

After one load, exact-ID and numeric-next lookups are map operations. Validation rejects missing or contradictory entries once the ledger is active, rejects malformed numeric/opaque fields, and preserves ADR-034's separate `supersedes:` pointer and ADR-007's criterion tombstone. An idempotent backfill seeds current live IDs and counters without renumbering history.

This expands M-052's scope from imported criteria to all native prefixes; it does not implement the ledger, allocator, validator rule, or backfill itself.

**Carrier:** the ledger format and loader, future allocation and validation commands, the backfill migration, and the ID-preservation guidance.
