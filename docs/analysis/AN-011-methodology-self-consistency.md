---
id: AN-011
type: analysis
status: active
provenance: inferred
reversal-cost: low
links: [G-001, P-008, AC-036, CAP-003, ADR-032, ADR-033, C-006, P-007]
title: A green corpus can carry contradictory methodology claims
reality: contradicted
---

# AN-011 — A green corpus can carry contradictory methodology claims

## Risk

Cliewen may pass every mechanical check while its public guide, generated skills, and normative corpus disagree about what the methodology requires.

## Evidence boundary

The audit inspected accepted `main` at `ef120203b816be259397d03a40b858fc4764e494` on 2026-07-28, using Windows 11 build 26200, PowerShell, Go 1.26.5, Node.js 24.0.0, and npm 11.18.0. Repository-source verification used `go run ./cmd/clue` because the installed `clue` was v0.6.0 while the repository and generated skills were v0.9.0; the released binary correctly reported that version drift and was not used to judge the v0.9 corpus.

Observed evidence is limited to committed repository state and reproducible local commands. The disposition of earlier analyses is an inference from their consumers and current carriers, not a claim about unrecorded maintainer intent. Deliberately deferred work is reported separately from contradictory current claims.

## What was tried

- Read every document under `docs/analysis` and traced its findings to current plans, decisions, constraints, capabilities, skills, guide pages, and implementation.
- Ran the v0.9 source validator, coverage and reality-gap reports, the complete Go suite, the guide build, and the CI-scope tests.
- Compared the evidence contract in ADR-032, ADR-033, CAP-002, `clue-delta`, `clue-extract`, and the public guide.
- Rejected treating every deferred capability as an inconsistency: an explicitly closed door is coherent when every carrier describes it as closed.

## Findings

### F1 — The public evidence boundary still describes the pre-v0.9 judge

`guide/adoption.md`, `guide/getting-started.md`, and `guide/operations.md` say that `clue validate` detects only one supported acceptance-criterion reference and cannot classify or count the positive and negative pair. The v0.9 judge does classify and count pairs for criteria declaring a test type. `guide/what-is-cliewen.md` simultaneously overstates that every active criterion requires positive and negative code tests, omitting the Human proof class and per-criterion `@draft`, and understates supported harvesting by omitting Cucumber.

This contradicts AC-036's claim that the public operations guide distinguishes shipped support from methodology intent. The guide build and its content-anchor tests pass because they check page form and selected phrases, not agreement with the current evidence contract.

### F2 — The extraction skill contradicts the per-criterion draft decision

The canonical `clue-extract` source and both generated copies say activation is per criteria file and that a capability is the smallest phasing unit. ADR-033 and the current `clue-delta` skill allow an active criteria file to carry an individual unproven criterion tagged `@draft`; Human-class criteria also require no code reference. Generation tests prove the copies match their canonical source, so they faithfully reproduce the wrong current instruction.

CAP-003's active criteria state namespaced acceptance IDs, JVM harvesting, and provenance linting; none of them states the activation-granularity rule the skill gets wrong. The reality edge is therefore recorded against the capability rather than against a criterion, and the repair must mint the criterion that names the rule before the contradiction has an executable carrier at all. F1's edge is narrower because AC-036 already claims the property the guide pages break.

### F3 — Mechanical consistency is narrower than semantic carrier consistency

At the pinned revision, `go run ./cmd/clue validate` passed 108 artifacts, `--coverage` reported all seven capabilities covered, `go test ./...` passed, the guide built, and CI-scope tests passed. Those checks establish corpus form, generated-copy identity, executable evidence, and guide buildability. None inventories every carrier affected when a methodology decision changes a shared contract. The green result was correct under the implemented judge and still insufficient to detect F1 and F2.

This is the foundation's acknowledged semantic-merge edge in a smaller form: Git and generators can keep text synchronized without keeping meaning synchronized.

The repair inherits that limit. Criteria and content-anchor tests that reject today's stale phrases prove F1 and F2 stay fixed; they do not detect the next carrier that falls behind a different shared contract, because no current mechanism derives the set of carriers a given contract has. Anything stronger needs a machine-readable statement of which carriers make which claim, which this audit does not have evidence to design. The general obligation therefore stays agent-enforced prose in the same sense as the constraints `clue validate` reports as awaiting machine checks.

### F4 — Not every unresolved analysis item is a present inconsistency

AN-002, AN-004, AN-006, AN-007, AN-009, and AN-010 have direct current consumers for their adopted conclusions. AN-005 deliberately called for no additional methodology adjustment. AN-003 remains partly consumed: evidence qualification landed, while clean-versus-prepared environment evidence, population eligibility, and adoption-governance prompts remain candidates. AN-008's core patterns are substantially consumed, while its adopter-hardening list remains explicitly outside completed P-007. Production operations, external catalogs, and foreign-document validation remain named future doors rather than contradictory promises.

## Finding and consumer

Cliewen needs a same-change carrier sweep whenever a methodology decision changes a shared contract, with mechanical guards for the stable claims that can be checked. P-008 consumes this incident: M-032 repairs F1 and F2 and records that rule; M-033 and M-034 consume the evidence-backed analysis and extraction gaps; M-035 and M-036 require adopter evidence before interfaces are chosen for the remaining integration boundaries.
