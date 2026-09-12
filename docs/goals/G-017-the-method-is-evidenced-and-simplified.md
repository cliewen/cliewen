---
id: G-017
type: goal
status: proposed
links: [G-001, VIS-001]
title: The method's value is evidenced, and the method is simplified against that evidence
---

# G-017 — The method's value is evidenced, and the method is simplified against that evidence

**Who wants it:** the maintainer deciding which obligations Cliewen should keep, and anyone weighing whether to adopt it (2026-09-12, from [AN-024](../analysis/AN-024-methodology-review-acceptance-and-learning.md)).

**Why:** Cliewen keeps adding obligations and has no evidence base for removing any. Each one was added for a reason that was good at the time, and the method's own record of whether it earned its cost afterwards is dogfooding plus two adopters, which cannot establish general effectiveness. The [draft vision](../vision.md) already says the uncertainty about scale and sustained maintenance is real.

The measurements that are easy here are the wrong ones. Artifact counts, green checks, and validator passes measure compliance with the method, not the outcome the method promises, and optimising them is how a method becomes ceremony. What matters is whether people misunderstood something consequential, how much acceptance actually cost, what got rediscovered, which documentation misled, and which obligations were bypassed and why — the last being the most informative signal available, because a bypassed rule is a rule someone judged not worth its cost in a real situation.

Two distinctions have to survive the evaluation. Contract impact and operational risk are not the same axis: unchanged-contract work can still break something badly, so verification depth should follow consequence rather than route. And practical task learning is not a product outcome: a working report export does not establish that anyone's reconciliation effort fell.

**Success looks like:**

- Real work is followed through planning, acceptance, and subsequent use, preferably with willing adopters rather than only here.
- Consequential misunderstandings, acceptance effort, rediscovery, misleading documentation, and bypassed obligations with their reasons are observed and recorded.
- Populations, conditions, and sampling are defined before any quantitative claim, and learning is measured rather than document counts or green checks.
- Each obligation is assessed by the failure it prevents and whether it can replace existing ceremony rather than add to it, and evidence is what justifies keeping, revising, or removing it.
- A goal whose benefit is uncertain names an observable benefit, its key assumption, and an owner and trigger for revisiting it.
- Observation starts alongside the first campaign rather than after the others finish.
