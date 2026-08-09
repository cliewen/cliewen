---
id: PDR-040
type: decision
status: inferred
links: [G-001, PDR-007, PDR-012, PDR-016, C-012, CAP-006]
title: A review of existing hosted work ends in a push exactly when it produced a repair
author: agent
accepted-by: []
---

# PDR-040 — A review of existing hosted work ends in a push exactly when it produced a repair

## Context and problem statement

[PDR-016](PDR-016-pr-state-carries-agent-handoffs.md) makes the updater commit, verify, review, push, and confirm the hosted head. Every carrier of that rule states it as a precondition of publishing or of reporting a change ready — both of which the agent decides when to attempt. An agent that repairs a reviewed pull request, commits the repair, and then simply reports that it has not published satisfies each of those sentences with the repair sitting in a local worktree.

The `Durable work state` rule already says an agent's private memory is never where work lives, and a repair no one can fetch is exactly that. What was missing is the consequence, stated as an obligation with its own trigger.

The converse was unstated too. Nothing said that a review finding nothing to repair leaves the reviewed commit alone. Without it, an agent may edit a commit that a clean pass already covered, which the exact-commit boundary then requires reviewing again — the loop restarts for no finding.

## Decision outcome

**A review of an existing branch or pull request ends in a commit and a push exactly when it changed something.**

- **Repair implies publication, in the same turn.** Repairs are committed, verified, reviewed under the loop, and pushed to that pull request before the turn ends. Publication is not deferred to a later report of readiness.
- **The stopping point is the human's to choose.** An agent never elects to leave a repair committed and unpushed. Only an explicit human request makes that a legitimate stopping point, and the agent then states that the repairs are unpublished and the pull request is not merge-ready.
- **No repair implies no commit and no push.** A review that produced nothing to repair leaves the reviewed commit exactly as reviewed. Advisories go to the verification handoff rather than to an edit.

This does not reorder the review loop. The push still follows the repaired commit's own local checks and its own review pass, so the published commit remains the reviewed one; "same turn" constrains when publication happens, not what it publishes. The context-isolated reviewer remains read-only, and the obligation falls on whichever context becomes the updater for that turn.

The rule is human-enforced. It governs what an agent does with a worktree and a remote, and no artifact in the corpus records either.

## Rejected: make CI detect unpublished local work

PDR-016 already rejected this and the reason is unchanged: CI receives hosted commits and events, and cannot inspect another machine's worktree or infer an unreported intention. Restating the rejection here keeps the question from being reopened as though the enforcement gap were new.

## Rejected: forbid a local stopping point outright

A human legitimately asks for work to be committed but not published — to inspect a diff, to hold a change over, to hand it to someone else. Removing the stopping point would make the honest report of that state impossible, so the repair is to name who may choose it rather than to delete it.

## Rejected: treat the unpublished repair as a reviewer failure

The reviewer returns findings and no edits; it has nothing to push. Locating the duty in the review step would leave it unowned whenever review and repair happen in the same context, which is the ordinary case.

## Carrier

The shared review-boundary source states the rule and the `clue-verify` loop and checklist carry it into verification (agent); [C-012](../constraints/C-012-agents-never-merge-own-changes.md) registers it; both pull-request templates carry the human-facing check (default/local); `CONTRIBUTING.md` and the public change-loop guide explain it; [AC-131](../capabilities/CAP-006-collaborative-handoffs/criteria.md) holds the acceptance meaning, proven by the acceptance brief. The generator test pins the clause text across every generated and scaffolded skill. This record amends PDR-016's updater clause; its other rules stand.
