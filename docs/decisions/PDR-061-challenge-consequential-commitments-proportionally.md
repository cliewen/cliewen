---
id: PDR-061
type: decision
status: inferred
links: [G-014, P-023, M-094, AN-024, PDR-056, PDR-042, ADR-062, C-013]
title: A consequential commitment is challenged before it is made, and the challenge scales with consequence rather than applying to everything
author: agent
accepted-by: pending human verification (proposed 2026-09-12)
---

# PDR-061 — Challenge a consequential commitment proportionally

## Context

Cliewen's thread proves that what was built matches what was stated; it cannot prove the statement was worth making. A criterion, its implementation, and its test are usually written by one author in one sitting, so one wrong assumption passes through all three and emerges looking verified ([G-014](../goals/G-014-riskiest-assumption-is-challenged.md)). [PDR-056](PDR-056-plan-health-checks-and-lightweight-replanning.md)'s plan-health check asks whether a plan still holds once work is under way; nothing asks whether it should have been adopted. Agents sharpen the gap, because an agent asked to plan plans fluently, and fluency reads as confidence.

## Decision

**Before a consequential commitment, the proposer challenges it where the reviewer will read it** — in the plan's prose for a plan, and in `proposal.md` for a change. The challenge names four things: the assumption most likely to undermine the work, a credible alternative course, the cheapest useful test of that assumption, and the result that would stop or revise the work. It also asks what an implementation could look like that met every criterion and still failed the person the work is for, because that is the failure the thread cannot catch.

**A commitment is consequential when being wrong would cost more than finding out first.** In practice that is adopting a plan or revising what it promises, starting a milestone whose approach is not yet settled, choosing between credible courses in a decision record, and any rule or behaviour that will reach adopters. It is not consequential when the work is simple under [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md), when it carries out a course already challenged and well understood, or when it can be undone inside the same change at no more cost than the challenge.

**The proportionality limit is part of the rule, not an exemption from it.** Work that is not consequential proceeds without the challenge and without a record explaining its absence. No prototype is ever universally required; a small path across the consequential boundaries is warranted only where the architecture is materially uncertain. A challenge demanded of everything is written by habit, and a challenge written by habit is the confident, unexamined text this rule exists to interrupt — so a rule without the limit would defeat itself, which is why the limit cannot be loosened while the rule is kept.

**Investigation is delivery.** A milestone whose deliverable is evidence that changes a decision, rather than a feature, is a legitimate use of the milestone.

**Repository experience is evidence, not authority.** A challenge cites what this repository has already learned where it applies, and treats older guidance as something to reassess against the case in hand rather than as settled.

## Scope

The rule binds this repository only, stated in [`AGENTS.md`](../../AGENTS.md). [P-023](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md) tries it on real work in both directions before any of it reaches the shipped carriers, so this record carries no `binds: adopter` and the carrier check of [ADR-062](ADR-062-repository-role-is-declared-machine-state.md) does not apply to it. Moving it to adopters, or recording why it should not move, is a later decision.

## Rejected: a mandatory challenge section in every proposal

A fixed section in the proposal skeleton or pull-request template would make the obligation impossible to forget. It would also be filled in on every simple-looking full change, so it would mostly collect boilerplate, and it would make the limit unenforceable because an empty field reads as a skipped duty. It remains the next move if the trial shows the stated rule being skipped where it applied.
