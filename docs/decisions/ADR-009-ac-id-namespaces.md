---
id: ADR-009
type: decision
status: verified
links: [ADR-005, ADR-007, AN-002, CAP-003, ADR-037]
title: AC IDs are namespaced — criteria declare an ac-prefix
author: agent
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-009 — AC ID namespaces

## Context and problem statement

Cliewen's AC grammar was global `AC-<digits>`. Brownfield repos arrive with their own AC ID schemes — model2diagram tags ~130 scenarios and ~270 tests with per-capability prefixes (`MG-010`, `CW-003`, [AN-002](../analysis/AN-002-model2diagram-extraction.md)). Renumbering to `AC-<digits>` would force a mass re-tag of every spec scenario and test annotation for zero semantic gain, and ADR-007 says IDs are meaning-immutable — an extraction that renames every ID is a poor start for a methodology built on provenance.

## Decision outcome

> **Partially superseded by [ADR-036](ADR-036-jvm-evidence-per-executable.md) and [ADR-037](ADR-037-brownfield-ac-id-grammar.md):** the namespace, evidence-normalization, and ID-preservation principles below remain current, while JVM evidence is per executable and the canonical grammar now admits segmented prefixes and lowercase letter suffixes.

**The AC grammar generalizes to `<PREFIX>-<digits><lowercase-suffix>`, namespaced per criteria file.** A `criteria.md` may declare `ac-prefix: MG` or a segmented prefix such as `SNAP-SQS` in its frontmatter; each segment is uppercase alphanumeric and starts with a letter, segments join with single hyphens, and the default is `AC`, so existing Cliewen corpora are unchanged. Consumer: `checkACTests`.

- **Declaration:** a criteria file declares ACs only in its own namespace — `@MG-010` tag lines in a file with `ac-prefix: MG`. A tag in AC-ID form with the *wrong* prefix for its file is a lint failure (wrong-namespace declaration), which also stops prose tokens from becoming accidental declarations.
- **Sharing and uniqueness:** several criteria files may share a prefix (this repo's capabilities all use `AC`); uniqueness is enforced at the full canonical ID level across the whole corpus, exactly as before (AC-013), and prefixes whose hyphen-stripped forms collide are rejected because named carriers would be ambiguous. A file whose scenarios need two prefixes is two capabilities — split it (the same instinct as splitting an over-broad AC).
- **The corpus is the registry.** Next-free-ID per prefix is the next numeric slot after the maximum numeric component already declared, ignoring lowercase letter suffixes; a source corpus's parallel registry retires at extraction. [AN-002](../analysis/AN-002-model2diagram-extraction.md) preserves the concrete extraction evidence without making the private target's registry file a normative reference. Retirement tombstones (ADR-007) work unchanged in any namespace: `@MG-010 @retired` or `@ADP-045b @retired`.
- **Test-side references** follow [ADR-005](ADR-005-test-reference-convention.md): framework-native tags where they exist — JVM `@Tag("MG_010")` or `@Tag("SNAP_SQS_001")`, underscores because JUnit discourages hyphens in tag values, normalized to the canonical ID at harvest — and normalized name prefixes in Go: `TestMG010_…` or `TestSNAPSQS001_…`. References must resolve to a declared AC (AC-010) and must not point at retired ACs (AC-012), in every namespace.
- **JVM harvesting is per executable.** `clue` harvests literal `@Tag("…")` strings from the contiguous JUnit annotation block of each supported `*Test.kt` / `*Test.java` executable, together with the stable named-executable form; unrelated class or method metadata receives no classified credit. ADR-036 defines the conservative parser and its unsupported-shape diagnostic.

**Carrier:** the generalized `checkACTests` in `clue` (machine); the namespace and tagging rules in the `clue-delta` and `clue-extract` skills (agent); and the per-executable JVM scanner defined by ADR-036 (machine).

### Rejected: renumber extracted ACs to global `AC-<digits>`

Mass churn across ~270 annotations, severed continuity with the source repo's history, and a violation of the spirit of meaning-immutable IDs — the criteria did not change meaning, so their names should not change either.

### Rejected: dual-tagging (keep `MG_010` as a runner tag, add `@Tag("AC-xxx")`)

Two IDs meaning the same thing on every test is exactly the redundancy the one-purpose-tag rule (ADR-006) exists to prevent.
