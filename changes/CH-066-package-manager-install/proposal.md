---
id: CH-066
type: change
status: open
links: [P-007, M-030]
title: clue installs from a package manager, and the release assets adopters depend on become an append-only contract
---

# CH-066 — `clue` installs from a package manager

## What

The release pipeline moves from a hand-rolled `GOOS`/`GOARCH` loop to **goreleaser**, and gains a Homebrew tap that serves macOS and Linux from one formula. `winget install --exact --id Cliewen.Clue` and `brew install cliewen/tap/clue` replace six manual steps — read a platform table, download, verify a checksum by hand, rename, `chmod`, edit `PATH` — as the documented first way to get the binary. The manual route stays, demoted to a fallback and unchanged in substance.

Nothing an adopter already depends on moves. The release publishes **both** shapes: the existing bare binaries under their exact current names (`clue-<version>-<os>-<arch>[.exe]`) for the CI wall vendored into every repository that ran `clue init`, and the archives Homebrew and winget consume — all covered by one `SHA256SUMS`. The wall's `sha256sum -c --ignore-missing` makes the extra checksum lines harmless, so `guide/ci-wall.md` and every adopted repository need no change at all.

This change carries a **declared plan revision** to [P-007](../../docs/plans/P-007-core-hardening.md), riding with its implementing change under [PDR-008](../../docs/decisions/PDR-008-plan-revisions-may-ride.md) and called out for deliberate human approval in the pull request. The revision appends two milestones — M-030 (this change) and M-031 (a Claude Code marketplace bootstrap) — and reopens a scope the corpus had parked.

## Why

[C-015](../../docs/constraints/C-015-onboarding-under-30-minutes.md) sets a 30-minute bar from a bare machine to a green `clue validate`, and [CAP-001](../../docs/capabilities/CAP-001-onboarding/README.md) makes onboarding a capability. Six manual steps stand at the front of that path, and they are the first thing every newcomer meets. A package manager removes all six and makes the upgrade story a command rather than a procedure.

The scope was parked, not rejected. The [decision log](../../docs/decisions/log.md) row of 2026-07-22 (CH-045/P-004) placed installer scripts and package-manager channels outside that campaign "until supported by their own evidence and scope", and P-004, P-005 and P-006 each restated them as out of scope. That is a reversal-cheap row, and the evidence it asked for now exists: a stated onboarding constraint, a released binary on six platforms, and an install procedure whose length is measurable against the bar it must clear. P-007's own out-of-scope list never named installation, so nothing in the active campaign has to be retracted.

Adopting goreleaser is not incidental to the channels — it is what makes them cheap. A tap formula and a winget manifest are generated from the tag by the same tool that builds the binaries, so neither becomes a second place a version must be bumped by hand.

## Decision boundary

Two decisions, both expensive to reverse once published, so both are ADRs under [C-011](../../docs/constraints/C-011-decision-records-typed.md):

- **[ADR-030](../../docs/decisions/ADR-030-package-manager-distribution.md)** — distribution is package-manager-first, and the release asset names are an append-only public contract. The bare-binary names are not an implementation detail: they are executed by a vendored workflow in every adopted repository, where a rename would surface as a broken CI run with no local change to explain it.
- **[PDR-014](../../docs/decisions/PDR-014-distribution-campaign-reopened.md)** — the plan revision itself: why the parked scope reopens now, and why the work belongs on the active campaign rather than a successor plan.

This changes neither what `clue validate` checks nor what a merge accepts, so it is **not** a red-line change under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md). No corpus rule, acceptance criterion, capability meaning, or skill text moves. The release-cutting ritual is unchanged: nothing here bumps a version, regenerates skills, or tags.

winget is deliberately **not** in this change. Its first submission to `microsoft/winget-pkgs` is human-reviewed and can bounce for reasons unrelated to the binaries, and a failure inside the release job would read as a failed release after the artifacts are already published. It follows in a patch tag once the publisher identity and token scope are proven.
