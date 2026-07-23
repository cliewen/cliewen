---
id: CH-051
type: change
status: open
links: [P-005, M-016, G-001, C-012]
title: Define the Cliewen core explicitly and open the P-005 simplification campaign
---

# CH-051 — Define the core, open P-005

## What

Create plan P-005 (the simplification campaign) and its founding decision: an explicit statement of what Cliewen's core is, a red-line rule protecting it, and the periphery/extension story for adopters. New artifacts: PDR-013 (core and red line), C-013 (core changes need a decision record), ARCH-003 `docs/architecture/core.md` (the durable core statement), P-005 (the campaign plan with milestones M-016..M-019).

## Why

Cliewen has grown enough surface that simplification needs a criterion. A Linux-kernel-style core declaration gives every future "should this exist?" debate its test (does the core need it?), protects the load-bearing parts behind an explicit red line, and tells adopters what they may freely extend. The campaign's remaining milestones (status-lifecycle unification, folding quality scenarios into constraints, tier-prose rewrite) all depend on this criterion existing first.

## Decision boundary

This change mints PDR-013 and C-013 and creates plan P-005 — all semantic, all requiring human acceptance at merge. It changes no code, no acceptance criteria, no skills, and no existing decisions. The proposal declares plan item M-016 of P-005, which lands in this same change (the plan's founding change carries its own plan — CH-045 precedent).
