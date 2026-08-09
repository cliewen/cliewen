---
id: C-021
type: constraint
status: active
links: [PDR-032, P-013]
title: A suggestion raised mid-change is triaged, never held in memory
source: PDR-032, the shared durable-work fragment
enforcement: human
---

# C-021 — A mid-change suggestion is triaged, never held in memory

A suggestion raised during a change is triaged immediately into one of two carriers. It belongs to the current work → it becomes a task in `tasks.md`, handled before merge. It does not → it becomes a goal with `status: proposed`, written in the digest so it survives the workspace's deletion. Neither carrier is optional, and "I will remember" is not a third: a suggestion that is neither actioned nor recorded has been declined without anyone deciding to decline it. The triage — which carrier, and why — is stated to the human when it is made.

**Residual:** all of it. A suggestion arriving in conversation is invisible to `clue validate`, which reads committed files: no machine can tell whether a suggestion was triaged, silently absorbed into the change, or dropped. The cost is a suggestion that was never accepted or declined by anyone, which looks from outside like it was heard.
