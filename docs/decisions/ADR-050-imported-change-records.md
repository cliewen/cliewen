---
id: ADR-050
type: decision
status: verified
links: [ADR-034, ADR-048, ADR-049, P-011, CAP-003, C-013]
title: In-flight source work becomes a durable imported-change record, never a transient workspace
author: agent
accepted-by: Flemming N. Larsen (2026-08-08, conversation)
---

# ADR-050 — In-flight source work becomes a durable imported-change record, never a transient workspace

## Context and problem statement

The current OpenSpec extraction mapping converts a source repository's pending change (`changes/<name>/`: `proposal.md`, `design.md`, `tasks.md`) into a milestone row in the target plan plus a `status: draft` capability holding its criteria. That row states plainly that the pending change's `tasks.md` "dies" — `clue-delta` regenerates tasks once implementation starts in the target corpus. A source `tasks.md` is not only a checklist: its ordering encodes a dependency graph between tasks, and the tasks that reference test evidence carry a proof link from a task to the criterion it satisfies. Discarding the file discards both, leaving a milestone row and a draft capability that state what remains to be done but not why the source work was designed the way it was, what it depended on, or which of its tasks already had proof before extraction. [P-011](../plans/P-011-truthful-brownfield-migration.md)'s M-054 names this "reduces in-flight work to an insufficient plan row" and commits to a record durable enough that a proposal, design, dependency, and proof task remain inspectable without keeping the source corpus around as a parallel registry.

What should hold that information, and does it behave like the transient `/changes/` workspace this same corpus already has, or like something that survives past the extracting change's own digest?

## Decision outcome

**A new native artifact type, `imported-change`, holds one record per source change extraction preserves. It is durable, not transient: once written it is never deleted by the digest that creates it, and its own lifecycle has no `retired` state to fall back to.**

- **Location:** `docs/imported-changes/`, one file per source change, mirroring `docs/decisions/`'s one-record-per-file layout with a `README.md` carrying a `clue:index` block (`checkIndexes` enforces it the same way it enforces every other taxonomy folder). `docs/README.md`'s Folders index and status-vocabulary table gain a row for it, per the same-change carrier convention every corpus-wide table already follows.
- **Identity:** minted through the ledger ([ADR-048](ADR-048-corpus-wide-id-ledger.md)) under a new prefix, `IC`. The ledger is prefix-agnostic — `clue id next IC` allocates through the same numeric counter and O(1) lookup path any other native prefix uses — so no ledger change is needed to support it.
- **Lifecycle:** `in-progress` → `complete`, added to `internal/corpus/rules.go`'s `statusVocabExceptions` (mirroring `docs/README.md`'s status table, per the existing same-change convention). `complete` names the milestone-level fact M-054 calls "completion state" directly as the status value, rather than inventing a second field that duplicates it. There is no `retired` value in this vocabulary at all — see the durability question below.
- **Frontmatter:** `source-revision` and `source-location` reuse the exact field names [ADR-048](ADR-048-corpus-wide-id-ledger.md)'s ledger entries and [ADR-049](ADR-049-migration-parity-manifests.md)'s source manifest already use for the same concept — the pinned commit/tag and path the source change was read from — so a reader who knows either carrier already knows this one. `links:` carries dependencies on other `imported-change` records and the capability or capabilities the change feeds, resolved by the ordinary `checkLinks` rule like any other artifact link — no new link-resolution mechanism is needed.
- **Body structure:** four sections every record carries — **Origin** (what source change this was, in prose), **Intent**, **Design rationale**, and a **Proof links** table mapping each source task to the criterion ID it proves and that criterion's evidence state. The proof-links table is the piece the old mapping discarded outright; it is what lets a reader trace "this source task" to "this target criterion" after the source repository is gone.
- **Machine check (`checkImportedChanges`, `internal/corpus/rules.go`):** a record whose status is `complete` must have every criterion ID in its proof-links table resolve to a declared, non-`@draft`, non-retired criterion (reusing `corpus.AcceptanceEvidence`'s declaration harvest — the same one `checkACTests` and `clue parity` already trust, so this rule invents no second reading of a criteria file). A `complete` record naming a criterion that does not exist, or is still `@draft` or retired, is an unjustified completion claim — the same failure shape [ADR-049](ADR-049-migration-parity-manifests.md) already gave parity's unjustified-disposition class, applied to this record's own claim instead of a source manifest's. An `in-progress` record has no such requirement: it is allowed to name proof links that are not yet satisfied, because it is still declaring what remains to be done.
- **What stays agent judgment, not a machine check:** `clue` never reads the source repository, so it cannot itself refuse to let a human delete source work. The rule that extraction must not delete an incomplete source change until its `imported-change` record reaches `complete` is `clue-extract` skill prose, exercised by the extracting agent during the rehearsal's mutate phase — the same split [ADR-049](ADR-049-migration-parity-manifests.md) already draws between what `clue parity` checks mechanically and what the rehearsal's human-reviewed resolution decides.

## Durability: why this is not a transient workspace artifact under ADR-034

[ADR-034](ADR-034-retirement-is-deletion.md) makes retirement of a default-lifecycle artifact mean deleting its file, and separately carves out `change`/`tasks`/`open-questions` as transient workspace types with a deliberately short custom vocabulary (`open` only) that dies at digest. An `imported-change` record is not that shape: it is written *by* a digest (the same commit that closes the extracting change) to be the permanent record of what a deleted source change once contained. Treating it as transient — deleting it once the milestone or capability it fed is itself done — would recreate exactly the defect this ADR closes: a reader loses the design rationale and proof-link trace the moment the record disappears, with no `git log` recovery path once the source repository (not just this record) is gone. So `imported-change` sits with `docs/decisions/`'s "decisions are never deleted" convention and [C-008](../constraints/C-008-completed-plans-immutable.md)'s frozen completed plans, not with the transient workspace types: once written, a record is kept forever, `complete` is a resting state a file may hold indefinitely, and there is no `retired` value because nothing about this type is ever retired — the source it describes already is, and that is the fact the record exists to survive.

## Rejected: keep the milestone-row-plus-draft-capability treatment, add only a narrative appendix

An appendix document describing "what the source tasks.md used to say" would carry the same prose a human could write once and never re-check. It gives no machine-checkable claim that a `complete` disposition is honest — exactly the unchecked-claim problem [ADR-049](ADR-049-migration-parity-manifests.md) closed for parity by deriving a target manifest instead of trusting an authored one. Making the proof-links table a first-class frontmatter-adjacent structure `checkImportedChanges` can read keeps the same discipline: an agent's claim that source work is done is checked against the corpus, not taken on prose alone.

## Rejected: fold imported-change into the extracted capability's design.md

A capability's `design.md` documents the capability as it stands today; a source change's dependency graph and proof-link history are about *how it got there*, spanning potentially several source changes that feed one capability or one source change that feeds several. Collapsing them would either duplicate the same dependency and proof-link content across every capability a multi-capability source change touches, or force an arbitrary single owner — the same one-to-many mismatch that keeps `docs/decisions/` and `docs/capabilities/` as separate folders today rather than folding decisions into the capability they informed.

## Rejected: reuse the `change` type with a widened status vocabulary

`change`'s existing `open` vocabulary and `ForbidChanges` gate exist specifically to keep `/changes/` empty on `main` — the digest-before-merge contract every full Cliewen change already follows. Widening that one type's vocabulary to include a durable `complete` state would either break the transient gate for every ordinary Cliewen change workspace or require the validator to distinguish two different `change` records by folder alone, which is exactly the kind of implicit distinction ADR-026's adopter-type default lifecycle is designed to avoid. A new type name keeps the transient contract and the durable contract each single-purpose.

## Carrier

This record, `internal/corpus/rules.go`'s `checkImportedChanges` rule and `statusVocabExceptions` entry, `internal/importedchange/`, `docs/imported-changes/`, `docs/README.md`'s Folders index and status table, `docs/capabilities/CAP-003-extract/criteria.md` (AC-115..AC-117), and the `clue-extract` skill's OpenSpec mapping and rehearsal guidance are the carriers of this decision.
