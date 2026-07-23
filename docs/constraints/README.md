# Constraints

C-xxx: rules you **must not break** — laws, licenses, policies, and since [ADR-017](../decisions/ADR-017-conventions-are-constraints.md) the **convention register**: every methodology rule that would otherwise live only in prose. Checked against every proposal; distinct from requirements.

Each constraint carries `source:` (the doc, law, or catalog that states the rule; external catalogs plug in here later) and `enforcement: machine|agent|human`. `machine` means `clue` enforces it; `agent` means the skills and agents hold it and the rule is **backlog for promotion** — each agent-enforced constraint states its promotion trigger, the condition under which it becomes a machine rule. `clue validate` lints these fields and reports the count of agent-enforced constraints on its OK line: the visible drift-risk backlog.

This index is the register table:

<!-- clue:index:start -->
- [C-001 — Markdown prose is never hard-wrapped](C-001-no-hard-wrapped-markdown.md) · `agent`
- [C-002 — Every release-relevant user-visible change adds a changelog entry](C-002-changelog-per-user-visible-change.md) · `agent`
- [C-003 — Tasks tick immediately; a skipped task carries its reason](C-003-tasks-tick-immediately.md) · `agent`
- [C-004 — Never weaken a test or a lint rule](C-004-never-weaken-checks.md) · `agent`
- [C-005 — Every Cliewen proposal declares its plan item or plan-less](C-005-proposal-declares-plan-item.md) · `agent`
- [C-006 — ADRs are timeless prose; method decisions name their carrier](C-006-adrs-timeless-with-carrier.md) · `agent`
- [C-007 — Diagrams are inline Mermaid](C-007-diagrams-inline-mermaid.md) · `agent`
- [C-008 — Completed plans are immutable](C-008-completed-plans-immutable.md) · `agent`
- [C-009 — Type-specific frontmatter fields are present](C-009-type-specific-frontmatter.md) · `agent`
- [C-010 — Milestone status values follow one vocabulary](C-010-milestone-status-vocabulary.md) · `agent`
- [C-011 — Decision records are routed by type: ADR, PDR, or log row](C-011-decision-records-typed.md) · `agent`
- [C-012 — Changes are reviewed locally, root at main, and remain human-merged](C-012-agents-never-merge-own-changes.md) · `agent`
- [C-013-core-changes-need-decision](C-013-core-changes-need-decision.md)
- [C-014-coverage-floor](C-014-coverage-floor.md)
- [C-015-onboarding-under-30-minutes](C-015-onboarding-under-30-minutes.md)
<!-- clue:index:end -->
