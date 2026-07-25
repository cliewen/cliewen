---
id: CH-062
type: change
status: open
links: [P-006]
title: Close P-006 and cut v0.7.0
---

# CH-062 — Close P-006 and cut v0.7.0

## What

Close [P-006](../../docs/plans/P-006-first-adoption.md) `completed` and prepare the v0.7.0 release pair. All four milestones (M-020, M-021, M-022, M-023) already show `done` with evidence, so this change flips the plan's top-level `status`, adds a completion note in [P-005](../../docs/plans/P-005-explicit-core.md)'s style, refreshes the plans index entry, and records the closure as a decision-log row per the plan's own mutation rules. It then bumps the shared skill frontmatter stamp from `0.6.0` to `0.7.0` and regenerates the five skills into `.agents/skills/` and `internal/scaffold/templates/skills/`, renames `CHANGELOG.md`'s `[Unreleased]` section to `[0.7.0]` with the `### Install` block the release workflow publishes verbatim, and refreshes the guide's binary-version examples in `guide/ci-wall.md` and `guide/getting-started.md` from `0.6.0` to `0.7.0`. Tagging `v0.7.0` is a human act after merge and is not part of this change.

## Why

P-006's exit criteria are satisfied and evidenced (CH-057, CH-059, CH-060, CH-061), but the plan still reads `status: active`, so `docs/plans` misstates what is in flight. The same campaign accumulated four user-visible entries under `[Unreleased]`, two of which change what the shipped skills say: `clue-extract` gained the MADR mapping and the corrected target contract, and every skill now states the `accepted-by:` boundary. Those skills still carry the `0.6.0` stamp that the last release already published, so an adopter vendoring skills from `Latest` cannot tell the two skill sets apart and the drift rule cannot either. Bumping the stamp, cutting the changelog section, and refreshing the guide's version examples closes that gap in one act, and the release workflow fails without a matching changelog section.

## Decision boundary

This change closes P-006, prepares the release pair, and records the closure decision. It does not start, scope, or imply a successor campaign; does not change `clue`, validation semantics, or the meaning of any skill text beyond its version stamp; does not reopen or re-evidence the already-`done` milestones; and does not push the tag or publish the release. Any proposal for what comes after v0.7.0 is a separate, later change.
