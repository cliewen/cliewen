---
id: CH-111-proposal
type: proposal
status: active
links: [P-010, M-051, CAP-001, CAP-002, CAP-005, ADR-019, ADR-034, ADR-041, C-016]
title: An index row says what the artifact is about, and a foreign index is absorbed rather than left beside it
---

# CH-111 — An index row says what the artifact is about

Serves [P-010](../../docs/plans/P-010-adopters-keep-current.md) milestone **M-051**.

## What

[ADR-041](../../docs/decisions/ADR-041-index-rows-state-their-record.md) made a generated index row state its record: `- [<id> — <title>](<file>) · \`<status>\``. That is most of what a reader needs in order to decide whether to open a file. The piece still missing is the one a title cannot carry: what the artifact actually says.

[C-016](../../docs/constraints/C-016-index-rows-state-their-record.md) already permits it — "curated text after the status is part of the contract rather than an exception to it" — but nothing produces it, nothing asks for it, and so nobody writes it. Across this corpus's 135 artifacts, not one index row carries a description.

This change makes the generator seed a description when it appends a row, extracted from the artifact's own body, and makes the judge count the rows that lack one. The description is curated after that, exactly as C-016 already says.

## Why an index without it produces a second index

Observed on 2026-08-05 while adopting v0.12.0 into a repository of 101 documentation files. That repository already maintained its own index file in each folder, one per folder, eight in total, because what it needed was not in a filename. After `clue init` it carried two indexes per folder, and the generated block listed the hand-maintained one as one of its own entries.

Those eight files did not agree on shape: five different column sets, including `Title, Status, Date, Design`, `Title, Scope, Related ADRs`, `Description`, and `Purpose`. What every one of them carried, under one name or another, was a sentence saying what the document is about. That is the column an adopter builds a parallel index to get, and it is the column Cliewen does not generate.

Two of those columns, `Design` and `Related ADRs`, are not description at all. They are graph edges, and Cliewen models edges in `links:` and walks them with `clue context`. Absorbing an index therefore drops them rather than preserving them.

## Why the row is seeded and never asserted

Extraction was measured against this corpus before being adopted. It reads a lede paragraph beneath the H1 where one exists, and otherwise the first sentence of the first paragraph under the first heading.

For capability READMEs the result is genuinely good, because they open with a purpose statement: CAP-001 yields "A new user goes from installing `clue` to their first green `clue validate` in under 30 minutes." For decisions it is systematically wrong, because an ADR opens under `## Context and problem statement` and the first sentence states the problem in the present tense. ADR-041's own row would read "Index generation appended every missing entry as the target's filename with the extension removed", advertising the defect it removed.

So extraction seeds a row and is never treated as truth. The author edits it, and regeneration never rewrites it, which is [ADR-041](../../docs/decisions/ADR-041-index-rows-state-their-record.md)'s existing rule and not a new one. Nothing is backfilled into the rows already in this corpus: a first draft an author corrects is a fair cost, and a wrong sentence written into 135 rows at once is not.

The remainder is handled the way this corpus already handles a backlog three times over — [ADR-035](../../docs/decisions/ADR-035-bounded-provenance-and-reality-edges.md) for costly inferred meaning, [ADR-017](../../docs/decisions/ADR-017-conventions-are-constraints.md) for the constraint register, and ADR-041 for filler rows. The judge counts rows with no description on its OK line and lists them on request. It never fails on them, because every such row was written by the generator.

## The shape

**One decision.** An ADR extends ADR-041: an appended row seeds a description from the body, the description is curated thereafter, extraction is a seed and not an assertion, and the rows without one are a counted population rather than an error. C-016 gains the description clause. ADR-019's contract is untouched, and no decision is retired.

**The generator.** `regenIndex` appends `- [<id> — <title>](<file>) · \`<status>\` — <description>`. An artifact with no extractable sentence keeps today's row exactly, because a row is one shape or another and never a third carrying an empty tail. `corpus.RowIdentity` already reads the frontmatter this needs.

**The judge.** The existing index-row backlog in `internal/corpus/index.go` gains a second population: rows that state their record but say nothing about it. Reported on the OK line, listable, never an `Issue`. The judge reads no new source and gains no dependency.

**Extraction guidance.** `clue-extract` learns to absorb a foreign index file into its sibling README block, map its description column onto the row tail, drop the columns that restate `links:`, and delete the file it absorbed under [ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md). This stays with the skill rather than becoming a migration because the five schemas observed in one adopter differ, and mapping them is judgment.

## Result

A folder README states, for every artifact beside it, its identity, its status, and what it is about. An adopter has no reason left to maintain a second index, and the one they already have can be absorbed and deleted rather than left to drift.

## Scope

Full tier. It extends a decision and adds a counted population to the judge's output, so [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) requires the decision record above and human acceptance at merge. The row contract is a methodology contract under [C-006](../../docs/constraints/C-006-adrs-timeless-with-carrier.md), so every live carrier moves in this change: C-016, CAP-002's and CAP-005's criteria and design, the generated `clue-extract` skill and its mappings, the CLI surface, and `[Unreleased]`.

## Not in this change

**No table.** A predefined column table was the first shape considered and was rejected once ADR-041 and C-016 were read: they settled the row's shape one release ago, and the information a table would carry is the information this change adds to the row.

**No backfill and no migration.** Nothing rewrites an existing row, so no corpus obligation is added or narrowed and the migration registry is untouched.

**No date column.** Sourcing it from the last commit measured inert — a single constant across the adopter's seventeen decision records, and three dates covering thirty-eight of this corpus's sixty-five — the creation date carries the same defect through renames, and no corpus-wide date field exists. Recency belongs in a `reviewed:` field decided in its own change.

**The link-less line drop stands.** `regenIndex` removes any line inside the markers that carries no link to a wanted target, so a hand-written table loses its header and separator while `checkIndexes` reports success. Nothing Cliewen emits is a table, so this change does not reach it. It is a real defect and is named here so it is not rediscovered as a surprise.
