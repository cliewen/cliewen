---
id: CH-136
type: change
status: open
links: [PDR-007, PDR-012, PDR-016, PDR-033, C-012, CAP-006]
title: Push is durability, ready is the explicit act
---

# CH-136 — Push is durability, ready is the explicit act

## What

Replace the rule that a push signals a reviewed, ready candidate with a simpler model carried by one new decision record and every live methodology carrier:

1. **Push is durability, not readiness.** Every working turn on a Cliewen change that changed anything ends in commit and push to the change branch. Nothing that exists only in a local worktree is a legitimate resting state for work.
2. **The pull request exists from first publication, as a draft.** A full change opens its draft PR right after the proposal commit; a light change after its first commit. Unfinished work lives on the hosted branch behind a draft, not in private agent memory.
3. **Ready is an explicit act.** The PR is marked ready for review only when local verification and a clean agentic review pass are bound to the exact hosted head SHA. A substantive edit after that returns the PR to draft until its new head earns the same binding. The acceptance brief belongs to the ready PR, not the draft.
4. **Hosted history is never rewritten.** Because publication begins with the first commit, the pre-publication rebase window closes; accepted `main` is always incorporated by merging it into the change branch.

The merge boundary does not move: the agent still never merges its own PR, never pushes to `main`, and the human merge still accepts the change. `clue validate` still judges. C-013 is not triggered; this reshapes agent-side publication timing inside the existing boundary.

## Why

On the pull request at <https://github.com/cliewen/cliewen/pull/141>, an agent reviewed the branch, committed seven repairs, truthfully reported them unpushed, and stopped; the human merged without them. The first attempted fix, <https://github.com/cliewen/cliewen/pull/143>, kept the assumption that a push means "reviewed and ready" and therefore had to add machinery on top of it: a biconditional publication rule, a human-requested stopping-point exception, and an unresolved tension with the review loop's no-publication-with-blocking-findings rule when the loop does not converge. The maintainer closed it and chose the simpler model: make the push carry no meaning at all, and move every signal into explicit, observable PR state. The stopping-point machinery then has nothing to govern — work is always as published as it can be, and readiness is a separate, deliberate act.

## How

- The canonical review-boundary fragment and the `clue-verify` and `clue-delta` templates under `internal/skills/source/` state the model; `go generate ./internal/skills` propagates it to the four rendered skills and the scaffolded adopter copies. The pause step's "whether the branch should be pushed" question is removed as moot, amending that clause of PDR-033.
- CI distinguishes draft from ready without weakening the merge gate: on a draft, the digest gate (`--forbid-changes`) and the acceptance-brief requirement do not apply, because a draft legitimately carries its workspace and no brief; the strict gate runs on `ready_for_review`, on every ready-PR synchronize, and on `main`. GitHub cannot merge a draft, so enforcement at the merge boundary is unchanged (C-004 assessed).
- PDR-040 records the decision; PDR-007, PDR-012, PDR-016, and PDR-033 gain amendment notes; C-012 and CAP-006 carry the contract; a new Human-class criterion states it; CONTRIBUTING, the change-loop guide, and both PR templates say it to humans.

## Plan

This change is explicitly plan-less, decided with the maintainer: P-013's remaining milestones (M-067, M-066) own other work, and PDR-026 protects a simplification campaign from absorbing a methodology redesign motivated by an incident. The maintainer approved the model and this change in conversation on 2026-08-09.
