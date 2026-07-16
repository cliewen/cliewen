---
id: CH-017-tasks
type: tasks
status: open
links: []
title: CH-017 task breakdown
---

# Tasks

- [x] Write the PDR (review-boundary rules; rejected: global serialization, per-plan serialization, local-merge-as-acceptance) and link it from the proposal
- [x] Write the constraint (agents never merge own changes; enforcement: agent; promotion trigger: branch protection / PR-provenance CI) and link it from the proposal
- [x] Update `clue-delta` step 1 (branch from current `main` tip, one change in flight per author, stacking escalates to the human)
- [x] Update `clue-delta` step 5 (agent boundary: never merge own PR, never local merge commits, never push `main`; review fixes stay on this branch)
- [x] Update `clue-delta` step 4 (the digest itself is never a task in `tasks.md` — the precondition applies to the work, not to the deletion act)
- [x] Update `clue-delta` small-deltas line (rebase onto `main` after a sibling merges, re-run `clue-verify`)
- [x] Update `clue-verify` with process checklist items (branch cut from `main` tip, no unreviewed predecessor, no foreign `/changes/`, hosted PR is the only route and the human merges)
- [x] Update AGENTS.md binding rules with the boundary sentence linking the constraint (skills stay ID-free per the repo-agnostic rule; AGENTS.md carries the link)
- [x] Update index blocks: `docs/constraints/README.md`, `docs/decisions/README.md`
- [x] Add CHANGELOG `[Unreleased]` entry for adopters
- [x] Verify: `go build ./... && go test ./...`, `go run ./cmd/clue validate`, re-read both skills for light-tier coherence
