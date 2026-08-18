---
id: PDR-017
type: decision
status: verified
links: [G-001, AN-008, P-007, PDR-012, PDR-013, PDR-039, C-012, C-013]
title: The merge gate carries an acceptance brief, and the review loop adds an advisory scenario-resolution step
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-017 — The merge gate has content

> **Extended by [PDR-033](PDR-033-planning-and-implementation-are-separate-steps.md):** an optional pre-implementation pause records where work may split or change hands and asks whether implementation begins.

> **Extended by [PDR-039](PDR-039-dependent-changes-carry-authorization.md):** a dependent change's acceptance brief names its unmerged base, authorization, and the meaning its merge would bind.

## Context and problem statement

The validator and agent review loop cannot replace the human judgment at merge. The boundary needed durable content telling the human what remains to verify, while semantic alignment checks had to remain advisory rather than pretending to be deterministic proof.

## Decision outcome

**The digest emits an acceptance brief at the top of the PR body, the review loop adds a non-blocking per-criterion scenario-resolution verdict, and Propose may opt into a pre-implementation pause.** The brief asks whether the plan item remains wanted, lists every added or changed criterion with its scenarios, names inferred decisions and invalidated or superseded records, and includes the competence warning and one-screen pressure.

For executable evidence the review records `verifies`, `verifies-something-adjacent`, or `undetermined` after comparing setup, action, and assertions. A genuine Human criterion is proved by its criterion-and-scenario entry in the acceptance brief, not by an invented code test. The verdict is informational; actual defects follow the normal finding path. The pause remains opt-in, and work or handoff authority lives in corpus artifacts, the change workspace, or the PR rather than private memory.

`clue-delta`, `clue-verify`, their durable-work fragment, the PR template, the reusable workflow and caller, CAP-006, and the public guide carry this merge content.
