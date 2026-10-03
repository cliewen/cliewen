---
id: CH-197
type: change
status: open
links: [G-024, CAP-004]
title: Generated skills describe when to use them
---

# CH-197 — Generated skills describe when to use them

This change is plan-less: it carries out [G-024](../../docs/goals/G-024-generated-skills-describe-when-to-use-them.md), which the maintainer asked to start after G-023 merged.

## Why

The generated skill entry points carry `cliewen-skill`, `version`, `type`, and `title`, but no `name` or `description`. Claude Code falls back to the first body line, so all six `clue-*` skills show the generated-file comment as their description, and a host that picks skills by description cannot tell them apart. The Agent Skills format that Claude Code and other agent hosts read expects `name` and `description` in a skill's frontmatter. Adopters receive the same frontmatter in `.agents/skills` and their `.claude/skills` copies.

## What changes

- Each generated entry point gains `name`, equal to its directory, and `description`, which states what the skill does and when to use it.
- The description comes from the skill definition the generator already holds. Its existing one-line summary stays the opening sentence of the routing body, and the definition gains a short when-to-use clause. The frontmatter joins the two. One source feeds both, so they cannot disagree.
- A new criterion in CAP-004 states the contract: every entry point names its directory and describes when to use it, and no two skills share a description. A description must not be the generated-file comment, and must not exceed the 1024 characters the format allows.
- The release ships the new skill bytes through the existing managed-carrier refresh, so `clue migrate` needs no new migration. A CHANGELOG entry names the generated skills.

## Scope boundary

- The entry-point file stays `skill.md` in `.agents/skills`. Renaming it to `SKILL.md` everywhere would change what `clue validate`, `clue migrate`, and every adopter's tree expect, and G-023 already met the Claude Code need in this repository by symlink. That rename is out of scope.
- The hubs' skill tables keep their own wording. They route a reader and are not a second copy of the description.

## Challenge

**Riskiest assumption:** adding `name` and `description` to an entry point that also carries `cliewen-skill`, `version`, `type`, and `title` is read by agent hosts the way the format intends. A host must not reject the extra keys or rename the skill.

**Credible alternative:** leave the frontmatter alone and move the generated-file comment below the title, so the first body line is the summary. That fixes the symptom for Claude Code's fallback alone, and depends on behaviour the format does not promise.

**Cheapest test:** in a scratch repository, give one skill the new frontmatter. Confirm that Claude Code lists it under its directory name with the description, and still loads its references.

**Stop or revise:** if a host lists the skill under another name or drops it, stop and bring the evidence to the maintainer.

**Meeting the criterion and still failing the person:** each description is non-empty and distinct, but restates the skill's name without saying when to reach for it, so a host still cannot choose. The when-to-use clause exists to prevent that, and the review checks each description against its skill's routing hub entry.
