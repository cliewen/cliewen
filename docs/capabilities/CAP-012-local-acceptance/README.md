---
id: CAP-012
type: capability
status: active
links: [G-016, PDR-067, PDR-069, PDR-070]
title: Human-controlled local acceptance
goal: G-016
---

# Local acceptance

A repository whose accepted policy selects local acceptance can accept a verified tracked branch without a pull request. A human reviews the acceptance brief and runs `clue accept`; the resulting merge preserves the candidate tree, proposal history, and complete brief. Acceptance defaults to local on main when policy is absent; PR is actively selected under PDR-069, and identity allocation remains independent of acceptance.

The local boundary is procedural: a prompt and a Git committer identity do not authenticate human presence. See [design](design.md) for policy selection, record, and recovery contract and [criteria](criteria.md) for executable evidence.
