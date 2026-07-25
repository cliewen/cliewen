---
cliewen-skill: true
version: 0.6.0
---

<!-- Generated from Cliewen's canonical skill sources; edit those sources, not this file. -->

# clue-extract

Brownfield adoption: transform an existing repository's specification corpus into a Cliewen `/docs` corpus. Use once per adopted repository; the extraction is that repository's first `clue-delta` loop. Apply the **Decision records**, **Repository-local conventions**, and **Review boundary** below.

## Target contract

The extraction PR is complete only when all of these hold:

1. **The full taxonomy exists:** `/docs` has goals, plans, capabilities, decisions, constraints, analysis, and architecture, with an indexed README in every folder and README, criteria, and design files in each capability folder. Extract meaning; do not invent it. A folder with nothing real to hold stays empty but indexed.
2. **Everything extracted is born inferred:** Use `provenance: inferred`; decisions instead use `status: inferred` and `author: agent`. Human review may promote records to `verified`, file by file or in bulk. Extracted decisions also follow **Decision records** below.
3. **Existing criterion IDs survive; a criterion with none is minted deterministically:** Declare each capability's namespace with `ac-prefix:` and keep source IDs verbatim — never renumber, since IDs are meaning-immutable and existing test tags must keep resolving. When a source requirement carries no stable ID of its own, mint one in its capability's `ac-prefix:` namespace by numbering that file's requirements in the source's own stated or file order, sequentially from `-001` up; the same source state always mints the same IDs, so minting is reproducible rather than ad hoc per extraction. Record the minted range in the extraction report below.
4. **Every test keeps or gains exactly one purpose:** Existing criterion tags remain. Untagged tests get `Unit`, `Sanity`, or `Arch` according to intent. Where a JVM test framework is present, install an ArchUnit or equivalent rule enforcing one purpose tag per test; `clue` only harvests at file level.
5. **Born-`draft` is the sanctioned way to phase an oversized corpus:** A capability whose criteria are extracted but not yet given positive and negative tests stays `status: draft`, exactly as an untestable capability does. This is the legitimate way to extract a corpus too large to tag-test in a single change: extract every capability's meaning now, and let later changes activate each criterion together with its tests. Extraction is never partial by silent omission — a capability left `draft` for size, not for untestability, says so in the extraction report and gets a named plan door (see below) for the change that activates it.
6. **`clue validate` is green before the ready PR opens:** The extracted corpus is judged by the same rules as a greenfield corpus.
7. **The source corpus dies in the same PR:** Delete parallel specification trees, registries, and source-format skills; Git history is their archive. Two systems of record is zero systems of record.
8. **Routing is rewritten and reconciled:** Point every assistant entry point the repository carries — `AGENTS.md` and any other assistant-specific entry file, such as `CLAUDE.md` or `.cursor/rules` — to `/docs/README.md` and the installed `clue-*` skills; `AGENTS.md` is the flagship instance of this class, not its only member. Absorb compatible pre-existing instructions as repository-local conventions; record conflicts as open questions.
9. **An extraction report lands in `/docs/analysis`:** Record what was found, what mapped where, and what was dropped and why.
10. **Unsolved adoption items become named plan doors:** Never leave a silent gap.
11. **Every converted file carries exactly one frontmatter block:** When a source file already has frontmatter, the conversion replaces it, folding retained fields into the Cliewen block — even when an invisible prefix such as a UTF-8 byte-order mark hides the opening fence. Strip BOMs; after conversion, sweep the tree for a complete frontmatter block at the start of an artifact body: leftover source frontmatter is an extraction failure, and `clue validate` rejects both shapes.
12. **The committed extraction receives a clean agentic review:** Run `clue-verify`, including its automatic review loop, before opening the ready extraction PR.

## Source mappings

Per-source mappings live in this skill's `mappings/` folder; the current mapping is [openspec.md](mappings/openspec.md). The target contract above governs every extraction; a mapping only describes one source format. A new source format adds a mapping file, never another skill. If no mapping exists, writing one is the extraction PR's first task.

## Boundaries

Never invent unstated requirements, renumber or rename IDs, leave the source corpus alive for reference, promote your own output to `verified`, or change test code beyond adding missing purpose tags and the purpose-enforcement rule.

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
