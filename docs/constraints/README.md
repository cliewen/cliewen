# Constraints

C-xxx: rules you **must not break** — laws, licenses, policies, and since [ADR-017](../decisions/ADR-017-conventions-are-constraints.md) the **convention register**: every methodology rule that would otherwise live only in prose. Checked against every proposal; distinct from requirements.

Each constraint carries `source:` (the doc, law, or catalog that states the rule; external catalogs plug in here later) and `enforcement: machine|partial|agent|human`. The classes say who holds the rule, and [ADR-045](../decisions/ADR-045-register-names-the-machine.md) gives two of them required declarations:

- **`machine`** — a machine holds the whole rule. The constraint names it under **Checked by** by convention here; there is no residual to state, so the lint asks for neither declaration.
- **`partial`** — a machine holds a stated subset, named under **Checked by**, and a stated **Residual** does not leave judgment. Most real rules are this shape: a `[-]` task must carry a reason, and the reason must be a real one.
- **`agent`** — the rule is **backlog for promotion**: no machine holds it yet, and the constraint states its **promotion trigger**, the condition under which one will. `clue validate` reports the count of these on its OK line, and says nothing when the count is zero.
- **`human`** — no machine can hold it, and the constraint declares its **Residual**: what stays with judgment and what it costs when that judgment fails.

`clue validate` lints the fields, the vocabulary, and the declarations. A rule leaves `agent` by gaining a real check or by being declared — never by relabelling, and never by weakening a check to make the count fall ([C-004](C-004-never-weaken-checks.md)).

Where the machines are: `clue validate` holds the rules about the corpus as it stands, and only those — it reads a repository state and never a transition ([ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md)). A rule about what a *change* did is held by a machine that is allowed to have a base: this repository's CI, the release gates, or the forge's branch protection.

This index is the register table. Its badge is the enforcement class rather than the artifact status, because the class is what the register exists to publish; the rows are curated and index regeneration preserves them.

<!-- clue:index:start -->
- [C-001 — Markdown prose is never hard-wrapped](C-001-no-hard-wrapped-markdown.md) · `machine` — Each prose paragraph and list item occupies one source line; structural boundaries alone create line breaks.
- [C-002 — Every release-relevant user-visible change adds a changelog entry](C-002-changelog-per-user-visible-change.md) · `partial` — Changes to shipped behavior, commands, workflows, generated skills, or adopter materialized files describe their user impact under `[Unreleased]`.
- [C-003 — A task marked infeasible carries its reason on the same line](C-003-skipped-tasks-carry-reasons.md) · `partial` — the tick-immediately half was withdrawn on 2026-08-08; a `[-]` with no reason is indistinguishable from a task nobody finished
- [C-004 — Never weaken a test or a lint rule](C-004-never-weaken-checks.md) · `human` — Repair a failing check at its cause or surface the conflict; never loosen the check merely to obtain green output.
- [C-005 — Every Cliewen proposal declares its plan item or plan-less](C-005-proposal-declares-plan-item.md) · `partial` — Proposal frontmatter links the real plan it serves, or the proposal body explicitly explains why no plan exists.
- [C-006 — Decision records are timeless prose; a method contract moves every live carrier together](C-006-adrs-timeless-with-carrier.md) · `human` — Decisions retain enduring context rather than episode narrative, and methodology changes update every current carrier in the same change.
- [C-007 — Diagrams use the clearest renderable form](C-007-diagrams-inline-mermaid.md) · `human` — Choose Mermaid, ASCII, or SVG by clarity while preserving every source link and referenced asset.
- [C-008 — Completed plans are immutable](C-008-completed-plans-immutable.md) · `partial` — Once a plan is completed, later changes may neither revise its content nor use it as a mutable backlog.
- [C-009 — Type-specific frontmatter fields are present](C-009-type-specific-frontmatter.md) · `machine` — Each artifact declares the additional metadata its type requires, beyond the common identity, status, links, and title fields.
- [C-010 — Milestone status values follow one vocabulary](C-010-milestone-status-vocabulary.md) · `partial` — Every recognized milestone table uses only `todo`, `doing`, `done`, or `dropped`, with skipped work explained.
- [C-011 — Future-shaping decisions route by subject into ADR, PDR, or IDR](C-011-decision-records-typed.md) · `human` — Only future-shaping choices earn records; software or corpus architecture routes to ADR, project or methodology work routes to PDR, and implementation routes to IDR.
- [C-012 — Full changes remain human-accepted while simple integration follows explicit user authority](C-012-agents-never-merge-own-changes.md) · `partial` — An agent never merges its own full change, while simple work integrates only through authority the user explicitly grants.
- [C-013 — Changes to a core carrier require an explicit decision record and human acceptance](C-013-core-changes-need-decision.md) · `human` — Altering the verifiable thread, full-loop acceptance boundary, or validator's meaning crosses the red line and must be decided openly.
- [C-014 — Total Go statement coverage stays at or above 80%](C-014-coverage-floor.md) · `machine` — Repository verification measures aggregate Go statement coverage and rejects a branch below the stated floor.
- [C-015 — A new user reaches their first green validate in under 30 minutes](C-015-onboarding-under-30-minutes.md) · `human` — A new user on a clean machine reaches a green validate within the stated bar reading only the quickstart, and a reviewer judges every quickstart or `clue init` change against that bar because no test can time the journey.
- [C-016 — A generated index row states the record it links and says what it is about, never just its filename](C-016-index-rows-state-their-record.md) · `machine` — Index generation emits the stated row, seeds its description, and the judge counts the rows still restating only their link or saying nothing about the artifact.
- [C-017 — The agentic review loop owns severity and stops within its bounded ordinary budget](C-017-agentic-review-loop-is-bounded.md) · `human` — Given a Cliewen change entering automatic agentic review, when a caller briefs the reviewer and review passes run, then the loop's own blocking/advisory model governs the verdict; computed counts and…
- [C-018 — A durable record never states a figure a command computes](C-018-no-computed-figures-in-prose.md) · `agent` — A durable record never states a figure a command computes — an artifact count, a coverage percentage, a reported population size.
- [C-019 — A reference names what it is about](C-019-references-name-what-they-are-about.md) · `human` — A report, status, or explanation written for a human or an adopter never leaves a bare `ADR-018` or `CH-008` standing alone: it says in a few words what the record holds, and a status is short and…
- [C-020 — An agent orients on the plan after a human reports a merge](C-020-orient-after-merge.md) · `human` — After a human reports that a Cliewen change's pull request merged, the agent orients before starting anything else: it describes the plan's next unfinished step in plain language and asks whether to…
- [C-021 — A suggestion raised mid-change is triaged, never held in memory](C-021-mid-change-suggestions-are-triaged.md) · `human` — A suggestion raised during a change is triaged immediately into one of two carriers.
- [C-022 — Every durable artifact serves one primary consumer](C-022-one-primary-consumer-per-file.md) · `partial` — `clue validate` reports the structural signs that one durable file or default context slice has grown past a focused read; a reader decides whether to split or accept it.
<!-- clue:index:end -->
