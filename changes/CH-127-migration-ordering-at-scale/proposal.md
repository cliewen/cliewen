---
id: CH-127
type: change
status: open
links: [P-012, M-061, CAP-003]
title: Prove extraction ordering, source-work preservation, and a pinned release at assessment scale
---

# CH-127 — Prove extraction ordering, source-work preservation, and a pinned release at assessment scale

## Why

P-012/M-061 is still `doing`. CH-126 closed its command-scale half only: AC-128 states in its own scenario that the fixture does not claim extraction ordering, source-work preservation, or pinned-release evidence, and CAP-003's status note repeats that limit. Three of the milestone's named path steps therefore remain unproven at the assessment's order of magnitude, and the campaign cannot close while its last milestone's text is broader than its evidence.

## What changes

One acceptance criterion, AC-129, and its focused positive and negative evidence in `internal/migration`, reusing the assessment-scale fixture CH-126 built.

The evidence runs the migration path in order through a `clue` binary stamped with a pinned release, against a fixture whose installed skills carry that same stamp — so the release-drift comparison that a `dev` build skips actually runs. The ordered path is: report-only rehearsal while the target does not yet exist, human-approved mutation, then verification against the byte-identical pins the rehearsal wrote. Source-work preservation is asserted for every in-flight imported change the fixture carries: pinned origin, intent and rationale sections, and proof links resolving to live criteria, with ADR-050's complete-versus-in-progress gate proven from both sides at that size.

No new validator rule is added. Ordering is proven as the sequence of verdicts the existing commands give at each phase, which is what P-012's "no milestone reopens P-011's mechanisms" boundary requires.

## Acceptance

CAP-003 carries a criterion for the ordered, pinned-release, source-work-preserving path at assessment scale, with executable positive and negative evidence that never runs a source suite or names a production adopter. M-061 has no unproven remainder, so this change's digest closes P-012 and designates P-013.
