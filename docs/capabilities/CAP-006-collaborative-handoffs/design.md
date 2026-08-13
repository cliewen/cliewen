---
id: CAP-006-design
type: design
status: active
links: [CAP-006, PDR-016, PDR-039, PDR-040, PDR-042, ADR-038]
title: Design for collaborative PR handoffs
---

# Design — CAP-006 Collaborative PR handoffs

Routing precedes coordination. The agent states a simple or full recommendation from accepted-contract impact, identifies its escalation conditions, and repeats the assessment against the complete diff. Paths select checks but do not select meaning. A user-directed simple override of a full recommendation is retained in Git trailers; it creates no corpus artifact. Route selection never supplies push authority, and repository-local integration policy remains controlling.

For a chosen full loop, the pull request is the coordination unit. Separate changes retain separate branches and proceed in parallel; agents sharing one change coordinate through the hosted head SHA and review conversations.

- **Publication cadence:** a push is durability, never a signal ([PDR-040](../../decisions/PDR-040-push-is-durability-ready-is-explicit.md)). Every working turn that changed anything ends with the change branch pushed; the pull request opens as a draft at first publication and holds unfinished work visibly. Marking it ready is the readiness claim, taken only with verification and a clean review pass bound to the exact hosted head, and a substantive edit returns it to draft.

- **Reviewer:** reads one hosted head, returns a clean result tied to that SHA or publishes actionable findings as unresolved, resolvable review conversations. A reviewer without the necessary host capability exposes the enforcement gap and cannot report merge readiness.
- **Updater:** any agent that mutates the change for that turn. It observes starting head `S`, fetches it, produces reviewed commit `R`, verifies immediately before publication that the remote still permits a normal fast-forward update, pushes `R`, and confirms the pull request head is `R`.
- **Race handling:** if the remote moved from `S`, the updater fetches and reconciles without force, then reruns checks and review because the candidate changed. If accepted `main` advances after the PR is published, the branch merges current `main` and publishes normally rather than rebasing and force-pushing rewritten history. A merged or closed pull request is terminal; unpublished work waits for explicit human scope.
- **Finding lifecycle:** a finding is unresolved while work is missing or local. Resolution follows the hosted fix and clean review, never merely an edit in a worktree.
- **Dependent work:** an exception to the main-root rule remains a committed answered blocking question until digest. It names the unmerged base, human authorization, and meaning the dependent merge would bind; the acceptance brief repeats it for the merge decision. Git history preserves the record after the exception ends, without treating a forge branch as corpus truth.
- **Enforcement boundary:** hosting can require resolved conversations and reject non-fast-forward pushes. It cannot discover a finding or edit that was never published, so the agent must report that state as not merge-ready.
- **CI handoff:** the generated caller keeps the stable `validate` job while the upstream reusable workflow owns the validation steps. Updating its immutable reference imports scope, warning, acceptance-brief, and digest-gate repairs without changing the caller's runner or approved binary-acquisition choices.
