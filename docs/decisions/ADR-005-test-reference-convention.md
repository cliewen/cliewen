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

The AC↔test lint (M-002) needs a deterministic way to know which acceptance criteria a test verifies. Where does the reference live?

## Decision outcome

**Use the test framework's native tag mechanism wherever one exists; fall back to a function-name convention where none does.**

- **Primary convention — framework tags**: JUnit 5 `@Tag("AC-004")`, pytest `@pytest.mark.ac004` / custom markers, NUnit `[Category("AC-004")]`, and equivalents. Tags are the framework's own metadata channel: filterable in runners, visible in reports, and attached to the test by the framework itself rather than by naming discipline.
- **Fallback — function names**: for frameworks with no tag mechanism (Go's `testing`), a test references `AC-<digits>` by carrying `AC<digits>` in its function name (`TestAC004_ValidCorpus…`). The name is what `go test -run`, CI logs and grep already display. (This fallback was the original agent proposal; the general tag-first rule is the human decision.)

Consequence for `clue`: the AC harvester is **per-language** — each language profile reads that framework's tag mechanism. The baseline ships the Go harvester (names); harvesters for tagged frameworks arrive with their language profiles, reading tags.

### Rejected: structured comments

`// verifies: AC-004` drifts — nothing attaches it to the test but proximity, and refactoring silently detaches it.

### Rejected: external mapping file

A manifest linking tests to AC-IDs is exactly the hand-maintained index the methodology forbids: the first artifact type to rot.
