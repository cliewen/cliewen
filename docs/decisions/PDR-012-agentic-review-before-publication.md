---
id: PDR-012
type: decision
status: verified
links: [G-001, ARCH-002, PDR-007, C-012]
title: Every Cliewen change receives an automatic agentic review before publication
author: agent
accepted-by: Flemming N. Larsen (2026-07-21, implementation conversation)
---

# PDR-012 — Agentic review before publication

## Context and problem statement

An implementing coding agent can complete its checks and still find actionable defects when asked to review the same work from a cleared context. Implementation context anchors the agent to its chosen approach, while a fresh review context changes the task from constructing the change to challenging a fixed candidate. Requiring the human to clear the context and request that review manually makes a useful quality step optional and repetitive. How should Cliewen make that challenge automatic without pretending that the same model is an independent human reviewer or creating a review loop that never converges?

## Decision outcome

**Every Cliewen change receives an automatic adversarial agent review of its committed candidate before publication, using a context-isolated reviewer when the coding-agent host supports it and a disclosed in-context fallback otherwise.**

- **Review is part of verification.** `clue-delta` hands the committed candidate to `clue-verify`; the human does not need to initiate a separate review request. A change is locally ready only after the current commit has received at least one clean review pass and applicable local checks pass.
- **Isolation is capability-dependent and explicit.** A host that can delegate without inherited implementation conversation starts a new read-only reviewer. The reviewer receives the repository, branch and base, the proposal recovered from branch history for a full change or the user's request and accepted clarifications for a light change, and no implementation rationale. A host without isolated delegation performs an explicit adversarial pass in the current context and reports that weaker mode in the handoff.
- **The reviewer challenges evidence and does not edit.** It inspects the complete base diff, durable corpus, tests, constraints, and quality scenarios for correctness, intent mismatch, regressions, security problems, missing evidence, and unjustified complexity. It returns only actionable findings with severity, location, evidence, and a concrete remediation; preferences and speculative scope expansion are excluded.
- **Fixes invalidate the review result.** The implementing context resolves each finding, reruns applicable local checks, and commits the result before a new review pass. Every substantive edit invalidates the previous clean pass. The loop ends only when a pass on the current commit has no actionable findings, or when a blocking question is recorded and the change stops.
- **The existing PR boundary remains.** The final handoff reports the review mode and reviewed commit. Hosted CI still checks the exact published head, a human still controls merge, and PDR-007's branch, PR, and no-agent-merge rules are unchanged.

**Carrier:** the automatic review loop in the generated `clue-verify` skill, invoked by `clue-delta`; the ready-state rule in scaffolded and repository routing hubs; and the review-loop assertion in the skill generator tests. The CLI does not attempt to prove agent context isolation.

### Rejected: require the human to start a clean-context review

This reproduces the motivating friction and makes the quality gate depend on a prompt the methodology can issue itself.

### Rejected: let the reviewer fix its own findings

Combining review and repair in one delegated context restores implementation anchoring and obscures whether the fixed commit received a clean challenge. The reviewer stays read-only and the implementing context owns repairs.

### Rejected: require a fixed number of review passes

Three ritual passes may all repeat the same blind spot, while one clean pass after a small change may be sufficient. Convergence is tied to the reviewed commit and actionable findings: substantive repair requires a fresh pass; stylistic preference cannot keep the loop alive.
