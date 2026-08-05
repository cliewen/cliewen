---
id: IC-002
type: imported-change
status: in-progress
links: [IC-001, CAP-003]
title: "Fixture: add-retry-budget — policy hooks"
source-revision: 4f1c9e2
source-location: github.com/example/adopter-repo/openspec/changes/add-retry-budget
---

# IC-002 — Fixture: add-retry-budget — policy hooks

The second half of M-054's worked fixture: an `in-progress` record naming a proof link that is not yet satisfied, proving that status is allowed to declare unproven work rather than being rejected for it.

## Origin

Same source change as [IC-001](IC-001-fixture-inspectable-in-flight-work.md): `example/adopter-repo`'s `changes/add-retry-budget/`, pinned at the same revision and location.

## Intent

Let a call site declare a per-endpoint retry policy (which errors are retryable, the backoff curve) on top of the shared allocator IC-001 provides.

## Design rationale

The source's `design.md` kept policy declaration separate from the allocator so a policy change never requires re-deriving the shared budget math — the same separation of concerns that motivated splitting this into two source tasks in the first place.

## Dependencies

- [IC-001](IC-001-fixture-inspectable-in-flight-work.md): the allocator this record's policy hooks wire into.

## Proof links

| Task | Criterion |
|---|---|
| reuse the existing rehearsal-before-mutation coverage | AC-056 |
| wire the per-endpoint policy hook | AC-999 |

`AC-999` is not yet declared anywhere in this corpus — the source task it names is still pending, which is exactly what `in-progress` is for: `clue validate` does not reject this record for it, unlike a `complete` record naming the same gap.
