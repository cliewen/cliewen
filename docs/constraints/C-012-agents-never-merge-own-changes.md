---
id: C-012
type: constraint
status: active
links: [PDR-007, PDR-011, PDR-012, PDR-016, PDR-021, LOG-001]
title: Changes are reviewed locally, root at main, and remain human-merged with preserved history
source: PDR-007, PDR-011, PDR-012, PDR-016, LOG-001, clue-delta steps 1 and 5
enforcement: agent
---

# C-012 — Changes are locally reviewed, root at main, and remain human-merged

Every change branches from the current tip of `main`, never from unaccepted work. An initiating author holds one Cliewen change in flight at a time: PR first, next light or full change after. Reviewing or helping update an existing PR does not mint another change or create a global lock; independent authors and plain changes under [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) may proceed from accepted `main`. An agent never merges its own PR, never creates a local merge commit into `main`, and never pushes to `main` — the merge is a human act, and until it happens the change is not accepted. A full Cliewen change is accepted through the supported merge-commit mode; squash and rebase-and-merge do not preserve the reviewed proposal, implementation, and digest chain.

Before publishing a Cliewen change, the agent automatically reviews the committed candidate under [PDR-012](../decisions/PDR-012-agentic-review-before-publication.md), preferring a context-isolated read-only reviewer and disclosing an in-context fallback. Every review of an existing PR names its hosted head; a clean result applies only to that commit, and actionable findings become unresolved hosted review conversations where supported. A reviewer unable to publish a resolvable finding reports the PR not merge-ready and discloses the enforcement gap.

Any agent that edits an existing PR becomes its updater for that turn under [PDR-016](../decisions/PDR-016-pr-state-carries-agent-handoffs.md). It fetches and records the hosted head, commits and verifies the complete repair, obtains a clean review of the resulting commit, pushes without force, and confirms that the PR head equals that reviewed commit before resolving satisfied findings. A changed head or non-fast-forward rejection requires reconciliation and renewed verification and review. If accepted `main` advances after publication, the updater merges current `main` into the PR branch and repeats verification and review instead of rewriting hosted history; before first publication, rebasing remains allowed. A merged or closed PR stops with local work reported as unpublished. A human-requested local stopping point is preserved work but is explicitly incomplete and not mergeable. Building on an unmerged change is a blocking open question for the human ([PDR-007](../decisions/PDR-007-review-boundary.md)).

**Promotion trigger:** the hosting plan permits branch protection or rulesets (direct pushes, self-merges, unsupported merge modes, and merge with unresolved review conversations become impossible), or CI gains a PR-provenance check that fails `main` when a commit is not reachable from a merged PR — then `enforcement: machine` for the detectable subset. An executable preflight that can inspect both local Git state and the hosting provider may promote the clean exact-head handoff; CI alone can never see uncommitted local edits, undisclosed findings, or prove that a context-isolated semantic review occurred. The review loop and one-in-flight-per-initiating-author rule stay agent-held.
