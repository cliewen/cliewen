---
id: CH-059-tasks
type: tasks
status: open
links: [CH-059]
title: Tasks for CH-059
---

# Tasks

- [x] Add the `@AC-038` scenario to `docs/capabilities/CAP-001-onboarding/criteria.md` and its note to `design.md`
- [x] Detect a symlinked ancestor below the root in `internal/scaffold` and skip the subtree instead of writing through it (AC-038)
- [x] Record the blocked directory in a new `Report.Linked` category, deduplicated and sorted (AC-038)
- [x] Render the `linked` line and its summary count in `cmd/clue/init.go`, and amend the `init` usage text (AC-038)
- [x] Cover AC-038 with a positive symlink case and a negative ordinary-directory case at package level, plus the edge cases
- [x] Cover the user-visible report at command level in `cmd/clue/init_test.go` (AC-038)
- [x] Confirm the AC-025 never-overwrite tests still pass unchanged
- [x] Run `go test ./...` and `clue validate`
