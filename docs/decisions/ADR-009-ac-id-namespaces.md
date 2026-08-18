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

Brownfield repositories bring namespaced criterion IDs. Renumbering them into a global `AC-<digits>` space would create churn and violate ADR-007's meaning-immutable identity.

## Decision outcome

> **Partially superseded by [ADR-036](ADR-036-jvm-evidence-per-executable.md) and [ADR-037](ADR-037-brownfield-ac-id-grammar.md):** the namespace and preservation rules remain; JVM evidence is per executable and the canonical grammar admits segmented prefixes and lowercase suffixes.

**Use `<PREFIX>-<digits><lowercase-suffix>` IDs, namespaced per criteria file.** A file declares `ac-prefix` (default `AC`); segments are uppercase alphanumeric and join with single hyphens. `checkACTests` enforces declarations in the file's namespace, corpus-wide canonical uniqueness, non-colliding prefixes, and references that resolve to non-retired criteria. The corpus is the registry and next-free allocation is per prefix.

Test references follow ADR-005: framework tags or normalized Go names preserve the namespace, while JVM tags are harvested only from each executable's contiguous annotation block as ADR-036 defines. Tombstones work unchanged in every namespace.

**Carrier:** `checkACTests`, the `clue-delta` and `clue-extract` skills, and ADR-036's JVM scanner. Renumbering extracted ACs or adding duplicate canonical tags is rejected because both sever identity or create redundant references.
