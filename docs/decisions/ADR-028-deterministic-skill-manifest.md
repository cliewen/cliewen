---
id: ADR-028
type: decision
status: verified
links: [P-006, G-001, C-013, ARCH-003, CAP-004, ADR-022]
title: A skill's manifest is resolved by case-folded name, so the judge reaches one verdict on every filesystem
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-028 — Deterministic skill-manifest resolution

## Context and problem statement

A fixed lowercase `skill.md` path can make an uppercase `SKILL.md` visible on case-insensitive hosts and invisible on Linux, giving the deterministic judge different verdicts for the same tree.

## Decision outcome

**Resolve the manifest by scanning each skill directory for a regular entry whose name case-folds to `skill.md`.** `skill.md`, `SKILL.md`, and other case variants therefore enroll identically; a symlink to a regular file counts, while a directory or special file does not. If more than one case variant exists, validation reports a named ambiguity rather than choosing by directory order.

Resolution precedes ownership, so ambiguity is reported even for an unmarked third-party directory; after resolution ADR-022 ownership, version, and drift rules are unchanged. The carrier is `corpus.checkSkillVersions` and AC-037. Requiring adopters to rename files or picking the first collision are rejected because they move a filesystem accident to adopters or reintroduce nondeterminism.
