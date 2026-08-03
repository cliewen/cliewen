---
id: C-006
type: constraint
status: active
links: [PDR-019]
title: Decision records are timeless prose; a method contract moves every live carrier together
source: docs/decisions/README.md
enforcement: agent
---

# C-006 — Decision records are timeless prose; a method contract moves every live carrier together

A record's context states the problem, not the episode: a motivating incident earns at most one sentence. This holds for ADRs and PDRs alike. A decision that changes a methodology contract must inventory every live carrier that states the affected contract — including the `clue` rule (machine), skill text (agent), init template (default), public guidance, current corpus truth, implementation explanation, or distribution metadata that applies — and update that inventory in the same change. A method decision without its complete current carrier set has not reached everyone it instructs.

The timeless and same-change carrier rules are restated by every workflow that writes decision records: the six `clue-*` skills carry them through the shared `decision-records` fragment in `internal/skills/source/shared/`, a single authoring point that the generator and drift tests hold identical across both distributed skill trees; the repository and scaffolded decisions READMEs carry separate human-facing copies. Stable content guards hold the evidence-model claims repaired with PDR-019, but no lint derives the complete carrier set for an arbitrary methodology contract, so that general completeness judgment remains agent-enforced.

**Promotion trigger:** timelessness is meaning and stays human-reviewed; the carrier half promotes when `clue` can derive an affected methodology contract's live carrier set and reject a change that leaves any member stale — then `enforcement: machine` for that detectable subset.
