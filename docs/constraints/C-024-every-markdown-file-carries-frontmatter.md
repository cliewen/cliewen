---
id: C-024
type: constraint
status: active
links: [PDR-065, G-022]
title: Every Markdown file this repository tracks carries frontmatter
source: PDR-065 (source-repository convention)
enforcement: machine
---

# C-024 — Every Markdown file this repository tracks carries frontmatter

Every tracked Markdown file opens with frontmatter that carries a non-empty `type` and `title`. A corpus artifact's existing fields satisfy this, and any other file carries the two-field document header [PDR-065](../decisions/PDR-065-markdown-carries-frontmatter-by-default.md) defines. This is stricter than what Cliewen asks of adopters, whose own Markdown outside `docs/`, `changes/`, and Cliewen's delivered files is never checked.

The named exceptions:

- The pull-request templates (`.github/pull_request_template.md` and its scaffold source), because the forge pastes them verbatim into every new pull request body.
- The fixtures under `internal/migrate/testdata/pre-contract/`, which model an adopter repository from before document headers existed and must stay header-less to test the migration that adds them.

**Checked by:** `TestSanity_EveryTrackedMarkdownFileCarriesFrontmatter` in `internal/corpus/documents_test.go`, over `git ls-files '*.md'`, with the same exception list.
