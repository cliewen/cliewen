---
id: G-022
type: goal
status: proposed
links: [VIS-001]
title: Every Markdown file carries YAML frontmatter, in this repository and for adopters
---

# G-022 — Every Markdown file carries YAML frontmatter

**Who wants it:** the maintainer, for this repository and for every adopter (2026-10-03), raised while CH-195 was in progress.

**Why:** today only corpus artifacts under `docs/` and change workspaces carry frontmatter. Everything else, such as guide pages, root files like `CONTRIBUTING.md`, folder READMEs, and `.github` templates, has none, so neither an agent nor `clue` can tell from the file itself what it is, what state it is in, or what it links to. The maintainer wants frontmatter to be the default for any Markdown file, so a file has to have a reason to go without it rather than the other way round.

**Open before a plan commits to it:**

- Which files are in scope, and which are deliberate exceptions: `CHANGELOG.md` (published verbatim as release bodies), `AGENTS.md`/`CLAUDE.md` (read raw by agent hosts), generated skills (which already carry their own frontmatter), `.github` templates (rendered by the forge), and guide pages (where VitePress already reads frontmatter for its own keys).
- What fields a non-corpus file needs. The corpus `id`/`type`/`status`/`links`/`title` set may be too much for a file that is not an artifact, and an identity would put it in the ledger.
- How existing adopters get there: a `clue migrate` notice or repair, and whether `clue validate` enforces it or only reports it.

**Success looks like:**

- A Markdown file without frontmatter is either a named exception or a finding.
- `clue init` materializes every Markdown file with frontmatter, and adopters have a supported migration path.
- No file that a tool reads raw breaks because frontmatter was added.
