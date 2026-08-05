---
id: CAP-003-design
type: design
status: active
links: [CAP-003]
title: Design notes for CAP-003
---

# CAP-003 — design notes

## Split of responsibility

The extraction itself is the `clue-extract` skill ([ADR-008](../../decisions/ADR-008-extraction-is-a-skill.md)) — meaning-level transform, agent-executed, human-reviewed. `clue` carries only the deterministic facets that make an extracted corpus validatable:

- **Namespace grammar** (`internal/corpus/actests.go`, ADR-037): the declaration harvest reads each criteria file's `ac-prefix` (default `AC`); canonical tags use `@<PREFIX>-<digits><lowercase-suffix>`, where `<PREFIX>` is one or more uppercase alphanumeric segments joined by single hyphens. Wrong-namespace tags fail, malformed near-matches are diagnosed, and prefixes whose hyphen-stripped forms collide are rejected. Prefixes may be shared across files; uniqueness stays at the full-ID level (AC-013 and AC-059).
- **Go test names** normalize a segmented prefix by removing its hyphens and retain the numeric component and lowercase suffix, so `SNAP-SQS-001` is `TestSNAPSQS001_…` and `ADP-045b` is `TestADP045b_…`; a classified reference adds its declared type and direction, for example `TestSNAPSQS001_IntegrationPositive_…`. `Unit`/`Sanity`/`Arch` are unchanged. A digit-bearing name whose prefix is in no declared namespace declares no purpose (AC-011) rather than silently passing.
- **JVM harvesting** is per executable: in `*Test.kt`, `*Test.java`, `*Tests.kt`, and `*Tests.java`, a supported Java method or Kotlin `fun` receives classified credit only when its own contiguous JUnit annotation block carries one literal AC tag, proof-type tag, and direction tag together with a supported executable annotation. Underscores in AC tags normalize to hyphens. Parameterized and nested tests are credited once per declared method; enclosing-class tags never flow into their methods. Frameworks without native tags use `test<PREFIX><digits><lowercase-suffix>_<Type><Direction>_<description>` with hyphens removed from segmented prefixes. Ambiguous or class-level syntax is diagnosed without classified credit, never deleted, deduplicated, or relocated by a migration tool; the rehearsal inventories it and human direction selects a split, an actual existing executable, a scoped new test, or an explicit `@draft` with source location and a named plan door when test work is out of scope (ADR-036, ADR-037, PDR-024). Structured comments remain ignored.
- **Provenance** (`internal/corpus/rules.go`): optional field `inferred|verified`, forbidden on decisions (they carry it in `status`); an inferred non-decision also declares `reversal-cost: low|high`, and high-cost inferred meaning blocks an active capability in its one-edge graph slice while low-cost meaning may remain deferred. The CLI reports activation blockers separately from inferred decisions awaiting verification (ADR-010, ADR-035).

## Migration proof parity

`internal/parity` (ADR-049) gives every source mapping a common, machine-checkable coverage check, independent of the mapping's own source-format logic:

- A **source manifest** is pinned and human/agent-authored, written during the `clue-extract` rehearsal (never derived): it names the exact source revision and location read, and per criterion a proof class, direction, and evidence location, or a declared exclusion, or a `draft`/`human`/`retired` disposition with a justification.
- A **target manifest** is always derived, never authored: `DeriveTargetManifest` reuses `corpus.AcceptanceEvidence` — the same declaration-and-evidence harvest `checkACTests` and `Coverage` already walk — so the comparison reads the identical classification the validator enforces, plus the identity ledger's state and recorded source revision for an imported ID.
- `Compare` reports exactly five failure classes: a missing criterion, an orphaned tag, changed direction or evidence location, a stale source fingerprint, and an unjustified `draft`/`Human`/retired disposition. The report is sorted and holds no wall-clock or environment-dependent content, so a CI artifact from the same inputs is reproducible.
- `clue parity <source-manifest> [root]` (`--out` writes the same report to a file) is the CLI entry point, excluded from the ambient release notice for the same reason `clue validate` is: a deterministic judge's output must not depend on another system's present state.

## In-flight source work

An `imported-change` record (ADR-050, `internal/importedchange/`, `docs/imported-changes/`) is the durable replacement for the old mapping's milestone-row-plus-draft-capability treatment of a source repository's still-open pending change:

- One record per source change, pinning its `source-revision`/`source-location` (the same field names ADR-048's ledger and ADR-049's source manifest use for the identical concept), holding the source proposal's intent and design rationale in prose, dependency links to other `imported-change` records or the capability it feeds (the ordinary `links:` field, resolved by `checkLinks`), and a proof-links table mapping each source task to the criterion ID it proves.
- Lifecycle is `in-progress` → `complete`, not the corpus default, and has no `retired` value: the record is never deleted once written, because it is the permanent evidence of what an extracted, now-gone source change once contained (ADR-050's durability argument, distinct from ADR-034's transient-workspace exception).
- `checkImportedChanges` (`internal/corpus/rules.go`) requires every `complete` record's proof-linked criteria to exist, be undrafted, and not be retired, reusing `AcceptanceEvidence`'s declaration harvest rather than inventing a second reading of a criteria file; an `in-progress` record is exempt, since it is still declaring pending work.
- Whether extraction may delete a source repository's incomplete pending change stays agent judgment in the `clue-extract` rehearsal, never a `clue` mechanism: `clue` never reads the source repository, so it cannot itself refuse a deletion there.

## Operational carrier reconciliation

`internal/carriers` (ADR-051) closes a third gap M-052 and M-053 leave open: material that governs how a source repository actually operates but is not itself a criterion or a piece of acceptance evidence — CI workflows, agent instructions, freshness inputs, cross-reference registries, links, and diagram assets. It mirrors ADR-049's split between an authored source side and a derived comparison:

- A **carrier inventory** is pinned and human/agent-authored, written during the `clue-extract` rehearsal (never derived): it names the exact source revision and location read, every source-repository path the migration will delete (`deleted-paths`), and per carrier an `id`, `kind` (one of `instruction`, `workflow`, `freshness-input`, `registry`, `link`, `diagram-asset`), and `source-path`, plus either a `target-path` and pinned `fingerprint`, or an explicit `blocked` marker with a `reason`. `LoadInventory` rejects an entry combining a mapped target with `blocked`, a duplicate ID, or one naming neither outcome.
- `Reconcile` recomputes each mapped entry's target fingerprint from the corpus currently on disk and compares it against the pinned value, reporting exactly three failure classes: a stale deleted-path reference (a Markdown link anywhere in the reconciled corpus still resolving to a `deleted-paths` entry), a lost fingerprint (a mapped target whose current content disagrees with what was pinned), and a missing asset (a mapped target path absent from the corpus). A `blocked` entry is never reconciled against a target — its presence alone is the record of a known gap. The report is sorted and holds no wall-clock or environment-dependent content, so a CI artifact from the same inputs is byte-identical.
- `clue carriers <inventory> [root]` (`--out` writes the same report to a file) is the CLI entry point, excluded from the ambient release notice for the same reason `clue parity` and `clue validate` are: a deterministic judge's output must not depend on another system's present state.
- No source-format parsing in `clue`: what counts as a carrier for a given source format, and how the rehearsal discovers it, stays a `clue-extract` mapping concern, the same split ADR-049 already draws for the source manifest.

## Deliberate limits (doors)

- **Rehearsal before mutation** (`clue-extract`): an extraction's first pass writes a branch-local report under `/changes/` and does not alter the target corpus, routing, tests, or hosted state. The report makes mappings, ID preservation or minting, confidence, test-purpose work, instruction conflicts, planned deletions, and plan doors inspectable; unresolved conflicts stop in `open-questions.md`. Only explicit human direction starts the same full change's mutate phase, which digests the rehearsal to `/docs/analysis` ([PDR-020](../../decisions/PDR-020-extraction-rehearsal-before-mutation.md)).
- Cucumber `.feature` tags are scenario-level proof carriers; other source formats still need a mapping section in the extraction skill.
- No source-format parsing in clue, ever — a new source is a new mapping section in the skill.
- No JVM compilation or runner discovery: the source scanner recognizes ADR-036's conservative annotation and named-executable forms and reports other shapes as unsupported.
- Binary distribution to adopting repos' CI (needed before `clue validate` can be a required check outside this repo) is unsolved and parked in the adopting repo's plan.
