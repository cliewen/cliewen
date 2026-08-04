---
id: C-003
type: constraint
status: active
links: []
title: Tasks tick the moment they complete; a skipped task carries its reason
source: clue-delta skill, step 2
enforcement: partial
---

# C-003 — Tasks tick immediately; `[-]` carries its reason

A task in `tasks.md` is marked `[x]` the moment it completes — never in batch at the end; the unticked list is what shows what is actually left. `[-]` marks addressed-but-not-feasible and must carry its reason on the same line. Digest precondition: every task is `[x]` or `[-]`-with-reason.

**Checked by:** `clue validate` ([AC-091](../capabilities/CAP-002-validate/criteria.md)) — a `[-]` line with nothing after its checkbox fails.

**Residual:** "immediately", and whether the prose after `[-]` is a reason at all. A batch of ticks applied at the end is indistinguishable, in the file, from ticks applied as the work happened; only whoever did the work knows. The cost is that the unticked list stops being a live account of what is left and becomes a formality completed at the end, which is the failure this constraint was written against.
