---
id: CH-111-proposal
type: proposal
status: active
links: [P-010, M-051, CAP-001, CAP-002, ADR-019, ADR-034, ADR-039]
title: The folder index is a table an agent can triage from
---

# CH-111 — The folder index is a table an agent can triage from

Serves [P-010](../../docs/plans/P-010-adopters-keep-current.md) milestone **M-051**.

## What

An index exists so a reader can see a link and decide, without opening the file, whether opening it is worth it. Cliewen's generated index cannot support that decision: a line reading `- [ADR-014-upload-delete-consistency](ADR-014-upload-delete-consistency.md)` restates the filename and says nothing else. An agent orienting in a decisions folder has to open every file to find the one it needs, which is the cost the index was supposed to remove.

This change makes the generated block a table with a predefined column prefix: `File`, `Title`, `Status`, `Description`. The first three are derived from frontmatter the corpus already requires, so they cannot drift from the artifacts they describe. The fourth is extracted from the artifact body. Columns an adopter adds after the prefix are preserved across regeneration.

## Why the current shape produces two indexes in practice

Observed on 2026-08-05 while adopting v0.12.0 into a real repository of 101 documentation files. That repository already maintained its own index files, one per folder, because the information it needed was not in a filename. After `clue init` it carried two indexes per folder, and the generated block listed the hand-maintained one as one of its own entries. The hand-maintained indexes were richer than the generated block in every folder, and they did not agree with each other on shape: five different column sets across eight files, including `Title, Status, Date, Design`, `Title, Scope, Related ADRs`, `Description`, and `Purpose`.

That divergence is the argument for a predefined prefix rather than a free-form table. Columns an adopter invents are unusable for triage precisely because they vary: an agent cannot rely on a column that exists in one folder and not the next. A guaranteed prefix in a fixed order is what makes a row readable without reading the corpus that produced it.

Two of those adopter columns, `Design` and `Related ADRs`, are not description at all. They are graph edges, and Cliewen already models edges in `links:` and walks them with `clue context`. Folding them away is a correction rather than a loss, and the extraction guidance says so.

## The decapitation

`regenIndex` keeps a line inside the markers only when it carries a link to a wanted target, and drops every other line (`internal/scaffold/scaffold.go:515`). A table's data rows each carry such a link and survive. Its header row and separator carry none and do not. A table written by hand therefore survives exactly until the next `clue init`, `scaffold`, or `migrate`, after which the rows remain with no header above them and render as broken prose.

`clue validate` does not notice, because `checkIndexes` judges three things — markers present, every link resolving, every sibling covered — and never shape. Both halves were reproduced against the adopter repository: a hand-written table validated with no new issue, and the next `clue scaffold` removed its header and separator while the verdict stayed green. The current contract is therefore not merely list-shaped by preference; it actively destroys the alternative while the judge reports success.

## Why there is no date column

The first draft of this change carried `Date`, sourced from the file's last commit. It was dropped on measurement. Across the adopter's seventeen decision records the last-commit date is a single constant, because one bulk frontmatter migration touched them all; across Cliewen's own sixty-five decision files, three dates cover thirty-eight of them. A column whose value is the date of the most recent sweep carries no triage signal.

The creation date repairs the clustering and reintroduces the same defect quietly, since `docs/` has five renames in its history and the adding commit for those files records the move. Nor is there a corpus-wide field to read instead: the only date-bearing frontmatter in the corpus is `accepted-by`, which exists on decisions alone and embeds its date inside prose.

The deeper reason is that a date is a weak proxy for something this corpus already states exactly. Every artifact carries `status`; the ones whose trustworthiness is genuinely in question carry `provenance`, `reversal-cost`, and `reality`, and `clue validate --reality-gaps` reports on them. Recency belongs in a `reviewed:` field decided in its own change, not smuggled in as a column that looks authoritative and is not.

## The shape

**One decision.** [ADR-019](../../docs/decisions/ADR-019-init-regenerates-indexes.md) fixed the contract as "hand-written single-line entries whose targets survive are kept". An ADR replaces that one clause with the predefined prefix, the split between the columns the writer owns and the columns the adopter owns, and the extraction rule. The rest of ADR-019 stands, so this is a clause supersession and not a retirement under [ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md).

**A table mode in the regeneration engine.** A block whose first two non-empty lines are a header and a separator is a table; anything else stays a list and behaves exactly as it does today. In table mode the header and separator are preserved, each row is keyed by the link target in its first cell, cells one through four are rewritten from the artifact, cells five onward are preserved verbatim, rows whose target no longer resolves are dropped as they are today, and a missing target is appended as a row with a derived prefix and empty extras.

**A description extracted from the body.** A lede paragraph directly beneath the H1 and before any subheading is used when one exists. Where there is none, the first sentence of the first paragraph under the first heading is used, truncated at a sentence boundary. Every fallback is named in the report, so a weak row is visible rather than silent, and writing a lede becomes a followable convention instead of an unstated one.

**One check in the judge.** When the block is a table, its first four headers must be the predefined names in the predefined order. Without it an adopter may reorder columns and quietly break the only guarantee that makes a row readable. `checkIndexes` stays otherwise shape-permissive, and the judge gains no new dependency: every derived cell comes from frontmatter already scanned, and nothing reads git.

**A migration.** MIG-008 converts a list block Cliewen itself emitted into the table shape, which is exactly the scope [ADR-039](../../docs/decisions/ADR-039-versioned-corpus-migrations.md) defines. Preview by default, `--apply` to write, idempotent on a corpus already converted, and prose outside the markers is untouched.

**Extraction guidance.** `clue-extract` learns to fold a foreign index file into its sibling README block, map its columns onto the prefix, keep what genuinely remains as extras, drop the ones that restate `links:`, and delete the file it absorbed. This stays with the skill rather than becoming a migration because the five schemas observed in one adopter differ, and mapping them is judgment.

## Result

A folder README states, for every artifact beside it, what the artifact is called, what it is titled, what its status is, and what it is about. An agent choosing what to read makes that choice from the index. An adopter who wants more columns keeps them, and the shape survives regeneration instead of being silently decapitated by it.

## Scope

Full tier. It changes what a green `clue validate` asserts and it amends a decision, so [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) requires the decision record above and human acceptance at merge. The index contract is a methodology contract under [C-006](../../docs/constraints/C-006-adrs-timeless-with-carrier.md), so every live carrier moves in this change: the eight scaffold templates, the generated `clue-extract` skill and its mappings, CAP-001's criteria and design, CAP-002's index rule, the migration inventory, and `[Unreleased]`.

## Not in this change

No `reviewed:` field and no recency column. No general table support outside the taxonomy README index blocks. No change to how `checkIndexes` judges coverage or link resolution, and no weakening of either. `clue` gains no git access, holding the boundary [ADR-042](../../docs/decisions/ADR-042-release-check-outside-the-judge.md) and [ADR-044](../../docs/decisions/ADR-044-judge-reads-state-not-transitions.md) drew.
