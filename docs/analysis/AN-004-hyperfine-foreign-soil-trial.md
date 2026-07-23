---
id: AN-004
type: analysis
status: active
provenance: inferred
links: [P-002, M-007, PDR-005]
title: Foreign-soil trial — hyperfine
---

# AN-004 — Foreign-soil trial: hyperfine

## Classification

This is the first qualifying foreign-soil trial for P-002/M-007. [`sharkdp/hyperfine`](https://github.com/sharkdp/hyperfine) is an external public open-source repository with no shared maintainer, was not designed for Cliewen, and was explicitly selected by the human. The trial was read-and-apply analysis, not adoption: no tracked file or hosted state in the target repository was changed.

The spike is complete and the finding stands, so the analysis status is `active`; its interpretation remains `provenance: inferred` because no hyperfine maintainer verified the reconstructed intent.

## Risk named

Cliewen may reconstruct a useful-looking goal, plan, and verification model from unfamiliar public evidence while silently promoting agent interpretation or environment-specific results to durable project truth.

## Evidence boundary

The trial used a detached disposable clone of `sharkdp/hyperfine` at [`f12f3d9f86f3643b3b7deace5e160b1f0f44d2b7`](https://github.com/sharkdp/hyperfine/commit/f12f3d9f86f3643b3b7deace5e160b1f0f44d2b7), the tip of `master` when CH-025 was proposed. Evidence came from committed files and history plus public repository metadata, issues, pull requests, and maintainer comments observed on 2026-07-18.

Verification ran on Windows 11 build 26200 with Rust and Cargo 1.93.1. Shell identity is recorded per result because it materially changed the outcome. Cargo populated the ignored `target/` directory and the local Cargo cache; the detached target checkout remained free of tracked changes.

## What was tried

- Reconstructed a candidate goal, capability model, acceptance boundary, quality boundary, and future direction from the README, Cargo manifest, source modules, tests, changelog, CI, history, and public GitHub activity.
- Applied `clue-analysis` and `clue-plan` without materializing a Cliewen corpus in the target.
- Assessed how `clue-delta` governance and `clue-verify` evidence would interact with the repository's existing pull requests, changelog, locked dependencies, tests, and cross-platform CI.
- Ran the repository's native format, test, lint, and release-build workflows from the clean detached checkout.
- Rejected `clue init`, a target branch, and a target pull request because PDR-005 defines this as findings work rather than adoption.

## Evidence-qualified interpretation

| Claim | Confidence | Basis |
|---|---|---|
| Hyperfine is a cross-platform command-line benchmarking tool that performs repeated runs, statistical analysis, comparison, and export | observed | `README.md`, `Cargo.toml`, CLI source, exporter modules |
| Its durable product goal is to help users make credible performance comparisons between commands | inferred | The repository says what the tool does; "credible" interprets its measurement controls, calibration, and warnings |
| Command definition, execution scheduling, platform timing, statistical comparison, presentation, and export are useful capability boundaries | inferred | Source module boundaries and focused unit/integration tests |
| Hyperfine 2.0 is a draft future direction centered on additional metrics, interleaved execution, and saved comparisons | observed as draft | [Owner-authored issue #788](https://github.com/sharkdp/hyperfine/issues/788) calls itself an unfinished planning document |
| Other open issues and pull requests are committed roadmap work | unverified | Activity shows demand and proposals, not maintainer acceptance |
| Git for Windows Bash is an intended contributor prerequisite | unverified | Windows CI exercises the suite with Bash semantics, but no contributor document declares that prerequisite |

The inferred goal is: let users compare command performance credibly across repeated measurements while exposing enough controls and raw evidence to diagnose environmental distortion.

Likely capabilities are command definition and parameter expansion; setup, warmup, preparation, benchmark, conclusion, and cleanup scheduling; platform-specific process timing; statistical summarization and comparison; terminal presentation and warnings; and structured result export. These boundaries are navigationally useful but are not presented as maintainer-approved taxonomy.

Issue #788 can map to a `draft` Cliewen plan because the owner explicitly names it as unfinished planning and describes three intended features plus two possible behavior changes. The rest of the issue and pull-request queue cannot safely become milestones merely because it is active.

## Native verification

| Invocation | Environment | Result |
|---|---|---|
| `cargo fmt -- --check` | PowerShell, Windows Rust 1.93.1 | passed |
| `cargo test --locked` | PowerShell, Windows Rust 1.93.1 | failed: 95 passed and two integration tests failed because `echo` was a PowerShell alias rather than an executable discoverable by the program under test |
| `cargo test --locked` | Git for Windows Bash invoked by exact executable path, Windows Rust 1.93.1 | passed: 97 tests |
| `cargo clippy --locked --all-targets` | Git for Windows Bash, Windows Rust 1.93.1 | passed with existing warnings, including the deprecated `assert_cmd::cargo_bin` use tracked by [issue #844](https://github.com/sharkdp/hyperfine/issues/844) |
| `cargo build --locked --release` | PowerShell, Windows Rust 1.93.1 | passed |

An unqualified `bash` command resolved to Windows Subsystem for Linux and exercised a Linux build. That result was discarded. Repeating the command with the exact Git for Windows Bash executable reproduced the relevant Windows environment and passed.

## Findings

### F1 — Confidence must survive the analysis handoff

The repository offers authoritative product description, executable behavior, owner-authored draft planning, issue discussion, pull-request activity, and agent interpretation. Those sources do not have equal authority. `clue-plan` correctly prevents an agent from activating semantic plan changes, but `clue-analysis` did not require confidence labels, so an unqualified analysis could hand a plausible inference to planning as if it were maintainer intent.

### F2 — A command string is incomplete verification evidence

The same checkout, revision, toolchain, operating system, and command failed under PowerShell and passed under Git for Windows Bash. The difference was executable discovery through the shell environment, not source behavior. A clean-clone rule alone would not preserve enough evidence to reproduce or interpret the result; relevant runtime, toolchain, operating-system, shell, and invocation deviations belong in the evidence boundary.

### F3 — Statistical product quality confirms the AC/quality split

Deterministic contracts fit acceptance criteria well: option validation, execution ordering, minimum-run and minimum-time rules, failure handling, export schemas, relative comparison, quoting, and platform command behavior have positive and negative test surfaces.

Measurement credibility is different. Timing distributions depend on caches, scheduling, command duration, shell calibration, and unrelated system activity. Hyperfine can deterministically calculate specified statistics and warnings, while usefulness, repeatability, and bias require environment-qualified quality scenarios and evidence. This independently confirms AN-003/F3 in a statistical benchmarking domain, but the split remains a candidate until the second qualifying trial tests whether standing skill guidance would help beyond compatibility and benchmarking projects.

### F4 — Repository activity is not a plan

The owner-authored Hyperfine 2.0 issue is explicitly a draft plan. Other active issues and pull requests are evidence of demand and proposed work, not of accepted direction. Treating recency or volume as roadmap authority would misrepresent the project.

### F5 — Adoption would change governance

Hyperfine already has pull requests, a user-facing changelog, a locked dependency graph, formatting, tests, linting, release builds, and a broad CI matrix. Cliewen would still add mandatory plan linkage or plan-less declarations, typed decisions, transient workspaces for non-light changes, artifact validation, and a human-only merge boundary. The trial can assess those consequences without describing scaffolding as neutral or asking the target to adopt them.

## Methodology adjustment

CH-025 updates `clue-analysis` to require an evidence boundary: pin source revisions when possible, record conditions relevant to reproduced results, distinguish observations from inferences and unverified intent, and avoid treating repository activity as maintainer intent without explicit evidence. The adjustment is carried by the generated analysis skill and a cheap, reversible decision-log row.

This adjustment is deliberately narrower than all candidate changes from AN-003. It is directly reproduced by a qualifying external trial and closes one real ambiguity without inventing a universal brownfield protocol before the second trial.

## Rejected interpretations

- **The README is the accepted goal.** It is an authoritative product description, but the credibility language in the candidate goal is an inference.
- **Every active issue or pull request is plan evidence.** Issue #788 is explicit draft planning; other activity remains demand or proposed work until accepted.
- **PowerShell proves the test suite is broken.** The repository's Windows CI uses Bash semantics, and the exact Git for Windows Bash environment passes.
- **A passing Bash run makes the PowerShell failure irrelevant.** The failure exposes an undeclared environmental dependency and proves why reproduction conditions matter.
- **The unqualified `bash` rerun is Windows evidence.** It resolved to WSL and tested Linux, so it was discarded.
- **A real trial requires initializing Cliewen in the target.** PDR-005 defines qualifying trials as read-and-apply findings, not adoptions.

## Trial result

Cliewen recovered a useful goal, capability model, deterministic contract surface, quality boundary, and explicit draft plan without modifying the target. The exercise also exposed a concrete evidence-provenance defect in `clue-analysis` and produced the first M-007-traceable methodology adjustment. M-007 remains open until the independently selected `toss/es-toolkit` trial is completed and both findings are considered together.
