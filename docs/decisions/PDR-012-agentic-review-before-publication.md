---
id: PDR-012
type: decision
status: verified
links: [G-001, ARCH-002, PDR-007, PDR-040, PDR-042, C-012]
title: Every chosen full change receives an automatic agentic review before readiness
author: agent
accepted-by: Flemming N. Larsen (2026-07-21, implementation conversation; finding-grounding amendment approved 2026-07-22, follow-up conversation)
---

# PDR-012 — Agentic review before publication

> **Amended by [PDR-040](PDR-040-push-is-durability-ready-is-explicit.md):** the review gate is before the ready mark; pushing is durability and claims no readiness.

> **Amended by [PDR-035](PDR-035-bounded-agentic-review-loop.md) and [PDR-036](PDR-036-review-loop-budget-and-human-checkpoint.md):** the loop is bounded, computed-figure findings are advisory, and outstanding blocking findings at the maximum are reported to the human rather than silently permitting readiness.

> **Scoped by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** mandatory agentic review applies to a chosen full loop; simple work has no mandatory Cliewen review.

## Context and problem statement

An implementing agent is anchored to its chosen approach, while a fresh read-only challenge can expose defects, missing evidence, or intent mismatch. Making that challenge automatic improves the committed candidate without pretending it is human review or requiring a non-converging ritual.

## Decision outcome

**Every chosen full change receives an automatic adversarial review of its committed candidate before the ready mark, using a context-isolated reviewer when the host supports it and a disclosed in-context fallback otherwise.** The reviewer inspects the complete base diff, durable corpus, tests, constraints, and quality scenarios, and reports blocking actionable findings or advisory observations with severity, location, evidence, and remediation.

A finding is actionable only when grounded in an operative requirement or declared intent and tied to a concrete consequence. The reviewer is read-only. The implementing context repairs blocking findings on the same change, reruns applicable checks against the repaired commit, and obtains a new review; each repair invalidates the prior clean pass. The bounded loop stops on a current-commit pass with no blocking findings, or reports the human checkpoint when the budget is exhausted with blocking findings outstanding.

The final handoff names the review mode, reviewed commit, pass count, and advisory findings left open. The generated `clue-verify` review loop, `clue-delta` handoff, boundary fragment, and generator tests carry this full-loop gate; it does not claim that the CLI proves reviewer independence.
