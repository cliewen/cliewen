---
version: 0.3.0
---

<!-- Generated from Cliewen's canonical skill sources; edit those sources, not this file. -->

# clue-verify

Run this pre-merge checklist before opening or updating any PR. When the `clue` CLI exists, `clue validate` performs the mechanical half; until then, check by hand. Never fix a failure by weakening the check.

- [ ] The change uses the correct workspace under **Change tiers** below.
- [ ] Every artifact touched has frontmatter `id`, `type`, `status`, `links`, and `title`, plus decision `author`/`accepted-by`, constraint `source`/`enforcement`, capability `goal`, and any other type-specific fields.
- [ ] Every `links` entry resolves to an existing ID.
- [ ] The proposal names a real plan item or explicitly declares the change plan-less.
- [ ] Plan bookkeeping reflects the merge, and no completed plan changed.
- [ ] Every active acceptance criterion has positive and negative tests, or its capability honestly stays `draft` with the gap stated.
- [ ] Every `/docs/**` folder has a README; index blocks list every sibling artifact and no deleted file.
- [ ] The change was assessed against every constraint and quality scenario.
- [ ] Repository-local conventions satisfy the contract below.
- [ ] Diagrams are inline Mermaid and readable when rendered.
- [ ] The full-change workspace is absent after digest; `main` never contains `/changes/`.
- [ ] Every decision satisfies **Decision records** below, including routing, timeless content, provenance, objections, and pending approval signatures.
- [ ] `git merge-base HEAD origin/main` equals `origin/main` after fetching; no other change workspace is visible on this branch.
- [ ] The branch and hosted PR satisfy the **Review boundary** below.

## Change tiers

A change is light only when all of these hold: no decision is made, no acceptance criterion or capability meaning changes, no semantic plan mutation occurs, and no methodology carrier such as a skill, AGENTS.md rule, or lint rule is touched. Typical light changes: typos, documentation clarity, dependency bumps, pure refactors, CI plumbing. A light change skips the transient workspace; its branch and ready PR remain mandatory, and the PR description is the proposal with a real plan item or an explicit plan-less declaration.

Every other change uses the full loop and a `/changes/CH-xxx-slug/` workspace. Escalate immediately if a decision, open question, meaning change, or methodology-carrier edit appears during work.

## Decision records

Route every decision by reversal cost. A cheap-and-local-to-reverse decision is a dated row in `docs/decisions/log.md` (columns `Date | Decision | Why | Change/PR`); otherwise write an ADR for software or corpus architecture, or a PDR for how the project works. A decision adopting a well-established practice cites it by name and records only the local why.

Agent-authored decisions start `status: inferred` and `author: agent`. Merging makes them binding without changing that status. Only explicit human approval promotes a decision to `verified`; record every approver in `accepted-by:`, use the first approval date, and cite the venue. An explicit objection keeps the decision `inferred` and becomes an open question.

Every decision record is timeless: state what is decided and only the enduring context and rationale needed to understand it. Keep triggering incidents, chronology, conversations, implementation details, and review history in findings, the change workspace, the PR, and Git history.

## Repository-local conventions

Apply the repository-local conventions declared in AGENTS.md, including digest requirements such as a user-facing changelog entry. Local conventions extend the methodology and never override it. If AGENTS.md conflicts with a skill, record the conflict in `open-questions.md` and stop for a human decision; never choose silently.

## Review boundary

Every change branches from the current tip of `main`, never from unaccepted work. Each author takes one change to its PR before starting another; independent authors may work in parallel from `main`. If work must build on an unmerged change, record a blocking open question and stop unless the human explicitly authorizes it. If another change merges first, rebase onto the new `main` tip and repeat verification.

Open the PR ready for review only after verification passes, never as a draft. The PR is the completed proposal's human review gate; unfinished work stays on the branch. An agent never merges its own PR, creates a local merge commit into `main`, or pushes to `main`. After opening the PR an agent stops and waits; it never starts the next change while the previous one is unreviewed. Review fixes stay on the same branch and PR; a follow-up change exists only when a human has accepted this one and explicitly scoped the follow-up.

After a human reports the merge, orient before starting anything else: describe the plan's next unfinished step in plain language and ask whether to start it, or say that the plan has nothing left and ask what comes next.
