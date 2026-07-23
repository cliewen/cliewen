---
id: CH-056
type: change
status: open
links: []
title: Open P-006 to digest the first foreign adoption
---

# CH-056 — Open P-006 to digest the first foreign adoption

## What

Create [P-006](../../docs/plans/P-006-first-adoption.md) `active`, a campaign that digests the adoption friction Robocode Tank Royale recorded when it adopted Cliewen (its `docs/analysis/AN-001-openspec-extraction.md`, merged as its PR #218 on 2026-07-19). The plan carries four milestones, M-020…M-023, continuing corpus-global milestone numbering from P-005. This change writes the plan file, its index row, and a decision-log row recording the campaign's opening — nothing else. It is plan-less itself: no plan can serve a change whose product is the plan, the same shape [P-004](../../docs/plans/P-004-first-try.md) was opened in.

Two of Tank Royale's eight findings arrived already closed and are recorded as such in the plan rather than scheduled: the skill-version drift check now scopes to skills marked `cliewen-skill: true` (ADR-022), and `clue init` ships in released binaries. The remaining six are the campaign.

## Why

Tank Royale is the first repository outside this maintainer's own toolchain to run Cliewen end to end, and it recorded its friction in the form the methodology asks for: a findings document naming the upstream candidates. That evidence has been sitting undigested since the adoption merged, and it is the only evidence of its kind — foreign-soil trials ([P-002](../../docs/plans/P-002-leaves-home.md), M-007) produced findings from the outside, but no repository had yet lived with the corpus. Leaving it unscheduled wastes the one signal that shows where Cliewen costs an adopter time, and each finding rediscovered by the next adopter costs that adopter the same hours.

One finding also reaches the core. The validator reads a skill only as lowercase `skill.md`, so on a case-insensitive filesystem it sees skills that are invisible to it on a Linux CI runner: the same corpus, the same binary, two verdicts. `clue validate` is Cliewen's deterministic judge ([ARCH-003](../../docs/architecture/core.md)), and a judge whose answer depends on the filesystem is not deterministic. That is why the campaign leads with it.

## Decision boundary

This change opens the campaign and fixes its milestone set; it implements no finding. It makes no change to `clue`, the validator, the skills, the scaffold templates, or the guide, and it does not promote, re-evidence, or reopen anything in P-005. Whether Tank Royale itself upgrades, and any change in that repository, stays outside this corpus. The two closed findings are recorded as history, not reopened. Each milestone is a separate full change, proposed and merged on its own.
