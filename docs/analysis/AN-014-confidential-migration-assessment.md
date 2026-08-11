---
id: AN-014
type: analysis
status: active
provenance: verified
reversal-cost: low
links: [P-009, CAP-003, AN-012, AN-013, ADR-005, ADR-009, PDR-020]
title: A confidential migration assessment exposes evidence, identity, and enterprise-operation gaps
---

# AN-014 — A confidential migration assessment exposes evidence, identity, and enterprise-operation gaps

## Risk

Cliewen may claim migration readiness while its JVM evidence can produce false classified coverage, its criterion grammar rejects existing stable identities, or its shipped CI and upgrade paths require enterprise adopters to maintain private forks.

## Evidence boundary

This finding records a migration assessment supplied to the project on 2026-07-30 for two confidential workplace repositories. Their identities are deliberately omitted. The supplied material is a coding agent's prose verdict, not source files, command output, a pinned source revision, or a Cliewen rehearsal. It reports that the agent inspected Cliewen 0.10.0's JVM harvester but did not run `clue` against either target. No target result is independently reproduced here, so target counts, policy claims, and migration viability are reported observations with unknown reproducibility, not deterministic Cliewen evidence.

The Cliewen-side mechanics were checked in the prepared local repository at accepted `main` commit `beca8904f6f9785a5dfd4f2cf9add81118792cee` on Windows with PowerShell and the repository's Go toolchain. This was not a clean disposable environment and establishes only the source-level behavior described below. No statistical population, sample, or percentage is claimed.

## Supplied target observations

The assessment reports that one target is structurally close to migration while the other is blocked by existing corpus debt and unsupported stable IDs. Across the targets it reports JVM files containing multiple criterion tags and separate positive and negative test methods, criterion identities with multi-segment prefixes such as `SNAP-SQS-001`, letter-suffixed forms such as `ADP-045b`, non-JUnit JVM tests whose existing carrier is a source comment, internal runner scale sets, action references pinned by SHA, and repository-specific binary installation. These claims describe the supplied targets only. Their frequency and exact counts are intentionally not used as plan acceptance evidence because the repositories and revisions cannot be inspected from this corpus.

The assessment also reports two workflow costs: extraction ultimately deletes the parallel source corpus after its human-authorized rehearsal, and Cliewen permits one in-flight change per initiating author rather than stacked dependent changes. Those are current deliberate methodology boundaries, not newly discovered implementation defects. The supplied material does not reproduce a failure that would justify changing either boundary.

## Reproduced Cliewen mechanics

`internal/corpus/actests.go` harvests every JUnit `@Tag` in a JVM test file, collects all proof types and directions in that file, and records their cross-product for every recognized AC tag. It does not associate an annotation block with a Java or Kotlin executable. A file containing several AC tags and separate positive and negative methods can therefore credit every recognized criterion with both directions even when no individual method carries that combination. This confirms the false-green mechanism, not that either confidential target currently receives a false-green verdict.

The same implementation accepts criterion declarations only through `@([A-Z][A-Z0-9]*)-(\d+)`, and ADR-009 defines the public grammar as `<PREFIX>-<digits>`. Multi-segment prefixes and letter suffixes are rejected or only partially recognized. Preserving the reported examples verbatim therefore requires a grammar decision and implementation change; current extraction cannot promise it.

AN-012 already demonstrates that editing the copied CI wall severs its upstream update path and that a corpus obligation change can create a large mechanical upgrade. Its one measured adopter did not need runner, action-reference, or install-path configuration and explicitly supplied no evidence for a general configuration file. The confidential assessment reports different enterprise policies, but without inspectable workflows. P-009 must therefore test those policy shapes with reviewable synthetic fixtures before binding a public input contract; the new report does not retroactively change AN-012's finding.

The operations guide does not select or distinguish merge, squash, and rebase strategies. The architecture calls the merge commit the acceptance and `git log docs/` the archive, while a squash of a branch that both creates and deletes `/changes/CH-xxx/` retains only the net tree. This confirms an unstated support boundary. It does not decide whether Cliewen must require merge commits or add a squash-compatible durable carrier.

## Findings

1. **JVM classified evidence is a correctness gap.** File-level attribution can over-credit evidence and must fail conservatively at executable scope before migration readiness is claimed.
2. **The criterion-ID grammar is a reported migration blocker with a reproduced implementation boundary.** The target prevalence is unverified, but supporting the two concrete shapes can be specified and tested without renumbering existing corpora.
3. **A non-forking CI wall and mechanical upgrades remain demonstrated adopter needs.** Enterprise runner, pinning, and installation inputs are hypotheses from a confidential report until synthetic fixtures reproduce the policy constraints.
4. **Merge strategy is an undocumented provenance boundary.** The supported history shape must be decided and made testable.
5. **Non-JUnit JVM evidence needs a stable executable carrier, not proximity comments.** ADR-005's rejection still applies; ADR-032's named-test fallback is the compatible route.

## Rejected

- **Treating the supplied verdict as a successful rehearsal.** Its author explicitly did not run Cliewen against the targets, and no pinned target revision or output is available.
- **Renumbering unsupported source IDs during extraction.** That breaks ADR-009's stable-identity rationale to avoid extending a parser.
- **Adding a general configuration file from this report.** AN-012 rejected that interface on measured evidence, and this report does not establish a second source of truth for all Cliewen settings.
- **Crediting `// AC:` comments.** Proximity cannot prove which executable a comment classifies; a stable method name can.
- **Making permanent dual-running or stacked changes migration prerequisites.** The supplied assessment names costs but does not reproduce a correctness failure in the current boundaries.

## Consumer

[P-009](../plans/P-009-migration-readiness.md) consumes these findings. M-037 removes the reproduced false-green mechanism and supplies the non-JUnit named-test route; M-038 extends stable identities; M-039 decides merge-mode support; M-040 tests the reported enterprise CI constraints before choosing inputs; M-041 builds the already-routed upgrade path; and M-042 requires the report-only rehearsal this evidence lacks.
