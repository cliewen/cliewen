---
id: CH-146-tasks
type: tasks
status: open
links: [CH-146-proposal]
title: Tasks
---

# Tasks — CH-146

- [x] Run the M-071 read-cost report at the proposal head.
- [ ] Classify every live multi-document artifact and every default slice over the eight-artifact budget.
- [ ] Split each repairable live multi-document artifact by primary consumer, or record a durable accepted reason for every count that remains.
- [ ] Repair only the reported default slice-budget findings, retaining all required durable links.
- [ ] Add or revise acceptance criteria and focused evidence for any behavior the repair changes, then update plan bookkeeping, indexes, and the user-facing `[Unreleased]` entry.
- [ ] Run the applicable local verification and prepare the change for review.

## Baseline — 2026-08-11

`go run ./cmd/clue validate --read-cost` reports one multi-document artifact and 33 identity slices over the eight-artifact budget. `AN-022` has two rendered H1 documents: its scoring of the non-prose surface and its Pattern C determination. Both serve M-064's reader, so the classification will assess whether the stated primary-consumer exception is sufficient rather than splitting the artifact.

The over-budget identities are ADR-038, ADR-039, ADR-040, ADR-041, ADR-042, ADR-044, ADR-046, ADR-049, ADR-055, AN-008, AN-011, AN-012, AN-013, AN-014, AN-015, AN-017, AN-018, AN-020, AN-021, AN-022, C-012, CAP-003, CAP-004-design, PDR-017, PDR-019, PDR-020, PDR-021, PDR-023, PDR-029, PDR-032, PDR-033, PDR-035, and PDR-040.
