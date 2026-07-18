---
id: AN-003
type: analysis
status: verified
links: [P-002, M-007, PDR-005]
title: Calibration trial — robocode-api-bridge
---

# AN-003 — Calibration trial: robocode-api-bridge

## Classification

This is a calibration trial for P-002/M-007, not one of the milestone's qualifying foreign-soil trials. The repository is brownfield and open source, but it shares a maintainer with Cliewen; that maintainer can compare Cliewen's blind interpretation with the repository's actual intent. No files or hosted state in `robocode-api-bridge` were changed.

## Risk named

Cliewen may turn scattered brownfield evidence into a confident but incorrect goal, plan, or acceptance model when the repository's intended compatibility boundary lives partly in maintainer knowledge and sibling repositories.

## Evidence boundary

The blind pass used committed files and history from `robocode-dev/robocode-api-bridge` at `5c2fac5`, plus the repository's public issues as observed on 2026-07-18. The maintainer's locally modified root `TODO.md` was deliberately excluded. The pass did not inspect or rely on the sibling `tank-royale` checkout, and it did not ask the maintainer for intent before recording the interpretation.

## What was tried

- Reconstructed the likely goal, active direction, boundaries, and verification model from the root and component READMEs, committed TODOs, Gradle build files, compatibility harness, issue history, and recent commits.
- Ran `gradlew.bat build` in the existing checkout to test whether the repository's mechanical path supplies an independent truth signal.
- Compared the observed repository shape with the assumptions carried by `clue-analysis`, `clue-plan`, `clue-delta`, and `clue-verify`.
- Rejected materializing a Cliewen corpus in the target repository: this experiment is read-and-apply calibration, not adoption, and PDR-005 prohibits trial PRs and new extraction mappings.

## Blind interpretation

### Goal

Allow as many legacy Java Robocode robots as feasible to run on Tank Royale through a compatibility bridge, while keeping legacy quirks and workarounds out of the Tank Royale platform itself.

The target is broad practical compatibility, not exact behavioral equivalence. The README's statement that the adapter library is "100% identical" appears to describe the legacy API surface or source compatibility, while the surrounding prose explicitly rejects an expectation of 100% backward-compatible behavior.

### Current direction

The active work appears to be evidence-driven compatibility improvement: run legacy robots on classic Robocode and Tank Royale, compare results and errors, then fix bridge or upstream Bot API defects that explain material differences. Recent work concentrated on event-dispatch semantics and Bot API 1.0.2. The next likely move is to rerun previously flagged robots, then address surviving score gaps, file-I/O sandboxing, and team support.

### Boundaries

- Java legacy robots are in scope; the abandoned .NET plugin is not.
- The `robocode` and `gl4java` packages are imported legacy API sources and are intended to remain read-only; bridge behavior belongs under `dev.robocode.tankroyale.bridge`.
- Compatibility work belongs in the bridge when it represents a legacy quirk; Tank Royale's Bot API or server should change only when a fix benefits the clean platform without embedding legacy behavior.
- The wrapper translates robot archives into Tank Royale bot directories; team metadata is acknowledged but not implemented.

### Verification model

The compatibility harness is the strongest available behavioral evidence. It runs mirror battles on classic Robocode and Tank Royale, compares aggregate scores, and checks for Tank Royale-only errors. Its default 25% score-difference threshold is a diagnostic heuristic affected by stochastic variance, not yet a stable product-level acceptance boundary.

The conventional Gradle lifecycle is not an independent verification gate: all project `test` tasks report `NO-SOURCE`. The harness is manual, depends on absolute local paths and artifacts from sibling repositories, and its referenced `compatibility_report.md` is not committed.

## Findings

### F1 — Brownfield truth has distinct confidence layers

The repository offers at least four kinds of evidence with different authority and freshness: explanatory README prose, committed TODOs, public issues, and executable behavior. Cliewen's corpus can represent the result, but the current skills do not tell an agent to label which statements are observed, inferred, or maintainer-verified while conducting a brownfield trial. Without those labels, extraction can launder a plausible inference into durable truth.

### F2 — A successful maintainer build can be false onboarding evidence

`gradlew.bat build` succeeds in the existing checkout, but the committed tree does not contain `gradle/wrapper/gradle-wrapper.jar`, and the required Bot API 1.0.2 is resolved from `mavenLocal()` with instructions to publish it from the sibling Tank Royale repository. The build therefore confirms this machine's prepared environment, not reproducibility from a clean clone.

### F3 — The strongest test does not fit deterministic acceptance-criterion granularity

Compatibility is probabilistic and population-based: score variance, robot-specific behavior, timeouts, protocol compatibility, and known noise all affect interpretation. Forcing the entire promise into one positive/negative AC test pair would either create an oversized criterion or pretend a heuristic threshold is deterministic. The methodology can model deterministic adapter behavior as ACs and the population-level compatibility promise as a quality scenario, but the skills do not currently prompt for that split during brownfield analysis.

### F4 — The repository's working plan is distributed and partially stale

The root TODO says the event-dispatch redesign and Bot API upgrade are uncommitted, although the referenced changes are present in the latest commits. Open issue #5 and the missing compatibility report also lag or omit the newest experiment state. A plan inferred from only one carrier would be wrong; the likely direction above required reconciling prose with history and executable structure.

### F5 — Applying the change skill would introduce governance, not merely structure

The visible history has direct commits on `main` and no pull requests. `clue-delta` requires branch-plus-PR review, structured decision records, plan linkage, traceable ACs, and changelog entries for user-visible changes. Those rules may be valuable, but applying them would change how this project is governed. A trial must surface that consequence rather than describe `clue init` as neutral scaffolding.

### F6 — Calibration and foreign-soil trials answer different questions

Maintainer comparison can reveal whether the blind interpretation is accurate, but it cannot reveal how a maintainer unfamiliar with Cliewen experiences the workflow or whether shared assumptions went unnoticed. This repository is useful as a control case and remains ineligible for M-007's count.

### F7 — A population target needs an explicit eligibility model

The maintainer's working outcome is that more than 90% of legacy bots should run through the bridge and score close to their classic Robocode behavior against the same competitors. Some robots may be unsuitable evidence because they depend on obsolete Java methods, reflection tricks, unsupported runtime behavior, or capabilities deliberately outside the bridge. Excluding such bots informally would make the success rate vulnerable to hindsight and selection bias. The harness needs versioned per-robot metadata, predetermined eligibility rules, recorded exclusion reasons, and both raw and eligible-population results.

## Candidate methodology adjustments

These are findings, not adopted decisions.

- Add a trial rubric that separates repository observation, agent inference, and maintainer verification.
- Require a clean-clone or clean-environment check before treating a build as onboarding evidence.
- Prompt brownfield analysis to distinguish deterministic acceptance criteria from statistical or population-level quality scenarios.
- Require a population-level quality target to declare its corpus, eligibility rules, exclusions, sampling method, and uncertainty before using a percentage as evidence.
- Use an inspection ladder for third-party programs: black-box behavior first, intentionally bundled source where its terms permit inspection, and decompilation only with a recorded legal basis, interoperability need, and information boundary.
- Require trial findings to name governance changes that adoption would introduce.

## Rejected interpretations

- **"100% identical" means full behavioral compatibility.** Rejected because the same README explicitly says full backward-compatible behavior is not expected and may be infeasible.
- **The Gradle build proves the bridge works.** Rejected because it executes no conventional tests and depends on uncommitted local prerequisites.
- **The root TODO is the current plan.** Rejected because its first priority is already present in recent commits and its evidence is not synchronized with the current tree.
- **This trial can count toward M-007 because the repository is brownfield and public.** Rejected because PDR-005 additionally requires no shared maintainer.

## Maintainer comparison

The maintainer compared the blind interpretation with project intent on 2026-07-18.

### Confirmed

- The goal is to run as many legacy Java Robocode robots as practical on Tank Royale without moving legacy compatibility behavior into Tank Royale.
- "100% identical" means a one-to-one mapping of the legacy API, not behavioral equivalence. Exact behavior is impossible in some cases because the engines differ, including bounding-circle versus bounding-box geometry.
- Legacy `robocode`/`gl4java` compatibility belongs in the bridge; classic Robocode has already adapted that API on its side and is not the target of this work.
- The harness's 25% score threshold is diagnostic. Lower divergence is desirable as the bridge improves, but the achievable floor is not yet known.
- The intended outcome is provisionally that more than 90% of eligible legacy bots run through the bridge and score close to their classic Robocode behavior against the same competitors. The maintainer expects users might prefer 99%, but the absent testing and contribution effort from legacy bot authors makes that an unsupported commitment.

### Clarified

The root `TODO.md` is the closest current plan carrier, but it is a working backlog rather than a formally accepted plan. The maintainer's local edit marks its first item, the event-dispatch redesign and Bot API upgrade, as done. The remaining ordering supports the blind inference that retesting score-gap bots comes next, followed by file-I/O sandboxing, team support, and a fuller sweep.

The locally generated compatibility report is not suitable for setting the percentage target. It contains only 16 robots, reports 8 passes at the 25% threshold, and predates the newest event-dispatch work. It establishes neither that 90% is realistic nor that it is impossible. A representative versioned population and repeated matched battles are needed before the target can graduate from aspiration to a quality scenario.

### Robot inspection boundary

The initial comparison overstated the boundary by treating reverse engineering as categorically prohibited without author permission. Under [§§ 36–37 of the Danish Copyright Act](https://www.retsinformation.dk/api/pdf/238630), a person entitled to use a program may observe, study, or test it during acts they are entitled to perform, and code reproduction or translation may be permitted without separate author permission when indispensable for interoperability and confined to the statutory conditions. [Articles 5–6 of the EU Software Directive](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:32009L0024) carry the same distinction. This is a jurisdiction-specific research finding, not legal advice.

A coding agent has no separate exception and creates no new prohibition: it can act on behalf of an entitled user only within the same permission, purpose, and information-use boundaries. The agent's identity is therefore not the deciding factor.

Some Robocode archives intentionally include their source code. Reading that bundled source is source inspection rather than decompilation and is the preferred white-box debugging path when the archive's license and terms permit it. Availability of source does not by itself grant permission to copy, modify, redistribute, or publish the bot's protected expression or algorithm.

Black-box runtime evidence remains the safest default and can distinguish many bridge failures from robot incompatibilities by running the robot in environments the tester is entitled to use, supplying controlled Robot API inputs, and comparing observable calls, events, outputs, errors, timing, and scores. The CJEU's [SAS Institute judgment, C-406/10](https://curia.europa.eu/juris/liste.jsf?num=C-406/10), confirms that software functionality is not itself copyright-protected and that an entitled user may observe, study, and test program behavior without accessing or decompiling its code, provided the user's acts stay within the permitted use and do not copy protected expression.

Per-robot metadata should record the robot and version, test environment, observed failure, suspected category, evidence, eligibility status, explicit exclusion reason, inspection mode (`black-box`, `bundled-source`, or `decompiled`), and the permission or statutory basis for any source-level inspection. Eligibility categories might include obsolete Java dependencies, reflection or runtime tricks incompatible with modern Java, unsupported team behavior, and failures demonstrably present outside the bridge. To respect the community's interest in protecting bot strategies even where functional testing is lawful, findings should expose bridge-facing evidence rather than publish algorithm details. The categories and inspection ladder are proposed analysis output, not adopted policy. No legacy robot binary was decompiled during this calibration.

## Calibration result

The blind pass recovered the project's goal, major boundary, current direction, and diagnostic nature of the score threshold with useful accuracy. It missed the explicit population-level success target and had no way to distinguish confirmed intent from plausible inference until the maintainer comparison. The calibration therefore supports F1 and F7: brownfield trials need confidence labels and an explicit eligibility model before their findings can become durable methodology or percentage-based evidence.
