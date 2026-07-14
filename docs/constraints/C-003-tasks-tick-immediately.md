---
id: C-003
type: constraint
status: active
links: []
title: Tasks tick the moment they complete; a skipped task carries its reason
source: clue-delta skill, step 2
enforcement: agent
---

# C-003 — Tasks tick immediately; `[-]` carries its reason

A task in `tasks.md` is marked `[x]` the moment it completes — never in batch at the end; the unticked list is what shows what is actually left. `[-]` marks addressed-but-not-feasible and must carry its reason on the same line. Digest precondition: every task is `[x]` or `[-]`-with-reason.

**Promotion trigger:** `clue validate` lints `tasks.md` under `/changes/`: a `[-]` line without prose after the checkbox fails — then `enforcement: machine` for the reason rule (the "immediately" half stays behavioral and agent-held).
