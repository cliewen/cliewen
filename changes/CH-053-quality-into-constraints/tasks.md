---
id: CH-053-tasks
type: tasks
status: open
links: [CH-053]
title: CH-053 ordered tasks
---

# CH-053 — Tasks

- [ ] Write ADR-027 — quality scenarios are constraints; supersession-reading clause for immutable texts; QS→C mapping
- [ ] Write C-014 — coverage floor ≥ 80% (`source: ADR-027`, `enforcement: machine`)
- [ ] Write C-015 — onboarding under 30 minutes (`source: ADR-027`, `enforcement: human`)
- [ ] Add the QS-001/QS-002 tombstone row to `docs/decisions/log.md`
- [ ] Delete `docs/quality/` and `internal/scaffold/templates/docs/quality/`
- [ ] Update `docs/README.md` and scaffold `docs/README.md`: cross-cutting line, folder index, types-on-default list, folder section
- [ ] Update `docs/architecture/architecture.md` and `guide/methodology.md` diagrams (drop the QS node), `guide/corpus.md` and `guide/adoption.md` taxonomy rows
- [ ] Repoint CAP-001 `README.md`/`criteria.md`/`design.md` QS-002 references to C-015
- [ ] Edit skill sources (clue-verify, clue-extract taxonomy list, openspec mapping); bump `frontmatter.md.tmpl` to 0.6.0; `go generate ./internal/skills`
- [ ] Update `.github/pull_request_template.md` and the `ci.yml` coverage-gate step name; re-type the `corpus_test.go` quality fixture; fix the `main_test.go` comment
- [ ] CHANGELOG entry; regenerate indexes; verify (`go test ./...`, `clue validate`, `npm run guide:build`, generate drift); digest; mark M-018 done; delete workspace; ready PR
