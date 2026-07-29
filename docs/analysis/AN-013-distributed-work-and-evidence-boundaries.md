---
id: AN-013
type: analysis
status: active
provenance: inferred
reversal-cost: low
reality: contradicted
links: [P-008, M-036, G-001, CAP-003, AC-011, ADR-009, ADR-006, ADR-011, ADR-030, C-011, C-012, C-013, PDR-007, AN-011, AN-012]
title: The corpus cannot say what is accepted, where its evidence actually lives, or which repository a reference means
---

# AN-013 — The corpus cannot say what is accepted, where its evidence actually lives, or which repository a reference means

## Risk

Cliewen may be about to build three distributed-work interfaces — stacked changes, cross-repository evidence, and external-tracker metadata — on speculation, and each has an obvious implementation that damages the core [C-013](../constraints/C-013-core-changes-need-decision.md) protects. Stacked changes invite treating a pull request's base as accepted meaning; cross-repository evidence invites importing a foreign forge's green check as proof; tracker metadata invites recording forge identity as corpus truth. The prior risk is the mirror image: all three boundaries are already being crossed in practice, so *not* looking is not neutral either.

## Evidence boundary

Cliewen was pinned at `980ec9b488682e6c3724b8daff61bee32fc9494c`, the accepted `main` this change branched from (the merge of PR #85, 2026-07-29). The judge used throughout was a source build of the branch commit `b3798bb16d61e4b319b9f946b04846fa97d73633`; it reports version `dev` and therefore **skips the skill-drift comparison by design** ([ADR-011](../decisions/ADR-011-version-stamping.md), AC-033). No released binary produced any result recorded here.

The adopter was the public Robocode Tank Royale repository at fetched `origin/main` commit `384d27d55176a2d2ad4668ac381852e629e4540a` (2026-07-29), the merge of its own pull request `ch-006-upgrade-cliewen-0.9.0`. **This is a newer revision than [AN-012](AN-012-adopter-configuration-cost.md) pinned**: the repository has since upgraded from `0.5.1` to `CLUE_VERSION: 0.9.0` with its five managed skills re-vendored at `0.9.0`, and head `clue validate` now reports it clean (116 artifacts). Where a claim here differs from AN-012's, the difference is the upgrade, not a correction.

**Every result below is from a prepared environment, and none of it is onboarding-reproducibility evidence.** The prerequisites were existing local clones of both repositories and a Go 1.26.5 toolchain used to build the judge from source. Reproductions ran on Microsoft Windows NT 10.0.26200 under Git Bash, on a case-insensitive filesystem.

Scenario 1 ran in a **disposable local clone** of the pinned Cliewen revision inside a scratch directory. Nothing was pushed, no pull request was opened, and no forge state was created — so nothing here measures GitHub's own stacked-pull-request behaviour. Every reading of the adopter was read-only; `clue validate` does not write, and no file in either repository was modified.

Counts below are deterministic outputs of one judge over pinned revisions, not samples. There is no population, no sampling method, and no uncertainty interval to report, and they say nothing about how another repository pair would measure. Cliewen accepts changes as merge commits (`main`'s history is a chain of `Merge pull request …`); the git mechanics in scenario 1 are specific to that and were not reproduced under squash-merge.

## What was measured

### Scenario 1 — an authorized dependent change is indistinguishable from accepted truth

Change A added a decision record (`ADR-900`, `status: inferred`, `author: agent`) and its index line. Change B branched from A and added an analysis record that links `ADR-900` and reasons from it in prose as settled corpus truth. This is exactly the situation [PDR-007](../decisions/PDR-007-review-boundary.md) and [C-012](../constraints/C-012-agents-never-merge-own-changes.md) permit only under explicit human authorization.

**B's corpus is green, and nothing in it says A is unaccepted.** `clue validate` on B reports `OK (115 artifacts, 21 inferred decision(s) awaiting verification …)`; `clue validate --forbid-changes` reports the same. `clue context AN-900` emits `ADR-900` with `status: inferred` — **the same status a merged agent-authored decision carries** until a human promotes it. The status vocabulary has no way to distinguish "proposed on a branch nobody has accepted" from "accepted by merge, awaiting human verification". Accepted-ness is `main`-membership, and `main`-membership is a git and forge fact that no corpus artifact records.

**Merging the dependent change alone binds the base's meaning.** `git diff main...ch-b-dependent` is four files and 33 insertions and includes both of A's files: a reviewer of B is shown A's decision inside B's diff. Merging B into a `main` that never accepted A puts `ADR-900` on `main`, and validation stays green. The human merge boundary was satisfied — for B — and A's decision became binding without ever being accepted on its own. Neither the corpus nor the judge reports anything.

**A base revised in review is absorbed silently.** Review then reversed `ADR-900`'s outcome on A's branch, and A merged. B did what the review boundary requires of a published change — merged the newer accepted `main` into its branch — and the merge succeeded with no conflict, one changed line, and a green `clue validate`, while B's analysis prose still reasons from the meaning A had before review. Git merged the text and could not see that the meaning had inverted. This is `clue-delta`'s own closing line arriving as a measured result.

What the branch-and-pull-request boundary does supply is real but narrower than it looks: linear ancestry, a reviewable diff, and — where the forge enforces it — a required check. What repository-local validation supplies is within-repository ID uniqueness and link resolution, both of which held. What neither supplies is any representation of the dependency itself. **No artifact anywhere records that B is rooted on unaccepted work, which change is its base, or that a human authorized it.** Today that authorization lives in conversation, or in the pull request's base branch — that is, in private agent memory or in forge state, the two places `clue-delta`'s durable-work-state rule and Cliewen's system-of-record rule respectively forbid.

### Scenario 2 — acceptance evidence already spans repositories, and cannot be expressed

Two live cases exist in the current corpus, in opposite directions.

**The emitted CI wall.** [CAP-001](../capabilities/CAP-001-onboarding/README.md) emits a workflow that only ever functions in another repository. The adopter's copy differs from the current template by 71 changed lines, and its header comment now states the situation plainly: the file is a copy, upstream improvements do not arrive by themselves, and every step but the deliberate install divergence is kept in step by hand on each upgrade. AN-012's candidate 5 — the missing acceptance-brief gate — was repaired that way, by the adopter's own `ch-006` change, three days after AN-012 measured its absence. The cross-repository consistency obligation is therefore real, currently discharged by human discipline, and carried mechanically by nothing on either side.

**The delegated JVM purpose rule, which was never installed.** [ADR-009](../decisions/ADR-009-ac-id-namespaces.md) names three carriers, the third being "the ArchUnit purpose rule shipped by extraction (machine, in the adopting repo)", and delegates [AC-011](../capabilities/CAP-002-validate/criteria.md) — one declared purpose per test — to it on the JVM, because `clue` harvests JVM tags only at file level. The instruction was live in the `v0.3.0` skill pair the adopter's extraction actually used: `git show v0.3.0:.agents/skills/clue-extract/skill.md` reads "Where a JVM test framework is present, install an ArchUnit (or equivalent) rule enforcing one purpose tag per test". In the extracted repository there are **95 `*Test.kt`/`*Test.java` files, 64 `@Tag("…")` uses, none of them in AC-ID form, and no ArchUnit dependency or rule anywhere** — the only occurrences of the word are the vendored skill's own text and unrelated principles documents. The obligation was stated, the extraction ran, and the rule was not installed. Both corpora are green and neither reports the absence.

Nothing is producing a wrong verdict *yet*, and that is the point: the adopter's 257 extracted criteria are deliberately `draft`, which exempts them from the AC↔test wall, so no evidence is due. Its coverage report is `gap` for all thirteen capabilities while validation passes. When it promotes criteria to `active` — its own tagging door — file-level JVM harvesting will require evidence and per-test purpose will still be unenforced, because the carrier ADR-009 delegated to does not exist there.

**Neither direction can be written down.** In the pinned corpus, any foreign identifier in `links:` is a hard failure: a foreign criterion ID gives `link BR-001 resolves to no artifact`, and a fully qualified forge reference fails identically (`link robocode-dev/tank-royale#223 resolves to no artifact`). A criterion whose real proof is a run in another repository — added as `AC-901` with `Test-type: Integration` — fails three ways at once: `has no test`, `has no Integration positive evidence`, `has no Integration negative evidence`. The only honest expressions available today are `Test-type: Human`, whose proof is the acceptance brief a human signs, and `@draft` for a criterion not yet proven. Both are correct, and both discard the one fact worth keeping: *where* the evidence lives.

The reverse dependency is equally unexpressed. The adopter's wall downloads Cliewen's release asset by URL, on the append-only asset-name contract [ADR-030](../decisions/ADR-030-verified-install-scripts.md) records. The only mechanical cross-repository check that exists anywhere is `checkSkillVersions`, and it runs entirely inside the adopter, comparing two local things — vendored skill stamps against the binary in hand.

### Scenario 3 — a bare forge reference never named a repository, so it cannot survive a move

The pinned corpus contains **50 bare `#N` references across 25 files**, spanning at least four repository namespaces: Cliewen's own pull requests, `robocode-api-bridge` (`#5`), `hyperfine` (`#788`, `#844`), and `tank-royale` (`#218`–`#223`).

**One already resolves to the wrong thing.** [AN-003](AN-003-robocode-api-bridge-calibration.md) writes "issue #5" meaning the calibration target's issue 5. Rendered in Cliewen's own namespace, `#5` is a real merged pull request there — `CH-006: P-001 completed; P-002 'Cliewen leaves home' active` (confirmed with `gh pr view 5`). The reference is not broken; it is confidently wrong. The tank-royale citations behave the other way round: Cliewen's highest number today is 85, so `#218`–`#223` are dead now and will become wrong live links when its own numbering passes them.

**The stable-ID rule does not reach this.** `docs/README.md`'s "identity is the ID, the path is only the current address" is a within-repository guarantee. It says nothing about forge references, and corpus IDs are not unique across repositories either: `CAP-001` is Cliewen's onboarding capability and the adopter's Battle Runner; `C-002` is Cliewen's changelog constraint and the adopter's review boundary. A bare ID becomes ambiguous the moment it crosses a repository, which is precisely what scenario 2's candidates would have it do.

**The judge sees none of it.** `checkLinks` resolves frontmatter IDs only. The sole body-link check is `checkIndexes`, which looks inside generated index blocks and skips any target matching a URL scheme. No rule inspects a `#N`, and 26 files hardcode `github.com/cliewen/cliewen`. A rename or transfer leaves forge redirects behind, which is forge behaviour rather than corpus truth and exactly the dependency Cliewen must not acquire; a mirror to another forge, or any reader outside that forge's web UI, gets no redirect at all. No move was performed here, and the finding does not need one: a reference that never named a repository cannot survive changing which repository it is read in.

## Findings

**F1 — Accepted-ness is not corpus data, so a dependent change cannot be judged locally and its merge can bind meaning nobody accepted.** Measured: a green dependent corpus, `inferred` used for two different things, A's decision landing on `main` through B's merge, and a reversed base absorbed by a clean merge. The gap is not in git and not in the forge; it is that no artifact records the base, the dependency, or the authorization, so nothing local can distinguish proposed meaning from accepted meaning.

**F2 — Cross-repository evidence exists, is discharged by hand, and has no expressible form.** Measured: a 71-line divergence in an emitted wall kept in step manually, and a delegated JVM carrier that the one extraction which triggered it never installed. `Test-type: Human` and `@draft` are the honest fallbacks and remain so; what is missing is a way to *name* foreign evidence and its revision without importing a foreign verdict.

**F3 — External references carry no repository identity, and one is already wrong.** Measured: 50 bare `#N` references across four namespaces, one resolving to an unrelated live pull request in Cliewen's own namespace, no validator rule that can see any of them, and cross-repository ID collisions with divergent meaning.

Because F2's second half is a claim the corpus states and reality does not bear out — ADR-009's third carrier has no instance in the only repository whose extraction was supposed to install it — this record carries `reality: contradicted` and links the capability that failed to deliver it ([CAP-003](../capabilities/CAP-003-extract/README.md)) alongside the delegated criterion ([AC-011](../capabilities/CAP-002-validate/criteria.md)) and the decisions that delegated it.

## The rejection boundary

Two rules bound every candidate below, and rule out the most convenient designs for each scenario.

**Nothing may weaken the human merge boundary.** Scenario 1 measured the concrete damage: merging a dependent change silently binds its base. Automating stacked merges, treating a pull request's base branch as accepted meaning, or letting an agent land a chain "in order" all make that damage systematic rather than accidental.

**Nothing may make forge state the system-of-record.** A foreign green check, a base-branch pointer, a pull-request number, and a rename redirect are all forge state. They may *enforce* — branch protection is exactly that, and Cliewen already relies on it — but they may not *mean*. A judge that resolves a reference by querying a forge is also no longer deterministic, offline-reproducible, or pinned to a revision.

## Rejected

- **Importing an adopter's green check as acceptance evidence.** Makes forge state the system-of-record, and is unpinned besides: a check is green about a revision, and Cliewen cannot see which revision, nor whether the criterion's scenario is what ran.
- **Cross-repository link resolution inside `clue validate`.** Network access in the deterministic judge, a verdict that depends on another repository's current state, and validation that stops working offline. Rejected regardless of how cross-repository evidence is eventually named.
- **Automated stacked-change merging, or treating a base branch as accepted meaning.** Measured to bind unaccepted meaning; see the boundary above.
- **Forbidding dependent changes outright.** The review boundary already permits them under explicit authorization, and the pressure is genuine — P-008 kept M-032 whole as one wide change precisely because splitting it would have put a self-contradictory state on `main`. Banning the pattern would push authors to bundle unrelated meaning into single changes instead, which is the same failure wearing different clothes.
- **Recording forge identity as corpus truth to fix scenario 3.** The defect is that references are unqualified, not that the corpus lacks a pull-request registry. Qualification is a text convention; a registry would be a second source of truth that ages.
- **Deriving accepted-ness by having `clue` query the forge for merge status.** Same two objections at once: forge state as meaning, and a non-deterministic judge.

## Candidates for a successor plan

Independently routable, routed by reversal cost, as candidates only — this change decides none of them.

1. **PDR candidate — a dependent change declares its base and its authorization as durable data.** The proposal names the change it is rooted on, records the human authorization that permitted it, and states what its own merge would bind that its base has not yet had accepted; the acceptance brief repeats that sentence where the merging human reads it. Process, so a PDR under [C-011](../constraints/C-011-decision-records-typed.md); expensive to reverse once skills, templates, and adopter workspaces carry it, and core-adjacent under C-013 because it touches what a merge binds.
2. **ADR candidate — foreign acceptance evidence gets a named form that is explicitly unverifiable here.** A criterion whose proof is a run in another repository names the repository, the pinned revision, and the identifier, and the judge treats it as *named but unproven locally* — never as coverage, never as an imported verdict. `Test-type: Human` stays the proof of record; the addition is the pointer. Corpus format and validator behaviour, so an ADR, and a high-cost one: it becomes a public contract for every adopter's criteria files.
3. **ADR candidate — external references are qualified, and a bare forge number is a lint failure.** `owner/repo#N` or a full URL, enforced by a new check, with the 50 existing references repaired in the same change and cross-repository corpus IDs qualified by repository. Corpus format plus a lint rule, so an ADR; the convention is cheap to write and expensive to unwind once adopters' corpora follow it.
4. **ADR candidate — extraction stops delegating an obligation nothing verifies.** ADR-009's third carrier either becomes something extraction demonstrably installs and reports, or the delegation is withdrawn and the JVM per-test purpose gap is stated openly as unenforced. Amends an existing accepted decision's carrier claim, so an ADR.
5. **For the adopter, not a Cliewen decision** — its JVM suites have no per-test purpose rule, and its wall remains a fork kept in step by hand. Recorded here because Cliewen's own evidence found it; candidate 4 addresses the delegation, and AN-012's candidate 1 addresses the fork.

**Named consumer:** the successor campaign **P-009**, to be proposed after P-008 closes. This analysis has no consumer inside P-008: M-036 forbids implementing an interface, and it is the campaign's last milestone. Candidates 1 and 2 are the load-bearing ones; candidate 3 may ride any change touching the corpus, and candidate 4 belongs with whatever next revisits extraction.

## What this analysis does not establish

One repository pair, one forge, one prepared environment, no clean-environment result, and an adopter maintained by Cliewen's own author. Scenario 1 was reproduced entirely in a local clone: nothing here measures GitHub's stacked-pull-request behaviour — base retargeting on merge, squash-merge mechanics, or required-check re-runs on a chain — and Cliewen's own merge-commit style is the only merge style exercised. No repository was renamed, transferred, or mirrored; the tracker finding rests on what a reference contains, not on measured redirect behaviour. The delegated-carrier gap is one extraction into one JVM repository, and it establishes that the obligation went unperformed and unnoticed, not how often that happens. Nothing here prices any of the three candidate interfaces; each still needs its own design work before it is proposed.
