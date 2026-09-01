---
id: CH-164-open-questions
type: open-questions
status: open
links: [CH-164, P-020, M-083, C-008, ADR-034]
title: Open questions for retiring the spent simplification registers
---

# Open questions

## 1. Three accepted rules make an artifact a completed plan links impossible to retire. Which one gives?

**Status:** blocking. This change stopped here.

Retiring AN-019, AN-020, and AN-021 was implemented and verified green, then reverted, because three rules the corpus already accepts contradict each other the moment a completed plan links the retired artifact:

- [ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md) — retirement is deletion; the file goes and Git history is the archive.
- `checkLinks` — a surviving link to a retired ID fails validation and demands it be repointed to the successor.
- [C-008](../../docs/constraints/C-008-completed-plans-immutable.md) and `.github/scripts/completed-plans.sh` — a completed plan is frozen, and the guard fails on *any* modification to it, with no exemption of any kind.

So repairing the link fails CI and not repairing it fails the judge. There is no third move. P-013's frontmatter links AN-006, AN-008, AN-010, AN-012, AN-013, AN-018, and AN-022, so this is not an edge case — it blocks most of P-020, and it will block any adopter whose completed plans cite the spikes that fed them.

The proposal claimed a "link-target repair" allowance with P-016's decision-log retirement as precedent. **That allowance does not exist.** It appears in neither C-008's text nor its guard, and ADR-034's own outcome says completed plans "remain frozen under C-008". CH-159 did edit P-013 while it was completed, so the precedent is an earlier instance of this same contradiction going unnoticed, not an authorised exception.

**Recommendation: make a completed plan's links historical, and leave C-008's freeze untouched.** A completed plan is a record of what a finished campaign referenced, in the same way AN-023's rows and AN-022's body are pinned records of what they observed at a revision. Reading its outgoing links as live navigation is the error; `checkLinks` should not demand a frozen plan repoint anything. That needs a decision record and a validator change, and it costs nothing that is true today.

The alternative — carving a link-repair exemption into C-008 and its guard — weakens a core freeze and still requires editing a finished campaign's file, which is the thing C-008 exists to prevent.

## 2. Evidence must exist in a commit, not in a working tree

**Status:** resolved by rework, recorded so the mistake is not repeated.

M-083 was marked `done` claiming the three spikes declared `carried-by:` and that `clue migrate` reported them under `MIG-013` before deletion. Both things did happen, and neither is in any commit: the field was added and the files deleted in the same working tree, so the branch never held a state a reader could re-run. The claim was true and unverifiable, which for an evidence column is the same as false.

When this change resumes, `carried-by:` lands in its own commit and the `MIG-013` output is captured against that commit, before any deletion.
