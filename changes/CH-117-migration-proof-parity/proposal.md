---
links: [P-011]
---

# CH-117 — Migration proof parity is reproducible at criterion and evidence-location level

## What

Closes M-053 of [P-011](../../docs/plans/P-011-truthful-brownfield-migration.md) by giving every brownfield extraction mapping a common, machine-checkable parity contract:

- A **source manifest** format ([ADR-049](../../docs/decisions/ADR-049-migration-parity-manifests.md)): a pinned, human/agent-authored YAML file, one per source mapping run, naming the source revision and location and, per criterion, its proof class, positive/negative direction, evidence location, and either a declared exclusion (with reason) or a disposition justification for a `@draft`, `Human`, or retired outcome.
- A **target manifest**, derived deterministically by `clue` from the current corpus and ledger state: per declared criterion, its ledger state (live/reserved/retired), proof class, direction(s), evidence location(s), and draft/Human/retired flags — reusing the existing AC-evidence harvest (`internal/corpus`) rather than re-parsing anything.
- `internal/parity/` — manifest types, a `Compare` that diffs source against derived target and reports every unmatched or altered entry, and the five required failure classes: missing criterion, orphaned tag, changed direction/location, stale source fingerprint, and unjustified `@draft`/`Human`/retirement disposition.
- `clue parity <source-manifest> [root]` — a new CLI command that derives the target manifest, runs the comparison, prints the report, and exits non-zero on any failure class; a `--out` flag writes the same report as a deterministic file for a migration workflow to upload as a CI artifact.

## Why

P-009's OpenSpec mapping proved extraction can preserve identities and evidence classification, but nothing today checks that a migration PR's claimed coverage still matches what the source actually declared once the source corpus is deleted. Without a reproducible parity contract, a mapping can silently drop a criterion, relocate evidence without a trace, or mark work `@draft`/`Human`/retired with no accountable reason — exactly the general source-to-corpus parity failures P-011 opened to close (see the plan's introduction). M-052's ledger closes the identity-reuse half of that gap; M-053 closes the coverage-and-evidence half.

## Scope

In scope: the four build items above, focused positive and negative fixtures for each of the five required failure classes plus one clean baseline, new acceptance criteria in `CAP-003-extract` (AC-109..AC-114), `clue-extract` skill and OpenSpec-mapping guidance updated to require emitting a source manifest during rehearsal and running `clue parity` before mutation, and closing M-053 in the P-011 plan table in the digest.

Out of scope: wiring `clue parity` into a real adopting repository's hosted CI workflow (no adopter is migrating right now) and the disposable end-to-end fixture-source rehearsal proving the full contract — both belong to M-056. Also out of scope: `imported-change` records for in-flight source work (M-054) and operational-carrier reconciliation (M-055).
