---
id: ADR-030
type: decision
status: inferred
links: [CAP-001, CAP-004, C-015, ADR-011, ADR-012]
title: Distribution is package-manager-first, and release asset names are an append-only contract
author: agent
accepted-by: []
---

# ADR-030 — Package-manager distribution and append-only release assets

## Context and problem statement

Getting `clue` onto a machine takes six manual steps: read a platform table, download the right asset, verify its SHA-256 by hand, rename it, make it executable, and put it on `PATH`. [C-015](../constraints/C-015-onboarding-under-30-minutes.md) gives the whole journey — install, `clue init`, `clue validate` — thirty minutes, and those six steps are what a newcomer meets first.

Package managers remove all six, but adopting them touches something delicate. The release currently publishes **bare binaries**, and `internal/scaffold/templates/github/workflows/clue.yml` — the CI wall vendored into every repository that has run `clue init` — verifies `SHA256SUMS` and then runs `install -m 0755 "clue-${CLUE_VERSION}-linux-amd64"`. Homebrew and winget both want archives. The naive migration renames the published assets, and every adopted repository's wall breaks at its next CI run, with no change in that repository to explain it.

## Decision outcome

**Package managers are the documented first install path, the release is built by goreleaser, and the published asset names are an append-only public contract.**

- **Channels.** A Homebrew *formula* in `cliewen/homebrew-tap` serves macOS and Linux; winget serves Windows. The manual download stays documented as a fallback, unchanged in substance, because it is also the only route for a machine with neither package manager and the only place the asset names are written down for a human.
- **goreleaser replaces the hand-rolled build loop.** Not for its own sake: it generates the tap formula and the winget manifest from the tag, so neither becomes a second place a version is bumped by hand. The version it injects is `{{ .Version }}` — the tag with its `v` stripped — which is the bare semver [ADR-011](ADR-011-version-stamping.md) requires for comparison against the skills' frontmatter stamp.
- **Assets are append-only.** Every release publishes both shapes under one `SHA256SUMS`: the bare binaries under their exact existing names `clue-<version>-<os>-<arch>[.exe]`, and the archives the package managers consume. A new channel may add assets; it may never rename or remove one. The adopters' wall reads its checksums with `--ignore-missing`, so the additional archive lines are inert.
- **Release notes are unaffected.** goreleaser's own changelog generation is disabled and the body is passed as the extracted `CHANGELOG.md` section, preserving [ADR-012](ADR-012-release-notes-from-changelog.md) exactly. Disabling the generator matters independently of the flag: it is the commit-log-and-@mentions dump ADR-012 exists to refuse, and a config that still asks for one would produce it the moment the flag was dropped.

The asset-name rule is the load-bearing half. A release artifact whose name appears in a file we ship into other people's repositories is a published interface, not an implementation detail, and it is invisible as one — nothing about `dist/clue-0.7.0-linux-amd64` announces that a workflow somewhere executes that string. Naming it a contract is what makes the constraint reviewable instead of rediscovered.

**Carrier:** `.goreleaser.yaml` (machine); `TestSanity_ReleaseKeepsTheAssetNamesTheAdopterWallInstalls`, which ties the emitted asset name to the scaffolded wall that installs it, so the two cannot drift apart silently.

### Rejected: migrate to archives only

The conventional goreleaser layout, and what both package managers want. It breaks every adopted repository at its next CI run — the failure lands on the person least able to diagnose it, caused by a release process they do not control and cannot see. A compatibility window would only move the break; publishing both shapes costs a second archive stanza and ends the problem permanently.

### Rejected: a `curl | sh` installer script

One command on macOS and Linux, no tap repository, no upstream review. But it is a third artifact to keep matched to the release, it teaches piping a network fetch into a shell, and Homebrew already covers both platforms with a formula goreleaser writes for us. [CAP-001](../capabilities/CAP-001-onboarding/design.md)'s design has held since 2026-07-22 that the first encounter does not present installer scripts; nothing here disturbs that.

### Rejected for now: a Homebrew cask instead of a formula

`brews:` is deprecated upstream in favour of `homebrew_casks:`, and `goreleaser check` exits non-zero on it, so a cask is the forward-looking choice. It is rejected because **Homebrew casks are macOS-only**. Adopting one now would drop Linux back to the six manual steps this decision exists to remove — a whole platform losing the outcome, to satisfy a deprecation that has not yet bitten. A cask also carries the quarantine attribute our unsigned, un-notarized macOS binaries do not survive, where a formula installs into the Cellar without it.

The deprecation is therefore managed rather than ignored: **the release pins an exact goreleaser version rather than a version range**, so an upstream removal cannot break a release with no change on our side, and bumping the pin is a deliberate act that reruns the release dry run. Two things must be true before the migration: `homebrew_casks:` (or another channel) must cover Linux, and the macOS binaries must be signed and notarized. Until both hold, migrating trades a working install path for a tidier config.

### Deferred: winget in the same release

The winget publisher submits a pull request to `microsoft/winget-pkgs`, which is human-reviewed and can bounce for reasons unrelated to the binaries. Inside the release job that failure arrives *after* the artifacts are published and reads as a failed release. It ships in a patch tag once the publisher identity and token scope are proven, so "did we ship" stays separate from "did the upstream review land".

### Deferred: signing and notarization

The macOS binaries remain unsigned, and Windows users will see SmartScreen warnings on the downloaded fallback binary. A Homebrew formula sidesteps Gatekeeper and winget's portable install sidesteps the installer-signing question, so neither channel is blocked — but the manual fallback still needs the documented Gatekeeper exception, and that stays true until signing is its own decision with its own certificate management.
