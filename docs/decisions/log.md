---
id: LOG-001
type: log
status: active
links: [PDR-003]
title: Decision log
---

# Decision log

Decisions that are cheap and local to reverse, one row each — newest first. **The cost of undoing a decision is the sole routing test** ([PDR-006](PDR-006-decision-records-are-typed.md)): cheap and local to reverse → a row here; expensive to reverse → a full record — ADR for architecture, PDR for project/process. A rule's reach does not change its routing: a decision can apply to every future change and still be a log row, as long as undoing it stays cheap and local (this table's coverage-gate row is the precedent). Demotion/promotion mechanics are [PDR-003](PDR-003-decision-log.md). Rows are never deleted; a reversed decision gets a new row, and a row whose reversal turns out to be expensive after all is promoted to a full record citing this table. Full text of demoted records lives in git history.

| Date | Decision | Why | Change/PR |
|---|---|---|---|
| 2026-07-17 | The scaffolded constraints register is seeded with exactly the conventions the generated AGENTS.md declares that no versioned skill carries — today C-001 (markdown never hard-wrapped); skill-carried rules are not duplicated into it | Skills are versioned carriers held by the drift lint, so registering copies would create a second carrier that drifts; a rule living only in the scaffolded AGENTS.md would be prose-only from day one, which is what the register exists to prevent. A template edit reverses this cheaply | CH-020 / PR #19 (review finding) |
| 2026-07-17 | AC-001 (30-minute onboarding) retired and split: the mechanical path became AC-002/AC-024/AC-025, the 30-minute human journey became QS-002 | The clock spans reading and installing — a human journey no focused test pair can verify; the granularity rule says split before testing, and re-minting ACs is cheap and local | CH-020 |
| 2026-07-17 | The init template tree lives at `internal/scaffold/templates/`, not the originally sketched `/.cliewen/templates` | The Go toolchain ignores dot-directories, putting them out of `go:embed`'s reach; moving a template folder later is a cheap local change (the embed decision itself is ADR-018) | CH-020 |
| 2026-07-17 | After a reported merge the agent orients before proceeding: it describes the plan's next step in plain language and asks the human to start it, or asks what to do next when the plan is empty | Agents treated "PR merged" as a go signal and silently started the next task, leaving the human without a picture of where the plan stands; carrier is `clue-delta` step 5 | CH-018 / PR #17 |
| 2026-07-16 | Litmus wording disambiguated: reversal cost is the sole routing test; "constrains future changes" means expensive to undo, not wide in reach | Two agent reviewers read the two phrasings as different tests and reached opposite routings for the same decision; the human directed that the reversal-cost reading takes precedence and the ambiguity be removed | CH-017 / PR #16 (review dispute) |
| 2026-07-16 | The digest is never a task in `tasks.md`; the digest precondition applies to the work only | A self-referential "digest" task can only be ticked falsely (before the deletion) or violate the precondition (unticked at deletion) — both misreadings have occurred; carrier is `clue-delta` step 4 | CH-017 / PR #16 (review finding) |
| 2026-07-12 | Default test-coverage gate: 80% total statement coverage, enforced in CI | A tripwire for code paths between the AC-covered ones — high enough to catch untested code landing, low enough not to invite meaning-free tests; total (not per-package) so thin entry points don't force contortions | demoted from ADR-004 by CH-010; carrier stays the CI workflow template |
| 2026-07-12 | Parse frontmatter with gopkg.in/yaml.v3 | A hand-rolled YAML subset would reject legitimate frontmatter with confusing errors, and the judge must never be wrong about form; one small pure-Go dependency doesn't threaten single-binary distribution | demoted from ADR-003 by CH-010 |
