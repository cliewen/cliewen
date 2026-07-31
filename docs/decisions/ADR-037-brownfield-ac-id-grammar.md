---
id: ADR-037
type: decision
status: inferred
links: [P-009, M-038, ADR-009, CAP-002, CAP-003]
title: Brownfield criterion IDs preserve segmented prefixes and letter suffixes
author: agent
accepted-by: []
---

# ADR-037 — Brownfield criterion IDs preserve segmented prefixes and letter suffixes

## Context and problem statement

ADR-009 preserves source criterion identities by allowing an uppercase namespace prefix followed by a numeric identity, but its grammar cannot express stable brownfield forms such as `SNAP-SQS-001` and `ADP-045b`. Renumbering those criteria would sever existing links and evidence for no change in meaning.

## Decision outcome

The canonical acceptance-criterion ID grammar is `<PREFIX>-<NUMBER><SUFFIX?>`, where `<PREFIX>` is one or more uppercase alphanumeric segments joined by single hyphens, `<NUMBER>` is one or more decimal digits, and `<SUFFIX>` is an optional run of lowercase ASCII letters. Thus `SNAP-SQS-001` and `ADP-045b` are valid, while underscores, empty segments, uppercase suffixes, and mixed-case forms are not. IDs are case-sensitive and declarations and links use the canonical spelling without normalization.

`ac-prefix` carries the `<PREFIX>` portion and accepts the same segmented uppercase grammar. The default remains `AC`. A declaration belongs to its criteria file only when its full ID starts with that file's exact prefix plus the final numeric separator; wrong-namespace tags fail. Full IDs remain unique across the corpus, and retirement tombstones keep their exact canonical IDs.

Evidence carriers normalize only their syntax into the canonical ID: literal JVM and Cucumber tags may use underscores in place of hyphens, while Go and JVM named executables remove hyphens from the prefix because language identifiers cannot contain them. The numeric component and lowercase suffix remain unchanged, so `SNAP-SQS-001` maps to `TestSNAPSQS001_UnitPositive_name`, `testSNAPSQS001_UnitPositive_name`, or `@Tag("SNAP_SQS_001")`, and `ADP-045b` maps to `TestADP045b_UnitPositive_name`, `testADP045b_UnitPositive_name`, or `@Tag("ADP_045b")`. Declared prefixes whose hyphen-stripped forms collide are rejected because the named carriers could not resolve them unambiguously.

Extraction preserves every source ID that is present. For requirements without a source ID, it assigns the next numeric slot after the maximum numeric component already declared in that namespace, ignoring letter suffixes for the maximum, in the source's stated order; an empty namespace starts at one. Minted IDs use the canonical prefix and decimal number, never rename preserved IDs, and are recorded in the extraction report.

**Carrier:** the canonical parser and namespace checks in `internal/corpus`, the supported evidence harvesters and context/coverage consumers, and the preservation and minting rules in the canonical and generated `clue-extract` skills.

### Rejected: renumber unsupported source criteria

Renumbering would break existing test annotations and repository history while changing no criterion meaning.

### Rejected: accept arbitrary punctuation or case and normalize it everywhere

Broad normalization would create collisions and hide malformed or wrong-namespace evidence; the judge keeps declarations and links exact and limits aliases to the documented evidence-carrier syntax.
