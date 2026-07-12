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

The `clue` CLI is the deterministic judge: stateless, no AI, no orchestration — "a CLI so boring it is finished." Which language do we build it in?

## Decision drivers

- **Single-binary distribution** — install must be one download; no runtime dependency (onboarding is CAP-001, <30 minutes to green).
- **Maintainer fluency** — the maintainer comes from a JVM background; the language must be quickly learnable and boring to maintain.
- **Agent fluency** — most code will be written by coding agents; the ecosystem's idioms must be ones agents produce reliably.

## Considered options

1. **Go**
2. **Rust**
3. **Kotlin Native / JVM**

## Decision outcome

**Go.** Static single binaries and trivial cross-compilation satisfy distribution; the standard-library-first culture and mature CLI/markdown ecosystem (cobra, goldmark + frontmatter extensions) are exactly the "boring and finished" target; agent fluency with idiomatic Go CLIs is the strongest of the three; and the jump from a JVM background is the shortest.

### Rejected: Rust

Single binary too, and best-in-class CLI ergonomics (clap, serde, comrak) — but slower iteration, a steeper fluency curve for both maintainer and reviewers, and correctness guarantees the problem (parsing frontmatter, walking a graph) does not need.

### Rejected: Kotlin Native / JVM

Maximum maintainer fluency, but single-binary distribution is the weak point (GraalVM/Kotlin-Native friction) and agent fluency for that toolchain is thinner. Fails the most important driver.
