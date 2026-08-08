---
id: ADR-045
type: decision
status: verified
links: [P-010, CAP-002, ADR-017, ADR-044, C-004, C-011, C-013]
title: Every constraint names the machine that holds it or the judgment that remains
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-045 — Every constraint names the machine that holds it or the judgment that remains

## Context and problem statement

[ADR-017](ADR-017-conventions-are-constraints.md) built the constraint register around a single idea: `agent` is a waiting room, every rule in it states the trigger that will get it out, and the count on `clue validate`'s OK line is the queue length. That was the right frame for a register whose rules had never been examined one at a time. It stops being right once they have.

Examined, most rules are not one thing. A `[-]` task must carry a reason, and the reason must be a real one: a machine can hold the first half exactly and can never hold the second. A changelog entry must exist at release and must be written for users. Changes must be human-merged, and reach `main` only through a pull request: a forge holds the second, and nothing holds the first but the agent declining to. The register had one field and three values, so every one of these was recorded as unheld — which is both wrong and demoralising, because it hides the machines that are working.

The waiting-room frame also outlived its truth in the other direction. Four rules were queued for a check that [ADR-044](ADR-044-judge-reads-state-not-transitions.md) has now decided should never be written, and some rules — routing a decision by what it is *about*, judging whether an edit weakened a check — are permanently judgment. A queue that contains items that will never be served is not a queue; it is a way of never having to say so.

## Decision outcome

**A constraint states the machine that holds it, the part that machine holds, and the judgment that remains. `enforcement` gains a fourth class for the ordinary case where both are true.**

*`partial` is added to the vocabulary.* The full set is `machine | partial | agent | human`. `partial` means a named machine holds a stated subset and a stated residual does not leave judgment. Most real rules are this shape, and before this decision the register could not say so — a constraint had to overclaim by calling itself `machine` or undersell by calling itself `agent`.

*A `partial` or `human` constraint declares two things in its body.* **Checked by** names the machine and the exact subset it holds — `clue validate`, a named workflow, the forge's branch protection. `human` carries it only when some machine holds a fragment worth naming. **Residual** names what stays with judgment and what it costs when that judgment fails: not an apology, a statement of the exposure a reader is accepting. `clue validate` requires of each class what that class owes — both from `partial`, the residual from `human` — so a constraint cannot claim a machine without naming it or claim permanence without pricing it.

*`agent` keeps exactly the meaning it has.* It is a rule awaiting a machine check, it states a promotion trigger, and it is what the OK-line count counts. This is deliberate: the count is a contract with every adopter's register and with the analyses that cite it, and redefining a number is not the same as reducing it. A rule leaves `agent` by gaining a real check or by being declared, never by relabelling.

*`human` widens.* It meant "only a person can verify this". It now means "held by judgment, and no machine can take it" — the judge does not care whether the judgment is a person's or an agent's, and the distinction the old reading drew is not one the register can observe. What matters is that nothing mechanical will catch a violation, which is the same fact either way.

*A constraint stating a permanent property states it as one.* "No machine can hold this and here is what that costs" is a finished answer. It is not a smaller version of a promotion trigger, and it is not entered in a backlog to be looked at again next campaign.

## Rejected: a `checked-by:` frontmatter field

More lintable than prose, and it would make the register machine-summarisable. It also narrows a corpus obligation for every adopter whose constraints already carry `enforcement: machine` — their registers would fail on the next binary over a field they have never heard of, to gain a summary nobody asked for. The declaration lives in the body, where the subset can be described precisely; a subset is a sentence, not an identifier.

## Rejected: drop the OK-line count once the backlog is empty

Tempting, because the number would be zero for this repository. But the count is the mechanism by which an adopter's own un-triaged rules stay visible, and this repository is not the audience for it. A count that only ever reads zero here is the correct outcome of using it, not a reason to remove it.

## Rejected: keep three classes and let `machine` mean "mostly"

The cheapest change, and it is how the register would have drifted on its own. It makes the strongest label the vaguest one: a reader could no longer tell a rule a machine fully holds from one where a machine holds a fragment, which is exactly the distinction the register exists to publish.

## Carrier

The constraints register README and its scaffolded template state the vocabulary and the two declarations; ADR-017's vocabulary sentence and promotion-trigger rule are amended here rather than restated; `clue validate` lints the class and the declarations, and CAP-002's criteria and design carry the checks. The architecture's `enforcement:`-classes door note names the widened set.
