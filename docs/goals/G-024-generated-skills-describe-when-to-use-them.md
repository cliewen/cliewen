---
id: G-024
type: goal
status: proposed
links: [VIS-001]
title: Generated skills describe when to use them, so an agent host can choose them
---

# G-024 — Generated skills describe when to use them

**Who wants it:** found by the agent on 2026-10-03 when the `clue-*` skills first appeared in a Claude Code session here, after G-023 was met.

**Why:** the generated skill frontmatter carries `cliewen-skill`, `version`, `type`, and `title`, but no `description`. Claude Code then shows the first line of the body as each skill's description, which is the generated-file comment `<!-- Generated from Cliewen's canonical skill sources; edit those sources, not this file. -->` for all six skills. A host that picks skills by description cannot tell them apart, so the skills reach an agent only through the routing hub's table. Adopters receive the same frontmatter in their `.claude/skills` copies.

**Open before a change commits to it:** whether the description is derived from each skill definition's existing one-line summary or written separately, and whether other agent hosts need the same field under another name.

**Success looks like:** each generated skill's description says when to use it, and a Claude Code session lists the six skills with distinct, useful descriptions.
