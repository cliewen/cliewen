---
id: PDR-007
type: decision
status: verified
links: [PDR-004, PDR-027, PDR-039, PDR-040, PDR-042, C-004, C-012]
title: The PR is the authorization boundary — changes root at main and humans merge
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# PDR-007 — The PR is the authorization boundary

> **Scope amended by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** this boundary governs a chosen full loop; simple integration follows explicit user authority and repository policy.

> **Amended by [PDR-040](PDR-040-push-is-durability-ready-is-explicit.md) and [PDR-039](PDR-039-dependent-changes-carry-authorization.md):** a full change is published as a draft from its first commit, readiness follows current-head review, and any authorized dependent base is disclosed without becoming accepted.

## Context and problem statement

Mechanical validation sees repository contents but not who authorized integration. The methodology needs a boundary that prevents an agent from accepting its own work while preserving parallel human review and the evidence carried by the pull request.

## Decision outcome

**For a chosen full loop, the pull request is Cliewen's authorization and protected-integration boundary: the agent prepares and publishes a candidate, and only a human-controlled merge accepts it.**

- Every full change branches from accepted `main`, never from another unaccepted change.
- One full Cliewen change may be in flight per initiating author; reviewing or updating an existing pull request does not create another change or a global lock.
- The agent never merges its own pull request, creates a local merge into `main`, or pushes directly to `main`; review fixes stay on the existing pull request.
- Stacking on an unmerged change requires an explicit human decision recorded with the dependent work.
- After publication, a sibling merge is incorporated with a normal merge and renewed checks rather than rewriting the reviewed branch; rebasing remains available before first publication.
- Hosted CI, branch protection, and required conversation resolution enforce admission where configured, but they are not acceptance evidence for a criterion.

The branch and boundary wording is carried by `clue-delta`, `clue-verify`, C-012, and the public change-loop guidance. Simple work is governed by PDR-042 instead of this full-loop boundary.
