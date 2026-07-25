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
- **`accepted-by` for acceptance that predates Cliewen**: a MADR record accepted by humans years before the corpus existed is real acceptance, but it is not the approval that the decision-records rule promotes on. The mapping states what the extraction may write into `accepted-by:` and what stays body prose, and keeps the original acceptance dates.
- **ID preservation**: MADR's numeric filename prefix (`0005-use-postgres.md`) is the record's stable ID and survives as `ADR-005`, never renumbered; the mapping covers zero-padding, gaps, duplicate numbers across mixed folders, and records with no number at all.

Two carriers follow the new file: the contract's **Source mappings** paragraph in `internal/skills/source/skills/clue-extract.md.tmpl` names one current mapping and must name both, and `guide/adoption.md` says "One extraction mapping ships today". The generated skills (`.agents/skills/clue-extract/` and `internal/scaffold/templates/skills/clue-extract/`) are regenerated with no drift.

The Tank Royale conversion is the worked case the mapping cites: its extraction converted a real `docs/decisions` folder of MADR records ahead of any mapping existing, which is the improvisation this file replaces.

## Why

`clue-extract` says that if no mapping exists for a source format, writing one is the extraction PR's first task. MADR is the most common decision format an adopter arrives with, and today every such adopter pays that cost alone and answers the same three questions differently — most consequentially the status collapse, where a wrong answer either invents acceptance the corpus cannot back or throws away acceptance history that really happened. Tank Royale already paid it once. A mapping makes the second adopter's conversion mechanical and reviewable against a stated rule.

## Scope

In scope: `internal/skills/source/resources/clue-extract/mappings/madr.md`, the **Source mappings** paragraph in `internal/skills/source/skills/clue-extract.md.tmpl`, the regenerated skill outputs, the extraction-mapping paragraph in `guide/adoption.md`, `docs/plans/P-006-first-adoption.md` bookkeeping, `CHANGELOG.md`, and a decision-log row for the status-collapse and `accepted-by` rules.

Out of scope: any change to `clue validate` or the decision status vocabulary itself — the mapping is written to fit the corpus rules, not to change them; re-litigating [ADR-010](../../docs/decisions/ADR-010-provenance-field.md); source formats other than MADR; any edit to `mappings/openspec.md`; automated migration tooling for an already-extracted corpus; and any change inside Tank Royale or another adopting repository.

## Open judgment calls

Two are settled inside this change and recorded as a decision-log row, not raised as blocking questions: what a `rejected` or `superseded` MADR record becomes when the target vocabulary has neither status, and whether `decision-makers` may be written into `accepted-by:` without a human approving the conversion. Both are periphery — a mapping describes one source format and does not touch the deterministic judge ([ARCH-003](../../docs/architecture/core.md)) — so a row, not an ADR.
