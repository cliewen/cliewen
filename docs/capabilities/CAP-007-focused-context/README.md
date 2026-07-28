---
id: CAP-007
type: capability
status: active
links: [G-001]
title: Focused corpus context
---

# CAP-007 — Focused corpus context

An engineer or agent can start from one known corpus identity and read only the durable artifacts that identity depends on. `clue context <id> [path]` resolves an artifact ID, acceptance-criterion ID, or milestone ID, then emits the owning artifact and the transitive closure of its outgoing `links` edges. It replaces manual path search and the parked `clue locate` idea with one deterministic read path.

The slice is intentionally directional. Following reverse links from a goal would pull most of a mature corpus into one result and recreate the mandatory full-corpus read this capability exists to remove. A caller starts from the artifact closest to the task — often a criterion, capability, decision, or milestone — and receives its declared durable dependencies.

`active`: AC-053 is implemented with focused positive and negative unit evidence and command-level coverage.
