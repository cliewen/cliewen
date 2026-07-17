# Changelog

All notable, user-visible changes to `clue` and the Cliewen skills. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow semver. Each GitHub release body is this file's matching version section, extracted verbatim by the release workflow — a release with no section here fails.

## [Unreleased]

### Added

- **`clue init` — the whole convention in one command.** Run it in a new or existing repository and it materializes the `docs/` corpus (folder READMEs that explain, in plain language, what each record type is and when a change updates it), an `AGENTS.md` routing hub, the five agent skills (`.agents/skills/` plus a `.claude/skills/` mirror for Claude Code), and a CI workflow template that runs `clue validate`. The scaffolding is embedded in the binary — no network or checkout needed — and `clue validate` is green on the result immediately. `init` is idempotent: re-running regenerates README index blocks from folder contents and touches nothing else, and it never overwrites an existing file (skips are reported), so your own `AGENTS.md` or READMEs survive.
- **Quickstart.** The README now takes a new user from install through `clue init` to their first change loop and green validate on one page — prerequisites stated up front, skills linked at the moment the first loop needs them.

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
