---
id: CH-093-tasks
type: tasks
status: open
links: [CH-093]
title: Tasks for CH-093
---

# Tasks

- [x] Record the pinned target revision, toolchain, runtime, operating system, shell, repository state, and clean-versus-prepared evidence boundary
- [x] Run a read-only source inventory covering formats, entry points, artifact mappings, stable and minted IDs, test-purpose work, unsupported carriers, instructions, deletions, CI constraints, migration shape, and merge compatibility
- [x] Run and separately record the target repository's own relevant test or build checks and Cliewen's deterministic validation results without changing the target
- [x] Produce the PDR-020 rehearsal report in this workspace with the full AC-056 inventory — source formats and entry points, proposed mappings, preserved and minted IDs, confidence and reversal cost, test-purpose work, instruction conflicts, governance effects, planned deletions, and named plan doors — plus observed facts, inferences, unverified intent, and rejected options
- [x] Record every unresolved conflict in open-questions.md and stop before mutation; otherwise state why the report supports or narrows migration readiness
- [x] Record the human selection of the target, both pins, and the retrospective limit of a rehearsal run against an already-adopted repository
- [x] State per M-042 exit conjunct which are met and which single conjunct is unmet, and separate the limits of this already-adopted target from the milestone's exit list
- [ ] Obtain the human answer to OQ-004: whether this narrower support boundary closes the migration-readiness phase, or whether a second rehearsal against an unconverted target carrying stable source IDs is needed first
- [ ] Obtain the human answer to OQ-003, and if it directs a durable `/docs/analysis` landing, amend `proposal.md` to bring that artifact inside this change's declared scope before writing it
- [ ] Land whatever OQ-003 directs: the report as a durable artifact and OQ-001/OQ-002 as named P-009 doors or recorded decisions
- [ ] Run the complete change verification with a clean worktree and a passing `clue validate`
