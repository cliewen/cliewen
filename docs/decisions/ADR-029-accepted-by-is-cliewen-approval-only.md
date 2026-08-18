---
id: ADR-029
type: decision
status: verified
links: [P-006, PDR-004, C-009]
title: accepted-by records only approval given under Cliewen's merge boundary
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-029 — `accepted-by` records only approval given under Cliewen's merge boundary

## Context and problem statement

Brownfield decision records may contain acceptance history from another system, but `accepted-by:` already means explicit approval under Cliewen's merge boundary.

## Decision outcome

**`accepted-by:` records only approval under a Cliewen pull-request review, review comment, or stated approval.** Pre-Cliewen acceptance such as MADR decision-makers, consulted, and informed is preserved in body prose with its names, roles, and dates; `accepted-by:` remains `[]`. Converted decisions are still born `author: agent` and `status: inferred`.

This extends PDR-004's signing mechanics without replacing them: a populated field always names an approval under that boundary. The carrier is the shared decision-record skill fragment, both decision READMEs, and the future `checkCoreFields` promotion trigger in C-009. Converting source decision-makers directly or adding a second frontmatter field is rejected because either gives one field two meanings or adds a structured field no consumer needs.
