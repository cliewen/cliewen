---
id: C-017
type: constraint
status: active
links: [P-013, PDR-012, PDR-035, C-004]
title: The agentic review loop owns severity and stops within its bounded ordinary budget
source: PDR-035
enforcement: human
---

# C-017 — The agentic review loop is bounded

Given a Cliewen change entering automatic agentic review, when a caller briefs the reviewer and review passes run, then the loop's own blocking/advisory model governs the verdict; computed counts and arithmetic disagreements remain advisory while wrong, missing, or reused identities remain blocking; the reviewer does not spend a pass re-deriving figures; and the loop stops on the first pass with no blocking findings.

Three passes are the ordinary budget. A fourth or later pass runs only when the immediately preceding pass returned at least one blocking finding. An advisory may be repaired before a pass already required by a blocking repair; one first reported by a pass with no blocking findings stays open for a later change so publication remains bound to the exact reviewed commit. The verification handoff reports the number of passes run together with the review mode, reviewed commit, and advisory findings left open.

**Residual:** no repository machine can observe the caller's brief, a delegated reviewer's classification, or the number and sequence of agent turns. The human at the merge boundary judges the reported review evidence against this rule. If that judgment fails, prompt wording can turn advisory work into a false gate or keep a clean candidate in an unbounded loop; treating the budget as a hard stop can instead let a genuinely blocking defect escape its required repair pass.
