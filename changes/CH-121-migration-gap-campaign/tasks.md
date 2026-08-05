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
- [x] Write `docs/decisions/PDR-026-…` — `status: inferred`, `author: agent`; the partial supersession is body prose naming PDR-025's superseded clause, following ADR-024's precedent, because ADR-034 reserves the `supersedes:` field for IDs a retirement deleted.
- [x] Add the forward pointer to `docs/decisions/PDR-025-brownfield-migration-precedes-simplification.md` so the register carries one answer to "which campaign is simplification".
- [x] Write `docs/plans/P-012-…` with M-057…M-061, `status: active`.
- [x] Write `docs/plans/P-013-…`, `status: draft`, milestone numbers unassigned until it opens.
- [x] Close P-011: `status: completed`, forward links added, and its out-of-campaign line corrected from P-012 to P-013 — a plan revision the mutation rules allow because PDR-026 backs it, and the last change permitted to touch the file.
- [x] Regenerate the `docs/plans/` and `docs/decisions/` README indexes; curate the seeded description sentences.
- [x] `go run ./cmd/clue validate` green; `go build ./...`, `go vet ./...`, `go test ./...` pass.
- [x] `.github/scripts/completed-plans.sh main HEAD` reports no completed plan modified — P-011 was `active` on the base, so this is the last change permitted to touch it.
- [-] Run `clue-verify` (local checks + agentic review loop) — cannot be ticked here: it reviews the complete candidate, which includes the digest commit that deletes this file. The PR body carries its result.
- [-] Push branch, open ready PR with the acceptance brief — same reason; it follows the digest, and the PR is its own evidence.
