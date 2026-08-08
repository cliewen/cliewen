---
id: PDR-012
type: decision
status: verified
links: [G-001, ARCH-002, PDR-007, C-012]
title: Every Cliewen change receives an automatic agentic review before publication
author: agent
accepted-by: Flemming N. Larsen (2026-07-21, implementation conversation; finding-grounding amendment approved 2026-07-22, follow-up conversation)
---

# PDR-012 — Agentic review before publication

> **Amended by [PDR-035](PDR-035-bounded-agentic-review-loop.md) and [PDR-036](PDR-036-review-loop-budget-and-human-checkpoint.md):** the loop exits on no blocking findings under its own severity model, computed-figure findings are advisory, the loop runs to a bounded maximum stated by [C-017](../constraints/C-017-agentic-review-loop-is-bounded.md) with a later pass allowed only after a blocking pass, reaching that maximum with blocking findings outstanding reports them to the human and asks whether to continue rather than earning further passes, and the handoff reports the pass count.

## Context and problem statement

An implementing coding agent can complete its checks and still find actionable defects when asked to review the same work from a cleared context. Implementation context anchors the agent to its chosen approach, while a fresh review context changes the task from constructing the change to challenging a fixed candidate. Requiring the human to clear the context and request that review manually makes a useful quality step optional and repetitive. How should Cliewen make that challenge automatic without pretending that the same model is an independent human reviewer or creating a review loop that never converges?

## Decision outcome

**Every Cliewen change receives an automatic adversarial agent review of its committed candidate before publication, using a context-isolated reviewer when the coding-agent host supports it and a disclosed in-context fallback otherwise.**

- **Review is part of verification.** `clue-delta` hands the committed candidate to `clue-verify`; the human does not need to initiate a separate review request. A change is locally ready only after the current commit has received at least one clean review pass and applicable local checks pass.
- **Isolation is capability-dependent and explicit.** A host that can delegate without inherited implementation conversation starts a new read-only reviewer. The reviewer receives the repository, branch and base, the proposal recovered from branch history for a full change or the user's request and accepted clarifications for a light change, and no implementation rationale. A host without isolated delegation performs an explicit adversarial pass in the current context and reports that weaker mode in the handoff.
- **The reviewer challenges evidence and does not edit.** It inspects the complete base diff, durable corpus, tests, constraints, and quality scenarios for correctness, intent mismatch, regressions, security problems, missing evidence, and unjustified complexity. Under PDR-035 it returns blocking actionable findings and advisory non-actionable observations for the publication gate, each with severity, location, evidence, and a concrete remediation; preferences and speculative scope expansion are excluded.
- **A finding is grounded in an operative requirement.** The reviewer identifies the current requirement or declared intent that is violated and explains the concrete consequence. It applies authoritative decisions and explicit lifecycle rules before escalating nearby wording: under [PDR-007](PDR-007-review-boundary.md), human-controlled merge is mandatory but duplicate human code review is not; under [ADR-012](ADR-012-release-notes-from-changelog.md), a release cut renames `[Unreleased]` to the versioned section, so the versioned section is the required evidence. Historical descriptions, optional activity, alternative readings, and lifecycle-correct state are not actionable defects by themselves.
- **Fixes invalidate the review result.** The implementing context resolves each blocking finding, commits the repaired candidate, reruns applicable local checks against that commit, and then starts a new review pass on the same commit. Every blocking repair invalidates the previous clean pass. Under PDR-035, the loop ends when a pass on the current commit has no blocking findings, or when a blocking question is recorded and the change stops. An advisory repair may ride before a pass already required by a blocking repair; an advisory first reported by a clean pass stays in the verification handoff for a later change, not in an unresolved repair-required conversation, so the published candidate remains the exact reviewed commit without making the advisory a merge gate.
- **The existing PR boundary remains.** The final handoff reports the review mode, reviewed commit, pass count, and advisory findings left open. Hosted CI still checks the exact published head, a human still controls merge, and PDR-007's branch, PR, and no-agent-merge rules are unchanged.

**Carrier:** the automatic review loop in the generated `clue-verify` skill, invoked by `clue-delta`; the generated review-boundary fragment; and the review-loop assertion in the skill generator tests. The routing hubs direct light and full changes into `clue-delta` without restating the handoff. The CLI does not attempt to prove agent context isolation.

### Rejected: require the human to start a clean-context review

This reproduces the motivating friction and makes the quality gate depend on a prompt the methodology can issue itself.

### Rejected: let the reviewer fix its own findings

Combining review and repair in one delegated context restores implementation anchoring and obscures whether the fixed commit received a clean challenge. The reviewer stays read-only and the implementing context owns repairs.

### Rejected: require a fixed number of review passes

Three ritual passes may all repeat the same blind spot, while one clean pass after a small change may be sufficient. Convergence is tied to the reviewed commit and actionable findings: substantive repair requires a fresh pass; stylistic preference cannot keep the loop alive.

PDR-035 preserves this rejection: its three-pass budget is a maximum for an ordinary loop, never a required number of passes, and a blocking third pass earns the repair's necessary follow-up review.
