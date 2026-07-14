# Changelog

All notable, user-visible changes to `clue` and the Cliewen skills. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow semver. Each GitHub release body is this file's matching version section, extracted verbatim by the release workflow — a release with no section here fails.

## [Unreleased]

### Added

- **Light change tier**: a change that decides nothing and changes no meaning (typos, doc clarity, dependency bumps, pure refactors, CI plumbing) no longer needs a `/changes/` workspace — branch, commit, and open a PR whose description is the proposal. The moment a decision, open question, or AC change appears, escalate to the full loop. The `clue-delta` skill carries the qualification test; `clue-verify` starts by checking the tier is right.
- **Decision log**: ADRs are now reserved for decisions that are expensive to reverse; everything else is a dated row in `docs/decisions/log.md` (litmus test: cheap and local to reverse → log row). `clue validate` lints the new `log` artifact type.

### Changed

- Release pages are now published from this changelog: each GitHub release body is the matching version section of `CHANGELOG.md`, written for users — no more auto-generated PR lists.
- The agent skills are repo-agnostic: repo-local conventions live in your AGENTS.md, which extends the methodology but never overrides it — a conflict between the two is surfaced as an open question. The OpenSpec extraction mapping moved from `clue-extract`'s skill text to `mappings/openspec.md` under the same skill, and the skills no longer cite cliewen's internal document IDs — every rule is stated in full where you read it.

## [0.1.0] - 2026-07-13

First release of `clue`, the command-line validator for the Cliewen methodology. It keeps a repository's documentation corpus — goals, plans, capabilities, decisions — and its agent skills consistent, traceable, and versioned.

### Added

- **`clue validate`** lints the corpus: frontmatter core fields, ID uniqueness, status vocabulary, cross-links, folder READMEs and index blocks, provenance, and acceptance-criteria-to-test traceability.
- **`clue version`** reports the release the binary was built from; untagged source builds report `dev`.
- **Versioned skills**: every agent skill declares a version in its frontmatter, and `validate` fails when skills disagree with each other or drift from the binary's release.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.1.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`.
