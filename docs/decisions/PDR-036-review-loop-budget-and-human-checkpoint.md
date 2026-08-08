---
id: PDR-036
type: decision
status: verified
links: [P-013, PDR-012, PDR-035, PDR-019, C-004, C-012, C-017]
title: The review loop runs five passes at most, and exhausting them is a human decision
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-036 — Five passes, then the human decides

## Context and problem statement

[PDR-035](PDR-035-bounded-agentic-review-loop.md) bounded the agentic review loop at three ordinary passes and kept the blocking gate intact by letting a blocking finding on the last permitted pass earn the next one. That combination is not actually bounded: as long as each pass finds something blocking, the loop continues, and the only thing the budget prevents is a *discretionary* extra pass after a clean one. A loop that may run indefinitely so long as it keeps finding defects has no upper cost, and nothing tells the human it is happening.

Three was also a guess. It came from one observed incident — a register whose prose figures made three consecutive passes find nothing but arithmetic — and that incident is the reason computed figures are advisory, not evidence about how many passes real convergence needs.

## Decision outcome

**Five passes is the ordinary budget, and exhausting it hands the decision to the human rather than to the loop.**

- **The loop stops on a pass that finds nothing blocking.** Unchanged from PDR-035, and it remains the ordinary exit: one clean first pass is sufficient and no pass is ever run to reach a quota.
- **A further pass runs only after a pass that returned at least one blocking finding.** Also unchanged. An advisory finding never earns another pass.
- **The budget is five, not three.** Three was inferred from a single incident. Five leaves room for a change that genuinely converges slowly without leaving the cost unbounded.
- **When the fifth pass still returns blocking findings, the loop stops and reports.** It tells the human what remains — the findings, not a count — and asks whether to run further passes. Only that answer extends the budget. The loop neither publishes nor continues on its own authority.
- **Nothing here permits publication with an unresolved blocking finding.** [C-012](../constraints/C-012-agents-never-merge-own-changes.md)'s boundary is untouched: a candidate with blocking findings is not ready, whether the budget is exhausted or not. Exhaustion produces a report and a question, never a publication.

The change of shape matters more than the number. Under PDR-035 the loop's own persistence decided when to stop, and a change that kept failing review consumed passes silently. Under this decision a change that cannot converge in five passes has stopped being a review problem and become a judgement about the change itself, which is the human's to make — the same reasoning that puts merge on the human's side of the boundary.

**Carrier inventory:** the amendment blockquote at the head of [PDR-035](PDR-035-bounded-agentic-review-loop.md); [C-017](../constraints/C-017-agentic-review-loop-is-bounded.md), which registers the standing rule; the canonical `internal/skills/source/skills/clue-verify.md.tmpl`, which states the operating procedure and generates the repository and scaffolded skill copies; `internal/skills/generate_test.go`, which pins the stable clauses; `guide/change-loop.md` and `CONTRIBUTING.md`, which carry the same rules for the human reading paths; and the repository and scaffolded pull-request templates, whose review line reports the pass count. [PDR-012](PDR-012-agentic-review-before-publication.md)'s rejection of a fixed pass count is untouched — a ceiling the human may raise is not a quota — and the routing hubs state no pass budget, so they are unaffected.

### Rejected: keep three and let a blocking finding always earn another pass

That is PDR-035's rule, and it is the one being replaced. Its defect is that it reads as a bound and is not one: the ceiling applies only to passes nobody needed. A rule whose stated limit never binds in the case that matters teaches its readers that the limit is decorative.

### Rejected: stop unconditionally at five and report the change as not ready

A hard stop is honest but wasteful. A change on its fifth pass with one small blocking finding left is usually one repair from done, and forcing it to become a new change moves the same work behind more bookkeeping. Asking costs one question and keeps the decision where the cost is felt.

### Rejected: pick the budget from measured convergence data

There is none. Cliewen has one recorded incident of a non-converging loop and no measurement of how many passes ordinary changes take, because the pass count only became part of the handoff in the change before this one. Five is a deliberate over-estimate chosen so the ceiling is felt rarely; when the handoff has accumulated real pass counts, that evidence can revisit the number, and revisiting it is cheap.
