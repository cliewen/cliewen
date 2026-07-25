---
id: CH-061
type: change
status: open
links: [P-006]
title: MADR is a mapping, not improvisation
---

# CH-061 — MADR is a mapping, not improvisation

## What

Add `mappings/madr.md` to `clue-extract`, written against the target contract CH-060 just settled, so a repository whose decisions are MADR records can be converted mechanically instead of by improvisation (M-023 of [P-006](../../docs/plans/P-006-first-adoption.md)). The canonical source is `internal/skills/source/resources/clue-extract/mappings/madr.md`, alongside the existing `openspec.md`.

The mapping covers what the milestone names, in the shape `openspec.md` already uses — a layout line, a source-to-Cliewen table, and a "watch for" line:

- **Layout**: the MADR shapes an adopter actually carries — `docs/decisions/` or `doc/adr/`, MADR 3.x/4.x frontmatter (`status`, `date`, `decision-makers`, `consulted`, `informed`) and the template headings (Context and Problem Statement, Decision Drivers, Considered Options, Decision Outcome, Consequences, More Information), plus the older heading-only Nygard form the same folders often mix in.
- **Status-vocabulary conversion**: MADR carries `proposed`, `rejected`, `accepted`, `deprecated`, and `superseded by …`; a Cliewen decision has only `inferred` and `verified`, because provenance lives in the status ([ADR-010](../../docs/decisions/ADR-010-provenance-field.md)). The mapping states where each source status goes — every converted record is born `inferred` regardless, so the source status has to survive as meaning in the body and in links rather than as a status value.
- **`accepted-by` for acceptance that predates Cliewen**: a MADR record accepted by humans years before the corpus existed is real acceptance, but it is not the approval that the decision-records rule promotes on. The mapping preserves that acceptance — names, roles, and the original dates — as body prose, and leaves the record carrying `accepted-by: []`, the empty list every unsigned record already carries, because the field stays reserved for approval given under Cliewen's merge boundary ([PDR-004](../../docs/decisions/PDR-004-merge-binds-approval-signs.md)). The rule itself is an ADR (below), because it says what the field means on every record, not only on converted ones. Keeping one meaning in the field is what makes it machine-checkable later: [C-009](../../docs/constraints/C-009-type-specific-frontmatter.md) already records the promotion trigger where `checkCoreFields` grows a required-fields map covering `accepted-by`, and a field carrying two kinds of acceptance is a field no such rule can be written against.
- **ID preservation**: MADR's numeric filename prefix (`0005-use-postgres.md`) is the record's stable ID and survives as `ADR-005`, never renumbered; the mapping covers zero-padding, gaps, duplicate numbers across mixed folders, and records with no number at all. Contract item 3 governs criterion IDs in `ac-prefix:` namespaces and does not reach decision records, so the mapping states its own minting rule for unnumbered records rather than citing item 3 by analogy — generalizing the contract item is a contract edit past what M-022 settled.

Four carriers follow the new file. The generator emits mappings by a single hardcoded read of `openspec.md` (`internal/skills/generate.go`), so a second file in the resources folder ships to nobody until that read becomes a walk of `mappings/` — and the no-drift test would stay green while the mapping was invisible, so the generator change comes with a test that fails on an unemitted mapping. The contract's **Source mappings** paragraph in `internal/skills/source/skills/clue-extract.md.tmpl` names one current mapping and must name both. `guide/adoption.md` says "One extraction mapping ships today". And `mappings/openspec.md` carries a one-line MADR row of its own, which becomes a pointer to `madr.md` rather than a second, thinner statement of the same conversion. The generated skills (`.agents/skills/clue-extract/` and `internal/scaffold/templates/skills/clue-extract/`) are regenerated with no drift.

The Tank Royale conversion is the worked case the mapping cites: its extraction converted a real `docs/decisions` folder of MADR records ahead of any mapping existing, which is the improvisation this file replaces.

## Why

`clue-extract` says that if no mapping exists for a source format, writing one is the extraction PR's first task. MADR is the most common decision format an adopter arrives with, and today every such adopter pays that cost alone and answers the same three questions differently — most consequentially the status collapse, where a wrong answer either invents acceptance the corpus cannot back or throws away acceptance history that really happened. Tank Royale already paid it once. A mapping makes the second adopter's conversion mechanical and reviewable against a stated rule.

## Scope

In scope: `internal/skills/source/resources/clue-extract/mappings/madr.md`; mapping emission in `internal/skills/generate.go` and its test in `internal/skills/generate_test.go`; the **Source mappings** paragraph in `internal/skills/source/skills/clue-extract.md.tmpl`; the `accepted-by:` clause in `internal/skills/source/shared/decision-records.md.tmpl`; the MADR row in `mappings/openspec.md`, reduced to a pointer; the regenerated skill outputs; the extraction-mapping paragraph in `guide/adoption.md`; an ADR for the `accepted-by` rule, its entry in the `clue:index` block of `docs/decisions/README.md`, and a decision-log row for the status collapse; `docs/plans/P-006-first-adoption.md` bookkeeping; and `CHANGELOG.md`.

The `accepted-by` rule has three carriers, not one. The shared block in `decision-records.md.tmpl` is the skills' statement of it, but `docs/decisions/README.md` restates it for this corpus and `internal/scaffold/templates/docs/decisions/README.md` ships it into every repository `clue init` touches. All three are in scope: each is read against the ADR's wording and amended where it now says less than the rule does. Leaving the two READMEs out would let the corpus and the scaffold say something the skills no longer mean.

Out of scope: any change to `clue validate` or the decision status vocabulary itself — the mapping is written to fit the corpus rules, not to change them; re-litigating [ADR-010](../../docs/decisions/ADR-010-provenance-field.md); source formats other than MADR; the rest of `mappings/openspec.md` beyond its MADR row; automated migration tooling for an already-extracted corpus; and any change inside Tank Royale or another adopting repository.

## Open judgment calls

Two are settled inside this change, not raised as blocking questions, and they route differently.

What a `rejected`, `deprecated`, or `superseded by …` MADR record becomes when the target vocabulary has neither status is a decision-log row, on one condition: the rule it settles has to be lossless. What a later revision reverses cheaply is the mapping file, not the corpora already converted under it — an extraction deletes its source in the same PR, so a rule that drops meaning is not reversible where it matters, and routing it to a log row would be wrong.

Part of the question is already answered, and the row cites that answer rather than re-deciding it: [`docs/decisions/README.md`](../../docs/decisions/README.md) holds that rejected alternatives stay in the corpus and that decisions are never deleted, so a `rejected` record converts to a kept record. What is genuinely open is narrower — where `deprecated` and `superseded by …` land, and whether the source status survives as body prose, as a `links:` edge to the superseding record, or as both. Every candidate answer keeps the record and its meaning and differs only in which carrier holds it, which is what leaves the choice local to one source format and cheap for a later mapping revision to reverse.

Whether `decision-makers` may be written into `accepted-by:` without a human approving the conversion is an ADR ([ADR-029](../../docs/decisions/), next free ID). The guide's own routing table names frontmatter fields and extraction mappings as the ADR examples, and the answer states what `accepted-by:` means on every record, not only on converted ones. That reaches a methodology carrier five skills render. It does not touch the deterministic judge ([ARCH-003](../../docs/architecture/core.md)), so [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) does not apply and the ADR is born `inferred` like any agent-authored decision.

The ADR is not writing on a blank field. [PDR-004](../../docs/decisions/PDR-004-merge-binds-approval-signs.md) is what binds `accepted-by:` to promotion today — the shared skill block only restates it — and it settled the signing mechanics without ever facing acceptance that happened outside Cliewen. The ADR therefore names PDR-004 in its `links:` and says plainly which of its clauses it extends: signing stays exactly as PDR-004 defined it, and the new sentence is that the field admits nothing else. Read the other way — letting a populated `accepted-by:` sit on an `inferred` record — the ADR would contradict PDR-004 rather than extend it, and would need to supersede it outright. This change takes the extending reading, so PDR-004 stands.
