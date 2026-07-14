---
id: C-002
type: constraint
status: active
links: [ADR-012]
title: Every user-visible change adds a changelog entry in its digest
source: AGENTS.md rule 6, ADR-012
enforcement: agent
---

# C-002 — Every user-visible change adds a changelog entry

A change with user-visible impact adds its entry to the `[Unreleased]` section of `CHANGELOG.md` in the digest — written for users, never a PR title or commit subject. The release workflow already fails a release whose version section is missing (`machine` at release time); per-change presence is agent-held.

**Promotion trigger:** `clue` gains git-diff context and can require a `CHANGELOG.md` hunk in changes that touch user-facing surfaces — then `enforcement: machine`.
