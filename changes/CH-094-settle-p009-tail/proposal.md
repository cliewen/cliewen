---
id: CH-094
type: change
status: open
links: [P-009, M-042, M-043, M-044, AN-013, AN-015, ADR-038]
title: Settle P-009's tail so migration readiness closes on migration alone
---

# CH-094 — Settle P-009's tail so migration readiness closes on migration alone

## What

This full change settles the three open milestones P-009 still carries, so the campaign can close on the work it was named for.

M-042 becomes `done`. Its one unmet exit conjunct — the proposed thin-caller wall running in the target CI shape — was closed by real evidence after AN-015 landed: the human-selected target replaced its copied wall with the upstream caller, upgraded the managed carriers through `clue migrate`, and merged that change green. The milestone's evidence records both AN-015 and that adopter merge.

M-043 is removed. Recording a dependent change's base and human authorization as durable data answers a gap AN-013 measured, but the review boundary already requires an agent to stop and ask before building on unaccepted work, and one Cliewen change is in flight per author. The residual gap is that the authorization stays conversational rather than readable in history. This change states that boundary rather than building an interface for a situation the supported path already discourages.

M-044 stays and is not implemented here. Naming foreign acceptance evidence is the gap the human meets in practice: an adopter's tests already prove Cliewen claims with no honest way to write that down. It gets its own change.

The removal is a plan revision, so it is routed through a decision record under P-009's mutation rules. No capability, criterion, validator behavior, skill, or shipped contract changes in this change.

## Why

P-009 bundled migration readiness with two distributed-work contracts that are not prerequisites for one repository to migrate truthfully — the campaign says so itself. Leaving them in keeps an otherwise finished campaign open and pushes its successor behind work the human has decided differently about: execute the foreign-evidence form, decline the dependent-change form.

## Decision boundary

This change adjusts plan scope and records the declined boundary. It does not implement M-044, weaken the review boundary's stop-and-ask rule for unaccepted bases, remove AN-013's finding, or claim the dependent-change gap does not exist.
