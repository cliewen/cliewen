---
id: CH-020-open-questions
type: open-questions
status: open
links: []
title: CH-020 open questions
---

# Open questions

None blocking at proposal time. Calls the agent will make as inferred decisions, reviewable in the PR:

- Templates embedded in the binary via `go:embed` (self-contained install; no network) — recorded as an ADR during implement.
- `init` emits skills to `.agents/skills/` (canonical) with a `.claude/skills/` mirror, matching the model2diagram fixture.
- AC-001 split per the granularity rule: mechanical ACs testable, 30-minute promise becomes a quality scenario. This retires an AC — flagged here so review looks at it deliberately.
