---
id: PDR-049
type: decision
status: verified
links: [ADR-046, C-018, C-019]
title: Durable communication is readable and nonderived
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# PDR-049 — Durable communication is readable and nonderived

## Context

A computed figure copied into durable prose becomes stale, while a bare identity in a human-facing report makes the reader open another artifact merely to learn whether it matters.

## Decision

A durable record names the command that computes a changing figure instead of restating the result, except when a measured result is the record's point and names its method and time. Human-facing reports, statuses, explanations, and handoffs describe what a referenced identity is about rather than leaving it bare.
