---
id: C-003
type: constraint
status: active
links: [PDR-002]
title: A task marked infeasible carries its reason on the same line
source: clue-delta skill, step 2 (Propose)
enforcement: partial
---

# C-003 — `[-]` carries its reason

A task in `tasks.md` marked `[-]` is addressed-but-not-feasible and carries its reason on the same line, whenever it is marked. Digest precondition: every task is `[x]` or `[-]`-with-reason.

**Checked by:** `clue validate` ([AC-091](../capabilities/CAP-002-validate/criteria.md)) — a `[-]` line with nothing after its checkbox fails.

**Residual:** whether the prose after `[-]` is a reason at all. A machine sees that something was written and cannot judge whether it explains anything, so `[-] no` passes the check and tells a later reader nothing.

The rule this constraint used to carry as well — that a task is ticked the moment it completes and never in batch at the end — was withdrawn by CH-133. It bound every change to a discipline about a file the digest deletes, its timing half was unobservable in the file, and nothing was lost by it. What survives is the part a reader depends on: a `[-]` with no reason is indistinguishable from a task nobody finished.
