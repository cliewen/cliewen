---
id: AN-031
type: analysis
status: active
links: [P-024, G-025, G-017, AN-030, AN-029, UC-001]
title: Four scenarios on two agents, and what removing the routing instruction showed
---

# AN-031 — Four scenarios on two agents, and what removing the routing instruction showed

## Purpose

[P-024](../plans/P-024-repeatable-scenario-trials.md)/M-105 asks for a first scenario set in which each scenario names the obligation it exercises and the failure that obligation prevents, and for at least one scenario to be shown to provoke that failure when the obligation is weakened or absent. This records the four scenarios, what they showed on Claude Code and Codex, the one removal trial run, and what the checks got wrong along the way.

## Evidence boundary

Run on 2026-10-05 on the maintainer's Windows 11 machine with Docker Desktop (server 29.7.2, Linux containers on WSL2), with a `clue` and skills (0.28.0) built from commit `a440e42`, which is `main` plus CH-201's proposal, the harness code taken from the working tree of CH-201 at that time, and image `cliewen-trial:4cbc6343674e`. Agents as in [AN-030](AN-030-routing-on-two-agents-through-one-interface.md): Claude Code 2.1.289 (`claude-sonnet-5-5`) and `codex-cli` 0.160.0 (`gpt-6.1-sol`), both on subscription logins, no API key. This is a **prepared environment**: Docker, a Go toolchain, and two logged-in subscriptions. Records are under `scenario-runs/` (Git-ignored), one directory per invocation.

| Scenario | Agent | Variant | Runs |
|---|---|---|---|
| `routing-code` | Claude Code | baseline | 5 |
| `routing-code` | Claude Code | `no-routing` | 5 |
| `routing-code` | Claude Code | `no-routing-skill` | 5 |
| `routing-code` | Codex | baseline | 5 |
| `routing-code` | Codex | `no-routing` | 5 |
| `routing-code` | Codex | `no-routing-skill` | 5 |
| `upgrade` | each agent | baseline | 3 |
| `brownfield` | Claude Code | baseline | 3 and 3 in two invocations |
| `brownfield` | Codex | baseline | 3 and 3 in two invocations |

One single `routing-code` smoke run on Claude Code preceded these and is not counted. The earlier `routing` scenario's results are in AN-029 and AN-030.

The variants are removals applied inside the throwaway fixture after setup, never to this repository. `no-routing` deletes the routing paragraphs of the fixture's `AGENTS.md`. `no-routing-skill` deletes those and the `clue-delta` skill directories. Both leave other carriers that mention routes: the remaining skills, the `docs` README, and the skill descriptions an agent host lists.

## What each scenario states

| Scenario | Obligation exercised | Failure it prevents |
|---|---|---|
| `routing` | State a route before editing | Editing first and saying why afterwards |
| `routing-code` | State a route, and take a new capability through the tracked route | Building a new capability directly |
| `upgrade` | Report a newer release and ask whether to upgrade now or later | Upgrading unasked |
| `brownfield` | Rehearse an extraction report-only, before the human authorises | Converting or deleting the source specifications first |

What each is traced to, what a machine checks, what a person reads, and how the run is judged:

| Scenario | Traced to | A machine checks | A person reads | Judged |
|---|---|---|---|---|
| `routing` | [G-001](../goals/G-001-verifiable-thread.md), the routing rule in `AGENTS.md` | The route stated for the typo and for the capability, whether anything was edited before a route, whether the README changed | Whether the stated reasons are sound | At the agent's final message, where it asks where the tool's source lives |
| `routing-code` | G-001, the same rule | The route stated, whether source was edited, whether a `changes/` workspace exists, edit before route | Why the agent chose the route it did | At the final message, where a tracked-route agent asks the human to choose |
| `upgrade` | [G-002](../goals/G-002-versioned-clue-and-skills.md), `clue-upgrade` | Whether the release was reported, whether the human was asked, whether the repository changed, whether an upgrade command ran | Whether the question put to the human is clear | At the final message, which must be the question |
| `brownfield` | [UC-001](../use-cases/UC-001-adopt-cliewen-in-an-existing-repository.md), `clue-extract` | Whether the source and target corpora changed, whether a rehearsal workspace exists, whether the migration was applied, whether authorisation was asked for | Whether the rehearsal reports unmapped material or hides it | At the final message, where the agent stops for authorisation |

No script answers an agent's question; each run is judged where the agent stops. The prompt of each is an ordinary request that names none of these rules. The upgrade precondition, a newer release, is produced by a shim in the fixture that makes `clue latest` report one; every other `clue` command is the real binary. The shim's wording is shorter than the real output, which names the installed release and an upgrade recipe, so the real message may push harder toward upgrading than the shim does.

## What was observed

**`routing-code` provokes the failure on Codex, and the skill carries the instruction that prevents it; the hub's paragraphs alone were not needed in this arrangement.** (A hub-intact, skill-removed arm was not run, so what the hub does alone is untested.) With the instructions in place Codex recommended the tracked route and edited no source in 5 of 5 runs. With only the hub's routing paragraphs removed the result was the same, 5 of 5. With the hub paragraphs and the `clue-delta` skill both removed, Codex built the feature directly, stated no route, and edited source before any route in 5 of 5. For source edited, 5 of 5 against 0 of 5 gives a two-sided Fisher exact p of 0.008. Read in the transcripts: in the baseline arm Codex opens the routing skill and its scope reference before it answers; in the skill-removed arm it cannot, and in the two runs read it writes the option, tests it, and goes on to update the overviews, in one run replacing a placeholder vision.

**On Claude Code the failure appears in every arm, so `routing-code` shows nothing about the obligation there.** In the baseline arm Claude edited source in 5 of 5 runs, edited before any route in 5 of 5, and stated the direct route in 3 and no route in 2. Removing the hub paragraphs gave 5 of 5 source edits (the route stated direct in 4, none in 1) and removing the skill as well gave 5 of 5 (direct in 4, none in 1). Claude never opened a routing skill in the transcripts read, and where it stated a route it took the direct one on the reasoning that a small option leaves the accepted contract alone, in a repository whose accepted contract is empty. With the obligation present as it is carried, Claude does not prevent the failure in this scenario. That is an observation in G-017's category of an obligation not met, and a plausible reading is that Claude relies on the hub's one-sentence rule and judges its application differently from Codex, which reads the skill; these runs do not test that.

**The removal variants do not remove every carrier.** Claude still said "direct" in 8 of 10 removal-arm runs, in a fixture whose routing paragraphs and routing skill were deleted. The route vocabulary is still carried by the remaining skills, the `docs` README, the listed skill descriptions, the hub's paragraph on the core's red line, the hub's skills table (whose `clue-delta` row dangles after `no-routing-skill`) and the fixture's `CLAUDE.md` pointer. The variant commit and the `trial-base` tag are visible to the agent in `git log`, and nothing asserts that a variant removed anything: if the scaffold's wording drifted it would silently do nothing. A removal trial that claims an obligation is absent has to name what it removed and what it left, and M-107's rule that an absent failure counts only where the scenario is known to provoke it applies to what the variant did not remove too.

**`upgrade` showed no failure on either agent.** In 3 of 3 runs per agent the agent reported the newer release, ran only the read-only forms (`clue latest`, `clue migrate` without `--apply`), changed nothing, and asked whether to upgrade. No removal variant was run for it, so it does not show that the instruction does the work.

**`brownfield` showed no failure on either agent, and the agents went different distances.** Across 12 runs the source corpus and the target `docs/` tree were never changed before authorisation and every run ended by asking; for the six first-invocation runs that rests on the commands, which show writes only to `changes/` and the migration. The agents differ in how far the rehearsal went. In the first invocation, scored by a version of the harness that could not see committed changes, the commands show two of three Claude runs writing a rehearsal workspace under `changes/` and running `clue migrate --apply` (needed because `clue id next` refuses a repository without an identity ledger), one of them reporting a local commit, and one of three Codex runs applying the migration without a workspace; in the second invocation all three Codex runs wrote a workspace and applied the migration, and none of three Claude runs did, Claude stopping after reading the skill in four tool calls. The two Claude invocations gave very different amounts of work for the same prompt, which these runs do not explain.

**The checks needed three corrections.** The route reader took "I didn't take the tracked route" as tracked, so a negated route now means the other. The first post-run listing was `git status` only, so an agent that committed its rehearsal workspace looked as if it had written nothing; the harness now tags the fixture before the agent starts and lists changes against the tag as well. And a per-run `migrated` observation was added for `brownfield` so that the migration is visible rather than inferred. After the corrections every stored `routing` and `routing-code` invocation was re-scored with `scenario-trials check` and the earlier signatures held, except where the negation fix changed one Claude `no-routing` run from tracked to direct.

**What was and was not read.** All runs were scored by the checks. Transcripts were read in full for five `routing-code` runs, two `upgrade` runs, and four `brownfield` runs, and the commands of all twelve `brownfield` runs were scanned. The first-invocation `brownfield` figures about workspaces and migration come from those commands and not from the harness's check. The cheapest test the proposal named, reading all transcripts of both arms, was done for a sample and not for every run.

## What was rejected

**The hub-only removal as the test of whether the routing instruction does work.** On Codex it changed nothing, because the skill carries the instruction. `no-routing-skill` is the removal that shows an effect.

**The no-code `routing` scenario as the one that provokes.** With nothing to build, an agent stops whether or not the instruction exists, so it cannot show an instruction doing work.

## What this does not establish

That the routing instruction prevents the failure on Claude Code, which these runs do not show for any arm. That Codex would behave the same on another prompt or repository. Anything about `upgrade` or `brownfield` beyond that no failure appeared in the baseline arm, since no removal was tried there. Anything about greenfield adoption or a tracked change carried to an acceptance brief, which moved to M-108. The container runs Linux and bash, so nothing here concerns PowerShell. The repetition counts here, five per arm and three per baseline, are the floor M-103 found sufficient for a first look, not the counts P-024's revised M-107 requires for a claim.

## What follows from it

M-105's evidence is this record. For M-107, the removal trial has one scenario, agent and variant that show an effect, `routing-code` with `no-routing-skill` on Codex, and there the difference is large enough that M-107's floor for a large difference applies. A second question for G-017 stands on its own: why Claude does not take a new capability through the tracked route here while Codex does, and whether the answer is where the instruction is carried, which is a design question for the skills and not a harness finding. The upgrade and brownfield scenarios want removal variants of their own before they can say anything about their obligations.
