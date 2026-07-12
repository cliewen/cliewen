---
id: CH-004
type: change
status: open
links: [P-001, M-003, G-001]
title: Brownfield extraction — machinery and the model2diagram run
---

# CH-004 — Brownfield extraction

## Why

M-003 is the last open exit criterion of [P-001](../../docs/plans/P-001-elaboration-baseline.md): brownfield extraction works on model2diagram. Without an extraction story, Cliewen only serves greenfield repos — the methodology must be able to adopt an existing codebase and its existing spec corpus without losing traceability that is already there.

## What

Two parts; this change is part one (the cliewen-side machinery). Part two is the extraction run itself, a separate PR in the model2diagram repository.

1. **AN-002** — analysis of the extraction target: model2diagram runs OpenSpec with 8 synced capability specs, one pending change (plantuml-generator), 5 ADRs, and tests already tagged with per-capability AC IDs (`@Tag("MG_010")`).
2. **ADR-008** — extraction is a generic skill (`clue-extract`) with per-source mappings; OpenSpec is the first mapping. clue stays source-agnostic.
3. **ADR-009** — AC ID namespaces: criteria.md may declare `ac-prefix`; the AC grammar generalizes to `<PREFIX>-<digits>`; prefix collisions are lint failures; the corpus replaces external AC registries. JVM tests declare ACs via `@Tag("<PREFIX>_<digits>")`, harvested by clue.
4. **ADR-010** — provenance field: extracted artifacts are born `provenance: inferred`; human review promotes to `verified`; clue lints the vocabulary and reports the inferred count.
5. **Code**: `checkACTests` learns prefixes and JVM tag harvesting; `rules.go` learns the provenance vocabulary; `clue validate` reports inferred counts.
6. **CAP-003-extract** — the extraction capability with ACs (AC-014…) and tests.
7. **`clue-extract` skill** — the target contract plus the OpenSpec mapping.

## Decisions already settled with the human (2026-07-12)

Extractor is one generic skill; AC IDs generalize (keep `MG-010`-style IDs, no renumbering); `/docs` replaces `/openspec` as system-of-record after extraction; this change plus the extraction run land in one arc.

## Out of scope (doors)

Per-AC required-test-type enforcement (`Test-type:` lines survive as scenario body text only); clue binary distribution / model2diagram CI wiring; `clue init`/`clue scaffold` automation; non-OpenSpec mappings; per-method JVM purpose enforcement in clue (a framework-native ArchUnit rule carries it); implementing the PlantUML generator.
