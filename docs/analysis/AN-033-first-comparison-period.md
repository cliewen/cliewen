---
id: AN-033
type: analysis
status: active
links: [P-024, G-025, G-017, AN-027, AN-031, AN-032]
title: A first comparison period across three agents, fixed before observing
---

# AN-033 — A first comparison period across three agents, fixed before observing

## Purpose

[P-024](../plans/P-024-repeatable-scenario-trials.md)/M-107 asks for a first comparison period, recorded as [AN-027](AN-027-baseline-observation-first-period.md) was: population, conditions and repetition fixed before any run, the scenario set run across the agents, each run read against [G-017](../goals/G-017-the-method-is-evidenced-and-simplified.md)'s categories, and the harness kept, revised or dropped on what that shows. It also asks for at least one removal trial with the number of runs per arm fixed in advance.

This document is committed with the protocol alone. Everything under "Methodology, fixed before observing" is binding for the runs that follow; results, readings and the decision are added later under their own headings, and the commit that adds the first of them comes after this one in `git log`.

## Methodology, fixed before observing

**Population.** Three agents by four scenarios, all at baseline: Claude Code (2.1.289, `claude-sonnet-5-5`), Codex (`codex-cli` 0.160.0, `gpt-6.1-sol`) and OpenCode (1.18.34, `openrouter/qwen/qwen3.8-27b`), on the scenarios `routing`, `routing-code`, `upgrade` and `brownfield`. A fresh set of runs, not those of [AN-029](AN-029-five-routing-runs-on-claude-code.md) to [AN-032](AN-032-one-opencode-run-on-a-third-model.md), because `clue` and its skills have changed since. Every invocation is built from one `clue` commit, the `main` commit at the merge of the pull request that carries this protocol, which the results section names; the agent versions are pinned by the Dockerfile.

**Conditions.** One maintainer's Windows 11 machine, Docker Desktop with Linux containers, subscription logins for Claude Code and Codex and an OpenRouter key for OpenCode, an empty agent home in every run. Reasoning effort: OpenCode at `medium`, Claude Code and Codex at their defaults, each recorded in `conditions.json`. The container runs Linux and bash, so nothing here concerns PowerShell. Fixtures are the harness's own and are small; a result describes these scenarios and not the method in general.

**Runs per arm, fixed now.**

| Cell | Runs | Basis |
|---|---|---|
| Baseline, each agent and scenario | 5 | The spread [AN-029](AN-029-five-routing-runs-on-claude-code.md) found enough for a first look; no claim about a rate is made from five runs |
| `routing-code` on Codex, baseline and `no-routing-skill` | 10 and 10 | A large expected difference, 5 of 5 against 0 of 5 in AN-031, so the floor of 10 |
| `routing-code` on OpenCode, baseline and `no-routing-skill` | 20 and 20 | The size of the effect is unknown, so about 20; the runs cost about $0.05 each, and the arms answer the question AN-032 left open, whether OpenCode finds and uses the repository's skills |

The Codex and OpenCode `routing-code` baseline arms are the baseline cells for those agents and are not run twice. About 110 runs in all. `routing-code` is run on Claude Code at baseline only: the failure appeared in every arm there in AN-031, so a removal arm could show nothing. `upgrade` and `brownfield` have no removal variant and none is built for this period; their obligations are recorded as not exercised at best.

**What each run is read for.** The deterministic checks score every run. A person reads, in full: every run whose signature differs from the most common one in its cell; two further runs per baseline cell, the two whose run numbers sort first by SHA-256 of `107-<invocation directory name>-<run number>`; and five per removal arm chosen the same way. The analysis states the number read against the number scored. Each run read is coded against G-017's categories: an obligation not met, a bypassed obligation with the reason the agent gave, documentation that misled (including what a removal left behind, as AN-031 found), and a consequential misunderstanding. Acceptance cost and rediscovery are not observable in these scenarios, and the analysis says so instead of reporting zero.

**Decision rules.**

- An obligation is proposed for removal only where a scenario known to provoke its failure ran without it and the failure did not appear. A scenario is known to provoke it only if the failure appeared in some arm of this period or of AN-031. Otherwise the result is recorded as not exercised, never as unneeded. A proposal goes to the tracked route as its own change.
- Stated before the run: `routing-code` on Codex again shows the effect; `routing-code` on OpenCode shows the effect or does not; nothing is proposed for removal, because only Codex has an arm in which the failure is known to appear and there it appears without the skill.
- The harness is **kept** if the period showed at least one difference between agents, or one fact about skill discovery, that hand observation would not have found for a similar cost. It is **revised** if its checks needed correction during the period; each correction, and each harness fault that stopped or spoiled a run, is counted and named, and the three check corrections and the mount fault of AN-031 and AN-032 are the point of comparison. It is **dropped** if the cost of the period per finding is higher than hand observation would have been, with the cost counted as runs, hours of attention and tokens.
- A skill or obligation changes as a result only through its own tracked change; this analysis states whether the period gave a reason to start one.

**What would make the period uninformative.** Spread so wide in a baseline cell that its modal behaviour is not a majority; checks that disagree with the transcripts read in more than one run in ten of those read, in which case the checks are corrected and every stored run re-scored before anything is concluded; or a harness fault that spoils a cell, which is re-run and recorded.

## Amendment before the Codex runs, on the maintainer's direction

Made after the protocol merged and after the OpenCode `routing-code` arms had started, and before any Codex or Claude Code run. The maintainer chose the narrower period, because the broad baseline cells were expected to show no failure and so to decide nothing about whether the harness earns its upkeep.

The period is now two removal trials: `routing-code` on OpenCode (20 baseline, 20 `no-routing-skill`) and on Codex (10 and 10), 60 runs, from the `main` commit `7498c4a`. The other baseline cells, every scenario on Claude Code and `routing`, `upgrade` and `brownfield` on Codex and OpenCode, are **not run in this period**, and nothing here says anything about them. The conditions, the reading protocol, the decision rules and the expected outcome above are unchanged, except that the "differs from the most common one" and the two-per-baseline-cell reading apply to the two trials' baseline arms only. A comparison across agents of `upgrade` and `brownfield` stays open as a later period.

## Observations

Run on 2026-10-06 on `main` at `7498c4a`, with `clue` and skills 0.28.0, Claude Code not run, Codex 0.160.0 and OpenCode 1.18.34 as pinned, in the four invocations under `scenario-runs/` (Git-ignored) named `20261006-174807`, `-175213`, `-175835` and `-181247`. Efforts as recorded in each `conditions.json`: OpenCode `medium`, Codex `default`. All 60 runs were scored by the checks. Cost: $0.96 for the 40 OpenCode runs; Codex runs are on a subscription and report none.

| Arm | Runs | Source edited | Edit before any route | Route stated |
|---|---|---|---|---|
| Codex, baseline | 10 | 0 | 0 | `tracked` 10 |
| Codex, `no-routing-skill` | 10 | 9 | 9 | none 10 |
| OpenCode, baseline | 20 | 19 | 10 | `direct` 18, `tracked` 1, none 1 |
| OpenCode, `no-routing-skill` | 20 | 20 | 18 | none 16, `tracked` 3, `direct` 1 |

**Codex shows the effect again, as stated before the run.** With the skill, 10 of 10 chose the tracked route and edited no source. With it removed, 9 of 10 edited source before any route and none stated one. Source edited, 0 of 10 against 9 of 10, gives a two-sided Fisher exact p of about 0.0001. In the six removal runs read, five built the feature and wrote architecture and design documentation; the sixth stopped and asked for the missing `clue-delta` skill to be restored, because the hub still points at it. In both baseline runs read, Codex opened the `clue-delta` skill or its scope reference before it answered.

**OpenCode on `qwen3.8-27b` shows the failure in both arms, so the obligation is not exercised for it.** With the skill present, 19 of 20 runs edited source and 18 stated the direct route; with it removed, 20 of 20 edited, so the difference is nothing (Fisher p = 1). The reasons given for `direct` in the runs read are the ones Claude gave in AN-031: the corpus holds no accepted contract, so an addition touches none. In none of the 40 runs did a tool call name a skill path: OpenCode did not open a skill in either arm, so the skill's presence could not have been what changed anything. That answers the question [AN-032](AN-032-one-opencode-run-on-a-third-model.md) left open for this scenario: on this model the repository's skills are not used, whether or not OpenCode lists them. The one visible effect of the removal is in what the agent says: it stated a route in 19 of 20 baseline runs and in 4 of 20 with the skill removed, which fits the hub's one-sentence rule doing some work on the wording.

**A harness fault spoiled 12 of the 40 OpenCode runs.** OpenCode asked permission to write under `/tmp`, treated as an external directory, and in a non-interactive run auto-rejected it (`permission requested: external_directory (/tmp/*); auto-rejecting`, in `stderr.txt` of all 12). The agent had announced a test of a file write and the run ended there, unfinished: 7 of 20 in the baseline arm and 5 of 20 in the removal arm. The scored measure, source edited, is not affected where the run had already edited, and the one baseline run that edited nothing finished normally, so no result above rests on an unfinished run having stopped early; but unfinished runs may have stopped before a route was stated, so the route counts for OpenCode are lower bounds. This was not caught by the checks. It is the period's one harness fault.

**What was read.** 25 of the 60 runs: every run differing from the most common signature of its arm (10 in the OpenCode baseline, 1 in the Codex removal arm, 4 in the OpenCode removal arm, none in the Codex baseline, with overlap), plus the two runs of each baseline arm and five of each removal arm that the fixed ordering gave. I read, per run, the list of tool calls, whether any named a skill path and the agent's last message, not each full transcript, so what the agents reasoned in between is not recorded here. The checks and the reading disagreed in none of the 25.

**G-017's categories.** *Obligation not met:* OpenCode in 19 of 20 baseline runs. *Bypass, with the reason:* "no accepted contract, so the addition is direct", in the OpenCode runs read that stated one. *Misleading documentation:* after `no-routing-skill` the hub still points to the `clue-delta` skill, which is a property of the variant and not of the method (AN-031 found the same dangling row); one Codex run treated it as a blocker and stopped. *Consequential misunderstanding:* none observed. Acceptance cost and rediscovery are not observable in these scenarios.

## What this does not establish

That the obligation is unneeded on OpenCode: it is not exercised, because the failure is present with the skill in place and the skill was never opened. Anything about `upgrade`, `brownfield` or `routing` on any agent in this period, or about Claude Code, which this period did not run. That Codex would behave the same on another prompt, repository or model; one scenario, one fixture. That the effect on Codex comes from the skill alone: the variant also removes the hub's routing paragraphs, and AN-031 found the hub paragraphs alone changed nothing, but an arm with only the skill removed was not run here. Anything about PowerShell.

## Decision

**No obligation is proposed for removal.** Codex is the one agent where the failure is known to appear without the skill, and there it did; on OpenCode the failure was already present, so the removal rule's "the failure did not appear" is not met anywhere. No skill or obligation changes as a result of this period.

**The harness is kept, with one revision.** It showed what hand observation would not have found for the same price: a clean 0 of 10 against 9 of 10 on Codex, and 40 OpenCode runs, for $0.96, in which no skill was opened. It also carried a fault the checks did not catch (12 of 40 OpenCode runs ended unfinished on the `/tmp` permission), which is the revision: OpenCode's non-interactive permissions for the container need setting so a run does not end on a harmless write, and the checks should report an unfinished run as its own signature. That follow-up is a harness change, not done in this period; until it is, an OpenCode result is read with the unfinished count beside it. One fault in 60 runs and none in the checks' readings of the 25 read is fewer corrections than M-105 and M-106 needed.

The cross-agent comparison of `upgrade`, `brownfield` and `routing`, and Claude Code, stays open as a later period; whether it is worth its cost is better judged after the revision above.