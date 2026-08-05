---
id: CH-120-tasks
type: tasks
status: open
links: [CH-120]
title: Tasks for CH-120
---

# CH-120 — tasks

- [x] Read the M-052–M-055 implementations, ADR-048 through ADR-051, CAP-003, and the canonical `clue-extract` skill.
- [x] Create the two fixture source inventories and write a report-only rehearsal for each without mutating the target corpus (serves AC-056).
- [x] Record the required mutation authorization in `open-questions.md` and stop before mutation.
- [x] Obtain explicit human direction to begin the approved fixture mutation after the rehearsals are complete.
- [x] Add or revise the criterion that defines the disposable end-to-end migration proof, then add focused positive and negative evidence in its declared test class.
- [x] Implement the approved fixture mutations and their deterministic validation, parity, and carrier-reconciliation failure paths.
- [-] Update canonical extraction guidance and regenerate generated skill carriers — the existing guidance already requires the composed ledger, parity, imported-work, and carrier contract; CAP-003 now records the fixture evidence, so no guidance contract changed.
- [x] Run focused checks, `go build ./...`, and `go test ./...`.
- [ ] Digest: update P-011 M-056, add the user-facing CHANGELOG entry, and delete this workspace.
- [ ] Run `clue-verify` and its agentic review loop.
- [ ] Push the reviewed branch, open a ready PR with the acceptance brief, and confirm its hosted head equals local `HEAD`.
