---
id: CH-211
type: change
status: open
links: [CAP-012, PDR-017, PDR-067, PDR-072, ADR-033, ARCH-003, C-013]
title: Acceptance needs the human's decision, not an acceptance brief
---

# CH-211 — Acceptance needs the human's decision, not an acceptance brief

This change is plan-less: the maintainer asked for it directly. It proceeds under draft, inferred VIS-001 without confirming it.

## What and why

Today a tracked change can only be accepted with an acceptance brief: a document with five required sections (intent, criteria and evidence, binding decisions, verification, review) that the agent writes after the work is done. In local mode `clue accept` refuses to run without a complete brief file, binds a recorded approval to the brief's hash and puts the whole brief in the merge message. In PR mode the validation workflow fails a ready pull request whose body lacks an `## Acceptance brief` section. The skills tell the agent to write the brief, show it, and wait for an approval given after the human has read it.

The maintainer, who is also the only user of the method so far, has decided that the brief is to go. In practice the brief repeats what the proposal, the commits, the test results and the conversation already hold, it costs a long writing step at the end of every change, and the ritual of approving only after reading it has twice stood between the maintainer and an integration they had already decided on. On 2026-11-09 the maintainer skipped it outright in an adopting repository and asked for it to be removed from Cliewen.

What stays: the human still decides acceptance, and the agent still never accepts its own change. `clue accept` stays and keeps every check that does not depend on the brief: full candidate and base IDs, a clean integration checkout, the candidate incorporating the current base, the same acceptance policy on base and candidate, the proposal introduced after the base in the candidate's history, a valid and digested corpus, safe snapshots, the exact candidate tree in a merge with ordered parents, and never pushing.

What goes: the brief file, its template, its parsing and section checks, its hash in the approval record, the `reviewed` field, the PR-body brief check, and every skill instruction to write a brief, show it and wait for approval after it.

## How it works afterwards

- **Local mode.** `clue accept <candidate> --base <base>` takes no brief. It finds the change from the one proposal introduced after the base in the candidate's history, as it already checks today, and refuses when there is none or more than one. It writes the merge message itself: `Accept CH-xxx locally`, the proposal's title, candidate and base, the human evidence described below, and the approval provenance. Interactive confirmation shows candidate, base, change and title.
- **The human's decision.** The agent records the human's actual words in the approval record, as today, but the record binds only candidate, base and decision. Nothing has to be read before deciding; the human decides when they are ready.
- **PR mode.** The pull request body is free. The validation workflow no longer requires an `## Acceptance brief` section.
- **Human proof.** A criterion with `Test-type: Human` is still proved by a person's observation, never by an invented test. Its carrier moves from the brief to a `## Human evidence` section in the change's proposal, written by the agent from what the human reported (observer, conditions, observation, result) before the digest. The proposal is already kept in history and already found by `clue accept`, which copies that section into the merge message, so the observation stays in the accepted history. In PR mode the section is in the proposal in the pull request's history.
- **Old use.** A `--brief` flag, or a brief file, is refused with a plain message that the brief is no longer used, so an adopter running an old skill against a new binary learns why instead of getting a silent difference.

## Challenge the commitment

The riskiest assumption is that the brief holds nothing a human needs at the moment of deciding. It has held one thing nothing else did: an honest list of what was not observed and what remains unproven, such as "the Python bot did not start on this computer". Without the brief that list would live only in the conversation and be lost. The change therefore keeps it in the proposal: the human evidence section also says what was not observed, and it is in the merge message. That is a small part of the brief kept where it is cheap, and the rest goes.

The credible alternative is to keep the brief but make it short and optional. It was rejected because an optional document becomes a habit agents keep producing, which is what the maintainer wants gone, and because the obligations it carried that matter (exact candidate, valid corpus, human proof, human decision) are either checked by `clue accept` or kept in the proposal.

The core is touched. The verifiable thread keeps Human proof, with a new carrier; the human acceptance boundary keeps the human's decision and the agent's exclusion, and loses only the document the decision was bound to. A decision record states this and supersedes the brief parts of PDR-017, PDR-067, PDR-072 and ADR-033, as the red line requires.

## Carriers and proof

- `internal/accept`: no brief, the change found from the proposal, the merge message written by the command, the approval record without the brief hash, `--brief` refused. Tests rewritten for each, including refusal of no proposal, two proposals, a stale or mismatched approval, and an old `--brief`.
- `internal/corpus/actests.go`: the comment and message that say a Human criterion is satisfied by the brief now name the proposal's human evidence section.
- The validation workflow and the repository's own CI: the PR-body brief check is removed.
- Skills sources (`clue-delta`, `clue-verify`, `clue-extract`, `clue-upgrade` and the shared review-boundary, change-tiers, intent-model and readable-references fragments), regenerated; the brief and approval templates in the scaffold.
- CAP-012 criteria AC-219, AC-237 and AC-238 reworded and the design updated; CAP-006 and CAP-002 where they name the brief; `docs/architecture/core.md`; a new decision record; the public guide's acceptance pages; `CHANGELOG.md` with the upgrade action for adopters.
- Proof: unit and integration tests in `internal/accept` and the skills tests on the generated text, `clue validate`, and the source repository's **Verify Locally** block.

Not in this change: the adopting repository where this was asked for follows by an ordinary upgrade after a release, and changes its own rules that repeat the brief there.

## Questions for the maintainer

1. **PR mode too?** Recommended yes, so the method has one rule and no brief anywhere; otherwise PR mode keeps its body check and its own Human-proof carrier.
2. **Human evidence in the proposal?** Recommended yes, as described above; the alternative is a separate small file per change, which is the brief again in miniature.
3. **An old `--brief` is refused rather than accepted and ignored?** Recommended refused with a plain message, so nothing differs in silence.
4. **The approval record stays, written by the agent from the human's words?** Recommended yes: it costs the human nothing and keeps the decision in history. The alternative is no record at all, only the interactive prompt.
