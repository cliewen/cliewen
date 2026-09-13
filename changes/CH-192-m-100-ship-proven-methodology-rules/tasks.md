## Tasks

- [x] Add `internal/skills/source/shared/challenge-commitments.md.tmpl` with the generalized challenge rule
- [x] Wire it into `clue-plan.md.tmpl` and `clue-delta.md.tmpl`, and register the new route in `internal/skills/generate.go`
- [x] Extend `internal/skills/source/shared/durable-work.md.tmpl` with the generalized guidance-capture rule
- [x] Regenerate with `go generate ./internal/skills` and confirm `go build ./...` and `go test ./internal/skills/...` pass
- [ ] Add acceptance criteria for both rules to `docs/capabilities/CAP-006-collaborative-handoffs/criteria.md`, with Go test evidence in `internal/skills/generate_test.go`
- [ ] Promote PDR-061 to `binds: adopter` with a Carrier section naming the shared template and its two inclusion sites
- [ ] Promote PDR-062 to `binds: adopter` with a Carrier section naming the shared template
- [ ] Replace AGENTS.md's repository-only challenge paragraph with a pointer to the shipped rule
- [ ] Add CHANGELOG `[Unreleased]` entry
- [ ] Digest: mark M-100 done in P-023 with evidence; set P-023 `status: completed` (last unfinished milestone, no successor named)
