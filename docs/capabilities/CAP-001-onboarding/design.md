---
id: CAP-001-design
type: design
status: draft
links: [CAP-001]
title: Design for CAP-001 onboarding
---

# Design — CAP-001 Onboarding

Stub — elaborated when `clue init` is implemented (CH-002 and onward).

Decided so far:

- `clue init` copies its scaffolding from `/.cliewen/templates`; the hand-built taxonomy created in CH-001 is the template source, not a separate artifact to keep in sync.
- `init` is idempotent on an already-initialized repo: it regenerates README index blocks (between `<!-- clue:index:start/end -->` markers) and touches nothing hand-written.
- The quickstart page lives in the repo README; the command is the guide's most important layer and must not require reading anything.

## Guide requirements (accumulating)

The guide itself is written when `clue init` exists — writing it earlier would document a tool nobody can run. Until then, every lesson about what it must contain is appended here the moment it is learned, so the guide is assembled from requirements, not from memory:

- **Prerequisites must be explicit** (learned 2026-07-12): git, the `clue` binary, and — for the PR-based change loop as practiced — an authenticated `gh` CLI. Note that the loop works with plain git and any forge; `gh` is the convenient path, not a dependency of the method.
- **Skills are guide layer 3, learned during use**: the guide links each skill at the moment the reader's first change loop needs it (delta when they branch, verify before their first PR) rather than explaining all four upfront. The system-level "why these skills" story lives in [architecture/skills.md](../../architecture/skills.md) and the guide references it, never duplicates it.
