---
id: CH-020-tasks
type: tasks
status: open
links: []
title: CH-020 task breakdown
---

# Tasks

- [x] Criteria first: revise `CAP-001/criteria.md` — retire AC-001 (tombstone, `@retired`), add testable ACs for the mechanical path (AC-024 idempotent re-run, AC-025 no-overwrite; AC-002/AC-003 kept), move the 30-minute end-to-end promise to a quality scenario (QS-002, task below)
- [x] Build the template tree from the verified fixture: docs taxonomy READMEs + decision log, AGENTS.md routing hub, CI workflow template (model2diagram's vendored pattern, generalized) — at `internal/scaffold/templates/` (not `/.cliewen/`: the Go toolchain ignores dot-directories, go:embed cannot reach them; recorded as a log row)
- [x] Implement `clue init` in Go: embed templates (`go:embed`), emit taxonomy + AGENTS.md + skills (`.agents/skills/` + `.claude/skills/` mirror) + CI template; idempotent re-run regenerates only index blocks; never overwrite existing files, report skips
- [x] Test pairs for every new AC (positive + negative), driving the scaffolder in a temp repo; purposes per the taxonomy (`AC<digits>` prefix in Go); plus a Sanity test holding the embedded skills byte-identical to `.agents/skills/`
- [x] Quickstart in README: prerequisites (git, clue, gh-as-convenient-path), install → init → first loop → green validate; skills linked at point of need
- [x] Document-type introduction in the generated `docs/README.md` + folder READMEs (plain-language what/when for ADR, PDR, log row, capability folder, architecture, constraints, quality)
- [x] Digest prep: CAP-001 README/criteria/design to `active`, design.md elaborated from stub (template source, embed decision, idempotence, skip behavior), architecture docs updated (`init` joins the CLI surface), decisions recorded (ADR-018; log rows for the template location and the AC-001 split), plan bookkeeping M-005, changelog entry
