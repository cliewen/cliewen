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
- [ ] Run the complete change verification and update P-009 M-042 bookkeeping in the digest, recording per exit item which are met (report-only rehearsal with no target mutation, stable-ID and minted-ID inventory, mixed JVM per-executable evidence sample, unsupported carriers, governance and deletion effects, CI-runner/action/install constraints, current-to-target corpus migration, merge-mode compatibility, separated target-versus-Cliewen results) and which are recorded as a narrower boundary (the proposed wall running in the target CI shape, blocked on OQ-001; a passing target test runner, blocked on OQ-002)
- [ ] Delete this whole `/changes/CH-093-migration-rehearsal/` workspace in the digest commit before review, because `main` never contains `/changes/` and the rehearsal is transient branch-local evidence under PDR-020 — only after the rehearsal handoff is authorized
