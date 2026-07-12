---
id: CAP-003-design
type: design
status: active
links: [CAP-003]
title: Design notes for CAP-003
---

# CAP-003 — design notes

## Split of responsibility

The extraction itself is the `clue-extract` skill ([ADR-008](../../decisions/ADR-008-extraction-is-a-skill.md)) — meaning-level transform, agent-executed, human-reviewed. `clue` carries only the deterministic facets that make an extracted corpus validatable:

- **Namespace grammar** (`internal/corpus/actests.go`): the declaration harvest reads each criteria file's `ac-prefix` (default `AC`); the general tag grammar is `@<PREFIX>-<digits>` with `<PREFIX>` uppercase alphanumeric starting with a letter. Wrong-namespace tags fail, which also keeps prose tokens from becoming accidental declarations. Prefixes may be shared across files; uniqueness stays at the ID level (AC-013).
- **Go test names** generalize to `Test<PREFIX><digits>_…` for every declared prefix; `Unit`/`Sanity`/`Arch` are unchanged. A digit-bearing name whose prefix is in no declared namespace declares no purpose (AC-011) rather than silently passing.
- **JVM harvesting** is file-level: `@Tag("…")` values in `*Test.kt` / `*Test.java` / `*Tests.kt` / `*Tests.java`, underscores normalized to hyphens; values in a declared namespace are AC references (coverage, unknown, retired), everything else is runner metadata clue ignores (ADR-006). Per-test purpose enforcement on the JVM is the adopting repo's ArchUnit rule, installed by extraction (ADR-009) — clue does not parse Kotlin/Java structure.
- **Provenance** (`internal/corpus/rules.go`): optional field, `inferred|verified`, forbidden on decisions (they carry it in `status`); `clue validate`'s OK line reports the inferred count (ADR-010).

## Deliberate limits (doors)

- No per-AC required-test-type enforcement — `Test-type:` metadata from OpenSpec survives as scenario body text only (door in ADR-006).
- No source-format parsing in clue, ever — a new source is a new mapping section in the skill.
- No JVM per-method attribution: file-level harvest is the contract; the ArchUnit rule owns granularity.
- Binary distribution to adopting repos' CI (needed before `clue validate` can be a required check outside this repo) is unsolved and parked in the adopting repo's plan.
