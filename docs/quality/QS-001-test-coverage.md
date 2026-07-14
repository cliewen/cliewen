---
id: QS-001
type: quality
status: active
links: [G-001, LOG-001]
title: Total Go statement coverage stays at or above 80%
---

# QS-001 — Test coverage ≥ 80%

**Scenario:** given the repository's Go module, when CI runs `go test ./... -coverprofile` and the coverage gate, then total statement coverage is at or above **80%**, and the build fails otherwise.

Rationale and boundaries in the [decision log](../decisions/log.md) (coverage-gate row, 2026-07-12; demoted from ADR-004 — full text in git history): the AC↔test contract is the binding behavioral gate; this is the backstop tripwire for code paths between the ACs. Total rather than per-package, so thin entry points don't force contortions. Baseline at introduction (2026-07-12): 85.1%.
