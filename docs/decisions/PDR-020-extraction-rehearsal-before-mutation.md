---
id: PDR-020
type: decision
status: verified
links: [P-008, CAP-003, AN-004, AN-005, ADR-008, PDR-019, C-013]
title: Extraction rehearses before it mutates
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-020 — Extraction rehearses before it mutates

## Context and problem statement

Brownfield extraction converts uncertain source meaning into a repository's durable system of record and deletes the parallel source corpus in the same pull request. AN-003 showed that source carriers can disagree, builds can depend on prepared environments, and adoption changes governance. A proposed extraction needs a reviewable account of what it would convert and delete before those changes become irreversible. How does the extraction workflow expose that account without treating an incomplete conversion as a second system of record?

## Decision outcome

**Every extraction begins with a mandatory report-only rehearsal after its full change is proposed and before it mutates the target.** The rehearsal writes a report under that change's `/changes/CH-xxx-slug/` workspace and changes no target source corpus, Cliewen `/docs` corpus, tests, routing, or hosted state.

The report inventories source formats and entry points; proposed artifact mappings; preserved and minted IDs; confidence and reversal cost; test-purpose work; instruction conflicts; planned deletions; and named plan doors. An unresolved conflict becomes an `open-questions.md` entry and stops the extraction before mutation.

Only explicit human direction may begin the existing full extraction change's mutate phase. That phase digests the rehearsal into the durable extraction report in `/docs/analysis`, carries out the accepted conversion, and deletes the transient workspace and parallel source corpus only when the resulting full change is ready for review.

The rehearsal is a checkpoint, not a new source of truth or a separate change: it is transient branch-local evidence that lets a human decide whether to authorize mutation. The full extraction still owns all conversion decisions, the final durable report, validation, and the human merge boundary.

**Carrier:** `clue-extract` and its generated and scaffolded copies instruct the checkpoint; CAP-003 and AC-056 define its capability and acceptance boundary; guide and README adoption explanations make the public workflow consistent; generator tests reject a rendered skill that omits or weakens the checkpoint.

### Rejected: mutate first and document the extraction afterward

The durable extraction report is valuable history, but cannot let a human inspect proposed ID mappings, deletions, and conflicts before the target corpus and governance are changed. Reconstructing that account after mutation is not reversible rehearsal.

### Rejected: make the rehearsal a separate analysis change

The report is specific to the pending extraction's source state and branch. Moving it to a different change would make its freshness and authorization ambiguous, while leaving a durable report before any decision to adopt was made.
