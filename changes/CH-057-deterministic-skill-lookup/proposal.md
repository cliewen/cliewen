---
id: CH-057
type: change
status: open
links: [P-006]
title: One verdict on every filesystem — resolve a skill's manifest by case-folded name
---

# CH-057 — One verdict on every filesystem

## What

Make `clue validate` locate a skill's manifest the same way on every filesystem. Today `checkSkillVersions` reads each skill folder's manifest at the fixed lowercase path `skill.md` (`internal/corpus/skillversions.go`). On a case-insensitive filesystem (Windows, default macOS) that call still finds a manifest named `SKILL.md`; on a case-sensitive one (a Linux CI runner) it does not, so the skill drops out of the managed set. The same corpus and the same binary then reach two different verdicts depending on where the check runs.

This change resolves the manifest by scanning each skill directory's entries for a name equal to `skill.md` under Unicode case folding, so `SKILL.md`, `Skill.md`, and `skill.md` are all found identically everywhere. If a single directory holds more than one case-variant of the name — only reachable on a case-sensitive filesystem — the validator reports that as a named ambiguity issue rather than silently choosing one, keeping the verdict deterministic in that corner too.

Because `clue validate` is Cliewen's deterministic judge, part of the core ([ARCH-003](../../docs/architecture/core.md)), the resolved behavior is recorded as a decision record under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) (ADR-028) and carried by a new acceptance criterion (AC-034) with a positive and a negative test that both fail against today's fixed-path lookup.

This serves [P-006](../../docs/plans/P-006-first-adoption.md) milestone **M-020**.

## Why

The first repository to live with Cliewen, Robocode Tank Royale, vendored skills whose manifests follow the Anthropic/Claude Code convention of `SKILL.md`. On the maintainer's macOS checkout `clue validate` saw those skills; on the Linux CI runner it did not, so the wall that is supposed to give one answer gave two. A judge whose conclusion depends on the host filesystem cannot carry a verifiable thread ([G-001](../../docs/goals/G-001-verifiable-thread.md)), and the divergence surfaces exactly in the skill-version drift check that [G-002](../../docs/goals/G-002-versioned-clue-and-skills.md) relies on.

Case-folded lookup fixes the divergence without forcing adopters to rename convention-named manifests, and the ambiguity report keeps the one filesystem that *can* hold two variants from reintroducing a silent, order-dependent choice.

## Decision boundary

This change touches only how a skill's manifest file is located and read. It does not change the ownership marker, the legacy-slot handling, the version-stamp rule, set consistency, or drift comparison — every existing acceptance criterion (AC-029…AC-033) keeps its meaning and tests. It renames nothing on disk and asks no adopter to rename anything. It does not alter the meaning of the verifiable thread or the merge boundary, and it stays inside `internal/corpus`; no other milestone of P-006 is touched.
