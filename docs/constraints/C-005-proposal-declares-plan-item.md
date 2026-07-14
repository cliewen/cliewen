---
id: C-005
type: constraint
status: active
links: []
title: Every proposal declares its plan item or declares itself plan-less
source: AGENTS.md rule 2, clue-delta skill step 2
enforcement: agent
---

# C-005 — Every proposal declares its plan item or plan-less

A full change's `proposal.md` references the plan item it serves (P/M-IDs in `links`) or states plan-less explicitly; a light change makes the same declaration in its PR description. No fake plan items.

**Promotion trigger:** `clue validate` on a branch requires `/changes/*/proposal.md` to carry a P/M link or the literal plan-less declaration — then `enforcement: machine` for the full tier (the light tier's PR description stays outside the tree, so agent/human there).
