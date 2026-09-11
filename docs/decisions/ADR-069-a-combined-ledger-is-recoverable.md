---
id: ADR-069
type: decision
status: verified
links: [G-013, CAP-010, ADR-068, ADR-048, ADR-044, C-013]
title: Allocation settings live outside the union-merged ledger, and a combined ledger is recoverable
binds: adopter
author: agent
accepted-by: Flemming N. Larsen (2026-09-11)
---

# ADR-069 — Settings live outside the merged ledger

## Context and problem statement

ADR-068 made the checked-in ledger an append-only event log and asked Git's union merge driver to combine independent additions. The driver takes both sides of a conflicting hunk, which is exactly right for the entry list. It combines line by line, so lines both branches wrote identically merge as ordinary context and only genuinely divergent lines are duplicated — which is why appended events, rewritten files, and two branches enabling coordination against the *same* remote all merge cleanly.

The reachable damage is therefore narrower than it first appears, and it is a duplicated key rather than a duplicated header. Two branches that each enable coordination against their own remote agree on `version` and on `mode: git`, so those merge as context and the result carries `remote` twice inside one `coordination` block. A YAML decode rejects a repeated key outright, wherever it sits.

That case is the one worth removing rather than surviving. Recovering from it still means a repository reaches its integration branch broken, every command stopped, and a person hand-editing an append-only registry to fix it — and the disagreement it encodes is one no reading can settle, because only the team knows which remote it meant. The union driver is right for events and wrong for settings, and the two were only ever in one file because the ledger came first.

The blast radius was the whole tool. Every command loads the ledger, so a file in that state failed `clue validate`, every `clue id` subcommand, and `clue parity` — including the `clue id sync` that ADR-068's own recovery path names. `clue migrate` skipped the ledger silently rather than failing, so no command could put the file right and the only remaining move was to hand-edit an append-only log, which is precisely what the guidance tells people not to do. The message a user saw named a YAML line number, not a merge they did not know had happened.

Refusing the file is defensible only if the content is genuinely lost. It is not: the entries are all present, and the log was designed so their order does not matter.

## Decision outcome

**Allocation settings move out of the merged ledger.** `.clue/id-coordination.yaml` holds the mode and the remote; the ledger keeps the version, the events, and the high-water claims. The union rule names the ledger by path, so the settings file merges normally and two branches that chose different remotes produce an ordinary conflict Git raises at merge time, for the person merging to resolve before anything is committed. Local allocation is the absence of the file rather than a file declaring it, so a repository that never coordinates carries no setting it does not have. An unresolved conflict left in the file is reported as one rather than as a parse error.

**Pointing an established repository at a different remote is refused unless forced.** It abandons every claim recorded on the remote it leaves, so the allocator would begin reissuing numbers that journal had already handed out. Coordinating again to the same remote remains an idempotent no-op.

**A repeated key is reconciled on read; only a contradiction is refused.** The ledger and its coordination block are both decoded through the YAML node API rather than into structs, so repeated keys survive to be folded wherever they appear. Repeated sequences concatenate, which is sound because events fold idempotently and advance state monotonically — their union is the same answer whatever the order, which is the property ADR-068 established. Repeated scalars must agree. Two values that disagree are refused with both named, because picking either would silently move a repository off the remote or the allocation mode its team agreed, and that is a decision only a person can make. A duplicated setting always disagrees in practice, since identical lines never duplicate.

Reconciling a repeated key at the top level of the file is defence rather than an observed case: no merge exercised in evidence produces one, because Git matches the surrounding lines. It is kept because it costs one branch in the same walk and covers a hand-edit, a different merge strategy, or a future driver.

**Saving the ledger rewrites it whole, so any command that writes repairs it.** `clue id repair` is the read-and-write pair on its own, for a repository that needs the file fixed without allocating anything, and `clue migrate` plans the same repair rather than skipping a ledger it cannot load. `clue validate` still fails a combined ledger — the file in the repository is malformed and must not stay — but names the cause and the repair instead of quoting a parser.

This narrows ADR-068's "rejects conflicting identity metadata" to what it was always for: contradictory *identity* metadata, where two events claim different immutable facts about one ID. A key written twice with the same value contradicts nothing.

**Recognizing the merge rule is a byte-level reading, and stays one.** The rule is matched by pattern and attribute rather than as one exact line, so a repository that wrote it without a leading slash is no longer reported non-compliant and refused coordination for a rule that was working. Validate does not ask Git what the attributes resolve to; ADR-044's repository-state boundary is why, and correcting a false negative is not weakening the check.

Ledgers written before the split still carry their settings inline and keep loading, and saving moves them out, so no migration step is required for a representation that has not shipped.

**Carrier:** `internal/ledger`'s loader, coordination file, and attribute reading, `internal/corpus`'s ledger rule, `clue id coordinate --force`, `clue id repair`, the ledger repair in `clue migrate`, `internal/skills/source/skills/clue-delta.md.tmpl`, and the team-allocation guide.
