---
id: LOG-001
type: log
status: active
links: [ADR-016]
title: Decision log
---

# Decision log

Decisions that are cheap and local to reverse, one row each — newest first. The litmus test and demotion/promotion mechanics are [ADR-016](ADR-016-decision-log.md); decisions that constrain future changes are full ADRs in this folder. Rows are never deleted; a reversed decision gets a new row, and a row that turns out to constrain future changes is promoted to an ADR citing this table. Full text of demoted ADRs lives in git history.

| Date | Decision | Why | Change/PR |
|---|---|---|---|
| 2026-07-12 | Default test-coverage gate: 80% total statement coverage, enforced in CI | A tripwire for code paths between the AC-covered ones — high enough to catch untested code landing, low enough not to invite meaning-free tests; total (not per-package) so thin entry points don't force contortions | demoted from ADR-004 by CH-010; carrier stays the CI workflow template |
| 2026-07-12 | Parse frontmatter with gopkg.in/yaml.v3 | A hand-rolled YAML subset would reject legitimate frontmatter with confusing errors, and the judge must never be wrong about form; one small pure-Go dependency doesn't threaten single-binary distribution | demoted from ADR-003 by CH-010 |
