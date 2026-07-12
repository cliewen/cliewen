---
id: CH-004-tasks
type: tasks
status: open
links: [CH-004]
title: Task breakdown for CH-004
---

# Tasks

- [ ] Write AN-002 (model2diagram extraction analysis) and update the analysis index
- [ ] Write ADR-008 (extraction is a generic skill), ADR-009 (AC ID namespaces), ADR-010 (provenance field) — born `inferred` — and update the decisions index
- [ ] Add AC-014…AC-018 to a new CAP-003-extract criteria.md (prefix declaration, prefix collision, JVM harvest, unknown/retired JVM reference, provenance vocabulary)
- [ ] Generalize `checkACTests` to `ac-prefix` namespaces — AC-014, AC-015
  - [ ] Harvest `@<PREFIX>-<digits>` per criteria file using its declared prefix (default `AC`)
  - [ ] Lint duplicate `ac-prefix` declarations across criteria files
  - [ ] Generalize the Go test-name grammar to declared prefixes (`TestMG010_…`)
- [ ] Add the JVM `@Tag` harvester for `*Test.kt` / `*Test.java` — AC-016, AC-017
- [ ] Add the `provenance` vocabulary check and the inferred count in `clue validate` output — AC-018
- [ ] Write positive + negative tests for AC-014…AC-018
- [ ] Write CAP-003-extract README and design.md; update capability index
- [ ] Write `.agents/skills/clue-extract/skill.md` (target contract + OpenSpec mapping)
- [ ] Update ARCH-002 (skills architecture) with clue-extract; update docs/README field/status tables for `ac-prefix` and `provenance`
- [ ] Verify: gofmt, `go test ./...` green, coverage ≥ 80%, `clue validate` clean
- [ ] Digest: delete `/changes/CH-004-extract/`, open PR, CI green
