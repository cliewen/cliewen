---
id: CH-014-tasks
type: tasks
status: open
links: []
title: CH-014 task breakdown
---

# Tasks

- [ ] Write PDR-006 (decision records are typed: ADR/PDR/log; well-established practices are cited, not re-derived), born `inferred`
- [ ] Rename ADR-014/015/016/018/019 → PDR-001/002/003/004/005 (git mv; ids, titles, cross-links; superseded notes repointed)
- [ ] Clerical signing: PDR-005 → `verified`, accepted-by Flemming N. Larsen (2026-07-15)
- [ ] Write C-011 (routing rule, `enforcement: agent`); update constraints README index
- [ ] Rewrite decisions README (three tiers) and regenerate its index block
- [ ] Update cross-references: AGENTS.md, docs/README.md status table, P-002 (M-007 link + mutation rule), CAP-002 design if it cites a moved ID
- [ ] Update skill wording: clue-delta, clue-verify, clue-plan, clue-analysis
- [ ] CHANGELOG entry under [Unreleased]
- [ ] Digest: delete /changes/CH-014-decision-record-taxonomy/
- [ ] Verify: clue validate green (--forbid-changes at digest), go test ./... green, no stale ADR-01[4-9] references, clue-verify checklist walked
