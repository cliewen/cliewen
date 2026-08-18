---
id: PDR-004
type: decision
status: verified
links: [PDR-001, ADR-010]
title: Merge makes a decision binding; approval verifies it — approvers sign, first signature dates it
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #11 review)
---

# PDR-004 — Merge binds, approval signs

## Context and problem statement

Merging and approving express different judgments: merge makes the reviewed change operative, while explicit approval says that a human endorses the decision. Conflating them would either block integration on a ritual or make `verified` mean only "appeared in a merged PR."

## Decision outcome

**Merging a pull request makes its `inferred` decisions binding; only explicit human approval makes them `verified`.** Approval may be a review approval, review comment, or recorded approval in conversation. The agent performs the clerical signing, appending each approver to `accepted-by:` while the first approval supplies the acceptance date. An explicit objection blocks verification of that decision but not necessarily the rest of the merge, and no decision self-promotes.

The promotion wording is carried by the decisions README, `clue-delta`, and `clue-verify`; the merge and frontmatter already carry the binding and provenance facts. The `inferred` count therefore remains the visible backlog of decisions in force without explicit human endorsement.
