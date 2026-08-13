---
id: C-013
type: constraint
status: active
links: [PDR-013, PDR-019, PDR-042, C-012]
title: Changes to a core carrier require an explicit decision record and human acceptance
source: PDR-013 and PDR-019, AGENTS.md ("The core is behind a red line")
enforcement: human
---

# C-013 — Core changes need an explicit decision record

Cliewen's core is defined in [ARCH-003](../architecture/core.md): the verifiable thread from goal to acceptance evidence, human acceptance of a chosen full loop, and `clue validate` as deterministic judge. Acceptance evidence includes classified executable evidence and the acceptance brief for a genuine Human-class criterion; `clue` validates declarations and references but does not execute tests or replace human judgment. An agent recommends the full loop for a change that alters a core element, with an explicit decision record and human acceptance; if the user explicitly chooses simple instead, PDR-042's override record makes that choice visible without pretending Cliewen can veto the repository owner. Periphery never constrains the core.

The protected carrier set includes this rule's own carriers: this constraint, its defining decisions [PDR-013](../decisions/PDR-013-explicit-core-red-line.md) and [PDR-019](../decisions/PDR-019-methodology-contract-carriers-move-together.md), and the shipped AGENTS.md red-line rule. Weakening or removing the red line is itself a core-meaning change and crosses the red line.

**Residual:** whether a change alters the *meaning* of a core element. A diff touching a core carrier is detectable and a decision record accompanying it is detectable, but the pairing proves nothing in either direction: a change can touch ARCH-003 while fixing a typo, and a change can redefine what a green validate asserts without editing a single core file — by changing a check that definition depends on. [ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md) keeps the judge out of the diff regardless, and the pairing check would have been the wrong answer even with one.

The cost is that the red line holds only as far as the classification and the human's reading of it. That is deliberate rather than regrettable: the rule exists to make a human look at a core change, and a machine that pre-approved the pairing would be a way of not looking.
