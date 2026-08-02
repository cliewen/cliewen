---
id: CH-106-tasks
type: change
status: open
links: [CH-106]
title: Tasks for CH-106
---

# CH-106 — Tasks

- [x] Confirm the defect against the code rather than the PR description: `checkIndexes` reads no label, `regenIndex` appends bare stems, curated rows survive regeneration.
- [x] Confirm the blast radius in adopted repositories: the same blocks ship through `clue init` and `clue scaffold`.
- [x] Record the blocking design question (Q1) with a researched recommendation.
- [ ] **Stopped for a human answer to Q1.** Implementation does not start until the judge's behaviour — count or fail — is decided, because it determines whether this change also ships a corpus migration.

Once Q1 is answered:

- [ ] Write the ADR recording the index-row contract, including the judge's decided behaviour and the curated-suffix boundary.
- [ ] Teach `regenIndex` to emit `- [<id> — <title>](<file>) · \`<status>\`` for missing entries, reading the target's frontmatter; curated rows keep surviving unchanged.
- [ ] Extend `checkIndexes` to read each row's label and status against the linked artifact's frontmatter, under Q1's decided behaviour. Serves the criterion added below; do not weaken any existing check to accommodate it (C-004).
- [ ] Add the acceptance criterion for the new judge behaviour to CAP-002's criteria with its `Test-type`, and paired positive/negative evidence in that class. A criterion must exist before the behaviour is implemented; if the order slips, add it first and re-tag.
- [ ] Add the criterion for the generator's emitted row format to CAP-005's criteria (or CAP-002's, wherever `scaffold`'s index regeneration is owned — confirm before writing), with paired evidence.
- [ ] Register the constraint in the register with `enforcement: machine` and its `source`.
- [ ] Only if Q1 is answered "fail": add the corpus migration, its preview/apply evidence, and the coordinated release note; re-check whether that pushes the change past what one review can hold.
- [ ] Update every live carrier that states the index-block contract — `docs/decisions/README.md` prose if it must state the rule, `docs/README.md`, and `internal/scaffold/templates/docs/decisions/README.md` — in this same change (C-006).
- [ ] `[Unreleased]` changelog entry: the emitted index row format changes for adopters, and the judge reports a new population. This is user-visible under C-002.
- [ ] Digest: regenerate indexes, tick M-049's evidence column with the honest qualification, delete this workspace.
