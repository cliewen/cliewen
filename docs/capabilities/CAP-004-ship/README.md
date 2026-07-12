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

`clue` reports a release version (`clue version` / `clue --version`), stamped at build time from the release tag. Every agent skill carries a matching `version:` in its frontmatter, and `clue validate` makes drift between them lintable: a skill without a stamp fails, skills that disagree on a version fail, and a *released* `clue` whose skills differ from the binary fails as drift (a `dev` build skips that last comparison — it has no release to drift from).

A tagged release (`vX.Y.Z`) builds cross-platform binaries — linux/darwin/windows × amd64/arm64 — each stamped with the version, published as a GitHub release for `go install` and `gh release download`.

## Why

Delivers [G-002](../../goals/G-002-versioned-clue-and-skills.md): `go install` builds whatever the checkout has, and nothing told an adopted repo whether its installed skills or binary had drifted behind cliewen's main. A version on the binary plus a version marker in each skill makes drift detectable — and, with the new rule, lintable. The carrier rule ships method decisions as binary rules and skill text; without versions, drift between the judge (`clue`), the guidance (skills), and the corpus conventions is invisible until something breaks.

Acceptance criteria: [criteria.md](criteria.md) · design and the release pipeline: [design.md](design.md).

## Status note

`active`: `clue version`, the skill stamps, and the drift rule ship with tests carrying the AC-IDs (CH-007). The remaining half of M-004 — an adopted repo's CI actually installing and running `clue` (the model2diagram door) — is a follow-up change in that repo; M-004 stays `wip` until it closes.
