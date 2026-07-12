---
id: CAP-002
type: capability
status: active
links: [G-001]
title: clue validate — the deterministic judge for the corpus
goal: G-001
---

# CAP-002 — `clue validate`

## What

`clue validate [--forbid-changes] [path]` scans `docs/` and `changes/`
for frontmatter artifacts and fails (exit 1) on any breach of the
corpus rules: missing frontmatter or core fields, duplicate IDs,
unresolvable `links` (milestones resolve inside plan files), status
values outside the per-type vocabulary, `/docs` folders without
README.md, index-block drift, and — with `--forbid-changes` — the
presence of a transient workspace (the digest-before-merge gate CI
runs).

## Why

The judge actor of [G-001](../../goals/G-001-verifiable-thread.md):
machines enforce form so humans only have to verify meaning. The same
binary runs locally and in CI — a green local run means a green PR.

Acceptance criteria: [criteria.md](criteria.md) · implementation notes:
[design.md](design.md).

## Status note

`active`: implemented and covered by Go tests whose names carry the
AC-IDs. The mechanical AC↔test link is the remaining half of
P-001/M-002 and lands as its own rule in a later change.
