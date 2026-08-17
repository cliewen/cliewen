---
id: C-020
type: constraint
status: active
links: [P-013, PDR-047]
title: An agent orients on the plan after a human reports a merge
source: PDR-047, the shared durable-work fragment
enforcement: human
---

# C-020 — An agent orients after a human reports a merge

After a human reports that a Cliewen change's pull request merged, the agent orients before starting anything else: it describes the plan's next unfinished step in plain language and asks whether to start it, or says that the plan has nothing left and asks what comes next. Treating "merged" as a go signal and silently starting the next task is the failure this rule exists to prevent.

**Residual:** all of it. Nothing observes a conversation, so no machine can tell whether an agent oriented, asked, or simply proceeded; a human who was not asked has to notice the omission themselves. The cost of that gap is a plan whose actual next step drifts from what the human believes is in flight.
