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

- `clue init` copies its scaffolding from `/.cliewen/templates`; the
  hand-built taxonomy created in CH-001 is the template source, not a
  separate artifact to keep in sync.
- `init` is idempotent on an already-initialized repo: it regenerates
  README index blocks (between `<!-- clue:index:start/end -->` markers)
  and touches nothing hand-written.
- The quickstart page lives in the repo README; the command is the
  guide's most important layer and must not require reading anything.
