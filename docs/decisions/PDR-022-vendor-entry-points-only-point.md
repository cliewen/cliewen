---
id: PDR-022
type: decision
status: verified
links: [G-001, CAP-001, ADR-018, ADR-031, PDR-019, C-006, C-013]
title: A scaffolded vendor entry point may exist, and it may only point at AGENTS.md
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-022 — A scaffolded vendor entry point may exist, and it may only point at AGENTS.md

## Context and problem statement

`AGENTS.md` is the cross-agent routing hub, but some assistants load a vendor-specific entry file instead. A missing bridge lets an assistant see skills while missing the routing that classifies work, and adding a second rule-bearing hub would make carriers diverge.

## Decision outcome

**The scaffold may emit a vendor entry point whose only methodology content is a pointer to `AGENTS.md` and an explanation of that pointer.** No rule lives in the vendor file, and adding one is a methodology-carrier change governed by PDR-019. The emitted file is adopter-owned: `clue init` never overwrites it and `clue migrate` never rewrites it, so adopter-specific instructions remain below the pointer.

A vendor qualifies by published evidence that its normal loading behavior otherwise makes the hub unreachable; Claude Code is the current case. `clue migrate` reports a missing or unrouted entry point and repairs neither; rerunning the non-destructive `clue init` is the available missing-file remedy, while an adopter-authored file is not migration-owned. No opt-out is decided here because the corpus cannot distinguish non-use of the vendor from an adopter who needs the warning.

CAP-001, the scaffolded `CLAUDE.md` pointer, ADR-018, and migration behavior carry the boundary. The decision rejects both vendor-specific methodology and an unbounded list of popularity-based entry points.
