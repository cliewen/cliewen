---
id: ADR-001
type: decision
status: verified
links: [G-001]
title: "Implementation language for the clue CLI: Go"
author: human
accepted-by: Flemming N. Larsen (2026-07-12)
---

# ADR-001 — Implementation language for the `clue` CLI: Go

## Context and problem statement

The deterministic `clue` judge needs a maintainable implementation language that produces a dependency-free single binary and is reliable for agent-authored work.

## Decision outcome

**Use Go.** Static binaries, straightforward cross-compilation, a standard-library-first ecosystem, and mature CLI and Markdown libraries satisfy distribution and maintenance. Go also has the strongest agent fluency and the shortest learning path for the maintainer.

Rust provides a single binary and strong CLI libraries, but its iteration and fluency costs are unnecessary for this parser and graph-walking problem. Kotlin Native/JVM fits maintainer fluency but makes single-binary distribution harder and has thinner agent fluency. Neither satisfies the primary distribution driver as well as Go.
