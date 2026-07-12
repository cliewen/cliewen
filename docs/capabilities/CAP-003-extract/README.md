---
id: CAP-003
type: capability
status: active
links: [G-001, ADR-008, ADR-009, ADR-010]
title: Brownfield extraction — adopt an existing corpus without losing its thread
goal: G-001
---

# CAP-003 — Brownfield extraction

## What

An existing repository's spec corpus (OpenSpec first, per [AN-002](../../analysis/AN-002-model2diagram-extraction.md)) is transformed into a Cliewen `/docs` corpus by the `clue-extract` skill ([ADR-008](../../decisions/ADR-008-extraction-is-a-skill.md)): every artifact born `provenance: inferred` ([ADR-010](../../decisions/ADR-010-provenance-field.md)), existing AC IDs and test tags preserved through namespaced prefixes ([ADR-009](../../decisions/ADR-009-ac-id-namespaces.md)), the source corpus and its parallel registries deleted in the same PR, `clue validate` green before review.

## Why

Serves [G-001](../../goals/G-001-verifiable-thread.md): the verifiable thread must be reachable from a brownfield start, not only from `clue init`. A repo that already traces tests to spec scenarios must keep that traceability through adoption — extraction that breaks existing IDs would destroy the very thread it claims to install.

Acceptance criteria: [criteria.md](criteria.md) · implementation notes: [design.md](design.md).

## Status note

`active`: the machine-checkable facets (namespace grammar, JVM tag harvesting, provenance vocabulary) are implemented in `clue validate` and covered by Go tests carrying the AC-IDs. The end-to-end extraction contract is meaning-level agent work judged by human PR review; its evidence is the P-001/M-003 extraction run.
