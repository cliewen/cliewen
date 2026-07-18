---
id: CAP-004
type: capability
status: active
links: [G-002]
title: clue ships — a versioned binary and versioned skills
goal: G-002
---

# CAP-004 — clue ships

## What

`clue` reports a release version (`clue version` / `clue --version`), stamped at build time from the release tag. Every Cliewen agent skill declares `cliewen-skill: true` and carries a matching `version:` in its frontmatter ([ADR-022](../../decisions/ADR-022-skill-ownership-marker.md)); `clue validate` scopes the version set to those marked skills, so unrelated skills can coexist under `.agents/skills/`. A marked skill without a stamp fails, marked skills that disagree on a version fail, and a *released* `clue` whose marked skills differ from the binary fails as drift (a `dev` build skips that last comparison — it has no release to drift from). The standalone skills are generated from skill-specific templates and shared instruction fragments; repository tests also reject drift between those canonical sources and either distributed skill tree.

A tagged release (`vX.Y.Z`) builds cross-platform binaries — linux/darwin/windows × amd64/arm64 — each stamped with the version, published as a GitHub release for `go install` and `gh release download`.

## Why

Delivers [G-002](../../goals/G-002-versioned-clue-and-skills.md): `go install` builds whatever the checkout has, and nothing told an adopted repo whether its installed skills or binary had drifted behind cliewen's main. Ownership and version markers on each Cliewen skill make the managed set explicit and its drift detectable — and lintable — without absorbing third-party skills that share the standard directory. The carrier rule ships method decisions as binary rules and skill text; without versions, drift between the judge (`clue`), the guidance (skills), and the corpus conventions is invisible until something breaks.

Acceptance criteria: [criteria.md](criteria.md) · design and the release pipeline: [design.md](design.md).

## Status note

`active`: `clue version`, the skill stamps, and the drift rule ship with tests carrying the AC-IDs (CH-007). The other half of M-004 — an adopted repo's CI actually installing and running `clue` — closed in model2diagram's CH-003 (its PR #3, merged 2026-07-13): the `validate` job runs a vendored, checksum-verified release binary, so an adopter needs no cross-repo credentials. M-004 is `done`.
