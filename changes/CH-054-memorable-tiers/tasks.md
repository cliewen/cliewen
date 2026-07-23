---
id: CH-054-tasks
type: tasks
status: open
links: [CH-054]
title: CH-054 ordered tasks
---

# CH-054 — Tasks

- [ ] Rewrite `internal/skills/source/shared/change-tiers.md.tmpl` as three rules plus two guards, preserving every obligation of the current text
- [ ] Check the rewrite against the old text clause by clause and confirm no work moves tier
- [ ] Align `AGENTS.md` classification prose and the light-change clause in rule 1
- [ ] Align `internal/scaffold/templates/AGENTS.md` the same way
- [ ] Align `guide/change-loop.md` and `guide/methodology.md` tier prose
- [ ] Run `go generate ./internal/skills` and confirm the regenerated skills carry version 0.6.0
- [ ] Add the decision-log row naming `change-tiers.md.tmpl` as the carrier
- [ ] Add the CHANGELOG entry under `[Unreleased]`
- [ ] Verify: `go build ./... && go test ./...`, `go run ./cmd/clue validate`, `npm run guide:build`, `go generate ./internal/skills && git diff --exit-code`
