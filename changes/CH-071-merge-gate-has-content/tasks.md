---
id: CH-071-tasks
type: tasks
status: active
links: [P-007]
title: Tasks for CH-071
---

# Tasks

- [x] Branch `ch-071-merge-gate-has-content` from `main` tip
- [x] Write proposal.md
- [ ] **Pause: draft PDR-017 and get human confirmation before implementing** (opt-in spec-first pause, exercised on this change itself)
- [ ] Record PDR-017 (docs/decisions/), update docs/decisions/README.md index
- [ ] Add shared skill fragment (system-of-record) under internal/skills/source/shared/
- [ ] Update internal/skills/source/skills/clue-delta.md.tmpl — opt-in pause after Propose, include new shared fragment
- [ ] Update internal/skills/source/skills/clue-verify.md.tmpl — per-criterion scenario-resolution step in the review loop, acceptance-brief handoff content
- [ ] Update .github/pull_request_template.md — acceptance-brief section at top, competence-heuristic warning, one-screen-cap prose, scenario-resolution table
- [ ] Add deterministic unfilled-brief check (CI workflow step/Go check) — serves new AC
- [ ] Mirror the check into the scaffolded wall (internal/scaffold/templates/github/workflows/clue.yml)
- [ ] Add acceptance criteria (AC-xxx) to the capability that owns the change-loop/PR contract, with positive and negative tests
- [ ] `go generate ./internal/skills` — confirm no drift
- [ ] Update guide pages describing the merge gate as empty / lacking scenario-resolution (guide/change-loop.md, guide/design.md, guide/methodology.md)
- [ ] CHANGELOG.md entry under [Unreleased]
- [ ] Update docs/plans/P-007-core-hardening.md M-024 row to done with evidence
- [ ] `go build ./...`, `go test ./...`, `go run ./cmd/clue validate --forbid-changes`, `git diff --check`
- [ ] Digest: delete /changes/CH-071-merge-gate-has-content/
- [ ] clue-verify agentic review loop, then open ready PR
