---
id: CH-008-tasks
type: tasks
status: open
links: [CH-008]
title: Task breakdown for CH-008
---

# Tasks — CH-008

- [ ] `CHANGELOG.md`: Keep a Changelog header, `[Unreleased]` with this change's own entry, back-filled `[0.1.0]` section
- [ ] `release.yml`: extract the tag's section verbatim → `body_path`; fail the release when the section is missing/empty; drop `generate_release_notes`
  - [ ] `TestSanity_ReleaseNotesComeFromChangelog` guards the workflow shape (extraction present, auto-generation gone)
- [ ] AGENTS.md: binding rule — release notes are user-facing and live in `CHANGELOG.md`
- [ ] `clue-delta` digest step: record user-visible impact in `[Unreleased]`; `clue-verify`: matching checkbox
- [ ] ADR-012 (release notes from CHANGELOG.md; rejected: auto-generation, per-release files, tag messages)
- [ ] README: one-line pointer to `CHANGELOG.md`
- [ ] Digest: decisions index (ADR-012), CAP-004 design refresh, delete `changes/CH-008-changelog/`
- [ ] clue-verify checklist, then PR
- [ ] After merge: re-sync the v0.1.0 release body to the `[0.1.0]` section verbatim
