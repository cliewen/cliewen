---
id: CAP-006-design
type: design
status: active
links: [CAP-006, PDR-016]
title: Design for collaborative PR handoffs
---

# Design — CAP-006 Collaborative PR handoffs

The pull request is the coordination unit. Separate changes retain separate branches and proceed in parallel; agents sharing one change coordinate through the hosted head SHA and review conversations.

- **Reviewer:** reads one hosted head, returns a clean result tied to that SHA or publishes actionable findings as unresolved, resolvable review conversations. A reviewer without the necessary host capability exposes the enforcement gap and cannot report merge readiness.
- **Updater:** any agent that mutates the change for that turn. It observes starting head `S`, fetches it, produces reviewed commit `R`, verifies immediately before publication that the remote still permits a normal fast-forward update, pushes `R`, and confirms the pull request head is `R`.
- **Race handling:** if the remote moved from `S`, the updater fetches and reconciles without force, then reruns checks and review because the candidate changed. If accepted `main` advances after the PR is published, the branch merges current `main` and publishes normally rather than rebasing and force-pushing rewritten history. A merged or closed pull request is terminal; unpublished work waits for explicit human scope.
- **Finding lifecycle:** a finding is unresolved while work is missing or local. Resolution follows the hosted fix and clean review, never merely an edit in a worktree.
- **Enforcement boundary:** hosting can require resolved conversations and reject non-fast-forward pushes. It cannot discover a finding or edit that was never published, so the agent must report that state as not merge-ready.
