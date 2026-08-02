---
id: CH-105-tasks
type: tasks
status: open
links: [CH-105]
title: Tasks for CH-105
---

# Tasks

- [x] Promote the 27 `inferred` decisions to `status: verified` with `accepted-by: Flemming N. Larsen (2026-08-02, conversation)`, leaving `author: agent` and every body unchanged
- [x] Confirm the judge reports zero inferred decisions and that no record was touched beyond its two frontmatter lines
- [x] Record the batch approval as a dated row in `docs/decisions/log.md`: what was approved, how the approval was given, and what it does not make automatic
- [x] Repair the decisions index, which carries hand-written status labels `clue scaffold` preserves: 7 entries still read `inferred`, and 20 of the promoted records had bare auto-appended lines showing no title and no status at all
- [x] Tick M-048 `done` with evidence in `docs/plans/P-010-adopters-keep-current.md` (digest bookkeeping; P-010 stays `active` — five milestones remain)
- [x] Verify: `clue validate --forbid-changes`, `go test ./...`, `go build ./...`, `go vet ./...`, `gofmt -l .`
