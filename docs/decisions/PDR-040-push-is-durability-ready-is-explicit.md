---
id: PDR-040
type: decision
status: inferred
links: [G-001, PDR-007, PDR-012, PDR-016, PDR-033, C-004, C-012, CAP-006]
title: Push is durability and the ready mark is the only readiness claim
author: agent
accepted-by: []
---

# PDR-040 — Push is durability, ready is the explicit act

## Context and problem statement

The review boundary let a push carry meaning: opening a PR asserted a reviewed candidate, so pushing was deferred until verification and review had passed, and unfinished work waited in a local worktree. An agent that repaired a branch, committed, truthfully reported the commit unpushed, and stopped broke no rule — and the human merged the pull request at <https://github.com/cliewen/cliewen/pull/141> without those repairs. The first attempted fix, closed unmerged at <https://github.com/cliewen/cliewen/pull/143>, kept the push-as-signal assumption and therefore needed exception machinery on top of it: a repair-implies-push biconditional, a human-requested stopping-point rule, and an unresolved conflict with the review loop's rule against publishing with blocking findings outstanding when the loop does not converge. When must work be pushed, and what may a push claim?

## Decision outcome

**A push claims nothing; the ready mark claims everything. Every working turn that changed anything ends by committing and pushing the change branch, the pull request exists as a draft from first publication, and marking it ready for review is the explicit act that binds green local verification and a clean agentic review pass to the exact hosted head.**

- **Push is durability.** A local worktree is private agent memory and no handoff survives it. A turn ends with the branch pushed, whatever state the work is in; a turn that changed nothing pushes nothing, so a reviewed commit stays exactly as reviewed. The one exception is a pull request found merged or closed: the turn stops without pushing and reports where the work stands.
- **The draft PR is where unfinished work lives.** A full change opens its draft right after the proposal commit, a light change after its first commit. The draft claims nothing, cannot be merged, and gives hosted CI, review conversations, and the human a continuous view of the work.
- **Ready is the readiness claim.** The PR is marked ready only when local verification and the review loop have passed on the current commit and the hosted head equals it. A substantive edit afterwards — including a review repair — returns the PR to draft until the new head earns the same binding. The stopping-point machinery disappears because there is nothing to govern: stopping anywhere is ordinary, the branch is pushed, and no readiness claim exists until the mark.
- **Hosted history is never rewritten.** Publication begins with the first commit, so the pre-publication rebase window closes; accepted `main` is always incorporated by merging it into the change branch.
- **CI distinguishes draft from ready without weakening the merge gate.** The digest gate (`--forbid-changes`) and the acceptance-brief requirement do not apply to a draft, which legitimately carries its workspace and no brief; they bind on `ready_for_review`, on every ready-PR update, and on `main`. The forge cannot merge a draft, so nothing reaches the merge boundary ungated, and [C-004](../constraints/C-004-never-weaken-checks.md) is satisfied: the checks moved to where their meaning applies, and the merge precondition is unchanged.

This amends [PDR-007](PDR-007-review-boundary.md) (the PR opens as a draft before the candidate exists, and the one-in-flight stop condition anchors to the ready mark, not the PR's creation), [PDR-012](PDR-012-agentic-review-before-publication.md) (review gates the ready mark, not the push), [PDR-016](PDR-016-pr-state-carries-agent-handoffs.md) (the updater pushes with the turn that repaired, and the handoff completes at the ready mark), and [PDR-033](PDR-033-planning-and-implementation-are-separate-steps.md) (the spec-first pause no longer asks whether to push, because the proposal is already pushed and visible). The human merge boundary, the deterministic judge, and the verifiable thread are untouched; [C-013](../constraints/C-013-core-changes-need-decision.md) is not triggered.

### Rejected: push as a signal with a publication-obligation rule

The closed <https://github.com/cliewen/cliewen/pull/143> approach: keep "push means reviewed and ready" and add a rule that a repair owes its push in the same turn. It required a human-requested stopping-point exception, left an unpublished local branch without a referent for "that PR", and contradicted the review loop's no-publication rule whenever repairs existed but the loop had not converged. Machinery grew to defend an assumption the model does not need.

### Rejected: keep never-draft and rely on agent discipline

The status quo. It makes every local worktree a place where finished-looking work can silently die, which is exactly what happened; discipline had already failed once under truthful reporting, because no rule was broken.

### Rejected: make CI detect unpublished local work

Carried forward from PDR-016: CI receives hosted commits and events; it cannot inspect another machine's worktree or infer an unreported intention. The always-push turn rule shrinks what only the agent can see instead of pretending a machine can see it.

**Carrier:** the shared review-boundary fragment and the `clue-delta`/`clue-verify` templates (agent); `ci.yml`, the reusable validation workflow, and the scaffolded caller (machine); [C-012](../constraints/C-012-agents-never-merge-own-changes.md), both pull-request templates, `CONTRIBUTING.md`, and the change-loop guide (human); CAP-006's criteria state the contract.
