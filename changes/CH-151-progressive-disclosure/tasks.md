---
id: CH-151-tasks
type: tasks
status: open
links: [CH-151]
title: Tasks for CH-151
---

# Tasks

- [ ] Record ADR-059's standalone-directory and progressive-disclosure decision, including rejected runtime sharing and optional-reference alternatives.
- [ ] Retire AC-028 with a tombstone and add AC-137 for deterministic standalone skill directories with routed references.
- [ ] Refactor the canonical skill sources and generator so each `skill.md` is a short router and every deferred rule is emitted inside that skill's `references/` directory (AC-137).
- [ ] Add focused positive and negative Unit evidence for routed references, complete standalone directories, byte-identical trees, and missing, changed, or unexpected generated files (AC-137).
- [ ] Regenerate `.agents/skills/` and `internal/scaffold/templates/skills/` exclusively from `internal/skills/source/`.
- [ ] Update every live carrier that describes generated skill completeness or loading, including CAP-004, ARCH-002, contributor/public guidance, and release notes where the shipped surface changes.

