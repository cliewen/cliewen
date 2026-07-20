---
id: ADR-012
type: decision
status: verified
links: [P-002, CAP-004, ADR-011]
title: Release notes are user-facing and come from CHANGELOG.md — extracted verbatim, missing section fails the release
author: agent
accepted-by: Flemming N. Larsen (2026-07-13, PR #7)
---

# ADR-012 — Release notes come from CHANGELOG.md

## Context and problem statement

Release pages are the first thing an adopter of [CAP-004](../capabilities/CAP-004-ship/README.md)'s pipeline sees; they need what-it-is, how-to-install, what-changed-for-me prose. GitHub's auto-generated release notes produce the opposite — a "What's Changed" list of PR titles with contributor @mentions: the repo's internal change history, not the product's (the first tagged release shipped exactly that). Where do user-facing release words live, who writes them when, and what stops the PR dump from coming back?

## Decision outcome

**A root `CHANGELOG.md` ([Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format) is the single source of truth for user-visible history; the release workflow extracts the tag's `## [X.Y.Z]` section verbatim as the release body and fails the release when the section is missing or empty.**

- **Written at merge time, not tag time.** Each change records its user-visible impact in the `[Unreleased]` section during the digest — phrased for a user of the tool, not a reviewer of the repo. Cutting a release renames `[Unreleased]` to the version and tags; the notes were already written and reviewed when the changes were fresh. Notes that do not exist until tag time get reconstructed from change history — which is the failure this decision removes.
- **The 1-1 map is structural.** The workflow publishes the section through `body_path`, so the release page cannot say anything the reviewed file does not. The extraction guard follows the wall philosophy ([ADR-011](ADR-011-version-stamping.md)'s drift rule, applied to prose): no section, no release — a rule that only warned would be ignored.
- **Auto-generation is banned, and lintably so.** `generate_release_notes` is removed and `TestSanity_ReleaseNotesComeFromChangelog` fails the build if it (or the extraction's absence) reappears in `release.yml`.

**Carrier:** the `release.yml` extraction + guard and the sanity test (machine); AGENTS.md rule 7 (agent — the repo-local convention layer). The skills are **not** a carrier: they ship verbatim to every adopting repo and stay generic, so they reference only "repo-local digest conventions in AGENTS.md", never `CHANGELOG.md` by name — the changelog is this repo's convention, not adopter methodology (the general principle is [ADR-013](ADR-013-ships-generic-vs-repo-local.md)). The `clue init` template may offer the convention as a default once M-005 ships (a door, not yet built).

### Rejected: GitHub's `generate_release_notes`

It describes PRs, not the product; it @-mentions contributors — internal credit, not product information; and it writes the page at tag time from data no human reviewed as release prose.

### Rejected: per-release note files (`docs/releases/vX.Y.Z.md`)

Workflow-simpler (`body_path` directly), but it scatters history across files nobody reads in sequence, puts user-facing prose in a nonstandard location, and loses the `[Unreleased]` accumulator — the mechanism that makes notes get written at merge time is the point, not the file layout.

### Rejected: annotated tag messages as the body

Tag messages bypass the PR review gate — the release page would be the one user-facing artifact no human verified. They are also invisible in the repo's rendered view.

### Deferred: a `clue` rule linting CHANGELOG.md structure

`clue validate` could check that `[Unreleased]` exists or that a stamped version has a section. The workflow guard already enforces the invariant at the moment it bites (the release), so a corpus rule is a door for when adopters carry the convention (M-005), not part of this decision.
