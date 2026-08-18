---
id: ADR-026
type: decision
status: verified
links: [PDR-013, ADR-025, C-004]
title: Unknown artifact types are adopter extensions, validated against the default lifecycle
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-026 — Adopter-defined types validate against the default lifecycle

## Context and problem statement

Adopters may add artifact types, but a closed validator registry rejected deliberate extensions such as `risk`, `experiment`, or `runbook`.

## Decision outcome

**An unrecognized type is an adopter extension and is validated against ADR-025's default lifecycle.** `statusVocabFor` returns the default for types without an exception and `checkStatusVocab` does not emit `unknown type`.

All other guarantees remain: core fields, unique IDs, resolving links, and valid lifecycle values are still required. This re-scopes the check rather than weakening it: it no longer asserts a closed set of type names, so a misspelled built-in type may instead fail when its status uses the wrong vocabulary. This is the deliberate boundary between core form and adopter domain. The carrier is `statusVocabFor`, the extension statement in ARCH-003, and the init red-line rule. A configurable type registry and warning-only rejection are rejected because they preserve the closed-world burden or train adopters to ignore failures.
