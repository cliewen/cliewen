---
id: PDR-023
type: decision
status: inferred
links: [P-010, CAP-001, CAP-004, ADR-018, ADR-042, ADR-043, PDR-019, PDR-022, C-013]
title: The routing hub is the session-start channel, and no vendor configuration is ever emitted
author: agent
accepted-by: []
---

# PDR-023 — The routing hub is the session-start channel, and no vendor configuration is ever emitted

## Context and problem statement

Cliewen can tell a repository it is behind ([ADR-042](ADR-042-release-check-outside-the-judge.md)) and can carry out the upgrade once someone asks for it ([ADR-043](ADR-043-upgrade-skill-is-a-managed-carrier.md)). Neither happens by itself. `clue latest` only helps someone who already knows the command exists, and `clue-upgrade` only runs when an agent is asked to upgrade — the request nobody makes, because nothing has told them there is anything to upgrade to. The quiet mode ADR-042 specified was built for a caller that runs without being asked.

The obvious caller is a session-start hook, and every assistant that has one configures it differently, in its own file, with its own schema. There is no cross-agent standard for executing anything when a session begins. So the obvious caller is also a commitment to one vendor's runtime behaviour, shipped into every adopted repository by a scaffold.

[PDR-022](PDR-022-vendor-entry-points-only-point.md) permits a vendor entry point that only *points* at the hub. That precedent does not carry here, and reading it as though it did is the error this decision exists to prevent. A pointer is inert: it is the minimum needed to make a cross-agent hub reachable from an assistant that would otherwise never load it, and it takes no position on what anything does. Configuration that runs a command is a position — on what an adopter's tool does, on their machine, for one vendor — taken by a project whose entire premise is that the method belongs to no assistant.

What, then, is the channel that reaches an agent when its session starts?

## Decision outcome

**The routing hub is the session-start channel, and it is the only one Cliewen uses. No executable vendor configuration is ever emitted, for any assistant, in any repository.**

*The hub is what an agent reads when it starts.* `AGENTS.md` is the cross-agent standard, it is read at the beginning of a session and again when a context window is rebuilt, and it is already the file that tells an agent how to work here. A line in it asking the agent to run `clue latest --quiet` before starting reaches exactly the moment the discovery is worth something, on every assistant, through a channel that already exists. Nothing is added to the repository, no schema is depended on, and no vendor is named.

*It is one line, and it is silent when there is nothing to say.* Where the output is a coding agent's context, every line is spent on something. The command is quiet by default in this use, cached for a day so repeating it costs no request, and silent when current, offline, or unable to tell. A check that greets a current repository every morning teaches its reader to skip it, and the same is true of the instruction to run it.

*The hub is the adopter's file, so migration reports and never rewrites.* `clue init` emits the line in the hub it materializes, and `clue migrate` reports a hub that never names the check without touching a byte of it. This is the pattern PDR-022 already chose for the entry point, and it is a better fit here than there: the file is vendor-neutral, so a team that reads the notice is never being told to accommodate a tool they do not use.

*Prose is the price, and it is the price Cliewen already pays everywhere.* A line in the hub can be skipped in a way a hook cannot. That is a real loss of determinism, and it buys agent-agnosticism, which is not negotiable. The hub already carries the rules that agents never merge their own pull requests and never weaken a check to make a build pass; those bind by being read and followed, and this one binds exactly as well as they do.

*No vendor configuration, and no exception for this repository either.* Cliewen's own repository is the reference implementation, so a vendor hook committed here would read as endorsement whatever the intent, and the first question an adopter asks about a file in this tree is whether they should have one too. The prohibition covers the product and the project alike.

### Rejected: emit a session-start hook for the vendor that has one

The mechanism that works, at a cost that is not payable. It ships executable configuration for one assistant into every adopted repository, takes a position on an adopter's tooling that no methodology should take, and creates a migration notice that nags teams about a vendor they may never have chosen. Determinism is the whole attraction — the check either ran or it did not — and it is not worth the premise.

### Rejected: read the cached answer inside `clue validate`

Circular. The cache is written by `clue latest` and by nothing else, so a repository where nobody runs the check has an empty cache and validate has nothing to report. Making validate populate it puts the network inside the deterministic judge, which is the boundary ADR-042 exists to hold.

### Rejected: a step in the emitted CI workflow

Agent-neutral, and it reaches every adopter. It also runs on ephemeral runners, where the cache never survives, so every push of every adopter would reach the release list from a shared address — spending an unauthenticated rate limit on every run to tell a human something the agent needed to know.

### Rejected: put it in the change loop's branch step instead

The skills are managed carriers that reach existing adopters mechanically, which is a genuine advantage over a notice. But the start of a change is a different moment from the start of a session: an agent asked a question, sent to read a file, or reviewing a pull request never reaches that step, and those are sessions too. The milestone is about starting up, and this would quietly answer a different question.

**Carrier:** CAP-001's criteria and design carry what `clue init` emits and what `clue migrate` reports; the scaffolded `AGENTS.md` template and this repository's own hub carry the instruction; CAP-004 owns the quiet check the hub names; PDR-022 governs the inert vendor pointer beside it.
