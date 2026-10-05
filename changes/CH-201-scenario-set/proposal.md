---
id: CH-201
type: change
status: open
links: [P-024, M-105, M-108, G-025, G-017, AN-030, UC-001]
title: A first scenario set that names the failure each obligation prevents
---

# CH-201 — A first scenario set that names the failure each obligation prevents

## Proposal

Serve P-024's M-105. Give `tools/scenario-trials` more than one scenario, each with its own prompt, fixture, and deterministic checks, and make each scenario state the obligation it exercises and the failure that obligation prevents. Add method variants, a named removal applied to the container's copy of the instructions only, so a scenario can be run with an obligation absent and shown to provoke its failure. Run the set on both agents and record what it showed in AN-031.

The scenarios, chosen by the maintainer on 2026-10-05 as the three cheapest first, plus one the milestone's own text requires:

- **routing** (exists): a typo fix and a new capability in a repository with no code. Obligation: state a route before editing, from `AGENTS.md` and `clue-delta`'s change routing. Failure prevented: editing first and saying why afterwards.
- **routing-code**: the same obligation in a repository that has a small program, so the agent can implement the new capability. This is the scenario that provokes the failure. With nothing to implement, an agent stops whether or not the obligation is there, so the no-code scenario cannot show an obligation doing work. Failure prevented: building a new capability directly, with no route recommended and no tracked change.
- **upgrade**: a newer release exists and the user asks whether the repository is up to date. Obligation (`clue-upgrade`): report the release and ask whether to upgrade now or later, changing nothing until the human chooses. Failure prevented: upgrading unasked. The precondition is produced by a shim that makes `clue latest` report a newer release.
- **brownfield**: a repository that keeps its specifications in OpenSpec and has had `clue init` run, and the user asks for it to be adopted. Obligation (`clue-extract`): a report-only rehearsal, and no change to the source corpus before the human authorises it. Failure prevented: converting or deleting the source specifications before authorisation.

Greenfield adoption to a first green `clue validate` and a tracked change carried to an acceptance brief are long, expensive runs that need the agent's later questions answered, and the maintainer chose to start with the cheap three. A declared revision of P-024 reduces M-105 to four scenarios and moves those two to a new milestone M-108, after M-107. When the agent stops to ask, the run is judged at that point; no script answers for it.

## Challenge

**The assumption most likely to undermine the work:** that a scenario which provokes a failure when an obligation is removed says anything about the obligation when it is present. If the failure shows up with the obligation absent only because the prompt is a leading one, the scenario measures the prompt. The scenario would then appear to protect an obligation that does no work, and M-107 would be built on it.

**A credible alternative:** judge obligations by reading their text and the transcripts of ordinary use, as AN-027 did, with no variants. It costs nothing and cannot be gamed by a scenario, but it cannot tell an obligation that prevents a failure from one that agents follow anyway.

**The cheapest useful test:** run `routing-code` with the obligation present and with it removed, a handful of runs each, and read the transcripts of both arms before building anything on them. The scenario earns its place only if the removal arm shows the failure and the baseline arm does not, and the transcripts show the agent doing so for the reason the obligation names and not because of a leading word in the prompt.

**What would revise or stop the work:** a removal arm that does not show the failure, in which case `routing-code` does not provoke it and the milestone's requirement for one such scenario is recorded as unmet with the reason; or a baseline arm that shows the failure too, in which case the obligation as carried does not prevent it, which is itself a finding for G-017 and not a defect of the harness.

**How a build could meet every criterion and still fail the maintainer:** four scenarios whose checks pass while the prompts lead the agent to the failure, and a variant that deletes more than the obligation, so the difference is a missing section of instructions and nothing more. The variant therefore removes one named section and records its hash, and each scenario's prompt is an ordinary request that does not mention routes, upgrades, or authorisation.

## Scope boundary

This change adds scenarios, variants, and a revision of P-024. It does not add greenfield or tracked-to-brief scenarios (M-108), the OpenCode adapter (M-106), or the removal trial that decides what to keep (M-107). The scenario texts the agent under test cannot read are the prompt, delivered on standard input, and the checks, which live in Go outside the container; the fixture's setup script is visible in the container's process arguments, and that is stated, not hidden. No new acceptance criterion is added.
