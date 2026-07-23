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

1. **Command (seconds):** `clue init` materializes the whole convention — `/docs` taxonomy, README indexes, skills, CI workflow — in one call.
2. **Quickstart (5 minutes):** one page — install, `clue init`, first change loop, watch `validate` go green.
3. **Skills** — learned during use.
4. **Book** — the why; depth, secondary.

`clue init` exists since CH-020; how the command and the guide realize the layers is [design.md](design.md), the mechanical path is held by the tested criteria in [criteria.md](criteria.md), and the 30-minute end-to-end promise is [C-015](../../constraints/C-015-onboarding-under-30-minutes.md).
