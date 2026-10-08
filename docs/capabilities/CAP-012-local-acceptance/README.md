---
id: CAP-012
type: capability
status: active
links: [G-016, PDR-067]
title: Human-controlled local acceptance
goal: G-016
---

# Local acceptance

An adopter that explicitly selects local acceptance can accept a verified tracked branch without a pull request. A human reviews the acceptance brief and runs `clue accept`; the resulting merge preserves the candidate tree, proposal history, and complete brief. PR acceptance remains the default, and identity allocation remains independent of the acceptance mechanism.

The local boundary is procedural: a prompt and a Git committer identity do not authenticate human presence. See [design](design.md) for the opt-in, record, and recovery contract and [criteria](criteria.md) for executable evidence.
