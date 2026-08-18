---
id: PDR-036
type: decision
status: verified
links: [P-013, PDR-012, PDR-035, PDR-019, C-004, C-012, C-017]
title: The review loop has a maximum pass count, and reaching it is a human decision
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# PDR-036 — A maximum pass count, then the human decides

## Context and problem statement

PDR-035's three-pass budget could be exceeded whenever each pass found a blocking defect, so it limited discretionary review but not the loop's actual cost.

## Decision outcome

**The review loop has a maximum pass count, and reaching it hands the decision to the human.** It stops on a pass with no blocking findings, runs another pass only after a blocking pass, and never runs merely to reach the number. The default maximum is five, stated in C-017 and the generated review skill; an adopter may set another in its repository conventions.

When the maximum is reached with blocking findings outstanding, the loop stops, reports the findings, and asks whether to run further passes. That answer is the only authority that extends the budget, and reaching the maximum never permits publication with an unresolved blocking finding.

## Rejected: let blocking findings exceed the maximum automatically

That is PDR-035's unbounded shape: a limit that never binds when review keeps failing is not a meaningful bound.

## Rejected: stop at the maximum without asking

A human question can finish a nearly converged repair without turning a hard ceiling into permission to publish an unresolved defect.

## Carrier

C-017 and the generated review skill are the only carriers of the number besides this record; PDR-012, PDR-035, contributor and guide text, templates, and skill sources carry the rule without duplicating the digit.
