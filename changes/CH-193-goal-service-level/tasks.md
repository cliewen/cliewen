---
id: CH-193-tasks
type: tasks
status: open
links: [CH-193]
title: Tasks for CH-193
---

# Tasks for CH-193

- [x] Promote G-018, G-019, G-020 from `proposed` to `accepted` and regenerate `docs/goals/README.md`.
- [x] Add `GoalState` (ID, Title, Status, the capabilities naming it via `goal:`, and the plans naming it via `links`, each ref with its own status) to `internal/corpus/intent.go`'s `IntentState`, populated by a bounded reverse lookup over every capability's required `goal:` field and every plan's `links`.
- [x] Print the goal listing in `clue validate --intent` (`cmd/clue/main.go`): each goal's capabilities and plans, each shown with its own status, "none" printed plainly for either list when empty — no percentage, ratio, or count anywhere in the addition. Serves AC-202.
- [ ] Add `@AC-202` to `docs/capabilities/CAP-009-product-intent/criteria.md` with focused positive and negative `Unit` evidence in `internal/corpus/intent_test.go` and `cmd/clue/intent_test.go`.
- [ ] Update `docs/capabilities/CAP-009-product-intent/README.md`'s `What`/`Why`/evidence-range prose and `links` (add G-020) for the new behavior.
- [ ] Write PDR-063 recording the decision and the two rejected alternatives (new goal status; a coverage ratio), linking G-020, ADR-025, PDR-054, ARCH-003, CAP-009.
- [x] Run this repository's own `clue validate --intent` and read the new goal lines against G-001, G-003, and G-018/G-019/G-020 to check the challenge — result recorded in `proposal.md`'s Challenge section.
- [ ] Add a CHANGELOG `[Unreleased]` entry naming the new `--intent` goal listing (shipped `clue` behavior).
- [ ] Digest: regenerate indexes, delete this workspace, update plan bookkeeping (none — plan-less).
