---
id: ADR-016
type: decision
status: verified
links: [ADR-010]
title: ADRs are for expensive-to-reverse decisions; the rest is a decision log
author: agent
accepted-by: Flemming N. Larsen (2026-07-14, PR #9)
---

# ADR-016 — ADRs for the expensive-to-reverse; a decision log for the rest

## Context and problem statement

When every recorded decision is an ADR, the decisions folder fills with library picks and tunable defaults alongside the choices that actually shape the system, and the reader can no longer tell which documents constrain future changes. A record that costs a full MADR document also discourages recording small decisions at all. Where is the line, and what is the cheap record on the other side of it?

## Decision outcome

**ADRs are reserved for decisions that are expensive to reverse: architecture, methodology mechanics, public contracts, AC semantics. Everything else is a row in [`docs/decisions/log.md`](log.md).**

- **Litmus test:** if reversing the decision later is cheap and local, it is a log row; if it constrains future changes, it is an ADR.
- **Log format:** one table, columns `Date | Decision | Why | Change/PR` — the same rows-in-a-table pattern plans use for milestones. The Why is one line; a why that will not fit one line is a sign the decision belongs in an ADR.
- **The log is a corpus artifact** (`type: log`, status `active`), linted like any other: `clue validate` enforces its frontmatter and status vocabulary.
- **Retention applies to decisions, not to file format.** The rule that decisions are never deleted stands; a demoted ADR's decision survives as a dated log row, and git history remains the provenance archive for the full text — the same philosophy that lets digests delete `/changes/`. Rejected *ADRs* still live as ADR files: a recorded rejection of an expensive-to-reverse option is itself expensive-to-reverse knowledge.
- **Demotion and promotion are ordinary change content:** an ADR failing the litmus test is demoted (file deleted, row added, inbound references repointed) in a change; a log row that turns out to constrain future changes is promoted to a full ADR citing the row's date.

**Carrier:** the ADR-vs-log routing in the `clue-delta` skill's digest step (agent); the `log` artifact type in `clue validate`'s status vocabulary (machine); the decisions folder README prose, which `clue init` will scaffold (default).

### Rejected: keeping everything as ADRs

Uniformity reads as rigor but costs signal: the folder stops answering "what must I not casually undo". The expensive format also suppresses recording cheap decisions, which is how undocumented lore forms.

### Rejected: a log outside the corpus (wiki, issue labels, PR descriptions)

Decisions scattered across venues are decisions lost; PR descriptions are already the light tier's proposal ([ADR-015](ADR-015-light-change-tier.md)) and do not aggregate. One linted file in the decisions folder keeps the register greppable, versioned, and inside the system-of-record.
