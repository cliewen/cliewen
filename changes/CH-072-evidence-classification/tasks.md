---
id: CH-072-tasks
type: tasks
status: open
links: [CH-072, P-007, M-025]
title: Ordered work for evidence classification
---

# Tasks

- [ ] Inspect the current AC/test harvesters, decision records, and generated-skill sources; record any decision-blocking ambiguity in `open-questions.md`.
- [ ] Write the ADR that supersedes ADR-006's deferred per-AC annotation door and resolves the no-tag-mechanism/QS-lane decisions.
- [ ] Add acceptance criteria for classified evidence and pair enforcement, with positive and negative tests.
- [ ] Implement criteria parsing and classified evidence harvesting for Go, JVM, and Cucumber `.feature` files. (AC IDs to be minted with the criteria.)
- [ ] Update capability designs, source skills, generated skills, indexes, changelog, and P-007 bookkeeping.
- [ ] Run formatting, generation, `clue validate`, and `go test ./...`; repair any failures without weakening checks.
- [ ] Commit the implementation candidate, conduct scenario-resolution and adversarial review, and repair all actionable findings.
