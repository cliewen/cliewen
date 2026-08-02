---
id: PDR-022
type: decision
status: inferred
links: [G-001, CAP-001, ADR-018, ADR-031, PDR-019, C-006, C-013]
title: A scaffolded vendor entry point may exist, and it may only point at AGENTS.md
author: agent
accepted-by: []
---

# PDR-022 — A scaffolded vendor entry point may exist, and it may only point at AGENTS.md

## Context and problem statement

Cliewen is deliberately cross-agent: `AGENTS.md` is the single routing hub and the skills are written for any assistant that can read them. Assistants, however, do not agree on which file they load. Claude Code documents that it reads `CLAUDE.md` and not `AGENTS.md`, and offers an import as the supported bridge; other tools name other files.

`clue init` therefore materializes a hub that a widely used assistant never opens, while ADR-018 already mirrors the skills into that same assistant's directory. The result is worse than either consistent choice: the adopter's agent lists the lifecycle skills and never receives the routing that says when to invoke them, so it can execute the method precisely while classifying the change wrongly. This repository ran with that gap and the errors it produced are recorded in its own history.

May a cross-agent methodology emit a vendor's flagship instruction file, and if so, what may that file contain?

## Decision outcome

**The scaffold may emit a vendor entry point, and the file may contain nothing but a pointer to `AGENTS.md` and the explanation of why it is a pointer.**

Three conditions bound it.

**No rule ever lives in a vendor file.** A rule written there is invisible to every other assistant, which recreates the split the single hub exists to prevent — and it would be invisible in the direction that matters least visibly, since the agent that reads it behaves correctly while the others silently do not. Adding a rule to a vendor entry point is a methodology-carrier change and is governed by PDR-019 like any other.

**The file is adopter-owned on arrival.** `clue init` never overwrites it and `clue migrate` never rewrites it. An adopter's own assistant-specific instructions sit below the import and are safe from every tool this project ships, which is what makes emitting the file honest rather than a claim on their repository.

**A vendor qualifies on evidence, not popularity.** An entry point is emitted for an assistant whose published behavior makes the hub unreachable without one, and Claude Code is that case today. This is the same evidence bar ADR-018 met for the skills mirror; it is not a standing invitation to accumulate one file per tool. Adding the next one is a decision, and the cost of the wrong answer is a repository root full of files nobody reads.

For an adopter who onboarded before this, `clue migrate` reports a missing or unrouted entry point and repairs neither: the missing case is already solved by re-running the non-destructive `clue init`, and a `CLAUDE.md` the adopter wrote themselves is their prose, which no migration edits.

### Rejected: emit nothing and keep the methodology purely neutral

Neutrality that an assistant cannot observe is not neutrality; it is a hub the tool never reads. The cost lands entirely on the adopter, who must discover the vendor's loading rule themselves, and it lands silently, because nothing reports that the routing never arrived. Cliewen already abandoned this position for the skills mirror, and holding it for the hub while conceding it for the manuals is the least defensible of the three options.

### Rejected: let the vendor file carry Claude-specific methodology

It splits the contract across two files with no mechanism keeping them in agreement, and the divergence is undetectable from inside any single assistant. PDR-019 exists because carriers drift; the fix is fewer carriers of the rule, not a sanctioned second one.

### Rejected: make the entry point a managed carrier that migration keeps current

A managed carrier is refreshed from Cliewen's own bytes, which would silently discard the assistant-specific instructions the adopter is explicitly invited to add below the import. A pointer that never changes needs no refresh; the only thing worth reporting is its absence, and a notice does that without touching a file.

**Carrier:** CAP-001's criteria and design carry what `clue init` emits and what `clue migrate` reports; the scaffolded `CLAUDE.md` template carries the pointer contract for adopters; ADR-018 names the emitted set.
