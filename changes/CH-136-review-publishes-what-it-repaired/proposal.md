---
id: CH-136
type: change
status: open
links: []
title: A review of existing hosted work publishes exactly what it repaired
---

# CH-136 — A review of existing hosted work publishes exactly what it repaired

**Plan-less.** This change serves no milestone of any live plan. [P-013](../../docs/plans/P-013-simplification.md) is a simplification campaign whose remaining milestones are M-067 and M-066, and [PDR-026](../../docs/decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) requires an argument before machinery is added inside one. A rule gap found in practice is not simplification, and stretching the campaign to hold it would misstate what the campaign is for. The human confirmed the plan-less declaration.

## What

State one rule in the canonical review-boundary source and carry it to every live carrier: **a review of an existing branch or pull request ends in a commit and a push exactly when it changed something.**

- Repairs are committed, verified, reviewed under the loop, and pushed to that pull request in the same turn. Publication is never deferred to a later report, and an agent never chooses a local stopping point for itself; only a human request makes committed-but-unpushed work legitimate, and the agent then says the repairs are unpublished and the pull request is not merge-ready.
- A review that produced no repairs commits nothing and pushes nothing. The reviewed commit stays exactly as reviewed, and advisories go to the verification handoff.

## Why

The obligation half-exists and never binds. [PDR-016](../../docs/decisions/PDR-016-pr-state-carries-agent-handoffs.md) clause 3 already makes the updater commit, review, push, and confirm, and [AC-041](../../docs/capabilities/CAP-006-collaborative-handoffs/criteria.md) states it as a criterion — but every carrier hangs the duty on *before publishing* or *before reporting ready*, which is a condition the agent controls. An agent that repairs, commits, and never claims readiness satisfies all of them while leaving the repairs where no handoff can reach them.

That is not hypothetical. On pull request [cliewen/cliewen#141](https://github.com/cliewen/cliewen/pull/141) an agent reviewed the branch, committed seven repairs, reported truthfully that it had not pushed, and stopped; the human merged the pull request without them. The `Durable work state` rule already says an agent's private memory is never where work lives, and an unpushed repair is exactly that — so this change states the consequence rather than inventing a principle.

The *only if* half is stated nowhere at all. The nearest sentence covers advisories from a clean pass, not a review that found nothing.

## Boundaries

The loop is not reordered: the push still follows the repaired commit's own review pass and local checks, so "same turn" is not "push before re-review". The context-isolated reviewer stays read-only; the duty falls on the updater turn. [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) is not triggered — nothing changes what the verifiable thread, the human merge boundary, or the deterministic judge means. No machine is proposed: PDR-016 already rejected making CI detect unpublished local work, because CI receives hosted commits and cannot inspect another machine's worktree.

The routing hubs are deliberately untouched. CH-132 decided they route rather than restate the detailed handoff, and duplicating this rule there would put it twice in one reading path.
