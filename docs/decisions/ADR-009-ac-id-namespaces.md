---
id: ADR-009
type: decision
status: inferred
links: [ADR-005, ADR-007, AN-002, CAP-003]
title: AC IDs are namespaced — criteria declare an ac-prefix
author: agent
accepted-by: pending
---

# ADR-009 — AC ID namespaces

## Context and problem statement

Cliewen's AC grammar was global `AC-<digits>`. Brownfield repos arrive with their own AC ID schemes — model2diagram tags ~130 scenarios and ~270 tests with per-capability prefixes (`MG-010`, `CW-003`, [AN-002](../analysis/AN-002-model2diagram-extraction.md)). Renumbering to `AC-<digits>` would force a mass re-tag of every spec scenario and test annotation for zero semantic gain, and ADR-007 says IDs are meaning-immutable — an extraction that renames every ID is a poor start for a methodology built on provenance.

## Decision outcome

**The AC grammar generalizes to `<PREFIX>-<digits>`, namespaced per criteria file.** A `criteria.md` may declare `ac-prefix: MG` in its frontmatter (uppercase letters/digits, starting with a letter); the default is `AC`, so existing Cliewen corpora are unchanged. Consumer: `checkACTests`.

- **Declaration:** a criteria file declares ACs only in its own namespace — `@MG-010` tag lines in a file with `ac-prefix: MG`. A tag in AC-ID form with the *wrong* prefix for its file is a lint failure (wrong-namespace declaration), which also stops prose tokens from becoming accidental declarations.
- **Sharing and uniqueness:** several criteria files may share a prefix (this repo's capabilities all use `AC`); uniqueness is enforced at the ID level across the whole corpus, exactly as before (AC-013). A file whose scenarios need two prefixes is two capabilities — split it (the same instinct as splitting an over-broad AC).
- **The corpus is the registry.** Next-free-ID per prefix is `max + 1` over the declared set; external registry files (model2diagram's `test/ac-registry.md`) retire at extraction. Retirement tombstones (ADR-007) work unchanged in any namespace: `@MG-010 @retired`.
- **Test-side references** follow [ADR-005](ADR-005-test-reference-convention.md): framework-native tags where they exist — JVM `@Tag("MG_010")`, underscores because JUnit discourages hyphens in tag values, normalized to `MG-010` at harvest — and name prefixes in Go: `TestMG010_…`. References must resolve to a declared AC (AC-010) and must not point at retired ACs (AC-012), in every namespace.
- **JVM harvesting is file-level.** `clue` harvests `@Tag("…")` strings from `*Test.kt` / `*Test.java` for AC coverage and unknown/retired references; it does not parse Kotlin/Java to attribute tags to methods. Per-test purpose enforcement (ADR-006's AC-011) on the JVM is delegated to a framework-native ArchUnit rule that the extraction installs in the adopting repo — the same "framework-native first" principle as ADR-005, and it keeps `clue` free of fragile source parsing.

**Carrier:** the generalized `checkACTests` in `clue` (machine); the namespace and tagging rules in the `clue-delta` and `clue-extract` skills (agent); the ArchUnit purpose rule shipped by extraction (machine, in the adopting repo).

### Rejected: renumber extracted ACs to global `AC-<digits>`

Mass churn across ~270 annotations, severed continuity with the source repo's history, and a violation of the spirit of meaning-immutable IDs — the criteria did not change meaning, so their names should not change either.

### Rejected: dual-tagging (keep `MG_010` as a runner tag, add `@Tag("AC-xxx")`)

Two IDs meaning the same thing on every test is exactly the redundancy the one-purpose-tag rule (ADR-006) exists to prevent.
