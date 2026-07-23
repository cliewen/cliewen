# The verifiable thread

Cliewen organizes system knowledge as a graph with one red thread from motivation to executable evidence.

```mermaid
flowchart TD
    G[Goal<br/>who needs what and why] --> P[Plan<br/>a bounded campaign]
    P --> CH[Change<br/>the transient proposal]
    CH --> CAP[Capability<br/>what the system can do]
    CAP --> AC[Acceptance criterion<br/>one verifiable behavior]
    AC --> T[Test<br/>positive and negative evidence]
    C[Constraints<br/>including verifiable quality bars] -. laws checked against every Cliewen change .-> CH
```

## Goal

A goal states who wants an outcome and why. Proposed goals form the inbox; accepting a goal says it is real, not that it must be built immediately.

## Plan

A plan is a finite campaign serving a goal. Its milestones have explicit exit criteria and evidence. Completed plans are frozen rather than rewritten, so the plan index also records what the project has achieved.

## Change

Cliewen does not own every repository edit. Three rules set the tier, by how deeply a change reaches into meaning; the first that matches wins. A change is **plain** when nothing about meaning changes: it affects no product behavior, intent, evidence, decision, plan, policy, or methodology, so it uses an ordinary branch, checks relevant to the edited surface, a pull request, and human merge, without a CH number or corpus work. It is **light** when meaning is touched but not changed: no decision, no acceptance meaning change, no semantic plan mutation, no methodology carrier touched. Everything else is **full**.

A Cliewen change is a branch-sized proposal. A full change uses a transient workspace under `/changes/CH-xxx-*` for the proposal, ordered tasks, and blocking questions. The workspace is deleted during the digest because the current system truth belongs in `/docs`, while Git keeps the proposal history. A light change skips that workspace and its ready pull-request description becomes the proposal, but the branch and human merge boundary remain.

Two guards hold above the rules. When the tier is unclear, take the higher one; and the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move to the full loop before continuing.

## Capability and acceptance criterion

A capability owns three views: a plain-language explanation, Gherkin acceptance criteria, and implementer-facing design. An active criterion has a stable ID and both positive and negative focused tests. If its meaning changes, the old ID is retired as a tombstone and a new one is minted.

That immutability matters. A test tagged `AC-042` should always mean the same promise, even years later.

## Constraints

Constraints are rules a Cliewen change must not break: a law, license, policy, project convention, or a verifiable quality bar such as a coverage floor or a maximum onboarding time. Each one names its source and whether a machine, agent, or human enforces it, and every Cliewen proposal is assessed against all of them.

## Four actors, one boundary

Skills carry process knowledge, `clue` is the deterministic judge, protected CI is the wall, and the human controls acceptance. The machine does not pretend to understand whether a criterion is valuable; the human does not have to repeat a locally completed code review, but the agent can never perform the merge that accepts its own work. CI becomes a wall only when its PR check is required and branch protection blocks integration without it.

## Next

[See where the durable artifacts live in the corpus.](./corpus)
