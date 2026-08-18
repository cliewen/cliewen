---
id: PDR-027
type: decision
status: verified
links: [PDR-007, PDR-017, PDR-019, C-012, P-012]
title: Branch protection enforces admission, not acceptance evidence
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-027 — Branch protection enforces admission, not acceptance evidence

## Context and problem statement

Required CI and branch protection can enforce admission to a merge boundary, but acceptance evidence must explain why a criterion is proven and cannot be delegated to mutable forge configuration.

## Decision outcome

**Branch protection enforces admission to the human merge boundary; it is never Cliewen acceptance evidence.** A supported executable reference or the Human-class acceptance-brief entry remains the evidence endpoint, while the protected PR and required checks enforce the conditions for human merge.

## Rejected: treat a green required check as proof

A green check reports only that a configured workflow succeeded; it does not express the criterion, scenario, or evidence relationship a reviewer must judge.

## Carrier

PDR-007 and C-012 distinguish protected integration from review and human judgment; PDR-017, PDR-019, and the public methodology guide carry the acceptance brief and classified-evidence boundary.
