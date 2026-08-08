---
id: PDR-035
type: decision
status: inferred
links: [P-013, PDR-012, PDR-016, PDR-019, PDR-029, C-004, C-012, C-017, ADR-021]
title: The agentic review loop owns its severity model and has a bounded ordinary budget
author: agent
accepted-by: []
---

# PDR-035 — The agentic review loop owns severity and budget

> **The budget and its exhaustion are replaced by [PDR-036](PDR-036-review-loop-budget-and-human-checkpoint.md):** the ordinary budget is **five** passes, not three, and when the fifth pass still returns blocking findings the loop stops, reports what remains to the human, and asks whether to run further passes — only that answer extends the budget. The "last permitted pass earns the next one" rule below is withdrawn, because it left the loop unbounded whenever it kept finding defects. Everything else this record decides is unchanged: the loop owns its blocking/advisory classification, a caller's brief cannot redefine severity, computed figures are advisory while identities stay blocking, "no blocking findings" is the exit condition, the exact-reviewed-commit boundary holds, and the handoff reports the pass count.

## Context and problem statement

[PDR-012](PDR-012-agentic-review-before-publication.md) requires an adversarial review loop before publication and rejects a fixed number of ritual passes, but its exit language says both that the loop ends on “no actionable findings” and, through its live carrier, that advisory findings do not gate publication. It also leaves a reviewer brief able to redefine severity, demand exhaustive checking of advisory material, and keep a correct candidate in repeated review. How does the loop retain a real blocking gate without allowing its cost or verdict to be rewritten by each invocation?

## Decision outcome

**The agentic review loop owns its blocking/advisory classification, and an ordinary review has a three-pass budget.**

- **The caller cannot redefine severity.** A brief may identify risks and intent, but asking for a different severity model or an exhaustive sweep of a class the loop calls advisory does not change the verdict or publication gate. Blocking means the change breaks the corpus, gets behaviour wrong, leaves a false normative claim in a live carrier, or ships a criterion its evidence does not hold; advisory material never becomes blocking because the caller labels it so.
- **Computed figures are advisory.** A finding whose substance is a count, total, population figure, or arithmetic disagreement is advisory regardless of the brief. A wrong, missing, or reused identity remains blocking: allocating an identity is not counting. The reviewer spends no pass re-deriving figures; the author remains responsible for computing and checking them before publication.
- **Three passes are the ordinary budget, not a ritual minimum.** The loop stops as soon as the current commit receives a pass with no blocking findings. A fourth or later pass runs only when the immediately preceding pass returned at least one blocking finding; when that pass returned none, the loop is over. This does not reinstate the fixed-pass alternative PDR-012 rejected: one clean first pass is sufficient, and nobody runs three merely to reach a quota.
- **“No blocking findings” is the exit condition without weakening the exact-commit boundary.** This settles PDR-012's older “no actionable findings” wording in favour of the severity model already carried by `clue-verify`. Under the hosted finding lifecycle, a blocking finding is actionable and becomes a repair-required unresolved conversation; an advisory is a non-actionable observation for the publication gate and stays in the verification handoff rather than becoming such a conversation. An advisory repair may ride before a pass already required by a blocking repair. An advisory first reported by a pass with no blocking findings stays visible in that handoff for a later change; editing the clean commit would create a new candidate that [PDR-016](PDR-016-pr-state-carries-agent-handoffs.md) and [C-012](../constraints/C-012-agents-never-merge-own-changes.md) require to be reviewed, while this decision says the bounded loop is over. The advisory is durably deferred, not silently repaired and not promoted into a gate.
- **The handoff reports cost.** Verification reports the review mode, reviewed commit, number of passes run, and any advisory findings left open, so the human merge gate can see both the verdict and what convergence cost.

The bounded budget does not permit publication with an unresolved blocking finding. If the last permitted ordinary pass finds one, its repair earns the next pass; the exceptional pass exists because the gate found a defect, not because a caller demanded more inspection.

**Carrier inventory:** [PDR-012](PDR-012-agentic-review-before-publication.md) carries the amended decision; [PDR-016](PDR-016-pr-state-carries-agent-handoffs.md), [C-012](../constraints/C-012-agents-never-merge-own-changes.md), and CAP-006's AC-040 keep their actionable-finding and exact-reviewed-commit rules, with “actionable” now supplied by this decision's blocking classification rather than widened to advisories; [C-017](../constraints/C-017-agentic-review-loop-is-bounded.md) registers the standing human-enforced rule; the canonical `internal/skills/source/skills/clue-verify.md.tmpl` states the operating procedure and generates the repository and scaffolded copies; the shared `internal/skills/source/shared/review-boundary.md.tmpl` states the advisory-deferral consequence in every skill that carries the exact-commit boundary; `internal/skills/generate_test.go` pins their stable clauses; `guide/change-loop.md` and `CONTRIBUTING.md` carry the same operating and handoff rules for their own reading paths; and the repository and scaffolded pull-request templates require the review mode, reviewed commit, pass count, and advisory findings left open, with the template assertion in `cmd/clue/main_test.go` holding that field. `clue-delta` and `clue-extract` invoke the review loop outside the generated review-boundary fragment, while the repository and scaffolded routing hubs direct light and full changes to `clue-delta` without stating the loop's severity or pass budget.

### Rejected: let each reviewer brief choose its own severity model

That makes publication depend on prompt wording rather than repository methodology. The same candidate could be clean under one caller and blocked under another, and advisory arithmetic could consume the same repair loop as a false normative claim.

### Rejected: stop unconditionally after three passes

A hard stop would make the budget override the gate and could publish a candidate whose third pass found a blocking defect. The exceptional-pass rule keeps the gate intact while preventing a clean pass from being followed by discretionary extra review.

### Rejected: make computed figures blocking when they appear in measured evidence

Authors still owe correct measured evidence, but making the reviewer re-derive it spends review capacity on arithmetic and recreates the non-converging repair cycle this decision closes. A figure can reveal a separate blocking defect; the defect is reported on its operative requirement, not promoted because its symptom is a number.
