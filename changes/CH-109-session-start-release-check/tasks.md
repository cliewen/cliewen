---
id: CH-109-tasks
type: tasks
status: active
links: [CH-109]
title: Tasks — a coding agent learns a release is available when its session starts
---

# Tasks — CH-109

- [ ] Write PDR-023: executable vendor configuration is a distinct category, bound to the entry point's ownership limits plus the limits execution adds
- [ ] Add AC-083 (init emits the session-start release check), AC-084 (migrate reports a settings file that never runs it), AC-085 (`Human`: the line arrives in a real session) to CAP-001's criteria
- [ ] Emit `.claude/settings.json` from the scaffold — new template, `claude/` target mapping, existing file skipped byte-for-byte (AC-083)
- [ ] Report a settings file that never runs the check, and an absent one, as MIG-006 notices that repair nothing (AC-084)
- [ ] Prove AC-083 and AC-084 in both directions with focused Unit tests
- [ ] Materialize this repository's own `.claude/settings.json` so the change is dogfooded (AC-085)
- [ ] Update CAP-001 README and design, CAP-004 design, `guide/operations.md`, `guide/getting-started.md`, and the migration inventory to state the same contract
- [ ] Add the `[Unreleased]` changelog entry
- [ ] Ask the human to observe the line in a real session and record the verdict in the acceptance brief (AC-085)
