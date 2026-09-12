---
id: CH-187-tasks
type: tasks
status: open
links: [CH-187]
title: Tasks — challenge a consequential commitment before making it
---

# CH-187 — Tasks

- [ ] Write `docs/decisions/PDR-061-challenge-consequential-commitments-proportionally.md`: what counts as a consequential commitment, the four questions and the satisfied-but-failing question, why the proportionality limit is part of the rule, that a milestone may deliver decision-changing evidence, that repository experience is evidence to reassess rather than standing authority, and that the rule binds this repository only until P-023/M-100; no `binds: adopter`.
- [ ] Add the obligation to `AGENTS.md` as one section reached before planning, stating the rule and its limit together and linking PDR-061 rather than restating its reasoning; keep it clear of the source-repository conventions table so M-100 can move it without rewriting.
- [ ] Link the new section from the "Work from durable context" planning paragraph in `AGENTS.md` only if the section's placement does not already put it in front of an agent about to plan.
- [ ] Regenerate `docs/decisions/README.md` with `clue scaffold` and write PDR-061's index row as a descriptive sentence.
- [ ] Assess documentation impact: `docs/design/README.md` describes cross-cutting methodology flow and may need one sentence locating the challenge before a commitment; `guide/` is the public description of the shipped method and stays unchanged because the rule does not ship yet. No `CHANGELOG.md` entry: nothing on the shipped surface changes.
