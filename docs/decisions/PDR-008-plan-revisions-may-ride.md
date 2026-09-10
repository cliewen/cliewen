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

Implementation is often where a plan promise proves inaccurate, and the implementing pull request already contains the evidence and review context needed to examine the revision. A separate planning pull request would add ceremony without changing the human acceptance boundary.

## Decision outcome

**A semantic plan revision may ride with its implementing change when the PR names it, shows the before and after, carries a correctly typed decision record, and asks for deliberate human approval of the revision.** An explicit objection removes the revision while leaving the implementation to stand on its own merits; an undeclared revision is not accepted.

The dedicated plan change remains the default. Under [PDR-004](PDR-004-merge-binds-approval-signs.md), merge makes a declared, unobjected revision binding and explicit approval verifies its record. The semantic-mutation rule in `clue-plan` carries this allowance to adopters.
