---
id: C-015
type: constraint
status: active
links: [G-001, CAP-001, ADR-027]
title: A new user reaches their first green validate in under 30 minutes
source: ADR-027, CAP-001 onboarding capability
enforcement: human
---

# C-015 — Onboarding under 30 minutes

Given a machine with git and no Cliewen tooling, when a new user follows the README quickstart — install `clue`, run `clue init`, run `clue validate` — then they reach a green validate in under **30 minutes**, reading nothing beyond the quickstart.

This is the successor of retired QS-002 ([ADR-027](../decisions/ADR-027-quality-scenarios-are-constraints.md)), itself the successor of retired AC-001 (CH-020): the clock spans a human journey (reading, installing) that no focused test pair can verify, so it lives here as the bar the quickstart is written against, while the mechanical path is held by AC-002/AC-024/AC-025. Every change that touches the quickstart or `clue init` is assessed against it: anything that adds a prerequisite, a manual edit, or a new required read before the first green validate violates it. Instant usability beats conceptual superiority — the method's first enforced requirement is its own accessibility ([CAP-001](../capabilities/CAP-001-onboarding/README.md)).

**Enforcement:** `human` — no test can time a human's first run; a reviewer judges whether a change to the quickstart or `clue init` keeps the journey under the bar.
