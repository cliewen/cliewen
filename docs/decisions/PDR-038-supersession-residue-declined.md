---
id: PDR-038
type: decision
status: verified
links: [P-013, PDR-026, PDR-029, ADR-034, AN-008, AN-022]
title: A surviving superseded record carries no machine edge, declined with its cost stated
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-038 — Supersession of a surviving record stays prose

## Context and problem statement

Retirement deletes an artifact, but decisions that are superseded and retained still need readable history; the corpus records those relationships in prose while `supersedes:` remains unused.

## Decision outcome

**The residue is declined with its cost stated, and `supersedes:` keeps its existing meaning.** It points from a live successor to a deleted artifact and is not widened to surviving superseded decisions; those relationships remain in index rows, successor prose, or amendment notes. The cost is that downstream relationships are not queryable and `clue context` does not reverse-walk them, but removing the unused field would remove the cheap future route without simplifying an active mechanism.

Widening the field is a named door for a successor campaign because it adds reverse-walk machinery, not a simplification. Pattern C is otherwise closed: its remaining residue points to this decision rather than being carried as an unanswered paragraph.

## Rejected: widen or remove the field

Widening would add a live edge, stale-state checks, and reverse traversal; removing it would make a future answer harder without reducing current read or validation work.

## Carrier

This record, P-013's successor door, and AN-022's re-derived evidence carry the disposition. ADR-034's retirement and `supersedes:` semantics, and no skill, hub, guide, or CLI carrier, are changed.
