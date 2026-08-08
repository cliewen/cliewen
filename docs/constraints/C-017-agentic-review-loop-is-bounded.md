---
id: C-017
type: constraint
status: active
links: [P-013, PDR-012, PDR-035, PDR-036, C-004, C-012]
title: The agentic review loop owns severity and stops within its bounded ordinary budget
source: PDR-035, PDR-036
enforcement: human
---

# C-017 — The agentic review loop is bounded

Given a Cliewen change entering automatic agentic review, when a caller briefs the reviewer and review passes run, then the loop's own blocking/advisory model governs the verdict; computed counts and arithmetic disagreements remain advisory while wrong, missing, or reused identities remain blocking; the reviewer does not spend a pass re-deriving figures; and the loop stops on the first pass with no blocking findings.

Five passes are the ordinary budget. A further pass runs only when the immediately preceding pass returned at least one blocking finding. When the fifth pass still returns blocking findings, the loop stops, reports those findings to the human, and asks whether to run further passes; only that answer extends the budget, and exhaustion never permits publication. A blocking finding is actionable and enters the unresolved hosted-conversation lifecycle; an advisory is a non-actionable observation for the publication gate and stays in the verification handoff. An advisory may be repaired before a pass already required by a blocking repair; one first reported by a pass with no blocking findings stays in the handoff for a later change so publication remains bound to the exact reviewed commit without making the advisory a merge gate. The verification handoff reports the number of passes run together with the review mode, reviewed commit, and advisory findings left open.

**Residual:** no repository machine can observe the caller's brief, a delegated reviewer's classification, or the number and sequence of agent turns. The human at the merge boundary judges the reported review evidence against this rule. If that judgment fails, prompt wording can turn advisory work into a false gate or keep a clean candidate in an unbounded loop; and an agent that treats the exhausted budget as permission to publish, rather than as a report and a question, defeats the boundary the budget was written to protect.
