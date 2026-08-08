---
id: CH-130-tasks
type: tasks
status: open
links: [CH-130, P-013]
title: Tasks for the M-062 statement register
---

# CH-130 tasks

- [x] Fix the segmentation rule: define what counts as a single statement, precisely enough for an independent second pass to segment the same prose the same way
- [x] Establish the evidence boundary: pin the revision, name the seven carriers read, and record which are generated from which sources
- [x] Walk `AGENTS.md` statement by statement; classify, trace, and flag duplication, checkability, and order
- [x] Walk `clue-plan` (smallest carrier; validates the segmentation rule before the large ones)
- [x] Walk `clue-analysis`
- [x] Walk `clue-upgrade`
- [x] Walk `clue-delta`
- [x] Walk `clue-verify`
- [x] Walk `clue-extract`
- [x] Aggregate the shared fragments once, and compute duplication per reading path rather than per file
- [x] Assess order per carrier: list every statement that binds absolutely but is read after the procedure it constrains
- [x] Report the population each class yields
- [x] Record compatible overlap candidates for M-063
- [x] Draft the escalations: untraceable rules, rules whose decision may have outlived its reason, and conflicting obligation pairs
- [x] Write `open-questions.md` from those escalations and stop for the human
- [x] Write `docs/analysis/AN-018-skill-statement-register.md` with the register, the populations, and the durable-form recommendation
- [ ] **Blocked on the human's answers to Q-01…Q-07.** The remaining tasks resume after they are recorded.
- [ ] Record each answer as its typed decision (PDR, constraint, or log row) per C-011
- [ ] Regenerate the analysis index and add the P-013 M-062 evidence row in the digest
- [ ] Add the `[Unreleased]` changelog entry
- [ ] Delete the change workspace in the digest commit
