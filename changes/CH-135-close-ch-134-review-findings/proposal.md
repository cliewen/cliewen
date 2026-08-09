---
id: CH-135
type: change-proposal
links: [P-013]
title: CH-134's review findings are closed on the merged corpus
---

# CH-135 — CH-134's review findings are closed on the merged corpus

## What

Repair the seven review findings raised against CH-134 after its pull request was accepted. None of them is in the rule PDR-039 states; all of them are places where the merged corpus now reads inconsistently or points a reader at the wrong record.

- `CHANGELOG.md` states the dependent-change entry once, inside the `[Unreleased]` section's existing `Changed` list, so the published release body carries one heading rather than two.
- P-013's prose stops contradicting the M-065 row it now carries, and each of the three determinations states what its machinery costs a simplification campaign, as M-065's exit criterion requires.
- AN-013 carries the three determinations and links the records that made them, so `clue context AN-013` reaches PDR-039 and its closing paragraph no longer reports settled findings as open.
- The 2026-08-01 decision-log row that declined to record dependent work names the record that reverses it.
- PDR-039's decisions index row states what was decided rather than restating PDR-007's rule.
- PDR-017 carries an amendment note and link for the fourth acceptance-brief item PDR-039 adds to the brief it owns.
- `clue-delta`'s step 5 enumeration of the acceptance brief names the authorized unmerged base, matching both pull-request templates.

## Why

CH-134 landed the rule and left the corpus around it disagreeing with itself. A plan that tells two stories about the same milestone, an analysis that reports answered findings as open, a decision log whose live row is silently reversed, and an index row that describes a problem instead of an outcome all cost the next reader exactly what the corpus exists to save them. The `clue-delta` clause is the same carrier-inventory obligation PDR-019 states: both pull-request templates were updated in place and the procedural checklist an agent actually fills the brief from was not.

## Boundaries

No rule changes. PDR-039's decision text, the shared review-boundary fragment, C-012, CAP-006's criteria, and both pull-request templates are untouched, and no new artifact, check, or obligation is introduced. AC-130 stays as CH-134 wrote it.
