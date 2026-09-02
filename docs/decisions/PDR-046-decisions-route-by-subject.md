---
id: PDR-046
type: decision
status: verified
links: [P-016, PDR-003, PDR-006, PDR-044, PDR-045, C-011, C-013]
title: Future-shaping decisions route by subject into concise records
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
supersedes: [LOG-001]
---

# PDR-046 — Future-shaping decisions route by subject into concise records

## Context

A reversal-cost gate sends small but durable choices into a shared log, where the record's subject and continuing force are hard to see. Long-form ADRs and PDRs then mix enduring rationale with chronology and implementation narrative.

## Decision

Only a future-shaping choice earns a decision record, and its subject selects exactly one type: an ADR for software or corpus architecture, a PDR for project or process, and an IDR for implementation. Reversal cost no longer routes a decision, routine facts and history are not records, and the legacy decision log is retired.

A new or modified record keeps the common provenance frontmatter and a compact body: enduring context, the decision, optional alternatives when they materially explain the choice, and optional consequences. Rejected future-shaping choices route by the same subject test. Triggering incidents, chronology, review history, carrier inventories, and implementation walkthroughs stay in analysis, the change workspace, the pull request, or Git history.

A legacy log is migrated only through a reviewed full change that inventories every row, creates or amends the routed destination for each future-shaping choice, explicitly accounts for narrative that is discarded, repairs live references, and removes the log. `clue migrate` reports the inventory and blocks rather than guessing any classification.

This decision supersedes PDR-003 and PDR-006's reversal-cost and decision-log clauses; their provenance, human-approval, retention, and subject definitions survive. Existing ADRs and PDRs remain valid until P-016's bounded compaction milestones touch them.

## Consequences

The live carriers are the decisions index and C-011; canonical decision guidance rendered into all lifecycle skills; the extraction target contract and MADR mapping; scaffolded decision guidance; the public corpus and change-loop guidance; validator taxonomy and migration registry behavior; capability criteria and explanations; and the adopter-facing changelog entry. They move together in this change.
