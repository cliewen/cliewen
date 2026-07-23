---
id: ADR-026
type: decision
status: inferred
links: [PDR-013, ADR-025, C-004]
title: Unknown artifact types are adopter extensions, validated against the default lifecycle
author: agent
accepted-by: []
---

# ADR-026 — Adopter-defined types validate against the default lifecycle

## Context and problem statement

[PDR-013](PDR-013-explicit-core-red-line.md) says adopters extend Cliewen by putting their own artifacts — including their own artifact *types* — into their corpus under `docs/`, and that the core fixes only what the thread, the merge boundary, and the judge mean, not which folders a corpus may contain. But `checkStatusVocab` held a closed-world type registry: any `type` it did not recognize produced an `unknown type` issue and a red validate. An adopter who added a `risk`, `experiment`, or `runbook` type could not get a green build, so the extension story was words without a mechanism. How should the validator treat a type it does not know?

## Decision outcome

**A type the validator does not recognize is an adopter extension. It is validated against the default lifecycle `draft → active → retired` ([ADR-025](ADR-025-one-status-lifecycle.md)), not rejected.** `checkStatusVocab` no longer emits `unknown type`; `statusVocabFor` returns the default for any type without a listed exception, and unknown types reach it the same way the built-in default types do.

**This re-scopes a check; it does not weaken one ([C-004](../constraints/C-004-never-weaken-checks.md)).** Every other guarantee still applies to an adopter-defined type, unchanged:

- core frontmatter fields (`id, type, status, links, title`) must be present and non-empty;
- IDs must be unique across the corpus;
- every `links` entry must resolve to an artifact or a milestone;
- the `status` must be a valid default-lifecycle value — a typo like `status: acive` on any type still fails.

What the check stops asserting is that the *set of type names* is closed. That assertion was never a form guarantee; it was a registry Cliewen had no authority to impose on an adopter's domain. Removing it narrows the check to what it can legitimately judge — the shape of each artifact — and hands the choice of vocabulary to the corpus owner, which is precisely PDR-013's division between core and periphery. The one guarantee genuinely given up is that a *misspelled built-in type* (`goad` for `goal`) is no longer caught by name; it is instead caught the moment its status falls outside the default (a misspelled `goal` carrying `status: accepted` still fails, since `accepted` is not a default-lifecycle value). This residual gap is the correct price for an open type system and is smaller than it looks.

**Carrier:** `statusVocabFor` and the removal of the `unknown type` branch in `internal/corpus/rules.go`; the extension paragraph in [ARCH-003](../architecture/core.md) and the `clue init` AGENTS.md red-line rule state the adopter-facing promise.

### Rejected: a configurable list of adopter types

An adopter could declare their extra types in a config file the validator reads. This keeps a closed world with an extra step: the adopter maintains a registry, forgets to add a type, and gets the same red build. The open default is simpler and matches how the core already refuses to enumerate corpus contents.

### Rejected: warn on unknown types instead of accepting them

A warning that never fails a build is noise that trains adopters to ignore validator output; a warning that fails a build is the rejection under another name. Neither serves an adopter whose type is deliberate.
