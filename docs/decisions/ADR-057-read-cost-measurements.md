---
id: ADR-057
type: decision
status: inferred
links: [P-015, CAP-002, CAP-007, PDR-029, ADR-056, C-022]
title: Read cost is reported as a structural backlog, never scored by size
author: agent
accepted-by: []
---

# ADR-057 — Read cost is reported as a structural backlog, never scored by size

> **How the report is worked is stated by [ADR-058](ADR-058-read-cost-is-a-backlog-not-a-target.md):** the over-budget population is a backlog judged artifact by artifact, a `links` entry is never deleted to move the reported count, and an accepted identity keeps its row, so the count is not expected to reach zero. Everything this record decides is unchanged, including both populations, the eight-artifact budget, the exclusion of completed plans, and the refusal of a size score, a failing check, and token estimates.

## Context and problem statement

P-015 has bounded what `clue context` prints, but a corpus can still become costly to read in two ways the bounded command alone does not reveal: one artifact can become several rendered documents with different primary readers, and an identity can directly name enough artifacts that its ordinary one-hop slice is no longer focused. Both shapes become harder to repair after another campaign adds dependent meaning, so they need to be visible before they become a rule-breaking check.

The capability-folder convention already separates what a reader understands, what a tester verifies, and what an implementer builds. The same rule has not been stated for the rest of the corpus, so a long artifact can silently accumulate independent documents while preserving valid frontmatter and links. PDR-029 rejects a size score as a simplification measure: a larger file can be one coherent document, while a short file can make two readers discover each other's material.

## Decision outcome

**Every durable artifact serves one primary consumer, and `clue validate` reports the two structural signs that it does not.**

The first reported population is a live artifact under `docs/` whose body contains more than one rendered H1 outside fenced examples, counting both forms Markdown renders as an H1: an ATX `# ` line and a setext title over a `=` underline. An H1 begins a rendered document; a second H1 is therefore the observable boundary that a split can preserve. Completed plans are excluded by their `completed` state because [C-008](../constraints/C-008-completed-plans-immutable.md) makes rewriting their historical record unavailable. This is a report, not a verdict: later work decides whether the documents should be split or whether one primary consumer genuinely needs them together.

The second reported population is every durable identity whose ordinary `clue context` slice prints more than eight artifacts: the root and its unique resolvable direct links. Acceptance-criterion and milestone identities count because the command accepts them as entry points, so the reported count is a count of entry paths, not of files: one plan that declares milestones contributes a row per identity, all measuring the same `links:` list, and repairing that list clears every one of them. Identities owned by a completed plan are excluded for the same reason as the first population: the slice is wide because of a frozen `links:` list, and C-008 leaves no permitted repair, so reporting them would name a backlog that cannot be worked. The measurement intentionally calculates only that bounded slice, not its frontier or wider closure; a validator should not traverse the rest of the graph merely to report the cost of what a reader sees by default.

Both populations appear as counts on every successful `clue validate` run and are named by `clue validate --read-cost`. Neither is an `Issue` and neither changes the command's exit code. The command derives them fresh from the repository state and writes no registry.

This is a structural navigation measure, not a size score. It does not claim that H1 count, artifact count, bytes, or words measure quality; it answers the narrower, checkable questions of whether one file visibly contains multiple rendered documents and whether the default entry path prints more than the stated focused-slice budget. PDR-029 remains the test for whether a carrier or rule should exist at all.

## Rejected: per-type byte budgets

A decision, an analysis, and a criterion have different legitimate densities. A byte cap would reward compression and punish a readable explanation without saying whether either file has more than one reader, reproducing the size score PDR-029 rejects.

## Rejected: a failing read-cost check

The existing corpus contains the population this change is meant to expose, and a hard failure would make repair order a condition of unrelated work. Reporting makes the backlog inspectable first; P-015/M-072 owns the deliberate repairs.

## Rejected: token estimates

Tokenizer choice, model changes, and language make a token total an unstable proxy for reader effort. The reported shapes are visible in the corpus itself and deterministic across installations.

## Rejected: measure only new artifacts

Exempting history would make the oldest and most expensive corpus material permanently invisible, while measuring it only from a diff would violate the judge's state-only boundary. The report reads the current durable corpus uniformly; completed plans are the one state-based exception, applied to both populations, because their immutability is already an explicit rule.

## Carrier

CAP-002's README, criteria, and design describe the `clue validate` report and flag. CAP-007 supplies the bounded one-hop slice the budget measures; its contract is unchanged. C-022 states the rule and the machine-visible subset beside the residual judgment. The CLI usage text names `--read-cost`.
