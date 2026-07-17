---
id: CH-020-tasks
type: tasks
status: open
links: []
title: CH-020 task breakdown
---

# Tasks

- [ ] Criteria first: revise `CAP-001/criteria.md` — retire AC-001 (tombstone, `@retired`), add testable ACs for the mechanical path (init on empty repo → validate green; idempotent re-run; refusal to overwrite hand-written files; broken-link failure naming file and ID), move the 30-minute end-to-end promise to a quality scenario
- [ ] Build `/.cliewen/templates` from the verified fixture: docs taxonomy READMEs + starter artifacts, AGENTS.md routing hub, CI workflow template (model2diagram's vendored pattern, generalized)
- [ ] Implement `clue init` in Go: embed templates (`go:embed`), emit taxonomy + AGENTS.md + skills (`.agents/skills/` + `.claude/skills/` mirror, stamped with binary version) + CI template; idempotent re-run regenerates only index blocks; never overwrite existing files, report skips
- [ ] Test pairs for every new AC (positive + negative), driving the built binary in a temp repo; purposes per the taxonomy (`AC<digits>` prefix in Go)
- [ ] Quickstart in README: prerequisites (git, clue, gh-as-convenient-path), install → init → first loop → green validate; skills linked at point of need
- [ ] Document-type introduction in the generated `docs/README.md` + folder READMEs (plain-language what/when for ADR, PDR, log row, capability folder, architecture, constraints, quality)
- [ ] Digest prep: CAP-001 README/criteria/design to `active`, design.md elaborated from stub (template source, embed decision, idempotence, skip behavior), architecture docs updated (`init` joins the CLI surface), decisions recorded (embed-templates ADR; log rows for smaller calls), plan bookkeeping M-005, changelog entry
