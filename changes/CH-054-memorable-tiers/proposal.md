---
id: CH-054
type: change
status: open
links: [P-005, M-019, G-001, C-013]
title: State change-tier routing as three rules and two guards
---

# CH-054 — Tier routing is memorable

## What

Rewrite the canonical change-tier text (`internal/skills/source/shared/change-tiers.md.tmpl`) from four dense enumerating paragraphs into three named rules — plain, light, full — plus two escalation guards. Semantics are unchanged: the same work lands in the same tier, and every obligation the old text imposed still holds. Align the prose copies that restate the same routing: `AGENTS.md`, the scaffold `AGENTS.md` template, `guide/change-loop.md`, and `guide/methodology.md`. Regenerate the skills at version 0.6.0 (already stamped in CH-053) and record the rewrite as a decision-log row.

## Why

The routing decision is the first thing an agent makes and the one it makes most often, so its text has to be recallable, not merely correct. The current version buries three rules under protected-surface inventories and obligation lists, and states each tier in a different shape — plain by what it excludes, light by a four-part conjunction, full as a leftover clause. One shape per tier, with the inventories kept as supporting detail rather than the definition, makes the rule usable from memory and the guards impossible to miss. This is the campaign's fourth spend of the P-005 criterion: the core needs the tier boundary, not the enumeration style.

## Decision boundary

The routing semantics do not change, so no ADR or PDR is warranted; a rewrite of canonical wording is cheap and local to reverse, which routes it to a dated row in `docs/decisions/log.md` ([C-011](../../docs/constraints/C-011-decision-records-typed.md)) naming `change-tiers.md.tmpl` as the carrier. The text is itself a methodology carrier, so by the red line ([C-013](../../docs/constraints/C-013-core-changes-need-decision.md)) the change is never plain and never light: it is a full change accepted only at human merge. If any reading of the new text would move work between tiers, that is a semantic change and becomes an open question instead of an edit.
