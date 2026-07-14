---
id: ADR-015
type: decision
status: verified
links: [ADR-016]
title: A light change tier — the PR description is the proposal
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #9)
---

# ADR-015 — A light change tier: the PR description is the proposal

## Context and problem statement

The change loop requires a `/changes/CH-xxx-slug/` workspace — proposal, tasks, open questions — for every mutation of `main`. For a typo fix or a dependency bump, that ceremony costs more than the change itself, and a methodology whose smallest honest step is expensive teaches people to batch or bypass it. What is the smallest loop that keeps the invariants — branch + PR, global change numbering, the review gate — without the transient workspace?

## Decision outcome

**Two tiers. A change qualifies as *light* when ALL of these hold; otherwise it is a full change:**

- no new decision is needed (no ADR, no decision-log entry — [ADR-016](ADR-016-decision-log.md));
- no acceptance criterion or capability meaning is added, changed, or retired;
- no semantic plan mutation (milestone-status bookkeeping is fine);
- no methodology carrier is touched (skills, AGENTS.md rules, lint rules).

Examples: typos, doc clarity, dependency bumps, pure refactors, CI plumbing.

**A light change is: branch `ch-xxx-slug` → commits → PR, where the PR description is the proposal.** No `/changes/` folder, no proposal.md/tasks.md/open-questions.md. The PR description states what and why, and names its plan item or declares itself plan-less — the same declarations a proposal.md carries, in the venue the reviewer already reads.

**Escalation is mandatory and immediate:** the moment an open question, a decision, or an AC change appears mid-work, create `/changes/CH-xxx-slug/` and continue the full loop from there. Qualification is re-judged by what the change *became*, not what it was expected to be.

**CH numbering stays global across both tiers** (next free number over `git log` and `/changes/`), so provenance grep works identically whether a change carried a workspace or not.

**Carrier:** the light-tier section of the `clue-delta` skill and the tier-correctness item in the `clue-verify` checklist (agent); the branch/PR invariants already enforced by `clue validate`'s digest gate (machine).

### Rejected: exempting small changes from branches and PRs entirely

Commit-to-main for "trivial" changes makes triviality a self-judged, unreviewable claim and breaks the audit trail. The branch + PR invariant is cheap; the workspace is the expensive part, and it is the only part the light tier removes.

### Rejected: a size threshold (lines/files changed) as the qualification test

Size is a proxy that fails in both directions — a one-line change can retire an AC, a 500-line rename decides nothing. The qualification test is about *meaning*: decisions, criteria, plans, carriers.
