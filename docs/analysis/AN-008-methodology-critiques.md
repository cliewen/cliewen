---
id: AN-008
type: analysis
status: active
links: [P-007, ARCH-003, G-001, C-012, ADR-005, ADR-006, ADR-008, C-006]
title: Two independent critiques land on the same four half-built patterns
---

# AN-008 — Two independent critiques land on the same four half-built patterns

## Evidence boundary

Two reviews were received by the maintainer on 2026-07-25 as model conversations and supplied to this analysis in full. The first is a methodology-level review of Cliewen as published (six findings). The second is a migration assessment produced while evaluating Cliewen for a production system in an enterprise setting (nineteen findings, prioritized P0–P2); that system's repository is not publicly accessible and, unlike the named-but-private target in [AN-002](AN-002-model2diagram-extraction.md), stays unnamed here for confidentiality; the analysis README's convention applies either way — readers should not expect its artifacts to resolve. The transcripts are evidence of the reviewers' reasoning, not versioned repository artifacts. Every claim below about Cliewen's current behavior was re-verified against this repository at accepted `main` commit `4373909`, on Windows 11 Pro amd64 with Go 1.26.5; claims about the private system are reported as the assessment states them and are not independently reproducible.

## The finding

Twenty-five findings reduce to four patterns, and all four are the same defect: each of the three core elements ([ARCH-003](../architecture/core.md)) is half-built, and in every case the built half is the agent-facing half.

### Pattern A — the human half of the merge boundary is empty

The corpus carries roughly 5,600 words of generated skill instruction for the agent side of a change and zero words telling the human what to verify at merge. [C-012](../constraints/C-012-agents-never-merge-own-changes.md) states that the merge is a human act; nothing states what the act checks that the validator and the review loop did not. Without stated content the boundary degrades to rubber-stamping, and "machines enforce form; humans verify meaning" becomes theatre. The migration assessment's requests for a home for human-verified evidence and for an optional pre-implementation review checkpoint are the same gap seen from the adopter's side: humans have no defined verification surface anywhere in the loop.

### Pattern B — the judge proves less than the prose promises

The `clue-verify` checklist asks for positive and negative tests per active criterion; `clue validate` enforces at least one reference. JVM evidence harvests at file level; Cucumber and Gatling evidence is not harvested at all; deliberately human-verified criteria have no machine-visible home, so their capabilities either stay `draft` or overstate. Cliewen already names this gap — `enforcement: agent` constraints are the declared promotion backlog — so the assessment's P0 list is not a new crack: it is that backlog, priced by a real prospective adopter, with a due date attached.

### Pattern C — the graph only accumulates

Nothing in the corpus bounds, reverses, or consumes state. The born-`inferred` population can only grow, so its counter stops being read. No supersession edge exists: criteria get tombstones, decisions do not, and when a decision is reversed nothing in the graph answers what was downstream of it. No edge returns from reality, so a fully green corpus can describe a wrong product. Retirement flips a status field but removes nothing from the mandatory read path, so read cost climbs with every artifact ever written. One family: a write-mostly corpus.

### Pattern D — everything is priced for the authoring repositories

Every rule in the corpus was born from an incident in a repository the maintainer controls, and the tooling shows it: the full loop is the modal path for an ordinary fix and the guide never defends that; extraction is destructive with no rehearsal pass; the CI wall hard-codes runner, action versions, and binary paths; managed skills reject local extension by design with no blessed companion mechanism; the `.claude/skills` mirror is written by `init` but never validated. The cost of the methodology has been asserted ("deliberately visible overhead") and never measured.

## Collisions the critiques did not see

Three requested fixes meet standing decisions and need explicit resolution rather than direct implementation — two as genuine collisions, one dissolving on inspection:

- Structured evidence comments (`// AC:`) were rejected by [ADR-005](../decisions/ADR-005-test-reference-convention.md) because nothing attaches a comment to a test but proximity. Cucumber Gherkin tags fit the existing tag philosophy; the no-tag-mechanism case (Gatling and kin) needs either an ADR revisiting ADR-005 for that class or a naming-convention answer.
- An extraction dry-run cannot be a `clue` flag because [ADR-008](../decisions/ADR-008-extraction-is-a-skill.md) made extraction a skill; the fix is a skill-contract change (a report-only first pass), not a command.
- An incident artifact type is not needed, and the reason is positive rather than a collision: episodes already live in analysis records — [AN-007](AN-007-review-handoff-gap.md) is an incident that fed C-012 and a LOG-001 readiness row through the carriers it links — and that division of labor is what keeps decision prose episode-free under [C-006](../constraints/C-006-adrs-timeless-with-carrier.md). The shape wants formalizing, not replacing.

## What the book already answers

Cliewen is born from the Intent Engineering book, and four of its published designs map directly onto the findings: the "after" verification gate's advisory checks are the missing human merge checklist (does the test match its GIVEN/WHEN/THEN; are invalidated records superseded); its review chapter carries the accountability principle ("review catches divergences between spec and implementation, not between the spec and reality") and names the competence-heuristic blind spot that agent output triggers; its coverage-pair convention includes the single-direction exemption that makes pair enforcement sane; its per-scenario test-type field is the design-time proof-class declaration the assessment asks for. Importing these keeps Cliewen and the book coherent instead of drifting.

## Finding and consumer

The critiques are one finding, not twenty-five: the core's guarantees are asserted by agent-facing machinery and unbacked on the human-facing and bounding sides. [P-007](../plans/P-007-core-hardening.md) consumes this analysis; its milestones M-024…M-029 answer patterns A through D in that order, and its out-of-campaign list holds the adopter-hardening items that a real pilot must price first.
