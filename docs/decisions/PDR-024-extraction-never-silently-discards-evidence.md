---
id: PDR-024
type: decision
status: verified
links: [PDR-020, ADR-036, ADR-037, CAP-003, C-013]
title: Extraction never silently discards acceptance evidence
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-024 — Extraction never silently discards acceptance evidence

## Context and problem statement

One source executable or class can carry several criterion identities even though ADR-036 credits one JVM executable to one identity; choosing a destination would change what the source evidence says.

## Decision outcome

**Extraction may diagnose non-attributable evidence but never silently discard, select, or relocate it.** The rehearsal inventories every multi-identity executable and class-level criterion tag with its location, identities, scope, and proposed resolution before mutation.

The human resolves the finding: an executable is split so each method carries one criterion, or one reviewed primary is retained while other identities receive dedicated-proof or `@draft` dispositions; a class tag moves only to an attributable existing or newly scoped test. An out-of-scope test may be `@draft` only with its source location and plan door recorded. Until resolved, the evidence remains unchanged and mutation is blocked; `clue validate` does not credit ambiguous or class-level JVM metadata.

## Rejected: discard a tag without executable credit

Annotation order does not establish which identity is primary, and removing a class-level or multi-identity tag hides the source evidence instead of making its non-attribution reviewable.

## Carrier

PDR-020's rehearsal boundary, the canonical `clue-extract` skill, and the OpenSpec mapping require inventory and human resolution; ADR-036 remains the deterministic attribution rule.
