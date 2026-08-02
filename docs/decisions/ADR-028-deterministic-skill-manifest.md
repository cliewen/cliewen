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

`corpus.checkSkillVersions` locates each skill's manifest at the fixed lowercase path `skill.md`. On a case-insensitive filesystem (Windows, default macOS) the read still succeeds when the file on disk is named `SKILL.md`; on a case-sensitive filesystem (a Linux CI runner) it fails and the skill silently drops out of the managed version set. A skill vendored under the Anthropic/Claude Code convention of an uppercase `SKILL.md` is therefore visible to `clue validate` on the maintainer's laptop and invisible in CI — the same corpus and binary reaching two verdicts on the same commit.

`clue validate` is Cliewen's deterministic judge ([ARCH-003](../architecture/core.md)); a conclusion that depends on the host filesystem is not deterministic and cannot carry the verifiable thread ([G-001](../goals/G-001-verifiable-thread.md)). Changing how the judge decides which files it reads alters what a green validate asserts, so the change crosses the core red line and needs this record ([C-013](../constraints/C-013-core-changes-need-decision.md)).

## Decision outcome

**A skill's manifest is the directory entry whose name equals `skill.md` under Unicode case folding, resolved by scanning the skill directory rather than reading a fixed path. When one directory holds more than one case-variant of that name, validation reports a named ambiguity issue instead of choosing one.**

- **Case-folded match, not a fixed path.** Each skill directory is scanned for an entry whose name case-folds to `skill.md` and which resolves to a regular file — a symlink to one counts, since a skills tree shared across checkouts is a supported shape, while a directory, FIFO, or device by that name is not a manifest and is never opened. `skill.md`, `SKILL.md`, and `Skill.md` resolve identically on every filesystem, so a convention-named manifest yields the same enrollment everywhere. A directory with no such entry is not a skill, unchanged from before.
- **Multiple variants are an ambiguity, not a silent pick.** Only a case-sensitive filesystem can hold two entries that both fold to `skill.md`. Rather than depend on directory-read order, validation names that directory and reports the ambiguity, keeping the verdict deterministic in the one place the divergence could re-enter.
- **The ambiguity is reported whoever owns the skill.** Resolution precedes ownership: the marker that would scope the directory out of Cliewen's set ([ADR-022](ADR-022-skill-ownership-marker.md)) lives inside the very file that cannot be chosen, so an unmarked third-party directory holding two case-variants is reported like any other. This is a deliberate exception to "an unmarked skill does not participate" (AC-029, AC-033), and the only one: it costs an adopter a named, one-rename fix, where staying silent would mean either picking a variant by read order or letting the ambiguity decide which skills Cliewen can see.
- **Nothing downstream of resolution changes.** Ownership marking ([ADR-022](ADR-022-skill-ownership-marker.md)), legacy-slot migration, the version-stamp requirement, set consistency, and binary drift all keep their meaning and operate on the manifest once resolved. This decision governs only which file is read.
- **Adopters rename nothing.** The rule accepts the manifest name adopters already ship; it does not impose a canonical spelling on disk.

**Carrier:** `corpus.checkSkillVersions` manifest resolution (machine); AC-037 under [CAP-004](../capabilities/CAP-004-ship/criteria.md) with a positive and a negative test, the negative case covering an unmarked third-party directory alongside a Cliewen one.

### Rejected: keep the fixed lowercase path and require adopters to rename

Renaming `SKILL.md` to `skill.md` in every adopting repository pushes Cliewen's platform accident onto adopters and fights an established assistant-skill convention. The judge, not the adopter, owns cross-filesystem determinism.

### Rejected: lowercase-compare but pick the first entry on collision

Choosing one of two case-variants by read order reintroduces the exact non-determinism this decision removes — a different entry could win on a different filesystem or listing. A collision is rare and always a real defect in the checkout, so naming it is the honest verdict.
