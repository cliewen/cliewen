---
id: CH-039
type: change
status: open
links: [PDR-007, PDR-011, C-012]
title: Close the review handoff
---

# CH-039 — Close the review handoff

Plan-less. This change repairs a methodology-carrier gap rather than advancing a current plan milestone.

## What

- Record how CH-037 review fixes remained outside its commits and PR, and distinguish the existing rule violation from the missing readiness tripwire.
- Define a change as ready only when every intended edit is committed, the worktree is clean, the branch is pushed, a ready hosted PR exists, and the hosted PR head equals the locally verified commit.
- Require review fixes to be committed, pushed, and reverified on the existing branch and PR; an explicitly requested local stopping point remains incomplete and is reported as not mergeable.
- Carry the rule through the constraint register, routing hubs, canonical skill sources, generated skills, PR template, tests, and user-facing release notes without adding a plain-change skill or changing PDR-011's boundary.

This is a full change because it changes methodology carriers and an agent-enforced constraint. It introduces no CLI or corpus-format API, changes no acceptance-criterion or capability meaning, and makes no semantic plan mutation.
