---
id: ADR-018
type: decision
status: verified
links: [P-002, CAP-001, ADR-011, ADR-013, ADR-021, ADR-038]
title: The init scaffolding is embedded in the clue binary
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, PR #20 review conversation)
---

# ADR-018 — Init templates embed in the binary

> **Skill-copy consequence refined by [ADR-021](ADR-021-generated-standalone-skills.md):** the canonical authoring sources now generate both skill output trees; the embedded tree and binary-delivery decision remain unchanged.

## Context and problem statement

`clue init` must materialize the taxonomy, routing hub, skills, and CI caller without network access or drift between the installed binary and its templates.

## Decision outcome

**Embed `internal/scaffold/templates/` with `go:embed`.** The binary is self-contained and template changes reach adopters only through a release. Go's embedding rules require dot-prefixed targets to live under ordinary source paths; the `.github` target is mapped at emit time.

The skills are generated into the embedded tree and checked against the canonical render. Init emits the `.agents/skills` tree, a `.claude/skills` mirror, and `CLAUDE.md` pointing at `AGENTS.md`; it substitutes the CI version and source reference from build metadata, falling back to the release tag. This read-only source-checkout lookup never inspects the target repository. Fetching templates at init or shipping a separate asset is rejected because it adds network/credential or version-pair failure modes.
