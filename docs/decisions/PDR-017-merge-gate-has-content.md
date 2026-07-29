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

[AN-008](../analysis/AN-008-methodology-critiques.md) pattern A identifies that the corpus carries agent-facing skill instruction and nothing telling the human what to check at merge. The human boundary degrades to rubber-stamping unless it names what remains beyond the validator and review loop. The book's after-gate uses an advisory test-matches-scenario check and supports an opt-in pre-implementation spec review. How can the merge boundary make that human responsibility visible without turning it into duplicate code review or claiming that semantic alignment is mechanically proven?

## Decision outcome

**The digest emits an acceptance brief at the top of the PR body, the review loop adds a non-blocking per-criterion scenario-resolution step that feeds it, and Propose gains an opt-in pre-implementation pause.**

- The acceptance brief is the human's verification surface. It asks whether the plan item remains wanted, lists every added or changed criterion verbatim with its scenarios, and names the inferred decisions that merge binds plus records invalidated or superseded. It carries a competence-heuristic warning and a one-screen cap as prose pressure toward small deltas, not a CI-enforced truncation rule.
- The review loop compares every added or changed criterion with executable evidence against its referenced tests' setup, action, and assertions and records one advisory verdict: `verifies`, `verifies-something-adjacent`, or `undetermined`. For a genuine `Test-type: Human` criterion, the acceptance brief's required criterion-and-scenario entry is its proof; the review confirms that entry rather than comparing it with a nonexistent test. The verdict is informational brief content, not an actionable finding or a `clue validate` gate. A real problem it reveals follows the ordinary actionable-finding lifecycle.
- A full change may opt into a spec-first pause after Propose; it records the pause in tasks and waits for human direction. The ordinary ready-PR loop remains the default.
- Work needed to implement, continue, review, or hand off a change belongs in a corpus artifact, change workspace, or pull request. An agent's private memory is never an authoritative carrier.

**Carrier:** the `clue-delta` and `clue-verify` source templates and their shared durable-work fragment (agent); the pull-request template (human/default); CI and the scaffolded CI wall (machine/default); CAP-006 criteria and tests (evidence); and the public guide (human explanation).

### Rejected: make the scenario-resolution verdict a blocking check

Turning an advisory semantic judgment into a failing validator result would claim that the deterministic judge proves meaning. Semantic drift is a human merge question; automating it away would re-empty the boundary this decision closes.

### Rejected: require every scenario verdict to become a hosted review conversation

Every criterion would mint a conversation on every change, most reporting `verifies` and closing immediately. That noise trains reviewers to resolve without reading; the brief carries the verdict at the point the human decides, while actual defects already use the existing conversation path.

### Rejected: make the pre-implementation pause the default

A mandatory pause before every full change adds a second gate even where the completed-candidate PR is sufficient. The pause is valuable when a human wants design review before code exists, not by default regardless of size or risk.
