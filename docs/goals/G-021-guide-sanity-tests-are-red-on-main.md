---
id: G-021
type: goal
status: accepted
links: []
title: The guide sanity tests pass again on main
---

# G-021 — The guide sanity tests pass again on main

**Who wants it:** anyone running `go test ./...` in this repository, and every tracked change's verification loop that relies on a clean baseline to isolate its own regressions (2026-09-15), found while verifying CH-193.

**Why:** `TestSanity_PRBoundaryExplainsAuthorizationAndCIEnforcement` and `TestSanity_AgenticFindingsRequireOperativeViolations` in `cmd/clue/main_test.go` pin wording from the review-boundary decisions ([PDR-007](../decisions/PDR-007-review-boundary.md), [PDR-012](../decisions/PDR-012-agentic-review-before-publication.md)) into the guide. A guide rewrite (<https://github.com/cliewen/cliewen/pull/225>) rephrased that wording, and both tests went red on `main`, so every change's verification had to set aside a failure that was not its own.

**How it is met:** commit `9ed9977` re-pinned both tests to the guide's current wording after checking that it still states the same rule, and `go test ./...` is green on `main`.

**Success looks like:**

- `go test ./...` is green on `main` again.
- Either the guide prose is restored to say what the pinned strings expect, or the tests are updated to pin whatever the reworded prose now says the same thing with — a human decision about which, since a test pinning the wrong words either way would misreport a real drift as none.
