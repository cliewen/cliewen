---
id: CH-009-tasks
type: tasks
status: open
links: [CH-009]
title: Task breakdown for CH-009
---

# Tasks — CH-009

- [x] ADR-013: ships-generic vs repo-local (skills verbatim, README prose as init template source, index blocks per-repo, AGENTS.md local; extends-never-overrides; refines ADR-008 to mapping files)
- [x] `clue-extract`: move OpenSpec mapping verbatim to `mappings/openspec.md`; generic "Source mappings" pointer in `skill.md`; contract item 7 gains AGENTS.md reconciliation
- [x] ADR-014: PR approval promotes the PR's ADRs (agent does the clerical flip citing approver/date/PR)
  - [x] Retroactive: ADR-011 → `verified` (PR #6), ADR-012 → `verified` (PR #7)
- [x] Decisions README: ADR style rule (timeless, problem-not-episode) + promotion mechanics; index += ADR-013, ADR-014
- [x] ADR-011: replace the "originally deferred / later closed" parenthetical with a timeless statement (meaning unchanged)
- [x] Digest: CHANGELOG `[Unreleased]` entry (skills are repo-agnostic), delete `changes/CH-009-generic-skills/` (deletion is the digest commit itself)
- [x] clue-verify checklist (validate OK 37 artifacts, tests green, skills sweep clean, mapping move verified verbatim against git history), then PR
