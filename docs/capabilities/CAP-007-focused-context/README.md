---
id: CAP-007
type: capability
status: active
links: [G-001]
title: Focused corpus context
goal: G-001
---

# CAP-007 — Focused corpus context

An engineer or agent can start from one known corpus identity and read only the durable artifacts that identity depends on. `clue context <id> [path]` resolves an artifact ID, acceptance-criterion ID, or milestone ID, then emits the owning artifact and the artifacts it links, out to a stated number of hops. It replaces manual path search and the parked `clue locate` idea with one deterministic read path.

The slice is intentionally directional. Following reverse links from a goal would pull most of a mature corpus into one result and recreate the mandatory full-corpus read this capability exists to remove. A caller starts from the artifact closest to the task — often a criterion, capability, decision, or milestone — and receives its declared durable dependencies.

The slice is also bounded, and defaults to the root and what it links to directly. A corpus densifies as campaigns close, and following every outgoing edge to exhaustion then returns most of the repository from an ordinary starting point, which is the read this capability exists to remove ([ADR-056](../../decisions/ADR-056-bounded-context-slice.md)). The bound is a default, never a limit: `--depth` widens it, `--depth=all` restores the whole closure, and what the bound held back is named rather than silently dropped, so widening stays a judgement the reader makes on evidence.

`active`: AC-133 is implemented with focused positive and negative unit evidence and command-level coverage. AC-053 is retired: its assertions describe the unbounded slice and survive as what `--depth=all` guarantees.
