---
id: ADR-033
type: decision
status: verified
links: [ADR-032, ADR-007, ADR-025, CAP-002, CAP-006, PDR-017, P-007]
title: Human proof class, a per-criterion draft exemption, and derived coverage
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-033 — Human proof class, a per-criterion draft exemption, and derived coverage

## Context and problem statement

Classified code evidence cannot represent a criterion that is permanently verified by human judgment, while using a capability-wide `draft` status for one not-yet-proven criterion confuses an unfinished proof with an unfinished capability. The corpus needs vocabulary for both cases without weakening retirement or review.

## Decision outcome

**The proof vocabulary gains `Test-type: Human`, a per-criterion `@draft` token exempts one live criterion from active-file evidence checks, and coverage is a derived report rather than a committed registry.** `Human` is satisfied by the acceptance brief already required by PDR-017, takes no direction pair, and is reviewed as a claim about how the criterion is always verified. `@draft` keeps a criterion alive and referenced while skipping its missing-evidence and pair checks; removing the token when evidence lands returns it to ordinary checking. The capability and criteria-file lifecycle remains governed by ADR-025.

`clue validate --coverage` derives `covered`, `partial`, or `gap` from current criteria and evidence state and writes nothing, so no stale coverage registry can exist.

**Carrier:** the `checkACTests` test-type and criterion-token parsing, `runValidate --coverage`, the acceptance-brief line in the `clue-delta` and `clue-verify` sources and generated copies, and CAP-002/CAP-006 criteria.
