---
id: CH-063
type: change
status: open
links: []
title: Open P-007 to harden the core against two independent critiques
---

# CH-063 — Open P-007 to harden the core against two independent critiques

## What

Create [P-007](../../docs/plans/P-007-core-hardening.md) `active`, a campaign that answers two independent methodology reviews received 2026-07-25: a model-led review of the methodology as published, and a migration assessment produced while evaluating Cliewen for a production system in an enterprise setting. Their combined findings are digested into [AN-008](../../docs/analysis/AN-008-methodology-critiques.md), which this change also writes: twenty-five findings that reduce to four patterns, each a half-built face of a core element ([ARCH-003](../../docs/architecture/core.md)). The plan carries six milestones, M-024…M-029, continuing corpus-global milestone numbering from P-006. This change writes the analysis record, the plan file, their index rows, and a decision-log row recording the campaign's opening — nothing else. It is plan-less itself: no plan can serve a change whose product is the plan, the same shape [P-006](../../docs/plans/P-006-first-adoption.md) was opened in.

## Why

Both reviews arrive at the same place from opposite directions. The methodology review finds that the human merge boundary has no stated content, that the corpus only accumulates, and that the change loop's cost is unpriced; the migration assessment finds that the deterministic judge proves less than the skill prose promises and prices that gap as its adoption blockers. Every finding lands on one of the three core elements, and in each case the agent-facing half is built while the human-facing or bounding half is prose. P-006 spent the first adoption's evidence; these reviews are the next evidence of the same kind — sharper, because one of them speaks for an adopter Cliewen has not yet won. Leaving them undigested wastes them, and several findings (an empty human gate, a judge weaker than its checklist) degrade the core's guarantees a little more with every merge they go unanswered.

## Decision boundary

This change opens the campaign and fixes its milestone set; it implements no finding. It makes no change to `clue`, the validator, the skills, the scaffold templates, the PR template, or the guide. Each milestone is a separate full change, proposed and merged on its own, and the milestones that touch a core element's meaning are flagged in the plan as red-line changes requiring their own decision records under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md). The adopter-hardening findings that depend on a real pilot are explicitly listed out of the campaign, not silently dropped.
