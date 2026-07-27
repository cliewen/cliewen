---
id: AN-009
type: analysis
status: active
links: [G-001, AN-007, PDR-007, PDR-012, C-012]
title: Multi-agent review repairs can remain outside the pull request
---

# AN-009 — Multi-agent review repairs can remain outside the pull request

## Risk

When several agents review and repair one pull request, a finding or intended fix may remain in one agent's private session while another agent and the human see the hosted pull request as complete.

## Evidence boundary

The maintainer reported on 2026-07-27 that this omission has been caught at least five times across single-agent and multi-agent review flows. The conversations and uncommitted worktrees are not versioned evidence, so the count and exact local states cannot be independently reconstructed. The repository was inspected at `a4d17d129020b078f753f41f2b0a0a3ff738e028` on Windows amd64 with PowerShell and GitHub CLI. At that revision, [AN-007](AN-007-review-handoff-gap.md) and [C-012](../constraints/C-012-agents-never-merge-own-changes.md) already require a clean, pushed, exact-head handoff, while the generated review skill keeps its delegated reviewer read-only. GitHub protects this repository's `main` with the required `validate` check and required conversation resolution.

## Observations

- The exact-head rule proves that a completed handoff reached the pull request, but only if the agent reaches and follows that final step. It gives another agent no durable signal that a finding or fix is pending in a private session.
- Agent identity is not a stable ownership boundary. A reviewer may become an updater when the human asks it to fix a finding, and a later reviewer or the original implementer may continue the same pull request.
- Independent pull requests do not need serialization. Existing rules already allow any number of authors to branch from accepted `main`; only collaborators updating the same remote branch contend.
- Git rejects a normal non-fast-forward push when the remote head changes. Reading the hosted head before work and again before publication therefore supplies an optimistic concurrency boundary without a global lock or force push.
- A resolvable hosted review conversation is durable across agent sessions and can be a merge prerequisite where the forge supports required conversation resolution. A chat-only finding cannot be discovered by CI or another agent.
- CI cannot inspect an uncommitted worktree or infer that an agent intended another edit. A host can block a known unresolved finding, not an undisclosed local intention.

## Options assessed

A global one-change-at-a-time rule would prevent ordinary team parallelism while adding no protection to the handoff inside one pull request. Assigning permanent ownership to the first implementing agent fails when another agent is explicitly asked to repair the branch. Treating final-response wording as the shared record remains private to one conversation and repeats the observed gap. A forge-specific `clue` command would move provider credentials and network state into the deterministic judge without making undisclosed local intent observable.

The proportionate boundary is per pull request and per update: findings become unresolved hosted conversations tied to the reviewed head; any agent that edits becomes the updater for that turn; the updater publishes against the observed head without force, verifies the resulting hosted SHA, and resolves findings only after the fix is hosted and reviewed.

## Finding and consumer

The missing unit is a durable, SHA-bound collaboration handoff, not another global author lock. CH-070 consumes this finding through PDR-016, CAP-006, C-012, the generated lifecycle skills, and the protected-host guidance.
