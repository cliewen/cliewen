---
id: CH-031-tasks
type: tasks
status: open
links: [CH-031]
title: Tasks for CH-031
---

# Tasks

- [x] Add AC-034 (BOM in a corpus file fails loudly) and AC-035 (second frontmatter fence in an artifact body fails loudly) to CAP-002 criteria
- [x] Implement the BOM and duplicate-frontmatter checks in `internal/corpus` and wire them into `Validate` (AC-034, AC-035)
- [x] Add positive and negative tests for both checks (AC-034, AC-035)
- [x] Extend the `clue-extract` target contract with the single-frontmatter conversion rule and regenerate the skills
- [x] Add the decision-log row and the user-visible changelog entry
- [ ] Run the full pre-merge verification suite
