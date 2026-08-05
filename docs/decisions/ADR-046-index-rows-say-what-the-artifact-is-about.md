---
id: ADR-046
type: decision
status: inferred
links: [ADR-017, ADR-034, ADR-035, ADR-041, C-004, C-016, CAP-002, CAP-005, P-010]
title: An index row says what its artifact is about; the sentence is seeded, curated, and counted when absent
author: agent
accepted-by: []
---

# ADR-046 — An index row says what its artifact is about

## Context and problem statement

[ADR-041](ADR-041-index-rows-state-their-record.md) made a generated index row state its record rather than restate its own filename: `- [<id> — <title>](<file>) · \`<status>\``. A reader scanning a folder README now learns each artifact's identity and provenance without opening it. What a title cannot carry is what the artifact says, and that is usually the thing the reader is deciding on.

[C-016](../constraints/C-016-index-rows-state-their-record.md) already permits the sentence: curated text after the status is part of the row contract rather than an exception to it. Permission was not enough. Nothing produces such a sentence, nothing asks for one, and across this corpus's 135 artifacts not a single index row carries one.

The consequence is visible in adopted repositories rather than here. A repository of 101 documentation files onboarded on 2026-08-05 already maintained a hand-written index file in each of eight folders, because what it needed was not in a filename. After `clue init` it carried two indexes per folder, and the generated block listed the hand-maintained one among its own entries. Those eight files agreed on nothing except that each carried, under one heading or another, a sentence saying what the document is about.

An index that omits the one column an adopter will build a parallel index to obtain is an index that guarantees a second index.

## Decision outcome

**An appended index row carries a description seeded from the artifact's own body. The sentence is curated after that, and a row that has none is counted rather than failed on.**

- **The generator seeds; it does not assert.** When `regenIndex` appends a row it extracts one sentence: a lede paragraph directly beneath the H1 where one exists, otherwise the first sentence of the first paragraph under the first heading, skipping tables, lists, blockquotes, and fenced blocks. The row becomes `- [<id> — <title>](<file>) · \`<status>\` — <description>`.
- **An artifact with no readable sentence keeps today's row.** A row is one shape or the other and never a third carrying an empty tail, which is the rule ADR-041 already applies to a missing id, title, or status.
- **The sentence stays on one line.** A continuation line is not an entry `checkIndexes` recognizes, so extraction truncates at a sentence boundary rather than wrapping.
- **Regeneration rewrites nothing that exists.** This is ADR-041's rule, unchanged and load-bearing here: the seed is a first draft, and the author's correction must outlive the next `clue init`.
- **Nothing is backfilled.** No command writes a description into the rows already in a corpus.
- **The judge counts the remainder.** `clue validate` reports rows that state their record but say nothing about it as a population on its OK line, listable on request, never an `Issue`.

## Why extraction seeds rather than asserts

Extraction was measured against this corpus before it was adopted, and it is good for one artifact shape and systematically wrong for another.

A capability README opens with a purpose statement, so CAP-001 yields "A new user goes from installing `clue` to their first green `clue validate` in under 30 minutes." That is exactly the row a reader wants.

A decision opens under `## Context and problem statement`, and its first sentence states the problem in the present tense. ADR-041's own row would read "Index generation appended every missing entry as the target's filename with the extension removed" — the defect it removed, phrased as though it were still true. An extractor cannot tell a purpose statement from a problem statement, because both are declarative prose in the same position.

A seed an author corrects costs a bad first draft. An assertion written into 135 rows at once, of which a whole artifact type would be actively misleading, costs the reader's trust in every row including the correct ones. That asymmetry, and not the extractor's accuracy, is what settles this: the same extractor is acceptable as a seed and unacceptable as a truth.

It is also why no backfill happens. Backfilling is precisely the operation the asymmetry forbids.

## Why the remainder is counted and not failed on

This corpus has taken the same reading three times: [ADR-035](ADR-035-bounded-provenance-and-reality-edges.md) reports costly inferred meaning as an actionable population, [ADR-017](ADR-017-conventions-are-constraints.md) reads the constraint register the same way, and [ADR-041](ADR-041-index-rows-state-their-record.md) counts filler rows instead of rejecting them. The inferred-decision counter ran that experiment to completion: visible for months, it drove a campaign milestone and reached zero.

Failing here would also invert where the fault lies. Every description-less row in a Cliewen corpus is a row the generator wrote before this decision existed. A judge that rejected the tool's own prior output, in a file the adopter owns, would turn a Cliewen defect into the adopter's red build on upgrade.

Under [C-004](../constraints/C-004-never-weaken-checks.md) the count is never softened to make the number look better, and no command clears it, because regeneration deliberately preserves rows whose targets still exist.

## Rejected: a predefined column table

This change began as one, and it was the wrong shape twice over.

It was proposed against v0.12.0, where a row genuinely did restate its filename and a table looked like the only way to carry title, status, and description together. On `main` ADR-041 had already put the first two in the row, so the table's remaining contribution was the description alone — which a row can carry, and which C-016 already said a row may carry.

A table would also have had to supersede ADR-041's row shape and rewrite C-016 one release after both were decided, and it would have needed a migration to convert every adopted corpus, an extras-preservation rule for adopter columns, and a repair to the generator's habit of dropping any line inside the markers that carries no link, which silently removes a table's header and separator while the judge reports success. That defect is real and remains; nothing Cliewen emits is a table, so nothing here reaches it.

Columns are easier to scan than prose. They are not worth superseding a one-release-old decision, adding a corpus obligation, and carrying a migration, to deliver a sentence the existing row already had room for.

## Rejected: a date column

The first shape carried `File | Title | Status | Date | Description`, with the date read from the file's last commit. It was dropped on measurement.

Across the adopter's seventeen decision records the last-commit date is a single constant, because one bulk frontmatter migration touched them all. Across this corpus's sixty-five decision files, three dates cover thirty-eight. The column reports when someone last ran a sweep. The creation date repairs the clustering and reintroduces the same defect quietly, since `docs/` carries five renames whose adding commit records the move rather than the authorship.

Nor is there a corpus-wide field to read instead: the only date-bearing frontmatter is `accepted-by`, which exists on decisions alone and embeds its date inside prose.

The deeper reason is that a date is a weak proxy for something this corpus states exactly. Every artifact carries `status`; those whose trustworthiness is genuinely in question carry `provenance`, `reversal-cost`, and `reality`, and `clue validate --reality-gaps` reports on them. Recency belongs in a `reviewed:` field decided in its own change, not smuggled in as a column that looks authoritative and is not.

## Consequences

- A folder README states, for every artifact beside it, its identity, its status, and what it is about, so an adopter has no reason left to maintain a second index in the same folder.
- An existing corpus improves one row at a time, as artifacts are added and as the count is worked down by hand. It does not improve on upgrade, and this decision accepts that.
- The judge's OK line grows a second index population. An adopter reading it sees two distinct backlogs against the same rows, and they are worked down by the same hand edit.

**Carrier:** the appending branch and the extractor in `regenIndex` (`internal/scaffold/scaffold.go`), the description population in `corpus.IndexRowBacklog` and its CLI surface (`internal/corpus/index.go`, `cmd/clue/main.go`), [C-016](../constraints/C-016-index-rows-state-their-record.md), the `clue-extract` absorption guidance, and CAP-002's and CAP-005's criteria and design.
