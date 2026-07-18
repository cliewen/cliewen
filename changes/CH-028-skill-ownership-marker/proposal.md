---
id: CH-028
type: change
status: open
links: [G-002, CAP-004]
title: Cliewen skills declare their ownership
---

# CH-028 — Cliewen skills declare their ownership

## What

Add a machine-readable `cliewen-skill: true` frontmatter marker to every generated Cliewen skill and scope version/drift validation to marked skills. Unmarked third-party skills under `.agents/skills/` coexist without being forced into Cliewen's version set; a present but invalid marker fails loudly instead of silently escaping validation.

Keep `clue init` non-destructive: it continues to create missing files and skip existing files. This change adds no uninstall or upgrade command and deletes no installed skill. Future removal is necessary only if the managed skill topology changes through a rename, move, replacement, or retirement.

## Why

The current validator treats every `.agents/skills/*/skill.md` as a Cliewen skill. That makes unrelated skills fail Cliewen's missing-version, set-consistency, or binary-drift checks and gives `clue` no machine-readable basis for deciding which files it owns. A frontmatter ownership marker makes that boundary explicit while preserving the existing per-skill version contract.

This change is explicitly plan-less. It refines the active CAP-004 shipping capability after P-002 completed; none of P-003's goes-public milestones claims this behavior.

## Decision boundary

ADR-022 records the durable ownership boundary approved by the maintainer. This change does not create an adopter-facing updater, overwrite existing files through `clue init`, or delete any installed skill.
