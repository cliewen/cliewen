---
id: CH-106-tasks
type: change
status: open
links: [CH-106]
title: Tasks for CH-106
---

# CH-106 — Tasks

- [x] Confirm the defect against the code rather than the PR description: `checkIndexes` reads no label, `regenIndex` appends bare stems, curated rows survive regeneration.
- [x] Confirm the blast radius in adopted repositories: the same blocks ship through `clue init` and `clue scaffold`.
- [x] Establish how a counted population is reported: `ProvenanceBacklog` and `agentConstraintCount` feed the OK line in `cmd/clue/main.go`; `Issue` has no severity, so a counter never goes through `checkIndexes`.
- [x] Narrow the scope — drop verbatim-title judging and `regenIndex` normalization, record both as resolved with their reasons.
- [x] Withdraw the M-049 declaration and state the change plan-less with its reasoning.
- [ ] **Stopped for a human answer to Q1** — count or fail. It changes what the judge fails on, which is a core carrier under C-013.

Once Q1 is answered:

- [ ] Write the ADR: the emitted row format, the judge's decided behaviour, and both rejections (title grading, row normalization) with the `scaffold.go` contract as the reason for the second.
- [ ] Add the acceptance criterion for the emitted row format to the capability that owns index regeneration — confirm whether that is CAP-002 or CAP-005 before writing — with its `Test-type` and paired positive/negative evidence.
- [ ] Teach `regenIndex` to emit `- [<id> — <title>](<file>) · \`<status>\`` for missing entries, reading the target's frontmatter. Curated rows keep surviving unchanged; assert that in the negative test.
- [ ] Add the acceptance criterion for the counted population to CAP-002's criteria, with paired evidence.
- [ ] Add `corpus.IndexRowBacklog` as a sibling of `ProvenanceBacklog` and one `notes` clause in `cmd/clue/main.go`. Do not touch `checkIndexes` (C-004: no existing check is weakened or repurposed to fit this).
- [ ] Skip rows whose target is a folder README — `docs/README.md`'s block links `goals/README.md` and its siblings, which carry no artifact identity of this shape.
- [ ] Compare against parsed frontmatter fields, never raw lines: ADR-001 and ADR-002 carry YAML-quoted titles because their values contain colons, and a line-level compare reports both as false positives. Cover that case with a negative test.
- [ ] List the counted rows behind a flag, following the `--reality-gaps` precedent, since no command clears this count.
- [ ] Register the constraint with `enforcement: machine` and its `source`.
- [ ] Update every live carrier that states the index-block contract in this same change (C-006): `docs/README.md`, `docs/decisions/README.md` prose if it must state the rule, and `internal/scaffold/templates/docs/decisions/README.md`.
- [ ] `[Unreleased]` changelog entry under C-002: the emitted index row format changes for adopters and the judge reports a new population.
- [ ] Digest: regenerate indexes, confirm no existing row moved, delete this workspace.
