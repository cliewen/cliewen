---
id: G-029
type: goal
status: proposed
links: [G-019, VIS-001]
title: The acceptance check and the change-loop skill agree on the proposal's identity
---

# G-029 — The acceptance check and the change-loop skill agree on the proposal's identity

**Who wants it:** every agent and maintainer taking a tracked change to local acceptance (2026-10-20), found while measuring the readable-handoff rule on repeated runs of the `tracked-brief` scenario.

**Why:** the change-loop skill and the validator tell a new change to give its proposal the identity `CH-xxx-proposal`, like the tasks and open-questions files. The local acceptance check, `clue accept`, looks for a proposal whose identity is exactly the change's own number, `CH-xxx`, and otherwise stops with "candidate history lacks the proposal introduced after its base". An agent that follows the skill therefore reaches the end of a tracked change and cannot pass the check. In the scenario runs this happened in every run of the current skills, three out of three: the agent did not find the cause and handed the failure to the person, and in one run it also tried changing the brief. The maintainer of this repository hit the same stop on two changes in a row and worked around it by writing the proposal's identity as `CH-xxx`.

**Success looks like:**

- The skill, the validator and the acceptance check name the same identity for a proposal, so a change written from the documentation passes the check on the first try.
- The acceptance check's message says which identity it looked for, so a failure explains itself.
