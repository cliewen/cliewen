---
id: PDR-007
type: decision
status: verified
links: [PDR-004, C-012]
title: The review boundary is real — changes root at main, one in flight per author, humans merge
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# PDR-007 — The review boundary is real

## Context and problem statement

The change loop calls the PR the review gate and the merge the acceptance ([PDR-004](PDR-004-merge-binds-approval-signs.md)), but describes the boundary without defending it: nothing states who performs the merge, where a change may start from, or how many changes an author may have open. An agent can satisfy every written sentence while producing a corpus no human ever reviewed — implementing several changes stacked on each other's unmerged tips, "merging" them with locally created merge commits, and pushing the result to `main`. The locally preserved [AN-002 adopter audit](../analysis/AN-002-model2diagram-extraction.md) observed exactly this. Mechanical validation cannot catch it: the corpus validator judges repository contents, not whether a review happened. What rules make the review boundary unfakeable by the agent, without serializing a whole team onto one change at a time?

## Decision outcome

**The boundary is defined by where changes start, who merges, and how many changes an author holds open — not by limiting team concurrency.**

1. **Every change branches from the current tip of `main`** — never from another change's branch, never from any commit not yet accepted into `main`. This is mechanically checkable (`git merge-base`) and leaves team parallelism unlimited: any number of changes may be in flight as long as each roots at `main` and carries its own reviewable PR.
2. **One change in flight per author.** An author — human or agent — takes a change to its PR before starting the next. For an autonomous agent this is the stop condition: after opening the PR it waits; it never starts the next change while its previous one is unreviewed.
3. **The merge is a human act.** An agent never merges its own PR, never creates a merge commit into `main` locally, and never pushes to `main`. A locally created merge commit is not an acceptance: no review object, no approval event, no human action exists behind it. Until a human merges, the change is not accepted.
4. **Review fixes belong to the reviewed change.** Corrections to an unaccepted change land on that change's existing branch and PR — never as a new CH. A new CH for follow-up work exists only when a human has accepted the current change and explicitly scoped the follow-up.
5. **Stacking is an explicit human decision.** Work that genuinely must build on an unmerged change is a blocking open question; the human's answer is recorded. The agent never chooses to stack.
6. **A sibling merge triggers a rebase.** When a parallel change merges first, the open branch rebases onto the new `main` tip and re-runs the pre-merge checklist before its PR proceeds — small deltas plus this re-check are the standing mitigation for parallel work.

**Carrier:** the branch and boundary wording in `clue-delta` (steps 1 and 5) and the process items in `clue-verify` (agent); the register entry [C-012](../constraints/C-012-agents-never-merge-own-changes.md), whose promotion trigger names the machine enforcement (branch protection or PR-provenance CI) where hosting permits.

### Rejected: one change in flight globally

Trivially prevents stacking, but serializes an entire team onto a single change — the methodology would forbid ordinary parallel feature work. The observed failure was dependency (changes rooted on unmerged work), not concurrency; a rule against concurrency punishes the wrong thing.

### Rejected: one change in flight per plan

Plan membership does not track conflict risk: two changes inside one plan can be fully independent while two changes in different plans collide on the same artifact. Per-plan serialization would have allowed the observed stack (spread across plans) and blocked harmless parallel work inside one plan — wrong in both directions. The plan layer is bookkeeping, not a locking domain.

### Rejected: treating a local merge commit as acceptance

A merge commit is a graph shape, not a review. Acceptance under PDR-004 is a human action on a PR; a commit the agent can fabricate alone cannot carry it.
