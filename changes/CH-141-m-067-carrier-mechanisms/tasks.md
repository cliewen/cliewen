# Tasks

- [ ] Verify the hub `clue-extract` row / generator test is already green (no new work; report as already-landed in the digest)
- [ ] Fix the five stale constraint `source:` fields (C-001, C-002, C-004, C-005, C-013) to point at their current live carriers
- [ ] Tighten `checkConstraints` so an ID-shaped `source:` segment must resolve to a live corpus artifact; add tests
- [ ] Write the no-computed-figures rule into constraint `C-018` (`enforcement: agent`, promotion trigger); update `docs/constraints/README.md`
- [ ] Write the reference-names-what-it-is-about rule into constraint `C-019` (`enforcement: human`, priced residual); state it in `clue-verify`'s agentic review loop finding shape
- [ ] Repair `clue-verify`'s unbounded "assessed against every constraint" checklist item
- [ ] Write the F-RB-09 decision (log row) and mint constraint `C-020` (`enforcement: human`)
- [ ] Move the orient-after-merge statement from `review-boundary` to `durable-work`; add `durable-work` to `clue-extract`'s render
- [ ] State PDR-032's mid-change triage rule in `durable-work` and in `clue-delta`'s Propose/Digest steps
- [ ] State PDR-006's rejected-record clause in the shared `decision-records` fragment
- [ ] Regenerate skills (`go generate ./internal/skills`) and taxonomy indexes (`clue scaffold`)
- [ ] Update P-013's M-067 row (`doing` → `done`) and CHANGELOG if user-visible
- [ ] Run local verification and the agentic review loop; mark PR ready
