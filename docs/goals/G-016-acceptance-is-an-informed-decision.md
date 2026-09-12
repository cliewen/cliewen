---
id: G-016
type: goal
status: proposed
links: [G-001, VIS-001]
title: Human acceptance is an informed decision about usable behaviour
---

# G-016 — Human acceptance is an informed decision about usable behaviour

**Who wants it:** the human at the merge gate of a full Cliewen change, in this repository and in every adopter running the loop (2026-09-12, from [AN-024](../analysis/AN-024-methodology-review-acceptance-and-learning.md)).

**Why:** the method places a human at the acceptance boundary and gives them an acceptance brief to decide from. That establishes who decides and when. It cannot establish that the decision was informed, because authorising a change and witnessing its behaviour are different acts and the current brief does not reliably distinguish them.

Two gaps follow. A `Human`-class criterion is proved by its line in the brief, which means the proof is a claim that someone observed something — with no statement of who, under what conditions, or what they saw. And a brief can be complete, green, and structurally valid while leaving the human no way to tell whether the changed behaviour is usable by the person it was built for. A green candidate that should be rejected is the case the boundary exists for, and rejecting one currently means repeating the review the brief was supposed to summarise.

Where a goal depends on someone adopting or operating the result, the gap widens: a guide existing is not evidence that an unfamiliar maintainer can follow it, and a passing test suite is not evidence that the failure path is recoverable.

**Success looks like:**

- The brief states what behaviour changed, what was observed, and which consequential uncertainty still needs human judgement.
- `Human` evidence identifies who observed what, under which conditions, and with what result.
- Where a goal depends on adoption or operation, the relevant task is demonstrated — an unfamiliar maintainer performing it and recovering from a documented failure — rather than inferred from the existence of documentation.
- A human can explain an acceptance or a rejection from the brief and its evidence alone, including rejecting an inadequate green candidate without repeating the code review.
- The brief stays one screen. More evidence must not mean a longer document.
