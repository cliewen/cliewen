---
id: P-020
type: plan
status: active
links: [G-011, PDR-052, ADR-034, P-019, CAP-002]
title: Cliewen's corpus forgets what it has finished with
---

# P-020 — Cliewen's corpus forgets what it has finished with

P-019 built the mechanism for a spike to leave the corpus and deliberately pruned nothing. This campaign applies it here, to the repository that wrote it, before it is pointed at any adopter.

The analysis folder is the corpus's largest population and none of it has ever left. Every plan this repository has run is `completed`, so the completed-plan half of the expiry test does not discriminate — each spike's disposition is a judgement someone makes and records, which is what [PDR-052](../decisions/PDR-052-a-spent-analysis-is-reported-not-retained.md)'s declared carrier exists to force. Retirement is deletion with Git history as the archive ([ADR-034](../decisions/ADR-034-retirement-is-deletion.md)); nothing here is hollowed out into a stub, because that keeps a document's index row and read cost while removing the only part with value.

Batches are ordered by how much a retirement costs a live reader, cheapest first. A spike that a live decision cites as its evidence base is a materially harder case than one only a completed campaign refers to, and this campaign does not prejudge it.

Milestone numbering continues corpus-global numbering from P-019's M-082. The identity ledger still does not cover milestones ([G-006](../goals/G-006-milestone-ids-in-the-ledger.md)), so these are assigned by reading the corpus maximum.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-083 | **The simplification campaign's transient registers that no live decision cites are retired.** Each of them declares the durable artifacts carrying its findings, is deleted rather than emptied, and is named in a `supersedes:` field on exactly one live successor. Every inbound reference is repaired, including a completed plan's link targets, which is the only edit a frozen plan accepts. `clue validate` is green and the retired identities are `retired` in the ledger. | `done` | CH-164: AN-019, AN-020, and AN-021 declared their carriers, `clue migrate` reported all three under `MIG-013` before any deletion, and each was deleted rather than emptied. PDR-033, PDR-016, and PDR-029 each name exactly one retired identity in `supersedes:`. P-013's frontmatter links and one prose link were repaired as link targets only, under C-008's narrow allowance; AN-022's link was repointed to PDR-029 as the validator's supersedes rule directed, and AN-023's pinned inspection rows were deliberately left intact because they are plain text recording what was inspected at `ef93cfa`. `go test ./...`, `go run ./cmd/clue validate`, and the CONTRIBUTING verification block passed before review. |
| M-084 | **Every analysis a live decision cites as its evidence base has a recorded disposition.** The campaign answers whether a decision's cited evidence may be deleted while the decision stands, records that answer as a decision record, and then either retires each spike with its citations repaired or retains it with the reason stated. No spike is retired on carrier grounds alone while a live decision rests on it. | `todo` | |
| M-085 | **Every remaining analysis has a recorded disposition, and the folder holds only spikes that still describe the system.** Findings that are still live truth but sit in the spike folder are moved to the durable artifact that owns them rather than being retired; the rest are retired under the same rule. The campaign closes on a re-derived `clue validate --read-cost` and a `clue migrate` preview whose spent-analysis report is empty or explained. | `todo` | |

## Mutation rules

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else changes only through a declared plan revision backed by a correctly typed decision record. Plan adjustments are decisions.
