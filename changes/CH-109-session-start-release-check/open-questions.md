---
id: CH-109-open-questions
type: open-questions
status: open
links: [CH-109, P-010, PDR-022, ADR-018]
title: Open questions — CH-109
---

# Open questions — CH-109

## Q1 — M-047's exit criterion requires emitting vendor configuration, which violates agent-agnosticism

**Raised:** 2026-08-03, during implementation, by the human.

M-047 states that "the scaffold emits that configuration for a new session and for a restarted context window alike". The only mechanism that satisfies that literally is a vendor's own settings file — for Claude Code, `.claude/settings.json` — because no cross-agent standard exists for running anything at session start. The implementation was written and then rejected: Cliewen must be agent-agnostic, and emitting executable configuration for one assistant crosses that line.

The distinction matters and is not the one [PDR-022](../../docs/decisions/PDR-022-vendor-entry-points-only-point.md) already settled. That decision permits a vendor entry point that only *points* at the hub — prose, inert, and the minimum needed to make a cross-agent hub reachable at all. Configuration that *executes* is a different thing: it takes a position on what an adopter's tool does, for one vendor, in a project whose whole premise is that the method belongs to no assistant.

The milestone cannot be delivered as written. The question is what replaces it.

**Recommendation: make the carrier a skill, and document the vendor wiring instead of shipping it.**

Two parts, and the first is the one that matters.

Cliewen already has an agent-agnostic carrier that reaches existing adopters: the generated skills, refreshed by `clue migrate`. `clue-delta` opens with **Branch**, and that is the honest moment to ask whether there is anything to upgrade to — better than session start, in fact, because a session start catches an agent that may be about to do something entirely unrelated, while the start of a change is when the answer changes what happens next. One line in the skill — run `clue latest --quiet` before branching, and route a non-empty answer to `clue-upgrade` — reaches every adopter on every assistant, through a carrier this project already owns and already upgrades. It is less deterministic than a hook, because an agent can skip a line of a skill in a way it cannot skip a hook; that is the price of agnosticism, and it is the same price every other Cliewen rule already pays.

The second part covers the literal wish without shipping it: `guide/operations.md` gains a short, vendor-neutral section saying that `clue latest --quiet` is safe to run at session start — it is silent when current, silent when offline, bounded, and exits 0 — and showing how an adopter wires it into whichever assistant they use, as their configuration in their file. Cliewen describes; the adopter decides. That also settles the migration question by dissolving it: with nothing emitted, there is nothing whose absence to report, so no MIG-006 and no notice about a vendor an adopter may not use.

Both parts change M-047's exit criterion, which P-010's mutation rules put behind a declared plan revision backed by a decision record.

**Answered:** 2026-08-03, by the human, who accepted the prohibition and asked whether it could be handled better than by documenting a recipe. It could. The recommendation above was superseded before implementation: the skills reach existing adopters mechanically, but the start of a change is not the start of a session, and the second part left the actual wish unimplemented. The routing hub is the session-start channel — the cross-agent file every assistant reads when a session begins — and the objection that it is adopter-owned is answered by the entry point's existing pattern of emitting into the template and reporting for everyone else. Recorded as [PDR-023](../../docs/decisions/PDR-023-tool-notice-and-hub-instruction.md), which carries the M-047 revision.

## Q2 — the hub's instruction reached a real session and did not convert into an action

**Raised:** 2026-08-03, during AC-085's Human verification, by the human.

AC-085 was run and it failed. The instruction was not missed on the way in: `AGENTS.md` was loaded into the session's context through `CLAUDE.md`, and the release-check paragraph was its first content. The agent read it and went straight to the requested `/review` task. So the channel delivered and the conversion from reading to acting is what did not happen.

Two contributing defects were visible in the wording. The trigger was vague — "when you start" was read as "start the task" by a session that arrived carrying an explicit slash-command task list — and the imperative was buried under about a hundred words of rationale in a single paragraph. Neither fully explains the failure. The instruction asks for an unprompted action *before* work begins, which is a weaker thing to ask than the rules around it: those constrain work already underway and are consulted at the moment they bind.

A second finding came from the same session. The installed `clue` predated `clue latest`, so the instruction, had it been obeyed, would have produced `unknown command "latest"` and a full usage dump — noisy, exit 2, and the opposite of the one line or silence [ADR-042](../../docs/decisions/ADR-042-release-check-outside-the-judge.md) promises. This is not a local quirk: the repositories most likely to be behind are exactly the ones whose binary may be too old to carry the command that reports it.

**Recommendation: make the tool volunteer the answer, and keep the prose as the fallback.**

The mechanism should not depend on an agent remembering to ask. The ordinary `clue` workflow commands can print the notice themselves, which is the update-notifier pattern `gh`, `npm`, and `brew` already established, and it reaches an agent through work it was already doing — `clue context` is named in the hub's own routing and in every lifecycle skill. The gate is where the care goes: never `clue validate`, so the deterministic judge's verdict and output stay a statement about the repository alone; never `clue version`, which answers instantly and offline forever; never off a terminal, in CI, or against a documented opt-out, so a script's captured output is byte-identical and an ephemeral runner never reaches the release list. This survives both rejections already recorded here — reading the cache inside validate, and a CI step — rather than reopening either.

The hub's line stays, because a session that runs no `clue` command at all is still a session and the notifier never reaches it. It gets a sharper trigger, its imperative separated from its rationale, and a sentence reading an unknown-command error as the answer it is. What changes is its job: fallback rather than mechanism.

**Answered:** 2026-08-03, by the human, who asked whether a checkmark, an assertion of criticality, or a common hook mechanism would fix it, and accepted this instead. A cross-agent hook does not exist — the `AGENTS.md` specification is standard Markdown with no lifecycle hooks or execution semantics, and every hook that exists belongs to one vendor, which PDR-023 already put out of bounds. A self-reported checkmark travels on the channel that just failed and would greet a current repository every morning, which ADR-042 refused. Marking the instruction critical addresses importance, and importance was not what failed. Recorded by revising [PDR-023](../../docs/decisions/PDR-023-tool-notice-and-hub-instruction.md) in place — it was unmerged, `inferred`, and unaccepted, so it is corrected rather than superseded — which carries the second M-047 revision and the restatement of AC-085.

## Q3 — the notifier's terminal gate excludes every coding-agent session, which is the audience it was built for

**Raised:** 2026-08-03, during review of this branch's implementation, by the agent.

[PDR-023](../../docs/decisions/PDR-023-tool-notice-and-hub-instruction.md) gates the notice on standard error being a terminal. Coding-agent harnesses capture a command's output through a pipe rather than a pseudo-terminal, so standard error is not a character device and `notifierAllowed` returns false. Measured in a Claude Code session on 2026-08-03: both streams report `charDevice=false` with `CI` unset. Every other assistant that runs shell commands and reads their output does the same thing, because reading the output is the point.

So the notifier never fires for an agent. The mechanism and the fallback have swapped jobs: the hub's line works in an agent session — `clue latest --quiet` writes to standard output and no gate touches it — and the notice that Q2 made the mechanism, on the finding that prose could not be relied on, reaches only a human at an interactive prompt. AC-085's retest is set up to fail exactly as it failed before, and for a reason the change already knows about.

The gate is not incidental. It carries the ephemeral-runner rejection and the byte-identical-capture promise, and it is stated in PDR-023, CAP-004's criteria and design, `guide/operations.md`, `guide/getting-started.md`, the `[Unreleased]` entry, and one case of `TestAC087_UnitNegative_TheGateExcludesTheJudgeAScriptAndARunner`.

**Recommendation: drop the terminal condition and keep the other three.**

The terminal check is doing two jobs and only one of them needs it. Keeping an ephemeral runner off the release list is the `CI` check's job — every runner worth naming sets it, and the code already calls the terminal check "belt-and-braces beside" it rather than the mechanism; Cliewen's own emitted wall runs `clue validate`, which carries no notice under any gate. What the terminal check uniquely buys is that a piped stream stays byte-identical, and that is the promise to narrow: a script consumes standard output and an exit code, and both stay identical with the notice and without it, which is what AC-086 already says in its own third clause. Standard error is where a tool speaks to a person, and in an agent session the person is the agent.

That leaves three gates — the command, `CI`, and `CLUE_NO_UPDATE_NOTIFIER` — and every rejection PDR-023 recorded survives untouched. Anyone who wants the old behaviour has the documented switch, which is the same answer the decision already gives to anyone who wants silence.

The alternative — detecting an agent harness by its own environment variable — is worse on this project's own terms: it is a per-vendor list that ages, and reading one assistant's variable to decide what the tool does is the position PDR-023 refuses to take.

Revising PDR-023 in place is still available: it is unmerged, `inferred`, and unaccepted, exactly as it was for Q2.

**Answered:** 2026-08-03, by the human, who accepted the recommendation as given. The terminal condition is removed and the gate keeps its other three checks; PDR-023 is revised in place again — still unmerged, `inferred`, and unaccepted — with the environment-sniffing alternative recorded as rejected for the reason above. AC-087 gains the direction that failed, driving the notice through a real pipe rather than a buffer, and AC-086's byte-identical claim is now proven against the built binary rather than left as prose.

**Status:** resolved. No open questions remain.
