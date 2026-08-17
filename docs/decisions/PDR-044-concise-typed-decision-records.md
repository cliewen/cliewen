---
id: PDR-044
type: decision
status: verified
links: [PDR-003, PDR-006, P-016]
title: Decision-record simplification runs as a staged campaign
author: agent
accepted-by: Flemming N. Larsen (2026-08-17, conversation)
---

# PDR-044 — Decision-record simplification runs as a staged campaign

## Context and problem statement

Replacing the decision log, changing record routing, migrating adopters, and compacting the existing corpus cannot be reviewed safely as one rewrite.

## Decision outcome

P-016 will deliver the simplification as a staged campaign. Its first implementation milestone will decide and ship the complete decision-record contract, migration path, and this repository's log conversion atomically. Later milestones will compact the existing ADR and PDR corpus in bounded batches without changing binding meaning.

## Consequences

- PDR-003 and PDR-006 remain the active routing contract until M-075 is accepted.
- Each implementing change records its consequential choices and updates their complete live-carrier inventory.
- Adopter behavior does not change in this planning change.
