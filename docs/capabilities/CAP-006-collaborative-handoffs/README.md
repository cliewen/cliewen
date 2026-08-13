---
id: CAP-006
type: capability
status: active
links: [G-001, G-008, PDR-016, PDR-039, PDR-040, PDR-042, C-012]
title: Agent routing and full-loop handoffs remain visible, authorized, and exact
goal: G-001
---

# CAP-006 — Agent routing and collaborative handoffs

## What

Before editing, an agent recommends simple work or the full Cliewen loop from whether the accepted contract changes, names what would alter that recommendation, and leaves route and integration authority with the user. When the full loop is chosen, agents may independently implement, review, and repair changes without losing work between sessions: every changed turn pushes the change branch, the pull request exists as a draft from first publication, review results name the hosted commit they cover, actionable findings remain visible, and marking the pull request ready binds verification and a clean review to its exact hosted head.

## Why

The human acceptance boundary of [G-001](../../goals/G-001-verifiable-thread.md) protects a full change only when the user chooses that workflow; Cliewen advises rather than acquiring authority over an adopter's repository. Inside that loop, a private finding or local fix is invisible to the next agent and to the human, while a global change lock would needlessly serialize independent work.

The stable `validate` handoff stays in the adopter-owned caller while the upstream reusable workflow carries validation and acceptance-brief repairs, so an update does not require copying the wall's logic ([ADR-038](../../decisions/ADR-038-upstream-validation-workflow.md)).
