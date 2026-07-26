---
id: CH-066-tasks
type: tasks
status: open
links: [CH-066]
title: Tasks for CH-066
---

# Tasks

- [ ] Write `PDR-014`: the plan revision reopening the parked package-manager scope onto the active P-007
- [ ] Write `ADR-030`: package-manager-first distribution, release asset names as an append-only contract, and the `brews:`-deprecation fallback
- [ ] Add `.goreleaser.yaml` emitting both asset shapes — bare binaries under the current names and archives for the package managers — under one `SHA256SUMS`
- [ ] Rewrite `.github/workflows/release.yml` onto goreleaser, keeping the tag guard, the changelog extraction, the stamped drift gate, and `go test` untouched
- [ ] Add `readReleaseWorkflow` and `readGoreleaserConfig` helpers to `cmd/clue/main_test.go`
- [ ] Retarget `TestSanity_ReleaseWorkflowIsCrossPlatform` to the parsed goreleaser config, judging the build matrix, the stamp, and the checksum name as structure rather than substrings
- [ ] Retarget `TestSanity_ReleaseNotesComeFromChangelog` to `--release-notes` and assert goreleaser's own changelog generator is disabled
- [ ] Fix the fail-open ordering guard in `TestSanity_ReleaseRunsTheJudgeStampedAsTheTag`: a marker that has vanished must fail, not silently pass
- [ ] Add `TestSanity_ReleaseKeepsTheAssetNamesTheAdopterWallInstalls`, tying the goreleaser asset name to the scaffolded wall that installs it
- [ ] Dry-run goreleaser locally: settle whether `formats: [binary]` appends `.exe` on Windows, and prove `sha256sum -c --ignore-missing SHA256SUMS` still passes
- [ ] Lead `guide/getting-started.md` with the package-manager commands and demote the manual download to a fallback
- [ ] Record in `guide/operations.md` that upgrading the binary alone produces a correct drift report, not a bug
- [ ] Mirror the install restructure in `README.md`
- [ ] Rewrite CAP-004's release-pipeline design for goreleaser and add the append-only asset rule to its deliberate limits
- [ ] Add the `[Unreleased]` changelog entry
- [ ] Run repository and corpus verification
