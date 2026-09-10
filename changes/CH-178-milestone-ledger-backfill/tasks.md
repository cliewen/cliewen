---
id: CH-178-tasks
type: tasks
status: open
links: [CH-178]
title: Tasks for CH-178
---

# CH-178 — tasks

- [x] Add `corpus.LedgerMilestoneIdentities` (`internal/corpus/plans.go`), sharing `PlanMilestones`'s table extraction across every plan regardless of the plan's own status. Serves AC-187.
- [x] Wire it into `planLedgerBackfill` (`internal/migrate/migrate.go`): `done`/`dropped` milestones mark retired, `todo`/`doing` mark live. Serves AC-187.
- [x] Add AC-187 to `docs/capabilities/CAP-010-team-safe-identity-allocation/criteria.md` and focused unit evidence (a corpus with non-contiguous milestone numbering including a withdrawn milestone, and a corpus with zero plans) in `internal/corpus/plans_test.go` and `internal/migrate/migrate_test.go`.
- [ ] Move P-022 from `draft` to `active` and set M-091 `todo` → `done` with this change's evidence in the digest.
- [ ] Regenerate the `docs/plans/README.md` index for P-022's new status.
- [ ] Run the CONTRIBUTING verification block (`clue validate`) against the committed candidate.
