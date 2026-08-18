---
id: PDR-042
type: decision
status: verified
links: [G-001, G-008, PDR-002, PDR-007, PDR-011, PDR-012, PDR-018, PDR-030, PDR-040, PDR-041, C-002, C-012, C-013, CAP-006]
title: Routing recommends effort from accepted-contract change while the user retains integration authority
author: agent
accepted-by: Flemming N. Larsen (2026-08-13, planning conversation)
---

# PDR-042 — Routing recommends effort from accepted-contract change

## Context and problem statement

Path and runtime heuristics made plain/light/full routing confuse semantic impact with file location, while Cliewen's pull-request workflow could be mistaken for authority over an adopter's integration policy.

## Decision outcome

**Cliewen recommends simple work when the accepted contract stays unchanged and full work when it changes; the user chooses the route and the repository chooses integration.** Before editing, the agent reads the smallest relevant context and states the recommendation, reason, and discovery that would change it, then reassesses on semantic discovery and the complete diff. Paths and diff size warn but do not decide.

Simple work includes observational analysis with a named consumer, defect correction restoring an unchanged criterion, regression evidence for an unchanged criterion, in-contract configuration, refactoring, maintenance, and editorial work; it carries no CH or full-loop artifacts. Full work adds, revises, or retires criteria, introduces behavior outside accepted criteria, changes capability, policy, plan promise, decision, or methodology meaning, or makes or rejects a consequential decision; uncertainty recommends full.

If simple work grows into full work, the agent pauses and recommends the full loop. If the user declines, the agent keeps the tree truthful and records one integration authorization in the commit with:

```text
Cliewen-Route: simple
Cliewen-Recommendation: full
Cliewen-Override: user chose simple; <concise semantic or evidence risk>
```

The trailers are Git history, not permanent corpus meaning. A route never authorizes a push; repository policy and explicit user authority govern integration, and this repository may require a human-merged pull request. A release is not an adopter route; this repository's version cut is a local simple-work specialization.

This replaces PDR-002's light tier, PDR-011's narrow plain boundary, and PDR-018's behavior-correction rule, and scopes PDR-007 and PDR-040 to a full loop the user chose without changing the human acceptance boundary.

## Rejected: route by diff size, paths, or a release process

Those heuristics confuse mechanical surface with accepted meaning, while release and integration policy belong to the repository owner rather than generated Cliewen methodology.

## Carrier

The core and system architecture, C-002/C-005/C-012/C-013, CAP-005/CAP-006 and AC-139, the amended decision records, routing hubs and scaffolded corpus, canonical and generated change-routing/review/local-convention/delta/verify skills, templates, contributor and public guides, CI scope detection, and their tests carry this recommendation and override boundary.
