---
id: CH-093-rehearsal
type: analysis
status: draft
links: [CH-093, P-009, M-042, PDR-004, PDR-020, PDR-021, ADR-037, AN-013, AN-014, CAP-003]
title: PDR-020 rehearsal report for Robocode Tank Royale
provenance: inferred
reversal-cost: low
---

# CH-093 — PDR-020 rehearsal report: Robocode Tank Royale

## Risk and scope

The risk tested was that a real brownfield repository can appear ready after its corpus is converted while stable source meaning, executable evidence, CI operation, or deletion and merge boundaries remain unproven.

This is a report-only rehearsal for P-009/M-042. It changes no file in the Tank Royale checkout, no durable Cliewen corpus, no target tests or routing, and no hosted state.

## Pins and evidence boundary

M-042 requires a human-selected repository: Robocode Tank Royale was named by the human initiating this change as the representative target, and no agent chose it.

Two revisions are pinned and every claim below names which one it reads. The **head pin** is `384d27d55176a2d2ad4668ac381852e629e4540a`, the public adopter revision also pinned by AN-013, inspected through a local clean mirror whose `main` is that commit; the revision is the reproducible identity of this evidence, not the mirror's local path.

The current source tree has already been adopted, so the **source pin** is `4e579878fd6667fab3b75515b6e68135a935c8df`, the first parent of adoption merge `b7fb320ccec1f4742ef923cb315e7dd84f7e824f`; this is the last source state containing `openspec/` before CH-001 removed it. The proposal named only the head pin, so the source pin is this report's addition to the declared scope and the digest carries it.

This rehearsal is therefore retrospective, not pre-mutation: the extraction it describes already happened between the source pin and the head pin. PDR-020's checkpoint — a live `openspec/` tree previewed for deletion, a mapping proposed before it is written, and a human authorization gate between report and mutation — is reconstructed from history here rather than exercised end to end. That limit is recorded again in the findings and is the reason no exit item below may be read as proving the checkpoint works on a repository that has not yet been converted.

Read-only inspection ran on Microsoft Windows NT 10.0.26200.0 with PowerShell 7.6.4, Go 1.26.5, Java 26, Node.js 24.0.0, Python 3.14.2, and .NET 8.0.423.

These are prepared-environment results, not clean-onboarding evidence: the existing checkout, local Gradle caches and build outputs, existing Node dependencies, and the local Go source tree were prerequisites. The target's contributor instructions require JDK 17+ (`CONTRIBUTING.md`, `DEVELOPMENT.md`, which also names JDK 11 for the Java/JVM Bot API toolchain), while its architecture report narrows the build range to JDK 17–21 (`docs/architecture/report/architectural-health-report.md`); this machine has JDK 26, inside the instructions' stated floor and outside the architecture report's ceiling.

The target was clean before inspection and clean after inspection, including after the failed Gradle test command; no tracked or staged target file changed.

## Source inventory at the source pin

Every count in this section was read at the source pin. The source corpus had `openspec/project.md`, `openspec/config.yaml`, twelve `openspec/specs/*/spec.md` capability specs, one pending `openspec/changes/add-typescript-bot-api-npm-publish`, and an `openspec/changes/archive` tree.

The twelve specs contain 107 `Requirement` headings and 257 `Scenario` headings; the scenarios carry no stable criterion IDs, so an extraction cannot preserve source scenario identities and must mint them deterministically in capability namespaces.

The source had 41 numbered MADR decisions (`0001` through `0041`) under `docs/decisions/`.

The source carried `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, `.junie/guidelines.md`, twelve `.agents/instructions/*` files, and `.principles` files at the repository and module levels; routing reconciliation therefore covered more than the flagship assistant file.

The source language inventory was 132 Java files, 356 Kotlin files, 133 Python files, 113 TypeScript files, and 118 C# files; 150 of those files were test-like by path or name (33 Java, 36 Kotlin, 35 Python, 17 TypeScript, and 29 C#).

The source had no per-test Cliewen purpose tags. Existing test metadata was a separate test-registry convention and could not be treated as Cliewen criterion evidence without a mapping and normalization pass.

## Current corpus and ID mapping

Every count in this section was read at the head pin. The adopted corpus has thirteen capabilities with prefixes `BR`, `BFD`, `BAU`, `GBP`, `GBC`, `PRO`, `PBA`, `RCV`, `TCS`, `TFS`, `TBA`, `UD`, and `TNP`.

The current criteria files contain 263 scenarios, all with `status: draft`; the first twelve capabilities carry the 257 source scenarios and `TNP` carries six scenarios derived from the pending TypeScript publication change.

The current corpus therefore demonstrates deterministic minting from a source with no IDs, but this rehearsal does not exercise preservation of an existing multi-segment or letter-suffixed source ID because the source pin did not declare one. M-042 asks the report to inventory stable-ID preservation, and this inventory meets that conjunct by establishing that the source had no identities to preserve. What it cannot do is exercise preservation: ADR-037's extended grammar surviving a migration is untested by this target, and the findings record that as a limit of this target rather than as an unmet exit item.

All current criteria remain draft, so current `clue validate --coverage` reports `gap` for every capability and there is no active classified AC coverage for the rehearsal to mistake for proof.

## JVM evidence sample

The sample was purposive, not random, and its population is the JVM test tree at the head pin. Eligible files were the Java and Kotlin test files carrying a legacy `TR-API-*` tag at both class and method scope, because a mixed carrier is the shape most likely to produce a false credit; two such files were read in full. The sample is therefore a worst-case probe, not a basis for a frequency claim over the 69 test-like Java and Kotlin files counted at the source pin.

The sample read `bot-api/java/src/test/java/dev/robocode/tankroyale/botapi/ProtocolConformanceTest.java`, which has a class-level `TCK` tag and eleven method-level `TR-API-TCK-007` through `TR-API-TCK-017` tags across eleven executable tests.

It also read `bot-api/java/src/test/java/dev/robocode/tankroyale/botapi/BotRunFirstTurnTest.java`, which repeats `TCK` and `TR-API-TCK-004` at class and method scope across two executable tests.

Both files contain mixed executable tags but no Cliewen criterion identity, declared `Test-type`, or positive/negative direction on the same executable; no `// AC:` proximity carriers were found in either file.

The current validator credits neither file as classified AC evidence. That result generalizes beyond the two files without a wider sample, because its cause is deterministic and corpus-wide rather than sample-specific: every criteria file in the adopted corpus is `status: draft`, and no `TR-API-*` tag resolves to an adopted namespace, so no JVM carrier anywhere in this target can currently produce a classified credit. What the sample adds is the confirmation that a future activation needs the per-executable evidence contract and cannot infer proof from class-level tags or unrelated legacy tags.

## Current validation and target test runner

Every run in this section was made against the target at the head pin. Running the current source judge without a release stamp produced `clue validate: OK (116 artifacts, 1 agent-enforced constraint(s) awaiting machine checks)` against the target and `gap` for all thirteen capabilities under `--coverage`; the development stamp intentionally skips skill-drift enforcement.

Running the same judge with `main.version=0.10.0` produced five deterministic skill-drift issues: the target's five managed `clue-*` skills are stamped `0.9.0`.

Running the target's own `.\gradlew.bat test --no-daemon` command reached the suite but failed after 23 seconds during configuration of `:bot-api:typescript:npmTest` because Gradle 9.6.1 reported four implicit task dependencies on `npmPack`. The run also warned about the JDK in use: JDK 26 is outside the JDK 17–21 build range documented in the target's architecture report, though its contributor instructions state only JDK 17+.

The failure is a task-configuration error, not a JDK rejection, so the toolchain gap and the task-dependency defect are separate concerns and only the second one stopped this run.

The target test result is therefore a prepared-environment failure, not evidence that the product suite is red on a supported toolchain; a resolution for the Gradle task dependency, and a JDK the target's own documents agree on, must be supplied before migration readiness can claim a passing target test runner.

## CI and installation boundary

Everything in this section was read at the head pin. The target's `.github/workflows/clue.yml` is a copied all-in-one validation wall pinned to `CLUE_VERSION: 0.9.0`, runs on `ubuntu-latest`, checks out with `fetch-depth: 0`, downloads the release binary and `SHA256SUMS`, installs the binary at `/usr/local/bin/clue`, detects the changed-file scope, runs `clue validate --forbid-changes`, and enforces a completed acceptance brief on Cliewen pull requests.

The target has 23 workflow action references and none are pinned to a 40-character commit SHA; the observed runner is GitHub-hosted Ubuntu and the binary installation assumes a writable `/usr/local/bin`, while no self-hosted or no-root runner constraint is evidenced in this checkout.

The target's local `.claude/skills` is a symbolic link to `..\.agents\skills`; the migration tool correctly reports seven mirror notices and does not write through that boundary.

The copied workflow is not a thin caller: current `clue migrate` preview reports the blocking finding `caller contains steps; copied validation logic is semantic and must be replaced by hand`. The tool previews six generated `.agents/skills` replacements and writes no file, but it cannot safely replace this wall without an adopter-owned manual conversion.

Six replacements and five drifting skills are two different counts over the same managed set, and neither equals its size. The release manifest owns seven generated carriers at `0.9.0`: five stamped `clue-*/skill.md` files plus the unstamped `clue-extract/mappings/madr.md` and `clue-extract/mappings/openspec.md`. Only the five stamped files can drift by version, which is what skill-drift reports. Six of the seven need replacing because `madr.md` is byte-identical between the target and the current release, so the migration has nothing to write for it.

The target's recent history contains merge commits for CH-001 through CH-006, including the adoption merge, so local history is compatible with Cliewen's supported merge-commit provenance shape. Hosted branch protection and its exact merge-mode enforcement were not read from local evidence and remain unverified here.

## Instruction and carrier reconciliation

Every carrier in this section was read at the head pin. `AGENTS.md` routes to the current corpus and says test-purpose tagging is a future adoption door, which agrees with the draft criteria state.

The target's generated `clue-extract` skill is stamped `0.9.0` and lacks the current report-only rehearsal section, segmented and suffixed ID grammar, and per-executable JVM evidence contract; its old text still delegates JVM purpose enforcement to an ArchUnit or equivalent rule.

No ArchUnit dependency or purpose rule was found in the target tree, so that delegated carrier is an unperformed obligation and must not be presented as successful extraction evidence.

No instruction-versus-build conflict on language or toolchain versions was found. The target's instruction carriers state no Kotlin version at all — the only language-version statement in them is `Java 11+` in `.agents/instructions/coding-conventions.md` — while the build source declares Kotlin 2.4.10 in `gradle/libs.versions.toml` against a JVM 11 target, so the two do not disagree.

The `Java 11+` line is worth separating from the JDK spread recorded above, because it sits in an instruction carrier rather than in contributor documentation. It states the language level the Bot API compiles against, which `DEVELOPMENT.md` also names as JDK 11 for that toolchain and which matches the build's JVM 11 target; the JDK 17+ and JDK 17–21 statements state the JDK the whole repository builds under. The two are different claims, so the target's documents spread on the build JDK without any instruction carrier disagreeing with the build.

The absence of a stated Kotlin version is itself worth naming: an instruction carrier that omits the toolchain version cannot go stale against the build, but it also gives an agent nothing to reconcile, so version coherence here is untested rather than proven.

## Governance effects

Adoption moves the accepted system of record into the Cliewen corpus, so the target's decisions, capabilities, and criteria become the artifacts a reviewer reads and approves rather than a parallel `openspec/` tree; that shift is the governance change this migration introduces, and it is already in force at the head pin.

Approval passes through the human merge boundary under PDR-004, and PDR-021 requires a merge commit for a full change: the target's local history for CH-001 through CH-006 shows merge commits, so its observed practice matches the supported shape.

The wall the target runs today also gates Cliewen pull requests on a completed acceptance brief, which makes the brief a required review step rather than a convention.

These governance effects are locally observed or proposed only. Hosted branch protection, the exact required check, and enforced merge mode were not read from local evidence, so the claim that the target's hosted governance matches its local practice is unverified here and is not evidence of readiness.

## Confidence and reversal cost

The OpenSpec-to-capability and MADR-to-decision mappings are high confidence and low reversal cost: both mappings ship in the current extraction skill, both were already performed between the two pins, and the target's own extraction report agrees with what the source tree contains.

Deterministically minted criterion IDs are high confidence at mint time but medium reversal cost once carriers reference them: the source declares no identities, so a later re-mint against a changed source renumbers criteria that tests, coverage reports, and links already point at.

Deleting the `openspec/` tree is low reversal cost while it stays a Git-history operation and no consumer outside the repository resolves those paths; it becomes high reversal cost the moment an external reference or published link depends on the deleted layout, which this rehearsal did not resolve over the network and therefore cannot rule out.

Replacing the copied CI wall with the upstream thin caller is high reversal cost and low confidence from here: the change is semantic, adopter-owned, and cannot be previewed mechanically, which is why `clue migrate` refuses it and OQ-001 holds it open.

Upgrading the six generated carriers the migration would replace is low reversal cost and high confidence: the transformation is deterministic, previewable, idempotent, and refuses locally modified inputs.

This report itself has no reversal cost. It writes nothing to the target and nothing durable here, which is what lets every judgment above stay revisable until a human authorizes mutation.

## Planned mappings and deletions

The OpenSpec mapping is available in the current Cliewen extraction skill and the MADR mapping is available as well, so no new source-format mapping is required for this target.

The twelve OpenSpec specs map to the twelve existing capabilities, the pending OpenSpec change maps to `CAP-013` and `P-002`, and the 41 numbered decisions map to the current numbered decision artifacts; these mappings are observed in the target's own extraction report at `docs/analysis/AN-001-openspec-extraction.md` in the Tank Royale tree — not this repository's `AN-001` — and are consistent with the source pin.

The planned source deletion is the complete `openspec/` tree, including its archive and pending-change records, plus the obsolete OpenSpec instruction and approval carriers; the head pin confirms that those parallel source files are already gone and that their history is retained in Git.

No target mutation is authorized by this rehearsal, so these mappings and deletions remain proposed historical evidence rather than a new extraction.

## Findings and honest support boundary

The report-only path is truthful for source inventory, deterministic ID minting where the source has no IDs, current-corpus validation, symlink-safe carrier detection, and local merge-history observation.

The target is not yet operationally ready for an end-to-end migration claim from this prepared run: the copied CI wall needs an adopter-owned manual replacement, managed skills are one release behind, and the full target test command fails on a task-configuration defect under a toolchain the target's own documents describe inconsistently.

One M-042 exit conjunct is explicitly not met and is recorded here rather than left silent.

**The proposed upstream thin-caller wall was not run in the target's CI shape.** The target still runs the copied wall, replacing it is a semantic adopter-owned change this report-only milestone may not make, and a caller cannot be exercised in the target's CI without writing to the target. What the rehearsal can say is that the target's observed CI shape — GitHub-hosted Ubuntu, a writable `/usr/local/bin`, no self-hosted or no-root constraint in evidence — is inside the shape the upstream unit already supports; whether the caller actually runs there stays unproven until OQ-001 is answered and the adopter makes that change.

Three further limits belong to this target rather than to M-042's exit list, and are named so the milestone's met conjuncts are not read as stronger than they are. M-042 asks the report to inventory these and to record the target's own results separately; it does not ask this rehearsal to exercise them.

**The target's own test result is a configuration failure, not a suite verdict.** The separated-results conjunct is met — the target's runner result and Cliewen's deterministic results are recorded apart — but the target half is a `:bot-api:typescript:npmTest` task-configuration defect under a toolchain the target's documents describe inconsistently, so nothing here says the product suite is green or red. OQ-002 holds that open.

**Stable-ID preservation was inventoried but not exercised.** The source pin declared no criterion identities at all, so ADR-037's multi-segment and letter-suffixed grammar surviving a migration is untested here. A target whose source carries stable IDs is still needed before Cliewen can claim preservation works on real brownfield input.

**The pre-mutation checkpoint itself was not exercised.** As the pins section records, this rehearsal reconstructs an extraction that already happened rather than previewing a pending one, so PDR-020's ordering guarantee — report, then human authorization, then mutation — is described from history and not demonstrated. M-042 requires that this rehearsal be produced with no target mutation, which held; it does not require that the ordering guarantee be demonstrated on an unconverted repository.

This rehearsal therefore supports a narrower boundary: Cliewen can describe and validate an already-adopted corpus without mutating the target, but it cannot claim that this target's current CI, supported build environment, or managed-carrier upgrade path is ready, that stable source identities survive migration, or that the pre-mutation checkpoint holds on an unconverted repository.

M-042's own exit criterion admits this ending: the phase closes when the rehearsal demonstrates the migration is truthful and operationally viable **or** records a narrower honest support boundary. Taking that second branch is compliance with the milestone as written, not a plan adjustment, so P-009's mutation rules permit the digest to set M-042's status and evidence without a plan revision. What still needs a human is whether this boundary is the honest one and whether the evidence cell states it plainly enough to close the migration-readiness phase; OQ-004 holds that narrower question open.

The result does not import foreign green checks, resolve links over the network, alter test purpose tags, infer unsupported stable IDs, or treat the target's current draft coverage gaps as evidence.

## Rejected alternatives

Running `clue init` or `clue migrate --apply` against Tank Royale was rejected because this milestone is report-only and the copied caller already produces a blocking semantic finding.

Treating the target's development-stamped validation as a release-ready result was rejected because it deliberately skips skill-drift checks.

Treating the failed Gradle run as a product-test verdict was rejected because the failure is a task-configuration and unsupported-JDK boundary, not a completed test assertion result.

Treating the current target's legacy `TR-API-*` tags as current Cliewen criterion evidence was rejected because the tags do not resolve to the adopted namespaces and do not carry the current proof-type and direction contract.

## Named follow-up doors

The target still needs its own capability-by-capability test-purpose migration door, beginning with draft criteria and a supported per-executable JVM carrier.

The target needs a manual adopter change that replaces the copied validation wall with the upstream thin-caller contract before a mechanical Cliewen upgrade can update its reference.

This report and both blocking questions live only in a workspace the digest deletes, and PDR-020 lands a durable report in `/docs/analysis` only in the mutate phase this change does not authorize. A narrower support boundary recorded solely in a deleted file is not recorded, so the digest must carry the unmet conjunct, the three named limits, and the open questions into durable form before the workspace goes; OQ-003 holds the shape of that landing open.

P-009 M-043 and M-044 remain independent Cliewen follow-up doors for durable dependent-change authorization and named but locally unproven foreign evidence; this rehearsal does not widen either contract.
