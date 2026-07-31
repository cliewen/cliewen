---
id: ADR-030
type: decision
status: inferred
links: [CAP-001, CAP-004, C-015, ADR-011, ADR-012, ADR-038]
title: Installation is a checksum-verifying script, and release asset names are an append-only contract
author: agent
accepted-by: []
---

# ADR-030 — Verified install scripts and append-only release assets

## Context and problem statement

Getting `clue` onto a machine takes six manual steps: read a platform table, download the right asset, verify its SHA-256 by hand, rename it, make it executable, and put it on `PATH`. [C-015](../constraints/C-015-onboarding-under-30-minutes.md) gives the whole journey — install, `clue init`, `clue validate` — thirty minutes, and those six steps are what a newcomer meets first.

The obvious answer is a package manager per platform. Only one of them is straightforward: winget covers Windows. Homebrew appears to cover macOS *and* Linux from a single formula, but formula generation is deprecated in the release tooling, and its replacement — casks — is macOS-only. Every remaining Linux channel needs infrastructure this project does not have: `.deb` and `.rpm` need a hosted apt or yum repository to become one command, Snap needs a store account and review, AUR and Nix reach one distribution each.

There is a second problem underneath. The release publishes **bare binaries**, and `internal/scaffold/templates/github/workflows/clue.yml` — the thin caller emitted into every repository that has run `clue init` — delegates to the upstream reusable workflow, which verifies `SHA256SUMS` and stages `clue-${CLUE_VERSION}-linux-amd64` without assuming a root-only install path. Any channel that repackages the release must not disturb those names.

## Decision outcome

**One command installs `clue` on all three platforms, through a checksum-verifying install script published with the guide. Native package managers are additive and follow. Release asset names are an append-only public contract.**

- **The scripts are the primary channel.** `guide/public/install.sh` serves macOS and Linux; `install.ps1` serves Windows. They deploy with the guide site, so `https://cliewen.dev/install.sh` needs no repository, account, or credential that does not already exist. Each detects the host, resolves the release, downloads the published asset, and **verifies it against `SHA256SUMS` before writing anything** — a failed comparison aborts with nothing installed. Neither needs elevation: both target a directory the user owns.
- **Verification is what makes the pattern acceptable.** Piping a network fetch into a shell is defensible only when what it installs is checked against checksums the project published separately, and when the script is served from the project's own domain over TLS and can be read before it is run. A script that skipped the checksum would be strictly worse than the manual steps it replaces, because it would remove the one step that was protecting the user while keeping the appearance of convenience.
- **Assets are append-only.** Every release publishes both shapes under one `SHA256SUMS`: the bare binaries under their exact existing names `clue-<version>-<os>-<arch>[.exe]`, and archives for the package managers that will follow. A new channel may add assets; it may never rename or remove one. The upstream reusable workflow selects the entry for the binary it is about to execute and fails when the file lists no such entry, so additional lines are inert without a present-but-unlisted binary being admitted unverified. The install scripts consume the same names, so the contract now has two kinds of dependents rather than one.
- **The manual download stays documented.** It is the route for a host where neither script can run, and it is where the asset names are written down for a human.
- **Release notes are unaffected.** The release tooling's own changelog generation is disabled and the body is the extracted `CHANGELOG.md` section, preserving [ADR-012](ADR-012-release-notes-from-changelog.md) exactly.

The asset-name rule is the load-bearing half. A release artifact whose name appears in a file shipped into other people's repositories — and now in a script strangers are asked to pipe into a shell — is a published interface, not an implementation detail. Nothing about `clue-0.7.0-linux-amd64` announces that; naming it a contract is what makes the constraint reviewable instead of rediscovered.

This revises one of [CAP-001](../capabilities/CAP-001-onboarding/design.md)'s accumulated design lessons. That lesson made release binaries the primary install path and kept installer scripts out of a newcomer's first encounter. It was written while a package manager was assumed reachable; the manual route becomes the documented fallback, and the lesson records the revision in place rather than being deleted.

**Carrier:** `guide/public/install.sh` and `guide/public/install.ps1`; `.goreleaser.yaml` (machine); `TestSanity_ReleaseKeepsTheAssetNamesTheAdopterWallInstalls` and `TestSanity_InstallScriptsUseTheReleaseAssetContract`, which tie the emitted asset names to both dependents so they cannot drift apart silently.

### Rejected: a Homebrew formula covering macOS and Linux

The original intent, and the reason this decision first read "package-manager-first". A single formula genuinely serves both platforms, and Homebrew is what a developer on either reaches for. It is rejected because the release tooling deprecated formula generation: keeping it would mean either pinning tooling against an announced removal, or hand-maintaining a formula and reintroducing the manual version bump this project removes everywhere else.

### Rejected: a Homebrew cask as the macOS and Linux channel

The supported successor, so it appears to be the drop-in replacement. Homebrew casks are macOS-only. Adopting one as *the* channel would leave Linux with the six manual steps this decision exists to remove — an entire platform losing the outcome in exchange for tidier configuration. A cask remains attractive **in addition**, once it is not the only thing standing between a Linux user and a working install.

### Rejected: `.deb` and `.rpm` packages instead of a script

Native, familiar, and generated by the same release tooling. Without a hosted apt or yum repository they are not one command — they are "download this file, then run `dpkg -i`", barely fewer steps than the manual route and with no automatic upgrade. Hosting a package repository is a service dependency and a running cost that a script on an already-deployed site does not incur. Such packages stay available as a later addition for users who prefer distribution-native formats.

### Rejected: an unverified `curl | sh`

The common shape, and the reason the pattern has a poor reputation. Downloading a binary over TLS and executing it without checking it against published checksums trusts the transport and the release host completely, and discards a verification step the manual instructions already required. The verification is not decoration on this decision; it is the condition under which the decision is defensible at all.

### Deferred: winget, and a Homebrew cask for macOS

Both are wanted and neither is blocking. winget submits a pull request to `microsoft/winget-pkgs`, reviewed upstream on its own schedule; a cask needs a tap repository and a credential. Adding them once the scripts already deliver a working install is additive, and it separates "can people install this" from "has an upstream review landed". Each becomes its own change with its own evidence.

### Deferred: signing and notarization

The macOS binaries remain unsigned, and Windows users may meet SmartScreen warnings. The install scripts avoid Gatekeeper's quarantine attribute because they download outside the browser, but the manual fallback still needs the documented exception, and a Homebrew cask would need real signing. That stays its own decision with its own certificate management.
