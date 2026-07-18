# Decisions

Why things are the way they are — every decision made during a change is recorded here in its digest, routed by how expensive it is to reverse.

- **ADR-xxx** — Architectural Decision Records: decisions about the structure of the software (or of this corpus). MADR format: context, options, outcome, consequences.
- **PDR-xxx** — Project Decision Records: expensive-to-reverse decisions about how the project works (process, review rules, release policy). Same format as ADRs.
- **[log.md](log.md)** — the decision log: cheap-and-local-to-reverse decisions as one-line dated rows (`Date | Decision | Why | Change/PR`). The litmus test: if reversing it later is a cheap, local edit, it is a log row, not a record.

Decision records carry provenance in their status: a record an agent writes during a change starts `inferred`; explicit human approval (review approval or a stated "approved") promotes it to `verified`, with each approver signed in `accepted-by:`. Merging a change makes its decisions binding either way — approval changes their standing, not their force. Rejected records stay in the corpus as history.

A decision that adopts a well-established practice cites it by name and records only the local why.

Decision records are timeless: state what is decided and only the enduring context and rationale needed to understand it. Triggering incidents, chronology, conversations, implementation details, and review history belong in the change workspace, PR, and Git history; include a historical fact only when removing it would make the decision unintelligible.

<!-- clue:index:start -->
- [log.md](log.md) — the decision log: cheap-to-reverse decisions as dated rows
<!-- clue:index:end -->
