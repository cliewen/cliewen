---
id: CH-018
type: change
status: open
links: []
title: Post-merge orientation — a merged PR is acceptance, not a go signal
---

# CH-018 — Post-merge orientation

## What

Extend `clue-delta` step 5 with a post-acceptance rule: when the human reports that the PR was merged, the agent does not silently start the next task. It first orients the human — if the plan this change served has steps left, it names the next one and describes in plain language what it is about (IDs are pointers, never the description), then asks the human whether to start it; if the plan has nothing left, it says so and asks the human what to do next. Record the decision as a decision-log row and add a CHANGELOG entry for adopters.

## Why

Observed while running multi-PR plans in the adopter repo (model2diagram): the moment the human reports a merge, the agent treats it as a go signal and continues with the next task, leaving the human without any picture of where the plan stands or what is about to happen. The existing boundary text already says follow-up work needs explicit human scoping, but nothing obliges the agent to *give the human the context needed to scope it* — the plan's current state and the next step, in words rather than bare IDs. This change closes that gap on the acceptance side of the loop, mirroring how the review boundary closed it on the merge side.

## Plan item

Plan-less, explicitly: a methodology correction from adopter feedback, same category as CH-017. P-002/M-007 (foreign soil) does not apply — model2diagram shares a maintainer with this repo.
