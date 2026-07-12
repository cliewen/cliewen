---
id: QS-001
type: quality
status: active
links: [G-001, ADR-004]
title: Total Go statement coverage stays at or above 80%
---

# QS-001 — Test coverage ≥ 80%

**Scenario:** given the repository's Go module, when CI runs `go test ./... -coverprofile` and the coverage gate, then total statement coverage is at or above **80%**, and the build fails otherwise.

Rationale and boundaries in [ADR-004](../decisions/ADR-004-coverage-gate-80-percent.md): the AC↔test contract is the binding behavioral gate; this is the backstop tripwire for code paths between the ACs. Baseline at introduction (2026-07-12): 85.1%.
