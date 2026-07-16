---
id: C-012
type: constraint
status: active
links: [PDR-007]
title: Changes root at main; one in flight per author; agents never merge their own PRs or push to main
source: PDR-007, clue-delta steps 1 and 5
enforcement: agent
---

# C-012 — Changes root at main; agents never merge their own

Every change branches from the current tip of `main`, never from unaccepted work. An author holds one change in flight at a time: PR first, next change after. An agent never merges its own PR, never creates a local merge commit into `main`, and never pushes to `main` — the merge is a human act, and until it happens the change is not accepted. Review fixes to an unaccepted change stay on its branch and PR; building on an unmerged change is a blocking open question for the human ([PDR-007](../decisions/PDR-007-review-boundary.md)).

**Promotion trigger:** the hosting plan permits branch protection or rulesets (direct pushes and self-merges become impossible), or CI gains a PR-provenance check that fails `main` when a commit is not reachable from a merged PR — then `enforcement: machine` for the detectable subset. The one-in-flight-per-author rule stays agent-held.
