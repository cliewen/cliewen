---
id: CH-007-tasks
type: tasks
status: open
links: [CH-007]
title: Task breakdown for CH-007
---

# Tasks — CH-007

- [x] Write CAP-004 criteria.md (AC-019…AC-022) — the anchors the tests trace to
- [x] `clue version` / `clue --version` reports the stamped version; `dev` when unstamped (AC-019)
  - [x] `runVersion(io.Writer)` in `cmd/clue/main.go`; `version` package var defaults to `dev`
  - [x] `TestAC019` positive + negative pair
- [x] `checkSkillVersions` rule in `internal/corpus` (AC-020, AC-021, AC-022)
  - [x] `Options.Version`; wire the rule into `Validate`; `runValidate` passes the binary version
  - [x] `TestAC020/AC021/AC022` positive + negative pairs
- [x] Stamp every `.agents/skills/*/skill.md` with `version: 0.1.0` frontmatter
- [x] `.github/workflows/release.yml`: `v*` tag → cross-platform stamped binaries → GitHub release
  - [x] `TestSanity_ReleaseWorkflowIsCrossPlatform` guards the workflow shape
- [x] CAP-004 README.md + design.md
- [x] ADR-011 (version stamping, per-skill granularity, drift-is-failure)
- [x] README install docs (`go install`, `gh release download`) + status refresh
- [x] Digest: index blocks (capabilities, decisions), plan bookkeeping (M-004 → wip), delete `changes/CH-007-clue-ships/`
- [x] clue-verify checklist, then PR
