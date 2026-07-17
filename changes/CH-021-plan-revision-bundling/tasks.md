---
id: CH-021-tasks
type: tasks
status: open
links: [CH-021]
title: Tasks for CH-021
---

# Tasks

- [x] Write PDR-008 (`inferred`, pre-acceptance venue cited) and index it in the decisions README
- [x] Revise the `clue-plan` skill's semantic-mutation rule with the bundling allowance (`.agents/skills/` canonical)
- [x] Copy the revised skill into `internal/scaffold/templates/skills/` (Sanity test holds the pair byte-identical)
- [x] Add the ADR-017 scope note (register scope in scaffolded repos)
- [x] Revise M-005's exit criterion per the AC-001 split and flip it `done` with CH-020 evidence
- [x] Changelog `[Unreleased]` entry for the skill-text change (C-002)
- [x] `go test ./...` and `clue validate --forbid-changes` green (validate run without the flag while the workspace exists; the flag gates the digest commit)
