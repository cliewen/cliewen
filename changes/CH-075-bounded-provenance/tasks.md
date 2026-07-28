---
id: CH-075-tasks
type: tasks
status: open
links: [CH-075]
title: Tasks for CH-075
---

# Tasks

- [x] Read M-028, AN-007, AN-008, ADR-010, ADR-025, ARCH-003, C-011, C-013, and the current provenance/coverage implementation
- [ ] Draft ADR-035 defining reversal-cost classification for inferred meaning, the active-capability blocking edge, actionable blocker output, and the incident-analysis edge back from reality
- [ ] Add or revise CAP-002 acceptance criteria for activation blocking and the derived “green met wrong” capability listing before implementing behavior
- [ ] Implement the reversal-cost field parsing and validation, including inferred decisions in the same provenance model
- [ ] Implement the active-capability dependency check with focused positive and negative Unit evidence
- [ ] Replace the monotonic inferred count with actionable blocker output and focused command evidence
- [ ] Implement the incident-analysis marker and derived affected-capability listing with focused positive and negative Unit evidence
- [ ] Update AN-007 to carry the failed capability or criterion edge and the incident marker required by the convention
- [ ] Update CAP-002 README/design, corpus/scaffold documentation, and public guide text that states the old inferred-counter contract
- [ ] Update generated skill source if the incident-analysis convention has an agent carrier, regenerate skills, and confirm no drift
- [ ] Regenerate affected indexes
- [ ] Update P-007 M-028 from `todo` to `done` with evidence
- [ ] Add a user-facing `[Unreleased]` CHANGELOG entry
- [ ] Run `go test ./...`, `clue validate`, the strict guide build, and generation drift checks
