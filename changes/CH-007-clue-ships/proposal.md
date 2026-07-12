---
id: CH-007
type: change
status: open
links: [P-002, M-004, G-002, CAP-004, ADR-011]
title: clue ships — versioned binary, per-skill version stamps, drift lint, release pipeline
---

# CH-007 — clue ships

## What

Serves **[P-002](../../docs/plans/P-002-leaves-home.md) / M-004** and accepts **[G-002](../../docs/goals/G-002-versioned-clue-and-skills.md)** (clue and the skills carry versions). This change:

1. Gives `clue` a release version: `clue version` (and `clue --version`) report a stamp injected at build time via `-ldflags -X main.version=…`; source/local builds report `dev`.
2. Stamps every agent skill with a `version:` in its frontmatter, and teaches `clue validate` to enforce three things (new capability **[CAP-004](../../docs/capabilities/CAP-004-ship/README.md)**):
   - every skill carries a version stamp (AC-020);
   - the skills agree on one version — "versioned as a set" via per-skill markers (AC-021);
   - a *released* `clue` whose skills' version differs from the binary fails as **drift**; a `dev` build skips the comparison (AC-022).
3. Adds a tagged-release pipeline (`.github/workflows/release.yml`): a `v*` tag builds cross-platform binaries (linux/darwin/windows × amd64/arm64), each stamped with the tag's version, and publishes them as a GitHub release for `gh release download`.
4. Documents the install story (`go install` and `gh release download`) in the README.

## Why

`go install` builds whatever the checkout has, and nothing tells an adopted repo whether its installed skills or binary have drifted behind cliewen's main (G-002). A version on the binary plus a version marker in each skill makes drift detectable — and, with the new `clue validate` rule, lintable.

## Doors resolved (recorded in ADR-011)

- **Version source:** git tags (`vX.Y.Z`) via ldflags; the stamped/skill string is bare semver `X.Y.Z` (the workflow strips the `v`).
- **Skill granularity:** per-skill frontmatter `version:` (not a single set-file), kept consistent by the AC-021 rule.
- **Drift policy:** a failure, not a warning (the wall philosophy — a warning gets ignored); `dev` builds are exempt.

## Scope

Cliewen's side of M-004: the release pipeline, the version surface, the drift lint, and the install docs. Actually wiring an adopted repo's CI to install and run `clue` (the model2diagram door) is a follow-up change in that repo; M-004 stays `wip` until it closes.

## Out of scope

A `clue`-driven version bump command (release still edits the skill frontmatter by hand); drift detection for skills installed outside `.agents/skills/`; public release / repo-visibility decisions (P-002 already parks these).
