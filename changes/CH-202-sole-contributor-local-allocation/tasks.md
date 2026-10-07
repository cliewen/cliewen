---
id: CH-202-tasks
type: tasks
status: open
links: [CH-202]
title: Sole-contributor local allocation tasks
---

# Tasks

- [ ] Add a CAP-010 criterion, `Test-type: Unit`, for the changed step: on the local-allocation warning the agent reads the remote with plain Git for other `ch-*` branches and change claims on the allocator branch, proceeds on serialized allocation when none exist, records what it read in the brief's ledger note, and stops where another contributor's branch exists or the remote cannot be read. Positive evidence asserts the generated change loop carries those obligations; negative evidence asserts it no longer tells the agent to stop unconditionally and still names the two stop cases.
- [ ] Edit the canonical source `internal/skills/source/skills/clue-delta.md.tmpl` for step 1, run `go generate ./internal/skills`, and inspect the generated and scaffolded `change-loop.md` as an adopter receives them; add the ledger note line to the acceptance brief guidance.
- [ ] Update every live carrier of the rule: `guide/first-change.md` (the line saying teams must serialize allocation), CAP-010's README and design, and an `[Unreleased]` entry in `CHANGELOG.md` naming the adopter-facing surface (generated `clue-delta` text).
- [ ] Record the decision as a PDR (inferred, agent-authored), routed by subject, stating the riskiest assumption and the stop cases; link it from CAP-010.
- [ ] Run the cheapest test from the proposal's challenge: the changed rule on the AN-036 fixture, five runs each on Claude Code and Codex, and on a fixture variant with another clone's `ch-*` branch pushed to the bare remote, three runs each. Fix the fixture's commit message so it does not name the variant. Record the result, and the rule's stated limit, in a new analysis; stop and revise if an agent proceeds past another contributor's branch.
- [ ] Digest: update permanent `/docs`, regenerate indexes and the evidence export, delete this workspace, run the full verification, and state the documentation impact in the pull-request handoff.
