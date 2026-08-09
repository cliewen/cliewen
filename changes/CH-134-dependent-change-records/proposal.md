---
id: CH-134
type: change-proposal
links: [P-013]
title: Dependent changes become durable, reviewable records
---

# CH-134 — Dependent changes become durable, reviewable records

## What

Serve P-013/M-065 by designing and adopting the durable record AN-013/F1 requires for work that depends on an unmerged change. The record must identify the unaccepted base, the dependency, and the human authorization, and make the dependent change's acceptance brief state what its merge would bind.

This change also obtains and records the required determinations for AN-013/F2's emitted-wall divergence and AN-013/F3's unqualified external references. It implements either a mechanism that the human directs, or the resulting explicit decline or successor route.

## Why

Today, a dependent change can make unaccepted meaning appear as accepted corpus truth, while the authorization and base exist only in conversation or forge state. The corpus needs a durable, locally reviewable representation without treating forge state as the system of record or weakening the human merge boundary. M-065 also requires F2 and F3 to be determined rather than silently repeated as open findings.

## Boundaries

This change does not treat a pull-request base, a foreign green check, or a forge query as corpus truth. It does not automate stacked merges or permit an agent to merge a change. Any F2 or F3 mechanism must be separately authorized by the human questions in this workspace.
