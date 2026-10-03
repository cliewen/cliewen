---
id: CH-196
type: change
status: open
links: [G-022, CAP-002, CAP-001, CAP-005]
title: Markdown files carry frontmatter by default
---

# CH-196 — Markdown files carry frontmatter by default

This change is plan-less: it carries out [G-022](../../docs/goals/G-022-every-markdown-file-carries-frontmatter.md), which the maintainer raised during CH-195 and asked to start once CH-195 merged.

## Why

Today only corpus artifacts and change workspaces carry frontmatter. Folder READMEs, guide pages, root files, the agent hubs, and every file `clue init` writes outside `docs/` have none, so neither an agent nor `clue` can tell from the file itself what it is. The maintainer wants frontmatter to be the default, so that a Markdown file needs a reason to go without it.

## What changes

The maintainer settled the goal's three open questions on 2026-10-03:

- **Fields.** A Markdown file that is not a corpus artifact carries a small document header, `type` and `title`, with no `id`. It never enters the identity ledger and is never an artifact. Corpus artifacts keep their existing field set.
- **Exceptions.** Only the pull-request template, because the forge copies it verbatim into every new pull request body. Everything else gets frontmatter, including `README.md` and `CONTRIBUTING.md`, accepting that GitHub renders the header as a small table at the top of those files.
- **Reach for adopters.** Cliewen's own files and the corpus, not the adopter's other Markdown:
  - `clue validate` requires a header on every Markdown file under `docs/` and `changes/`. A folder README carries `type: index` and stays a non-artifact. It also requires a header on the files Cliewen delivers when they are present: `AGENTS.md`, `CLAUDE.md`, the `.clue/evidence/` guides, `.github/cliewen-wall.md`, and every generated skill file.
  - `clue init` and `clue scaffold` write every Markdown file with frontmatter except the pull-request template, and index regeneration preserves it.
  - `clue migrate` adds the header to managed folder READMEs (generated skills already update as a set). It reports, and never rewrites, adopter-owned delivered files that lack one: the hubs, under [PDR-023](../../docs/decisions/PDR-023-tool-notice-and-hub-instruction.md), the evidence guides, and the wall checklist.

This repository goes further than adopters, as its source-repository convention: every tracked Markdown file carries frontmatter. A new machine-enforced constraint holds it. The named exceptions are the pull-request templates and the test fixtures that model pre-contract adopter repositories, which must stay header-less to keep testing migration.

A new process decision records the default, the two-field header, and the exception rule. New criteria cover validation, scaffolding, generated skills, and migration.

## Scope boundary

- No `id` on non-artifact files. No new ledger kind, no `links`, no `status`.
- The vocabulary of `type` values is open for adopters' own files. Cliewen fixes only the values it writes (for example `index`, `agent-hub`, `skill`, `skill-reference`, `guide`), and only `index` is checked, on corpus folder READMEs.
- Adopter Markdown outside `docs/`, `changes/`, and Cliewen's delivered paths is untouched and unchecked.
- Completed plans and analyses already carry artifact frontmatter and are not touched. `CHANGELOG.md` gains a header at the top of the file; released sections stay verbatim.

## Challenge

**Riskiest assumption:** every tool that reads these files raw tolerates a leading YAML block. The tools are:

- GitHub's renderer, for the repository front page and the community-profile files.
- Claude Code loading `CLAUDE.md` and following its `@AGENTS.md` import, and other agent hosts reading `AGENTS.md`.
- Agent hosts reading generated skill files whose frontmatter they interpret.
- VitePress, where `title` in frontmatter overrides the page title.
- The release workflow, which extracts a `CHANGELOG.md` section verbatim.
- `clue scaffold`, which rewrites index blocks inside the READMEs.

If a host treats the hub's header as instructions, or drops the import, adopters would lose routing.

**Credible alternative:** keep header-less files and record each file's kind in one registry, such as a manifest under `.clue/`. That avoids touching raw-read files, but the file still cannot say what it is, and the registry drifts from the files, which is the problem the goal names.

**Cheapest test:** before the bulk rollout, add the header to one file of each consumer class:

- `CLAUDE.md` and `AGENTS.md`, checked in a fresh Claude Code session that still reports the routing hub.
- A skill entry point and a reference.
- A guide page, checked with `npm run guide:build` and the rendered title.
- `CHANGELOG.md`, run through the release-notes extraction test.
- A folder README, run through `clue scaffold`.

**Stop or revise:** if a raw-reading host loses or misreads routing, or a skill host rejects the extra keys, stop and bring the maintainer the evidence before exempting that class. Do not exempt it silently.

**Meeting every criterion and still failing the person:** every file gets a header whose values are boilerplate, such as `type: doc`, that no reader or tool uses. Then the repository carries noise instead of meaning. Against that, Cliewen writes specific `type` values for what it delivers, `clue` uses `type: index` to tell folder READMEs from artifacts, and every `title` is the file's real title, not a placeholder.
