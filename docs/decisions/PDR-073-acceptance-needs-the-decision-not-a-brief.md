---
id: PDR-073
type: decision
status: inferred
links: [G-016, CAP-012, CAP-006, ARCH-003, C-012, C-013, C-023, PDR-017, PDR-039, PDR-067, PDR-072, ADR-033]
title: Acceptance needs the human's decision, not an acceptance brief
author: agent
accepted-by: []
binds: adopter
---

# PDR-073 — Acceptance needs the decision, not a brief

## Context

A tracked change could only be accepted with an acceptance brief: a document of five required sections that the agent wrote after the work, which `clue accept` parsed, hashed into the approval record and copied into the merge, and which the validation workflow required in a ready pull request's body. In use the brief repeated what the proposal, the commits, the test results and the conversation already held, cost a long writing step at the end of every change, and made reading it a condition of a decision the human had often already taken. The one thing it carried that nothing else did was the record of a person's observation for a `Test-type: Human` criterion, including what was not observed.

## Decision

**There is no acceptance brief. The human's decision on an exact candidate is what accepts a tracked change, in either acceptance mode.**

- `clue accept <candidate> --base <base>` takes no brief. It finds the change from the one proposal the candidate's history introduces after the base, refusing none or several, and writes the merge message from the change's identity, the proposal's title, candidate and base, the proposal's human evidence and any approval record. `--brief` is refused with a message saying the brief is no longer used.
- An approval record binds candidate, base and decision. A record still carrying `brief-sha256` is refused with the same plain explanation.
- A pull request's body is free; no shipped workflow reads it.
- Human-class proof moves to a `## Human evidence` section of the change's proposal, written by the agent before the digest only from what the human reported: observer, conditions, observation, result and what was not observed. The proposal stays in history and `clue accept` copies the section into the merge. An authorized dependency and the vision a change proceeds under are stated in the proposal too.

Every check `clue accept` made that did not depend on the brief remains: exact candidate and base, a clean integration checkout, the candidate incorporating the base, the same policy on both, a valid digested corpus, safe snapshots, the exact candidate tree with ordered parents, and never pushing. The agent still never accepts its own change.

This supersedes the brief parts of [PDR-017](PDR-017-merge-gate-has-content.md), [PDR-039](PDR-039-dependent-changes-carry-authorization.md), [PDR-067](PDR-067-local-human-acceptance.md), [PDR-072](PDR-072-recorded-human-approval-allows-delegated-local-execution.md) and [ADR-033](ADR-033-human-proof-and-draft-criteria.md). The verifiable thread keeps genuine Human proof with a new carrier, and the human acceptance boundary keeps the decision and the agent's exclusion; only the document the decision was bound to goes ([ARCH-003](../architecture/core.md)).

## Alternative considered

A short, optional brief was rejected: an optional document becomes a habit agents keep producing, and the obligations it carried that matter are either checked by `clue accept` or kept in the proposal.

**Carrier:** `internal/accept`, `cmd/clue/accept.go`, `internal/skills/source/shared/review-boundary.md.tmpl` and the `clue-delta`, `clue-verify`, `clue-extract` and `clue-upgrade` sources under `internal/skills/source/skills/`, `internal/scaffold/templates/acceptance/approval.yaml` and `internal/scaffold/templates/github/PULL_REQUEST_TEMPLATE/tracked.md`, the reusable validation workflow, and CAP-012 and CAP-006. The [design overview](../design/README.md) describes the acceptance flow.
