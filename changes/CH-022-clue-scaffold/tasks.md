---
id: CH-022-tasks
type: tasks
status: open
links: [CH-022]
title: Tasks for CH-022
---

# Tasks

- [x] Export a materialization-free entry point in the scaffold package (`Regen`): docs-tree precheck, index regeneration, missing-README reporting (AC-026, AC-027)
- [x] Add the `clue scaffold` command: report + summary, errors to stderr, usage text in main.go (AC-026, AC-027)
- [x] CAP-005 capability folder: README, criteria (AC-026/AC-027 as Gherkin), design; indexed in the capabilities README
- [x] Tests: positive + negative pair per AC, at package and command level
- [x] Correct the stale CAP-001 index annotation (`draft` → `active`)
- [x] Changelog `[Unreleased]` entry (C-002)
- [x] `go test ./...` and `clue validate` green (plus an end-to-end smoke: init → add artifact → scaffold → green validate; no-docs path exits 1)
