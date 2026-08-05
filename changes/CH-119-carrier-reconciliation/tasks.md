---
id: CH-119-tasks
type: tasks
status: open
links: [CH-119]
title: Tasks for CH-119
---

# CH-119 — tasks

- [x] Read precedent: ADR-048, ADR-049, `internal/parity`, `internal/ledger`, `internal/corpus/rules.go`, `clue-extract` skill and openspec mapping.
- [x] Write `proposal.md` and this `tasks.md`.
- [x] Design and write ADR-051 (carrier reconciliation contract) — `status: inferred`, `author: agent`.
- [x] Implement `internal/carriers/carriers.go`: `Inventory`, `Entry`, `LoadInventory`, `Fingerprint`, `Reconcile`, `Report`, `Finding` (serves AC-118..AC-122 — renumbered from the proposal's AC-115..AC-119, since CH-118 (M-054) merged into `main` first and already claims AC-115..AC-117; same renumbering `f684f79` already applied to ADR-050→ADR-051).
- [x] Implement `internal/carriers/carriers_test.go`: clean baseline + one fixture per failure class + inventory-validation fixtures (serves AC-118..AC-122).
- [x] Wire `clue carriers <inventory> [root]` into `cmd/clue/main.go`: usage text, dispatch, `notifierCommands` exclusion (deterministic judge, same reasoning as `parity`).
- [ ] Add AC-118..AC-122 to `docs/capabilities/CAP-003-extract/criteria.md`.
- [x] Update `docs/capabilities/CAP-003-extract/design.md` and `README.md` with the carrier reconciliation section and status note.
- [ ] Update `internal/skills/source/skills/clue-extract.md.tmpl` (target contract item) and `internal/skills/source/resources/clue-extract/mappings/openspec.md` with the systematic carrier inventory; run `go generate ./internal/skills` to regenerate `.agents/skills/` and the `.claude/skills/` mirror.
- [ ] `go build ./...` and `go test ./...` pass.
- [ ] Digest: update `docs/plans/P-011-truthful-brownfield-migration.md` M-055 row to `done`; add `CHANGELOG.md` `[Unreleased]` entry; delete `/changes/CH-119-carrier-reconciliation/`.
- [ ] Run `clue-verify` (local checks + agentic review loop).
- [ ] Push branch, open ready PR with acceptance brief; confirm hosted head equals local `HEAD`.
