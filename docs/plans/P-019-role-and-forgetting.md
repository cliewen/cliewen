---
id: P-019
type: plan
status: completed
links: [G-011, ADR-062, PDR-052, ADR-013, ADR-034, CAP-001, CAP-002]
title: A repository knows its role, and a spent analysis can leave
---

# P-019 — A repository knows its role, and a spent analysis can leave

Cliewen's corpus grows but never forgets, and its two audiences are not distinguishable from inside a repository. This campaign makes the role of a repository an observable fact that tooling and skills can branch on, makes an adopter-binding rule checkably present on the surface adopters receive, and gives a spike whose findings have reached durable form a reported route out of the corpus.

Retirement remains what [ADR-034](../decisions/ADR-034-retirement-is-deletion.md) already decided it is: deletion in a reviewed change, with Git history as the archive. Nothing in this campaign deletes on a tool's own judgment, and nothing turns an unexpired analysis into an invalid corpus.

Milestone numbering continues corpus-global numbering from P-018's M-080. The allocator cannot yet issue these numbers, because the identity ledger does not cover milestones ([G-006](../goals/G-006-milestone-ids-in-the-ledger.md)); they are assigned by reading the corpus maximum until that goal is served.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-081 | **A repository declares whether it is Cliewen's source or an adopter, and an adopter-binding rule is checkably shipped.** A role marker under `.clue/` records the role alone; `clue init` materializes it as `adopter` without overwriting an existing declaration, and an undeclared repository is an adopter that no rule blocks. `clue validate` applies a source-repository-only rule that a decision or constraint binding adopter behaviour names a carrier on the shipped surface, and never applies that rule to an adopter's corpus. This repository's local rules are separated from the shared routing text, and the generated skills direct an agent to read the role before applying a rule that differs by repository kind. Focused positive and negative evidence covers declaration, defaulting, the carrier rule, and its non-application. | `done` | CH-163: AC-153 covers the marker's declaration, the adopter default for an undeclared repository, rejection of an unknown or malformed role at both the reader and its validator consumer, `clue init` materializing an adopter marker and never overwriting a declared one, and migration reporting `MIG-012` as a notice that never blocks or writes; AC-154 covers the source-only carrier rule passing on a canonical skill source or scaffolded template, failing without one, and never judging an adopter or undeclared corpus. `internal/role` holds the reader, `clue init` materializes `adopter` through the existing non-overwriting path, the routing hub separates the repository-local layer under its own heading, and the generated skills direct an agent to read the role. `go test ./...`, `go run ./cmd/clue validate`, and the CONTRIBUTING verification block passed before review. |
| M-082 | **A spent analysis is reported for retirement instead of accumulating.** An analysis declares the durable artifacts now carrying its findings; validation keeps that declaration honest without making its absence a failure. `clue migrate` reports each analysis whose plans are complete and whose declared carriers resolve, emitting a non-blocking notice that never becomes a finding or a write, and the upgrade skill walks the reported list for explicit human approval per document. Retirement stays a reviewed deletion under ADR-034 and migration gains no deletion code. Focused positive and negative evidence covers the declaration, the derivation's two required conditions, and the report-only boundary. | `done` | CH-163: AC-155 covers an honest `carried-by` passing, an analysis declaring nothing passing unchanged, and rejection of the field on a non-analysis artifact, an empty list, an unresolvable ID, a self-reference, and an ID naming another analysis; AC-156 covers the derivation requiring both a completed plan and a resolving carrier, and `MIG-013` reporting a spent analysis as a notice while `Apply` leaves the file byte-identical. `clue validate` gained no expiry rule and `internal/migrate` gained no deletion code. `go test ./...`, `go run ./cmd/clue validate`, and the CONTRIBUTING verification block passed before review. |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a correctly typed decision record. Plan adjustments are decisions.
