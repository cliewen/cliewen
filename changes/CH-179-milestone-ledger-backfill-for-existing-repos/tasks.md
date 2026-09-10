---
id: CH-179-tasks
type: tasks
status: open
links: [CH-179]
title: Tasks for CH-179
---

# CH-179 — tasks

- [x] Add `planMilestoneLedgerBackfill` (MIG-016, `internal/migrate/migrate.go`), guarded to run only against an already-existing ledger this plan is not itself creating or rewriting. Serves AC-188.
- [x] Add AC-188 to `docs/capabilities/CAP-010-team-safe-identity-allocation/criteria.md` and disposable-Git integration evidence (`internal/migrate/migrate_test.go`) exercising migrate + Git-coordinated `clue id next M` end to end against a repository shaped like this one: an existing coordinated ledger missing milestone entries.
- [x] Add supplementary focused unit evidence for the backfill and no-op cases, and update the registry-order test for the new migration ID.
- [x] Add a user-facing `[Unreleased]` `CHANGELOG.md` entry for the new migration.
- [ ] Move P-022/M-092 `todo` → `done` with this change's evidence in the digest.
- [ ] Run the CONTRIBUTING verification block (`clue validate`) against the committed candidate.
