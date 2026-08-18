---
id: ADR-013
type: decision
status: verified
links: [P-002, ADR-008, ADR-011, ADR-012]
title: What ships to adopters is generic; AGENTS.md is the repo-local layer
author: agent
accepted-by: Flemming N. Larsen (2026-07-13, PR #8)
---

# ADR-013 — Ships-generic vs repo-local

## Context and problem statement

Skills and folder READMEs ship to every adopter, while each repository has legitimate local conventions. Naming one repository's convention in a shared artifact would turn methodology into a fork.

## Decision outcome

**Shipped artifacts are generic; `AGENTS.md` is the one repo-local layer; local rules extend but never override the methodology.** Skills are self-contained, carry no corpus document IDs, and use generic hooks for local behavior. Folder README prose is generic and this repository's READMEs are the `clue init` sources; index blocks are generated per repository. `AGENTS.md` routes the corpus and carries local conventions.

Contradictions become an open question whose human answer is a decision. Extraction reconciles existing instructions the same way. ADR-008's per-source mapping sections are refined into mapping files under the one `clue-extract` skill. The carrier is the shipped skill and README text, the doc-ID sanity test, AGENTS.md, and init templates. A second config file and per-repository forked skills are rejected because they split or drift the source of truth.
