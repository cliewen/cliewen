# Tasks

- [ ] Record PDR-043: an upgrade routes as simple work, with its escalation condition, carrier inventory, and rejected alternatives.
- [ ] Retire AC-081 with a tombstone and mint AC-142 for the upgrade skill's contract including its route, serving PDR-043.
- [ ] State the route and escalation condition in `internal/skills/source/skills/clue-upgrade.md.tmpl` and regenerate the skill trees with `go generate ./internal/skills`.
- [ ] Retag the AC-081 evidence to AC-142 and extend it with focused positive and negative coverage of the route statement.
- [ ] Move the remaining live carriers: `docs/architecture/skills.md`, CAP-004's design, `guide/operations.md`, and CAP-001's design where it states the upgrade's shape.
- [ ] Write the `[Unreleased]` changelog entry for the adopter-visible skill-text change.
- [ ] Run the full local verification block, complete the digest, and prepare the reviewed PR handoff.
