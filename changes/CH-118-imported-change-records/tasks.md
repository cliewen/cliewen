---
id: CH-118-tasks
type: tasks
status: open
links: [CH-118]
title: Tasks for CH-118
---

# Tasks

- [ ] Write ADR-050 deciding the `imported-change` artifact type, its frontmatter shape, its lifecycle vocabulary, and why it is durable rather than transient (ADR-034 boundary)
- [ ] Add `imported-change` to `internal/corpus/rules.go`'s `statusVocabExceptions` (`in-progress`, `complete`)
- [ ] Implement `internal/importedchange/` — parsing of the proof-links table plus the `complete` completeness check, following `internal/parity`'s fixture-test shape
- [ ] Wire a new validator rule (`checkImportedChanges`) into `internal/corpus/rules.go`'s `Validate` rule list
- [ ] Add AC-115..AC-117 to `docs/capabilities/CAP-003-extract/criteria.md` (continuing after AC-114) with focused Go fixtures
- [ ] Create `docs/imported-changes/README.md` with a `clue:index` block, and add the folder to `docs/README.md`'s Folders index and status-vocabulary table
- [ ] Write one worked `docs/imported-changes/IC-001-*.md` fixture record proving a proposal, design, dependency, and proof task remain inspectable (M-054's stated fixture requirement)
- [ ] Update `.agents/skills/clue-extract/mappings/openspec.md`'s pending-change row and `internal/skills/source/resources/clue-extract/mappings/openspec.md` (regenerate via `go generate ./internal/skills`), plus the scaffold template copy
- [ ] Update CAP-003 `README.md` and `design.md` to describe the imported-change record
- [ ] `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...`
- [ ] `clue validate --forbid-changes --coverage .` clean
- [ ] Digest: set M-054 `done` in `docs/plans/P-011-truthful-brownfield-migration.md`, add CHANGELOG entry, delete `/changes/CH-118-imported-change-records/`
- [ ] Run `clue-verify` and open the ready PR
