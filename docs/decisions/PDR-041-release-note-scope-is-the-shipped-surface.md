---
id: PDR-041
type: decision
status: verified
links: [G-007, C-002, ADR-012, ADR-013, PDR-042]
title: A change owes a release note when it changes what an adopter receives
author: agent
accepted-by: Flemming N. Larsen (2026-09-02, conversation)
---

# PDR-041 — Release-note scope is the shipped surface

> **Terminology amended by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** simple/full routing replaces plain/light/full; route does not decide release-note scope.

## Context and problem statement

Because `CHANGELOG.md` is published as the release body, its scope must distinguish adopter-visible behavior from this repository's own corpus, contributor guidance, commands, and workflows.

## Decision outcome

**A change owes a release note when it changes what an adopter receives.** The test is whether it changes `clue` behavior, generated skill text, or an artifact materialized by `clue init` or `clue scaffold`; if none changes, no entry is owed. The test governs over category shortcuts, so a shipped file under `.github/` owes a note even when neighboring files are repository-local. A change spanning both layers describes only the adopter-visible part.

Whether an edit changes the meaning received by a user remains human judgment; release gates continue to require a reviewed version section rather than deriving the obligation mechanically from paths.

## Rejected: say only "user-visible" or list paths

Those answers leave adopter and contributor audiences ambiguous, and path lists go stale while comments, corpus changes, or generated output can change the shipped surface without changing their directory label.

## Carrier

C-002, the `AGENTS.md` release-note convention, and the digest guidance in `CONTRIBUTING.md` carry the test; generated skills remain generic and do not carry this repository's scope.
