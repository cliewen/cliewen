---
id: CH-195-tasks
type: tasks
status: open
links: [CH-195]
title: CH-195 tasks
---

# CH-195 tasks

- [ ] Record the PDR naming the routes `direct` and `tracked`, the unchanged route meaning, and the dual-spelling trailer rule
- [ ] Retire AC-139 and AC-142 with tombstones and mint successors using the new route names
- [ ] Add a criterion that the shipped validation workflow accepts an override in either spelling, with positive and negative evidence
- [ ] Add a criterion and a `clue migrate` notice for an adopter hub that still names the old routes, with positive and negative evidence
- [ ] Update canonical skill sources under `internal/skills/source/` and run `go generate ./internal/skills`
- [ ] Update `internal/scaffold/templates/` (AGENTS.md, PR template) and the shipped `clue-validation.yml`
- [ ] Update the local CI scope script and its tests to accept both spellings
- [ ] Update this repository's `AGENTS.md`, `CONTRIBUTING.md`, `.github/pull_request_template.md`, and live decisions, constraints, criteria, architecture, and design overviews
- [ ] Update `/guide` prose with the humanizer skill, leading each route's first mention with the question that decides it
- [ ] Retag and update tests that assert route wording; regenerate the evidence manifest
