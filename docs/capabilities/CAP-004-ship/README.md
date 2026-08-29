---
id: CAP-004
type: capability
status: active
links: [G-002, ADR-043]
title: clue ships — a versioned binary and versioned skills
goal: G-002
---

# CAP-004 — clue ships

## What

`clue` reports a release version (`clue version` / `clue --version`), stamped at build time from the release tag. The six managed Cliewen skills — `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, `clue-upgrade`, and `clue-verify` — declare `cliewen-skill: true` and carry a matching `version:` in their frontmatter ([ADR-022](../../decisions/ADR-022-skill-ownership-marker.md), [ADR-043](../../decisions/ADR-043-upgrade-skill-is-a-managed-carrier.md)); `clue validate` scopes the version set to those marked skills, so unrelated skills can coexist under `.agents/skills/`. A marked skill without a stamp fails, marked skills that disagree on a version fail, and a *released* `clue` whose marked skills differ from the binary fails as drift (a `dev` build skips that last comparison — it has no release to drift from). Each standalone skill directory is generated from skill-specific templates and shared instruction fragments: a short `skill.md` routes to required local references only when their workflow condition is reached, while repository tests reject drift between every canonical entry point or reference and either distributed skill tree ([ADR-059](../../decisions/ADR-059-progressive-standalone-skill-directories.md)).

Whether the installed release is still the newest one cannot be answered from the repository, so `clue latest` asks the published release list and reports: both versions, and — when behind — the installation route for the machine it is running on plus the `clue migrate` sequence that moves the repository with it. It reaches the network, so it stays outside the deterministic judge and must never be a required check; it writes no file in the repository and never replaces the binary ([ADR-042](../../decisions/ADR-042-release-check-outside-the-judge.md)). Not being able to tell is reported calmly and exits 0, `--quiet` is silent unless there is something to say, and the answer is cached in the user's cache directory for a day. The drift report names the same two routes: forward, or staying on the release the repository carries.

A tagged release (`vX.Y.Z`) builds cross-platform binaries — linux/darwin/windows × amd64/arm64 — each stamped with the version, published as a GitHub release for `go install` and `gh release download`. When a release changes corpus obligations, its notes name the supported `clue migrate` path before adopters upgrade.

## Why

Delivers [G-002](../../goals/G-002-versioned-clue-and-skills.md): `go install` builds whatever the checkout has, and nothing told an adopted repo whether its installed skills or binary had drifted behind cliewen's main. Ownership and version markers on each Cliewen skill make the managed set explicit and its drift detectable — and lintable — without absorbing third-party skills that share the standard directory. The carrier rule ships method decisions as binary rules and skill text; without versions, drift between the judge (`clue`), the guidance (skills), and the corpus conventions is invisible until something breaks.

Acceptance criteria: [criteria.md](criteria.md) · design and the release pipeline: [design.md](design.md).

## Status note

`active`: `clue version`, the skill stamps, release artifacts, and the drift rule ship with tests carrying the AC-IDs. AC-147 additionally holds the shared ADR/PDR/IDR subject-routing and compact-shape instruction across every generated lifecycle skill and both distributed trees. AC-148 keeps upgrade discovery honest by requiring a no-apply migration preview after the release availability check, even when the installed release is newest. P-002's completed M-004 row preserves the historical adopter-CI evidence; the analysis index identifies the private-repository boundary around that evidence. The capability's normative contract and verification remain entirely in this repository.
