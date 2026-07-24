---
id: CH-057-tasks
type: tasks
status: open
links: [CH-057]
title: Tasks for CH-057
---

# Tasks — CH-057

- [x] Write ADR-028 recording case-folded manifest resolution as the deterministic-judge rule (routed by C-013, core change)
- [x] Add acceptance criterion AC-037 to `docs/capabilities/CAP-004-ship/criteria.md` — a case-variant manifest is found identically on every filesystem; two case-variants in one directory are reported as ambiguity
- [x] Resolve the manifest by case-folded directory scan in `internal/corpus/skillversions.go` (serves AC-037)
- [x] Add the positive test (a `SKILL.md` skill joins the managed set) and the negative test (two case-variants in one directory are named) — both AC-037
- [x] Run `clue validate` and `go test ./...`
