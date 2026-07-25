---
id: CH-060
type: change
status: open
links: [P-006]
title: The extraction contract matches a real repository
---

# CH-060 — the extraction contract matches a real repository

## What

Edit `clue-extract`'s target contract, the canonical source at `internal/skills/source/skills/clue-extract.md.tmpl`, so it matches what Tank Royale's real extraction found missing (M-022 of [P-006](../../docs/plans/P-006-first-adoption.md)):

- Contract item on routing generalizes from `AGENTS.md` alone to assistant entry points as a class — `AGENTS.md` plus any other assistant-specific entry file a repository carries (`CLAUDE.md`, `.cursor/rules`, and similar) — keeping `AGENTS.md` as the flagship example.
- A new contract item states that a criterion born `draft` is the sanctioned way to phase extraction of a corpus too large to give every criterion its positive and negative tests in one change, rather than requiring the whole corpus tag-tested in a single PR.
- The ID-survival contract item gains concrete, deterministic minting guidance for source requirements that carry no stable existing ID at all, so re-running the same extraction reproduces the same minted IDs.

The generated skills (`.agents/skills/clue-extract/skill.md` and `internal/scaffold/templates/skills/clue-extract/skill.md`) are regenerated from the edited source with no drift.

## Why

Robocode Tank Royale's extraction (`docs/analysis/AN-001-openspec-extraction.md`) found the contract assumed a single `AGENTS.md` entry point, gave no legitimate way to phase a large corpus's extraction across changes, and had no guidance for minting IDs where the source carried none. Leaving the gap open means every future adopter re-discovers and re-solves it at the same cost.

## Scope

In scope: `internal/skills/source/skills/clue-extract.md.tmpl` (the target contract), the regenerated skill outputs, `docs/plans/P-006-first-adoption.md` bookkeeping, `CHANGELOG.md`, and a decision-log row.

Out of scope: writing the MADR mapping (M-023), any mapping-file edits, and any change to `clue validate` behavior — this is a documentation/phrasing clarification of an agent-facing contract, not a change to the deterministic judge or the core verifiable thread ([ARCH-003](../../docs/architecture/core.md)).
