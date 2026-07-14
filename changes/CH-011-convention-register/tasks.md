---
id: CH-011-tasks
type: tasks
status: open
links: []
title: CH-011 task breakdown
---

# Tasks

- [x] AC-023 added to CAP-002 criteria.md: constraint artifacts are linted (source + enforcement present, vocabulary machine|agent|human, agent-count reported)
- [x] `checkConstraints` rule in `internal/corpus/rules.go` (serves AC-023)
- [x] `cmd/clue/main.go` OK line reports the `enforcement: agent` constraint count (serves AC-023)
- [x] Tests: AC-023 positive + negative in corpus tests, count reporting in cmd tests
- [x] Seed constraints C-001…C-010 written, each with `source`, `enforcement: agent` (C-009 notes its partial machine promotion), and a promotion trigger
- [x] `docs/constraints/README.md` rewritten as the register (intro + index)
- [x] ADR-017 — conventions register as constraints with enforcement classes (born `inferred`)
- [x] `docs/architecture/architecture.md` — "enforcement classes beyond machine" removed from the deliberately-out list
- [x] `docs/capabilities/CAP-002-validate/design.md` — constraint rule added to the rules list
- [x] `AGENTS.md` rules 4/5/6 gain pointers to C-004/C-001/C-002
- [ ] Digest: CHANGELOG `[Unreleased]` entry, delete `/changes/CH-011-convention-register/`
- [ ] Verify: `go test ./...` green, `go run ./cmd/clue validate --forbid-changes` green, clue-verify checklist walked
