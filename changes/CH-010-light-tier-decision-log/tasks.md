---
id: CH-010-tasks
type: tasks
status: open
links: []
title: CH-010 task breakdown
---

# Tasks

- [x] ADR-015 — light change tier (qualification criteria, PR-description-as-proposal, escalation rule, global CH numbering)
- [x] ADR-016 — decision log (litmus test, log format, amended retention convention, demotion mechanics)
- [x] `docs/decisions/log.md` created with frontmatter (`type: log`, `status: active`)
- [x] `log` type added to `statusVocab` in `internal/corpus/rules.go` AND the status table in `docs/README.md` (the code comment mandates changing them together)
- [x] Unit test: `log` status vocabulary enforced (positive: `active` passes; negative: other status fails)
- [x] Demote ADR-003 → log row (delete file, repoint CAP-002 design.md link + prose)
- [x] Demote ADR-004 → log row (delete file, repoint QS-001 link + prose, fix ADR-013 prose citation)
- [x] `docs/decisions/README.md` — ADR vs log split explained, retention wording amended, index updated (drop ADR-003/004, add log.md/ADR-015/ADR-016)
- [x] `.agents/skills/clue-delta/skill.md` — light-tier section + escalation rule; ADR-vs-log routing in digest step
- [x] `.agents/skills/clue-verify/skill.md` — first checklist item: is this change correctly light / correctly full?
- [x] `AGENTS.md` rule 1 mentions the light tier
- [x] Digest: CHANGELOG `[Unreleased]` entry, delete `/changes/CH-010-light-tier-decision-log/`
- [x] Verify: `go test ./...` green, `go run ./cmd/clue validate` green, clue-verify checklist walked
