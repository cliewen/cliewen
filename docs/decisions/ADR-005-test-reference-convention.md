---
id: ADR-005
type: decision
status: verified
links: [CAP-002]
title: Tests reference ACs via framework-native tags; function names where no tags exist
author: human
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-005 — Tests reference ACs via framework-native tags

## Context and problem statement

The AC↔test lint needs a deterministic reference from each test to the acceptance criterion it verifies.

## Decision outcome

> **Current carrier boundary:** [ADR-071](ADR-071-framework-neutral-evidence-export.md) replaces framework parsing inside the judge with repository-owned export to a common evidence manifest. Native metadata and per-executable attribution remain; exporter examples carry framework-specific conventions.


**Use the test framework's native tag mechanism wherever one exists, and a function-name convention where none does.** Tags such as JUnit `@Tag("AC-004")`, pytest markers, and NUnit categories remain filterable and visible to the runner. The canonical identity is `<PREFIX>-<digits><lowercase-suffix>`, with documented carrier aliases normalized by the harvester. Go's `testing` fallback uses names such as `TestAC004_ValidCorpus`.

The AC harvester is per-language: each profile reads its framework's tag channel or fallback naming convention. Structured comments drift because proximity is not attached to the test, and an external mapping file would be a hand-maintained index that can rot.
