---
id: CH-075
type: change
status: open
links: [M-028]
title: Provenance is bounded, and reality gets an edge back
---

# CH-075 — Provenance is bounded, and reality gets an edge back

M-028 of [P-007](../../docs/plans/P-007-core-hardening.md). Today `clue validate` reports every non-decision artifact carrying `provenance: inferred` as one ever-growing population, while omitting decisions whose `status: inferred` expresses the same unverified meaning. The count gives no user a next action, and an active capability can depend on expensive-to-reverse inferred meaning without the judge objecting. At the same time, an analysis can record that the corpus was green while reality was wrong, as AN-007 did, but no stated shape connects that incident back to the capability or criterion whose claim failed and no derived view lists those gaps.

This change bounds both omissions. ADR-035 will route inferred meaning by reversal cost using a machine-readable classification: expensive-to-reverse inferred artifacts and decisions are blockers when an active capability depends on them, while cheap inferred artifacts may remain legitimately deferred. The validator will enforce that activation boundary and its successful output will report the actionable blockers, including inferred decisions, instead of the monotonic all-inferred count. The same decision will define a machine-visible incident-analysis marker: an analysis recording “green met wrong” links the failed capability or criterion as well as the rule carriers that failed to prevent the incident. A derived CLI listing will roll those incident edges up to the affected capabilities without creating another committed registry.

This is a [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) red-line change because it changes what a green `clue validate` asserts. It carries ADR-035, routed as architecture under [C-011](../../docs/constraints/C-011-decision-records-typed.md) because the decision defines corpus fields, graph semantics, and the validator contract. The V3 production door remains explicitly out of scope.
