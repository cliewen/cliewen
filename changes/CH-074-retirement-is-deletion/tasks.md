---
id: CH-074-tasks
type: tasks
status: open
links: [CH-074]
title: Tasks for CH-074
---

# Tasks

- [x] Read PDR-003, ADR-025, `docs/decisions/log.md`, `docs/decisions/README.md`, and the `clue-extract` MADR mapping to establish exact current text before drafting
- [ ] Draft ADR-034 (`Retirement is deletion; supersedes carries the pointer forward`): defines the `supersedes:` frontmatter field, states retirement = delete-in-the-same-change, states the two exceptions (criteria tombstones, completed plans), and states which clauses of PDR-003 and ADR-025 it supersedes
- [ ] Add `Supersedes []string` to `internal/corpus.Artifact`, parsed from frontmatter like `Links`
- [ ] Add `checkSupersedes` validator rule (`internal/corpus/rules.go`): a `supersedes:` entry whose ID still resolves to a live artifact is rejected; enrich `checkLinks`' dangling-link message when the unresolved ID matches a `supersedes:` entry elsewhere, naming the successor
- [ ] Add AC criteria on CAP-002 (validate) for the new checks, `@AC-xxx` tagged, `Test-type: Unit`, positive and negative Go tests (`TestAC0xx_*`)
- [ ] Revise PDR-003: demotion mechanic uses `supersedes:` on the successor (`docs/decisions/log.md`'s own frontmatter for demotions with no successor record) instead of ad hoc "row added, inbound references repointed" prose
- [ ] Revise ADR-025: state that `retired` is not a resting status any default-lifecycle type's file is ever observed carrying on `main` — retirement is the deletion event, not a status value a file holds — and point to ADR-034
- [ ] Update `docs/decisions/README.md`'s retention paragraph and `docs/decisions/log.md`'s header preamble to describe the generic `supersedes:` mechanism
- [ ] Update `internal/skills/source/resources/clue-extract/mappings/madr.md`'s `superseded by` row if its meaning changed, regenerate skills (`go generate ./internal/skills`), confirm no drift
- [ ] Update `clue-delta` skill source if the retirement/tombstone sentence needs the `supersedes:` field named, regenerate skills
- [ ] Regenerate `docs/decisions/README.md` index if entries change
- [ ] Update `docs/plans/P-007-core-hardening.md`: M-027 status `todo` → `done` with evidence
- [ ] Add CHANGELOG entry if this is user-visible (it changes `clue validate` behavior — yes)
- [ ] `go test ./...` and `clue validate` green; regenerate skills with no drift
