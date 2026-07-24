# Tasks — CH-057

- [ ] Write ADR-028 recording case-folded manifest resolution as the deterministic-judge rule (routed by C-013, core change)
- [ ] Add acceptance criterion AC-034 to `docs/capabilities/CAP-004-ship/criteria.md` — a case-variant manifest is found identically on every filesystem; two case-variants in one directory are reported as ambiguity
- [ ] Resolve the manifest by case-folded directory scan in `internal/corpus/skillversions.go` (serves AC-034)
- [ ] Add the positive test (a `SKILL.md` skill joins the managed set) and the negative test (two case-variants in one directory are named) — both AC-034
- [ ] Run `clue validate` and `go test ./...`
