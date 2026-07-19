---
cliewen-skill: true
version: 0.3.0
---

<!-- Generated from Cliewen's canonical skill sources; edit those sources, not this file. -->

# clue-extract

Brownfield adoption: transform an existing repository's specification corpus into a Cliewen `/docs` corpus. Use once per adopted repository; the extraction is that repository's first `clue-delta` loop. Apply the **Decision records**, **Repository-local conventions**, and **Review boundary** below.

## Target contract

The extraction PR is complete only when all of these hold:

1. **The full taxonomy exists:** `/docs` has goals, plans, capabilities, decisions, constraints, quality, analysis, and architecture, with an indexed README in every folder and README, criteria, and design files in each capability folder. Extract meaning; do not invent it. A folder with nothing real to hold stays empty but indexed.
2. **Everything extracted is born inferred:** Use `provenance: inferred`; decisions instead use `status: inferred` and `author: agent`. Human review may promote records to `verified`, file by file or in bulk. Extracted decisions also follow **Decision records** below.
3. **Existing criterion IDs survive:** Declare each capability's namespace with `ac-prefix:` and keep source IDs verbatim. Never renumber; IDs are meaning-immutable and existing test tags must keep resolving.
4. **Every test keeps or gains exactly one purpose:** Existing criterion tags remain. Untagged tests get `Unit`, `Sanity`, or `Arch` according to intent. Where a JVM test framework is present, install an ArchUnit or equivalent rule enforcing one purpose tag per test; `clue` only harvests at file level.
5. **`clue validate` is green before the ready PR opens:** The extracted corpus is judged by the same rules as a greenfield corpus.
6. **The source corpus dies in the same PR:** Delete parallel specification trees, registries, and source-format skills; Git history is their archive. Two systems of record is zero systems of record.
7. **Routing is rewritten and reconciled:** Point AGENTS.md to `/docs/README.md` and the installed `clue-*` skills. Absorb compatible pre-existing instructions as repository-local conventions; record conflicts as open questions.
8. **An extraction report lands in `/docs/analysis`:** Record what was found, what mapped where, and what was dropped and why.
9. **Unsolved adoption items become named plan doors:** Never leave a silent gap.
10. **Every converted file carries exactly one frontmatter block:** When a source file already has frontmatter, the conversion replaces it, folding retained fields into the Cliewen block — even when an invisible prefix such as a UTF-8 byte-order mark hides the opening fence. Strip BOMs; after conversion, sweep the tree for a complete frontmatter block at the start of an artifact body: leftover source frontmatter is an extraction failure, and `clue validate` rejects both shapes.

## Source mappings

Per-source mappings live in this skill's `mappings/` folder; the current mapping is [openspec.md](mappings/openspec.md). The target contract above governs every extraction; a mapping only describes one source format. A new source format adds a mapping file, never another skill. If no mapping exists, writing one is the extraction PR's first task.

## Boundaries

Never invent unstated requirements, renumber or rename IDs, leave the source corpus alive for reference, promote your own output to `verified`, or change test code beyond adding missing purpose tags and the purpose-enforcement rule.

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
