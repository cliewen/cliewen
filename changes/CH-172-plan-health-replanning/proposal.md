---
id: CH-172
type: change
status: proposed
links: []
title: Plan-health checks pause invalid work without treating every replan as a decision
---

# CH-172 — Plan-health checks pause invalid work without treating every replan as a decision

## Why

The current plan workflow says that every substantive plan revision needs a decision record. That contradicts the rule that only future-shaping choices earn one and turns ordinary delivery replanning into ceremony. It also gives an agent no explicit point to check whether a milestone's plan still holds before it continues work.

## What

This plan-less full change adds an evidence-triggered plan-health check at milestone start or resumption and whenever new evidence challenges the remaining campaign. A failed check records the mismatch in the change workspace and pauses affected work for human direction. A passing check leaves no durable record. A selected replan is declared and human-directed; it creates a typed decision record only when the selected course is future-shaping.

The change updates the canonical generated-skill sources, their generated adopter outputs, the methodology decision and live documentation, and focused evidence.
