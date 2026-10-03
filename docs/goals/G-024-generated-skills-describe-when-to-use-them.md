---
id: G-024
type: goal
status: accepted
links: [VIS-001]
title: Generated skills describe when to use them, so an agent host can choose them
---

# G-024 — Generated skills describe when to use them

**Who wants it:** found by the agent on 2026-10-03 when the `clue-*` skills first appeared in a Claude Code session here, after G-023 was met.

**Why:** an agent host that picks skills by description, such as Claude Code, falls back to the first body line when the frontmatter has none. For the generated skills that line is the generated-file comment, identical for all six, so the host cannot tell them apart, and adopters receive the same files.

**How it is met:** each generated entry point carries `name`, equal to its directory, and a `description` built from the skill definition's routing summary and its when-to-use clause ([AC-217](../capabilities/CAP-004-ship/criteria.md)). One definition feeds both the description and the routing body, so they cannot disagree. Other agent hosts that read the Agent Skills format use the same two fields.

**Success looks like:** each generated skill's description says when to use it, and a Claude Code session lists the six skills with distinct, useful descriptions.
