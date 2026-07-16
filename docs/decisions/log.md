---
id: LOG-001
type: log
status: active
links: [PDR-003]
title: Decision log
---

# Decision log

Decisions that are cheap and local to reverse, one row each — newest first. The litmus test and demotion/promotion mechanics are [PDR-003](PDR-003-decision-log.md); decisions that constrain future changes are full records in this folder — ADRs for architecture, PDRs for project/process ([PDR-006](PDR-006-decision-records-are-typed.md)). Rows are never deleted; a reversed decision gets a new row, and a row that turns out to constrain future changes is promoted to a full record citing this table. Full text of demoted records lives in git history.

| Date | Decision | Why | Change/PR |
|---|---|---|---|
| 2026-07-16 | The digest is never a task in `tasks.md`; the digest precondition applies to the work only | A self-referential "digest" task can only be ticked falsely (before the deletion) or violate the precondition (unticked at deletion) — both misreadings have occurred; carrier is `clue-delta` step 4 | CH-017 / PR #16 (review finding) |
| 2026-07-12 | Default test-coverage gate: 80% total statement coverage, enforced in CI | A tripwire for code paths between the AC-covered ones — high enough to catch untested code landing, low enough not to invite meaning-free tests; total (not per-package) so thin entry points don't force contortions | demoted from ADR-004 by CH-010; carrier stays the CI workflow template |
| 2026-07-12 | Parse frontmatter with gopkg.in/yaml.v3 | A hand-rolled YAML subset would reject legitimate frontmatter with confusing errors, and the judge must never be wrong about form; one small pure-Go dependency doesn't threaten single-binary distribution | demoted from ADR-003 by CH-010 |
