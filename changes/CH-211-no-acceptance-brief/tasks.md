---
id: CH-211-tasks
type: tasks
status: open
links: [CH-211]
title: CH-211 tasks
---

# CH-211 tasks

- [x] Retire the criteria whose meaning names the brief and add their successors: AC-219 → AC-242, AC-237 → AC-243, AC-238 → AC-244 (CAP-012); AC-042 → AC-245, AC-048 → AC-246, AC-239 → AC-247, AC-063 → AC-248, AC-130 → AC-249 (CAP-006); AC-054 → AC-250 (CAP-003); AC-218 → AC-251 (CAP-010).
- [x] `internal/accept` and `cmd/clue/accept.go`: no brief; change found from the one proposal introduced after the base; merge message from the proposal's title and its `## Human evidence` section; approval record bound to candidate, base and decision; `--brief` and a `brief-sha256` field refused with a plain message (AC-242, AC-243, AC-244, AC-246).
- [x] `internal/corpus/actests.go`: the Human-class comment names the proposal's human evidence section (AC-246).
- [x] Shipped and local workflows: remove the PR-body brief check; tests for AC-245 and AC-248.
- [x] Scaffold: delete the brief template, drop `brief-sha256` from the approval template, remove the brief from the tracked pull-request template and its checklist; tests for AC-247.
- [x] Skills sources (`clue-delta`, `clue-verify`, `clue-extract`, `clue-upgrade`, shared fragments) rewritten without the brief, then `go generate ./internal/skills`; tests for AC-246, AC-249, AC-250, AC-251.
- [x] Decision record PDR-073 superseding the brief parts of PDR-017, PDR-067, PDR-072 and ADR-033; C-023 moved to the proposal; affected decisions, constraints, goals, architecture, design and capability documents updated.
- [x] Public guide pages and `AGENTS.md` (repository and scaffolded) updated; `CHANGELOG.md` `[Unreleased]` entry with the upgrade action.
- [x] Run the **Verify Locally** block from `CONTRIBUTING.md`.
