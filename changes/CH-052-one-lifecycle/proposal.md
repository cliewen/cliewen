---
id: CH-052
type: change
status: open
links: [P-005, M-017, G-001, C-013]
title: One default status lifecycle, and adopter-defined types validate against it
---

# CH-052 — One status lifecycle, adopter types tolerated

## What

Collapse the validator's fourteen per-type status vocabularies to one default lifecycle `draft → active → retired` plus a short list of semantically necessary exceptions, and let artifact types the validator does not recognize — adopter extensions — validate against that same default instead of erroring "unknown type". Mints ADR-025 (the default lifecycle) and ADR-026 (adopter-type tolerance). Flips the architecture and analysis records that used `verified`-in-status to `active`, and promotes ARCH-003 `core.md` from `draft` to `active`. Rewrites the status tables in `docs/README.md` and the `clue init` template as "default plus exceptions".

## Why

Eleven of the fourteen vocabularies were near-duplicates of one `draft → active → retired` lifecycle; only goal, plan, decision, and the transient types carry a genuinely different shape. Stating one default plus the real exceptions is smaller to hold and smaller to teach. Tolerating unknown types is the mechanical half of PDR-013's extension story: adopters must be able to add their own artifact types under `docs/` without the core's permission, and a closed-world type registry blocked exactly that.

## Decision boundary

This change alters what a green `clue validate` asserts (the status field's allowed values and the removal of the unknown-type rejection), so by the red line ([C-013](../../docs/constraints/C-013-core-changes-need-decision.md)) it is a full change carrying explicit decision records and taking effect only at human merge. ADR-026 must show the unknown-type change is a re-scoping, not a weakening ([C-004](../../docs/constraints/C-004-never-weaken-checks.md)): core fields, unique IDs, resolvable links, and a valid default status are all still enforced on adopter types. No skill template carries a lifecycle table, so no skill regeneration or version bump is needed — confirmed by a generate drift check.
