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

Adopters need release prose about the product, not an unreviewed list of pull requests or mentions.

## Decision outcome

**Root `CHANGELOG.md` in Keep a Changelog format is the source of truth.** Each change writes user-facing impact under `[Unreleased]`; a release renames that section to `## [X.Y.Z]`, extracts it verbatim, and fails if it is missing or empty. The workflow verifies the published release body matches the extracted section and rejects GitHub's `generate_release_notes`.

**Carrier:** release extraction, post-publish body verification, the sanity test, and this repository's AGENTS.md convention. Skills stay generic and do not name this repository's changelog. Per-release files, annotated tag messages, and warning-only structure lint are rejected or deferred because they bypass review, scatter the history, or add a second carrier.
