---
id: G-013
type: goal
status: accepted
links: [VIS-001]
title: Parallel contributors receive identities that cannot collide
---

# G-013 — Parallel contributors receive identities that cannot collide

## Who wants it

Developers and coding agents working concurrently in branches, clones, and worktrees of the same repository.

## Why

A repository-local counter gives every caller the same answer when they start from the same accepted commit. Git can merge identical ledger edits without warning, leaving two changes to discover their shared identity only after integration. A team cannot rely on a verifiable thread whose identities are safe only when one person works at a time.

## Success looks like

- Concurrent supported allocations return distinct sequential IDs.
- A failed coordination attempt never silently falls back to a checkout-local claim.
- A reservation remains unavailable even when the work that requested it is abandoned.
- Contributors without allocator write access can receive a preallocated ID and synchronize its reservation.
