---
id: PDR-044
type: decision
status: verified
links: [PDR-003, PDR-006, P-016]
title: Decision records are concise and routed by subject
author: agent
accepted-by: Flemming N. Larsen (2026-08-17, conversation)
---

# PDR-044 — Decision records are concise and routed by subject

## Context and problem statement

The decision log and reversal-cost routing obscure whether a durable choice concerns architecture, project process, or implementation, while long records hide the decision itself.

## Decision outcome

Cliewen will record only future-shaping choices as concise MADR-inspired records, routed by subject: ADR for architecture, PDR for project or process, and IDR for implementation. `docs/decisions/log.md` will be removed. Older adopter repositories must classify each durable log entry before migration can continue; routine implementation facts and history are not decision records.

New and modified records contain only enduring context, the outcome, optional considered options when they materially aid understanding, and consequences. Verification treats narrative history, review chronology, carrier inventories, implementation walkthroughs, and repeated rationale as blocking verbosity.

This supersedes PDR-003 and PDR-006's reversal-cost routing and decision-log clauses; their provenance and retention rules remain.

## Consequences

- Shipped scaffolding, validation, migration, extraction, and generated skills must carry the same taxonomy and concision rule.
- Existing records remain valid until touched, while this corpus is compacted in reviewable batches.
- Superseded and rejected decisions remain findable but lose nonessential narrative when compacted.
