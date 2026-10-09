---
id: CAP-012
type: capability
status: active
links: [G-016, PDR-067]
title: Human-controlled local acceptance
goal: G-016
---

# Local acceptance

An adopter whose accepted policy selects local acceptance can accept a verified tracked branch without a pull request. A human reviews the acceptance brief and runs `clue accept`; the resulting merge preserves the candidate tree, proposal history, and complete brief. Fresh adoption defaults to local; PR is actively selected. Existing workflows are preserved during upgrade under PDR-068, and identity allocation remains independent of acceptance.

The local boundary is procedural: a prompt and a Git committer identity do not authenticate human presence. See [design](design.md) for policy selection, record, and recovery contract and [criteria](criteria.md) for executable evidence.
