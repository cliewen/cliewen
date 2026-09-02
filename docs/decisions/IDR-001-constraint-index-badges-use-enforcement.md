---
id: IDR-001
type: decision
status: verified
links: [ADR-041, ADR-045, ADR-046, C-016, CAP-005]
title: Constraint index badges show enforcement
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# IDR-001 — Constraint index badges show enforcement

## Context

An index badge should expose the fact a reader uses to judge the record. A constraint's lifecycle status does not say who or what holds its rule.

## Decision

Generated constraint index rows use the constraint's `enforcement` value as their badge; every other artifact uses `status`. `clue validate --index-rows` reports a direct constraint row whose badge disagrees without rewriting it or failing ordinary validation.
