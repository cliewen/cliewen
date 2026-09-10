---
id: CAP-010
type: capability
status: active
links: [G-013]
title: Team-safe identity allocation
goal: G-013
---

# CAP-010 — Team-safe identity allocation

## What

Cliewen allocates readable sequential identities safely across concurrent clones and worktrees by serializing permanent claims through a dedicated Git remote branch. A checked-in append-only ledger carries identity lifecycle state and merges independent events without making the deterministic judge depend on network access.

Repositories opt into coordination explicitly. Checkout-local allocation remains available for single-user or deliberately serialized work and warns that it is not safe for concurrent callers. Batch reservation and read-only synchronization cover assigned identities, fork contributors, and planned offline work. `clue id retire <id>` records the terminal transition after an artifact is deleted, so no contributor has to edit the append-only ledger by hand.

## Why

Serves [G-013](../../goals/G-013-team-safe-identities.md): two callers starting from the same committed counter must not receive the same identity merely because their working trees cannot see each other.

Acceptance criteria: [criteria.md](criteria.md) · implementation notes: [design.md](design.md).

## Status note

`active`: the event ledger, Git compare-and-swap protocol, batch and synchronization commands, migration, and fail-closed behavior are implemented and covered by focused unit and disposable-Git integration evidence.
