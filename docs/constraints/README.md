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
- [C-001 — Markdown prose is never hard-wrapped](C-001-no-hard-wrapped-markdown.md) · `machine`
- [C-002 — Every release-relevant user-visible change adds a changelog entry](C-002-changelog-per-user-visible-change.md) · `partial`
- [C-003 — Tasks tick immediately; a skipped task carries its reason](C-003-tasks-tick-immediately.md) · `partial`
- [C-004 — Never weaken a test or a lint rule](C-004-never-weaken-checks.md) · `human`
- [C-005 — Every Cliewen proposal declares its plan item or plan-less](C-005-proposal-declares-plan-item.md) · `partial`
- [C-006 — Decision records are timeless prose; a method contract moves every live carrier together](C-006-adrs-timeless-with-carrier.md) · `human`
- [C-007 — Diagrams are inline Mermaid](C-007-diagrams-inline-mermaid.md) · `machine`
- [C-008 — Completed plans are immutable](C-008-completed-plans-immutable.md) · `partial`
- [C-009 — Type-specific frontmatter fields are present](C-009-type-specific-frontmatter.md) · `machine`
- [C-010 — Milestone status values follow one vocabulary](C-010-milestone-status-vocabulary.md) · `machine`
- [C-011 — Decision records are routed by type: ADR, PDR, or log row](C-011-decision-records-typed.md) · `human`
- [C-012 — Changes are reviewed locally, root at main, and remain human-merged](C-012-agents-never-merge-own-changes.md) · `partial`
- [C-013 — Changes to a core carrier require an explicit decision record and human acceptance](C-013-core-changes-need-decision.md) · `human`
- [C-014 — Total Go statement coverage stays at or above 80%](C-014-coverage-floor.md) · `machine`
- [C-015 — A new user reaches their first green validate in under 30 minutes](C-015-onboarding-under-30-minutes.md) · `human`
- [C-016 — A generated index row states the record it links, never just its filename](C-016-index-rows-state-their-record.md) · `machine`
<!-- clue:index:end -->
