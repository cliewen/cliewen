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
- [x] Ask the human to observe a real session reading the hub and running the check (AC-085) — **run 2026-08-03 and failed.** The hub reached the session's context and the agent read past the instruction into the requested work; a second finding came with it, that the installed binary predated `clue latest`, so the instruction would have produced an unknown-command error rather than the promised one line or silence. `open-questions.md` Q2 holds both.
- [x] Write Q2 and revise PDR-023 in place: the tool carries the notice, the hub carries the instruction (unmerged and unaccepted, so revised rather than superseded; renamed to match)
- [x] Add the update notifier — `release.Notice` and `release.QuietLine` beside the existing check, gated in `main` by command, terminal, `CI`, and `CLUE_NO_UPDATE_NOTIFIER`, printed to standard error after the command it rode in on
- [x] Restate AC-085 as an observable criterion about a session learning it is behind by either channel, and add AC-086 and AC-087 to CAP-004 for the notice and its gate
- [x] Prove AC-086 and AC-087 in both directions, including that the two callers share one phrasing and that the gate still lets every workflow command through
- [x] Sharpen the hub's trigger in both carriers, name the notifier beside it, and read an unknown-command error as the answer
- [x] Revise M-047 a second time and record why prose alone did not carry it; update CAP-001 README and design, CAP-004 design, both guides, and the `[Unreleased]` entry
- [x] Answer Q3 — the notifier's terminal gate excluded every coding-agent session, so the mechanism reached nobody it was built for
- [x] Remove the terminal condition from the gate, revise PDR-023 in place, and restate AC-086 and AC-087; prove the notice survives a pipe, that standard output and the exit code do not move, and that no repository file is written
- [x] Repair review findings for AC-086 and AC-087: skip unusable ambient requests for incomparable stamps, and make the built-binary test exercise a real notice on standard error across supported cache locations
- [x] Name the notice and its opt-out in `clue help` — the switch was documented everywhere except the tool a user asks first
- [x] Bound what a check that cannot answer costs: only a successful fetch writes the cache today, so every workflow command in an offline or blackholed session pays the ambient budget again for the same non-answer, and "cached for a day so repeating it costs no request" is true only when the release list answers. Remember the failure too, with a short lifetime of its own, and decide in the same step what that does to ADR-042's "a cache that cannot be read or written is simply ignored" — this changes what the cache file means, so it needs a criterion of its own, not a quiet repair under AC-086
- [x] Ask the human to observe a fresh session learning the repository is behind, and record the verdict in the acceptance brief (AC-085 retest — not this session, which knows the criterion is under test) — **run 2026-08-04 and passed.** A fresh Codex session in the Robocode Tank Royale adopter repository, asked only "hvad er status på CH-012?", ran `clue latest --quiet` as its first action unprompted, reported `clue 0.11.2 is behind 0.12.0`, and stopped to ask whether to upgrade now or continue — the routing the hub prescribes, and the human said nothing about releases. The binary was built from this branch and stamped 0.11.2 so it carried the notifier and was genuinely behind the real published 0.12.0; the hub line was added by hand, since MIG-006 reports it rather than repairing it. A 0.12.0-stamped build of the same code is silent in the same repository, so the absence of the line is informative. One observation, not a defect: the hub routes to `clue-upgrade`, which that repository's 0.11.2-era skills predate, so the link dangled — a state the test rig created and `clue migrate` cannot leave behind, because the migration that reports the missing line is the one that refreshes the skills. The session reported the gap and asked rather than inventing an upgrade
