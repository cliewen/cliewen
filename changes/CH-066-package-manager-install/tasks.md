---
id: CH-066-tasks
type: tasks
status: open
links: [CH-066]
title: Tasks for CH-066
---

# Tasks

- [x] Write `PDR-014`: the plan revision reopening the parked package-manager scope onto the active P-007
- [x] Write `ADR-030`: package-manager-first distribution, release asset names as an append-only contract, and the `brews:`-deprecation fallback
- [x] Add `.goreleaser.yaml` emitting both asset shapes — bare binaries under the current names and archives for the package managers — under one `SHA256SUMS`
- [x] Rewrite `.github/workflows/release.yml` onto goreleaser, keeping the tag guard, the changelog extraction, the stamped drift gate, and `go test` untouched
- [x] Add `readReleaseWorkflow` and `readGoreleaserConfig` helpers to `cmd/clue/main_test.go`
- [x] Retarget `TestSanity_ReleaseWorkflowIsCrossPlatform` to the parsed goreleaser config, judging the build matrix, the stamp, and the checksum name as structure rather than substrings
- [x] Retarget `TestSanity_ReleaseNotesComeFromChangelog` to `--release-notes` and assert goreleaser's own changelog generator is disabled
- [x] Fix the fail-open ordering guard in `TestSanity_ReleaseRunsTheJudgeStampedAsTheTag`: a marker that has vanished must fail, not silently pass
- [x] Add `TestSanity_ReleaseKeepsTheAssetNamesTheAdopterWallInstalls`, tying the goreleaser asset name to the scaffolded wall that installs it
- [x] Dry-run goreleaser locally: settle whether `formats: [binary]` appends `.exe` on Windows, and prove `sha256sum -c --ignore-missing SHA256SUMS` still passes — both confirmed; `SHA256SUMS` records upload names, so the adopter wall verifies unchanged
- [x] Prove each new guard bites by mutating the config it judges: renamed checksum, renamed bare asset, removed bare archive, re-enabled changelog generator, dropped `-trimpath` — all five refused
- [-] Ship a Homebrew tap — withdrawn: `brews:` is deprecated upstream, and its replacement `homebrew_casks:` is macOS-only, so Homebrew can no longer serve macOS and Linux from one channel
- [x] Lead `guide/getting-started.md` with the package-manager commands and demote the manual download to a fallback
- [x] Record in `guide/operations.md` that upgrading the binary alone produces a correct drift report, not a bug
- [x] Mirror the install restructure in `README.md`
- [x] Rewrite CAP-004's release-pipeline design for goreleaser and add the append-only asset rule to its deliberate limits
- [x] Add the `[Unreleased]` changelog entry

## Second pass — the install channel changed

- [x] Remove `brews:` from `.goreleaser.yaml` and restore the goreleaser version pin's stated reason — `goreleaser check` is now clean, with no deprecated properties
- [x] Write `guide/public/install.sh`: detect OS and architecture, resolve the release, download the bare asset and `SHA256SUMS`, verify the checksum before installing, and install to a `PATH` directory
- [x] Write `guide/public/install.ps1` with the same contract for Windows
- [x] Prove both scripts against the live v0.7.0 release, including the failure path: a tampered binary is refused and nothing is installed
- [x] Rewrite `ADR-030` for the actual decision: a checksum-verifying install script is the primary channel; native package managers are additive and follow
- [x] Revise CAP-001's design lesson that says the first encounter does not present installer scripts, stating what changed rather than contradicting it silently
- [x] Revise M-030's exit criterion in `P-007` so it describes the channel actually being built
- [x] Rewrite the install sections of `guide/getting-started.md`, `README.md`, `guide/operations.md`, CAP-004's design, and the `[Unreleased]` changelog entry so every documented command works on the day it ships
- [x] Guard the scripts against the asset contract they depend on, so a renamed asset fails a test rather than a user's install; proven to bite by removing the verification and by requiring elevation
- [x] Run repository and corpus verification
