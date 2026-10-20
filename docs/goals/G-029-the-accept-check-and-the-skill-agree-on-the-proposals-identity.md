---
id: G-029
type: goal
status: proposed
links: [G-019, VIS-001]
title: The acceptance check and the change-loop skill agree on the proposal's identity
---

# G-029 — The acceptance check and the change-loop skill agree on the proposal's identity

A tracked change's proposal should carry one identity that both the change-loop skill and the local acceptance check expect, so that following the documentation reaches acceptance.

**Who wants it:** every agent and maintainer taking a tracked change to local acceptance (2026-10-20), found while measuring the readable-handoff rule on repeated runs of the `tracked-brief` scenario.

**Why:** the change-loop skill tells a new change to give its proposal the identity `CH-xxx-proposal`, like the tasks and open-questions files. The validator is silent about the proposal's identity: its hint covers only the tasks and open-questions files. The local acceptance check, `clue accept`, looks for a proposal whose identity is exactly the change's own number, `CH-xxx`, and otherwise stops with "candidate history lacks the proposal introduced after its base". An agent that follows the skill therefore reaches the end of a tracked change and cannot pass the check. In the scenario runs it happened in both runs made after the readable-handoff wording was added (a baseline run made before it did not reach the check): the agent did not find the cause and handed the failure to the person, and in one run it also tried changing the brief's status and change value. In this repository, change CH-208 hit the same stop and was worked around by rewriting the local history so the proposal's identity read `CH-208`; the change before it, CH-205, had written the proposal's identity that way from the start.

**Success looks like:**

- The skill and the acceptance check name the same identity for a proposal, and the validator says which identity a proposal takes, so a change written from the documentation passes the check on the first try.
- The acceptance check's message says which identity it looked for, so a failure explains itself.
