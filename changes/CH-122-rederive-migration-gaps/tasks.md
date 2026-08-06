---
id: CH-122-tasks
type: tasks
status: open
links: [CH-122]
title: Tasks for CH-122
---

# CH-122 — tasks

- [x] Re-derive the assessment's gaps and blocking gates from the current corpus and commands; record the proposed dispositions in `proposal.md`.
- [x] Write `proposal.md` and this task list.
- [-] Run the existing migration to seed this repository's identity ledger, then inspect the resulting allocation and removal of MIG-008 — blocked: `go run ./cmd/clue migrate --apply` refuses every write while `.github/workflows/clue.yml` is absent (MIG-003); see `open-questions.md`.
- [ ] Write the sanitized assessment analysis and its gate-by-gate PDR-026 disposition register (after the migration boundary is decided, so its disposition can be stated truthfully).
- [ ] Record the branch-protection acceptance-evidence boundary and update its methodology carriers.
- [ ] Repair CAP-003's stale adopter-CI binary-distribution claim against the reusable validation workflow and verified installation scripts.
- [ ] Update M-057 bookkeeping, README indexes, and `[Unreleased]`.
- [ ] Run focused tests and corpus validation.
- [ ] Run `clue-verify` and resolve its review findings.
- [ ] Push the reviewed commit and open a ready PR with its acceptance brief; confirm its hosted head equals `HEAD`.
