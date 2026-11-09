---
id: G-016
type: goal
status: accepted
links: [G-001, VIS-001]
title: Human acceptance is an informed decision about usable behaviour
---

# G-016 — Human acceptance is an informed decision about usable behaviour

> Amended by [PDR-073](../decisions/PDR-073-acceptance-needs-the-decision-not-a-brief.md): there is no acceptance brief. The human is told in plain words what changed, what was checked and what remains uncertain, and Human-class observations are recorded in the change proposal.

**Who wants it:** the human at the acceptance boundary of a tracked Cliewen change, in this repository and in every adopter running the loop (2026-09-12, from [AN-024](../analysis/AN-024-methodology-review-acceptance-and-learning.md)).

**Why:** the method places a human at the acceptance boundary. That establishes who decides and when. It cannot establish that the decision was informed, because authorising a change and witnessing its behaviour are different acts.

Two gaps follow. A `Human`-class criterion's proof is a claim that someone observed something, and is worth little without a statement of who, under what conditions, and what they saw. And a candidate can be complete, green, and structurally valid while leaving the human no way to tell whether the changed behaviour is usable by the person it was built for. A green candidate that should be rejected is the case the boundary exists for, and rejecting one should not mean repeating the review.

Where a goal depends on someone adopting or operating the result, the gap widens: a guide existing is not evidence that an unfamiliar maintainer can follow it, and a passing test suite is not evidence that the failure path is recoverable.

**Success looks like:**

- The handoff to the human states, in a few plain sentences, what behaviour changed, what was observed, and which consequential uncertainty still needs human judgement.
- `Human` evidence identifies who observed what, under which conditions, with what result, and what was not observed.
- Where a goal depends on adoption or operation, the relevant task is demonstrated — an unfamiliar maintainer performing it and recovering from a documented failure — rather than inferred from the existence of documentation.
- A human can explain an acceptance or a rejection from the handoff, the proposal and its evidence, including rejecting an inadequate green candidate without repeating the code review.
- Informing the decision costs no separate document: more evidence must not mean a longer ritual before deciding.
