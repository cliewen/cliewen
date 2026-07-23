---
id: C-014
type: constraint
status: active
links: [G-001, LOG-001, ADR-027]
title: Total Go statement coverage stays at or above 80%
source: ADR-027, decision log (coverage-gate row, 2026-07-12)
enforcement: machine
---

# C-014 — Test coverage ≥ 80%

Given the repository's Go module, when CI runs `go test ./... -coverprofile` and the coverage gate, then total statement coverage is at or above **80%**, and the build fails otherwise.

The AC↔test contract is the binding behavioral gate; this is the backstop tripwire for the code paths between the ACs. Total rather than per-package, so thin entry points don't force contortions. Rationale and boundaries are in the [decision log](../decisions/log.md) (coverage-gate row, 2026-07-12; demoted from ADR-004 — full text in git history); baseline at introduction (2026-07-12) was 85.1%. This constraint is the successor of retired QS-001 ([ADR-027](../decisions/ADR-027-quality-scenarios-are-constraints.md)).

**Enforcement:** `machine` — the coverage-gate step in the CI workflow is the check; a build below the floor fails.
