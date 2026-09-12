---
id: G-015
type: goal
status: proposed
links: [G-001, VIS-001]
title: What the work teaches is preserved and found again
---

# G-015 — What the work teaches is preserved and found again

**Who wants it:** the maintainer and every agent that starts a task in a repository someone has already solved this problem in (2026-09-12, from [AN-024](../analysis/AN-024-methodology-review-acceptance-and-learning.md)).

**Why:** a repository accumulates non-obvious operational knowledge that its corpus has nowhere to put. Which command actually works on this platform, which documented flag form fails, what the recovery is when a release check rejects a candidate — these are neither product meaning nor architecture nor a future-shaping decision, so today they are learned during a change and lost when it ends. The next session rediscovers them, and an agent rediscovers them every session by construction, because a fresh context is the normal case rather than the exception.

The corpus is designed to be shared memory between humans and agents, and this is the class of memory it drops. [ADR-026](../decisions/ADR-026-adopter-types-default-lifecycle.md) already permits an adopter to add a type such as a runbook, and the design overview already asks changes to keep useful documentation current, but nothing establishes that a discovery gets noticed, given a home, found again, or corrected when it goes stale.

The risk runs the other way too, which is why this is a goal rather than an instruction. A capture obligation with no eligibility bar produces a folder of near-duplicate notes that nobody reads and an agent must read all of, which costs more than the rediscovery it was meant to prevent. Writing down a workaround can also make a workaround permanent when the honest fix was to remove the confusing step.

**Success looks like:**

- A reusable discovery made during a change reaches an appropriate home within that same change, and the next relevant task finds it without reading everything.
- Reuse that exposes a problem corrects or retires the guidance, so what the repository states stays true.
- An existing home is corrected in preference to adding a new document, and removing or automating a confusing step is preferred to documenting a workaround.
- Guidance that records a recurring task states its trigger, prerequisites, procedure, expected result, and recovery, so a reader can tell whether it applies before following it.
- Nothing is written when nothing useful was learned, and the method does not reward volume.
- The trial is an agent with no prior session finding the guidance, completing the task, and correcting it when conditions changed — not the existence of the document.
