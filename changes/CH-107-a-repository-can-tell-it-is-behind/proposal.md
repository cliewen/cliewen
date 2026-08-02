---
id: CH-107
type: change
status: open
links: [P-010, M-045, CAP-002, CAP-004, ADR-011, ADR-039, ADR-040, C-013]
title: A repository can tell it is behind, and is shown how to catch up on the machine it is running on
---

# CH-107 — A repository can tell it is behind, and is shown how to catch up on the machine it is running on

**Serves P-010's M-045**, the first milestone of that campaign's first arc, and it implements the milestone's exit criterion in full rather than a part of it.

## What is wrong

Nothing in Cliewen distinguishes *being up to date* from *being unable to tell*.

An adopted repository whose binary and skills were installed together stays internally consistent forever. `clue validate` passes: the skill stamps agree with each other and with the binary, which is exactly what [ADR-011](../../docs/decisions/ADR-011-versioned-clue-and-skills.md)'s drift rule is for. `clue migrate` on that same pair reports `no changes needed` — a sentence an adopter reasonably reads as "you are current", when all it says is that the binary in hand knows of no migration past itself. A release cut yesterday is invisible to both. The corpus can be a year behind and every command is green.

The one moment a repository does learn something is wrong is the drift message, and it names no way out: `skill version 0.11.2 != clue 0.12.0 (drift — reinstall the skills or clue)`. It states the fact and then stops. "Reinstall" is not a command; which of the two sides to move is not said; and the legitimate answer of *staying* on the older release — pinning — is not mentioned at all, so the message reads as an instruction to upgrade when it is actually a report of disagreement.

There is also no place a machine could learn the answer without violating something. The obvious home for "is there a newer release?" is the judge, and the judge is the one place it can never live: a verdict that depends on another system's present state is not deterministic, not reproducible against a pinned revision, and not available offline. [ADR-040](../../docs/decisions/ADR-040-qualified-external-references.md) already settled that argument for external references and put resolution in `clue refs`, outside the judge and outside branch protection. The release check is the same shape of problem and needs the same answer stated for itself — a reporting command is not a background updater, and the difference has to be written down before either exists.

## What changes

1. **An ADR states where the network may go.** Why a release check may reach the network while `clue validate` never may; why the answer is advisory and never a required check; and what separates a *reporting* command from a *background updater* — the check writes no file anywhere in the repository, with or without flags, and never replaces the binary. It extends ADR-040's two-layer boundary to a new subject rather than reopening it. This touches what `clue validate` may depend on, a core carrier under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md), so it is an explicit decision record with human acceptance.

2. **`clue latest` reports whether a newer release exists, and how to get it on this machine.** It prints the running version, the newest published release, and — when behind — the installation route *for the platform it is actually running on*: the PowerShell script on Windows, the shell script on macOS and Linux, `go install` where no prebuilt asset exists. One text that named all three would make every adopter read past two wrong lines to find theirs. It then prints the coordinated sequence — `clue migrate` to preview, `clue migrate --apply` to write — because moving the binary alone is what produces the drift report, not what resolves it.

   The name is deliberate. `clue release` would read as "cut a release" in a project that cuts them, and folding the check into `clue version` would put a network call inside the one command guaranteed never to make one.

3. **`--quiet` prints one line when behind and nothing at all when current.** M-047 will run this inside a coding agent's session start, where every line of output is context spent. A check that greets a current repository every morning is a check that gets deleted.

4. **Every degradation is silent and exit 0.** Offline, timeout, rate limit, and a response the command does not recognize all mean "could not tell", and "could not tell" is not a failure of the repository. A check that complains when the network is down is a check that gets removed; and once M-047 puts it in a session hook, a noisy failure mode becomes a broken session.

5. **The answer is cached outside the repository, with a bounded lifetime.** In the user's cache directory, never in the repository: the repository is reviewed content and a cache is machine state. A repeated call inside the lifetime costs no request, which is what makes it safe to run on every session start. An unreadable, corrupt, or unwritable cache is treated as absence — it falls back to asking, or to silence — never as an error.

6. **The drift message names the way out and the way to stay.** It gains the command that resolves it and the pinning route for a repository that is not ready to move. The message's content becomes a criterion of its own; ADR-011's drift rule and AC-033 keep their meaning untouched.

7. **The platform branch and the network are both injected.** All three installation recipes and every degradation are proven without a test reaching a live service — no test may depend on GitHub being up, and no test may depend on which machine ran it.

8. **The contract is stated where adopters read it**: `guide/operations.md`'s upgrade section, `guide/getting-started.md`, the command's own usage text, CAP-004, and `CHANGELOG.md`'s `[Unreleased]`.

## What does not change

**`clue validate` stays offline.** No milestone in this campaign puts a network call where a validation verdict can read it, and this change adds none. The judge's inputs remain files in the repository.

**`clue` does not replace its own binary.** The install scripts verify each download against the release's published checksums; a self-updating binary would have to reimplement that verification inside itself. `clue latest` prints the recipe and stops. The convenience is delivered in M-046 by an agent following a skill, under the human merge boundary.

**No routing reaches an already-onboarded repository yet.** M-045 builds the thing to point at; M-046 does the pointing. This change deliberately does not touch the skills.

## Reversal cost

High, which is why it is an ADR. A new public command is a name adopters put in scripts and hooks, and the network boundary it draws is the precedent the two remaining milestones of this arc build on. Reversing it after M-046 and M-047 have shipped means retracting a command that skills and vendor configuration invoke.
