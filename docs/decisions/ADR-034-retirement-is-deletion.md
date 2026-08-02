---
id: ADR-034
type: decision
status: verified
links: [ADR-025, PDR-003, C-008, CAP-002, P-007]
title: Retirement is deletion; supersedes carries the pointer forward
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-034 — Retirement is deletion; supersedes carries the pointer forward

## Context and problem statement

ADR-025 gives every default-lifecycle type a `draft → active → retired` vocabulary, but no artifact in the corpus has ever actually reached `retired`: the conventions that exist in practice are a criteria tombstone (`@retired` tag, the file stays so a stale test tag keeps failing loudly) and a demoted decision (the file is deleted and its content folded into a dated row in `docs/decisions/log.md`, per PDR-003). Both are deletion in substance. A `retired` status value that no committed file is ever observed holding is not a reachable terminal state — it is a word in a vocabulary table with no corresponding fact on disk. Deletion also currently leaves no machine-visible trace of what replaced a removed artifact beyond prose ("supersedes ADR-006's …") and git history, so a reader who has not memorized the corpus's history cannot tell, from the files alone, that a dangling-looking reference in an old PR or a changelog line once pointed at something real.

## Decision outcome

**Retiring any artifact means deleting its file in the same change that retires it.** There is no committed state in which a file carries `status: retired` on `main` — that value never needs to exist, because the file it would describe is gone. `retired` is removed from the default lifecycle ADR-025 states; retirement is an event (deletion), not a status a surviving file holds.

**A `supersedes:` frontmatter field, optional on any artifact type, names the IDs a retirement deleted.** It is carried by whichever artifact is the reader's best next stop: the successor record where a direct one exists, or — for a decision demoted with no successor record — `docs/decisions/log.md` itself, whose own frontmatter accumulates every ID it has absorbed since this decision (a decision demoted before this one keeps its existing plain-prose row; the field is not retrofitted onto history). Git history remains the archive of the full retired text; `supersedes:` is only the pointer, not a copy.

**The validator enforces the half of this it can see:** an ID named in any artifact's `supersedes:` list must not resolve to a live artifact — if it does, the retirement was declared but not actually done, and `clue validate` rejects it. This is new: `checkLinks` already rejects any other artifact's `links:` entry that resolves to nothing, which is the enforcement half that already existed and needed no new rule — a live artifact still pointing at a deleted ID was already a validation failure before this decision; `supersedes:` gives that failure a resolution path (repoint the link to the ID naming the successor in its `supersedes:` list) instead of a bare "resolves to no artifact."

**Two exceptions stand, because they are not this kind of retirement:**

- **Criteria tombstones.** A retired acceptance criterion keeps its `@retired` tag and its file, because a stale test tag has to keep failing loudly rather than disappear silently (ADR-007). It is deliberately the one case where "retired" is a status a surviving artifact holds.
- **Completed plans.** [C-008](../constraints/C-008-completed-plans-immutable.md) keeps a completed plan frozen and on disk forever; it is not retired at all, and nothing here changes that.

**This decision supersedes two specific clauses, not the records that carry them:**

- **PDR-003's demotion mechanic** — "file deleted, row added, inbound references repointed" — is unchanged in effect but now names its mechanism: the row is (or, for a demoting change, gains) the `supersedes:` pointer on `log.md`, and "inbound references repointed" is what `checkLinks` was already enforcing.
- **ADR-025's claim that `retired` is a reachable terminal state** is corrected for the default lifecycle and goals: no file was ever observed reaching it, and this decision states why that was never a bug to fix — the state was never meant to be a resting place. `ADR-025`'s other exception vocabularies (plan, decision, log, transient workspace types) are unchanged.

**Carrier:** the `Supersedes` field on `internal/corpus.Artifact` and the `checkSupersedes` rule in `internal/corpus/rules.go` (machine); the retirement sentence in `clue-delta`'s Implement step, extended to name `supersedes:` alongside the existing criteria-tombstone sentence (agent); the decisions folder README and `docs/decisions/log.md`'s header preamble (default).

## Consequences

- `clue validate` can tell "this file was deleted and its retirement is machine-visible" from "this file was deleted and something still links to it" — the second case fails, naming a fix.
- ADR-025's status table shrinks by one value for the types that shared it; nothing that ever relied on `retired` existing (nothing does) breaks.
- A reader following a `supersedes:` pointer never has to grep git log to find out what happened to a missing ID mentioned in old prose.

### Rejected: keep `retired` as a real status and require a tombstone file for every type

This is what criteria already do, and it is the right shape for a test tag that must keep failing. Generalizing it to every type would mean every capability, design, decision, and analysis record still exists forever, just inert — which is exactly the "graph only accumulates" pattern [AN-008](../analysis/AN-008-methodology-critiques.md) named as one of the four patterns this campaign works core-first. Deletion plus a machine-visible pointer keeps the corpus's live size honest while keeping the successor findable.

### Rejected: no structured field, prose supersession only

Prose already carries "supersedes ADR-006" today and is not going away for partial-clause supersession between two records that both stay on disk (ADR-024's "supersedes only its published-URL clause" of ADR-023 is exactly that, unaffected by this decision). But for a fully retired-and-deleted ID, prose alone gives a human reader a pointer and gives the validator nothing to check — the corpus could accumulate dangling prose claims with no way to verify any of them stayed true as the corpus grew. A field the validator reads closes that gap for the one case that matters most: an ID that no longer exists anywhere.
