---
id: PDR-065
type: decision
status: verified
links: [G-022, PDR-023, CAP-002, CAP-001, CAP-005]
title: Markdown files carry frontmatter by default, and a non-artifact carries only type and title
author: agent
accepted-by: Flemming N. Larsen (2026-10-08, conversation)
---

# PDR-065 — Markdown files carry frontmatter by default

## Context and problem statement

Corpus artifacts and change workspaces carry frontmatter, and `clue` reads it. Every other Markdown file has none: folder READMEs, guide pages, root files, the agent hubs, and most of what `clue init` writes. Neither an agent nor `clue` can tell from such a file what it is. A folder README under `docs/` was recognised as a non-artifact only because it lacked frontmatter, so it could never gain any ([G-022](../goals/G-022-every-markdown-file-carries-frontmatter.md)).

## Decision outcome

**A Markdown file carries YAML frontmatter unless it has a named reason not to.** The one reason Cliewen names is a file whose consumer copies it verbatim into other content: the pull-request template, which the forge pastes into every new pull request body.

**A file that is not a corpus artifact carries a document header of exactly two required fields:**

```yaml
---
type: index
title: Goals
---
```

- `type` says what the file is.
- `title` is the file's real title.

The header has no `id`, so the file never enters the identity ledger, cannot be linked by ID, and is never an artifact. It has no `status` and no `links`. Corpus artifacts keep their existing field set.

- A folder README under `docs/` carries `type: index`. That value is what makes it a non-artifact, and `clue validate` checks it.
- Cliewen writes fixed values for the other files it delivers: `agent-hub`, `skill`, `skill-reference`, `evidence-guide`, and `checklist`.
- For an adopter's own non-artifact files, the `type` vocabulary is open.

**Reach for adopters is Cliewen's own files and the corpus:**

- `clue validate` requires the header on every Markdown file under `docs/` and `changes/`. It also requires it on Cliewen's delivered files when they are present: `AGENTS.md`, `CLAUDE.md`, the `.clue/evidence/` guides, `.github/cliewen-wall.md`, and every generated skill file.
- `clue init` and `clue scaffold` write every Markdown file with the header, except the pull-request template.
- `clue migrate` adds the header to managed folder READMEs. It reports, and never rewrites, adopter-owned delivered files that lack one ([PDR-023](PDR-023-tool-notice-and-hub-instruction.md)).
- An adopter's other Markdown is neither changed nor checked.

## Consequences

- GitHub renders a header as a small table at the top of a Markdown file, including on a repository's front page. The maintainer accepted that cost for `README.md` and the other root files rather than exempt them.
- Raw-reading consumers tolerate the header. On 2026-10-03, probes showed the following:
  - Claude Code follows `@AGENTS.md` and strips the header before the model sees it.
  - Codex reads the header raw and still follows the hub.
  - Skill hosts load skill files that carry the extra keys.
  - VitePress keeps the page title.
  - Release-notes extraction from `CHANGELOG.md` is unchanged.
- This repository's own convention is stricter: every tracked Markdown file carries frontmatter. That convention is a source-repository rule, not an adopter one ([C-024](../constraints/C-024-every-markdown-file-carries-frontmatter.md)).

## Rejected: full artifact identity on every file

Giving each file an `id`, `status`, and `links` would put hundreds of identities in the ledger for files that nothing links by ID, and would make every README an artifact with a lifecycle it does not have.

## Rejected: a registry of file kinds instead of headers

A manifest under `.clue/` naming each file's kind would leave raw-read files untouched. But the file would still not say what it is, and the registry would drift from the files it describes.

## Rejected: exempting the root community files

Exempting `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, and `CHANGELOG.md` would keep the forge's front page unchanged. But it would make the most-read files the exception to a rule meant to have few.
