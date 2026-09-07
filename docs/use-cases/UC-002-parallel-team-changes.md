---
id: UC-002
type: use-case
status: active
links: [G-013, CAP-010, CAP-002]
title: A team starts independent Cliewen changes in parallel
---

# UC-002 — A team starts independent Cliewen changes in parallel

## Actors

The **contributors or coding agents** starting independent work, the **repository maintainer** who enables and protects coordination, and the **Git remote** that serializes claims without deciding what the changes mean.

## Trigger

Several contributors start from the same accepted repository state and each needs a new Cliewen identity before their branches can diverge.

## Preconditions

The repository's version-two identity ledger is committed in Git-coordinated mode. The configured remote carries the protected `clue/id-allocator` branch and permits ordinary fast-forward pushes while rejecting force-push and deletion.

## Main flow

1. Each contributor runs `clue id next CH` from their own clone or worktree.
2. Each command reads the same remote claim history, proposes a distinct next claim, and attempts an ordinary fast-forward push.
3. One concurrent proposal wins. Every loser fetches the new head, chooses the next unclaimed number, and retries within the bounded timeout.
4. After a successful remote claim, the command appends every missing reservation to the contributor's checked-in ledger and prints only the IDs this invocation obtained.
5. The contributor creates the change artifact and promotes its reservation locally. Parallel branches merge their one-line ledger events through Git's union driver, and `clue validate` judges the resulting checkout without contacting the remote.

## Alternative and failure flows

- **The contributor cannot push to the allocator branch.** A maintainer reserves a batch with `--count`, assigns one ID, and the contributor runs `clue id sync` through a readable upstream remote before using it.
- **The remote is unavailable or rejects the operation.** Allocation exits nonzero without changing the local ledger. There is no local fallback.
- **The remote claim succeeds but the local ledger write fails.** The error names the permanently claimed IDs and directs the contributor to `clue id sync`.
- **The established allocator ref is missing or malformed.** Allocation stops. It never recreates a ref that may have held abandoned claims.

## Outcome

Every contributor holds a distinct durable reservation, abandoned numbers remain unavailable, and the accepted repository still carries a deterministic ledger that can be validated offline.
