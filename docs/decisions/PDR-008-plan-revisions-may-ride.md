---
id: PDR-008
type: decision
status: verified
links: [P-002, PDR-004, PDR-007]
title: A declared plan revision may ride with its implementing change
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# PDR-008 — A declared plan revision may ride with its implementing change

## Context and problem statement

The plan rules route semantic mutations — anything altering what a plan promises — through a dedicated change/PR, human-accepted and backed by a decision record. In practice the two often arrive together: implementing a milestone is exactly when its exit criterion's wording turns out to be wrong, and the implementing PR already contains the evidence and the reviewer's full attention. CH-020 hit this head-on: it revised M-005's criterion inside the feature PR, declared the revision loudly, and still had to withdraw it, because the rule as written gives no bundling affordance — leaving a merged feature, an open milestone, and a follow-up PR whose only content is wording the reviewer had already read. The separate-PR rule was stricter than the acceptance boundary the human actually wants to hold. What does that boundary actually protect, and can a bundled revision satisfy it?

## Decision outcome

**The boundary is deliberate human acceptance of the revision itself, not the vehicle that carries it. A semantic plan revision may ride with its implementing change when ALL of these hold:**

1. **The PR explicitly declares the revision as a plan revision** — named as such, with before and after discernible, never buried in a diff.
2. **The correctly typed decision record backs it** — normally a PDR (plan adjustments are process direction), an ADR if architectural; a decision-log row only where the existing reversal-cost litmus routes it there.
3. **The PR calls the revision out for deliberate human approval** — the reviewer is asked to accept the revision specifically, not just the change as a whole.
4. **An explicit objection reverts the revision without blocking the rest of the change** — the milestone stays open, the implementation still merges on its own merits.

The default vehicle remains a dedicated plan change/PR: bundling is the exception for revisions that surface during implementation, not a license to fold planning into feature work. Acceptance semantics are unchanged from [PDR-004](PDR-004-merge-binds-approval-signs.md): the merge makes a declared, unobjected revision binding; explicit approval promotes its record to `verified`.

**Carrier:** the semantic-mutation rule in the `clue-plan` skill (agent; ships verbatim to adopters, so the allowance travels with the methodology).

### Rejected: keeping the strict separate-PR rule

It protects nothing the four conditions do not: the reviewer sees the same words either way, and the objection path (milestone stays open) preserves exactly the outcome a rejected dedicated PR would produce. What it adds is a second PR whose content the reviewer has already read — ceremony without safety, and an incentive to leave milestones stale.

### Rejected: revisions ride freely as bookkeeping

Silent plan drift is the failure the semantic/bookkeeping split exists to prevent; a revision nobody declared is a promise nobody accepted. The declaration and the deliberate-approval callout are what make the bundled route reviewable at all.
