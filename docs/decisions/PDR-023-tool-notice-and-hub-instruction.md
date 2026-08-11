---
id: PDR-023
type: decision
status: verified
links: [P-010, CAP-001, CAP-004, ADR-011, ADR-018, ADR-042, ADR-043, PDR-019, PDR-022, C-013]
title: The tool carries the notice and the hub carries the instruction, and no vendor configuration is ever emitted
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-023 — The tool carries the notice and the hub carries the instruction, and no vendor configuration is ever emitted

## Context and problem statement

Cliewen can tell a repository it is behind ([ADR-042](ADR-042-release-check-outside-the-judge.md)) and can carry out the upgrade once someone asks for it ([ADR-043](ADR-043-upgrade-skill-is-a-managed-carrier.md)). Neither happens by itself. `clue latest` only helps someone who already knows the command exists, and `clue-upgrade` only runs when an agent is asked to upgrade — the request nobody makes, because nothing has told them there is anything to upgrade to. The quiet mode ADR-042 specified was built for a caller that runs without being asked.

The obvious caller is a session-start hook, and every assistant that has one configures it differently, in its own file, with its own schema. There is no cross-agent standard for executing anything when a session begins: the `AGENTS.md` specification is standard Markdown with no lifecycle hooks and no execution semantics, and every hook that exists belongs to one vendor. So the obvious caller is also a commitment to one vendor's runtime behaviour, shipped into every adopted repository by a scaffold.

[PDR-022](PDR-022-vendor-entry-points-only-point.md) permits a vendor entry point that only *points* at the hub. That precedent does not carry here, and reading it as though it did is the error this decision exists to prevent. A pointer is inert: it is the minimum needed to make a cross-agent hub reachable from an assistant that would otherwise never load it, and it takes no position on what anything does. Configuration that runs a command is a position — on what an adopter's tool does, on their machine, for one vendor — taken by a project whose entire premise is that the method belongs to no assistant.

This decision first answered that the routing hub was the channel, and that a line of prose in it was enough. **That answer was tested and it failed.** AC-085 put the instruction in front of a real session on 2026-08-03, and the agent read past it and began the requested work. The failure is not the one that was expected: the channel delivered, because the hub was loaded into the session's context and the instruction was in its first paragraph. What did not happen is the conversion from reading to acting. A second observation followed from the same session — the installed binary predated `clue latest`, so the instruction, had it been obeyed, would have produced an unknown-command error and a full usage dump rather than the one line or the silence the contract promises.

Both observations point the same way. An instruction that asks an agent to perform an unprompted action before it starts is weaker than the rules around it, which constrain work already underway and are consulted at the moment they bind. So the question is not which file reaches the agent, but what does not depend on the agent remembering.

## Decision outcome

**The tool carries the notice and the hub carries the instruction. `clue` itself reports a newer release as a side effect of ordinary commands, and the hub's line remains as the fallback for sessions that run none. No executable vendor configuration is ever emitted, for any assistant, in any repository.**

*The tool volunteers the answer, so nothing has to remember to ask.* Any `clue` workflow command — `init`, `scaffold`, `context`, `migrate`, `refs` — prints one line to standard error when a newer release exists. This is the update-notifier pattern that `gh`, `npm`, and `brew` already established, and it inverts the dependency this decision got wrong the first time: an agent learns it is behind as a consequence of work it was already doing, rather than by performing an action nobody prompted. `clue context` is named in the hub's own routing and in every lifecycle skill, so the notice reaches the ordinary path of ordinary work.

*The judge is excluded, and so is `version`.* `clue validate` never carries the notice. The verdict, its exit code, and its output stay exactly what they were, because a deterministic judge that prints a line depending on another system's present state is no longer reporting only on the repository ([C-013](../constraints/C-013-core-changes-need-decision.md)'s red line, held rather than approached). `clue version` is excluded for ADR-042's own reason: it is the one command guaranteed to answer instantly, offline, and identically forever. `clue latest` is excluded because it *is* the check.

*The notice is gated so it can never become noise or cost.* It appears only outside CI, only when `CLUE_NO_UPDATE_NOTIFIER` is unset, and only when the answer is known and behind. The ephemeral-runner objection that rules out a CI step rules itself out here: a runner sets `CI`, so it never reaches the network. The answer is cached for a day as ADR-042 specified, and the ambient budget is shorter than the requested check's, because a check the user did not ask for must never be something they can feel.

*The gate has no terminal condition, and that is the third thing this decision got wrong before it got right.* The first implementation required standard error to be a terminal, on the reasoning that a piped stream should stay byte-identical and a runner has no terminal anyway. It was measured during review and it excluded the entire audience: a coding agent runs a command and captures its output through a pipe, because reading the output is the point, so the notice was switched off in exactly the sessions this decision built it to reach — leaving only the prose fallback that the section above records failing. The condition was doing two jobs and needed neither. `CI` keeps runners off the release list on its own. What a terminal check uniquely bought is narrower than it looked: a script consumes standard output and an exit code, and both are byte-identical with the notice and without it, so what a script captures is unchanged in the way that matters to it. Standard error is where a tool speaks to a person, and in an agent session the person is the agent.

*Nothing about the boundary moves.* No repository file is written, the binary never replaces itself, no verdict changes, and every degradation is silence at exit 0. The notice is advisory text on standard error; the exit code and standard output of every command are byte-identical with it and without it.

*The hub keeps the instruction, as the fallback and not the mechanism.* A session that runs no `clue` command at all — a question, a review, a conversation — is still a session, and the notifier never reaches it. The hub's line covers that case, with a sharper trigger than the one that failed: before the first tool call, whatever the request is. It is prose, and prose is now carrying only what nothing else can.

*An unknown-command error is itself the answer.* A binary old enough to lack `clue latest` is a binary that is behind, and the repositories most likely to be behind are exactly the ones whose binary may not have the command that reports it. The hub says so, so an agent reads the error as the signal it is rather than as a broken instruction. This is the one gap no mechanism closes from inside: a notifier that does not exist in the installed binary cannot announce itself. It self-heals after one upgrade, and until then the error carries the meaning.

*Prose is still the price where prose is still the carrier, and it buys agent-agnosticism.* The hub already carries the rules that agents never merge their own pull requests and never weaken a check to make a build pass. Those bind by being read, and they bind well because they are consulted when the work reaches them. This decision no longer asks a line of prose to do more than that.

*No vendor configuration, and no exception for this repository either.* Cliewen's own repository is the reference implementation, so a vendor hook committed here would read as endorsement whatever the intent, and the first question an adopter asks about a file in this tree is whether they should have one too. The prohibition covers the product and the project alike, and it is unchanged by everything above.

### Rejected: emit a session-start hook for the vendor that has one

The mechanism that works, at a cost that is not payable. It ships executable configuration for one assistant into every adopted repository, takes a position on an adopter's tooling that no methodology should take, and creates a migration notice that nags teams about a vendor they may never have chosen. Determinism is the whole attraction — the check either ran or it did not — and it is not worth the premise. Documenting the wiring for an adopter who wants it in their own file is not this, and stays available.

### Rejected: the hub's prose alone

This decision's first answer, rejected by its own acceptance test rather than by argument. The instruction reached the session and did not convert into an action, and no rewording makes an unprompted pre-task action as reliable as a rule consulted at the moment it binds. The line is kept for the sessions the notifier cannot reach; it is no longer asked to be the mechanism.

### Rejected: keep the terminal gate and detect an agent harness by its environment

The narrow repair for the failure above, and it fails this decision's own test. Every assistant announces itself with its own variable, so this is a per-vendor list maintained inside `clue` — one that is wrong for every harness not on it and ages the moment a new one ships. Reading one assistant's environment to decide what the tool does is the position this decision refuses to take three paragraphs later, and it would buy nothing the plain removal does not: the notice belongs on standard error whoever is reading, and the opt-out already serves anyone who wants silence.

### Rejected: a self-reported confirmation line in every session

An agent that skips the check will skip the instruction to report the check, because both travel on the channel that already failed. What it would buy is visibility rather than compliance — the human sees the missing line — and it would cost a line in every session forever, which is precisely what ADR-042 refused: a check that greets a current repository every morning teaches its reader to ignore it.

### Rejected: read the cached answer inside `clue validate`

Circular, and now doubly excluded. The cache is written by the release check and by nothing else, so a repository where nobody runs it has an empty cache and validate has nothing to report; making validate populate it puts the network inside the deterministic judge, which is the boundary ADR-042 exists to hold. The notifier does not rescue this: it deliberately does not run under validate, so that the judge's output and verdict remain a statement about the repository alone.

### Rejected: a step in the emitted CI workflow

Agent-neutral, and it reaches every adopter. It also runs on ephemeral runners, where the cache never survives, so every push of every adopter would reach the release list from a shared address — spending an unauthenticated rate limit on every run to tell a human something the agent needed to know. The notifier's gate encodes this rejection rather than contradicting it.

### Rejected: put it in the change loop's branch step instead

The skills are managed carriers that reach existing adopters mechanically, which is a genuine advantage over a notice. But the start of a change is a different moment from the start of a session: an agent asked a question, sent to read a file, or reviewing a pull request never reaches that step, and those are sessions too. The milestone is about starting up, and this would quietly answer a different question.

**Carrier:** CAP-004 owns the release check and the notifier's gate; CAP-001's criteria and design carry what `clue init` emits and what `clue migrate` reports; the scaffolded `AGENTS.md` template and this repository's own hub carry the fallback instruction; `guide/operations.md` documents the notice and its opt-out; ADR-011's drift message and PDR-022's inert vendor pointer sit beside it unchanged.
