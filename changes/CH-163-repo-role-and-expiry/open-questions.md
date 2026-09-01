---
id: CH-163-open-questions
type: open-questions
status: resolved
links: [CH-163, G-011, M-081]
title: Open questions for the repository role marker and analysis expiry
---

# Open questions

## 1. May this change edit `AGENTS.md` to separate the repository-local rules?

**Answered 2026-09-01 (conversation): lift the instruction and split the hub in this change.** The split landed here — the hub now states the repository's role first, keeps the shared methodology routing under one heading, and gathers the release process, generated-skill pipeline, local verification commands, guide-prose rule, and integration stricture under an explicitly repository-local heading that says none of it reaches an adopter.

M-081 includes separating this repository's local rules from the shared routing text, and the approved plan states it as work item 1.4. `AGENTS.md` is the carrier: the scaffolded template is 989 words against this repository's 1,311, and that delta is the repository-local layer sitting inside the shared routing prose with nothing marking it as local.

The conflict is that the user opened this session with an explicit instruction not to read `AGENTS.md` or `CLAUDE.md` in this repository. Approving a plan that names the split is plausibly authorization to lift that instruction for this task, but it is not the same as saying so, and the file cannot be split without reading it.

The rest of M-081 does not depend on the answer and is complete. The role marker, the source-only carrier rule, the generated skills' role-reading instruction, and the scaffolded guidance all landed without touching the hub.

**Options:** lift the instruction for this change and split the hub here; or leave the hub alone, close M-081 on what shipped, and carry the split as its own change with the instruction still standing.
