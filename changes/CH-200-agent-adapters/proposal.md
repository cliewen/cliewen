---
id: CH-200
type: change
status: open
links: [P-024, M-104, G-025, AN-029, AN-028]
title: Two agents run the same scenario through one adapter interface
---

# CH-200 — Two agents run the same scenario through one adapter interface

## Proposal

Serve P-024's M-104. Turn the single Claude Code function in `tools/scenario-trials` into an adapter interface and add Codex as the second adapter, so the same routing scenario runs on both from a clean container home. Each adapter owns what is specific to its agent: how it logs in, the command that runs it, how its event stream becomes a vendor-neutral transcript, and the probe that reads its version and model back from inside the container. The checks, the conditions record and the spread summary read only the neutral transcript, so a vendor name appears nowhere outside its own adapter file.

The change records two decisions in IDR-010: agents are reached through adapters over one interface, and a terminal-driving orchestrator such as Orca is not used, the choice AN-028 set aside on a two-page reading and M-104 asks to be recorded. It also carries the plan revision P-024's challenge promised: M-107's repetition count is stated in terms of what the claim needs and fixed before the run, with AN-029's own arithmetic as the reason, instead of "repeated enough to read the spread M-103 found". The revision is declared in P-024's Declared revisions and rides with this change, as the planning workflow allows when an implementing change uncovers it.

The first Codex runs are read against the checks before any claim is made, and what they show is recorded in AN-030.

## Challenge

**The assumption most likely to undermine the work:** that Codex's event stream carries enough for the same checks to mean the same thing. It has no model name, no turn count, and no cost in the stream the way Claude's does, and its command events wrap each command in a login shell. If the checks read a neutral transcript but the Codex mapping silently loses the ordering of "said a route" and "edited", the two agents would be compared on different measurements and the comparison would be an artefact of the mapping.

**A credible alternative:** a second copy of the Claude path with Codex-specific parsing, no interface. It is faster to write and nothing else changes. It keeps vendor names in the checks and the runner, which is exactly what M-104's exit criterion forbids, and the third adapter (M-106) would be a third copy.

**The cheapest useful test:** run Codex once, read its transcript by hand next to the neutral transcript the adapter produced, and compare the order of text and tool steps. Do this before running more than one Codex run. If the order differs, fix the mapping first.

**What would revise or stop the work:** a Codex stream that cannot give the order of text and edits, in which case the adapter reports "not comparable" for the edit-before-route check and the plan says so, rather than reporting a value it cannot support; or a login that cannot be supplied from outside the repository without writing into an image layer or a log, which is M-104's own exit criterion.

**How a build could meet every criterion and still fail the maintainer:** an interface so thin that adding a third agent needs a change in the runner, and an agent that can read its own credential from the container. The second is true of both adapters, since the agent runs as a user who can read the file or environment variable that logs it in, and it is stated, not hidden: the credential is the maintainer's own subscription login in a throwaway container, and the container is the boundary.

## Scope boundary

This change adds an adapter interface and the Codex adapter and runs the same scenario on both. It does not add scenarios (M-105), the OpenCode adapter (M-106), or a removal trial (M-107, beyond the declared revision of its repetition count). The Codex agent runs with its own sandbox off because it cannot start inside the container, and the container is the boundary. No new acceptance criterion is added.
