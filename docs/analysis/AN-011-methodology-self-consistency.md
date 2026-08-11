---
id: AN-011
type: analysis
status: active
provenance: verified
reversal-cost: low
links: [G-001, P-008, AC-036, CAP-002, CAP-003, ARCH-003, PDR-013]
title: A green corpus can carry contradictory methodology claims
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

### F1 — Public, contributor, and capability guidance still describes the pre-v0.9 evidence boundary

`guide/adoption.md`, `guide/getting-started.md`, and `guide/operations.md` say that `clue validate` detects only one supported acceptance-criterion reference and cannot classify or count the positive and negative pair. The v0.9 judge does classify and count pairs for criteria declaring a test type. `guide/what-is-cliewen.md` simultaneously overstates that every active criterion requires positive and negative code tests, omitting the Human proof class and per-criterion `@draft`, and understates supported harvesting by omitting Cucumber.

The stale public inventory extends beyond those four pages. `guide/design.md` says Human proof has no home, omits Cucumber from the supported carriers, and repeatedly describes test evidence as mandatory for every active criterion. `guide/methodology.md` omits Human from the proof-class vocabulary, while `guide/change-loop.md` gives every active criterion the unconditional positive-and-negative-test rule. `CONTRIBUTING.md` repeats that rule as a normative contributor instruction. CAP-001's `design.md` preserves the old one-reference/counting boundary even though AC-036 claims that the guide distinguishes shipped support from methodology intent. The scaffolded capabilities README tells every adopter that each active criterion has positive and negative tests and otherwise the whole capability stays draft. All of these are live instructions or current design claims, not pinned historical descriptions.

CAP-002's active status note makes another false live capability claim: it says the mechanical AC↔test link remains future work even though AC-009, `checkACTests`, and its positive and negative Go tests already implement that rule. The incident therefore contradicts CAP-002 in addition to AC-036 and CAP-003.

The active ARCH-001 actor explanation and frontmatter-graph diagram also present a test and positive/negative pair as the only acceptance-criterion edge. The `internal/corpus/actests.go` implementation header says every criterion in an active file has a test and documents only Go and JVM carriers even though the code below it implements Human, `@draft`, classified/single-direction evidence, and Cucumber. These implementation-facing explanations are current carriers and need the same content guards as the public prose.

This contradicts AC-036's claim that the public operations guide distinguishes shipped support from methodology intent. The guide build and its content-anchor tests pass because they check page form and selected phrases, not agreement with the current evidence contract.

P-008/M-032 repaired this carrier set. CAP-001's `design.md` now states the corrected claim directly — "the guide distinguishes machine and method guarantees ... corrected after the v0.9 carrier audit" — and CAP-002's active status note now says the mechanical AC↔test link is implemented rather than future work. F1 is resolved.

### F2 — Both lifecycle skills contradict the Human and per-criterion draft decision

The canonical `clue-extract` source and both generated copies say activation is per criteria file and that a capability is the smallest phasing unit. ADR-033 and the current `clue-delta` skill allow an active criteria file to carry an individual unproven criterion tagged `@draft`; Human-class criteria also require no code reference. Generation tests prove the copies match their canonical source, so they faithfully reproduce the wrong current instruction.

The canonical `clue-verify` source also requires positive and negative tests for every active criterion or a draft capability. Its `.agents` and scaffold-template copies faithfully repeat that obsolete choice even though the same skill later recognizes Human proof in the acceptance brief. The repository pull-request template likewise recognizes Human proof in the acceptance brief but later requires positive and negative executable evidence for every changed active criterion; its scaffold source ships the same contradiction to adopters. The repair set therefore includes both canonical skill source templates and all generated copies, plus both pull-request templates and the scaffold synchronization guard; editing only a generated or repository-local copy would leave future scaffolds stale.

CAP-003's active criteria state namespaced acceptance IDs, JVM harvesting, and provenance linting; none of them states the activation-granularity rule the skill gets wrong. The reality edge is therefore recorded against the capability rather than against a criterion, and the repair must mint the criterion that names the rule before the contradiction has an executable carrier at all. F1's edge is narrower because AC-036 already claims the property the guide pages break.

P-008/M-032 minted that criterion: AC-054 (`docs/capabilities/CAP-003-extract/criteria.md`) now states the per-criterion activation-granularity rule the skills previously got wrong. F2 is resolved.

### F3 — The protected core still excludes the Human evidence path

PDR-013 and its current core-architecture carrier, ARCH-003, define the verifiable thread as ending in a test and require every durable claim to trace to executable evidence. G-001's title and body say that chain runs from goal to test and is mechanically enforced; its title propagates into the generated goals index. The repository `README.md` calls it the enforced thread, and `docs/README.md` presents the linter's red thread as ending at a test tag. The repository and scaffold `AGENTS.md` routing hubs repeat the test-ending core definition as a protected rule, while the scaffolded corpus README emits both that graph and the claim that every criterion is tied to tests. The CLI usage text, plugin-marketplace description, guide-site description, and guide landing page expose the same executable-only promise at shipped entry points; ARCH-001, `guide/design.md`, and `guide/methodology.md` repeat it for architecture and public readers. ADR-033 subsequently made the acceptance brief the proof carrier for a genuine Human-class criterion, so the current core wording excludes an evidence path the judge and lifecycle skill accept. Because C-013 protects what the thread connects, repairing this wording is a core-meaning change and needs the explicit human-accepted PDR already required by M-032, not an incidental prose edit.

PDR-019 is that PDR: it refines the thread's endpoint to acceptance evidence, and PDR-013 and `docs/architecture/core.md` now carry the refinement. F3 is resolved.

### F4 — Mechanical consistency is narrower than semantic carrier consistency

At the pinned revision, `go run ./cmd/clue validate` passed 108 artifacts, `--coverage` reported all seven capabilities covered, `go test ./...` passed, the guide built, and CI-scope tests passed. Those checks establish corpus form, generated-copy identity, executable evidence, and guide buildability. None inventories every carrier affected when a methodology decision changes a shared contract. The green result was correct under the implemented judge and still insufficient to detect F1 and F2.

This is the foundation's acknowledged semantic-merge edge in a smaller form: Git and generators can keep text synchronized without keeping meaning synchronized.

The repair inherits that limit. Criteria and content-anchor tests that reject today's stale phrases prove F1…F3 stay fixed; they do not detect the next carrier that falls behind a different shared contract, because no current mechanism derives the set of carriers a given contract has. Anything stronger needs a machine-readable statement of which carriers make which claim, which this audit does not have evidence to design. The general obligation therefore stays agent-enforced prose in the same sense as the constraints `clue validate` reports as awaiting machine checks.

### F5 — Not every old statement or unresolved analysis item is a present inconsistency

AN-002, AN-004, AN-006, AN-007, AN-009, and AN-010 have direct current consumers for their adopted conclusions. AN-005 deliberately called for no additional methodology adjustment. AN-003 remains partly consumed: evidence qualification landed, while clean-versus-prepared environment evidence, population eligibility, and adoption-governance prompts remain candidates. AN-008's core patterns are substantially consumed, while its adopter-hardening list remains explicitly outside completed P-007. Production operations, external catalogs, and foreign-document validation remain named future doors rather than contradictory promises.

AN-008, P-007's completed milestones, and the context sections of ADR-032 and ADR-033 describe the evidence model at their pinned historical stages. They are not live instructions and remain accurate history, so the carrier repair must not rewrite them as though the newer contract already existed.

## Finding and consumer

Cliewen needs a same-change carrier sweep whenever a methodology decision changes a shared contract, with mechanical guards for the stable claims that can be checked. P-008 consumes this incident: M-032 repairs the complete live inventory in F1…F3 and records that rule; M-033 and M-034 consume the evidence-backed analysis and extraction gaps; M-035 and M-036 require adopter evidence before interfaces are chosen for the remaining integration boundaries.
