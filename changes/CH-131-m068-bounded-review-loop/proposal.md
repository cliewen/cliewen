---
id: CH-131
type: change
status: open
links: [P-013, M-068, PDR-012, PDR-029, ADR-021, C-004]
title: M-068 — bound the agentic review loop
---

# CH-131 — M-068: bound the agentic review loop

## Plan item

This change serves [P-013](../../docs/plans/P-013-simplification.md)'s **M-068**. The milestone remains wanted: every remaining campaign change pays the review loop's cost, and bounding that cost before M-063 keeps the campaign's largest carrier edit affordable.

## What this change does

The change amends [PDR-012](../../docs/decisions/PDR-012-agentic-review-before-publication.md) through PDR-035 so the agentic review loop has its own severity model and a finite ordinary budget:

- The caller's reviewer brief cannot redefine blocking and advisory findings or turn an advisory class into a publication gate.
- A finding whose substance is a count, total, population figure, or arithmetic disagreement is advisory; wrong, missing, or reused identities remain blocking because allocating an identity is not counting.
- The reviewer spends no pass re-deriving figures. Authors remain responsible for computed figures before publication.
- Three passes are the ordinary budget. A fourth or later pass runs only when the immediately preceding pass returned a blocking finding; a pass with no blocking findings ends the loop.
- The verification handoff reports the number of passes run as well as the review mode and reviewed commit.

The standing rule is registered as C-017 with `enforcement: human`. The canonical `clue-verify` source states it, `go generate ./internal/skills` publishes it into both generated trees, and focused generator assertions pin the stable contract.

## What this change does not do

It does not restate the rule in `AGENTS.md`, the scaffolded routing hub, or another skill. Those surfaces already route to `clue-verify`; another copy would duplicate one reading path.

It does not change what `clue validate` judges, the human merge boundary, or acceptance-criterion evidence. No executable mechanism observes a running agentic loop, so C-017 honestly records human enforcement rather than claiming a machine gate.

It does not perform M-063's broader skill trim or reorder unrelated `clue-verify` prose.

## Tier

Full. The change amends project methodology, adds a decision and constraint, and edits a generated methodology carrier.

## Reversal cost

High. The review exit condition controls which candidate may be published, applies to every future Cliewen change, and ships to adopters. Reversal therefore requires another PDR and coordinated carrier update.
