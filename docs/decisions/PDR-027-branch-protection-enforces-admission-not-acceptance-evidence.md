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

Required CI and branch protection can prevent an unvalidated pull request from merging. That is enforcement of an admission condition. Acceptance evidence answers a different question: why a particular acceptance criterion is proven. Treating hosted protection as that proof would make mutable forge configuration stand in for the reviewed criterion and its evidence.

## Decision outcome

**Branch protection enforces admission to the human merge boundary; it is never Cliewen acceptance evidence.** A supported executable reference or the Human-class acceptance-brief entry remains the evidence endpoint. The protected PR makes the human merge and required checks enforceable, but it does not prove that a test or brief matches a criterion.

**Carrier:** PDR-007 and C-012 distinguish protected integration from review and human judgment; the public methodology guide states the same division. PDR-017 and PDR-019 retain the acceptance brief and classified evidence as the evidence carriers.

### Rejected: treat a green required check as proof

A green check establishes only that the configured workflow reported success. It cannot express the criterion, scenario, or evidence relationship a reviewer must judge, and its configuration is forge state rather than corpus meaning.
