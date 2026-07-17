---
id: CH-021
type: change
status: open
links: [P-002, M-005]
title: Plan revisions may ride with their implementation; M-005 criterion revised and closed
---

# CH-021 — Plan-revision bundling, and the M-005 close-out it unblocks

## What

Three pieces, one decision thread:

1. **PDR-008**: a declared semantic plan revision may ride with its implementing change when all of four conditions hold — the PR explicitly declares the revision, the correctly typed decision record backs it, the PR calls the revision out for deliberate human approval, and an explicit objection reverts the revision (milestone stays open) without blocking the rest of the change. The default vehicle remains a dedicated plan PR. The `clue-plan` skill is revised to carry the allowance (canonical + embedded template).
2. **M-005 criterion revision + close**: the exit criterion's "install→first-green-validate-in-30-min criteria tested" wording predates the AC-001 split (decision-log row 2026-07-17); it becomes "criteria tested, the 30-minute install→green promise held as quality scenario QS-002". With the revision in place, M-005 flips `done` with CH-020 (PR #19, merged 2026-07-17) as evidence. This change *is* the dedicated plan-revision PR, so the revision is compliant under the existing rule; PDR-008 governs future bundling.
3. **ADR-017 scope note**: in scaffolded repos the register holds repo-local conventions and the methodology conventions no versioned skill carries; skill-carried rules are not registered. This makes CH-020's seeded-C-001 behavior explicit doctrine instead of a log-row inference.

## Why

CH-020 bundled the M-005 revision into its feature PR and had to withdraw it: `clue-plan` unambiguously requires a separate PR, which review showed is stricter than the acceptance boundary the human actually wants. The direction of PDR-008 was accepted by Flemming N. Larsen on 2026-07-17 in the PR #19 review conversation ("I accept the methodology proposal"); this change turns that acceptance into the record and the carrier, and completes the M-005 bookkeeping the withdrawal left open.

## Plan item

Serves P-002 / M-005 (the revision and close-out); the PDR itself is process infrastructure the same campaign surfaced.
