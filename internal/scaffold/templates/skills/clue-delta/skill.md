---
cliewen-skill: true
version: 0.6.0
---

<!-- Generated from Cliewen's canonical skill sources; edit those sources, not this file. -->

# clue-delta

Use for every Cliewen change: features, fixes, docs, and plans whose meaning belongs in the corpus or methodology. Plain changes are classified by AGENTS.md before the corpus is loaded and do not invoke this skill. Apply the **Change scope and tiers**, **Decision records**, **Repository-local conventions**, and **Review boundary** below throughout the loop.

1. **Branch:** Follow the review boundary and name the branch `ch-xxx-slug`. Take the next free CH number by searching Git history and `/changes/` for the highest used number.
2. **Propose:** For a full change, create `/changes/CH-xxx-slug/` and commit it before implementation:
   - `proposal.md` states what and why; its frontmatter `links` names the real plan item it serves or explicitly declares the change plan-less.
   - `tasks.md` is an ordered `- [ ]` checklist with dependencies first and at most one nested level. Mark `[x]` the moment a task completes, never in batch at the end. Mark an infeasible task `[-]` with its reason on the same line. A behavior-changing task names the acceptance-criterion IDs it serves; if none exists, add the criterion before implementation. Tests trace to criteria, never transient tasks.
   - `open-questions.md` records blocking questions. When one appears, write it and stop; the human answer becomes a decision record.
3. **Implement:** Update the permanent corpus. Capabilities own README, criteria, and design files. Write criteria as Gherkin tagged `@AC-xxx`; every active criterion gets positive and negative tests in the same change, while an untestable capability stays `draft`. Split a criterion that cannot be verified by a focused pair. Every test declares exactly one purpose: the criterion ID, `Unit`, `Sanity`, or `Arch`, using framework tags where available and the test-name prefix in Go. When a criterion's meaning changes, retire it with `@retired`, keep the tombstone, mint a new ID, and remove or retag its tests.
4. **Digest:** After every task is `[x]` or `[-]` with a reason, update permanent `/docs`, regenerate README indexes, apply repository-local digest conventions, record decisions, and update plan bookkeeping. Delete the change workspace. The digest is never a task in `tasks.md`; deletion is the digest, so a self-referential digest task cannot be completed honestly.
5. **Verify, review, and propose for acceptance:** Run `clue-verify`, including its automatic agentic review loop on the verified committed candidate, then open the PR under the review boundary. Never ask the human to initiate the review. Merging accepts the change; decision provenance follows **Decision records** below.

Keep deltas small: Git merges text, not meaning.

## Change scope and tiers

Classify scope before using the Cliewen loop. Three rules decide the tier, by how deeply the change reaches into meaning; take the first rule that matches.

1. **Plain — nothing about meaning changes.** The change has no effect on product behavior, intent, executable evidence, decisions, plans, policy, or methodology. Protected surfaces are never plain: `/docs`, `/changes`, product code, tests, configuration, build and release machinery, security and governance policy, AGENTS.md rules, skills, and lint rules; changes to commands, contracts, user workflow, or normative instructions are not editorial. A plain change stays outside this skill: an ordinary branch from the current tip of `main`, checks relevant to the changed surface, a ready PR, and human merge. It takes no CH identity, corpus read, plan item, proposal artifacts, Cliewen verification, plan bookkeeping, or Cliewen-mandated changelog entry, does not consume the one-Cliewen-change-in-flight slot, and never builds on unmerged work.
2. **Light — meaning is touched but not changed.** The change makes no decision, changes no acceptance criterion or capability meaning, mutates no plan semantically, and touches no methodology carrier such as a skill, AGENTS.md rule, or lint rule. Typical: typos or documentation clarity on a protected surface, dependency bumps, pure refactors, and CI plumbing. Its branch and ready PR are mandatory and it skips the transient workspace, so the PR description is the proposal, naming a real plan item or explicitly declaring the change plan-less.
3. **Full — everything else.** The whole loop, with a `/changes/CH-xxx-slug/` workspace.

Two guards hold above the rules. **Uncertainty escalates:** when the tier is unclear, take the higher one. **Discovery escalates immediately:** the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move to the full loop before continuing.

## Decision records

Route every decision by reversal cost. A cheap-and-local-to-reverse decision is a dated row in `docs/decisions/log.md` (columns `Date | Decision | Why | Change/PR`); otherwise write an ADR for software or corpus architecture, or a PDR for how the project works. A decision adopting a well-established practice cites it by name and records only the local why.

Agent-authored decisions start `status: inferred` and `author: agent`. Merging makes them binding without changing that status. Only explicit human approval promotes a decision to `verified`; record every approver in `accepted-by:`, use the first approval date, and cite the venue. An explicit objection keeps the decision `inferred` and becomes an open question.

Every decision record is timeless: state what is decided and only the enduring context and rationale needed to understand it. Keep triggering incidents, chronology, conversations, implementation details, and review history in findings, the change workspace, the PR, and Git history.

## Repository-local conventions

For a Cliewen change, apply the repository-local conventions declared in AGENTS.md, including digest requirements such as a user-facing changelog entry. Plain changes follow only the repository conventions that apply to their changed surface. Local conventions extend the methodology and never override it. If AGENTS.md conflicts with a skill, record the conflict in `open-questions.md` and stop for a human decision; never choose silently.

## Review boundary

Every change branches from the current tip of `main`, never from unaccepted work. Each author takes one Cliewen change to its PR before starting another; independent authors may work in parallel from `main`, and plain changes do not consume this slot. If work must build on an unmerged change, record a blocking open question and stop unless the human explicitly authorizes it. If another change merges first, rebase onto the new `main` tip and repeat verification.

Open the PR ready for review only after local verification and the automatic agentic review loop pass on the current commit, never as a draft. The PR is the completed proposal's authorization and protected-integration boundary, not a demand for duplicate human code review: the agent may publish the candidate, but only a human-controlled PR merge accepts it. Unfinished work stays on the branch. An agent never merges its own PR, creates a local merge commit into `main`, or pushes to `main`.

A PR alone displays hosted CI but does not enforce it. Where hosting supports enforcement, the PR triggers CI, branch protection makes its required status check a merge precondition, and the agent cannot silently skip the gate. Never weaken the workflow or required-check policy to make a change pass.

Ready means the hosted PR contains the exact locally reviewed and verified state. Before reporting any change ready, commit every intended edit, run the applicable local verification and a clean agentic review pass against that commit, require `git status --porcelain` to be empty, push the reviewed commit, and confirm that the ready hosted PR's head branch and SHA equal the current local branch and `HEAD`. Perform the hosted check immediately after opening or updating the PR; if either side differs, reconcile without silently overwriting remote work, rerun verification and review on the resulting commit, require a clean worktree, push the reviewed commit, and check again. A human-requested local stopping point such as "commit only" is preserved work, not a completed or mergeable change, and the agent says that no ready PR exists.

After opening and confirming the PR, an agent stops and waits on further Cliewen work; independent plain changes may still proceed from accepted `main`. Review fixes stay on the same branch and PR: commit them, rerun local verification and the automatic agentic review loop against that commit, require a clean worktree, push the reviewed commit there, and repeat the hosted-head check before reporting ready again. A follow-up Cliewen change exists only when a human has accepted this one and explicitly scoped the follow-up.

After a human reports the merge, orient before starting anything else: describe the plan's next unfinished step in plain language and ask whether to start it, or say that the plan has nothing left and ask what comes next.
