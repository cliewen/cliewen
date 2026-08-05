---
id: CH-120-tasks
type: tasks
status: open
links: [CH-120]
title: Tasks for CH-120
---

# CH-120 — tasks

- [ ] Read the M-052–M-055 implementations, ADR-048 through ADR-051, CAP-003, and the canonical `clue-extract` skill.
- [ ] Create the two fixture sources and write a report-only rehearsal for each without mutating the target corpus (serves AC-056).
- [ ] Record any unresolved source-evidence or mapping decision in `open-questions.md` and stop before mutation.
- [ ] Obtain explicit human direction to begin the approved fixture mutation after the rehearsals are complete.
- [ ] Add or revise the criterion that defines the disposable end-to-end migration proof, then add focused positive and negative evidence in its declared test class.
- [ ] Implement the approved fixture mutations and their deterministic validation, parity, and carrier-reconciliation failure paths.
- [ ] Update CAP-003 and canonical extraction guidance with the composed fixture workflow; regenerate generated skill carriers if changed.
- [ ] Run focused checks, `go build ./...`, and `go test ./...`.
- [ ] Digest: update P-011 M-056, add the user-facing CHANGELOG entry, and delete this workspace.
- [ ] Run `clue-verify` and its agentic review loop.
- [ ] Push the reviewed branch, open a ready PR with the acceptance brief, and confirm its hosted head equals local `HEAD`.
