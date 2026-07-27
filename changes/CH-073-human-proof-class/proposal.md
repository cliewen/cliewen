---
id: CH-073
type: change
status: open
links: [P-007, M-026, ADR-032, ADR-007, CAP-002, CAP-006]
title: Give human-verified criteria a machine-visible home
---

# CH-073 — Give human-verified criteria a machine-visible home

M-026 closes the gap M-025 deliberately left open: a criterion that is deliberately verified by a human, never by a code test, currently has no honest way to validate. Today it either forces its whole capability to stay `draft` or borrows the tombstone channel it does not mean. The change gives the proof-class vocabulary ADR-032 established a `human` class, satisfied by a named item in the acceptance brief rather than executable evidence, and generalizes ADR-007's `@retired` tombstone convention with a per-criterion `@draft` token so one criteria file can honestly hold both proven and unproven criteria without the file or its capability leaving `active`.

The change will record the extension in its own decision record — extending ADR-032's vocabulary rather than riding inside it, per M-026's exit criterion — teach `checkACTests` the `human` class and the `@draft` tag-line token, extend the acceptance-brief template and `clue-verify` so a newly or materially declared `Test-type: Human` criterion is named and confirmed at merge, and add a derived coverage report in the book's three states (`covered`, `partial`, `gap`) computed from the corpus rather than committed as a registry file.

## Scope

- Extend ADR-032's proof-class vocabulary with a `human` class in a new decision record.
- Parse and enforce a per-criterion `@draft` tag-line token in `checkACTests`, exempting that criterion from the active-file test requirement without retiring it or drafting its capability.
- Extend the acceptance-brief template and its digest so a criterion newly or materially declaring `Test-type: Human` is named and the human confirms its verification at merge.
- Add a derived coverage report (`covered` / `partial` / `gap` per capability) computed from corpus state, not a committed file.
- Update the corpus, capability contracts, source skills, generated skills, and user-facing changelog where they describe the former limit.

## Out of scope

- Retiring corpus artifacts, bounded provenance, and task-proportionate context remain later P-007 milestones (M-027…M-029).
- The evidence-classification vocabulary and harvesters M-025 already shipped are not reopened except to add the `human` class.
- Any change to how `draft` behaves for a whole capability or criteria file outside the new per-criterion token.
