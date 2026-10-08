---
id: PDR-066
type: decision
status: inferred
links: [CAP-010, ADR-068]
title: A sole contributor's tracked change proceeds on local allocation and says so
author: agent
accepted-by: []
---

# PDR-066 — A sole contributor's tracked change proceeds on local allocation and says so

## Context and problem statement

The change loop told an agent to stop at `clue id next`'s local-allocation warning and have a maintainer either serialize allocation or enable Git coordination. In ten first tracked changes with a seeded ledger and local allocation, nine stopped and asked, so every first change in a repository without coordination cost one maintainer answer, and the warning returns on each later change until coordination is enabled. Only a human could settle a question that, in a one-person repository, has the same answer each time.

## Decision outcome

**On the warning, the agent reads the remote with plain Git for another contributor and goes on if there is none.** It looks for any `ch-*` branch other than its own, merged or not, because the remote does not say who wrote it, and for the allocator branch, whose existence is enough because a branch listing shows no claims. With none, it continues on serialized allocation, commits and pushes the ledger with the proposal, and records in the acceptance brief's ledger note that allocation was local, what it read, and that `clue id coordinate` comes before a second contributor. It stops, as before, where such a branch or the allocator branch exists or the remote cannot be read; the answer there is to enable coordination, which ends the question.

**The check reads Git, not the forge.** A pull request needs a pushed branch, so a pull-request list adds nothing a branch listing does not show, and a repository hosted where a forge command cannot reach is treated like any other.

**The limit is stated and not closed:** a clone that has not pushed, or an identity allocated locally and not yet pushed, is invisible to the check. The rule names it, tells the agent not to present the check as proof of being alone, and keeps coordination as the step before a second contributor. `clue validate` reports a duplicate identity where two artifacts hold one, but it was not shown to catch a collision whose first change has already merged, so it is not relied on as a net.

The alternative of keeping the stop and having `clue init` say what to do was weighed and not chosen: it leaves the first change of every new repository waiting on an answer.
