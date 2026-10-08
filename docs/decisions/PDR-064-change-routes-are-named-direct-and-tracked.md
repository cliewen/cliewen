---
id: PDR-064
type: decision
status: verified
links: [PDR-042, PDR-023, CAP-006, CAP-001]
title: Change routes are named direct and tracked, and override trailers are read in either spelling
author: agent
accepted-by: Flemming N. Larsen (2026-10-08, conversation)
---

# PDR-064 — Change routes are named direct and tracked

## Context and problem statement

[PDR-042](PDR-042-routing-recommends-contract-aware-effort.md) routes work by whether the accepted contract changes, but named the routes `simple` and `full`. Those words describe how big the work is. A newcomer reads a large refactor as not "simple" and a one-line criterion edit as not "full", and reads "full" as the thorough option a careful person should pick by default.

## Decision outcome

**The route that leaves the accepted contract unchanged is `direct`, and it is the default. The route that changes the accepted contract is `tracked`.** `direct` names how the work proceeds: make the change and run the checks relevant to its surfaces. `tracked` names what the work gets: a CH workspace followed from proposal through digest, verification, and human acceptance. What each route requires, and which work belongs to which route, is exactly as PDR-042 states.

Agents say `Recommended route: direct` or `Recommended route: tracked`. A user override is recorded with:

```text
Cliewen-Route: direct
Cliewen-Recommendation: tracked
Cliewen-Override: user chose direct; <concise risk>
```

Every reader of override trailers — the shipped validation workflow and any local CI scope detection — accepts `direct` or the legacy `simple` as the route and `tracked` or the legacy `full` as the recommendation, permanently. Old commits stay in history, adopters may run older skills, and recognising a legacy spelling costs one alternation.

An adopter's routing hub is their own prose and is never rewritten ([PDR-023](PDR-023-tool-notice-and-hub-instruction.md)). `clue migrate` reports a hub that still names the old routes, so the adopter can align it with the skills their agents read.

## Rejected: `planned`

"Plan" is already a Cliewen artifact. A tracked change may be plan-less, and "a plan-less planned change" would mislead the readers the rename is for.

## Rejected: keep the names and explain them

Explaining `simple` as "the default" leaves names that describe size when the rule is about the contract, so every new reader still has to unlearn their first reading.
