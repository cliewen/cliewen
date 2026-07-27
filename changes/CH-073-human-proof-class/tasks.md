---
id: CH-073-tasks
type: tasks
status: open
links: [CH-073, P-007, M-026]
title: Ordered work for human-verified criteria
---

# Tasks

- [ ] Inspect `checkACTests`, ADR-032, ADR-007, the acceptance-brief template and digest, and the coverage-report idea in P-007; record any decision-blocking ambiguity in `open-questions.md`.
- [ ] Write the ADR extending ADR-032 with the `human` proof class and generalizing ADR-007's tombstone convention with a per-criterion `@draft` token.
- [ ] Add or revise acceptance criteria for the `human` class, the `@draft` exemption, and the coverage report, with positive and negative tests.
- [ ] Implement `human`-class parsing and the `@draft` tag-line exemption in `checkACTests`.
- [ ] Extend the acceptance-brief template and its digest so a newly or materially declared `Test-type: Human` criterion is named and confirmed at merge.
- [ ] Implement the derived coverage report (`covered` / `partial` / `gap` per capability), computed from corpus state.
- [ ] Update capability designs, source skills, generated skills, indexes, changelog, and P-007 bookkeeping.
- [ ] Run formatting, generation, `clue validate`, and `go test ./...`; repair any failures without weakening checks.
- [ ] Commit the implementation candidate, conduct scenario-resolution and adversarial review, and repair all actionable findings.
