---
id: PDR-011
type: decision
status: verified
links: [PDR-002, PDR-007, PDR-042, C-002, C-005, C-012, AN-006]
title: Plain changes stay outside Cliewen while retaining human merge
author: agent
accepted-by: Flemming N. Larsen (2026-07-20, planning conversation)
---

# PDR-011 — Plain changes stay outside Cliewen

> **Superseded by [PDR-042](PDR-042-routing-recommends-contract-aware-effort.md):** simple work now includes the accepted-contract-preserving cases named there, and route selection does not itself authorize integration.

## Context and problem statement

The light tier still charged purely editorial work with Cliewen identity, proposal metadata, corpus loading, and unrelated verification even when the change produced no relevant evidence. The methodology needed a boundary for edits that carry no product, contract, decision, plan, policy, or methodology meaning.

## Decision outcome

**The original plain route kept meaning-free editorial changes outside Cliewen while retaining an ordinary human-reviewed branch and pull request.** Plain work carried no CH identity, plan declaration, proposal, corpus read, Cliewen verification, plan bookkeeping, or mandated changelog entry. Protected surfaces and uncertainty failed closed, and paths or diff size could not decide meaning.

PDR-042 replaces the narrow plain/light hierarchy with the simple/full recommendation and defines current simple work. The surviving boundary is that Cliewen does not own every repository edit, while integration authority remains with the user and repository policy; agents still do not push to `main` or merge their own full changes. The routing hubs, canonical tier and boundary skills, contributor guidance, and focused CI selection carry the rule.
