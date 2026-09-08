---
id: ADR-069
type: decision
status: inferred
links: [G-013, CAP-010, ADR-068, ADR-048, ADR-044, C-013]
title: A ledger Git's union merge combined is recoverable, not fatal
binds: adopter
author: agent
accepted-by: []
---

# ADR-069 — A combined ledger is recoverable

## Context and problem statement

ADR-068 made the checked-in ledger an append-only event log and asked Git's union merge driver to combine independent additions. The driver takes both sides of a conflicting hunk, which is exactly right for the entry list and wrong for the few lines above it: a hunk covering the top of the file yields two `version` and `coordination` keys, and a YAML decode rejects a repeated key outright.

The blast radius was the whole tool. Every command loads the ledger, so a file in that state failed `clue validate`, all four `clue id` subcommands, and `clue parity` — including the `clue id sync` that ADR-068's own recovery path names. `clue migrate` skipped the ledger silently rather than failing, so no command could put the file right and the only remaining move was to hand-edit an append-only log, which is precisely what the guidance tells people not to do. The message a user saw named a YAML line number, not a merge they did not know had happened.

Refusing the file is defensible only if the content is genuinely lost. It is not: the entries are all present, and the log was designed so their order does not matter.

## Decision outcome

**A repeated header is reconciled on read; only a contradiction is refused.** The ledger is decoded through the YAML node API rather than into a struct, so repeated keys survive to be folded. Repeated sequences concatenate, which is sound because events fold idempotently and advance state monotonically — their union is the same answer whatever the order, which is the property ADR-068 established. Repeated scalars must agree. Two halves declaring different coordination modes are refused with the choice stated, because picking either would silently move a repository off the allocation mode its team agreed, and that is a decision only a person can make.

**Saving the ledger rewrites it whole, so any command that writes repairs it.** `clue id repair` is the read-and-write pair on its own, for a repository that needs the file fixed without allocating anything, and `clue migrate` plans the same repair rather than skipping a ledger it cannot load. `clue validate` still fails a combined ledger — the file in the repository is malformed and must not stay — but names the cause and the repair instead of quoting a parser.

This narrows ADR-068's "rejects conflicting identity metadata" to what it was always for: contradictory *identity* metadata, where two events claim different immutable facts about one ID. A header written twice with the same content contradicts nothing.

**Recognizing the merge rule is a byte-level reading, and stays one.** The rule is matched by pattern and attribute rather than as one exact line, so a repository that wrote it without a leading slash is no longer reported non-compliant and refused coordination for a rule that was working. Validate does not ask Git what the attributes resolve to; ADR-044's repository-state boundary is why, and correcting a false negative is not weakening the check.

**Carrier:** `internal/ledger`'s loader and attribute reading, `internal/corpus`'s ledger rule, `clue id repair`, the ledger repair in `clue migrate`, `internal/skills/source/skills/clue-delta.md.tmpl`, and the team-allocation guide.
