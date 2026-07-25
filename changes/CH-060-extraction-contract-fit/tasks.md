---
id: CH-060-tasks
type: tasks
status: open
links: [CH-060]
title: Tasks for CH-060
---

# Tasks

- [x] Generalize the routing contract item from `AGENTS.md` alone to assistant entry points as a class, keeping `AGENTS.md` as the flagship example.
- [x] Add a contract item stating that a criterion born `draft` is the sanctioned phasing lever for a corpus too large to tag-test in one extraction change.
- [x] Add deterministic ID-minting guidance to the ID-survival contract item for source requirements with no stable existing ID.
- [x] Regenerate the skills (`go generate ./internal/skills`) and confirm no drift.
- [x] Run `go test ./...` and confirm green.
- [x] Run the strict guide build (`npm run guide:build`) and confirm green.
- [x] Update `docs/plans/P-006-first-adoption.md`: mark M-022 `done` with an evidence cell citing CH-060 and the concrete files/tests.
- [x] Add a `CHANGELOG.md` entry under `[Unreleased]` describing the contract change from a user's perspective.
- [x] Record the decision as a `docs/decisions/log.md` row citing CH-060.
- [ ] Run `clue-verify`'s local checks and agentic review loop (in-context fallback, disclosed) before opening the PR.
