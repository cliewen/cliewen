---
id: CAP-010-design
type: design
status: active
links: [CAP-010, ADR-068]
title: Design for team-safe identity allocation
---

# Design — Team-safe identity allocation

## Two records with different jobs

`.clue/id-ledger.yaml` is the checked-in lifecycle truth. Version two is an append-only event set: commands add one flow-style YAML line per state transition, loading folds duplicate events idempotently, and state only advances `reserved → live → retired`. The `.gitattributes` union driver combines concurrent additions; contradictory identity metadata remains a validation error.

`refs/heads/clue/id-allocator` is the remote allocation journal. Its orphan history contains numeric claims only, so a feature branch's speculative lifecycle cannot make accepted `main` claim that an unmerged artifact is live. The branch is permanent because an abandoned claim still prevents reuse.

## Allocation transaction

The Git transport uses plumbing commands and temporary refs rather than checking out the allocator branch or touching the contributor's index. It fetches and validates the current claim file, unions numeric identities already known by the local ledger, allocates after the prefix high-water mark, creates a child commit with a random transaction identity, and performs an ordinary push. The transaction identity prevents two same-second contenders from creating the same commit object. A concurrent winner therefore makes the other push non-fast-forward; the loser fetches the new head and retries until the operation-wide timeout.

No force option is used. A push error is retried only when a follow-up remote read proves the head changed; authentication and transport failures surface immediately. Local ledger bytes change only after the claim is remote-durable.

## Initialization and recovery

`clue id coordinate` is the only operation that may create the allocator branch. It seeds claims from every numeric ledger identity, wins or retries a concurrent initialization race, then records `mode: git` and the remote name locally. An established coordinated repository treats a missing ref as possible history loss and stops.

`clue id sync` needs read access only. It imports missing claims as reservations and never downgrades a local live or retired state. This is also the recovery after a remote push succeeded but the local atomic replacement failed.

Allocation settings are kept in `.clue/id-coordination.yaml` rather than in the ledger, and the union rule names the ledger by path, so the settings file merges normally. Two branches that each coordinated to their own remote therefore produce an ordinary conflict that Git raises at merge time and stops on, rather than a combined file that reaches the integration branch already broken. The absence of the file means local allocation. Pointing an established repository at a different remote is refused unless forced, because it abandons the claims recorded on the remote it leaves.

Damage from the merge driver itself is narrower than it sounds. Union merge combines line by line, so anything both branches wrote identically merges as ordinary context: appended events, rewritten files, and two branches enabling coordination against the same remote all merge cleanly. Only genuinely divergent lines are duplicated, and with the settings moved out there is no longer a divergent line the ledger can carry.

A repeated key is reconciled on read wherever it appears: sequences concatenate, because events fold idempotently and advance state monotonically, and repeated scalars must agree. Saving rewrites the file whole, so a command that allocates also repairs; `clue id repair` does it on its own, and `clue migrate` plans the same repair for a repository being brought up to date. Two values that disagree stop instead, naming both, because only a person can decide which the team meant ([ADR-069](../../decisions/ADR-069-a-combined-ledger-is-recoverable.md)).

## Boundary

`clue validate` reads only repository bytes, as [ADR-044](../../decisions/ADR-044-judge-reads-state-not-transitions.md) requires. Network access belongs only to the explicit allocation, coordination, and synchronization commands. Remote branch protection remains a repository administration responsibility: the branch permits normal pushes and forbids force-push and deletion.
