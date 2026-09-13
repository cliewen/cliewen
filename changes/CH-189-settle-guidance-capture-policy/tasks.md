# Tasks

- [ ] Write `docs/decisions/PDR-062-guidance-capture-has-an-eligibility-bar-and-no-new-type-yet.md`: the eligibility bar, the two existing homes and why no new corpus type ships yet, evidence sufficiency scoped to what was observed, discovery riding existing read paths, retirement as ordinary document correction, and the three guards named as guards; link `P-023`, `M-096`, `G-015`, `AN-024`, `ADR-026`; no `binds: adopter` (this repository only, per M-100's ordering).
- [ ] Regenerate `docs/decisions/README.md` with `clue scaffold` and write PDR-062's index row as a descriptive sentence.
- [ ] Assess documentation impact: `docs/design/README.md`'s cross-cutting flow is unaffected — this decision governs where a non-core artifact is filed, not the delivery thread itself — so no addition there. No `CHANGELOG.md` entry: nothing an adopter receives changes.
- [ ] Update P-023's M-096 row to `done` with evidence pointing at PDR-062.
- [ ] Digest: give `tasks.md` and `open-questions.md` their frontmatter, mark `CH-189` and `PDR-062` live and hand-append `M-096`'s retirement in `.clue/id-ledger.yaml` — no command writes any of the three yet ([G-019](../../docs/goals/G-019-workspace-identities-need-no-hand-edited-ledger.md)), matching CH-186/187/188's practice — then delete this change workspace.
