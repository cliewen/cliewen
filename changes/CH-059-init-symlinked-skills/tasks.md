# Tasks — CH-059

- [ ] Add the `@AC-038` scenario to `docs/capabilities/CAP-001-onboarding/criteria.md` and its note to `design.md`
- [ ] Detect a symlinked ancestor below the root in `internal/scaffold` and skip the subtree instead of writing through it (AC-038)
- [ ] Record the blocked directory in a new `Report.Linked` category, deduplicated and sorted (AC-038)
- [ ] Render the `linked` line and its summary count in `cmd/clue/init.go`, and amend the `init` usage text (AC-038)
- [ ] Cover AC-038 with a positive symlink case and a negative ordinary-directory case at package level, plus the edge cases
- [ ] Cover the user-visible report at command level in `cmd/clue/init_test.go` (AC-038)
- [ ] Confirm the AC-025 never-overwrite tests still pass unchanged
- [ ] Run `go test ./...` and `clue validate`
