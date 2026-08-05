---
id: C-012
type: constraint
status: active
links: [PDR-007, PDR-011, PDR-012, PDR-016, PDR-021, LOG-001]
title: Changes are reviewed locally, root at main, and remain human-merged with preserved history
source: PDR-007, PDR-011, PDR-012, PDR-016, LOG-001, clue-delta steps 1 and 5
enforcement: partial
---

# C-012 — Changes are locally reviewed, root at main, and remain human-merged

Every change branches from the current tip of `main`, never from unaccepted work. An initiating author holds one Cliewen change in flight at a time: PR first, next light or full change after. Reviewing or helping update an existing PR does not mint another change or create a global lock; independent authors and plain changes under [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) may proceed from accepted `main`. An agent never merges its own PR, never creates a local merge commit into `main`, and never pushes to `main` — the merge is a human act, and until it happens the change is not accepted. A full Cliewen change is accepted through the supported merge-commit mode; squash and rebase-and-merge do not preserve the reviewed proposal, implementation, and digest chain.

Before publishing a Cliewen change, the agent automatically reviews the committed candidate under [PDR-012](../decisions/PDR-012-agentic-review-before-publication.md), preferring a context-isolated read-only reviewer and disclosing an in-context fallback. Every review of an existing PR names its hosted head; a clean result applies only to that commit, and actionable findings become unresolved hosted review conversations where supported. A reviewer unable to publish a resolvable finding reports the PR not merge-ready and discloses the enforcement gap.

Any agent that edits an existing PR becomes its updater for that turn under [PDR-016](../decisions/PDR-016-pr-state-carries-agent-handoffs.md). It fetches and records the hosted head, commits and verifies the complete repair, obtains a clean review of the resulting commit, pushes without force, and confirms that the PR head equals that reviewed commit before resolving satisfied findings. A changed head or non-fast-forward rejection requires reconciliation and renewed verification and review. If accepted `main` advances after publication, the updater merges current `main` into the PR branch and repeats verification and review instead of rewriting hosted history; before first publication, rebasing remains allowed. A merged or closed PR stops with local work reported as unpublished. A human-requested local stopping point is preserved work but is explicitly incomplete and not mergeable. Building on an unmerged change is a blocking open question for the human ([PDR-007](../decisions/PDR-007-review-boundary.md)).

**Checked by:** branch protection on this repository's `main`, which makes the detectable half impossible rather than merely forbidden: pull requests are required, administrators are included, `validate` is a required status check, review conversations must be resolved before merge, and force pushes and branch deletion are refused. The repository-wide pull-request settings hold the merge shape, allowing merge commits and disabling squash and rebase-and-merge, so the reviewed proposal, implementation, and digest chain survives the merge ([PDR-016](../decisions/PDR-016-pr-state-carries-agent-handoffs.md)). CI additionally requires a completed acceptance brief on a full change.

**Residual:** everything that happens before the push. No forge can see an uncommitted local edit, a review that was skipped, a finding an agent kept to itself, or whether a context-isolated semantic review actually occurred; nor can it count how many changes one author has in flight, or tell an initiating author's second change from a review of somebody else's. The cost is that the wall stops an agent from *accepting* its own work but not from *proposing* work it never really reviewed — which is why the human merge is a judgment about the change rather than a formality, and why an agent's disclosure of what it could not enforce is part of the report rather than an optional courtesy.
