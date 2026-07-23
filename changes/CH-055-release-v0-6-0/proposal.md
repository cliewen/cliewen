---
id: CH-055
type: change
status: open
links: [P-005]
title: Close P-005 and cut v0.6.0
---

# CH-055 — Close P-005 and cut v0.6.0

## What

Close [P-005](../../docs/plans/P-005-explicit-core.md) `completed` and prepare the v0.6.0 release pair. All four milestones (M-016, M-017, M-018, M-019) already show `done` with evidence, so this change flips the plan's top-level `status`, adds a completion note in [P-004](../../docs/plans/P-004-first-try.md)'s style, and records the closure as a decision-log row per the plan's own mutation rules. It then renames `CHANGELOG.md`'s `[Unreleased]` section to `[0.6.0]` with the `### Install` block the release workflow publishes verbatim, and refreshes the guide's binary-version examples in `guide/ci-wall.md` and `guide/getting-started.md` from `0.5.1` to `0.6.0`. The skills already carry the 0.6.0 stamp from CH-053 and CH-054, so no version bump is needed there. Tagging `v0.6.0` is a human act after merge and is not part of this change.

## Why

P-005's exit criteria are satisfied and evidenced (CH-051, CH-052, CH-053, CH-054), but the plan still reads `status: active`, so `docs/plans` misstates what is in flight. The same four changes accumulated seven user-visible entries under `[Unreleased]` and bumped the skills to a 0.6.0 stamp that no release carries yet: a published 0.6.0 skill set with no 0.6.0 binary means an adopter following the guide vendors a 0.5.1 binary that rejects those skills as drift. Cutting the section and refreshing the guide's version examples closes that gap in one act, and the release workflow fails without a matching changelog section.

## Decision boundary

This change closes P-005, prepares the release pair, and records the closure decision. It does not start, scope, or imply a successor campaign; does not touch `clue`, validation semantics, skills, or scaffold templates; does not reopen or re-evidence the already-`done` milestones; and does not push the tag or publish the release. Any proposal for what comes after v0.6.0 is a separate, later change.
