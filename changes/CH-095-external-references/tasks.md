---
id: CH-095-tasks
type: tasks
status: open
links: [CH-095]
title: Tasks for CH-095
---

# Tasks

- [x] Record the ADR: the qualified-reference notation, the judge's offline form rule, the separated resolver and its outcomes, the foreign-evidence pointer, and the rejections that stay closed
- [x] Revise M-044's exit criterion to the broader boundary this change implements, citing the ADR as its backing decision
- [x] Add the criteria this change is judged by, with positive and negative evidence for each: AC-066 and AC-067 for the form rule and the pointer, AC-068 and AC-069 for the resolver's classification and its rewrite boundary
- [x] Implement the judge's form rule: a bare forge reference fails with the line a reader opens to, while fenced code, code spans, link targets, labelled links, anchors, colour literals, `clue:` identities, and every local citation pass untouched
- [x] Repair every unqualified reference in this corpus, naming the repository each one actually meant
- [x] Implement the resolver command: preview by default, explicit write, five outcomes, a credential sent only where it is honoured, and pinned history reported but never rewritten
- [x] Implement the foreign-evidence pointer: repository, pinned revision, and identifier, parsed as an identity the resolver never follows
- [x] Add the report-only adopter migration and confirm it against both live adopters without writing to either
- [x] Move every live carrier together: the verify checklist, the analysis skill, both generated skill trees, and `[Unreleased]`
- [x] Run the complete change verification: build, vet, the full test suite, and `clue validate`
