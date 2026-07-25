---
id: ADR-029
type: decision
status: inferred
links: [P-006, PDR-004, C-009]
title: accepted-by records only approval given under Cliewen's merge boundary
author: agent
accepted-by: []
---

# ADR-029 — `accepted-by` records only approval given under Cliewen's merge boundary

## Context and problem statement

Brownfield extraction converts decision records from formats that carry their own acceptance — MADR's `decision-makers`, `consulted`, and `informed` frontmatter name humans who approved a decision years before the corpus existed. That acceptance is real, but [PDR-004](PDR-004-merge-binds-approval-signs.md) defines `accepted-by:` as the record of an explicit human approval given under Cliewen's own merge boundary — a PR review approval, a review comment, or a stated "approved" — dated by the first such signature. PDR-004 never faced a record with acceptance history that predates the corpus, so nothing settled whether converting `decision-makers` into `accepted-by:` is a permitted shortcut or a second, incompatible meaning for the same field. Left unstated, an extraction could populate `accepted-by:` from source metadata alone, promoting a record to `verified` without any human ever approving it under Cliewen — the exact silent inflation PDR-004 was written to prevent — and [C-009](../constraints/C-009-type-specific-frontmatter.md)'s promotion trigger, a `checkCoreFields` required-fields map covering `accepted-by`, can only be built machine-checkable once the field carries one meaning everywhere.

## Decision outcome

**`accepted-by:` records only approval given under Cliewen's merge boundary and nothing else. Acceptance that predates the corpus — a MADR record's `decision-makers`, `consulted`, and `informed`, or any other pre-Cliewen acceptance — is preserved as body prose carrying the original names, roles, and dates, while the record's `accepted-by:` stays `[]`, the same empty list every unsigned record already carries.**

- **One field, one meaning.** `accepted-by:` names a human who reviewed and approved the decision through a Cliewen PR review approval, a review comment, or a stated "approved" — never a name recovered from a converted source's own acceptance metadata. A record with real, dated, pre-Cliewen acceptance and a record with none look identical in `accepted-by:` (both `[]`) and are distinguished only by what the body prose says, which is exactly the shape an unsigned record already has.
- **This extends PDR-004; it does not replace it.** PDR-004 states how a signature is made and dated once it happens. This decision adds one boundary condition on top: a populated `accepted-by:` always names an approval given under the merge boundary PDR-004 describes, so a converted record's real acceptance history can never be mistaken for that approval. Reading a populated `accepted-by:` on an `inferred` record any other way would contradict PDR-004 rather than extend it.
- **Extracted acceptance is not discarded, only relocated.** The names, roles, and original acceptance dates a source record carried are preserved as body prose — real history, not deleted, but carried where it cannot be read as a Cliewen approval.
- **Every converted decision is still born `inferred`, `author: agent`**, exactly like any other agent-authored decision; the boundary this record states changes nothing about that starting status.

**Carrier:** the `accepted-by:` clause in the shared skill fragment `internal/skills/source/shared/decision-records.md.tmpl`, rendered into `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, and `clue-verify`; the restated rule in `docs/decisions/README.md`; and `internal/scaffold/templates/docs/decisions/README.md`, which ships the same statement into every repository `clue init` touches. No machine carrier yet — [C-009](../constraints/C-009-type-specific-frontmatter.md) names the promotion trigger that would let `checkCoreFields` enforce it.

### Rejected: convert decision-makers directly into accepted-by

Writing a MADR record's `decision-makers` straight into `accepted-by:` would promote a converted record to `verified`-looking status without any human ever approving it under Cliewen's own merge boundary — the field would stop meaning what PDR-004 defines it to mean, and the corpus would carry two kinds of acceptance indistinguishable by the field alone. This also forecloses the field ever becoming machine-checkable, since `checkCoreFields` cannot enforce "approved under Cliewen" against a field that sometimes means something else.

### Rejected: a separate frontmatter field for pre-Cliewen acceptance

A second field (e.g. `prior-acceptance:`) would carry the same information the body prose already can, at the cost of a new required-fields shape for exactly one extraction case. Body prose already has the room to state names, roles, and dates in full sentences, and nothing downstream needs the pre-Cliewen acceptance in a structured field — only `accepted-by:` is ever checked for promotion.
