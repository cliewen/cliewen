# Changelog

All notable, user-visible changes to `clue` and the Cliewen skills. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow semver. Each GitHub release body is this file's matching version section, extracted verbatim by the release workflow — a release with no section here fails.

## [0.4.0] - 2026-07-19

### Changed

- **The public guide now gives newcomers a practical starting point.** It separates required and optional tooling, explains how agents maintain the corpus alongside code, routes cheap decisions to the log and expensive decisions to ADRs or PDRs, and provides prompts for greenfield work, routine changes, and brownfield adoption without claiming live synchronization or cross-repository validation.
- **`clue validate` now catches hidden and duplicated frontmatter.** A UTF-8 byte-order mark anywhere in a corpus markdown file fails validation with guidance to strip it — a BOM can hide a frontmatter fence from the parser. A complete second frontmatter block opening an artifact body fails as a leftover from a brownfield conversion that prepended instead of replaced; the `clue-extract` contract now spells out that a converted file carries exactly one frontmatter block.
- **Public-readiness documentation now states the current access boundary.** The root README explains that the repository remains private during the readiness campaign and distinguishes collaborator installation from the final anonymous-install milestone. Historical evidence from a private adopter repository remains preserved, while the Cliewen corpus now carries every rule derived from it without relying on inaccessible target files.
- **Cliewen skills now declare their ownership without claiming neighboring skills.** Generated skills carry `cliewen-skill: true`; `clue validate` applies version agreement and binary-drift checks only to marked skills, so unrelated skills can coexist under `.agents/skills/`. The five pre-marker Cliewen skill names fail with reinstall guidance instead of being silently ignored, while `clue init` remains non-destructive and never overwrites or deletes an existing skill.
- **Analysis now preserves its evidence boundary.** `clue-analysis` asks agents to pin source revisions when possible, record conditions relevant to reproduced results, distinguish observed facts from inferences and unverified intent, and avoid treating repository activity as maintainer intent without explicit evidence.
- **The five standalone skills are now generated from shared canonical instructions.** Their public names and independent installation stay unchanged, while repeated rules such as decision routing, change tiers, repository-local conventions, and the review boundary are authored once and composed into each complete skill. Repository tests reject changed, missing, or unexpected generated files in either the agent or `clue init` template tree.
- **Pull requests open ready for review, never as drafts.** The change loop now keeps unfinished work on its branch, runs the verification checklist, and only then opens the PR as the human review gate; `clue-delta`, `clue-extract`, and `clue-verify` carry the rule.
- **Decision records stay focused and timeless in every workflow.** The analysis, planning, change, extraction, and verification skills now keep triggering incidents, chronology, conversations, implementation details, and review history out of ADRs/PDRs; a decision records its outcome and only the enduring context and rationale needed to understand it, while findings, change workspaces, PRs, and Git retain the history.
- **A declared plan revision may ride with its implementing pull request** (`clue-plan`). The default stays a dedicated plan PR, but a semantic plan revision that surfaced during implementation may travel with the change that implements it when four conditions hold: the PR explicitly declares the revision, a correctly typed decision record backs it, the PR calls it out for deliberate human approval, and an explicit objection reverts the revision — the milestone stays open — without blocking the rest of the change.

### Added

- **The public Cliewen guide is deployment-ready.** A handwritten site now teaches newcomers the methodology, corpus taxonomy, change loop, and lifecycle skills without requiring them to read the internal corpus first. Its locked VitePress build fails on broken internal links in CI, renders inline Mermaid diagrams, and carries a GitHub Pages workflow that stays safely skipped while the repository is private and becomes eligible when visibility flips.
- **Cliewen now has a clear community front door.** Contributors can use structured forms for reproducible bugs and proposed outcomes, follow one contribution guide from first report through human-reviewed merge, and rely on private, published routes for security vulnerabilities and conduct concerns. A proposed-goal issue records demand without silently changing the accepted plan, while blank public issues are disabled to keep intake actionable.
- **`clue scaffold` — regenerate the README index blocks, nothing else.** The index-regeneration engine `clue init` runs is now a standalone command: run it in any Cliewen repo and the taxonomy README index blocks are rebuilt from folder contents — your hand-written entry lines survive as long as their targets do, missing entries are appended, and prose outside the markers is never touched. It materializes nothing: missing folder READMEs are reported rather than invented, and a path without a `docs/` tree is an error.
- **`clue init` — the whole convention in one command.** Run it in a new or existing repository and it materializes the `docs/` corpus (folder READMEs that explain, in plain language, what each record type is and when a change updates it), an `AGENTS.md` routing hub, the five agent skills (`.agents/skills/` plus a `.claude/skills/` mirror for Claude Code), and a CI workflow that runs `clue validate` — passing with a visible warning until you vendor the pinned binary it expects (the arming commands are in its comments), so a fresh repo is never red before its first change. The scaffolding is embedded in the binary — no network or checkout needed — and `clue validate` is green on the result immediately. `init` is idempotent: re-running regenerates README index blocks from folder contents and touches nothing else, and it never replaces an existing file (skips are reported) — your own `AGENTS.md` survives, and a taxonomy README of your own just gains an index block.
- **Quickstart.** The README now takes a new user from install through `clue init` to their first change loop and green validate on one page — prerequisites stated up front, skills linked at the moment the first loop needs them.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.4.0`, or download a prebuilt binary for your platform from the release assets and verify it against `SHA256SUMS`. Update vendored Cliewen skills from this release's `.agents/skills/`; a 0.4.0 binary rejects older Cliewen skill versions as drift.

## [0.3.0] - 2026-07-17

### Changed

- **The review boundary is now explicit in the change loop.** The skills previously described the PR and merge without defending them, and an autonomous agent could satisfy the letter of the loop while bypassing review entirely — stacking changes on unmerged work, fabricating local merge commits, and pushing straight to `main`. The `clue-delta` skill now states the rules as prohibitions: every change branches from the current tip of `main` (never from unaccepted work), each author takes one change to its PR before starting the next, review fixes stay on the reviewed branch, and the agent never merges its own PR, never creates a merge commit into `main`, and never pushes to `main` — after opening the PR it stops and waits for the human. `clue-verify` gains matching checklist items, including a rebase-and-recheck step when a parallel change merges first. Team parallelism is untouched: any number of changes may be in flight, each rooted at `main` with its own PR.
- **A merged PR is acceptance, not a go signal.** When the human reports that a PR was merged, the agent no longer silently continues with the next task. It first says where the plan stands: the next step described in plain language — what it is about, not just document IDs — followed by the question whether to start it. When the plan has nothing left, the agent says so and asks the human what to do next. The `clue-delta` skill carries the rule.
- **The digest is never a task in `tasks.md`.** The digest precondition (every task `[x]` or `[-]`-with-reason before the workspace is deleted) applies to the work; a self-referential "digest the change" task could only be ticked falsely or left violating the precondition, so `clue-delta` now forbids it.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.3.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`. Update the vendored skills to this release's `.agents/skills/` — the binary fails validation against 0.2.0 skills (drift check). The `clue` binary itself is unchanged from 0.2.0; this release ships the skill changes, and the version bump keeps the binary+skills pair matched.

## [0.2.0] - 2026-07-15

### Added

- **Light change tier**: a change that decides nothing and changes no meaning (typos, doc clarity, dependency bumps, pure refactors, CI plumbing) no longer needs a `/changes/` workspace — branch, commit, and open a PR whose description is the proposal. The moment a decision, open question, or AC change appears, escalate to the full loop. The `clue-delta` skill carries the qualification test; `clue-verify` starts by checking the tier is right.
- **Decision log**: full decision records are now reserved for decisions that are expensive to reverse; everything else is a dated row in `docs/decisions/log.md` (litmus test: cheap and local to reverse → log row). `clue validate` lints the new `log` artifact type.
- **Typed decision records**: ADR keeps its industry meaning — Architectural Decision Records cover the structure of the software and the corpus format, nothing else. Expensive-to-reverse decisions about how the project works (change tiers, decision acceptance, validation strategy) get their own series: **PDR**, Project/Process Decision Records — same MADR template, same `inferred`/`verified` provenance. Decisions that adopt a well-established practice cite it by name instead of re-deriving it. Existing records filed under the wrong type were renamed into the right series.
- **Convention register**: methodology rules that lived only in prose are now constraint artifacts in `docs/constraints/`, each naming its `source` and an `enforcement` class (`machine|agent|human`). `clue validate` requires both fields, checks the vocabulary, and reports the count of agent-enforced constraints on its OK line — the visible backlog of rules awaiting machine checks, each carrying its promotion trigger.

### Changed

- **Merge binds, approval signs**: merging a PR makes the decision records (ADRs and PDRs) it introduces binding — in force immediately, no approval ritual blocks shipping. A decision is marked `verified` only when a human explicitly approves it; each approver signs `accepted-by`, approvals accumulate, and the acceptance date is the first approval. The `inferred` count now honestly means "in force but unapproved". The `clue-delta` and `clue-verify` skills carry the rule.
- Release pages are now published from this changelog: each GitHub release body is the matching version section of `CHANGELOG.md`, written for users — no more auto-generated PR lists.
- The agent skills are repo-agnostic: repo-local conventions live in your AGENTS.md, which extends the methodology but never overrides it — a conflict between the two is surfaced as an open question. The OpenSpec extraction mapping moved from `clue-extract`'s skill text to `mappings/openspec.md` under the same skill, and the skills no longer cite cliewen's internal document IDs — every rule is stated in full where you read it.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.2.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`. Update the vendored skills to this release's `.agents/skills/` — the binary fails validation against 0.1.0 skills (drift check).

## [0.1.0] - 2026-07-13

First release of `clue`, the command-line validator for the Cliewen methodology. It keeps a repository's documentation corpus — goals, plans, capabilities, decisions — and its agent skills consistent, traceable, and versioned.

### Added

- **`clue validate`** lints the corpus: frontmatter core fields, ID uniqueness, status vocabulary, cross-links, folder READMEs and index blocks, provenance, and acceptance-criteria-to-test traceability.
- **`clue version`** reports the release the binary was built from; untagged source builds report `dev`.
- **Versioned skills**: every agent skill declares a version in its frontmatter, and `validate` fails when skills disagree with each other or drift from the binary's release.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.1.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`.
