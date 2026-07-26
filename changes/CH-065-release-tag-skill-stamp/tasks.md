---
id: CH-065-tasks
type: tasks
status: open
links: [CH-065]
title: Tasks for CH-065
---

# Tasks

- [x] Add the gate to `.github/workflows/release.yml`: run `clue validate --forbid-changes` from source, stamped with the tag's version, after the release notes are extracted and before the build
- [x] Rewrite the workflow header comment so it points at the gate instead of asserting the agreement it never checked
- [x] Add `TestSanity_ReleaseRunsTheJudgeStampedAsTheTag` to `cmd/clue/main_test.go`, covering the stamped invocation and its position ahead of the build and the publish step
- [x] Prove the gate empirically: a mismatched stamp fails and names both versions, a matching stamp passes
- [x] Record the release pipeline's gate in `docs/capabilities/CAP-004-ship/design.md`
- [x] Record the workflow-not-corpus-rule choice as a decision-log row
- [x] Add the `[Unreleased]` changelog entry for the fix
- [x] Run repository and corpus verification
