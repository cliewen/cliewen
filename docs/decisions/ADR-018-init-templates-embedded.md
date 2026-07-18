---
id: ADR-018
type: decision
status: verified
links: [P-002, CAP-001, ADR-011, ADR-013, ADR-021]
title: The init scaffolding is embedded in the clue binary
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# ADR-018 — Init templates embed in the binary

> **Skill-copy consequence refined by [ADR-021](ADR-021-generated-standalone-skills.md):** the canonical authoring sources now generate both skill output trees; the embedded tree and binary-delivery decision remain unchanged.

## Context and problem statement

`clue init` materializes the foundation every adopting repo starts from: the `/docs` taxonomy, the AGENTS.md routing hub, the skills, and a CI workflow template ([CAP-001](../capabilities/CAP-001-onboarding/README.md)). That content has to travel with the binary somehow, and the delivery mechanism decides whether onboarding needs a network, credentials, or a checkout of this repo — and whether the scaffolding can drift from the binary that emits it.

## Considered options

1. **Fetch from this repo at init time** — always current, but requires network access and (while the repo is private) credentials; a user's first command can fail for reasons that have nothing to do with them, and the fetched scaffolding may not match the binary's version.
2. **Ship templates as a separate release asset** — no network at init time, but two artifacts to download and keep matched; the pair-version discipline ADR-011 built for skills would need a third leg.
3. **Embed the templates in the binary via `go:embed`** — the installed binary is self-contained; scaffolding and binary version as one unit, exactly like the skill stamps.

## Decision outcome

**Option 3.** The template tree lives at `internal/scaffold/templates/` and is compiled in via `go:embed`. Two placement facts follow from the toolchain, not from preference: the Go tool ignores directories whose names start with `.` or `_`, so neither the originally sketched `/.cliewen/templates` nor the canonical `.agents/skills` can be embedded directly. Consequences:

- The **skills are duplicated** into the template tree as generated distribution artifacts; a Sanity test holds both trees to their shared canonical render, so drift between the authored sources, canonical skills, and what `init` emits fails the build.
- The `.github/` workflow template is stored under a `github/` path and mapped to its dotted target at emit time.
- The CI template's `CLUE_VERSION` pin is substituted from the embedded skills' version stamp — the pair version (ADR-011) has a single carrier.
- `init` emits the skills twice in the target repo: `.agents/skills/` as the canonical location and a `.claude/skills/` mirror in the Claude Code spelling (`SKILL.md`), matching the layout proven in the first adopter repo.
- Template changes reach users only through a release — the same doctrine as the skills: released binaries and their scaffolding never disagree.
