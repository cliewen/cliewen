---
id: CH-118-tasks
type: tasks
status: open
links: [CH-118]
title: Tasks for CH-118
---

# Tasks

- [x] Write ADR-050 deciding the `imported-change` artifact type, its frontmatter shape, its lifecycle vocabulary, and why it is durable rather than transient (ADR-034 boundary)
- [x] Add `imported-change` to `internal/corpus/rules.go`'s `statusVocabExceptions` (`in-progress`, `complete`)
- [x] Implement `internal/importedchange/` — parsing of the proof-links table plus the `complete` completeness check, following `internal/parity`'s fixture-test shape
- [x] Wire a new validator rule (`checkImportedChanges`) into `internal/corpus/rules.go`'s `Validate` rule list
- [x] Add AC-115..AC-117 to `docs/capabilities/CAP-003-extract/criteria.md` (continuing after AC-114) with focused Go fixtures
- [x] Create `docs/imported-changes/README.md` with a `clue:index` block, and add the folder to `docs/README.md`'s Folders index and status-vocabulary table
- [x] Write worked `docs/imported-changes/IC-001-*.md` and `IC-002-*.md` fixture records proving a proposal, design, dependency, and proof task (both proven and in-progress-unproven) remain inspectable (M-054's stated fixture requirement)
- [x] Update `internal/skills/source/resources/clue-extract/mappings/openspec.md`'s pending-change row and `internal/skills/source/skills/clue-extract.md.tmpl`'s deletion-gating clause, regenerated via `go generate ./internal/skills` into `.agents/skills/` and `internal/scaffold/templates/skills/`
- [x] Update CAP-003 `README.md` and `design.md` to describe the imported-change record
- [x] `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...`
- [x] `clue validate --forbid-changes --coverage .` clean (only the expected `/changes` present issue remains, resolved at digest)
- [ ] Digest: set M-054 `done` in `docs/plans/P-011-truthful-brownfield-migration.md`, add CHANGELOG entry, delete `/changes/CH-118-imported-change-records/`
- [ ] Run `clue-verify` and open the ready PR
