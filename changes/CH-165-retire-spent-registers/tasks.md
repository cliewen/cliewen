---
id: CH-165-tasks
type: tasks
status: open
links: [CH-165]
title: Tasks for retiring the spent simplification registers
---

# Tasks

- [ ] Declare `carried-by:` on AN-019, AN-020, and AN-021 and commit that alone, so the claim that their findings survived exists in history rather than in a working tree.
- [ ] Capture `clue migrate`'s `MIG-013` report against that commit as the milestone's evidence.
- [ ] Repoint AN-022's link, which is an active artifact and repairs like anything else.
- [ ] Delete the three files and name each retired identity in a `supersedes:` field on exactly one live successor.
- [ ] Retire the three identities in the ledger and regenerate the analysis index.
- [ ] Confirm P-013 needs no edit and its CI guard passes, which is what ADR-063 bought.
- [ ] Assess documentation impact and record whether any overview changed; add a release note only if the shipped surface changed.
- [ ] Run focused and full verification, complete P-020's M-084 digest, and prepare the reviewed pull-request handoff.
