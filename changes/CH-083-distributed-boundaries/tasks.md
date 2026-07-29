---
id: CH-083-tasks
type: tasks
status: open
links: [CH-083, M-036]
title: CH-083 tasks
---

# CH-083 — tasks

- [x] Establish the evidence boundary: pin the Cliewen revision, the adopter revision, the `clue` binaries used and their version stamps, the host, and the prepared-versus-clean classification of every reproduction
- [x] Scenario 1 — authorized dependent change
  - [x] Build the pinned dependent-change situation in a disposable clone: change B's work rooted on unmerged change A, never pushed
  - [x] Measure what repository-local `clue validate` reports on B's corpus, and whether anything distinguishes A's unaccepted meaning from accepted truth
  - [x] Measure what merge order does to B: base accepted, base superseded, base rejected
  - [x] Record which guarantees the branch/PR boundary supplies and which it cannot
- [x] Scenario 2 — acceptance evidence that spans repositories
  - [x] Identify the live cases in the current corpus where a criterion's carrier or proof is in another repository
  - [x] Measure what `clue validate`, `--coverage`, and the AC↔test rules do with a foreign carrier, in both directions between the two repositories
  - [x] Record what a foreign green check would and would not prove, and at which revision
- [x] Scenario 3 — external-tracker reference across a repository move
  - [x] Inventory the corpus's existing forge references and the repository namespaces they belong to
  - [x] Demonstrate the mis-resolution a bare reference already produces, and what a rename, transfer, or mirror does to it
  - [x] Measure what the stable-ID rules and `clue validate` guarantee for external references today
- [x] Write `docs/analysis/AN-013-*.md` with risk, evidence boundary, per-scenario measurements, findings, rejected options, and what the analysis does not establish
- [x] State the rejection boundary explicitly: no option that weakens the human merge boundary or makes forge state the system-of-record
- [x] End the analysis with independently routable candidates for stacked changes, cross-repository evidence, and tracker metadata, each routed by reversal cost, plus the named successor-plan consumer
- [x] Route this change's own decisions in `docs/decisions/log.md`; propose no interface
- [x] Regenerate the analysis index and update M-036's status and evidence in P-008
- [x] Verify: `clue validate --forbid-changes`, coverage and reality reports, `go test ./...`, guide build
