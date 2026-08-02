---
id: CH-107-tasks
type: change
status: open
links: [CH-107]
title: Tasks for CH-107
---

# CH-107 — tasks

- [x] Record the decision: an ADR stating why a release check may reach the network while the judge never may, and what separates a reporting command from a background updater (no repository write with or without flags, no self-replacement, never a required check) — ADR-042
- [x] Add the acceptance criteria to CAP-004 before implementation
  - [x] AC-075 — a newer release is reported with the installation route for the platform the command is running on, and the coordinated preview-and-apply sequence
  - [x] AC-076 — quiet mode prints one line when behind and nothing when current
  - [x] AC-077 — offline, timeout, rate limit, and an unrecognized response are all silent and exit 0
  - [x] AC-078 — the answer is cached outside the repository with a bounded lifetime, and an unreadable cache is absence rather than error
  - [x] AC-079 — the drift message names the command that resolves it and the way to pin instead
- [x] Implement the release check with the network and the platform branch both injected (AC-075, AC-076, AC-077, AC-078)
- [ ] Wire `clue latest [--quiet] [--timeout=<duration>]` into the CLI and its usage text (AC-075, AC-076)
- [x] Extend the drift message with the resolving command and the pinning route (AC-079)
- [x] Prove all five criteria with focused positive and negative unit evidence, no test reaching a live service
- [ ] State the same contract in `guide/operations.md`, `guide/getting-started.md`, and CAP-004's README and design
- [ ] Add the `[Unreleased]` changelog entry
- [ ] Run `clue-verify`, including its agentic review loop, and open the ready pull request
