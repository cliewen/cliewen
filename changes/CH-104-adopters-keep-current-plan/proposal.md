---
id: CH-104
type: change
status: open
links: []
title: P-010 opens — Cliewen makes staying current something an adopter notices
---

# CH-104 — P-010 opens

**Plan-less by construction.** P-009 is `completed` and designated no successor; this change proposes one. A plan-creating change cannot name a plan item it serves.

## What is wrong

Cliewen has a working upgrade path that the people who need it cannot see.

The path is real: re-run the install script to move the binary, then `clue migrate` to preview and `clue migrate --apply` to write. It is stated once, in `guide/operations.md`, on the public site. An adopted repository does not carry that document.

Four gaps follow from this, each verified rather than assumed:

- **No skill mentions upgrading.** The word does not occur once across the five generated skills or their canonical sources. An agent asked to upgrade Cliewen has no routing at all.
- **The emitted `AGENTS.md` never mentions upgrading, migrating, or versions** — and that file is the adopter's, never overwritten by `init` and never rewritten by `migrate`. So routing added there reaches new repositories only. Only skills reach existing adopters, through MIG-003.
- **Nothing tells anyone they are behind.** An old binary with matching old skills is fully green, and `clue migrate` on that pair reports "no changes needed" — which reads as "you are current". The reassurance is false.
- **The drift message does not name the command that fixes it.** `internal/corpus/skillversions.go` says `(drift — reinstall the skills or clue)` at the one moment an adopter learns something must be done, and points at no command.

Meanwhile three backlogs have been accumulating across P-007 through P-009 without a campaign that owns them: 27 decisions still `inferred`, constraints whose enforcement is agent-memory rather than a machine check, and four capabilities that AN-007, AN-011, and AN-013 contradict.

## What changes

**P-010 opens `active` as P-009's successor**, with six sequential milestones in two arcs.

The first arc makes staying current something an adopter notices and can get done: a command that reports a newer release and shows the recipe for the machine it is running on, routing that reaches existing adopters through a skill, and discovery at the moment a coding agent starts or its context is restarted.

The second arc clears the three backlogs the earlier campaigns routed forward and never scheduled.

Milestone numbering continues corpus-global numbering from P-009 at M-045.

## What does not change

The deterministic judge stays offline. `clue validate` reaches no network in this campaign or after it; a milestone that puts a network call anywhere near a validation verdict is out of scope by construction.

The human merge boundary is untouched. Nothing here lets an agent upgrade a repository and merge the result.

Simplification remains outside this campaign and is deliberately deferred to P-011: folding it in would let "we should simplify first" postpone three backlogs that have already waited three campaigns.

## Reversal cost

Cheap and local. Opening, reordering, or closing an unstarted campaign is a status flip, and every milestone routes its own lasting contract through its own ADR or PDR. A dated row in `docs/decisions/log.md` is the right weight, following P-007, P-008, and P-009.
