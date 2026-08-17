---
id: P-016
type: plan
status: active
links: [PDR-044, PDR-045, PDR-003, PDR-006]
title: Decision records are concise, typed, and migration-safe
---

# P-016 — Decision records are concise, typed, and migration-safe

This campaign retires the decision log in favour of subject-typed decision records; keeps future records focused on enduring decisions; migrates older adopter repositories without silently misclassifying their log rows; and compacts this repository's existing decisions in bounded review batches.

The contract and migration mechanism land first. Existing ADRs and PDRs then remain valid until their assigned compaction batch, so each review can distinguish preserved meaning from discarded narrative. A record may be reclassified or split when its subject requires it, but its binding decision cannot silently change.

Milestone numbering continues corpus-global numbering from P-015's M-074.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-075 | **The concise typed-record contract ships and this repository no longer has a decision log.** An accepted decision record states which record types exist, the subject test that routes a decision to exactly one of them, and the compact shape a new or modified record takes; `docs/decisions` then holds only those types, and only future-shaping choices earn a record; validation rejects legacy logs and decision filenames outside the settled types; `clue init`, extraction, generated skills, guidance, and release notes agree; `clue migrate` inventories a legacy log and blocks with guided full-change classification until every durable row has a reviewed destination and the log is removed; this repository's rows are classified, live references repaired, LOG-001 retired, and completed-plan link repairs are limited to link targets. Focused positive and negative evidence covers validation, migration, scaffolding, extraction, and generated parity. | `todo` | |
| M-076 | **ADR-001 through ADR-030 are compacted without changing their decisions.** A transient audit accounts for every removed sentence as retained core, moved binding carrier, or discarded narrative; incorrectly typed records are reclassified and independent decisions are split atomically; verification finds no lost live contract. | `todo` | |
| M-077 | **ADR-031 through ADR-060 are compacted under the same audit and preservation rule.** | `todo` | |
| M-078 | **PDR-001 through PDR-022 are compacted under the same audit and preservation rule.** | `todo` | |
| M-079 | **PDR-023 onward, together with every decision record earlier milestones created, are compact under the same audit and preservation rule.** The plan closes only after re-derived validation and review confirm the whole live decision corpus follows the new contract. | `todo` | |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a correctly typed decision record. Plan adjustments are decisions.
