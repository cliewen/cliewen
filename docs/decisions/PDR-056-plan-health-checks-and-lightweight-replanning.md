---
id: PDR-056
type: decision
status: verified
links: [PDR-033, PDR-046, CAP-006, C-011]
supersedes: [PDR-008]
title: Plan-health checks pause invalid work; replanning is not automatically a decision
author: agent
accepted-by: Flemming N. Larsen (2026-09-10, conversation)
---

# PDR-056 — Plan-health checks pause invalid work; replanning is not automatically a decision

## Context

Implementation can show that a milestone is no longer wanted or achievable, or that its remaining dependencies no longer support the campaign. The former rule required a decision record for every plan adjustment even though C-011 limits records to future-shaping choices. It also left the point for discovering a failed plan implicit.

## Decision

Before starting or resuming milestone work, and whenever new evidence challenges the campaign, an agent assesses whether the plan still serves its goal, the milestone remains wanted and achievable, and the remaining dependencies and order still hold. A passing assessment is not recorded.

If the assessment fails, the agent records the mismatch and options in the change workspace and pauses affected work for human direction. The chosen revision is declared in the plan and may ride with the implementing change; a plan-only change remains available when it makes review clearer. Only a selected course that is future-shaping earns a typed decision record. Human direction permits work to resume, while merge binds the revised plan.
