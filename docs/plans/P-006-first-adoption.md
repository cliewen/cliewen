---
id: P-006
type: plan
status: active
links: [G-001, G-002]
title: Cliewen digests its first adoption — one verdict everywhere, an honest init, an extraction contract that survives a real repository
---

# P-006 — Cliewen digests its first adoption

P-002 trialed the methodology on foreign soil without changing those repositories; P-003 and P-004 made the public path reachable and worth a first try; P-005 stated the core. This campaign spends the first evidence of a repository that actually lives with Cliewen: Robocode Tank Royale adopted it through `clue-extract` and recorded eight upstream findings in its own extraction report (`docs/analysis/AN-001-openspec-extraction.md`, merged 2026-07-19 — the repository is public, so citing it here is allowed). The campaign converts what is still open into fixes: one verdict on every filesystem, an `init` that does not fight a real skills folder, and an extraction contract that matches what a real repository looks like. It serves [G-001](../goals/G-001-verifiable-thread.md) — a judge that answers differently per platform cannot carry a verifiable thread — and [G-002](../goals/G-002-versioned-clue-and-skills.md), whose drift check is where the divergence shows.

Two of the eight findings arrived already closed and are recorded here as history, not scheduled: the skill-version drift check no longer demands one stamp across every installed skill but scopes to skills marked `cliewen-skill: true` ([ADR-022](../decisions/ADR-022-skill-ownership-marker.md)), and `clue init` now ships in released binaries, so the layout no longer has to be hand-materialized. The remaining six findings are the four milestones below.

The campaign is sequential. Each milestone is a separate full Cliewen change rooted at accepted `main`, published as a ready pull request, and accepted by human merge before the next begins. Milestone numbering is corpus-global and continues from P-005 (M-016…M-019). M-020 leads because it reaches the core; M-023 follows M-022 because a source mapping is written against the contract it serves.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-020 | **One verdict on every filesystem**: `clue validate` reaches the same conclusion about an adopted repository's skills whether the checkout sits on a case-preserving or a case-insensitive filesystem — today it locates a skill only as lowercase `skill.md`, so a skill visible to it on Windows or macOS is invisible on a Linux runner; the resolved behavior is stated as a decision record under [C-013](../constraints/C-013-core-changes-need-decision.md) because the deterministic judge is core ([ARCH-003](../architecture/core.md)), carried by an acceptance criterion with a positive and a negative test that fail without the fix; `clue validate` and `go test ./...` green | `todo` | |
| M-021 | **`clue init` does not fight a symlinked skills folder**: initialization into a repository whose `.claude/skills` is a symlink — the usual shape when an assistant's skills are shared across checkouts — detects the link and skips the mirror instead of writing through it, and reports that skip; `init`'s existing guarantee that no pre-existing file is overwritten is unchanged and is not what this milestone adds; the behavior is covered by an acceptance criterion with a positive and a negative test; `clue validate` and `go test ./...` green | `todo` | |
| M-022 | **The extraction contract matches a real repository**: `clue-extract` names assistant entry points as a class rather than `AGENTS.md` alone, states born-`draft` criteria as the sanctioned phasing lever for a corpus too large to tag tests in one change, and gives ID-minting guidance for sources whose requirements carry no stable IDs; the source edits land in `internal/skills/source/` and the generated skills are regenerated with no drift; strict guide build green where the guide restates any of it | `todo` | |
| M-023 | **MADR is a mapping, not improvisation**: `clue-extract` ships `mappings/madr.md` covering status-vocabulary conversion, `accepted-by` for acceptance that predates Cliewen, and ID preservation, sufficient for a MADR corpus to be converted mechanically; the mapping is written against the M-022 contract and cites the Tank Royale conversion as its worked case; skills regenerated with no drift | `todo` | |

## Explicitly out of this campaign

Any change inside Tank Royale, model2diagram, or another adopting repository, including whether they upgrade; reopening the two findings that arrived closed, or re-litigating [ADR-022](../decisions/ADR-022-skill-ownership-marker.md)'s ownership scoping; folding a `/dot-*` principles system into the methodology; automated migration of an already-extracted corpus; source formats other than MADR; package-manager distribution channels; and changes to the meaning of the verifiable thread or the merge boundary.

## Mutation rules (lintable)

Status and evidence fields in the milestone table may mutate in an implementing change's merge digest. Everything else in this file changes only through a declared plan revision backed by a decision record routed by reversal cost ([C-011](../constraints/C-011-decision-records-typed.md)). Plan adjustments are decisions.
