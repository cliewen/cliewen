---
id: CH-008
type: change
status: open
links: [P-002, M-004, G-002, CAP-004]
title: Release notes come from CHANGELOG.md — user-facing, reviewed, drift-proof
---

# CH-008 — Release notes come from CHANGELOG.md

## What

Serves **[P-002](../../docs/plans/P-002-leaves-home.md) / M-004** as a refinement of the release pipeline that milestone delivered: the first tagged release (v0.1.0) exposed that `generate_release_notes: true` publishes a PR dump with contributor @mentions — internal change history, not something a user of `clue` can read. This change makes release notes a reviewed, user-facing artifact:

1. **`CHANGELOG.md` at the repo root** ([Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format) becomes the single source of truth for user-visible history: an `[Unreleased]` section that accumulates entries as changes merge, and one section per released version. The v0.1.0 section back-fills the current release.
2. **The release workflow derives the release body from it**: the `## [X.Y.Z]` section matching the tag is extracted verbatim and published via `body_path`; `generate_release_notes` is dropped. A missing or empty section fails the release before anything is built — the same wall philosophy as the version-drift lint: a release cannot ship without user-facing notes, and the page can never drift from the file.
3. **The agent instructions carry the convention**: AGENTS.md gains the binding rule, `clue-delta`'s digest step adds "record the user-visible impact in `[Unreleased]`", and `clue-verify` gets the matching checkbox.

## Why

Release pages are the first thing an adopter reads; v0.1.0's auto-generated body described cliewen's PR history to its own (sole) author. Writing notes at tag time from memory is how that happens — accumulating them in `[Unreleased]` at merge time, then renaming the section at release, means the notes are written when the change is fresh and reviewed like everything else. Deriving the page from the file makes the 1-1 map structural rather than disciplined.

## Doors resolved (recorded in ADR-012)

- **Source of truth:** one root `CHANGELOG.md`, not per-release files under `docs/` and not annotated tag messages (both rejected in the ADR).
- **Enforcement:** extraction failure is a release failure, not a fallback to auto-generated notes.

## Scope

`CHANGELOG.md` (with `[Unreleased]` and the back-filled `[0.1.0]`), the `release.yml` extraction + guard, a sanity test pinning the workflow shape, the AGENTS.md rule, the `clue-delta`/`clue-verify` skill text, ADR-012, and CAP-004 design bookkeeping. After merge (operational, outside the diff): re-sync the published v0.1.0 release body to the changelog section so the map is 1-1 from the first release.

## Out of scope

A `clue` rule that lints `CHANGELOG.md` structure (a door noted in ADR-012 — the workflow guard already enforces the invariant where it bites); shipping the convention to adopters via `clue init` templates (that is M-005's job); any skill version bump (the skills still say 0.1.0; the next release bumps them as usual).
