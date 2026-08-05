---
id: CH-121-tasks
type: tasks
status: open
links: [CH-121]
title: Tasks for CH-121
---

# CH-121 — tasks

- [x] Re-derive the assessment's gaps from the corpus and the tool, not from P-011's evidence column (done during planning; recorded in `proposal.md`).
- [x] Write `proposal.md` and this `tasks.md`.
- [ ] Write `docs/decisions/PDR-026-…` — `status: inferred`, `author: agent`, `supersedes` naming PDR-025's superseded clause.
- [ ] Add the forward pointer to `docs/decisions/PDR-025-brownfield-migration-precedes-simplification.md` so the register carries one answer to "which campaign is simplification".
- [ ] Write `docs/plans/P-012-…` with M-057…M-061, `status: active`.
- [ ] Write `docs/plans/P-013-…`, `status: draft`, milestone numbers unassigned until it opens.
- [ ] Close P-011: `status: completed` in frontmatter only.
- [ ] Regenerate the `docs/plans/` and `docs/decisions/` README indexes; curate the seeded description sentences.
- [ ] `go run ./cmd/clue validate` green; `go build ./...`, `go vet ./...`, `go test ./...` pass.
- [ ] `.github/scripts/completed-plans.sh` passes with `main` as base.
- [ ] Run `clue-verify` (local checks + agentic review loop).
- [ ] Push branch, open ready PR with the acceptance brief; confirm the hosted head equals local `HEAD`.
