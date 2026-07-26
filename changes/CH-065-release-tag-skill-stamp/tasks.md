---
id: CH-065-tasks
type: tasks
status: open
links: [CH-065]
title: Tasks for CH-065
---

# Tasks

- [ ] Add the gate to `.github/workflows/release.yml`: run `clue validate --forbid-changes` from source, stamped with the tag's version, after the release notes are extracted and before the build
- [ ] Rewrite the workflow header comment so it points at the gate instead of asserting the agreement it never checked
- [ ] Add `TestSanity_ReleaseRunsTheJudgeStampedAsTheTag` to `cmd/clue/main_test.go`, covering the stamped invocation and its position ahead of the build and the publish step
- [ ] Prove the gate empirically: a mismatched stamp fails and names both versions, a matching stamp passes
- [ ] Record the release pipeline's gate in `docs/capabilities/CAP-004-ship/design.md`
- [ ] Record the workflow-not-corpus-rule choice as a decision-log row
- [ ] Add the `[Unreleased]` changelog entry for the fix
- [ ] Run repository and corpus verification
