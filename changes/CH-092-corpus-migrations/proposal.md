---
id: CH-092
type: change
status: open
links: [P-009, M-041, AN-012, ADR-025, ADR-035, ADR-038, CAP-001, CAP-004]
title: Give corpus releases a safe mechanical migration path
---

# CH-092 — Give corpus releases a safe mechanical migration path

## What

This full change serves P-009 milestone M-041 by adding a versioned `clue migrate` command for deterministic corpus and managed-carrier upgrades. It will preview the exact files and line-level transformations before writing, apply only transformations whose preconditions are known, preserve unrelated content, and stop with actionable diagnostics for semantic choices, ambiguous syntax, local edits, or incomplete targets.

The first ordered migration set will cover the historical `reversal-cost` field and the narrowed default status vocabulary, the five generated skills and their Claude mirror, and the immutable upstream CI workflow reference in the thin caller. Applying it will be an atomic coordinated operation; a second run will be a no-op, and a user can resolve a reported case and resume without a hidden state file or a destructive `init` overwrite.

The release and operations contracts, capability criteria and designs, fixture tests, generated guidance, and `[Unreleased]` notes will describe the supported migration boundary and the cases that remain human decisions. `clue init` will remain a non-destructive materializer rather than an updater.

## Why

AN-012 measured an adopter upgrade blocked by dozens of missing `reversal-cost` fields, a status value no longer accepted by the validator, five stale managed skills, and a CI wall whose release reference could not be updated as one safe operation. Manual instructions leave semantic choices and partial writes to each adopter, so a release can change the corpus contract without giving the adopter a truthful recovery path.

## Decision boundary

This change makes the mechanical migration contract explicit and implements the first migrations. It does not infer a reversal-cost class without an explicit user choice, delete ambiguous artifacts, overwrite locally modified managed files, add a general repository configuration file, or make `clue validate` dependent on network access.
