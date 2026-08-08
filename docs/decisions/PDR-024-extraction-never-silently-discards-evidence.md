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

ADR-036 limits one JVM executable to one credited acceptance-criterion identity so unrelated evidence cannot cross-product. A source executable can nevertheless carry several source identities, and a class-level identity can describe structural evidence that belongs to no individual inherited test method. Selecting, moving, or deleting either form changes what source evidence says, and a normalizer cannot determine its correct executable destination.

## Decision outcome

**An extraction normalizer may diagnose non-attributable evidence but never silently discard, select, or relocate acceptance evidence.** A rehearsal inventories every multi-identity executable and class-level criterion tag with its source location, identities, behavior or structural scope, and proposed resolution before mutating source tests.

**The human selects the resolution after the rehearsal.** A multi-identity executable either splits into executable methods that each carry one criterion and retain the relevant assertions, or keeps one reviewed primary criterion while every other identity receives a recorded dedicated-proof or `@draft` disposition. A class-level criterion tag moves only to an existing executable that actually proves it, or to a new attributable test introduced in scoped implementation work. When that test work is explicitly out of scope, the converted criterion may be `@draft` only with the original class-level source location and a named plan door recorded in the extraction report. A source tag may be removed only as part of one of these reviewed, recorded resolutions.

**No resolution is a finding, not a default.** The source evidence remains unchanged, is recorded as an open question, and blocks mutation that would alter it. `clue validate` continues to give ambiguous or class-level JVM metadata no classified credit under ADR-036; it does not rewrite source.

## Rejected: discard a tag because no executable credit is available

Annotation order is not a statement that the final identity is primary, and an architectural class tag can state a structural claim without naming one inherited method. Removing either form makes the discarded criterion appear untested rather than making the non-attribution visible.

## Carrier

PDR-020's rehearsal boundary and the canonical `clue-extract` skill require the inventory and human resolution. The OpenSpec mapping applies the same rule to JUnit tags. ADR-036 remains the deterministic attribution rule.
