---
id: CAP-011-design
type: design
status: active
links: [CAP-011, PDR-058]
title: Design for repository orientation
---

# Design — Repository orientation

## Derived position

`corpus.OpenChanges`, the existing plan milestone parser, and `corpus.ProposedGoals` provide four ordered populations to `clue next`: open changes, active milestones, draft milestones, and proposed goals. Paths and table rows provide stable ordering. An open change wins because it records work already underway; an active `doing` milestone precedes `todo`; draft milestones and proposed goals remain choices for human review.

The default command prints the strongest available choice and counts relevant alternatives. `--all` exposes every population. An empty result points back to the goal inbox. The command reads the same local corpus scan as `clue context`, performs no Git inspection, and writes nothing.

## Fresh-context handoff

The cross-agent `AGENTS.md` hub is the session-start carrier. Its first tool call remains the release check. The orientation read follows once, before substantive work, and is conditionally disclosed: an open-ended request gets a brief position and recommendation, while a concrete request is not delayed by unrelated project-management prose. A recommendation follows a bounded `clue context` read and cannot promote a proposed goal, activate a draft plan, claim a milestone, or start a change.

No assistant hook or settings file is emitted. Nothing identifies a Git user or saves a focus outside the corpus, so a stale private preference cannot override repository truth.
