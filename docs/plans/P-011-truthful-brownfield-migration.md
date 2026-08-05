---
id: P-011
type: plan
status: active
links: [P-010, G-001, G-002, PDR-019, PDR-025, C-013, ADR-048]
title: Cliewen makes brownfield migration truthful
---

# P-011 — Cliewen makes brownfield migration truthful

P-009 established a report-only rehearsal and source-specific extraction mappings, but its completed OpenSpec mapping still treats source archives and test registries as deletable once Git history remains. A migration assessment found that this can reuse historic identities, hide existing evidence debt, reduce in-flight work to an insufficient plan row, and leave operational carriers stale. These are general source-to-corpus parity failures, not one repository's special case.

P-011 is the successor campaign for that boundary. It opened `active` once P-010 completed (CH-114, 2026-08-05); PDR-025 moves simplification to P-012 rather than combining it with migration credibility. Every milestone is a separate full change rooted at accepted `main`, published as a ready pull request, and accepted by human merge before the next begins. Milestone numbering continues corpus-global numbering from P-010.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-052 | **Every native and imported criterion identity remains unique and meaningful after source deletion**: [ADR-048](../decisions/ADR-048-corpus-wide-id-ledger.md) defines a permanent validated identity ledger, scoped to all native ID prefixes as well as imported criteria, holding — for imported entries — the pinned source revision, source location, exact canonical identifier, and for every entry a live/reserved/retired state; numeric namespaces retain their existing allocator but read and write it through the ledger instead of a live corpus scan, while declared opaque-token namespaces preserve existing identifiers exactly and require their documented source-owned generator before minting a new identifier; the validator rejects a repeated live, tombstoned, or reserved identity and a source mapping cannot delete an incomplete inventory; new greenfield corpora retain the current numeric default; focused positive and negative fixtures prove an archived numeric ID above the live maximum and a UUID-like opaque ID cannot be reused | `todo` | |
| M-053 | **Migration proof parity is reproducible at criterion and evidence-location level**: each source mapping emits a pinned source manifest and Cliewen derives a target manifest containing criteria, exact identity state, proof class, positive/negative direction, evidence locations, and declared exclusions; a deterministic parity command and CI artifact compare them, report every unmatched or altered entry, and fail deliberate missing criterion, orphaned tag, changed direction/location, stale source fingerprint, and unjustified `@draft`, `Human`, or retirement disposition; the report stays derived rather than becoming an editable coverage registry | `todo` | |
| M-054 | **In-flight source work survives as durable normalized work**: an adopter-defined `imported-change` record holds each source change's pinned origin, intent, design rationale, dependencies, task-to-criterion proof links, and completion state; extraction rejects deletion of active source work until its record is complete, and a fixture proves a proposal, design, dependency, and proof task remain inspectable without retaining a parallel source corpus | `todo` | |
| M-055 | **Migration reconciles every operational carrier before deletion**: the rehearsal inventories instructions, workflows, freshness inputs, registries, local/external links, and diagram assets; every carrier maps to a retained target or blocks mutation; a mapped carrier must be current at the pinned source revision, and CI proves stale deleted-path references, lost fingerprints, and missing assets fail before a migration PR can merge | `todo` | |
| M-056 | **A disposable brownfield migration proves the full contract**: numeric-archive and opaque-identifier fixture sources execute report-only rehearsal, approved mutation, target validation, parity generation, and required-CI failure paths; the result proves the contract without naming a production adopter and distinguishes deterministic Cliewen checks from the fixture source's own test results | `todo` | |

## Explicitly out of this campaign

P-011 does not make `clue validate` execute source test suites, use Git history as an implicit registry, accept arbitrary prose as an opaque identifier, or introduce a permanent parallel source corpus. Source-specific interpretation remains in extraction mappings; the common ledger and parity contract only establishes what every mapping must prove before source deletion. Simplification belongs to P-012.

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a decision record routed by reversal cost. Plan adjustments are decisions.
