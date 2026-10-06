---
id: AN-033
type: analysis
status: draft
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

## Observations

Not yet run.
