# Changelog

All notable, user-visible changes to `clue` and the Cliewen skills. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow semver. Each GitHub release body is this file's matching version section, extracted verbatim by the release workflow — a release with no section here fails.

## [Unreleased]

### Changed

- Release pages are now published from this changelog: each GitHub release body is the matching version section of `CHANGELOG.md`, written for users — no more auto-generated PR lists.
- The agent skills are repo-agnostic: repo-local conventions live in your AGENTS.md, which extends the methodology but never overrides it — a conflict between the two is surfaced as an open question. The OpenSpec extraction mapping moved from `clue-extract`'s skill text to `mappings/openspec.md` under the same skill.

## [0.1.0] - 2026-07-13

First release of `clue`, the command-line validator for the Cliewen methodology. It keeps a repository's documentation corpus — goals, plans, capabilities, decisions — and its agent skills consistent, traceable, and versioned.

### Added

- **`clue validate`** lints the corpus: frontmatter core fields, ID uniqueness, status vocabulary, cross-links, folder READMEs and index blocks, provenance, and acceptance-criteria-to-test traceability.
- **`clue version`** reports the release the binary was built from; untagged source builds report `dev`.
- **Versioned skills**: every agent skill declares a version in its frontmatter, and `validate` fails when skills disagree with each other or drift from the binary's release.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.1.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`.
