---
id: T-082
type: tasks
status: open
links: [CH-082, P-008]
title: Tasks for pricing adopter configuration
---

# CH-082 — tasks

- [x] Pin the evidence boundary: Tank Royale source revision, Cliewen revision, and the toolchain, OS and shell of every reproduction, with each result classified clean-disposable or prepared
- [x] Inventory the three assumption families as they stand today
  - [x] CI wall: runner label, pinned actions, download mechanism, asset name, install path, and the version the adopter's wall actually pins
  - [x] Agent-directory placement: what already occupies `.agents` and `.claude` in the adopter, and what Cliewen writes or expects there
  - [x] Skills: the managed `clue-*` set against any repository-local skill or command the adopter keeps beside it
- [x] Reproduce each assumption against the adopter and record the concrete edits or failures it forces, on a disposable clone wherever reproduction would mutate anything — no disposable clone was needed: every reproduction was read-only
- [x] Separate mandatory needs from maintainer preferences, listing every case that cannot be settled from the repository alone as an open question
- [x] Evaluate candidate configuration locations against ADR-013's AGENTS-local-layer boundary and its stated `clue`-needs-repo-local-settings condition
- [x] Write the findings to `docs/analysis/AN-012-adopter-configuration-cost.md`, including what was rejected and why
- [x] Route the outcome: explicit ADR/PDR candidates by reversal cost, and a named successor-plan consumer
- [x] Regenerate the analysis index and verify `clue validate --forbid-changes`, coverage, and `go test ./...`
