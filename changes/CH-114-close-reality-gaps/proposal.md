---
id: CH-114
type: change
links: [P-010, M-050]
title: Close the four reported reality gaps
---

# CH-114 — Close the four reported reality gaps

## What

`clue validate --reality-gaps` reports four capabilities as contradicted:

- CAP-001, contradicted by AN-011 (through AC-036)
- CAP-002, contradicted by AN-011
- CAP-003, contradicted by AN-011 and AN-013
- CAP-006, contradicted by AN-007

Every one of these findings already has a consuming change on `main`: P-008/M-032 (CH-078) repaired the live-carrier inventory AN-011's F1–F3 found stale; ADR-036 withdrew the ArchUnit delegation AN-013's F2 second half found never installed; CH-039 shipped the exact-handoff rule AN-007 proposed. In each case the capability text the analysis originally contradicted already states the corrected claim — verified directly against `docs/capabilities/CAP-001-onboarding/design.md`, `docs/capabilities/CAP-002-validate/README.md`, `docs/capabilities/CAP-003-extract/criteria.md` (AC-054), `docs/capabilities/CAP-006-collaborative-handoffs/criteria.md` (AC-040, AC-041), and `docs/decisions/ADR-009-ac-id-namespaces.md`. What is stale is only the `reality: contradicted` marker and its capability-link edge on AN-007, AN-011, and AN-013, which `internal/corpus/reality.go` derives the report from.

This change retires those three markers, closing plan P-010's milestone M-050 with the reality-gap report at zero. No capability README, criteria, or design file needs a content change: the gap was already closed by the changes named above, only the bookkeeping never caught up.

## Why

M-050's exit criterion is that the reality-gap report reaches zero or names what remains and why. All four current gaps are stale: `--reality-gaps` derives its report from live `links:` edges (`internal/corpus/reality.go:100-130`), so a resolved incident keeps reporting until the marker is removed. Leaving it reporting after the fix already shipped defeats the report's purpose — it stops distinguishing a genuine unresolved contradiction from bookkeeping debt.

## AN-013's remaining findings

AN-013's F1 (accepted-ness is not corpus data) and F3 (unqualified forge references) are untouched by this change. They never carried a `reality:` edge — AN-013's own prose (`docs/analysis/AN-013-distributed-work-and-evidence-boundaries.md:76`) explains the marker was scoped only to F2's second half, the ArchUnit delegation — and remain the analysis P-011's M-052–M-056 are scoped against. Removing the `reality: contradicted` field does not touch `status: active` or the F1/F3 findings; the analysis stays live and linked from P-011.
