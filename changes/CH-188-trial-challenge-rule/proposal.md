---
links: [P-023, G-014]
---

# CH-188 — Trial the challenge rule on real work, in both directions

## What

Closes P-023/M-095. Applies [PDR-061](../../docs/decisions/PDR-061-challenge-consequential-commitments-proportionally.md)'s challenge to two real, currently unsettled commitments in this campaign, and records the result in a new analysis, `AN-026`:

1. **A commitment redirected by the challenge:** how M-095 itself should be satisfied. The unexamined path was to write two illustrative examples of the rule working. The challenge caught that the milestone's own text excludes exercises, found that two real unsettled decisions (M-096 and M-101) already existed to test the rule on, and redirected the milestone's delivery to use them instead of inventing a case.
2. **A commitment that proceeds without extra investigation:** whether M-101 (reconciling `guide/methodology.md` and `guide/intent.md` against ARCH-003) needs a decision record before it can proceed. The challenge's cheapest test — reading all three documents — found no core-meaning conflict: ARCH-003 does not name vision or use-case among its three core elements, and both guides state the delivery thread identically. The routing conclusion (editorial, no ADR/PDR, no red-line change) is recorded so M-101 can proceed straight to the edit when it is next picked up.

## Why

M-095 is the only milestone that produces evidence for G-014: a rule that has only ever said yes has not been demonstrated. Both trials are genuine campaign decisions, not manufactured exercises, so the evidence is about how the rule performs on the work this repository actually has queued.

## Scope

This change adds `AN-026` and updates P-023's M-095 row (`todo` → `done`, evidence field). It does not execute M-096 or M-101 themselves — those remain separate `todo` milestones; this change only records the routing conclusion that will let M-101 proceed without a decision record when it is taken up, and the observation that redirected M-095's own approach.

No shipped carrier changes; no CHANGELOG entry owed.
