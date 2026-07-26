---
id: CH-066
type: change
status: open
links: [P-007, M-030]
title: One verified command installs clue on every supported platform, and the release assets adopters depend on become an append-only contract
---

# CH-066 — One verified command installs `clue`

## What

Installing `clue` becomes one command on macOS, Linux and Windows, through a checksum-verifying script published with the guide: `curl -fsSL https://cliewen.dev/install.sh | sh` or `irm https://cliewen.dev/install.ps1 | iex`. Each script detects the host, resolves the release, downloads the published binary, and verifies it against the release's `SHA256SUMS` **before writing anything**; a mismatch aborts with nothing installed. Neither needs elevation. They deploy with the guide site, so they require no repository, account, or credential that does not already exist. The manual download stays documented beneath them, unchanged.

The release pipeline moves from a hand-rolled `GOOS`/`GOARCH` loop to goreleaser, which publishes **both** asset shapes under one `SHA256SUMS`: the existing bare binaries under their exact current names, and archives a later package-manager channel can consume. The adopters' CI wall reads its checksums with `--ignore-missing`, so the extra lines are inert and no adopted repository changes — `guide/ci-wall.md` needs no edit.

This change carries a **declared plan revision** to [P-007](../../docs/plans/P-007-core-hardening.md), riding with its implementing change under [PDR-008](../../docs/decisions/PDR-008-plan-revisions-may-ride.md) and called out for deliberate human approval in the pull request. The revision appends M-030 (this change) and M-031 (a Claude Code marketplace bootstrap), reopening a scope the corpus had parked.

## Why

[C-015](../../docs/constraints/C-015-onboarding-under-30-minutes.md) gives the whole journey — install, `clue init`, `clue validate` — thirty minutes, and six manual steps stand at the front of it: read a platform table, download, verify a checksum by hand, rename, `chmod`, edit `PATH`. They are the first thing every newcomer meets.

The scope was parked, not rejected. The [decision log](../../docs/decisions/log.md) row of 2026-07-22 (CH-045/P-004) placed installer scripts and package-manager channels outside that campaign "until supported by their own evidence and scope", and P-004, P-005 and P-006 each restated them. That is a reversal-cheap row, and the evidence it asked for now exists: a stated onboarding constraint, a released binary on six platforms, and an install procedure whose length is measurable against the bar it must clear.

No package manager leads, and that is the finding rather than a preference. Homebrew appeared to cover macOS and Linux from one formula, but formula generation is deprecated upstream and its replacement — casks — is macOS-only. Every remaining Linux channel needs infrastructure this project does not have: `.deb` and `.rpm` need a hosted repository to become one command, Snap needs a store account and review, AUR and Nix reach one distribution each. So the real choice was a verified script or six manual steps, not a script or `brew install`. winget and a macOS cask remain wanted; they are additive and follow as their own changes, which also keeps "can people install this" separate from "has an upstream review landed".

## Decision boundary

Two decisions, both expensive to reverse once published, so both are records under [C-011](../../docs/constraints/C-011-decision-records-typed.md):

- **[ADR-030](../../docs/decisions/ADR-030-verified-install-scripts.md)** — a checksum-verifying install script is the primary channel, and the release asset names are an append-only public contract. The bare names are not an implementation detail: they are executed by a workflow vendored into every adopted repository and now by a script strangers are asked to pipe into a shell. The verification is the condition under which that pattern is defensible at all.
- **[PDR-014](../../docs/decisions/PDR-014-distribution-reopens-on-the-active-campaign.md)** — the plan revision: why the parked scope reopens now, and why the work belongs on the active campaign rather than a successor plan.

ADR-030 supersedes one clause of [CAP-001](../../docs/capabilities/CAP-001-onboarding/design.md)'s accumulated design lessons — that the first encounter does not present installer scripts. That lesson is revised in place rather than contradicted silently: it was written while a package manager was assumed reachable, and it now records what changed and on what terms a script may lead.

This changes neither what `clue validate` checks nor what a merge accepts, so it is **not** a red-line change under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md). No corpus rule, acceptance criterion, capability meaning, or skill text moves. The release-cutting ritual is unchanged: nothing here bumps a version, regenerates skills, or tags.
