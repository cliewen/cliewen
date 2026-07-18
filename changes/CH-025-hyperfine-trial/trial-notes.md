---
id: CH-025-notes
type: notes
status: active
links: [CH-025]
title: Hyperfine trial working notes
---

# Hyperfine trial working notes

## Risk

Can Cliewen reconstruct useful and appropriately qualified project intent, product contracts, and verification boundaries from an unfamiliar public repository without turning agent interpretation into maintainer intent or imposing Cliewen on the target?

## Evidence boundary

- Target: `sharkdp/hyperfine` at `f12f3d9f86f3643b3b7deace5e160b1f0f44d2b7`, the tip of `master` when CH-025 was proposed.
- Checkout: detached, disposable clone under the local temporary directory.
- Sources: committed repository files and history, public GitHub repository metadata, issues, pull requests, and owner comments.
- Mutations: no tracked target file, commit, branch, issue, pull request, or other hosted state was changed. Cargo populated the ignored `target/` directory and the local cargo cache.
- Verification host: Windows 11 build 26200, Rust and Cargo 1.93.1. Shell identity is part of each result because it changed the outcome.

## Evidence register

| Claim | Confidence | Evidence |
|---|---|---|
| Hyperfine is a cross-platform command-line benchmarking tool that performs repeated runs, statistical analysis, comparison, and export | observed | `README.md`, `Cargo.toml`, CLI source, exporter modules |
| Its durable product goal is to help users make credible performance comparisons between commands | inferred | The repository says what the tool does, but it has no explicit goal artifact and the credibility wording is an interpretation of its measurement controls and warnings |
| Command execution lifecycle, CLI parsing, result calculation, formatting, and export are separate capabilities | inferred | Source module boundaries and focused unit/integration tests |
| Hyperfine 2.0 is a draft future direction centered on additional metrics, interleaved execution, and saved comparisons | observed as draft | Owner-authored issue #788 explicitly calls itself an unfinished planning document and links a milestone |
| Any open issue or pull request outside that owner-authored draft is committed roadmap work | unverified | Activity shows demand and proposed work, not maintainer acceptance |
| Bash is an intended contributor prerequisite on Windows | unverified | CI runs test steps through Bash and the suite passes there, but no contributor documentation declares the prerequisite |
| The statistical warning policy represents the maintainer's present view that environmental outliers should remain visible rather than silently discarded | observed for the cited discussion, not a timeless decision | Owner discussion in issue #528; no durable decision record in the repository |

## Native verification

| Invocation | Environment | Result |
|---|---|---|
| `cargo fmt -- --check` | PowerShell, Windows Rust 1.93.1 | passed |
| `cargo test --locked` | PowerShell, Windows Rust 1.93.1 | failed: 95 passed and 2 integration tests failed because `echo` was a PowerShell alias rather than an executable discoverable by the program under test |
| `cargo test --locked` | Git for Windows Bash invoked by exact executable path, Windows Rust 1.93.1 | passed: 97 tests |
| `cargo clippy --locked --all-targets` | Git for Windows Bash, Windows Rust 1.93.1 | passed with existing warnings, including the publicly tracked deprecated `assert_cmd::cargo_bin` use |
| `cargo build --locked --release` | PowerShell, Windows Rust 1.93.1 | passed |

An unqualified `bash` command resolved to Windows Subsystem for Linux and exercised a Linux build rather than the intended Windows CI shell. That result was discarded. Repeating the command with `C:\Program Files\Git\bin\bash.exe` reproduced the relevant Windows execution environment and passed.

The target remained detached with no tracked changes after verification. The ignored Cargo build directory is the only repository-local output.

## Read-and-apply result

### Analysis

The repository makes most deterministic behavior discoverable through a compact source tree and focused tests. It does not state the authority level of prose, issues, owner discussion, code, and executable results, so a brownfield analyst must add those confidence distinctions or risk promoting a plausible interpretation to fact.

### Hypothetical goal and capability model

The inferred goal is to let users compare command performance credibly across repeated measurements while exposing enough controls and raw output to diagnose environmental distortion.

Likely capability boundaries are command definition and parameter expansion; setup, warmup, preparation, benchmark, conclusion, and cleanup scheduling; platform-specific process timing; statistical summarization and comparison; terminal presentation and warnings; and structured result export. These boundaries are useful navigation, but only a maintainer could accept them as the project's durable taxonomy.

### Acceptance and quality

Deterministic contracts fit Cliewen acceptance criteria well: option validation, execution ordering, minimum-run and minimum-time rules, failure handling, output schemas, relative comparison, quoting, and platform-specific command behavior all have positive and negative test surfaces.

Measurement credibility does not fit a binary acceptance criterion. Timing distributions depend on caches, scheduling, command duration, shell calibration, and unrelated system activity. The product can deterministically calculate documented statistics and warnings, while the usefulness, repeatability, and bias of measurements require environment-qualified quality scenarios and evidence. This independently confirms AN-003/F3 on a repository whose core domain is statistical measurement rather than population compatibility.

### Plan

Issue #788 can map to a `draft` Cliewen plan because the owner explicitly identifies it as unfinished planning and names three intended features plus two possible behavior changes. Open issues and pull requests cannot safely become plan milestones merely because they are active. `clue-plan` correctly requires human acceptance before a semantic plan becomes active, but the analysis handoff currently lacks a rule that preserves the confidence distinction used to reach the draft.

### Change and verification

Hyperfine already uses pull requests, a locked dependency graph, a user-facing changelog, formatting, tests, clippy, release builds, and a broad cross-platform CI matrix. Cliewen adoption would still introduce mandatory plan linkage or plan-less declarations, typed decisions, transient change workspaces for non-light changes, Cliewen artifact validation, and a human-only merge boundary. Those are governance changes, not neutral documentation scaffolding.

The clean-checkout experiment also shows that a command string is incomplete evidence. The same repository, revision, toolchain, operating system, and nominal command passed or failed depending on the invoking shell and executable search path. Verification findings need to preserve the relevant environment and identify deviations from the repository's declared or CI workflow.

## Candidate methodology adjustment

Add one evidence-boundary step to `clue-analysis`: identify the inspected revision and the relevant execution environment for reproduced results; label material claims as observed, inferred, or unverified intent; and do not treat repository activity as maintainer intent. This is small, general, and directly supported by the two distinct failures exposed here: confidence loss at the analysis-to-plan handoff and a verification result that reverses with shell identity.

The deterministic-acceptance versus statistical-quality split is confirmed as useful, but it is not added to the skill in this change. A second qualifying trial can show whether it generalizes beyond compatibility and benchmarking domains before Cliewen turns it into standing guidance.

## Rejected interpretations

- The README is the accepted goal. It is an authoritative product description, but the goal wording above adds an inference about credible comparison.
- Every open issue is plan evidence. Issue #788 is explicit draft planning; other issues and pull requests are demand, proposals, or activity until the maintainer accepts them.
- PowerShell proves the test suite is broken. The repository's CI uses Bash for the test step, and the exact Git for Windows Bash environment passes.
- A passing Bash run makes the PowerShell failure irrelevant. The failure exposes an undeclared environment dependency and demonstrates why a verification record must include more than a command string.
- The unqualified `bash` rerun was Windows CI evidence. It resolved to WSL and tested Linux, so it was discarded.
- Cliewen should be initialized in the clone to make the trial real. PDR-005 defines the trial as read-and-apply findings and prohibits adoption work against the foreign repository.
