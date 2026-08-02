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
- [x] Q1 answered by the human: the judge counts and does not fail. Recorded as ADR-041.
- [x] Write [ADR-041](../../docs/decisions/ADR-041-index-rows-state-their-record.md): the emitted row format, the counted population, and both rejections with the `scaffold.go` contract as the reason for the second.
- [x] Add AC-073 to CAP-005's criteria with `Test-type: Unit` — the appended row states its record.
- [x] Teach `regenIndex` to emit `- [<id> — <title>](<file>) · \`<status>\`` for missing entries via the shared `corpus.RowIdentity`. Curated rows survive unchanged; asserted in the negative test.
- [x] Add AC-074 to CAP-002's criteria with `Test-type: Unit` — filler rows are counted and listed, never failed on.
- [x] Add `corpus.IndexRowBacklog` as a sibling of `ProvenanceBacklog`, plus the `--index-rows` flag and one `notes` clause in `cmd/clue/main.go`. `checkIndexes` keeps its behaviour; the only edit to it lifts its inline taxonomy-README predicate into `isTaxonomyReadme` so the judge and the count read one definition rather than two.
- [x] Skip rows whose target is a subfolder README, and rows carrying more than one link.
- [x] Compare against parsed frontmatter, never raw lines. The AC-073 positive test uses a YAML-quoted title containing a colon, which is the case that produced false positives when this was checked by hand.
- [x] Paired evidence in the declared class for both criteria: `TestAC073_UnitPositive_…` / `TestAC073_UnitNegative_…` in `internal/scaffold/regen_test.go`, `TestAC074_UnitPositive_…` / `TestAC074_UnitNegative_…` in `cmd/clue/main_test.go`.
- [x] Update `TestUnit_CrlfReadmeKeepsItsLineEndings`, whose expected string ended at the closing link paren and so encoded the old row format. Its purpose is line endings; the assertion is now the full new row plus `\r\n`, equally strict and still proving CRLF (C-004: the check is not weakened, and no other test was touched to make this pass).
- [x] Register [C-016](../../docs/constraints/C-016-index-rows-state-their-record.md) with `enforcement: machine` and its `source`.
- [x] Repair this repository's own three filler rows — `AN-009` in the analysis index, `ARCH-003` in the architecture index, `G-003` in the goals index. CH-105 repaired the decisions index only; the new count found the rest, which is the count doing its job on its first run.
- [x] Update every live carrier stating the index contract in this same change (C-006): CAP-005's README and design, CAP-002's design, CAP-001's design. `architecture.md` names the command without claiming a row format and needs no edit; the analysis records are pinned history.
- [x] `[Unreleased]` changelog entry under C-002 — the emitted row format and the new counted population are both user-visible.
- [x] Verify: `go build`, `go vet`, `gofmt -l`, `go test ./...`, coverage 82.8% against C-014's 80% floor, and `clue validate` reporting 0 index rows on this corpus.
- [x] Digest: regenerate indexes (`clue scaffold` appended ADR-041 and C-016 in the new format, which is the change proving itself), confirm no pre-existing row moved, delete this workspace.
