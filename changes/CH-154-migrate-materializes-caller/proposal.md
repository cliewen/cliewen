---
id: CH-154
type: change
status: open
links: [CAP-001, ADR-052]
title: Migration materializes a missing thin CI caller and reports a competing validation wall
---

# CH-154 — Migration materializes a missing thin CI caller and reports a competing validation wall

## Proposal

This change is **plan-less**: no active plan owns adopter upgrade ergonomics, and inventing a milestone under P-014 or P-015 would be a fake plan item.

An adopted repository that predates the thin CI caller cannot reach the current CI shape through `clue migrate`. The caller template first shipped in v0.10.0, so every repository adopted before it has no `.github/workflows/clue.yml` by construction rather than by choice. [ADR-052](../../docs/decisions/ADR-052-missing-optional-carriers-do-not-block-safe-migrations.md) classifies that absence as a notice naming `clue init` as the materialization route, and [AC-124](../../docs/capabilities/CAP-001-onboarding/criteria.md) proves that migration neither creates nor rewrites the caller. The upgrade therefore ends with the adopter hand-writing a file whose exact bytes the binary already carries, and `clue init` is a poor substitute because it also materializes corpus stubs an established repository may deliberately not have.

A second gap compounds it. A repository whose pre-caller CI installs and runs `clue validate` in its own workflow keeps that job after the caller arrives, and migration is silent about the resulting duplicate wall. The stale job then fails work the caller was configured to treat leniently.

This change will supersede ADR-052's rejected alternative on a narrower argument, retire AC-124 and mint its successor with the reversed materialization clause, add a criterion for a competing-wall finding, and carry the reconciliation step into the upgrade workflow guidance.

## Scope boundary

Migration gains authority to create the caller from the embedded template at its default adopter choices. It gains no authority to rewrite a caller that already exists beyond the reference and version updates it already performs, and no authority to edit a repository-owned workflow: a competing validation wall is reported as a finding for a human to resolve, never rewritten. `clue init` keeps its role for a repository that has no corpus yet.
