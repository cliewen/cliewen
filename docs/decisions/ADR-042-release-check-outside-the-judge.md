---
id: ADR-042
type: decision
status: inferred
links: [P-010, CAP-002, CAP-004, ADR-011, ADR-030, ADR-039, ADR-040, C-004, C-013]
title: A release check reaches the network, reports, and writes nothing
author: agent
accepted-by: []
---

# ADR-042 — A release check reaches the network, reports, and writes nothing

## Context and problem statement

A repository running Cliewen cannot tell whether it is current. A binary and its skills installed together stay internally consistent forever: the version stamps agree, so [ADR-011](ADR-011-version-stamping.md)'s drift rule is silent, and `clue migrate` reports `no changes needed` because the binary in hand knows of no migration past itself. Both answers are correct and neither is the one the adopter wants. A release cut yesterday is invisible to every command that ships.

The question "is there something newer?" cannot be answered from files in the repository. It has exactly one source — the published release list — and that source is on the network.

The judge is the one place the answer can never live. A verdict that depends on another system's present state is not deterministic, is not reproducible against a pinned revision, and is not available offline; and it would let someone else's outage condemn a corpus that has not changed. [ADR-040](ADR-040-qualified-external-references.md) settled that argument for external references and put resolution in a separate advisory command. This decision states the same boundary for the release check, and states the second boundary that one does not need: a command that reports a newer release sits one short step from a command that installs it.

## Decision outcome

**A release check is a separate command that reaches the network, reports what it found, and changes nothing.**

*It is outside the judge, permanently.* `clue validate` reads files in the repository and nothing else. The release check never contributes to a validation verdict, never runs as part of one, and — like `clue refs` — must never be a required status check. Making it one would put the release host's uptime between a change and its merge.

*It reports; it does not update.* The command writes no file in the repository, with or without flags, and it does not replace the binary it is running as. The published install scripts verify every download against the release's `SHA256SUMS` and refuse to install on a mismatch ([ADR-030](ADR-030-verified-install-scripts.md)); a self-updating binary would have to reimplement that verification inside itself, and a program that can overwrite itself from the network is a different kind of program than an offline judge. The upgrade stays what [ADR-039](ADR-039-versioned-corpus-migrations.md) made it: a reviewed change, moving the coordinated set, accepted by a human.

That is the whole line between a reporting command and a background updater, and it is drawn where it can be checked rather than promised: no repository write, no self-replacement, no verdict.

*It names the route for the machine it is running on.* Cliewen publishes three installation routes — the PowerShell script on Windows, the shell script on macOS and Linux, `go install` where no prebuilt asset exists. Printing all three would make every reader skip past two wrong lines to find theirs, and the two wrong ones are wrong in a way that costs time to discover. The command resolves the platform itself and prints one. It then prints the coordinated sequence, because moving the binary alone *produces* the drift report rather than resolving it: the repository moves with `clue migrate`, previewed and then applied inside a reviewed change.

*Every degradation is silence and success.* Offline, a timeout, a rate limit, and a response the command does not recognize all mean the same thing — *could not tell* — and not being able to tell is not a defect in the repository. Each exits 0 and writes nothing to standard error. A check that complains when the network is down is a check that gets removed, and this one is meant to run where a failure is expensive: in a session hook, on every start. Reporting "could not reach the release list" on the ordinary output is honest and calm; failing is neither.

*A quiet mode says one line when behind and nothing at all when current.* Where the output is a coding agent's context, every line is spent on something. A check that greets a current repository every morning teaches its reader to ignore it.

*The answer is cached outside the repository.* The repository holds reviewed content; a cache is machine state, and committing "what the release list said on Tuesday" would put a fact with an expiry date under review. It lives in the user's cache directory with a bounded lifetime, so repeated calls inside that window cost no request — which is what makes running it on every session start acceptable to the host as well as to the user. An unreadable, corrupt, or unwritable cache is treated as absence: ask again, or stay silent. It is never an error, because a cache that can fail a command is worse than no cache.

*The drift message names the way out and the way to stay.* ADR-011's rule is unchanged; only what it says changes. A message that states a disagreement and stops leaves the reader to guess which side to move, and silently implies that moving is the only option. It names the command that reports the upgrade, the command that moves the repository, and the route for a repository that is deliberately staying where it is.

## Rejected: fold the check into `clue validate`

The convenient design, and the one ADR-040 already rejected for external references. It puts the network in the deterministic judge, makes a verdict depend on another system's present state, and stops validation working offline. Nothing about the release check makes that trade better than it was.

## Rejected: fold the check into `clue version`

`clue version` prints what this binary is. It is the one command guaranteed to answer instantly, offline, and identically forever, and it is used in scripts that assume exactly that. Adding a network flag to it means the guarantee now depends on which flags were passed, and a script that reads its output inherits a timeout it never asked for.

## Rejected: a self-updating binary

Rejected in [P-010](../plans/P-010-adopters-keep-current.md) before this decision and kept rejected here. The convenience is real and it is delivered by an agent following a skill, under the human merge boundary, not by a judge that can overwrite itself.

## Rejected: cache the answer in the repository

It would survive a clone and be shared across a team, which sounds like a feature until the file is a reviewed artifact recording a fact that expires. It would appear in diffs, invite merge conflicts about nothing, and make a green repository depend on when someone last ran a network command.

## Carrier

CAP-004 owns the command's shipped behaviour and the drift message's content; CAP-002 keeps the judge offline. Focused positive and negative tests cover all three installation routes, the quiet mode in both directions, every degradation, the cache's lifetime and its unreadable state, and the drift message. The network and the platform are both injected, so no test reaches a live service and no verdict depends on which machine ran it.

The standing rejections this decision does not reopen: the judge stays offline and deterministic, no forge state becomes corpus meaning, and no command accepts a change on a human's behalf.
