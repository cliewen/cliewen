---
id: PDR-018
type: decision
status: verified
links: [G-001, P-007, M-029, AN-008, AN-010, PDR-002, PDR-011, PDR-042]
title: Behavior changes remain full until adopter evidence supports a narrower loop
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-018 — Behavior changes remain full until adopter evidence supports a narrower loop

> **Superseded by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** a defect correction restoring an unchanged accepted criterion may be simple; changing the accepted contract remains the full-loop recommendation.

## Context and problem statement

An existing criterion states intended behavior but does not prove that an implementation change leaves the criterion, evidence boundary, and product reality aligned. Before adopter evidence existed, widening the light tier for product behavior would have removed the proposal and digest on an untested assumption.

## Decision outcome

**The original rule kept product behavior changes full, including behavior under an existing criterion, until adopter history could measure their cost and failure modes.** Pure behavior-preserving refactors remained light. PDR-042 now owns the route: simple includes a defect correction that restores an unchanged criterion, while behavior or contract changes remain full.

The canonical change-tier source, generated lifecycle skills, routing hubs, and public change-loop guidance carry the distinction. M-029's adopter evidence door remains the evidence-based way to reopen a route boundary rather than using file size or Cliewen's methodology-heavy history as a proxy.
