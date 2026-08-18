---
id: ADR-037
type: decision
status: verified
links: [P-009, M-038, ADR-009, CAP-002, CAP-003]
title: Brownfield criterion IDs preserve segmented prefixes and letter suffixes
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-037 — Brownfield AC ID grammar

## Context and problem statement

Brownfield criteria need to preserve existing identities, including segmented prefixes and lowercase suffixes, while evidence carriers need syntax that language and tag rules can represent without making declarations ambiguous.

## Decision outcome

**The canonical acceptance-criterion ID is `<PREFIX>-<NUMBER><SUFFIX?>`: uppercase alphanumeric prefix segments joined by single hyphens, decimal digits, and an optional lowercase suffix.** IDs are case-sensitive, declarations and links use the canonical spelling, `ac-prefix` uses the same grammar with default `AC`, full IDs remain unique, and wrong-namespace declarations fail.

Evidence carriers normalize syntax only: JVM/Cucumber tags may replace prefix hyphens with underscores, while Go and JVM names remove them; numeric components and lowercase suffixes remain unchanged. Hyphen-stripped prefix collisions are rejected. Malformed-identity diagnostics apply only inside declared namespaces, leaving unrelated runner tags and prose untouched.

Extraction preserves every source ID. Requirements without one receive the next numeric component after the namespace maximum, ignoring suffixes for that maximum; minted IDs use the canonical prefix and number and are recorded in the extraction report.

**Carrier:** the canonical parser and namespace checks, evidence harvesters, context/coverage consumers, and the preservation and minting rules in `clue-extract`.
