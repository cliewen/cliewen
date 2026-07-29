# CH-082 — tasks

- [ ] Pin the evidence boundary: Tank Royale source revision, Cliewen revision, and the toolchain, OS and shell of every reproduction, with each result classified clean-disposable or prepared
- [ ] Inventory the three assumption families as they stand today
  - [ ] CI wall: runner label, pinned actions, download mechanism, asset name, install path, and the version the adopter's wall actually pins
  - [ ] Agent-directory placement: what already occupies `.agents` and `.claude` in the adopter, and what Cliewen writes or expects there
  - [ ] Skills: the managed `clue-*` set against any repository-local skill or command the adopter keeps beside it
- [ ] Reproduce each assumption against the adopter and record the concrete edits or failures it forces, on a disposable clone wherever reproduction would mutate anything
- [ ] Separate mandatory needs from maintainer preferences, listing every case that cannot be settled from the repository alone as an open question
- [ ] Evaluate candidate configuration locations against ADR-013's AGENTS-local-layer boundary and its stated `clue`-needs-repo-local-settings condition
- [ ] Write the findings to `docs/analysis/AN-012-<slug>.md`, including what was rejected and why
- [ ] Route the outcome: explicit ADR/PDR candidates by reversal cost, and a named successor-plan consumer
- [ ] Regenerate the analysis index and verify `go run ./cmd/clue validate --forbid-changes`, coverage, and `go test ./...`
