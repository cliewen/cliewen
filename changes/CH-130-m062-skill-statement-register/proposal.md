---
id: CH-130
type: change
status: proposed
links: [P-013, M-062, PDR-029, PDR-013, PDR-028, ADR-021, C-011, C-013]
title: M-062 — the shipped skills and routing hub get a statement register
---

# CH-130 — M-062: the shipped skills and routing hub get a statement register

## Plan item

This change serves [P-013](../../docs/plans/P-013-simplification.md)'s **M-062**. The milestone remains wanted: it is the campaign's first milestone and M-063 cannot start without its register, because every previous simplification attempt argued case by case and stalled.

## What this change does

A spike walks the six generated skills (`clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, `clue-upgrade`, `clue-verify`) and `AGENTS.md` statement by statement, and produces a register as a findings document under `/docs/analysis/`.

For each statement the register records:

- **Class** — rule-bearing or connective. A heading, navigational sentence, worked example, or other statement carrying no obligation is recorded as connective and is not forced through the three carrier conditions.
- **Trace** — for a rule-bearing statement, the narrowest live corpus artifact that *states* the rule, or the fact that it traces to nothing. Per [PDR-029](../../docs/decisions/PDR-029-simplification-tests-by-surface.md), derivability is not a trace: a goal or general constraint counts only when the statement directly restates an obligation that artifact states.
- **Duplication** — whether a live carrier already states the same rule in the same reading path. Counted per reading path, not per file, so [ADR-021](../../docs/decisions/ADR-021-generated-standalone-skills.md)'s deliberate file-level repetition is not scored as a defect.
- **Checkability** — whether a reader can determine that it has been satisfied.
- **Order** — whether a statement that binds absolutely is read after the procedure it constrains.

Before any of that, the register fixes **what counts as a single statement**, precisely enough that an independent second pass which did not write the definition segments the same prose the same way.

Three classes escalate to the human rather than being decided here: a rule-bearing statement that traces to nothing, a statement whose decision may itself have outlived its reason, and a pair of obligations over one situation that could pull a reader in different directions. Each becomes an open question naming what it traces to or fails to, what removing or retaining it would cost, and what human judgment is required.

The register also records compatible overlap candidates for M-063 to consolidate, and recommends what durable form the register should take — including whether and how a row stays bound to a statement it does not annotate, and what leaving it unbound costs.

## What this change does not do

**Nothing is removed, reworded, or reordered.** M-063 does that, against this register and its answered questions.

**No trace, citation, or marker is added to a skill or to `AGENTS.md`.** PDR-029 puts the burden on the register, never on the carrier: the reader of a skill meets the rule, not its provenance.

**The register's durable form is recommended, not built.** A recommendation that `clue` enforce it changes what a green `clue validate` asserts and is core-adjacent under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md); building it would be a plan revision adding a milestone, backed by its own decision record and human acceptance.

## Tier

Full. The change writes a permanent corpus artifact and mutates plan bookkeeping.

## Reversal cost

Low and local. The register is a new analysis artifact and a milestone status flip; nothing it describes is edited. Reversing it deletes one file and restores one table cell.
