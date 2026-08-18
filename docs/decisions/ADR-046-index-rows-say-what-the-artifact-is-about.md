---
id: ADR-046
type: decision
status: verified
links: [ADR-017, ADR-034, ADR-035, ADR-041, C-004, C-016, CAP-002, CAP-005, P-010]
title: An index row says what its artifact is about; the sentence is seeded, curated, and counted when absent
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-046 — Index rows say what the artifact is about

## Context and problem statement

ADR-041 makes a generated row identify its record, but a reader still needs a first description of what that record is about. The generator can seed that description without claiming that its extraction is final or rewriting adopter-owned prose.

## Decision outcome

**A newly appended row may add one single-line description sentence seeded from the artifact body; the author may curate it, and a row without one is counted rather than failed.** The generator prefers a lede beneath the H1, otherwise the first sentence under the first heading, while skipping headings, tables, lists, quotes, code, HTML blocks, and rules. It reduces links and code spans to labels, declines unsafe candidates such as marker/comment residue, and keeps the shorter ADR-041 row when no safe sentence exists.

Regeneration rewrites nothing that already exists and never backfills descriptions. A dead target drops its row regardless of its description. `clue validate` reports rows that identify their record but say nothing about it as a visible population, not an `Issue`.

**Carrier:** `regenIndex`, index-row validation and backlog reporting, C-016, and CAP-002/CAP-005 guidance.
