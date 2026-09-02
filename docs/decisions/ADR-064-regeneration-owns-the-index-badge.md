---
id: ADR-064
type: decision
status: inferred
links: [G-010, ADR-041, ADR-046, IDR-001, CAP-005, CAP-002]
title: Regeneration owns an index row's badge; the author owns its description
binds: adopter
author: agent
accepted-by: []
---

# ADR-064 — Regeneration owns the badge

## Context and problem statement

[ADR-041](ADR-041-index-rows-state-their-record.md) gives a generated row its record's identity and badge; [ADR-046](ADR-046-index-rows-say-what-the-artifact-is-about.md) adds a description that is seeded and then belongs to the author. `regenIndex` protects both the same way: a row that already covers its target is kept byte for byte.

That is right for the description and wrong for the badge. The description is a judgement a person made and a regenerator cannot improve; the badge is a copy of a frontmatter field, made once, and never checked again. Change an artifact's status and its row keeps the old value indefinitely. Nothing reports it — `clue validate` counts rows that say *nothing* about their artifact, not rows that say something *untrue* about it — and running the regenerator changes nothing, which is what makes a stale row convincing rather than merely wrong.

## Decision outcome

**A kept row's badge is refreshed from the artifact it links; everything else on the row is left exactly as it stands.** This amends [ADR-041](ADR-041-index-rows-state-their-record.md)'s clause that every existing row is preserved and [ADR-046](ADR-046-index-rows-say-what-the-artifact-is-about.md)'s that regeneration rewrites nothing existing: both now cover the row apart from its badge, and both say so. The value is whatever the row would have been appended with, so a constraint keeps showing its enforcement rather than its status ([IDR-001](IDR-001-constraint-index-badges-use-enforcement.md)).

Three boundaries keep this from trading a silent staleness for a silent overwrite:

- **A row carrying no badge gains none.** Its author removed it, and a row that claims nothing is never wrong.
- **A row covering more than one artifact is left alone.** No single artifact owns its badge, and the generator does not guess.
- **A row whose artifact has unreadable frontmatter is left alone.** Naming a malformed artifact is the judge's work, not the regenerator's.

**Rejected: make `clue validate` fail on a disagreement instead.** It would also catch a corpus where nobody runs the regenerator, and the goal accepts either answer. It is rejected because a new required check turns an adopter's next upgrade red for rows they did not write and could not have known about, in order to report something the regenerator can simply fix. The badge came from the generator, so the generator is where it stays current.

**Accepted cost:** regenerating an index now writes a file that previously would have been left alone, so a corpus with stale badges reports newly regenerated blocks on the first run after upgrading. That is the defect being repaired becoming visible once, not a new source of churn.

## Carrier

`refreshBadge` and the kept-row path in `regenIndex`, its criterion, this record, and the index guidance in `internal/scaffold/templates/docs/README.md` stating what a regenerated row keeps and what it refreshes.
