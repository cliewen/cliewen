# Changelog

All notable, user-visible changes to `clue` and the Cliewen skills. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow semver. Each GitHub release body is this file's matching version section, extracted verbatim by the release workflow — a release with no section here fails.

## [Unreleased]

### Added

- **Cliewen installs from inside Claude Code.** `/plugin marketplace add cliewen/cliewen` followed by `/plugin install cliewen@cliewen` adds one skill, `/cliewen:setup`, which works out your platform, installs `clue` with the same checksum-verifying script the guide documents, checks that the binary reports a release version, and then asks before running `clue init` — nothing is written into your repository until you agree. It installs the binary and nothing else: the five skills that run the Cliewen loop are not in the plugin, and will not be. Those are committed files that `clue init` writes into your repository, each stamped with the version of the binary that wrote it, and `clue validate` fails if the two ever disagree. A plugin's files live in a per-user cache instead, outside where that check can see them and shared across every repository you open, so a bundled copy would quietly contradict any repository pinned to a different release. The new *Install from Claude Code* page explains the boundary and says which route to use where; CI should keep using the pinned, checksum-verified binary, never a plugin.

## [0.8.0] - 2026-07-27

The release where installing `clue` stops being a six-step errand, and where the guide finally argues its own case.

### Added

- **`clue` installs with one command.** `curl -fsSL https://cliewen.dev/install.sh | sh` on macOS and Linux, or `irm https://cliewen.dev/install.ps1 | iex` on Windows, replaces reading a platform table, downloading a binary, checking its SHA-256 by hand, renaming it, making it executable and editing your `PATH`. The script works out your platform, fetches the matching binary, and checks it against the release's `SHA256SUMS` before it writes anything — if that check fails, nothing is installed. It needs no administrator rights, and `CLUE_INSTALL` and `CLUE_VERSION` let you choose the directory and the release. Upgrading is the same command again. The manual download is still documented, unchanged, for anyone who prefers it or whose machine cannot run the script. If you already run Cliewen in a repository, note that upgrading the binary does not touch the skills committed there, so `clue validate` will correctly report drift until you update those too; the guide's *Operate safely* page explains how.
- **The guide now explains why Cliewen works the way it does.** If you are deciding whether to adopt it, the other pages tell you the rules; the new page, *The design of Cliewen*, gives you the argument. It says what problem the method is answering — that accepting a change, not writing it, is what agents made scarce — which failures of the earlier approach each rule repairs, which three parts cannot be changed without an explicit decision, what each principle costs you in tokens and time, how you extend the method with your own rules, and what Cliewen does not solve, including the gaps still open today.

### Fixed

- **A release can no longer ship a `clue` and skills that disagree about their version.** The version came from the git tag, while the skills carried a version written by hand, and nothing compared the two — so a forgotten bump published a binary and a skill set that contradicted each other. You would have met that as a drift error in your own repository, after installing both halves, with nothing you had done to explain it. Cutting a release now runs the same drift check you run, stamped as the tag, before anything is built: a mismatch stops the release and names both versions, so the pair you download always agrees with itself.
- **A guide link pointed at a finished plan.** The corpus page invited you to browse Cliewen's "active public campaign" and linked one that finished several campaigns ago. It now links the plan actually being worked on.

### Install

`curl -fsSL https://cliewen.dev/install.sh | sh` on macOS and Linux, `irm https://cliewen.dev/install.ps1 | iex` on Windows, or `go install github.com/cliewen/cliewen/cmd/clue@v0.8.0`. You can still download a prebuilt binary from the release assets and verify it against `SHA256SUMS` by hand; those asset names are unchanged, so a vendored CI wall pinned to an earlier release keeps working exactly as before. Update vendored Cliewen skills from this release's `.agents/skills/`; a 0.8.0 binary rejects older Cliewen skill versions as drift.

## [0.7.0] - 2026-07-25

Mostly about moving a repository into Cliewen when it already keeps decision records of its own, plus two fixes for `clue` judging the same files differently depending on where it ran.

### Added

- **Decision records you already have convert automatically.** If your decisions live in `docs/decisions/` or `doc/adr/` as MADR records (3.x or 4.x) or in the older Nygard style, the `clue-extract` skill now converts them by stated rules instead of improvising. Your numbering survives as each record's ID; a record with no number, or one whose number is already taken, gets one assigned the same way on every run, so re-running the conversion reproduces the same result. A record you rejected still reads as rejected rather than as something you adopted, and a proposal you never acted on becomes an open question instead of a decision that merging would make binding. The original status is kept in the record's text, and a superseded record links to the one that replaced it. Titles, slugs and dates all have a stated destination.

### Changed

- **`accepted-by:` now means one thing: a human approved this record here.** It records approval given at Cliewen's own merge gate — a pull request approval, a review comment, a written "approved" — and nothing else. If an imported record already named who decided, who was consulted and who was informed, those names, roles and dates are kept in the record's text, but the field itself starts empty, exactly as it does on any record nobody has signed yet. The five shipped skills and both copies of the decisions README say this the same way now.
- **Converting an existing repository copes with what real repositories look like.** Agents are pointed at the converted documentation from every instruction file your repository actually has — `AGENTS.md` plus files like `CLAUDE.md` or `.cursor/rules` — instead of assuming `AGENTS.md` is the only one. A repository too large to give every acceptance criterion its tests in one go may now leave those capabilities marked `draft` on purpose, said plainly in the conversion report, rather than the whole conversion stalling until coverage is complete. And where a source requirement has no stable ID of its own, one is minted the same way each run, so re-running does not renumber everything.

### Fixed

- **`clue init` no longer writes into a skills folder you share between checkouts.** If `.claude/skills` — or `.claude`, or `.agents/skills` — is a symlink to a skills tree kept outside the repository, `init` used to follow it and install Cliewen's skills into that shared tree, quietly changing a directory you were not initializing. It now writes nothing through a symlinked directory, leaves it untouched, and tells you so on its own line (`linked .claude/skills (symlink — mirror skipped, nothing written through it)`) with a count in the summary, so an empty mirror is explained rather than mysterious. Everything the link does not block is still created, and the promise that no existing file is ever overwritten is unchanged. If you want the skills in the shared tree, install them there yourself. Reading a skill through a symlink works as before.
- **`clue validate` now reaches the same verdict on every operating system.** A skill whose file is named `SKILL.md` — the spelling Claude Code and Anthropic skills use — was found on Windows and macOS but silently ignored on a Linux CI runner, so the same commit could pass on your laptop and fail in CI, or the reverse. The validator now finds the file whatever the capitalisation, and if a directory somehow contains two spellings at once it says so instead of quietly picking one. Only a real file counts, so an unusual file with that name cannot stall the check. You do not need to rename anything.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.7.0`, or download a prebuilt binary for your platform from the release assets and verify it against `SHA256SUMS`. Update vendored Cliewen skills from this release's `.agents/skills/`; a 0.7.0 binary rejects older Cliewen skill versions as drift.

## [0.6.0] - 2026-07-24

The release where Cliewen says which of its parts are non-negotiable, and where the guide became something you can adopt from without reading the project's own documentation first.

### Changed

- **Deciding how much process a change needs is now three rules you can remember.** Nothing about meaning changes: an ordinary branch and pull request, no paperwork. Meaning is touched but not changed: still no workspace, and the pull request description is the proposal. Everything else: the full loop. Two guards sit above them — if you cannot tell which applies, take the heavier one, and if a decision or a change of meaning turns up while you work, move up immediately. The skills, the `AGENTS.md` template and the guide all state it this way now. The boundaries themselves did not move, so work that was light or full before is light or full still.
- **Quality bars are now recorded as constraints.** The separate `quality` record type and its `docs/quality/` folder are gone. A measurable bar — a coverage floor, a maximum onboarding time — is now a constraint like every other cross-cutting rule, naming where it comes from and whether a machine, an agent or a human enforces it. `clue init` no longer creates a `quality/` folder, and conversion maps coverage and quality gates to constraints. If you already have a `docs/quality/` folder it keeps validating: the type was dropped from Cliewen's own vocabulary, not from what your documentation may contain.
- **`clue validate` now accepts record types you invented.** A type the validator does not recognise is treated as yours and checked against the default lifecycle (`draft` → `active` → `retired`) instead of being rejected, so you can keep your own kinds of document under `docs/` without waiting for the tool to learn about them. Everything else still applies: required fields, unique IDs, links that resolve, a valid status. Statuses are now one default lifecycle plus a short list of types that differ for a stated reason, rather than a separate vocabulary per type.
- **Cliewen now states its core, and protects it.** Three parts carry the method: the traceable line from goal to test, the rule that a human performs the merge, and `clue validate` as a judge that gives the same verdict everywhere. Changing what any of them means always takes an explicit decision record and human acceptance — it can never ride along inside another change. Everything else is yours to extend.
- **Running Cliewen over time is documented.** The guide now says exactly what the released tool supports, how to upgrade the binary, the skills and CI together so they stay in step, and how to recover from version drift, files that initialization skipped, a failing validation, a rollback or removing Cliewen altogether — none of it by going around review. Every guide page ends with one clear next step.
- **The guide has a visual identity.** The thread-and-check logo appears in the header and on the homepage, and browsers use it as the site icon.
- **You can try Cliewen without risking a repository you care about.** The front page says who the method is for and separates the three things people conflate — Cliewen the method, `clue` the tool, and the documentation itself — before any taxonomy appears. Prebuilt binaries are now the main way to install, with a download table and short checksum and `PATH` steps. The throwaway demo now activates a criterion so you see the real "no test for this" diagnostic before you delete it again. The guide also states the boundary honestly: `clue` detects whether supported test evidence exists; whether the tests are the right ones is the review loop's job and yours. The site now loads correctly from `https://cliewen.dev/`.
- **Turning CI into a wall now has complete instructions and a proof.** The guide shows how to vendor and verify a pinned `clue` binary, require the check, protect your default branch without a bypass, and then prove to yourself that an unfinished change genuinely cannot merge. Adoption advice now starts from the smallest useful thread — one goal, one capability, one criterion, a positive and a negative test — and says when the fuller structure earns its keep and when Cliewen is the wrong tool for you.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.6.0`, or download a prebuilt binary for your platform from the release assets and verify it against `SHA256SUMS`. Update vendored Cliewen skills from this release's `.agents/skills/`; a 0.6.0 binary rejects older Cliewen skill versions as drift.

## [0.5.1] - 2026-07-22

A fix for review agents reporting problems that were not problems.

### Changed

- **A review finding now has to point at a rule it actually breaks.** The agent that reviews a change before you see it was raising findings from its own preferences. It now checks the project's recorded decisions and lifecycle rules first, which stops two recurring false alarms: treating optional human code review as if it were missing, when the required gate is the human merge, and reporting a released version's changelog section as a missing `[Unreleased]` section.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.5.1`, or download a prebuilt binary for your platform from the release assets and verify it against `SHA256SUMS`. Update vendored Cliewen skills from this release's `.agents/skills/`; a 0.5.1 binary rejects older Cliewen skill versions as drift.

## [0.5.0] - 2026-07-21

Agents now review their own work before asking you to, and Cliewen became installable without credentials.

### Changed

- **Your agent's change is challenged before it reaches you.** Once the work is committed, a second agent reviews it — on hosts that support it, in a fresh context with no memory of writing the code and nothing to defend; on other hosts, in the same context, disclosed to you as the weaker fallback. Findings go back to the implementing agent, and any real fix invalidates the pass and starts another. Only a clean review on the current commit may open the pull request. This is not a substitute for your judgement at the merge, and it is not a demand that you review the code twice: CI plus branch protection is what stops an agent skipping the gate quietly.
- **Cliewen is publicly reachable.** The repository, the v0.4.0 release assets and the guide no longer need credentials. Installing from source and downloading release binaries both work anonymously.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.5.0`, or download a prebuilt binary for your platform from the release assets and verify it against `SHA256SUMS`. Update vendored Cliewen skills from this release's `.agents/skills/`; a 0.5.0 binary rejects older Cliewen skill versions as drift.

## [0.4.0] - 2026-07-20

The big one: `clue init` puts the whole convention in place with one command, there is a public guide, and the rules agents follow got a great deal stricter about not cutting corners.

### Added

- **`clue init` — the whole convention in one command.** Run it in a new or existing repository and it creates the `docs/` structure (with folder READMEs that explain in plain language what each kind of record is and when a change updates it), an `AGENTS.md` that routes agents to it, the five agent skills (in `.agents/skills/`, mirrored to `.claude/skills/` for Claude Code), and a CI workflow that runs `clue validate` — passing with a visible warning until you vendor the pinned binary it expects, so a fresh repository is never red before its first change. Everything is embedded in the binary, so no network or checkout is needed, and `clue validate` passes on the result straight away. Running it again is safe: it regenerates the index blocks in folder READMEs and touches nothing else. It never replaces a file you already have, and it tells you what it skipped — your own `AGENTS.md` survives, and your own folder README simply gains an index.
- **Quickstart.** The README takes you from install through `clue init` to your first change and a passing validation on one page, with prerequisites stated up front and the skills linked at the point you first need them.
- **`clue scaffold` — regenerate the README index blocks and nothing else.** The index generator that `init` runs is now its own command. Run it in any Cliewen repository and the index blocks are rebuilt from what is actually in each folder: lines you wrote by hand survive as long as the files they point to do, missing entries are added, and prose outside the markers is never touched. It invents nothing — a missing folder README is reported rather than created, and a path with no `docs/` tree is an error.
- **A public guide.** A handwritten site teaches the method, the kinds of record, the change loop and the skills without sending you into the project's own documentation first. Its build fails on broken internal links, it renders diagrams inline, and it ships with a Pages workflow.
- **A front door for contributors.** Structured forms for reproducible bugs and for proposing outcomes, one contribution guide that runs from first report to human-reviewed merge, and private published routes for security vulnerabilities and conduct concerns. Proposing a goal records that someone wants it without silently changing the accepted plan; blank issues are off, so every report arrives actionable.

### Changed

- **What you review is exactly what the agent verified.** Agents now commit every edit, confirm the working tree is clean, push, and check that the pull request's head really is the commit they verified before handing it to you. Fixes after review repeat that check on the same pull request. If you asked the agent to stop locally instead, it says clearly that the work is incomplete and not mergeable, rather than implying otherwise.
- **Editorial changes stay out of Cliewen's way.** Agents now classify work before loading anything: a change that affects no behaviour, intent, evidence, decision, plan, policy or method wording gets an ordinary branch, the checks that are relevant, a pull request and your merge — with no change identity, no proposal, no full verification, no plan bookkeeping and no release note.
- **`clue validate` catches two ways a file's header goes wrong.** An invisible byte-order mark anywhere in a documentation file now fails with instructions to strip it, because it can hide the header from the parser entirely. A second complete header block at the top of a file also fails: it is the signature of a conversion that prepended a header instead of replacing one. The conversion contract now says outright that a converted file carries exactly one.
- **Cliewen's skills say which ones are Cliewen's.** Generated skills are marked, and `clue validate` applies its version and drift checks only to marked ones, so your unrelated skills can live in the same folder undisturbed. The five older Cliewen skill names that predate the marker fail with instructions to reinstall rather than being silently ignored, and `clue init` still never overwrites or deletes a skill you have.
- **Analysis keeps its evidence honest.** The analysis skill now asks agents to pin the exact revision they looked at where possible, record the conditions behind any result they reproduced, keep observed facts apart from inference and from unverified intent, and stop treating repository activity as proof of what a maintainer meant.
- **The five skills are generated from one shared source.** Their names and independent installation are unchanged, but rules that appear in several of them — decision routing, how much process a change needs, repository-local conventions, the review boundary — are now written once and composed in. Tests reject a generated file that was edited, is missing, or should not exist.
- **Pull requests open ready for review, never as drafts.** Unfinished work stays on the branch; the agent runs its checklist and only then opens the pull request, because the pull request is the review gate rather than a place to think out loud.
- **Decision records stay focused and readable years later.** The skills now keep the triggering incident, the chronology, the conversation, implementation detail and review history out of decision records. A record states what was decided and the lasting reasons for it; the story stays where stories belong — in findings, in the pull request and in Git.
- **A plan revision may travel with the change that implements it.** The default is still a separate pull request for the plan, but when implementation reveals that the plan itself was wrong, the correction may ride along if four things hold: the pull request declares the revision, a properly typed decision record backs it, the pull request asks you to approve it deliberately, and your objection reverts just the revision — leaving the milestone open — without blocking the rest of the change.
- **The guide gives newcomers a practical starting point.** It separates required from optional tooling, explains how agents maintain the documentation alongside the code, shows which decisions go in the one-line log and which earn a full record, and offers prompts for a new project, a routine change and an existing codebase — without claiming live synchronisation or cross-repository checking that does not exist.
- **The README states the current access boundary.** It explains that the repository is still private during this campaign and distinguishes installing as a collaborator from the anonymous install that comes later. Evidence gathered from a private repository stays preserved, and every rule derived from it is now written down here rather than depending on files nobody else can open.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.4.0`, or download a prebuilt binary for your platform from the release assets and verify it against `SHA256SUMS`. Update vendored Cliewen skills from this release's `.agents/skills/`; a 0.4.0 binary rejects older Cliewen skill versions as drift.

## [0.3.0] - 2026-07-17

Closes the ways an autonomous agent could follow the change loop to the letter and still avoid being reviewed. The `clue` binary is unchanged from 0.2.0; this release is the skills.

### Changed

- **An agent can no longer route around review.** The skills described the pull request and the merge without defending them, which left an agent free to stack work on unmerged changes, create merge commits locally, or push straight to `main` while appearing to comply. The rules are now stated as prohibitions: branch from the current tip of `main` and never from work nobody has accepted; take one change to its pull request before starting the next; keep review fixes on the branch under review; never merge your own pull request, never make a merge commit into `main`, never push to `main`. After opening the pull request the agent stops and waits for you. The verification checklist gained matching items, including rebasing and re-checking when someone else's change merges first. Team parallelism is untouched: any number of changes can be in flight, each branched from `main` with its own pull request.
- **A merged pull request is acceptance, not permission to continue.** When you say a pull request merged, the agent no longer starts the next task on its own. It tells you where the plan stands and what the next step is about — in plain language, not just record IDs — and asks whether to start it. If nothing is left in the plan, it says so and asks what you want next.
- **"Finish the change" can never be one of the change's own tasks.** Finishing requires every task to be ticked or explicitly dropped with a reason, so a task that stands for finishing could only be ticked untruthfully or left blocking itself. Agents are now forbidden to write one.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.3.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`. Update the vendored skills to this release's `.agents/skills/` — the binary fails validation against 0.2.0 skills (drift check).

## [0.2.0] - 2026-07-15

Ceremony becomes proportionate: small changes stop paying for machinery they do not need, and cheap decisions stop being written as if they were expensive.

### Added

- **A lighter tier for changes that decide nothing.** A typo fix, a clearer paragraph, a dependency bump, a pure refactor or CI plumbing no longer needs a workspace of its own: branch, commit, and open a pull request whose description is the proposal. The moment a decision, an open question or a change to an acceptance criterion appears, it moves up to the full loop. The `clue-delta` skill carries the test for which applies; `clue-verify` starts by checking the answer was right.
- **A decision log for cheap decisions.** Full decision records are now reserved for decisions that would be expensive to undo. Everything else is a dated row in `docs/decisions/log.md`. The test is simply whether reversing it later would be cheap and local. `clue validate` checks the new log.
- **Two kinds of decision record, named for what they are about.** ADR keeps its industry meaning — the structure of the software and the format of the documentation, nothing else. Expensive decisions about how the project itself works get their own series, PDR, using the same template and the same approval rules. A decision that simply adopts a well-known practice cites it instead of re-arguing it. Records that were filed under the wrong kind were renamed.
- **Rules that lived only in prose are now records.** Each one names where it comes from and who enforces it — a machine, an agent or a human. `clue validate` requires both, checks the vocabulary, and reports how many rules are still waiting on an agent rather than a machine check, so the backlog of not-yet-automated rules is visible instead of implied. Each carries the trigger that would promote it.

### Changed

- **Merging makes a decision binding; approval signs it.** Merging a pull request puts the decision records it introduces into force immediately — no approval ritual blocks shipping. A record is marked approved only when a human explicitly approves it; each approver is recorded, approvals accumulate, and the acceptance date is the first one. The count of unapproved records now honestly means "in force, but nobody has endorsed the reasoning".
- **Release pages come from this file.** Each GitHub release body is the matching section of `CHANGELOG.md`, written for readers rather than generated from a list of pull requests.
- **The skills no longer assume they are in this repository.** Your own conventions live in your `AGENTS.md`, which may add to the method but never override it; a conflict between the two is raised as an open question instead of being resolved silently. The conversion mapping for OpenSpec moved into its own file under the same skill, and the skills no longer cite this project's internal document IDs — every rule is stated in full where you read it.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.2.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`. Update the vendored skills to this release's `.agents/skills/` — the binary fails validation against 0.1.0 skills (drift check).

## [0.1.0] - 2026-07-13

The first release of `clue`, a command-line checker for the Cliewen method. It keeps a project's documentation — goals, plans, capabilities, decisions — and its agent skills consistent, traceable and versioned.

### Added

- **`clue validate`** checks the documentation: required header fields, IDs that are unique, a valid status, links that resolve, folder READMEs and their generated indexes, where each record came from, and whether every acceptance criterion has a test that references it.
- **`clue version`** reports which release the binary was built from; a build from untagged source says `dev`.
- **Versioned skills.** Every agent skill declares its version, and `validate` fails when the skills disagree with each other or with the binary's release.

### Install

`go install github.com/cliewen/cliewen/cmd/clue@v0.1.0`, or download a prebuilt binary for your platform from the release assets (checksums in `SHA256SUMS`). While the repo is private, `go install` needs `GOPRIVATE=github.com/cliewen` and git authentication for github.com; `gh release download` authenticates through `gh auth login`.
