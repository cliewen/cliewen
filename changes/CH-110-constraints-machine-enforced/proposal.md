---
id: CH-110-proposal
type: proposal
status: active
links: [P-010, M-049, C-001, C-002, C-003, C-004, C-005, C-006, C-007, C-008, C-009, C-010, C-011, C-012, C-013, CAP-002, ADR-017, ADR-042]
title: Every constraint names the machine that holds it or declares what judgment costs
---

# CH-110 — Every constraint names the machine that holds it, or declares what judgment costs

Serves [P-010](../../docs/plans/P-010-adopters-keep-current.md) milestone **M-049**.

## What

Thirteen of Cliewen's sixteen constraints — C-001 through C-013 — carry `enforcement: agent`: nothing but an agent remembering to read them holds them. `clue validate` has counted them on its OK line since [ADR-017](../../docs/decisions/ADR-017-conventions-are-constraints.md) as the visible promotion backlog. This change empties that backlog, and the count reaches zero.

Six of the thirteen get a real check inside the judge. One gets a real check in CI. Two already had a machine nobody had written down. The remaining four are declared: what part can never leave judgment, and what that costs.

## Why the backlog cannot be emptied by building what the triggers ask for

Every one of the five remaining triggers asks for the same thing: "`clue` gains git-diff context." C-002 wants to know whether *this change* touched user-visible behavior; C-004 whether *this change* loosened a test; C-008 whether *this change* edited a completed plan; C-013 whether *this change* touched a core carrier without a decision. They ask what a change did, not what the repository is.

`clue validate` is the deterministic judge ([ARCH-003](../../docs/architecture/core.md)): the same corpus yields the same verdict for everyone, everywhere. A verdict that reads history is not that. It would depend on which branch the caller stands on, how deep the clone is, whether `main` was fetched, and whether the base is even present — a shallow CI checkout and a developer's worktree would disagree about a corpus whose bytes are identical. This is the same boundary [ADR-042](../../docs/decisions/ADR-042-release-check-outside-the-judge.md) drew around the network, for the same reason: the judge answers from what it can see in front of it.

So the triggers were a design error, not a backlog. The transition rules do have machines — just not that one. Branch protection on `main` holds C-012's detectable subset today, and has since before the constraint was written. The release workflow holds C-002 at release time. A base-state comparison is exactly what a CI job may do, and C-008 gets one there.

## The shape

**Two decisions.** An ADR states that `clue validate` judges a repository state and never a transition, naming where transition rules live instead and retiring the git-diff promotion triggers. A second ADR amends [ADR-017](../../docs/decisions/ADR-017-conventions-are-constraints.md): the register stops carrying a promotion backlog as its organising idea and starts naming, for every constraint, the machine that holds it and the judgment that remains. A decision-log row declares the milestone status vocabulary that C-010 needs before it can be linted.

**A fourth enforcement class.** `machine | agent | human` cannot say "a machine holds this part and a person holds that part", and most real rules are that shape. `partial` is added: a named machine holds a stated subset, and the constraint declares both **Checked by** and **Residual**. `agent` keeps its exact current meaning — awaiting a machine check, the thing the OK line counts — so an adopter's register and the reported count still mean what they meant. `human` widens from "a person verifies" to "held by judgment, human or agent, and no machine can take it" — which is what C-004, C-006, C-011, and C-013 actually are.

**Six checks in the judge**, each for a rule about state:

| | Check |
|---|---|
| C-001 | prose-layout lint: a paragraph or list item broken across lines in `docs/**` and `changes/**` markdown, outside fences, tables, and frontmatter |
| C-003 | a `[-]` task in `changes/**/tasks.md` with no prose after the checkbox |
| C-005 | `changes/*/proposal.md` carries a P/M link or the literal `plan-less` declaration |
| C-007 | image links and image files under `docs/**` |
| C-009 | per-type required frontmatter: `author` and `accepted-by` on decisions, `goal` on capabilities |
| C-010 | milestone status cells in plan tables against the declared vocabulary |

**One check in CI.** C-008 gets a workflow step that fails a pull request modifying a plan file whose `main`-side status is `completed`. It reads the base state, which is the workflow's business and not the judge's; the change that closes a plan is untouched, because the file it edits is `active` on `main`.

**Two rules already machine-held, now written down.** C-012's protection is live on this repository: pull requests required, administrators included, `validate` a required status check, conversation resolution required, force pushes and deletions blocked. C-002's release-time half is the release workflow, which already fails a release with no version section. Neither constraint said so.

**Four declared.** C-004 (weakening a check versus refactoring it is meaning), C-006 (timelessness, and deriving an arbitrary contract's carrier set), C-011 (routing a decision by what it is *about*), C-013 (whether a change alters core *meaning*) each state their residual and its cost in their own carrier.

## Result

`clue validate` reports no agent-enforced constraints. The scaffolded corpus a new adopter receives reports none either, because its only shipped constraint is C-001 and C-001 is now checked.

## Scope

Full tier, and it crosses the red line: adding checks changes what a green `clue validate` asserts, so [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) requires the explicit decision records above and human acceptance at merge. The enforcement vocabulary is a methodology contract under [C-006](../../docs/constraints/C-006-adrs-timeless-with-carrier.md), so every live carrier moves in this change: the register and its scaffolded template, ADR-017, the architecture door note, CAP-002's criteria and design, the validator and its CLI text, the scaffolded `AGENTS.md` note, and the guide pages that print the OK line.

## Not in this change

No check is weakened, and no constraint is reclassified to make the number fall — a rule moves out of `agent` only by gaining a real check or by declaring, truthfully, what cannot be mechanized. `clue` gains no git access. C-014, C-015, and C-016 keep their current classes.
