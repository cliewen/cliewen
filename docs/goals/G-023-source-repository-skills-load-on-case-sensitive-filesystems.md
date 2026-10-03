---
id: G-023
type: goal
status: proposed
links: [VIS-001]
title: This repository's generated skills load in Claude Code on case-sensitive filesystems
---

# G-023 — This repository's generated skills load in Claude Code on case-sensitive filesystems

**Who wants it:** found by the agent on 2026-10-03 while probing raw-reading consumers for CH-196. It is recorded here because it is unrelated to that change's frontmatter work.

**Why:** `CLAUDE.md` says the `.claude/skills` symlink to `.agents/skills` makes Claude Code list the `clue-*` skills. Claude Code 2.1.288 on Linux discovers a skill only through an upper-case `SKILL.md`, and the symlink exposes the generated lower-case `skill.md`, so a session in this repository on a case-sensitive filesystem lists none of them. In a scratch repository, the same skill was not listed as `skill.md` and was listed and used as `SKILL.md`. Adopters are not affected: `clue init` writes `.claude/skills/<name>/SKILL.md` copies rather than a symlink.

**Success looks like:** a Claude Code session in this repository on Linux lists and can load every `clue-*` skill, and nothing about it depends on the filesystem's case sensitivity.
