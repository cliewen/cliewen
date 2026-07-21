---
id: PDR-007
type: decision
status: verified
links: [PDR-004, C-004, C-012]
title: The PR is the authorization boundary — changes root at main and humans merge
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# PDR-007 — The PR is the authorization boundary

## Context and problem statement

The change loop calls the PR the review gate and the merge the acceptance ([PDR-004](PDR-004-merge-binds-approval-signs.md)), but that shorthand obscures the safeguard: the PR is where an agent-prepared candidate stops being under the agent's sole control. Without that boundary, an agent can implement several changes on unmerged work, manufacture local merge commits, push them to `main`, and produce a corpus no human authorized. An adopter audit in the private `model2diagram` repository (`AN-002 — Cliewen workflow audit`) observed exactly this; the [analysis evidence boundary](../analysis/README.md) explains why that historical artifact is not publicly resolvable. Mechanical validation judges repository contents, not who authorized integration. What rules make acceptance unfakeable by the agent without requiring duplicate human code review or serializing a whole team onto one change at a time?

## Decision outcome

**The PR is Cliewen's authorization and protected-integration boundary: the agent may prepare and publish a candidate, but only a human-controlled merge accepts it.** It is a safeguard, not a requirement that a solo developer repeat a code review already completed locally. Where hosting supports enforcement, the PR also gives required hosted CI a candidate it can block before integration; a PR without a required status check and branch protection displays CI but does not enforce it.

1. **Every change branches from the current tip of `main`** — never from another change's branch, never from any commit not yet accepted into `main`. This is mechanically checkable (`git merge-base`) and leaves team parallelism unlimited: any number of changes may be in flight as long as each roots at `main` and carries its own reviewable PR.
2. **One change in flight per author.** An author — human or agent — takes a change to its PR before starting the next. For an autonomous agent this is the stop condition: after opening the PR it waits; it never starts the next change while its previous one is unreviewed.
3. **The PR merge is the human authorization act.** An agent never merges its own PR, never creates a merge commit into `main` locally, and never pushes to `main`. A human may inspect and accept the candidate locally before publication, but that does not authorize the agent to integrate it; until a human controls the PR merge, the change is not accepted.
4. **Review fixes belong to the reviewed change.** Corrections to an unaccepted change land on that change's existing branch and PR — never as a new CH. A new CH for follow-up work exists only when a human has accepted the current change and explicitly scoped the follow-up.
5. **Stacking is an explicit human decision.** Work that genuinely must build on an unmerged change is a blocking open question; the human's answer is recorded. The agent never chooses to stack.
6. **A sibling merge triggers a rebase.** When a parallel change merges first, the open branch rebases onto the new `main` tip and re-runs the pre-merge checklist before its PR proceeds — small deltas plus this re-check are the standing mitigation for parallel work.
7. **Hosted CI is enforced by the protected PR, not by convention.** The PR triggers the hosted checks; branch protection makes the required check a merge precondition. The PR alone is insufficient, and a locally reported check is evidence rather than enforcement because the agent controls that environment. A change to the workflow or required-check policy remains subject to [C-004](../constraints/C-004-never-weaken-checks.md): never weaken a test or lint rule to make work pass.

**Carrier:** the branch and boundary wording in `clue-delta` (steps 1 and 5) and the process items in `clue-verify` (agent); the register entry [C-012](../constraints/C-012-agents-never-merge-own-changes.md), whose promotion trigger names the machine enforcement (branch protection or PR-provenance CI) where hosting permits; and the public change-loop guide, which explains why the PR remains mandatory after local acceptance.

### Rejected: one change in flight globally

Trivially prevents stacking, but serializes an entire team onto a single change — the methodology would forbid ordinary parallel feature work. The observed failure was dependency (changes rooted on unmerged work), not concurrency; a rule against concurrency punishes the wrong thing.

### Rejected: one change in flight per plan

Plan membership does not track conflict risk: two changes inside one plan can be fully independent while two changes in different plans collide on the same artifact. Per-plan serialization would have allowed the observed stack (spread across plans) and blocked harmless parallel work inside one plan — wrong in both directions. The plan layer is bookkeeping, not a locking domain.

### Rejected: treating a local merge commit as acceptance

A merge commit is only a graph shape. A human can perform a local merge, but the repository cannot distinguish that act from a merge the agent fabricated and pushed; an agent-controlled local commit therefore cannot carry the authorization boundary. The protected PR supplies the externally visible human merge event and, where configured, the required hosted check.
