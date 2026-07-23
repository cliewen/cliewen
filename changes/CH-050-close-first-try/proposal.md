---
id: CH-050
type: change
status: open
links: [P-004]
title: Close P-004 — the first-try campaign is complete
---

# CH-050 — Close P-004 — the first-try campaign is complete

## What

Mark [P-004](../../docs/plans/P-004-first-try.md) `completed`: all three milestones (M-013, M-014, M-015) already show `done` with evidence, so no further product or guide work is needed. This change flips the plan's top-level `status` field, adds a completion note in the same style as [P-003](../../docs/plans/P-003-goes-public.md), and records the closure as a decision-log row per the plan's own mutation rules (the plan's top-level status is outside "the milestone table" the rules let an implementing change mutate freely). No successor plan is designated in this change.

## Why

P-004's exit criteria are all satisfied and evidenced (CH-046, CH-047, CH-049), but the plan document still reads `status: active`, leaving the corpus out of sync with reality and leaving the campaign without a recorded end. Closing it keeps `docs/plans` an accurate record of what's actually in flight.

## Decision boundary

This change only flips P-004's status and records the closure decision. It does not start, scope, or imply any successor campaign; does not touch `clue`, validation semantics, skills, or any other plan; and does not reopen or re-evidence the already-`done` milestones. Any proposal for what comes next is a separate, later change.
