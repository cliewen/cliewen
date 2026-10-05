---
id: CH-199-tasks
type: tasks
status: open
links: [CH-199]
title: Pinned scenario run tasks
---

# Tasks

- [ ] Record IDR-009: the scenario-trials harness is Go under `tools/scenario-trials`, repository tooling that is not shipped and not started by CI.
- [ ] Build `tools/scenario-trials`: pinned container image, `clue` built from a named commit, the routing scenario, a Claude Code run in the container, conditions recorded for every run, deterministic checks over the event stream, and a spread summary. Cover the pure logic with Go tests.
- [ ] Run the scenario five times on Claude Code, read all five transcripts against the checks, and write the findings document with the spread and how often the checks and the reading agreed.
- [ ] Write M-103's evidence in P-024, and revise M-105 and M-106 through a declared plan revision only if the spread says so.
- [ ] Regenerate indexes and the evidence export, run the full verification, and state the documentation impact in the pull-request handoff.
