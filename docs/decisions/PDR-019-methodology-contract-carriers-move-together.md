---
id: PDR-019
type: decision
status: verified
links: [G-001, P-008, AN-011, PDR-013, ADR-032, ADR-033, C-006, C-013]
title: Methodology contract changes update every live carrier in the same change
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-019 — Methodology contract changes update every live carrier in the same change

## Context and problem statement

Mechanical validation can be green while skills, guidance, templates, architecture, metadata, or implementation explanations still state an older methodology contract. The corpus needs a complete carrier repair without pretending a machine can infer every semantic relationship.

## Decision outcome

**A methodology contract change inventories every live carrier that states the affected contract and updates that inventory in the same change.** Live carriers include permanent corpus truth, canonical and generated skills, scaffolded copies, public and contributor guidance, templates, implementation explanations, tests, CLI text, and distribution metadata. Pinned analyses, completed plans, changelog history, and decision context remain historical unless a current note needs a pointer to a refinement.

The implementing change records the inventory, repairs all named carriers together, regenerates derived copies, and adds focused guards for stable machine-recognizable claims. The general completeness obligation remains agent-enforced; a carrier registry requires its own evidence and is not inferred here.

**Refining [PDR-013](PDR-013-explicit-core-red-line.md), the protected thread ends in acceptance evidence:** goal → plan → change → capability → acceptance criterion → acceptance evidence. Unit, Integration, E2E, and Performance criteria use classified Go, JVM, or Cucumber evidence, with `(single-direction)` as the explicit exception; Human criteria use the pull-request acceptance brief; `@draft` marks one unproven promise; legacy criteria without a proof type retain the one-reference contract. `clue validate` checks declarations and references, while normal test runners execute tests. This refinement crosses the red line and therefore requires the decision and human boundary already named by C-013.

C-006, the shared decision-record fragment in lifecycle skills, the decisions READMEs, ARCH-003, C-013, routing hubs, acceptance criteria, and content guards carry the rule and its evidence-model refinement.
