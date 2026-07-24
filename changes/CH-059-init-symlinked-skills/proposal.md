---
id: CH-059
type: change
status: open
links: [P-006, M-021]
title: clue init does not write through a symlinked skills folder
---

# CH-059 — clue init does not write through a symlinked skills folder

## What

Implement [P-006](../../docs/plans/P-006-first-adoption.md) M-021. Before emitting a file, `clue init` inspects the target's ancestor directories below the repository root; if any of them is a symlink, the subtree is not written. The blocked directory is reported as its own category — `linked` — with its own count in the summary line, so a shared skills tree reads differently from a file `init` refused to overwrite. A new acceptance criterion, AC-038 under CAP-001, carries the behavior with a positive and a negative test.

The direction — skip and report, never follow, no override — is recorded as a row in `docs/decisions/log.md`.

## Why

`clue init` emits every skill twice: canonically into `.agents/skills/` and mirrored into `.claude/skills/` in the Claude Code `SKILL.md` spelling ([ADR-018](../../docs/decisions/ADR-018-init-templates-embedded.md)). Sharing an assistant's skills across checkouts by making `.claude/skills` — or `.claude`, or `.agents/skills` — a symlink is a normal shape, and it is the shape a repository living with Cliewen runs into. `writeIfAbsent` probes with `os.Stat`, which follows links, so when the shared tree carries no Cliewen skill yet, `init` creates files **inside it**: a command run to initialize one repository silently mutates a directory outside it, and the report says only `created`.

Writing into a directory the user deliberately shared across checkouts is outside `init`'s mandate. Naming the skip is what turns an inexplicably empty mirror into an informed choice the user can act on.

## Decision boundary

The never-overwrite guarantee (AC-025) is untouched and is explicitly not what this change adds. `clue validate`'s own symlink handling stays as it is — `checkSkillVersions` follows links by design ([ADR-028](../../docs/decisions/ADR-028-deterministic-skill-manifest.md)). Index regeneration is unchanged: it rewrites only READMEs the user already owns. No `--force` or follow-the-link escape hatch is introduced. Detection is bounded to ancestors *below* the root; the root itself may legitimately be reached through a link and is never inspected.
