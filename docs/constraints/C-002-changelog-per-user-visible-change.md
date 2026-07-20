---
id: C-002
type: constraint
status: active
links: [ADR-012]
title: Every release-relevant user-visible change adds a changelog entry
source: AGENTS.md rule 6, ADR-012
enforcement: agent
---

# C-002 — Every release-relevant user-visible change adds a changelog entry

A Cliewen change that affects shipped behavior, a capability, a contract, a command, or a user workflow adds its entry to the `[Unreleased]` section of `CHANGELOG.md` in the digest — written for users, never a PR title or commit subject. A plain editorial change under [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) is not release history and adds no entry; prose that changes normative instructions or user workflow is not plain. The release workflow already fails a release whose version section is missing (`machine` at release time); per-change presence is agent-held.

**Promotion trigger:** `clue` gains git-diff context and a reliable signal that distinguishes release-relevant behavior, contract, command, and workflow changes from plain editorial work, then can require the `CHANGELOG.md` hunk for the former — at that point this rule becomes `enforcement: machine`.
