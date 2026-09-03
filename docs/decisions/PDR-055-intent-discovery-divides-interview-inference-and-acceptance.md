---
id: PDR-055
type: decision
status: inferred
links: [G-012, P-021, M-089, CAP-003, CAP-009, ADR-044, ADR-010, ADR-035, PDR-020]
title: The agent interviews and infers, the CLI detects state, and the human accepts meaning
binds: adopter
author: agent
accepted-by: []
---

# PDR-055 — Who interviews, who detects, and who accepts

## Context

Producing a vision needs a conversation on a greenfield repository and an evidence sweep on a brownfield one. Neither is a deterministic operation, and both are easy to put in the wrong place: a `clue interview` command would put a semantic conversation inside the judge, and leaving everything to the agent's discretion would let inferred content reach the corpus looking like fact.

## Decision

**The division follows what each participant can be held to.** `clue` detects and reports structure and state — whether a vision exists, what its status and provenance are, which use cases name what — and conducts no interview and draws no inference. The agent interprets: it interviews on greenfield, infers from repository evidence on brownfield, cites its sources, and recommends. The human accepts meaning, and acceptance is the existing merge boundary rather than a new step.

**Greenfield discovery is a conversation, not a questionnaire.** The agent asks only questions whose answers would change the vision, the initial goals, the system boundary, or a candidate use case; it asks a few at a time, adapts to the answers, does not make the user learn Cliewen's vocabulary to answer, and stops when another question would not change anything. It then summarizes what it understood and asks for correction before treating any of it as accepted.

**Brownfield discovery reads before it asks.** Existing documentation, architecture and design material, any existing corpus, source and public APIs, tests, CLI help, configuration, package metadata, deployment definitions, existing decision records, in-repository change history, and examples are inspected first. The draft that follows cites the repository sources behind each material claim, separates what was observed from what was interpreted, and names contradictions and stale documents rather than choosing between them. Only what the repository genuinely cannot answer is asked.

**Code never establishes intent.** It demonstrates behaviour. A statement about why the product exists or whom it serves is asked or marked as an assumption; it is never derived from implementation structure. Where sources disagree and the disagreement changes durable meaning, the conflict is recorded and the decision is the human's.

**Drafted content is visibly drafted.** Whatever its origin, agent-produced intent is `status: draft` with `provenance: inferred` and a stated reversal cost, and open questions and assumptions stay in the artifact until answered. There is a documented escape hatch for a user who explicitly wants a draft from little — it changes how much is written, never how honestly it is labelled.

## Rejected: a `clue interview` or `clue discover` command

It would put a conversation inside the deterministic judge, and the judge's value is that its answer does not depend on who ran it. The agent already has the conversation; what it needed from `clue` was the state, which `clue validate --intent` supplies.

## Carrier

`internal/skills/source/shared/intent-discovery.md.tmpl` and the skills that route it, the `--intent` report, and the greenfield and brownfield walkthroughs in `guide/intent.md`.
