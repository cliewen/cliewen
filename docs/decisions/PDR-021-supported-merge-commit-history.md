---
id: PDR-021
type: decision
status: verified
links: [G-001, P-009, M-039, AN-014, PDR-007, PDR-016, ARCH-001]
title: Full Cliewen changes are accepted with a merge commit that preserves their branch history
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-021 — Full Cliewen changes are accepted with a merge commit that preserves their branch history

## Context and problem statement

Cliewen creates a proposal workspace before implementation and deletes that workspace in the digest. The human-controlled pull-request merge is acceptance, while repository history is the durable provenance archive. A pull request can be integrated by creating a merge commit, squashing its commits, or rebasing and fast-forwarding them; those modes can leave the same final tree while retaining different evidence about the accepted change. Which integration shape preserves the proposal, implementation, digest, and corpus history without making the forge the system of record?

## Decision outcome

**A full Cliewen change is accepted only through a human-controlled merge commit.** The repository's hosting configuration must allow merge commits and disable squash and rebase-and-merge for the protected default branch. The required `validate` check, pull-request requirement, resolved-conversation requirement, no-bypass rule, deletion restriction, and force-push block remain part of the same branch-protection setup.

The supported merge-commit path keeps the exact proposal, implementation, and digest commits as ancestors of `main`. The proposal and implementation may be transient in the final tree, but their committed contents remain inspectable from repository history; the digest commit also keeps the durable corpus changes in that same reachable chain. `git log` and ancestor queries therefore provide the archive without consulting pull-request pages or review metadata.

The three outcomes reproduced by the focused history fixture have different retention:

| Integration mode | What reaches `main` | Cliewen status |
|---|---|---|
| Merge commit | The original proposal, implementation, digest, and durable corpus commits remain reachable as the merged branch's ancestors | Supported |
| Squash and merge | One new commit carries the net tree, while the proposal, implementation, and digest commits are not reachable from `main` | Unsupported |
| Rebase and merge | Replayed commits may carry similar content, but the reviewed branch commits and their identities are rewritten before the fast-forward | Unsupported |

An agent may rebase an unpublished local branch onto newly accepted `main` before its first publication, as the review boundary already permits. That preparation does not change the required integration mode: once a full-change pull request is ready, the human accepts it with a merge commit. After publication, incorporating newer `main` continues to use the normal non-rewriting merge required by [PDR-016](PDR-016-pr-state-carries-agent-handoffs.md).

The support boundary must be established before adoption. The setup and operations path verifies both surfaces that decide the merge shape — the protected branch's allowed merge methods and the repository-wide pull-request settings — and the disposable branch-protection probe verifies that an undigested change fails the `--forbid-changes` gate and is blocked by the required check. A forge that cannot expose this merge-commit-only configuration is outside the supported full-change adoption boundary rather than an equivalent implementation.

This decision refines [PDR-007](PDR-007-review-boundary.md)'s human merge boundary and [ARCH-001](../architecture/architecture.md)'s repository-native provenance statement. It does not authorize agents to merge, make forge state durable meaning, or change the one-change-in-flight rule.

Plain changes under [PDR-011](PDR-011-plain-changes-bypass-cliewen.md) carry no proposal or digest chain worth preserving, so nothing here is aimed at them. They are nonetheless affected in practice: both enforcement surfaces — a ruleset's allowed merge methods and the repository-wide pull-request settings — are branch- or repository-scoped rather than per-pull-request, so a default branch restricted to merge commits restricts every pull request into it. This decision accepts that cost rather than claiming an exemption the enforcement cannot express. A repository that needs squash for plain changes must route them to a branch the rule does not target, and no full Cliewen change may be accepted through such a path.

**Carrier:** [C-012](../constraints/C-012-agents-never-merge-own-changes.md), [ARCH-001](../architecture/architecture.md), [ARCH-003](../architecture/core.md), the review-boundary source rendered into the lifecycle skills and scaffold, contributor routing, the public change-loop and operations/setup guidance, and the branch-protection probe all carry this support boundary. The repository and scaffolded routing hubs direct light and full changes to that source rather than restating it. The focused history fixture and content guards keep the three outcomes from being described as equivalent.

### Rejected: allow squash and rely on the pull request as the proposal archive

Squashing preserves the final tree but removes the branch commits that carried the proposal and implementation. Pull-request pages and review state are forge records, not the repository's system of record, and can disappear or become inaccessible independently of `main`.

### Rejected: allow rebase-and-merge because the logical commits remain visible

Rebase-and-merge rewrites the reviewed branch commits and their identities before the fast-forward. The resulting tree and commit messages can look complete while the exact accepted proposal, implementation, and digest chain is no longer the chain that was reviewed. The supported contract chooses the conservative, identity-preserving outcome.
