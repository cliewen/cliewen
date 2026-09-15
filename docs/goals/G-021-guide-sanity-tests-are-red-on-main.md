---
id: G-021
type: goal
status: proposed
links: []
title: The guide sanity tests pass again on main
---

# G-021 — The guide sanity tests pass again on main

**Who wants it:** anyone running `go test ./...` in this repository, and every full change's verification loop that relies on a clean baseline to isolate its own regressions (2026-09-15), found while verifying CH-193.

**Why:** `TestSanity_PRBoundaryExplainsAuthorizationAndCIEnforcement` and `TestSanity_AgenticFindingsRequireOperativeViolations` in `cmd/clue/main_test.go` pin literal strings from `docs/decisions/PDR-007-review-boundary.md` and `docs/decisions/PDR-012-agentic-review-before-publication.md` into `guide/what-is-cliewen.md` and `guide/change-loop.md`. Both fail on unmodified `origin/main` (verified in a scratch worktree at commit `d43ed18`, independent of any change on this branch): `guide/what-is-cliewen.md` no longer contains "pull request is the authorization boundary" or "does not require repeating a code review", and `guide/change-loop.md` no longer contains "Release is not a Cliewen route". The recent guide-prose rewrite (<https://github.com/cliewen/cliewen/pull/225>, "center the adopter's skills and corpus") likely rephrased these without updating the pinned strings or the tests that pin them.

This is not itself evidence that the guide misstates the rule — the rewording may still say the same thing in different words — but until someone checks, `go test ./...` is red on `main` for a reason no open change caused, and every subsequent change's verification has to notice and set aside a failure that is not its own.

**Success looks like:**

- `go test ./...` is green on `main` again.
- Either the guide prose is restored to say what the pinned strings expect, or the tests are updated to pin whatever the reworded prose now says the same thing with — a human decision about which, since a test pinning the wrong words either way would misreport a real drift as none.
