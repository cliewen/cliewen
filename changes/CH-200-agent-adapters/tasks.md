---
id: CH-200-tasks
type: tasks
status: open
links: [CH-200]
title: Agent adapters tasks
---

# Tasks

- [x] Introduce a vendor-neutral transcript and an adapter interface, move the Claude Code specifics into one file, and keep every check, the conditions record and the summary on the neutral transcript. Cover the interface with Go tests.
- [x] Add the Codex adapter: login from a directory outside the repository, the command that runs it with its own sandbox off, a mapping of its event stream to the neutral transcript, and a probe for its version and model. Install both agents in the image with pinned versions.
- [x] Run Codex once and read its transcript beside the neutral one to check the order of text and edits, then run the scenario on both agents and read the transcripts against the checks.
- [x] Record IDR-010 (adapters over one interface, a terminal-driving orchestrator not used) and AN-030 (the runs and what the adapters could and could not see of each agent's host configuration).
- [x] Declare the revision of M-107's repetition count in P-024, and write M-104's status and evidence.
- [x] Regenerate indexes and the evidence export, run the full verification, and state the documentation impact in the pull-request handoff.
