---
version: 0.3.0
---

<!-- Generated from Cliewen's canonical skill sources; edit those sources, not this file. -->

# clue-delta

Use for every mutation of `main`: features, fixes, docs, and plans. Apply the **Change tiers**, **Decision records**, **Repository-local conventions**, and **Review boundary** below throughout the loop.

1. **Branch:** Follow the review boundary and name the branch `ch-xxx-slug`. Take the next free CH number by searching Git history and `/changes/` for the highest used number.
2. **Propose:** For a full change, create `/changes/CH-xxx-slug/` and commit it before implementation:
   - `proposal.md` states what and why; its frontmatter `links` names the real plan item it serves or explicitly declares the change plan-less.
   - `tasks.md` is an ordered `- [ ]` checklist with dependencies first and at most one nested level. Mark `[x]` the moment a task completes. Mark an infeasible task `[-]` with its reason on the same line. A behavior-changing task names the acceptance-criterion IDs it serves; if none exists, add the criterion before implementation. Tests trace to criteria, never transient tasks.
   - `open-questions.md` records blocking questions. When one appears, write it and stop; the human answer becomes a decision record.
3. **Implement:** Update the permanent corpus. Capabilities own README, criteria, and design files. Write criteria as Gherkin tagged `@AC-xxx`; every active criterion gets positive and negative tests in the same change, while an untestable capability stays `draft`. Split a criterion that cannot be verified by a focused pair. Every test declares exactly one purpose: the criterion ID, `Unit`, `Sanity`, or `Arch`, using framework tags where available and the test-name prefix in Go. When a criterion's meaning changes, retire it with `@retired`, keep the tombstone, mint a new ID, and remove or retag its tests.
4. **Digest:** After every task is `[x]` or `[-]` with a reason, update permanent `/docs`, regenerate README indexes, apply repository-local digest conventions, record decisions, and update plan bookkeeping. Delete the change workspace. The digest is never a task in `tasks.md`; deletion is the digest, so a self-referential digest task cannot be completed honestly.
5. **Verify and propose for acceptance:** Run `clue-verify`, then open the PR under the review boundary. Merging accepts the change; decision provenance follows **Decision records** below.

Keep deltas small: Git merges text, not meaning.

## Change tiers

A change is light only when all of these hold: no decision is made, no acceptance criterion or capability meaning changes, no semantic plan mutation occurs, and no methodology carrier such as a skill, AGENTS.md rule, or lint rule is touched. A light change skips the transient workspace; its branch and ready PR remain mandatory, and the PR description is the proposal with a real plan item or an explicit plan-less declaration.

Every other change uses the full loop and a `/changes/CH-xxx-slug/` workspace. Escalate immediately if a decision, open question, meaning change, or methodology-carrier edit appears during work.

## Decision records

Route every decision by reversal cost. A cheap-and-local-to-reverse decision is a dated row in `docs/decisions/log.md`; otherwise write an ADR for software or corpus architecture, or a PDR for how the project works. A decision adopting a well-established practice cites it by name and records only the local why.

Agent-authored decisions start `status: inferred` and `author: agent`. Merging makes them binding without changing that status. Only explicit human approval promotes a decision to `verified`; record every approver in `accepted-by:`, use the first approval date, and cite the venue. An explicit objection keeps the decision `inferred` and becomes an open question.

Every decision record is timeless: state what is decided and only the enduring context and rationale needed to understand it. Keep triggering incidents, chronology, conversations, implementation details, and review history in findings, the change workspace, the PR, and Git history.

## Repository-local conventions

Apply the repository-local conventions declared in AGENTS.md, including digest requirements such as a user-facing changelog entry. Local conventions extend the methodology and never override it. If AGENTS.md conflicts with a skill, record the conflict in `open-questions.md` and stop for a human decision; never choose silently.

## Review boundary

Every change branches from the current tip of `main`, never from unaccepted work. Each author takes one change to its PR before starting another; independent authors may work in parallel from `main`. If work must build on an unmerged change, record a blocking open question and stop unless the human explicitly authorizes it. If another change merges first, rebase onto the new `main` tip and repeat verification.

Open the PR ready for review only after verification passes, never as a draft. The PR is the completed proposal's human review gate; unfinished work stays on the branch. An agent never merges its own PR, creates a local merge commit into `main`, or pushes to `main`. Review fixes stay on the same branch and PR.

After a human reports the merge, orient before starting anything else: describe the plan's next unfinished step in plain language and ask whether to start it, or say that the plan has nothing left and ask what comes next.
