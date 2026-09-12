---
id: CH-181
type: change
status: open
links: [G-012, CAP-002, CAP-011, AC-190, AC-191, AC-192]
title: Fresh contexts receive useful repository orientation
---

# CH-181 — Fresh contexts receive useful repository orientation

This plan-less change makes a fresh agent context useful between campaigns as well as during one. Today `clue next` can name active and draft milestones, but an adopter with no active plan receives only an empty result even when an open change or proposed goals provide the relevant choices. The scaffolded routing hub also waits for the user to ask what is next, so a newly started agent does not establish where the repository stands.

Add a dedicated orientation capability serving G-012. `clue next` will derive a stable ladder from repository artifacts: open change workspaces, active `doing` and `todo` milestones, draft milestones, then proposed goals. Its existing interface remains read-only; `--all` exposes all categories, while the default recommends the first candidate at the strongest available level. When the corpus offers none, it says that a new goal is the remaining direction-setting path.

The scaffolded and repository routing hubs will tell an agent to run `clue next --all` once after the mandatory release check in a fresh context. The agent mentions the result only when the request leaves direction open or existing work materially affects it, reads bounded context before recommending a candidate, and never treats a proposal as authorization. No local focus, Git-user namespace, assistant hook, or state outside the corpus is introduced.

The change is plan-less because no campaign is active. It proceeds under VIS-001 and serves G-012's goal that an arriving agent can orient from a bounded repository read.
