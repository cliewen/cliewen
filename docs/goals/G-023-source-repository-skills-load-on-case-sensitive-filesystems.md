---
id: G-023
type: goal
status: accepted
links: [VIS-001]
title: This repository's generated skills load in Claude Code on case-sensitive filesystems
---

# G-023 — This repository's generated skills load in Claude Code on case-sensitive filesystems

**Who wants it:** found by the agent on 2026-10-03 while probing raw-reading consumers for CH-196, and taken up by the maintainer the same day.

**Why:** Claude Code discovers a skill only through an upper-case `SKILL.md` on a case-sensitive filesystem; observed with Claude Code 2.1.288 on Linux. The generated skills are named `skill.md`, so a single `.claude/skills` symlink to `.agents/skills` listed none of them in a session here. Adopters are not affected, because `clue init` writes `.claude/skills/<name>/SKILL.md` copies.

**How it is met:** `.claude/skills` holds one directory per generated skill, whose `SKILL.md` and remaining entries are symlinks to the generated files, so the mirror keeps no copy that could drift. `TestSanity_ClaudeMirrorLinksEveryGeneratedSkill` in `internal/skills` fails when a generated skill has no mirror, a link points elsewhere, or a mirror outlives its skill.

**Success looks like:** a Claude Code session in this repository on Linux lists and can load every `clue-*` skill.
