---
id: CH-109-tasks
type: tasks
status: open
links: [CH-109]
title: Tasks — a coding agent learns a release is available when its session starts
---

# Tasks — CH-109

- [-] Emit vendor configuration that runs the release check at session start, and report its absence — rejected by the human: shipping executable configuration for one assistant violates Cliewen's agent-agnosticism. The implementation was written and reverted; `open-questions.md` Q1 holds the finding.
- [x] Write PDR-023: the routing hub is the session-start channel, and no vendor configuration is ever emitted
- [x] Revise M-047 in P-010 and add the campaign's standing no-vendor-configuration boundary
- [x] Add AC-083 (init emits a hub that asks), AC-084 (migrate reports a hub that never asks), AC-085 (`Human`: a real session reads the hub and runs the check) to CAP-001's criteria
- [x] Put the instruction in the emitted `AGENTS.md` template and in this repository's own hub (AC-083)
- [x] Report a hub that never asks, and an absent hub, as MIG-006 notices that repair nothing (AC-084)
- [x] Prove AC-083 and AC-084 in both directions with focused Unit tests, including that no vendor configuration is emitted
- [x] Update CAP-001 README and design, CAP-004 design, `guide/operations.md`, `guide/getting-started.md`, and the migration inventory to state the same contract
- [x] Add the `[Unreleased]` changelog entry and its migration note
- [ ] Ask the human to observe a real session reading the hub and running the check, and record the verdict in the acceptance brief (AC-085)
