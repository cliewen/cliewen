---
id: ADR-051
type: decision
status: verified
links: [ADR-008, ADR-048, ADR-049, PDR-019, PDR-025, P-011, C-013]
title: A pinned carrier inventory and a deterministic reconciliation check close the operational-carrier gap
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-051 — Carrier reconciliation manifests

## Context and problem statement

Migration can preserve a source decision while losing an operational carrier such as an instruction, workflow, freshness input, registry, link, or diagram asset. A pinned inventory is needed to make every carrier's retained target or known block inspectable and reproducible.

## Decision outcome

**The extraction rehearsal writes a pinned carrier inventory, and `clue carriers` re-derives target fingerprints and reports reconciliation findings without writing either side.** Each entry names its fixed kind, source path, and either a retained target plus fingerprint or `blocked: true` with a reason; `deleted-paths` records source paths that must no longer be referenced. The fixed vocabulary can grow only through a plan decision.

Reconciliation reports stale deleted-path references, lost fingerprints, and missing assets; blocked entries remain explicit gaps until a later inventory revision maps them. Reports are sorted and deterministic. Source-format discovery remains `clue-extract` work, and the command is separate from `clue validate` like parity because it needs a source-side manifest; the migrating repository's workflow makes it required.

**Carrier:** the inventory schema, `clue carriers`, fingerprint and deleted-path checks, the extraction rehearsal, and migration CI guidance.
