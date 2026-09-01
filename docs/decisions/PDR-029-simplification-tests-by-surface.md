---
id: PDR-029
type: decision
status: verified
links: [P-013, PDR-013, PDR-019, PDR-025, PDR-026, ADR-021, ADR-034, AN-006, AN-008, AN-010, AN-012, G-001, C-011, C-013, C-016]
title: Simplification is judged by two tests, chosen by surface
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-029 — Simplification is judged by two tests, chosen by surface

> **Amended by [ADR-059](ADR-059-progressive-standalone-skill-directories.md), [PDR-037](PDR-037-tooling-is-judged-by-what-it-holds.md), and [PDR-031](PDR-031-architecture-artifacts-are-traces.md):** the generated skill may be a progressive standalone directory; tooling is judged by whether removal hands work to a human; and architecture is a valid trace when it states the rule. The other tests and boundaries below remain unchanged.

## Context and problem statement

PDR-013's *does the core need it?* test is right for meaning but passes nearly every rule-bearing carrier statement and most tooling because those surfaces cost repetition, weak checkability, missing traces, or human memory rather than core guarantees.

## Decision outcome

**Simplification is judged by two tests, and the surface chooses the test.** *Does the core need it?* governs rules, artifact types, required fields, commands, checks, and anything whose existence changes what the method obliges. On the tooling surface, PDR-037's *does removing it move an obligation from a machine back to a human?* keeps a mechanism when the answer is yes and makes it a removal candidate when the answer is no; where both apply, the core test is decisive.

Carrier prose is judged statement by statement: every rule-bearing statement must trace to the narrowest live decision, constraint, goal, acceptance criterion, or architecture artifact that states it, state its rule once in the reading path, and admit checking. Connective prose is outside this test. Derivability is not tracing, and a statement that traces nowhere becomes a human open question rather than an automatic deletion or retention.

Binding statements come first. Compatible overlap is consolidated; conflicting obligations are escalated for human decision. Word count is not the measure: splitting an unreadable sentence or adding a missing decision can increase the corpus. Removing or weakening corpus meaning, the verifiable thread, the human merge boundary, or what `clue validate` asserts requires its own decision under C-013; this record authorizes the tests, not a removal.

## Rejected: apply the core test to every surface

It returns "needed" for almost every carrier rule and misses repetition, checkability, tracing, order, and the human-holder test for tooling, so it defers simplification rather than judging it.

## Rejected: measure success in words or artifacts removed

A deletion count rewards compressing a rule until it cannot be checked and cannot measure reordering or writing a missing decision; the surviving rule set, not its size, is the outcome.

## Carrier

This record and its amendment note, P-013's milestones, PDR-013's refinement, and the curated index rows carry the tests. The skills, hubs, guides, contributor guidance, and CLI text are the surfaces P-013 applies them to; their carrier edits belong to that campaign.
