---
id: IDR-002
type: decision
status: verified
links: [ADR-013, ADR-022, ADR-043, ADR-046]
title: The managed skill roster is guarded locally and reported to adopters
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# IDR-002 — The managed skill roster is guarded locally and reported to adopters

## Context

Cliewen controls its generated skill set but does not own an adopter's routing hub, where omitting a managed skill may be deliberate.

## Decision

A repository test requires every managed skill name in both Cliewen routing-hub sources. In an adopter corpus, a missing hub row is counted and reported without becoming a validation issue or triggering a rewrite.
