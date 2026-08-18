---
id: ADR-022
type: decision
status: verified
links: [G-002, CAP-004, ADR-011, ADR-018]
title: Cliewen skills declare ownership in frontmatter
author: agent
accepted-by: Flemming N. Larsen (2026-07-18, planning conversation — explicitly approved the marker and non-deletion boundary)
---

# ADR-022 — Cliewen skills declare ownership in frontmatter

## Context and problem statement

The shared `.agents/skills/` directory may contain third-party skills, so location alone cannot identify the files Cliewen owns or require its version marker.

## Decision outcome

**Every distributed Cliewen skill declares `cliewen-skill: true` in YAML frontmatter.** Only marked skills participate in the version set; an absent marker is unmanaged and a non-boolean marker is malformed. Marked skills still require a string `version:`, mutual agreement, and released-binary equality under ADR-011.

The six canonical directories are reserved: an unmarked `skill.md` in one is reported as a managed skill needing migration. `clue init` remains non-destructive, creating missing files but neither overwriting nor deleting skills; a future explicit upgrade may replace marked files. The carrier is generated frontmatter, `corpus.checkSkillVersions`, and init's never-overwrite behavior. Directory-only ownership, a separate manifest, and delete-and-reinstall initialization are rejected.
