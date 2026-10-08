---
id: CH-203-tasks
type: tasks
status: open
links: [CH-203]
title: Local allocation stops on any trace tasks
---

# Tasks

- [ ] Edit the canonical source `internal/skills/source/skills/clue-delta.md.tmpl` step 1: the existence of `clue/id-allocator` on the remote is a stop, and any `ch-*` branch other than your own stops, merged or not, with the answer ending the question; run `go generate ./internal/skills` and inspect the generated and scaffolded `change-loop.md` (AC-218).
- [ ] Revise AC-218's scenario and its Unit tests in `internal/skills/generate_test.go` to the new wording, keeping a negative test; regenerate the evidence export.
- [ ] Amend PDR-066, the `[Unreleased]` changelog entry, and CAP-010's `design.md` sentence to match.
