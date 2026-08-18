---
id: PDR-037
type: decision
status: verified
links: [P-013, PDR-013, PDR-029, ARCH-003, AN-022, C-013]
title: Tooling is judged by whether removing it hands work back to a human
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-037 — Tooling is judged by what it holds, not by what the core names

## Context and problem statement

PDR-013's core test under-determines commands, checks, fields, and artifact types because the core names guarantees while tooling holds the mechanisms that keep people from forgetting their obligations.

## Decision outcome

**For tooling, the test is whether removing it moves an obligation from a machine back to a human.** A yes keeps the item; a no makes it a removal candidate. PDR-013's *does the core need it?* remains decisive for anything whose existence changes meaning, and neither test authorizes removing a check that changes what the thread, merge boundary, or `clue validate` asserts. PDR-029's shared-memory rule also remains: reducing what the corpus remembers requires its own decision.

The test explains why a campaign may find that almost nothing should be removed; that is an honest result, not a requirement to manufacture deletions.

## Rejected: apply the core test to tooling or score the surface by removals

The core test would discard mechanisms that hold obligations, while a removal count would reward changes that make the repository smaller without making it safer or clearer.

## Carrier

The amendment note at PDR-029, P-013's M-064 and its measure, and this record carry the tooling test. No skill, hub, guide, or CLI carrier states the campaign test itself.
