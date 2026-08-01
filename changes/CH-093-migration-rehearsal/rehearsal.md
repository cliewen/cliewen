---
id: CH-093-rehearsal
type: analysis
status: draft
links: [CH-093, P-009, M-042, PDR-020, AN-013, AN-014, CAP-003]
title: PDR-020 rehearsal report for Robocode Tank Royale
provenance: inferred
reversal-cost: low
---

# CH-093 — PDR-020 rehearsal report: Robocode Tank Royale

## Risk and scope

The risk tested was that a real brownfield repository can appear ready after its corpus is converted while stable source meaning, executable evidence, CI operation, or deletion and merge boundaries remain unproven.

This is a report-only rehearsal for P-009/M-042. It changes no file in the Tank Royale checkout, no durable Cliewen corpus, no target tests or routing, and no hosted state.

## Pins and evidence boundary

The selected repository is the public Robocode Tank Royale checkout at `C:\Code\tank-royale` with clean `main` at `384d27d55176a2d2ad4668ac381852e629e4540a`, the same public adopter revision pinned by AN-013.

The current source tree has already been adopted, so the source snapshot for this rehearsal is `4e579878fd6667fab3b75515b6e68135a935c8df`, the first parent of adoption merge `b7fb320ccec1f4742ef923cb315e7dd84f7e824f`; this is the last source state containing `openspec/` before CH-001 removed it.

Read-only inspection ran on Microsoft Windows NT 10.0.26200.0 with PowerShell 7.6.4, Go 1.26.5, Java 26, Node.js 24.0.0, Python 3.14.2, and .NET 8.0.423.

These are prepared-environment results, not clean-onboarding evidence: the existing checkout, local Gradle caches and build outputs, existing Node dependencies, and the local Go source tree were prerequisites; the target's own instructions require JDK 17–21 while this machine has JDK 26.

The target was clean before inspection and clean after inspection, including after the failed Gradle test command; no tracked or staged target file changed.

## Source inventory at the pinned pre-adoption snapshot

The source corpus had `openspec/project.md`, `openspec/config.yaml`, twelve `openspec/specs/*/spec.md` capability specs, one pending `openspec/changes/add-typescript-bot-api-npm-publish`, and an `openspec/changes/archive` tree.

The twelve specs contain 107 `Requirement` headings and 257 `Scenario` headings; the scenarios carry no stable criterion IDs, so an extraction cannot preserve source scenario identities and must mint them deterministically in capability namespaces.

The source had 41 numbered MADR decisions (`0001` through `0041`) under `docs/decisions/`.

The source carried `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, `.junie/guidelines.md`, twelve `.agents/instructions/*` files, and `.principles` files at the repository and module levels; routing reconciliation therefore covered more than the flagship assistant file.

The source language inventory was 132 Java files, 356 Kotlin files, 133 Python files, 113 TypeScript files, and 118 C# files; 150 of those files were test-like by path or name (33 Java, 36 Kotlin, 35 Python, 17 TypeScript, and 29 C#).

The source had no per-test Cliewen purpose tags. Existing test metadata was a separate test-registry convention and could not be treated as Cliewen criterion evidence without a mapping and normalization pass.

## Current corpus and ID mapping

The current adopted corpus has thirteen capabilities with prefixes `BR`, `BFD`, `BAU`, `GBP`, `GBC`, `PRO`, `PBA`, `RCV`, `TCS`, `TFS`, `TBA`, `UD`, and `TNP`.

The current criteria files contain 263 scenarios, all with `status: draft`; the first twelve capabilities carry the 257 source scenarios and `TNP` carries six scenarios derived from the pending TypeScript publication change.

The current corpus therefore demonstrates deterministic minting from a source with no IDs, but this rehearsal does not exercise preservation of an existing multi-segment or letter-suffixed source ID because the pinned source did not declare one.

All current criteria remain draft, so current `clue validate --coverage` reports `gap` for every capability and there is no active classified AC coverage for the rehearsal to mistake for proof.

## JVM evidence sample

The manual sample included `bot-api/java/src/test/java/dev/robocode/tankroyale/botapi/ProtocolConformanceTest.java`, which has a class-level `TCK` tag and eleven method-level `TR-API-TCK-007` through `TR-API-TCK-017` tags across eleven executable tests.

The sample also included `bot-api/java/src/test/java/dev/robocode/tankroyale/botapi/BotRunFirstTurnTest.java`, which repeats `TCK` and `TR-API-TCK-004` at class and method scope across two executable tests.

These files contain mixed executable tags but no Cliewen criterion identity, declared `Test-type`, or positive/negative direction on the same executable; no `// AC:` proximity carriers were found in the sampled Java and Kotlin tree.

The current validator credits neither file as classified AC evidence because the tags do not resolve to the current criteria namespaces and every criteria file is draft. The sample found no false classified coverage, but it does confirm that a future activation needs the per-executable evidence contract and cannot infer proof from class-level tags or unrelated legacy tags.

## Current validation and target test runner

Running the current source judge without a release stamp produced `clue validate: OK (116 artifacts, 1 agent-enforced constraint(s) awaiting machine checks)` against the target and `gap` for all thirteen capabilities under `--coverage`; the development stamp intentionally skips skill-drift enforcement.

Running the same judge with `main.version=0.10.0` produced five deterministic skill-drift issues: the target's five managed `clue-*` skills are stamped `0.9.0`.

Running the target's own `.\gradlew.bat test --no-daemon` command reached the suite but failed after 23 seconds during configuration of `:bot-api:typescript:npmTest` because Gradle 9.6.1 reported four implicit task dependencies on `npmPack`; the build also warned that JDK 26 is outside the documented JDK 17–21 range.

The target test result is therefore a prepared-environment failure, not evidence that the product suite is red on a supported toolchain; a supported JDK and a resolution for the Gradle task dependency must be supplied before migration readiness can claim a passing target test runner.

## CI and installation boundary

The target's `.github/workflows/clue.yml` is a copied all-in-one validation wall pinned to `CLUE_VERSION: 0.9.0`, runs on `ubuntu-latest`, checks out with `fetch-depth: 0`, downloads the release binary and `SHA256SUMS`, installs the binary at `/usr/local/bin/clue`, detects the changed-file scope, runs `clue validate --forbid-changes`, and enforces a completed acceptance brief on Cliewen pull requests.

The target has 23 workflow action references and none are pinned to a 40-character commit SHA; the observed runner is GitHub-hosted Ubuntu and the binary installation assumes a writable `/usr/local/bin`, while no self-hosted or no-root runner constraint is evidenced in this checkout.

The target's local `.claude/skills` is a symbolic link to `..\.agents\skills`; the migration tool correctly reports seven mirror notices and does not write through that boundary.

The copied workflow is not a thin caller: current `clue migrate` preview reports the blocking finding `caller contains steps; copied validation logic is semantic and must be replaced by hand`. The tool previews six generated `.agents/skills` replacements and writes no file, but it cannot safely replace this wall without an adopter-owned manual conversion.

The target's recent history contains merge commits for CH-001 through CH-006, including the adoption merge, so local history is compatible with Cliewen's supported merge-commit provenance shape. Hosted branch protection and its exact merge-mode enforcement were not read from local evidence and remain unverified here.

## Instruction and carrier reconciliation

`AGENTS.md` routes to the current corpus and says test-purpose tagging is a future adoption door, which agrees with the draft criteria state.

The target's generated `clue-extract` skill is stamped `0.9.0` and lacks the current report-only rehearsal section, segmented and suffixed ID grammar, and per-executable JVM evidence contract; its old text still delegates JVM purpose enforcement to an ArchUnit or equivalent rule.

No ArchUnit dependency or purpose rule was found in the target tree, so that delegated carrier is an unperformed obligation and must not be presented as successful extraction evidence.

The target's `AGENTS.md` describes Kotlin 2.2.0 while the current `gradle/libs.versions.toml` declares Kotlin 2.4.10; this is a target instruction-versus-build-source conflict that needs maintainer resolution before a migration report can call the repository's operating instructions coherent.

## Planned mappings and deletions

The OpenSpec mapping is available in the current Cliewen extraction skill and the MADR mapping is available as well, so no new source-format mapping is required for this target.

The twelve OpenSpec specs map to the twelve existing capabilities, the pending OpenSpec change maps to `CAP-013` and `P-002`, and the 41 numbered decisions map to the current numbered decision artifacts; these mappings are observed in the existing `docs/analysis/AN-001-openspec-extraction.md` and are consistent with the pinned source tree.

The planned source deletion is the complete `openspec/` tree, including its archive and pending-change records, plus the obsolete OpenSpec instruction and approval carriers; the current head confirms that those parallel source files are already gone and that their history is retained in Git.

No target mutation is authorized by this rehearsal, so these mappings and deletions remain proposed historical evidence rather than a new extraction.

## Findings and honest support boundary

The report-only path is truthful for source inventory, deterministic ID minting where the source has no IDs, current-corpus validation, symlink-safe carrier detection, and local merge-history observation.

The target is not yet operationally ready for an end-to-end migration claim from this prepared run: the copied CI wall needs an adopter-owned manual replacement, managed skills are one release behind, the full target test command fails under the available toolchain, and one repository instruction conflicts with the build version source.

This rehearsal therefore supports a narrower boundary: Cliewen can describe and validate the adopted corpus without mutating the target, but it cannot claim that this target's current CI, supported build environment, or managed-carrier upgrade path is ready until the blocking questions below are resolved.

The result does not import foreign green checks, resolve links over the network, alter test purpose tags, infer unsupported stable IDs, or treat the target's current draft coverage gaps as evidence.

## Rejected alternatives

Running `clue init` or `clue migrate --apply` against Tank Royale was rejected because this milestone is report-only and the copied caller already produces a blocking semantic finding.

Treating the target's development-stamped validation as a release-ready result was rejected because it deliberately skips skill-drift checks.

Treating the failed Gradle run as a product-test verdict was rejected because the failure is a task-configuration and unsupported-JDK boundary, not a completed test assertion result.

Treating the current target's legacy `TR-API-*` tags as current Cliewen criterion evidence was rejected because the tags do not resolve to the adopted namespaces and do not carry the current proof-type and direction contract.

## Named follow-up doors

The target still needs its own capability-by-capability test-purpose migration door, beginning with draft criteria and a supported per-executable JVM carrier.

The target needs a manual adopter change that replaces the copied validation wall with the upstream thin-caller contract before a mechanical Cliewen upgrade can update its reference.

P-009 M-043 and M-044 remain independent Cliewen follow-up doors for durable dependent-change authorization and named but locally unproven foreign evidence; this rehearsal does not widen either contract.
