# The skills

A skill is a folder of Markdown instructions that a coding agent loads when a task matches it. Cliewen puts its process knowledge in six skills. The `clue` binary stays a small judge that only checks the repository, and your prompts can stay in ordinary words. You rarely name a skill yourself. You describe the work, and the agent picks the skill.

## How the agent finds its way

`clue init` writes the files an agent reads before it does anything:

| File | What it is for |
|---|---|
| `AGENTS.md` | The routing hub at the repository root. Most coding agents read this file name. It says what to do first and which skill fits which task. |
| `.agents/skills/clue-*/` | The six skills. Each has a short `skill.md` and a `references/` folder the agent reads only when a step needs it, so a small task does not load the whole method. |
| `.claude/skills/` and `CLAUDE.md` | A mirror of the skills and a file that imports `AGENTS.md`, because Claude Code reads those names instead. Both only point at the hub, so every rule stays in one place. |

`AGENTS.md` and `CLAUDE.md` are yours to extend with your repository's own conventions. `clue init` never overwrites them and `clue migrate` never rewrites them.

Every session starts the same way. The agent runs `clue latest --quiet`, which prints one line only when a newer release exists. Then it reads the smallest relevant part of the corpus and, before editing, tells you which route it recommends, why, and what it might still discover that would change the recommendation. If you never see a sentence starting `Recommended route:`, the agent has not read `AGENTS.md`, so point it there.

**Simple** work keeps every promise the repository has already made. Examples are a bug fix that makes the code meet an existing acceptance criterion again, refactoring, maintenance, configuration within what was agreed, and editing prose. The agent makes the change and runs the checks that apply to it.

**Full** work changes a promise: a new or changed acceptance criterion, a capability, a decision, a policy, a plan's goals, or behavior no criterion covers yet. The agent uses the full change loop, with a written proposal and a pull request that you merge. When the agent is unsure, it recommends full. You can overrule it, and it then records your choice and the risk in the commit message. [The change loop](./change-loop) walks through both routes.

## The six skills

| Skill | The agent uses it when | What you get |
|---|---|---|
| `clue-plan` | A goal needs several steps, or a running plan must change what it promises | A plan with milestones you can check off |
| `clue-analysis` | Something is too uncertain to plan or build yet | A findings document that a plan or change then uses |
| `clue-delta` | You accept a full route | A branch, a proposal, the code, tests, and updated `docs/`, in one pull request |
| `clue-verify` | A full change is about to be marked ready for review | A checked and independently reviewed pull request |
| `clue-extract` | An existing repository adopts Cliewen, once | A reviewed first corpus built from your existing specifications and tests |
| `clue-upgrade` | A newer Cliewen release exists | The repository moved to that release in one reviewed change |

### `clue-plan`: shape a campaign

Say something like `Plan how we get to <outcome>.` The agent first runs `clue validate --intent` to see what direction the corpus states. If the repository has no vision yet, it offers to work one out with you. It then writes a plan under `docs/plans/` whose milestones each have an exit criterion someone can check. A plan goes through a pull request like any other change.

It asks you whether the goal is right, whether the milestones are in the right order, and what to do when a plan stops making sense.

### `clue-analysis`: find out first

Say `Before we build this, find out whether <risk>.` The agent runs a time-boxed investigation (a spike) such as a prototype, a measurement, or a reading of the sources. It ends with a findings document under `docs/analysis/` that records what was tried, what was rejected and why, and the revisions it looked at. It does not write findings that no plan or change will use.

### `clue-delta`: make a full change

This is the skill behind most real work. After you accept a full route, the agent:

1. branches from `main`;
2. writes `proposal.md` in a temporary `/changes/CH-xxx-slug/` folder, pushes it, and opens a draft pull request;
3. changes the code, the tests, and `docs/` together;
4. folds what the change means into `docs/` and deletes the temporary folder (the *digest*);
5. hands the result to `clue-verify`.

It stops and asks you when a question blocks the work, when a plan no longer holds, and at the end, because only you merge.

### `clue-verify`: check before ready

Before a full pull request is marked ready, the agent runs the repository's tests and `clue validate --forbid-changes` on the exact commit. Where the agent host supports it, it starts a fresh reviewer that sees the declared intent but not the implementation conversation. Blocking findings are fixed and reviewed again. The pull request opens with an *acceptance brief*: one screen saying what merging would accept, with a verdict on whether each changed criterion's tests really check its scenario.

### `clue-extract`: adopt an existing repository

Say *"Bring this repository into Cliewen."* The first result is a report-only rehearsal: an inventory of what exists and a proposed mapping. Nothing in the repository changes until you say so. The corpus it then builds keeps your existing IDs and test links. Everything it inferred is marked `inferred` until a human confirms it. [Greenfield and brownfield](./adoption) covers the whole path.

### `clue-upgrade`: take a new release

When `clue latest` reports a newer release, the agent tells you and does nothing until you say *now*. On a yes, it previews `clue migrate` and shows you how to install the new binary on your machine. It then moves the skills, the CI workflow, and any corpus obligations to that release together in one branch. An upgrade is simple work, because the release's own changes were reviewed before it was published. If the preview reports analysis documents that may have served their purpose, the agent asks you about them one at a time and retires any you approve in a separate change.

## Habits the skills build in

A few rules run across several skills. Knowing them helps you read what the agent writes.

### It challenges a big commitment before making it

Before a plan's promise is adopted or revised, or before the proposal for a change that would be expensive to get wrong, the agent writes a short challenge where you will review it, in the plan or in `proposal.md`. The challenge names the assumption most likely to sink the work, a credible alternative, the cheapest test of that assumption, and the result that would stop or change the work. It also asks how an implementation could meet every criterion and still fail the person it is for. Tests written in the same sitting as the criteria cannot catch that failure.

Simple work skips this. So does work that follows a course already challenged, or that could be undone within the same change for less than the challenge would cost. The agent adds no note to explain the skip.

### It checks that the plan still holds

Before starting or resuming a milestone, the agent checks that the plan still serves its goal and that the milestone is still wanted and achievable. A passing check leaves no record. A failing one pauses the affected work, and the agent asks you to choose how to replan.

### It keeps work out of private memory

An agent's memory does not survive a new session or a different agent, so nothing important may live there. A suggestion made during a change becomes a task in the change if the change would be incomplete without it, and otherwise a goal with `status: proposed` for later.

A useful discovery gets written down only when it cost real effort to find and a fresh agent doing similar work would probably hit it again, for example an undocumented flag or a non-obvious recovery step. The agent first asks whether the confusing step can be fixed instead. If not, the note goes into the relevant capability's `design.md`, or into your own contributor guide such as `CONTRIBUTING.md` when it is about running the repository's tooling. The note claims only what was actually observed.

### It asks whether merging is really protected

Branch protection is a setting at your Git host, so no file in the repository records it and `clue validate` cannot see it. Before the first pull request in a repository is marked ready, the agent asks the host what the branch enforces. If something is missing, it tells you what and offers the exact commands, and changes nothing without your go-ahead. When it cannot ask the host, it reports the state as unknown rather than assuming you are covered. [Enforce CI](./ci-wall) has the settings.

## The skills are versioned files you own

Each generated skill carries a `version:` stamp that matches the `clue` release it came from. `clue validate` reports when the skills disagree with each other or with a released binary, and ignores skills that are not Cliewen's. Do not edit the generated skills by hand. `clue migrate` replaces them when you upgrade. [Operate safely](./operations#upgrade-one-coordinated-set) explains how to move the binary, skills, and CI workflow together.

Because the skills are committed files, the guidance that shaped a branch is part of that branch's history. You can see what the agent was told on the day it did the work.

## Next

[See the full change loop step by step.](./change-loop)
