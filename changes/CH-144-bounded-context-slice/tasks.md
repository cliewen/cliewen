---
id: CH-144-tasks
type: tasks
status: open
links: [CH-144]
title: Tasks
---

# Tasks

- [ ] Open P-015 with M-070 through M-074 and its index row, so this change serves a real plan item.
- [ ] Record the bounded-slice decision as an ADR that engages PDR-034's rejected read cap and makes the frontier report a condition of the bound.
- [ ] Retire AC-053 with a tombstone and mint AC-133 for the bounded slice, its widening, and its frontier report.
- [ ] Bound the slice in `internal/corpus`: a depth option, a frontier return, and unfollowed edges reported only from included artifacts.
- [ ] Add `--depth` and `--stats` to `clue context`, printing the frontier and keeping deterministic order.
- [ ] Add focused positive and negative Unit evidence for AC-133 and retag the tests that survive as `--depth=all` behaviour.
- [ ] Update every live carrier that describes what `clue context` emits: CAP-007's README, criteria, and design, the CLI text, the routing hub, the generated skills and their scaffolded copies.
- [ ] Add the user-facing changelog entry for the bounded slice.
- [ ] Run local verification and the agentic review loop on the committed candidate, complete the M-070 digest, and prepare the ready PR handoff.
