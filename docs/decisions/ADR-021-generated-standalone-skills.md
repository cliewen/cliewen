---
id: ADR-021
type: decision
status: verified
links: [G-002, CAP-004, P-002, ADR-011, ADR-018]
title: Skills are generated as standalone artifacts from shared canonical sources
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, conversation)
---

# ADR-021 — Generated standalone skills

## Context and problem statement

Cliewen's skills must be independently installable, yet several skills carry the same cross-cutting methodology rules. Authoring those rules separately makes their consistency depend on every future edit finding every copy. The embedded `clue init` templates add another output tree that must remain identical to the installed skill set.

## Considered options

1. **Keep independently authored standalone skills** — preserves portability and named entry points, but shared rules can drift.
2. **Collapse the workflows into one composite skill with references** — gives shared rules one home, but removes the separate lifecycle entry points and loads unrelated workflow guidance.
3. **Reference another skill or a shared directory at runtime** — reduces copied text, but an individually installed skill is incomplete and depends on a repository layout outside its own folder.
4. **Generate standalone skills from skill-specific templates and shared source fragments.**

## Decision outcome

**Option 4.** The authored source of the Cliewen skill set lives under `internal/skills/source/`. A deterministic repository generator composes each skill-specific template with shared fragments and writes the complete, independently usable artifacts to `.agents/skills/` and `internal/scaffold/templates/skills/`.

- Shared instructions are authored once; skill-specific workflow instructions remain separate.
- Generated skills contain no runtime include, inheritance, symlink, or dependency on another skill. Copying one skill folder preserves its complete instructions and version stamp.
- The five public skill names and lifecycle entry points remain unchanged.
- Every generator-owned file in both output trees is checked against the canonical rendering. A missing, changed, or unexpected file within an owned skill directory fails the build.
- Per-skill version markers remain in every generated artifact, preserving ADR-011's portable version surface while making the marker itself single-source at authoring time.
- The generated embedded tree replaces ADR-018's manually synchronized skill-copy consequence; embedding the generated output in the binary remains unchanged.

**Carrier:** `internal/skills` owns the source fragments, composition manifest, generator, and drift tests; `.agents/skills/` and `internal/scaffold/templates/skills/` are generated distribution artifacts.
