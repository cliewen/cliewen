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

An existing repository's spec corpus (OpenSpec first, per [AN-002](../../analysis/AN-002-model2diagram-extraction.md)) is transformed into a Cliewen `/docs` corpus by the `clue-extract` skill ([ADR-008](../../decisions/ADR-008-extraction-is-a-skill.md)): every non-decision born `provenance: inferred` and classified by low/high reversal cost ([ADR-010](../../decisions/ADR-010-provenance-field.md), [ADR-035](../../decisions/ADR-035-bounded-provenance-and-reality-edges.md)), existing AC IDs and test tags preserved through namespaced prefixes ([ADR-009](../../decisions/ADR-009-ac-id-namespaces.md)), the source corpus and its parallel registries deleted in the same PR, `clue validate` green before review.

The brownfield reading that precedes that transform is governed by `clue-analysis`, which is why this capability also holds the analysis-evidence criteria: a verification result is classified as a clean disposable environment or a prepared environment with local prerequisites, a statistical or percentage claim names its versioned corpus and population, and adoption analysis names the governance or process changes it introduces ([AN-003](../../analysis/AN-003-robocode-api-bridge-calibration.md), [AN-004](../../analysis/AN-004-hyperfine-foreign-soil-trial.md), [AN-005](../../analysis/AN-005-es-toolkit-foreign-soil-trial.md)).

## Why

Serves [G-001](../../goals/G-001-verifiable-thread.md): the verifiable thread must be reachable from a brownfield start, not only from `clue init`. A repo that already traces tests to spec scenarios must keep that traceability through adoption — extraction that breaks existing IDs would destroy the very thread it claims to install.

Acceptance criteria: [criteria.md](criteria.md) · implementation notes: [design.md](design.md).

## Status note

`active`: the machine-checkable facets (namespace grammar, JVM tag harvesting, provenance vocabulary) are implemented in `clue validate` and covered by Go tests carrying the AC-IDs. The end-to-end extraction contract is meaning-level agent work judged by human PR review; its evidence is the P-001/M-003 extraction run. The guidance criteria (AC-054, AC-055) are proven against the rendered canonical skills by generator tests, not by a `clue validate` rule: `clue` judges what a skill says, never how an investigator applies it.
