---
id: CH-123-tasks
type: tasks
status: open
links: [CH-123]
title: Tasks for CH-123
---

# CH-123 — tasks

- [x] Define the structured deferred-disposition contract and its compatibility boundary in ADR-053, serving AC-125.
- [x] Add AC-125 and focused positive and negative Unit evidence for accountable dispositions and the derived population report.
- [x] Implement manifest validation, parity comparison, and deterministic reporting for AC-125.
- [x] Update CAP-003 guidance and generated indexes to carry the new contract.
- [x] Run focused and full verification — `go test ./internal/parity ./cmd/clue`, `go run ./cmd/clue validate`, `go test ./...`, and `git diff --check` pass.
