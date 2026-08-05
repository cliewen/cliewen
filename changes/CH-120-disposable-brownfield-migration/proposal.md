---
id: CH-120
type: change
status: open
links: [P-011]
title: Prove brownfield migration in disposable fixtures
---

# CH-120 — Prove brownfield migration in disposable fixtures

Serves [P-011](../../docs/plans/P-011-truthful-brownfield-migration.md) M-056.

## What

Create disposable numeric-archive and opaque-identifier brownfield fixture sources and use the canonical extraction workflow to exercise report-only rehearsal, explicitly approved mutation, target validation, parity generation, and the required failure paths. The fixtures will demonstrate the complete migration contract without referring to a production adopter and will keep Cliewen's deterministic checks distinct from tests that belong to the source fixture itself.

## Why

M-052 through M-055 establish the ledger, parity, imported-work, and carrier-reconciliation components independently. M-056 must show that those components compose into a safe source-deletion workflow for both numeric and opaque criterion identity families. A disposable fixture makes the claim repeatable without exposing a real adopter or treating a one-off migration as a product test.

## Scope

- Establish the report-only extraction rehearsal for each disposable fixture before any target mutation.
- Add the end-to-end fixture material, extraction outputs, and focused tests needed to prove the clean and mandatory failure paths.
- Add or revise acceptance criteria before implementation if the current criteria do not state the end-to-end fixture contract.
- Update the canonical extraction guidance and durable CAP-003 documentation with the proven workflow.
- Record the M-056 result, release note, and plan bookkeeping in the digest.

## Boundary

This change does not make `clue validate` run source-test suites, use Git history as an identity registry, accept arbitrary prose as an opaque identifier, or retain a permanent parallel source corpus.
