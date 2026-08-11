---
id: AN-020
type: analysis
status: active
links: [P-013, PDR-012, PDR-016, PDR-018, PDR-021, PDR-029, PDR-035]
title: Routing-hub carrier placement after the M-063 decision
---

# AN-020 — Routing-hub carrier placement after the M-063 decision

## Question

AN-018 found the same detailed review boundary in the routing hubs and the skills. The hubs are read before a skill, while the skills are the versioned standalone carriers. The question was whether safety required retaining the full procedure in both places.

## Observed boundary

The repository and scaffolded hubs can safely own the pre-corpus decision: they classify plain work locally, retain its ordinary branch/PR/human-merge route, and direct light and full work to `clue-delta`. The canonical `clue-delta` skill renders change tiers and the review boundary before its procedure, so a light or full reader encounters the detailed branch, hosted-head, review, and merge-mode rules before implementation steps.

The detailed handoff is therefore not an additional rule the hubs need to restate. Keeping it there makes the ordinary hub-to-skill path read the same contract twice. Moving it into the canonical skill sources preserves the contract while making those sources the one place a repair must be made and generated outputs must be checked.

## Result

The human approved that placement in the CH-132 conversation. The decision-log row records it; the affected PDR carrier inventories name the skills and the hubs' routing function; and the scaffold tests now prove both sides of the relationship: the hub routes light and full work to `clue-delta`, and the generated `clue-delta` carrier contains the detailed handoff and merge-history rules.
