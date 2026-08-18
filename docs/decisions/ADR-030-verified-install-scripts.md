---
id: ADR-030
type: decision
status: verified
links: [CAP-001, CAP-004, C-015, ADR-011, ADR-012, ADR-038]
title: Installation is a checksum-verifying script, and release asset names are an append-only contract
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-030 — Verified install scripts and append-only release assets

## Context and problem statement

Manual installation is a six-step barrier, while package-manager coverage is incomplete and the reusable adopter workflow depends on stable bare binary names and `SHA256SUMS`.

## Decision outcome

**Make checksum-verifying install scripts the primary one-command channel, and treat release asset names as append-only.** `guide/public/install.sh` serves macOS and Linux; `install.ps1` serves Windows. Each resolves a release, downloads the asset, verifies it against `SHA256SUMS` before writing, and installs into a user-owned directory. The manual route remains documented. Native package managers may add channels later.

Every release keeps the existing bare names `clue-<version>-<os>-<arch>[.exe]` and may add archives, never rename or remove an existing name. The upstream workflow and both scripts consume the same names. Release notes continue to follow ADR-012. The carrier is the scripts, `.goreleaser.yaml`, and the two release-contract sanity tests. Homebrew-only, repository-less native packages, and unverified `curl | sh` are rejected; winget, a macOS cask, signing, and notarization are deferred.
