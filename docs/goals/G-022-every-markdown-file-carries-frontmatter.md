---
id: G-022
type: goal
status: accepted
links: [VIS-001]
title: Every Markdown file carries YAML frontmatter, in this repository and for adopters
---

# G-022 — Every Markdown file carries YAML frontmatter

**Who wants it:** the maintainer, for this repository and for every adopter (2026-10-03), raised while CH-195 was in progress.

**Why:** a Markdown file without frontmatter cannot say what it is, so neither an agent nor `clue` can tell its kind from the file itself. The maintainer wants frontmatter to be the default for any Markdown file, so that a file needs a reason to go without it.

**Settled by [PDR-065](../decisions/PDR-065-markdown-carries-frontmatter-by-default.md):** a non-artifact file carries only `type` and `title`; the pull-request template is the one exception; for adopters the rule covers the corpus and Cliewen's delivered files, while this repository applies it to every tracked Markdown file ([C-024](../constraints/C-024-every-markdown-file-carries-frontmatter.md)).

**Success looks like:**

- A Markdown file without frontmatter is either a named exception or a finding.
- `clue init` materializes every Markdown file with frontmatter, and adopters have a supported migration path.
- No file that a tool reads raw breaks because frontmatter was added.
