---
id: C-012
type: constraint
status: active
links: [PDR-007, PDR-011, LOG-001]
title: Changes root at main; Cliewen work is one in flight; agents never merge or push main
source: PDR-007, PDR-011, LOG-001, clue-delta steps 1 and 5
enforcement: agent
---

# C-012 — Changes root at main; agents never merge their own

Every change branches from the current tip of `main`, never from unaccepted work. An author holds one Cliewen change in flight at a time: PR first, next light or full change after. Plain changes under [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) do not consume that slot and may proceed independently from accepted `main`. An agent never merges its own PR, never creates a local merge commit into `main`, and never pushes to `main` — the merge is a human act, and until it happens the change is not accepted.

Before reporting any change ready, the agent commits every intended edit, requires a clean worktree, pushes the branch, and confirms that the ready hosted PR's head branch and SHA equal the locally verified branch and `HEAD`. Review fixes to an unaccepted change are committed and pushed to that change's existing branch and PR, then verification and the hosted-head check repeat. A human-requested local stopping point is preserved work but is explicitly incomplete and not mergeable. Building on an unmerged change is a blocking open question for the human ([PDR-007](../decisions/PDR-007-review-boundary.md)).

**Promotion trigger:** the hosting plan permits branch protection or rulesets (direct pushes and self-merges become impossible), or CI gains a PR-provenance check that fails `main` when a commit is not reachable from a merged PR — then `enforcement: machine` for the detectable subset. An executable preflight that can inspect both local Git state and the hosting provider may promote the clean exact-head handoff; CI alone can never see uncommitted local edits. The one-in-flight-per-author rule stays agent-held.
