---
id: CH-136-tasks
type: tasks
status: open
links: [CH-136]
title: Tasks for CH-136
---

# Tasks — CH-136

- [x] Write PDR-040 (push is durability, ready is the explicit act), amend PDR-007, PDR-012, PDR-016, and PDR-033 with blockquote notes, and update the decisions index (serves AC-131)
- [x] Rewrite the review-boundary fragment and the `clue-verify`/`clue-delta` templates to the new model, run `go generate ./internal/skills`, and update the generator and scaffold guards so the new sentences bind (serves AC-131)
- [x] Gate the strict digest and acceptance-brief checks on non-draft in `ci.yml`, `clue-validation.yml`, and the scaffolded caller, keeping the merge gate intact (serves AC-131)
- [x] Update C-012, CAP-006 (README, criteria with new AC-131, design), both PR templates, CONTRIBUTING.md, and guide/change-loop.md (serves AC-131)
- [x] Add the CHANGELOG entry and update the identity ledger
