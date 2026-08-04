---
id: C-006
type: constraint
status: active
links: [PDR-019]
title: Decision records are timeless prose; a method contract moves every live carrier together
source: docs/decisions/README.md
enforcement: human
---

# C-006 — Decision records are timeless prose; a method contract moves every live carrier together

A record's context states the problem, not the episode: a motivating incident earns at most one sentence. This holds for ADRs and PDRs alike. A decision that changes a methodology contract must inventory every live carrier that states the affected contract — including the `clue` rule (machine), skill text (agent), init template (default), public guidance, current corpus truth, implementation explanation, or distribution metadata that applies — and update that inventory in the same change. A method decision without its complete current carrier set has not reached everyone it instructs.

The timeless and same-change carrier rules are restated by every workflow that writes decision records: the six `clue-*` skills carry them through the shared `decision-records` fragment in `internal/skills/source/shared/`, a single authoring point that the generator and drift tests hold identical across both distributed skill trees; the repository and scaffolded decisions READMEs carry separate human-facing copies. Stable content guards hold the evidence-model claims repaired with PDR-019, but no lint derives the complete carrier set for an arbitrary methodology contract, so that general completeness judgment cannot be mechanized at all.

**Residual:** both halves. Timelessness is a judgment about prose — whether a sentence states an enduring rule or narrates the week it was written. Carrier completeness is worse: nothing derives the set of live carriers that state a given methodology contract, so no check can report the one that was missed. The focused content guards this repository adds bind the specific claims they name and nothing more; presenting them as proof of a complete inventory would be the mistake [PDR-019](../decisions/PDR-019-methodology-contract-carriers-move-together.md) warns against.

The cost is drift that stays invisible until someone reads the stale carrier and believes it — which, for a skill or a template, means an adopter follows a rule this project no longer holds. A mechanism could exist, but it needs a machine-readable statement of which carrier makes which claim, and that is a decision with its own evidence to gather rather than a check to write.
