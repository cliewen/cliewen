---
id: PDR-016
type: decision
status: inferred
links: [G-001, AN-009, PDR-007, PDR-012, C-012, CAP-006]
title: Hosted PR state carries review findings and updater handoffs across agents
author: agent
accepted-by: []
---

# PDR-016 — Hosted PR state carries review findings and updater handoffs across agents

## Context and problem statement

The review boundary requires an exact published head, but collaboration may move between agents before that final handoff. Private conversation state cannot tell the next agent or the human that a finding or repair remains unpublished. How can the boundary make known incomplete work durable without serializing independent authors?

## Decision outcome

**Collaboration is scoped to the pull request, and its hosted head and resolvable review conversations are the durable handoff state.**

- **Reviews are SHA-bound.** A reviewer names the hosted head it inspected. A clean result applies only to that commit; an edit invalidates it.
- **Actionable findings become unresolved hosted review conversations.** Where the forge supports resolvable conversations, each finding remains unresolved until a hosted commit contains the verified repair. Required conversation resolution turns known findings into a merge prerequisite. If the reviewer cannot publish a resolvable finding, it reports the pull request as not merge-ready and discloses that the host cannot enforce the finding.
- **The updater role follows mutation, not identity.** Any agent asked to edit an open pull request becomes the updater for that turn. It fetches the current hosted head, commits and verifies the complete repair, obtains a clean review of the repaired commit, pushes without force, confirms the hosted head equals that reviewed commit, and only then resolves satisfied findings.
- **Same-PR updates use optimistic concurrency.** The updater records the starting hosted head and rechecks it before publication. A changed head or rejected non-fast-forward push requires fetch, reconciliation, and renewed verification and review of the resulting commit. Remote work is never silently overwritten.
- **Independent work remains parallel.** One in-flight Cliewen change per initiating author remains unchanged. Reviewing or helping update an existing pull request does not mint another change or create a global lock; separate changes continue from accepted `main`.
- **Closed boundaries stop updates.** If the pull request merged or closed before publication, the updater stops, preserves and reports any local work as unpublished, and does not create a follow-up change without explicit human scope.

Premature human merge remains a human error. This contract makes known findings visible and enforceable where hosting permits; it does not claim to detect private intent that no agent published.

**Carrier:** the shared review-boundary source and the `clue-verify` review loop (agent); the scaffolded and repository routing hubs and pull-request template (default/local); the protected-host guide's required-conversation-resolution setting (human/machine where supported); CAP-006's tests hold the generated contract and guidance.

### Rejected: serialize all Cliewen changes

Independent branches do not share a remote head. A global lock would punish ordinary team parallelism without protecting the contested branch more effectively than its own SHA.

### Rejected: permanent ownership by the first implementer

Humans legitimately ask reviewers or later agents to repair the same pull request. Mutation creates the publication duty for that turn; an agent name does not.

### Rejected: keep findings only in the agent conversation

Private conversation state is unavailable to other agents, hosted protection, and often the human at merge. It cannot carry a shared boundary.

### Rejected: make CI detect unpublished local work

CI receives hosted commits and events. It cannot inspect another machine's worktree or infer an unreported intention.
