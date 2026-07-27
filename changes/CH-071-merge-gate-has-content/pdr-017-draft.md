---
id: PDR-017
type: decision
status: inferred
links: [G-001, AN-008, P-007, PDR-012, PDR-013, C-012, C-013]
title: The merge gate carries an acceptance brief, and the review loop adds an advisory scenario-resolution step
author: agent
accepted-by: []
---

# PDR-017 — The merge gate has content

## Context and problem statement

[AN-008](../analysis/AN-008-methodology-critiques.md) pattern A: the corpus carries thousands of words of agent-facing skill instruction and nothing telling the human what to check at merge. [C-012](../constraints/C-012-agents-never-merge-own-changes.md) makes the merge a human act but never states what that act verifies beyond what the validator and the review loop already covered — the human boundary degrades to rubber-stamping. The book Cliewen is born from already answers this with an "after" gate whose advisory checks include whether a test matches its scenario, and separately supports an opt-in pause before implementation for spec review. What should the human verify at merge, how does an agent surface it, and how does an advisory semantic check coexist with the review loop's existing every-finding-resolved rule ([PDR-012](PDR-012-agentic-review-before-publication.md), [PDR-016](PDR-016-pr-state-carries-agent-handoffs.md))?

## Decision outcome

**The digest emits an acceptance brief at the top of the PR body, the review loop adds a non-blocking per-criterion scenario-resolution step that feeds it, and Propose gains an opt-in pre-implementation pause.**

- **The acceptance brief is the human's verification surface.** Emitted by the digest step, placed first in the PR body (before Summary, in reading order). It states: which plan item this change serves and whether that is still wanted — the human approves the criteria against reality, not just the code against the criteria; the added or changed acceptance criteria, listed verbatim with their scenarios; and what becomes binding on merge — the `inferred` decisions this change mints, and any records it invalidates or supersedes. The template carries the book's competence-heuristic warning (fluent agent output is not evidence of correctness) and a one-screen length cap stated as prose pressure toward small deltas — a cap the template asserts, not a check the CI enforces.
- **A per-criterion scenario-resolution step is advisory, not a gate.** For every added or changed criterion, the review loop (the same context-isolated or in-context-fallback reviewer [PDR-012](PDR-012-agentic-review-before-publication.md) already runs) compares the criterion's scenario against its referenced tests' setup, action, and assertions, and records one of: verifies the scenario, verifies something adjacent, or undetermined. This relocates the book's advisory test-matches-scenario check from unaided human work to agent work the human confirms. The verdict is informational content the digest carries into the brief; it is never an actionable finding and never blocks `clue validate` — the judge stays form-only, per [PDR-013](PDR-013-explicit-core-red-line.md)'s core boundary.
- **The advisory verdict does not enter the every-finding-resolved lifecycle.** [PDR-012](PDR-012-agentic-review-before-publication.md) requires actionable findings — correctness, intent mismatch, regression, security, missing evidence, unjustified complexity — to become unresolved hosted review conversations until repaired. A scenario-resolution verdict is not a defect claim; it is the same category of information the brief already carries (criteria listed verbatim for human judgment), not a machine verdict of wrong or right. If the human reads an "adjacent" or "undetermined" verdict as revealing a real problem, they raise it as an ordinary review comment, which then follows PDR-012's actionable-finding lifecycle like any other finding. The brief surfaces the signal; it does not manufacture a second class of enforced finding alongside actionable ones.
- **The pre-implementation pause is opt-in.** After Propose, the human may ask for a review of proposal.md before implementation starts; the agent stops and waits. The default is unchanged: a full change proceeds straight through Propose → Implement → Digest → Verify without pausing unless the human asks.
- **Private agent memory is never where work lives.** A new shared skill fragment states this once: anything needed to implement, continue, or hand off work lives in a corpus artifact, the change workspace, or the pull request — never only in an agent's conversation, because private memory does not survive a change of agent. AGENTS.md's existing system-of-record line is kept as the sole repo-local pointer to this rule; the fragment is the one place the rule itself is carried, not duplicated into AGENTS.md.

**Carrier:** this PDR; the `clue-delta` and `clue-verify` source templates and the new shared fragment (agent); `.github/pull_request_template.md` and its scaffolded copy (default); a deterministic unfilled-brief check in this repository's CI and the scaffolded wall (machine) — the check confirms the brief section was edited away from its template placeholder, not that its content is correct, because correctness is exactly what the human act of merging verifies.

### Rejected: make the scenario-resolution verdict a blocking check

Turning "verifies something adjacent" into a failing `clue validate` result would make a semantic judgment — whether a test really matches its Gherkin scenario — a mechanical gate, which [PDR-013](PDR-013-explicit-core-red-line.md) reserves for corpus form. Semantic drift between a scenario and its test is exactly the kind of meaning question the human merge act exists to catch; automating it away would re-empty the boundary this decision closes.

### Rejected: require the scenario-resolution verdict to become a hosted review conversation like an actionable finding

Every criterion would mint a conversation on every change, most reporting "verifies" and closing immediately — noise that trains reviewers to resolve without reading, the opposite of the goal. The brief already carries the verdicts where the human reads them at the point of decision; escalating a clean verdict to conversation-per-criterion overhead buys no additional safety PDR-012's existing actionable-finding path doesn't already provide when something is actually wrong.

### Rejected: make the pre-implementation pause the default

A mandatory pause before every full change reintroduces exactly the friction [PDR-012](PDR-012-agentic-review-before-publication.md) rejected for review: making a useful step optional-in-practice by making it required-by-default just moves the friction earlier. The book's spec-first pause is valuable when a human wants to see the shape of the change before code exists, not on every change regardless of size or risk.
