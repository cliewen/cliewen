---
id: PDR-024
type: decision
status: inferred
links: [PDR-020, ADR-036, ADR-037, CAP-003, C-013]
title: Extraction never silently discards acceptance evidence
author: agent
accepted-by: []
---

# PDR-024 — Extraction never silently discards acceptance evidence

## Context and problem statement

ADR-036 limits one JVM executable to one credited acceptance-criterion identity so unrelated evidence cannot cross-product. A source executable can nevertheless carry several source identities. Selecting one identity or deleting the others changes what source evidence says, and a normalizer cannot determine which criterion the assertions genuinely prove.

## Decision outcome

**An extraction normalizer may diagnose evidence ambiguity but never silently discard or select acceptance evidence.** A rehearsal inventories every multi-identity executable with its source location, identities, test behavior, and proposed resolution before mutating source tests.

**The human selects the resolution after the rehearsal.** The accepted conversion either splits the test into executable methods that each carry one criterion and retain the relevant assertions, or selects one primary criterion after reviewing the test behavior and records every other identity as requiring dedicated proof or an explicit `@draft` criterion. A source tag may be removed only as part of that reviewed, recorded resolution.

**No resolution is a finding, not a default.** The ambiguity remains unchanged, is recorded as an open question, and blocks mutation that would alter the test or its evidence. `clue validate` continues to give ambiguous JVM metadata no classified credit under ADR-036; it does not rewrite source.

## Rejected: keep the last tag

Annotation order is not a statement that the final identity is primary. Keeping it by position silently changes source evidence and makes the discarded criteria appear untested rather than making the ambiguity visible.

## Carrier

PDR-020's rehearsal boundary and the canonical `clue-extract` skill require the inventory and human resolution. The OpenSpec mapping applies the same rule to JUnit tags. ADR-036 remains the deterministic attribution rule.
