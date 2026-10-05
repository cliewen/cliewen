---
id: CH-201-tasks
type: tasks
status: open
links: [CH-201]
title: Scenario set tasks
---

# Tasks

- [x] Generalise the harness: a scenario names its prompt, fixture, the obligation it exercises, the failure that obligation prevents, its observed fields and its checks; the signature is built from the scenario's own fields. Keep `routing` working and its recorded results re-checkable. Cover the generalisation with Go tests.
- [x] Add method variants: a named removal applied to the container's copy after setup, never to the repository, recorded with its hash in the conditions. Add the `no-routing` removal.
- [x] Add the `routing-code`, `upgrade` and `brownfield` scenarios with fixtures and checks, and the shim and image change the upgrade scenario needs.
- [x] Run the cheapest test first: `routing-code` with and without the routing obligation, read all transcripts of both arms, and revise a check wherever a transcript contradicts it.
- [x] Run each scenario on both agents, read the transcripts against the checks, and write AN-031 with what each scenario showed, what each check could not see, and whether the provoking scenario provoked.
- [x] Declare the revision of P-024 (M-105 reduced to four scenarios, M-108 added), and write M-105's status and evidence.
- [x] Regenerate indexes and the evidence export, run the full verification, and state the documentation impact in the pull-request handoff.
