---
id: CH-083
type: change
status: open
links: [P-008, M-036, G-001, AN-011, AN-012, ADR-009, C-012, C-013]
title: Reproduce the distributed work and evidence boundaries before they become contracts
---

# CH-083 — Reproduce the distributed work and evidence boundaries before they become contracts

## What

Serves **M-036** of [P-008](../../docs/plans/P-008-self-consistency.md), the campaign's last milestone. A spike exercises three pinned scenarios and ends in one analysis record under `/docs/analysis/`:

1. **One authorized dependent change** — a change whose work builds on another change that `main` has not accepted, which the review boundary permits only under explicit human authorization.
2. **One capability whose acceptance evidence spans repositories** — a Cliewen criterion whose real carrier or proof lives in an adopting repository rather than in this one.
3. **One external-tracker reference that must survive a repository move** — a corpus reference to forge state that has to keep meaning the same thing after the referenced repository is renamed, transferred, or mirrored.

The analysis records which guarantees today's branch-and-pull-request boundary, repository-local `clue validate`, and the stable-ID rules can and cannot supply for each scenario; rejects any option that weakens the human merge boundary or makes forge state the system-of-record; and ends with independently routable candidates for stacked changes, cross-repository evidence, and tracker metadata, plus a named successor-plan consumer.

## Why

M-036 is the second half of the same discipline M-035 applied to configuration: Cliewen deferred three distributed-work doors behind P-007's "a real pilot must price it first" line, and an interface built now would be built on speculation. [AN-012](../../docs/analysis/AN-012-adopter-configuration-cost.md) has just shown what that discipline is worth — two of the three configuration assumptions it priced turned out to be solved already and measured at zero cost, and the leading hypothesis was rejected on evidence. The three boundaries here are more dangerous than configuration, because each has an obvious implementation that would quietly damage the core: stacked changes invite treating a pull-request base as accepted meaning, cross-repository evidence invites trusting a foreign forge's green check, and tracker metadata invites recording forge identity as corpus truth. [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) makes all three core-adjacent, so the investigation must run before any interface is proposed.

## Scope

No interface is implemented. This change adds one analysis record, routes its outcomes as candidates (recording only decisions this change actually makes), and updates M-036's plan bookkeeping in the digest. It changes no skill contract, no validator behavior, and no adopting repository.
