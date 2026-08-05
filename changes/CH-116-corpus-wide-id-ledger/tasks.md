---
id: CH-116-tasks
type: tasks
status: open
links: [CH-116]
title: Tasks for implementing the corpus-wide identity ledger
---

- [x] Add AC-101..AC-107 to `docs/capabilities/CAP-002-validate/criteria.md` and `docs/capabilities/CAP-001-onboarding/criteria.md`: ledger-backed numeric allocation, opaque verbatim preservation, checkLedger retired-id rejection, checkLedger shape rejection (numeric-without-component / opaque-with-component), the two M-052 fixtures (archived numeric ID above live max, UUID-like opaque ID reuse), and the migrate backfill step
- [x] Implement `internal/ledger/` (`Ledger`, `Entry`, `Load`, `Save`, `NextNumeric`, `ReserveOpaque`, `MarkLive`, `IsUsed`) with positive/negative unit tests naming AC-101/AC-102/AC-105
- [x] Implement `clue id next <prefix>` in `cmd/clue/main.go` (+ `runID`), wired through `internal/ledger/`
- [x] Implement `checkLedger` in `internal/corpus/rules.go`, registered in `Validate`, gated on `.clue/id-ledger.yaml` presence, with positive/negative unit tests naming AC-103/AC-104/AC-106
- [x] Implement the ledger backfill migration step in `internal/migrate/migrate.go` (new `MIG-008` constant + `planLedgerBackfill`), idempotent, with a test proving a second run produces zero changes
- [ ] Update `docs/capabilities/CAP-002-validate/design.md` with the ledger's role in `checkLedger`
- [ ] Update `docs/capabilities/CAP-002-validate/README.md`'s rule list to mention the ledger cross-check
- [ ] Add `[Unreleased]` CHANGELOG entry: `clue validate` now cross-checks IDs against a persisted ledger, and `clue id next` allocates through it
- [ ] Set M-052 to `done` in `docs/plans/P-011-truthful-brownfield-migration.md`, citing this change's evidence
- [ ] Digest: delete this change workspace
