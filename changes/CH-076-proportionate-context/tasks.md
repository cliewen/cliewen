---
id: CH-076-tasks
type: tasks
status: open
links: [CH-076]
title: Tasks for CH-076
---

# Tasks

- [ ] Pin the adopting repository evidence boundary and read M-029, AN-008, PDR-002, the parked `clue locate` record, context-relevant capabilities, and current CLI/graph implementation
- [ ] Run a disposable history-measurement spike against the adopting product repository and record observed facts, inferences, rejected measurements, and limits in a new analysis artifact
- [ ] Write the PDR deciding the change-tier boundary from the measured evidence
- [ ] Add or revise acceptance criteria for `clue context <id>` before implementing behavior
- [ ] Implement deterministic transitive context slicing with focused positive and negative evidence
- [ ] Expose `clue context <id>` through the CLI and consume the parked `clue locate` idea
- [ ] Consolidate tier routing at its canonical template source and align AGENTS.md, generated skills, scaffold output, and public guide carriers with the PDR
- [ ] Update capability documentation and command-facing guide material for focused context discovery
- [ ] Regenerate affected skills and indexes and confirm no generated drift
- [ ] Update P-007 M-029 from `todo` to `done` with evidence
- [ ] Add a user-facing `[Unreleased]` CHANGELOG entry
- [ ] Run focused tests, `go test ./...`, `clue validate`, and the strict guide build
