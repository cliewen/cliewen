---
id: ADR-011
type: decision
status: verified
links: [G-002, CAP-004, P-002]
title: clue and the skills are versioned — tag-stamped binary, per-skill markers, drift is a failure
author: agent
accepted-by: Flemming N. Larsen (2026-07-13, PR #6)
---

# ADR-011 — Version stamping for clue and the skills

## Context and problem statement

The binary and distributed skills need a comparable version so a released judge cannot silently run with drifted guidance.

## Decision outcome

**Stamp releases from git tags, put a `version:` marker in every skill, require the markers to agree, and make drift a failure; `dev` builds are exempt from the binary comparison.** Release tags are `vX.Y.Z`; the workflow injects the bare semver with ldflags, while checkout and pseudo-version builds report `dev`.

`clue validate` checks each skill marker, set consistency, and released-binary equality. There is no separate set file, so a copied skill carries its own stamp. The carrier is `clue version`, `corpus.checkSkillVersions`, the release workflow, and skill frontmatter. A warning-only drift check is rejected; skills outside `.agents/skills/` remain a future configuration door.
