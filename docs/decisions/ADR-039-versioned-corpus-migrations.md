---
id: ADR-039
type: decision
status: verified
links: [P-009, M-041, ADR-025, ADR-035, ADR-038, CAP-001, CAP-004, C-004, C-011]
title: Versioned corpus migrations plan safe mechanical upgrades
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# ADR-039 — Versioned corpus migrations plan safe mechanical upgrades

## Context and problem statement

A released `clue` can narrow the corpus contract while an adopted repository still carries the older shape. `clue init` must remain non-destructive, but manual edits do not provide a preview, a safe boundary around semantic choices, or one operation that keeps the corpus, generated skills, and CI caller together. A migration mechanism must also distinguish an old generated carrier from a locally edited one without keeping hidden repository state or asking the deterministic judge to use the network.

## Decision outcome

**`clue migrate` uses an ordered, versioned registry of mechanical migrations and plans the complete target before writing.** The first registry entries are `MIG-001` for the explicit `reversal-cost` field, `MIG-002` for the old architecture and analysis `status: verified` lifecycle value, and `MIG-003` for generated skills, their Claude mirror when present, and the thin CI caller. Each entry is idempotent: current content is a no-op, and a later release appends a new migration rather than rewriting the meaning of an earlier one.

**Preview is the default and `--apply` is explicit.** A plan lists every changed path and transformation. A preflight finding — an absent semantic choice, unsupported syntax, a missing required carrier, a copied workflow, or a locally modified generated file — prevents all writes. Applying a plan rechecks the planned source bytes immediately before writing. The command preserves line endings and all content outside the named field or managed carrier, and it has no state file to corrupt or consult on a later run.

**Missing inferred-artifact routing is a human choice.** `MIG-001` reports an inferred non-decision without `reversal-cost` until the adopter supplies `--reversal-cost=low|high`; it never guesses from an artifact's type or prose. The status migration is deterministic only for the historical `architecture` and `analysis` values. Retired default-lifecycle artifacts, unusual status forms, and other semantic cases remain unchanged and are reported for a reviewed manual repair.

**Generated carriers are replaced only from a release manifest.** The manifest carries SHA-256 digests for the supported prior generated skill releases. A matching old digest may be replaced by the current embedded generated content; a marked file with any other bytes is reported as locally modified or unsupported and is never overwritten. New generated files introduced inside a recognized managed skill may be added. A whole newly introduced canonical skill directory may be added only when every carrier a supported preceding release shipped exactly matches that release and the preview names both the recognized preceding release and the target release whose bytes it writes; a partial or modified set remains a finding. Missing required skill files and a missing thin caller are reported. A missing or symlinked Claude mirror is outside the repository's write boundary and remains an explicit notice. The caller update changes only its upstream workflow reference and pair version; runner, binary-source, and install-directory choices remain untouched.

**Release workflow changes to the migration registry require a migration section in that release's changelog entry.** The release gate checks the tag's first-parent diff for `internal/migrate/` and refuses to publish when the extracted notes lack a `### Migration` section. The release notes therefore tell adopters when a corpus obligation was added or narrowed, while ordinary releases keep the existing notes contract.

The judge remains offline and deterministic. `clue migrate` is an explicit reviewed repository operation, not a background updater, and `clue init` continues to skip existing files rather than replacing them.

## Rejected: rewrite every frontmatter block

Marshalling YAML would discard comments, ordering, quoting, and local formatting unrelated to the migration. Line-preserving edits keep the adopter's prose and make the preview an exact account of the mutation.

## Rejected: overwrite every marked skill

The ownership marker identifies Cliewen's carrier, not proof that the local bytes are unmodified. Replacing unknown bytes could erase a repository's deliberate repair or local convention. The release digest manifest gives the updater a conservative positive identity test and reports the rest.

## Carrier

`internal/migrate`, `cmd/clue migrate`, the embedded carrier renderer in `internal/scaffold`, the release workflow guard, and the operations and capability guidance carry this decision.
