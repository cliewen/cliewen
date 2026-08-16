---
id: CH-155-tasks
type: tasks
status: open
links: [CH-155]
title: Tasks for CH-155
---

# Tasks

- [x] Record PDR-043: an upgrade routes as simple work, with its escalation condition, carrier inventory, and rejected alternatives.
- [x] Retire AC-081 with a tombstone and mint AC-142 for the upgrade skill's contract including its route, serving PDR-043.
- [x] State the route and escalation condition in `internal/skills/source/skills/clue-upgrade.md.tmpl` and regenerate the skill trees with `go generate ./internal/skills`.
- [x] Retag the AC-081 evidence to AC-142 and extend it with focused positive and negative coverage of the route statement.
- [x] Move the remaining live carriers: `docs/architecture/skills.md`, CAP-004's design, and `guide/operations.md`; CAP-001's design states the hub's routing rather than the upgrade's shape and needed no change.
- [x] Write the `[Unreleased]` changelog entry for the adopter-visible skill-text change.
- [ ] Run the full local verification block, complete the digest, and prepare the reviewed PR handoff.
