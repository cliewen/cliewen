---
id: PDR-039
type: decision
status: verified
links: [P-013, AN-013, PDR-007, PDR-017, C-012, C-013, CAP-006]
title: An authorized dependent change records its base and what its merge would bind
author: agent
accepted-by: Flemming N. Larsen (2026-09-02, conversation)
---

# PDR-039 — An authorized dependent change records its base and what its merge would bind

## Context and problem statement

PDR-007 permits work based on an unmerged change only after human authorization, but a branch or green validation cannot show which unaccepted meaning it relied on or what its merge would bind.

## Decision outcome

**An authorized dependent change keeps the answered dependency in its committed workspace until digest.** It names the unmerged base branch or commit, the human authorization, and the specific unaccepted meaning at risk; the ready pull request repeats the same dependency and binding in its acceptance brief. This discloses the issue for human merge judgment without accepting or merging either change.

Once the base is accepted and the dependent branch incorporates accepted `main`, the dependency disappears and the workspace is digested. Git history preserves the authorization; no validator infers it from Git or forge state.

## Rejected: keep a permanent registry or infer the relation

The relation is transient and Git history already preserves it after acceptance, while graph or forge state cannot supply the human authorization or the meaning at risk.

## Carrier

PDR-007's stacking clause is amended; the review-boundary source and generated skills, C-012, the public change-loop guide, pull-request templates, and CAP-006 carry the disclosure boundary.
