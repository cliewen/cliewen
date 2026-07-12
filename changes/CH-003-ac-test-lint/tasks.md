---
id: CH-003-tasks
type: tasks
status: open
links: [CH-003]
title: Task breakdown for CH-003
---

# Tasks — CH-003

Ordered by dependency; ticked the moment each completes.

- [x] ADR-005: test-reference convention (test name carries ACnnn) — decides the rule's shape (AC-009, AC-010)
- [x] Extend CAP-002 criteria.md with AC-009 (active AC without test fails) and AC-010 (test referencing unknown AC fails)
- [x] Implement rule in internal/corpus: harvest @AC tags from criteria bodies, harvest ACnnn from *_test.go names, check both directions with draft exemption (AC-009, AC-010)
- [x] Tests for the rule: TestAC009_…, TestAC010_…, draft-exemption case (AC-009, AC-010)
- [x] Self-check: repo's own corpus passes (CAP-002 active ACs all have tests; CAP-001 draft exempt) — the check itself caught a missing ADR-005 index entry first, which was fixed, not suppressed
- [x] Update CAP-002 design.md: rule description, pair-enforcement door
- [ ] Digest: decisions index + M-002 → done in P-001, delete /changes/CH-003-ac-test-lint/
- [ ] Push, PR, CI green
