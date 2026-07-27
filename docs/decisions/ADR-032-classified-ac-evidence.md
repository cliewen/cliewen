---
id: ADR-032
type: decision
status: inferred
links: [ADR-005, ADR-006, CAP-002, CAP-003, P-007]
title: Acceptance criteria declare classified proof and paired directions
author: agent
accepted-by: []
---

# ADR-032 — Acceptance criteria declare classified proof and paired directions

## Context and problem statement

ADR-006 deliberately left a per-criterion test-type annotation unenforced because a declaration with no consumer would double bookkeeping. M-025 gives it a consumer: `checkACTests` can count evidence by the classes each scenario declares. A single unclassified reference proves only that a string occurred in a test suite; it cannot show that the intended boundary was exercised or that both directions were considered.

## Decision outcome

An acceptance-criterion scenario that opts into classified evidence declares its required proof class on the first non-blank line of the scenario body: `Test-type: Unit`, `Test-type: Integration`, `Test-type: E2E`, or `Test-type: Performance`. This is a design-time declaration reviewed with the scenario, not a path to a test file. An unannotated legacy scenario keeps ADR-006's one-reference rule; every newly added or materially revised scenario declares a test type.

Each declared class requires one `positive` and one `negative` reference. A scenario whose statement itself has one direction may state `Test-type: <class> (single-direction)` and requires one reference in that class. The exemption is explicit because it is a claim about the scenario, not an absence the checker can infer.

Go names carry all three parts: `TestAC042_IntegrationPositive_…` and `TestAC042_IntegrationNegative_…`. JVM files use `@Tag("AC_042")` with `@Tag("integration")` and `@Tag("positive")` or `@Tag("negative")`; Cucumber `.feature` scenarios use the equivalent `@AC-042 @integration @positive` tags. The JVM harvester remains file-level as ADR-009 established. Cucumber tags are scenario-level. `checkACTests` counts only references that match a declared class and reports a missing class or direction.

The purpose taxonomy remains `AC`, `Unit`, `Sanity`, and `Arch`: an AC reference is still a test's one purpose. Test type is runner metadata; positive and negative are evidence direction metadata. Their new consumer is the evidence checker, so they are no longer ignored.

Gatling and similar profiles without native tags use the same stable test-name fallback as Go, not structured comments. A profile that cannot attach an AC ID, type, and direction to a named executable test is unsupported for automated AC evidence until it supplies an equally stable framework-native carrier. This preserves ADR-005's rejection of proximity comments and supersedes ADR-006's deferred `QS<digits>` purpose lane: performance is the `Performance` test type attached to a named AC, rather than a second purpose namespace.

## Consequences

- The validator distinguishes an ordinary legacy reference from classified coverage and can fail an incomplete pair deterministically.
- A test type remains visible where it is cheapest to challenge: the scenario review.
- Existing accepted criteria do not need a mechanical rewrite merely to retain their existing contract.
- Human-only proof and per-criterion draft exemptions remain M-026 work.
