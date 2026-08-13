---
id: C-012
type: constraint
status: active
links: [PDR-007, PDR-012, PDR-016, PDR-021, PDR-027, PDR-039, PDR-040, PDR-042, LOG-001]
title: Full changes remain human-accepted while simple integration follows explicit user authority
source: PDR-007, PDR-012, PDR-016, PDR-040, PDR-042, LOG-001, clue-delta steps 1 and 5
enforcement: partial
---

# C-012 — Full changes are human-accepted and simple integration is user-authorized

Every full change branches from the current tip of `main`, never from unaccepted work. An initiating author holds one full Cliewen change in flight at a time; reviewing or helping update an existing pull request does not mint another change or create a global lock. An agent never accepts its own full change: the human-controlled merge commit is the acceptance act, and squash and rebase-and-merge do not preserve the reviewed proposal, implementation, and digest chain.

Simple work remains outside the loop. Route selection does not authorize integration: an agent pushes directly to an integration branch only when the user explicitly authorizes that push and repository permissions allow it; otherwise it follows the repository's requested workflow. A human acting independently may integrate by any mechanism the repository permits, and a repository may impose stricter local rules. If the agent recommended full and the user chose simple, the agent records that one-integration authorization and its risk in the complete Git trailers PDR-042 defines.

Inside the full loop, [PDR-040](../decisions/PDR-040-push-is-durability-ready-is-explicit.md) makes a push durability rather than readiness: every changed turn commits and pushes the change branch, the pull request exists as a draft from first publication, and hosted history is never rewritten. Before marking it ready, the agent automatically reviews the committed candidate under [PDR-012](../decisions/PDR-012-agentic-review-before-publication.md), preferring a context-isolated read-only reviewer and disclosing an in-context fallback; the ready mark binds verification and a clean review pass to the exact hosted head.

Any agent that edits an existing PR becomes its updater for that turn under [PDR-016](../decisions/PDR-016-pr-state-carries-agent-handoffs.md). It fetches and records the hosted head, commits and pushes its repairs with the turn that made them — a repair pushed to a ready PR returns it to draft — verifies the complete repair, obtains a clean review of the resulting commit, and confirms that the PR head equals that reviewed commit before marking it ready again and resolving satisfied findings. A changed head or non-fast-forward rejection requires reconciliation and renewed verification and review. If accepted `main` advances, the updater merges current `main` into the PR branch and repeats verification and review instead of rewriting hosted history. A merged or closed PR stops without pushing — the one turn that ends unpushed — and reports where the work stands. Building on an unmerged change is a blocking open question for the human; once authorized, its committed answer names the base, authorization, and meaning at risk, and the ready acceptance brief repeats it ([PDR-039](../decisions/PDR-039-dependent-changes-carry-authorization.md)).

**Checked by:** branch protection on this repository's `main`, whose repo-local policy requires pull requests for every route, includes administrators, requires `validate`, resolves review conversations, and refuses force pushes and branch deletion. Repository-wide settings preserve full-change history with merge commits. CI requires a completed acceptance brief only when branch history carries a full proposal without a later complete simple override. These checks enforce admission only: they are not acceptance evidence for a criterion ([PDR-027](../decisions/PDR-027-branch-protection-enforces-admission-not-acceptance-evidence.md)).

**Residual:** route recommendation, explicit push authority, semantic growth, and the user's refusal are judgments no forge can derive. Git trailers make an agent-recorded override inspectable but do not prove the conversation or survive every history-rewriting integration mechanism. Inside a full loop, no forge can see an uncommitted edit, skipped review, or private finding. These limits are disclosed rather than converted into claims the wall cannot hold.
