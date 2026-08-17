---
id: CAP-001
type: capability
status: active
links: [G-001]
title: Onboarding — install to green validate in under 30 minutes
goal: G-001
---

# CAP-001 — Onboarding

## What

A new user goes from installing `clue` to their first green `clue validate` in under 30 minutes.

## Why

The method's first enforced requirement is its own accessibility — instant usability beats conceptual superiority (the Spec-Kit lesson). Serves [G-001](../../goals/G-001-verifiable-thread.md): a thread nobody can pick up enforces nothing.

The layered guide, layers kept strictly separate:

1. **Command (seconds):** `clue init` materializes the whole convention — `/docs` taxonomy, README indexes, skills, CI workflow — in one call; `clue migrate` later upgrades a released corpus through a previewed, safe plan. A legacy decision log is the semantic exception: MIG-010 inventories every row and blocks every write until a reviewed full change classifies future-shaping choices by subject, accounts for narrative, repairs references, and removes the log.
2. **Quickstart (5 minutes):** one page — install, `clue init`, first change loop, watch `validate` go green.
3. **Skills** — learned during use.
4. **Book** — the why; depth, secondary.

`clue init` exists since CH-020; how the command and the guide realize the layers is [design.md](design.md), the mechanical path is held by the tested criteria in [criteria.md](criteria.md), and the 30-minute end-to-end promise is [C-015](../../constraints/C-015-onboarding-under-30-minutes.md).

What init emits for a vendor is bounded to an inert pointer at the hub ([PDR-022](../../decisions/PDR-022-vendor-entry-points-only-point.md)); nothing executable is ever emitted for an assistant. Discovery of a newer release therefore runs through the tool itself, which reports one line from any workflow command, with the hub carrying the instruction for sessions that run none ([PDR-023](../../decisions/PDR-023-tool-notice-and-hub-instruction.md)). Both files are the adopter's, and `clue migrate` reports rather than repairs either one.

The emitted CI caller delegates its validation wall to Cliewen's immutable upstream reusable workflow while keeping runner and binary-acquisition choices local to the adopter ([ADR-038](../../decisions/ADR-038-upstream-validation-workflow.md)).
