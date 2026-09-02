---
id: PDR-043
type: decision
status: verified
links: [CAP-004, ADR-039, ADR-043, ADR-060, PDR-042]
title: An upgrade routes as simple work unless the release makes the adopter decide
author: agent
accepted-by: Flemming N. Larsen (2026-09-02, conversation)
---

# PDR-043 — An upgrade routes as simple work unless the release makes the adopter decide

## Context and problem statement

The upgrade workflow rewrites managed surfaces and therefore looked full by path, even though its release contract had already been accepted upstream and the adopter's own accepted contract had not changed. Its release-availability check also cannot see a repository whose installed binary is current while its managed carriers are older, so treating that check as proof of repository currency leaves a coordinated migration undiscovered. When the release list is unavailable, stopping before the local preview leaves the same drift undiscovered and a clean preview cannot prove that no newer release exists.

## Decision outcome

**An upgrade is simple work.** It adopts an upstream release whose contract was accepted before publication, so the adopter's capabilities, criteria, decisions, plans, and constraints remain unchanged; it uses relevant surface checks and the repository's integration rules without full-change bookkeeping. Upgrade discovery runs `clue latest` for published release availability and then a no-apply `clue migrate` preview for repository carrier currency, even when the installed release is the newest or the release list is unavailable. A clean preview proves only that the managed carriers match the installed binary; the repository may be called current only when the reachable release list reports no newer release as well. If applying the release forces a decision about the adopter's own obligation, criteria, CI wall, or migration finding, that decision is routed as full work while the mechanical upgrade remains simple.

Changed-file count, diff size, and `/docs` or skill paths do not change the route. The human is still asked whether to upgrade before any write, and the human merge boundary remains unchanged.

## Rejected: route every upgrade through the full loop

That would make an adopter re-accept a release contract it did not author and turn a mechanical maintenance action into a false plan and evidence exercise.

## Carrier

The canonical and generated `clue-upgrade` skills, CAP-004, the skills architecture hand-off, routing hubs, `guide/operations.md`, and their tests carry this exception; PDR-042 and ADR-043 are linked rather than rewritten.
