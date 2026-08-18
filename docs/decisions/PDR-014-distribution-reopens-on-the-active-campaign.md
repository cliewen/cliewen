---
id: PDR-014
type: decision
status: verified
links: [P-007, CAP-001, C-015, ADR-030, PDR-008]
title: Installation distribution reopens on the active campaign, not a successor plan
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-014 — Installation distribution reopens on the active campaign

## Context and problem statement

Installation distribution had been parked until its onboarding evidence and scope existed, but the corpus did not say who could reopen it or where the work would land once that condition was met. The evidence now includes the checkable onboarding bar, a released binary on six platforms, and a written install procedure whose six manual steps expose the gap.

## Decision outcome

**Installation distribution joins the active P-007 campaign as M-030 and M-031 rather than waiting for a successor plan.** M-030 creates the install commands and M-031 builds the marketplace bootstrap from them, so the milestones remain ordered. CAP-001, C-015, CAP-004, the guide, and P-007 carry the revision under PDR-008.

A successor plan, a duplicate goal, or a log row alone would either delay work whose shape and evidence are already known or fail to carry the deliberate plan-revision boundary and rejected alternatives.
