---
id: AN-002
type: analysis
status: active
links: [P-001, M-003, G-001]
title: Extraction target analysis — model2diagram (OpenSpec)
---

# AN-002 — Extraction target analysis: model2diagram

Findings from the read-only survey of `model2diagram` (2026-07-12), the guinea pig for P-001/M-003 (brownfield extraction). Verified by human review of the CH-004 PR.

## The repository

Kotlin 2.2 / Gradle multi-module (`api` → `core` → `cli`), version 0.1.0, main clean. Converts JSONSchema to Mermaid class diagrams via an internal `ModelGraph`. Quality gates already in place: detekt, ktlint, JaCoCo with an 80% coverage threshold, GitHub Actions CI. 32 production Kotlin files, 14 test files (13 Kotlin + 1 Java interop).

## The existing spec corpus: OpenSpec

The repo runs OpenSpec (`openspec/` with `config.yaml`, schema `spec-driven`) — its ADR 0002 records the adoption. Three artifact families exist:

1. **Synced capability specs** — `openspec/specs/<capability>/spec.md`, one per capability, eight in total: `model-graph` (MO), `jsonschema-parser` (JP), `mermaid-generator` (MG), `mermaid-validator` (MV), `config-system` (CS), `diagram-service` (DS), `cli-wrapper` (CW), `graph-validator` (GV). Structure: `# Spec:` header → `## ADDED Requirements` → `### Requirement: <title>` with a SHALL statement → one or more `#### Scenario: <name> [MG-010]` blocks, each carrying a `Test-type:` line and `- **WHEN**` / `- **THEN**` / `- **AND**` bullets. Roughly 100 requirements / 130 scenarios across the eight files.
2. **A pending change** — `openspec/changes/plantuml-generator/`: proposal.md, design.md (decisions D1–D5, risks R1–R3), tasks.md (~55 checkboxes, **all unchecked** — not implemented), and a spec delta `specs/plantuml-generation/spec.md` with 17 requirements spanning two scenario prefixes: `PG` (generator) and `PV` (PlantUML test validator). No API or CLI changes — `OutputFormat.PlantUML` and the factory hook already exist.
3. **One archived change** — `changes/archive/2026-05-25-mvp-jsonschema-to-mermaid/`, the fully applied MVP that seeded all eight capabilities (218/218 tasks checked). History, not live truth.

## Traceability already present

Every spec scenario carries a per-capability AC ID, and every scenario-verifying test carries that ID as a JUnit tag plus a test-type tag: `@Tag("MG_010")` + `@Tag("UNIT")`. Observed tag volume: 201 UNIT, 35 E2E, 34 INTEGRATION, 5 ARCHITECTURAL, roughly two tests (positive + negative) per AC. Governance lives outside the specs in `test/ac-registry.md` (monotone per-prefix counters, PG/PV already reserved) and `test/scenario-template.md`. ADR 0005 there records the AC-ID/tag convention.

**Three notations coexist for the same logical ID** — the main parsing hazard for any converter: `[MG-010]` (brackets, synced specs), `` `PG-001` `` (backticks, the pending delta), `MG_010` (underscores, JUnit tags — JUnit tag values discourage hyphens).

## Other corpus material

`docs/architecture/` (12 files: vision-and-scope, module/component architecture, test-strategy, …), `docs/decisions/` (5 Nygard-style ADRs), `AGENTS.md` (engineering + OpenSpec workflow guide), 10 `openspec-*` skills under `.agents/skills/`.

## Mapping conclusions (consumed by ADR-008, ADR-009, ADR-010 and the clue-extract skill)

- OpenSpec **capability spec → Cliewen capability** (`README.md` / `criteria.md` / `design.md`); Requirement + Scenario → Gherkin AC. The scenario IDs are already acceptance-criterion IDs in everything but name — they must survive extraction unchanged or ~270 test tags break (ADR-009: namespaced AC prefixes).
- The **pending change → plan milestone + `status: draft` capability** whose criteria carry the PG/PV ACs; design decisions D1–D5 survive in its design.md. Draft criteria are exempt from the AC↔test contract until implementation (AC-009), which is exactly the pending state.
- The **tag scheme maps 1:1 onto ADR-005/ADR-006**: the AC tag is the purpose tag, `UNIT`/`INTEGRATION`/`E2E` are runner tags clue ignores. Only the handful of tests with no AC tag (the 5 ARCHITECTURAL ones) need a purpose re-tag during extraction — the feared untagged-suite phase-in problem does not exist here.
- The **AC registry and scenario template retire**: uniqueness and next-free-ID are derivable from the corpus, which `clue validate` already lints (AC-013). The corpus is the registry.
- `Test-type:` lines are per-AC required-test-type metadata — a door Cliewen has deliberately not built (ADR-006). They survive as scenario body text so the information is not lost.
- Archived OpenSpec changes map to **git history only** — Cliewen's transient `/changes` model has no archive folder by design.

## Risks

- The semantic transform (SHALL prose → Gherkin, architecture docs → capability design.md, five ADRs → Cliewen ADRs) is meaning-level work — agent judgment reviewed by a human PR, not deterministic conversion. This is why extraction is a skill, not clue code (ADR-008).
- clue cannot run in model2diagram's CI yet (private-repo binary distribution unsolved) — extraction evidence is a local `clue validate` pass plus the reviewed PR; CI wiring is a named door in that repo's plan.
