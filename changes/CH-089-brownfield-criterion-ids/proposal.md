---
id: CH-089
type: change
status: open
links: [P-009, M-038, ADR-009, CAP-002, CAP-003]
title: Preserve brownfield criterion identities
---

# CH-089 — Preserve brownfield criterion identities

## What

This full change serves P-009 milestone M-038 by extending Cliewen's acceptance-criterion identity contract beyond a single prefix segment and digits. It supports stable source IDs such as `SNAP-SQS-001` and `ADP-045b` through declaration, namespace validation, normalization, uniqueness, retirement, coverage, and the supported Go, JVM, and Cucumber evidence carriers.

The change adds an ADR, focused acceptance criteria and regression fixtures, updates the validator and coverage paths, and aligns the extraction contract, generated skills, implementation explanations, public guidance, and `[Unreleased]` release note. Existing `<PREFIX>-<digits>` corpora continue to validate unchanged, while extraction preserves source IDs verbatim and mints only missing IDs deterministically.

## Why

Brownfield repositories can already have multi-segment or letter-suffixed criterion identities whose meaning and test references depend on their exact spelling. Rejecting those identities forces renumbering during migration, severing provenance and creating avoidable evidence churn. The current grammar is therefore a migration blocker identified by P-009 and cannot be solved as a documentation-only change.

## Decision boundary

The change makes criterion identity syntax and its local evidence handling more expressive without changing what an acceptance criterion, classified evidence, extraction, or the human merge boundary means. It does not add a general configuration interface, foreign evidence resolution, distributed-work authorization, or a new source-format mapping.
