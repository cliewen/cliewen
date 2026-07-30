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

- **Namespace grammar** (`internal/corpus/actests.go`): the declaration harvest reads each criteria file's `ac-prefix` (default `AC`); the general tag grammar is `@<PREFIX>-<digits>` with `<PREFIX>` uppercase alphanumeric starting with a letter. Wrong-namespace tags fail, which also keeps prose tokens from becoming accidental declarations. Prefixes may be shared across files; uniqueness stays at the ID level (AC-013).
- **Go test names** generalize to `Test<PREFIX><digits>_…` for every declared prefix; a classified reference adds its declared type and direction, for example `TestMG010_IntegrationPositive_…`. `Unit`/`Sanity`/`Arch` are unchanged. A digit-bearing name whose prefix is in no declared namespace declares no purpose (AC-011) rather than silently passing.
- **JVM harvesting** is per executable: in `*Test.kt`, `*Test.java`, `*Tests.kt`, and `*Tests.java`, a supported Java method or Kotlin `fun` receives classified credit only when its own contiguous JUnit annotation block carries one literal AC tag, proof-type tag, and direction tag together with a supported executable annotation. Underscores in AC tags normalize to hyphens. Parameterized and nested tests are credited once per declared method; enclosing-class tags never flow into their methods. Frameworks without native tags use `test<PREFIX><digits>_<Type><Direction>_<description>`. Ambiguous or unsupported syntax is diagnosed without classified credit, and structured comments remain ignored (ADR-036).
- **Provenance** (`internal/corpus/rules.go`): optional field `inferred|verified`, forbidden on decisions (they carry it in `status`); an inferred non-decision also declares `reversal-cost: low|high`, and high-cost inferred meaning blocks an active capability in its one-edge graph slice while low-cost meaning may remain deferred. The CLI reports activation blockers separately from inferred decisions awaiting verification (ADR-010, ADR-035).

## Deliberate limits (doors)

- **Rehearsal before mutation** (`clue-extract`): an extraction's first pass writes a branch-local report under `/changes/` and does not alter the target corpus, routing, tests, or hosted state. The report makes mappings, ID preservation or minting, confidence, test-purpose work, instruction conflicts, planned deletions, and plan doors inspectable; unresolved conflicts stop in `open-questions.md`. Only explicit human direction starts the same full change's mutate phase, which digests the rehearsal to `/docs/analysis` ([PDR-020](../../decisions/PDR-020-extraction-rehearsal-before-mutation.md)).
- Cucumber `.feature` tags are scenario-level proof carriers; other source formats still need a mapping section in the extraction skill.
- No source-format parsing in clue, ever — a new source is a new mapping section in the skill.
- No JVM compilation or runner discovery: the source scanner recognizes ADR-036's conservative annotation and named-executable forms and reports other shapes as unsupported.
- Binary distribution to adopting repos' CI (needed before `clue validate` can be a required check outside this repo) is unsolved and parked in the adopting repo's plan.
