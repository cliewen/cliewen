---
id: CH-196-tasks
type: tasks
status: open
links: [CH-196]
title: CH-196 tasks
---

# CH-196 tasks

- [ ] Probe one file per raw-reading consumer class (agent hubs in a fresh Claude Code session, a skill entry point and reference, a guide page, `CHANGELOG.md` release extraction, a folder README through `clue scaffold`); stop and report if any loses or misreads content
- [ ] Record PDR-065: Markdown carries frontmatter by default, the `type`/`title` document header for non-artifacts, and the pull-request-template exception
- [ ] Add AC-213 (CAP-002): `clue validate` requires a header on every Markdown file under `docs/` and `changes/`, with folder READMEs as `type: index` non-artifacts; positive and negative evidence
- [ ] Add AC-214 (CAP-002): `clue validate` requires a header on Cliewen-delivered files when present, never on the pull-request template or other adopter Markdown; positive and negative evidence
- [ ] Add AC-215 (CAP-001): `clue init` and `clue scaffold` write every Markdown file with frontmatter except the pull-request template, generated skill files included, and index regeneration preserves it; positive and negative evidence
- [ ] Add AC-216 (CAP-001): `clue migrate` adds the header to managed folder READMEs and reports, never rewrites, delivered adopter-owned files lacking one; positive and negative evidence
- [ ] Implement validation, scaffold templates, skill generator output, index regeneration, and the migrate repair and notice (AC-213 to AC-216)
- [ ] Add C-024, this repository's machine-enforced constraint that every tracked Markdown file carries frontmatter, with its named exceptions and the test that holds it
- [ ] Add headers to this repository's Markdown files (root files, guide pages, folder READMEs, skill resources); regenerate skills, indexes, and the evidence manifest
- [ ] Update `/guide` prose with the humanizer skill where it describes artifact frontmatter or the README exemption
