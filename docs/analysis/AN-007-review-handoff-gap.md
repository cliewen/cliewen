---
id: AN-007
type: analysis
status: verified
links: [PDR-007, PDR-011, C-012, LOG-001]
title: A ready PR can omit intended local review fixes
---

# AN-007 — A ready PR can omit intended local review fixes

## Risk

An agent may report a change ready while intended review fixes exist only in its local worktree, allowing the human to merge a PR that does not contain the state the agent verified.

## Evidence boundary

The repository was inspected at `c1b26057a9f46d36399c3393e5df87c34cb860dc` on Windows 11 Pro amd64 with PowerShell 7.6.3, GitHub CLI 2.87.3, Go 1.26.5, Node.js 24.0.0, and npm 11.18.0. PR #36 merged CH-037 as `f5029a83301bdac61399d51b76793236885f3ca5`; its three commits end at `e538042d317d719144ff41bfeba26ad0422ab535` and do not contain the later workflow gate or scaffold classifier test. PR #37 merged those two files as `c1b26057a9f46d36399c3393e5df87c34cb860dc` from commit `27b358b8b1ab5270a0933509e2532aecd7b555bf`, and its PR description states that the review fixes missed PR #36 because they were never committed.

The maintainer supplied the complete agent conversation on 2026-07-20. It records the agent finding the edits in the local worktree after PR #36 merged, offering a new local branch and commit as one stopping point, and opening PR #37 only after the maintainer asked how that commit could merge. The transcript is evidence of the reported interaction and maintainer expectation, but it is not a versioned repository artifact. Git cannot reconstruct an earlier uncommitted worktree, so the exact local state before the merge is not independently reproducible. No inference about intent is drawn from repository activity alone.

## Observations

- PDR-007 already says review corrections to an unaccepted change land on that change's existing branch and PR. Leaving the fixes uncommitted violated the rule; it did not reveal a missing plain-change classification rule.
- `clue-verify` runs before opening or updating a Cliewen PR and checks that the branch roots at current `main`, but it does not require a clean worktree or compare local `HEAD` with the hosted PR head.
- PDR-011 already gives plain work its complete route before any Cliewen skill is loaded: branch from accepted `main`, relevant checks, ready PR, and human merge. The incident changes CI and executable evidence, both protected surfaces, so the plain route was not applicable.
- CI and hosted review can inspect only pushed commits. No server-side check can discover intended edits that remain solely in an agent's local worktree.
- After PR #36 merged, PDR-007 required explicit human scope before a new follow-up change. The maintainer supplied that scope by selecting a new branch and local commit, but that stopping point was preservation rather than a complete or mergeable change.

## Options assessed

Leaving the carriers unchanged would treat the incident only as agent noncompliance. That preserves the existing compact checklist but leaves no explicit tripwire for the invisible state that caused the incomplete handoff.

Adding a Cliewen skill for plain changes or another classify-before-edit front door would address a different problem and would reverse PDR-011's deliberate boundary. AGENTS.md already classifies before the corpus and already requires a branch and ready PR for plain work.

Adding a CI check for uncommitted edits is infeasible because CI receives only hosted commits. CI can prove properties of a PR head, not that an agent has no additional local intent.

The proportionate carrier change is to define readiness as an exact handoff: every intended edit committed, a clean worktree, the branch pushed, a ready hosted PR, the local branch matching the PR head branch, and local `HEAD` matching the PR head SHA. Review updates repeat verification against that exact commit. A human-requested earlier stop remains allowed but is reported as incomplete and not mergeable.

## Finding and consumer

PDR-007 and PDR-011 already decide the branch, PR, review-fix, and plain-change boundaries. CH-039 consumes this finding by adding a cheap local readiness decision to LOG-001, strengthening C-012 and the routing and skill carriers, and adding regression coverage without creating a new PDR or a plain-change skill.
