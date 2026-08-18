---
id: PDR-001
type: decision
status: verified
links: [ADR-010, ADR-011, ADR-012]
title: PR approval is ADR acceptance — the agent performs the clerical promotion
author: agent
accepted-by: Flemming N. Larsen (2026-07-13, PR #8)
---

# PDR-001 — PR approval promotes the PR's ADRs

> **Superseded by [PDR-004](PDR-004-merge-binds-approval-signs.md):** merge makes a decision binding, while only explicit human approval makes it `verified`; the surviving rule here is that the agent performs the clerical provenance update.

## Context and problem statement

Agent-authored decisions are born `inferred`, and human acceptance must remain distinct from authorship. Requiring a second manual edit after review made accepted decisions remain unbookkept and made the inferred count unreliable.

## Decision outcome

**The agent performs the clerical update after a recorded human approval; no decision promotes itself.** The update records the approver, date, and PR or conversation in `accepted-by:` and changes `status: inferred` to `verified` only for the explicitly approved decision. An objection keeps the decision `inferred` until a later human approval; approval of unrelated change content is not enough.

The promotion wording belongs in the decisions README and the `clue-verify` guidance. A decision is never born `verified`, and the PR or conversation named in `accepted-by:` remains its acceptance record.
