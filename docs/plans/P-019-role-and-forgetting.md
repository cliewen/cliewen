---
id: P-019
type: plan
status: active
links: [G-011, PDR-052, ADR-013, ADR-034, CAP-001, CAP-002]
title: A repository knows its role, and a spent analysis can leave
---

# P-019 — A repository knows its role, and a spent analysis can leave

Cliewen's corpus grows but never forgets, and its two audiences are not distinguishable from inside a repository. This campaign makes the role of a repository an observable fact that tooling and skills can branch on, makes an adopter-binding rule checkably present on the surface adopters receive, and gives a spike whose findings have reached durable form a reported route out of the corpus.

Retirement remains what [ADR-034](../decisions/ADR-034-retirement-is-deletion.md) already decided it is: deletion in a reviewed change, with Git history as the archive. Nothing in this campaign deletes on a tool's own judgment, and nothing turns an unexpired analysis into an invalid corpus.

Milestone numbering continues corpus-global numbering from P-018's M-080. The allocator cannot yet issue these numbers, because the identity ledger does not cover milestones ([G-006](../goals/G-006-milestone-ids-in-the-ledger.md)); they are assigned by reading the corpus maximum until that goal is served.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-081 | **A repository declares whether it is Cliewen's source or an adopter, and an adopter-binding rule is checkably shipped.** A role marker under `.clue/` records the role alone; `clue init` materializes it as `adopter` without overwriting an existing declaration, and an undeclared repository is an adopter that no rule blocks. `clue validate` applies a source-repository-only rule that a decision or constraint binding adopter behaviour names a carrier on the shipped surface, and never applies that rule to an adopter's corpus. This repository's local rules are separated from the shared routing text, and the generated skills direct an agent to read the role before applying a rule that differs by repository kind. Focused positive and negative evidence covers declaration, defaulting, the carrier rule, and its non-application. | `todo` | |
| M-082 | **A spent analysis is reported for retirement instead of accumulating.** An analysis declares the durable artifacts now carrying its findings; validation keeps that declaration honest without making its absence a failure. `clue migrate` reports each analysis whose plans are complete and whose declared carriers resolve, emitting a non-blocking notice that never becomes a finding or a write, and the upgrade skill walks the reported list for explicit human approval per document. Retirement stays a reviewed deletion under ADR-034 and migration gains no deletion code. Focused positive and negative evidence covers the declaration, the derivation's two required conditions, and the report-only boundary. | `todo` | |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a correctly typed decision record. Plan adjustments are decisions.
