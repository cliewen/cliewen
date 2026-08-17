---
id: C-014
type: constraint
status: active
links: [G-001, ADR-027, IDR-008]
title: Total Go statement coverage stays at or above 80%
source: ADR-027, IDR-008
enforcement: machine
---

# C-014 — Test coverage ≥ 80%

Given the repository's Go module, when CI runs `go test ./... -coverprofile` and the coverage gate, then total statement coverage is at or above **80%**, and the build fails otherwise.

The AC↔test contract is the binding behavioral gate; this is the backstop tripwire for the code paths between the ACs. Total rather than per-package, so thin entry points do not force contortions ([IDR-008](../decisions/IDR-008-go-coverage-is-an-aggregate-tripwire.md)). This constraint is the successor of retired QS-001 ([ADR-027](../decisions/ADR-027-quality-scenarios-are-constraints.md)).

**Checked by:** the coverage-gate step in the CI workflow; a build below the floor fails.
