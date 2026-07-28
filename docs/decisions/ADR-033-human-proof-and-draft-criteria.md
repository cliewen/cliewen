---
id: ADR-033
type: decision
status: inferred
links: [ADR-032, ADR-007, ADR-025, CAP-002, CAP-006, PDR-017, P-007]
title: Human proof class, a per-criterion draft exemption, and derived coverage
author: agent
accepted-by: []
---

# ADR-033 — Human proof class, a per-criterion draft exemption, and derived coverage

## Context and problem statement

ADR-032 gave the proof-class vocabulary a consumer but only for evidence a framework can name: `checkACTests` requires classified code references for every criterion in an `active` criteria file. Two cases fall outside that shape. First, a criterion the corpus deliberately verifies by human judgment — not because no test exists yet, but because none ever will — has no honest class to declare; today it either goes untested and fails the build, or its whole capability is held in `draft` to excuse it. Second, a criterion that is genuinely not yet proven, in a file that otherwise is, has no channel short of ADR-007's `@retired` tombstone, which means the opposite thing. Neither gap is provable meaning; both are missing vocabulary. M-025 named both as its own follow-on work rather than reopening its own scope.

## Decision outcome

**The proof-class vocabulary gains `Human`, satisfied by the acceptance brief rather than a code reference. A per-criterion `@draft` tag-line token exempts one criterion from the active-file test requirement without retiring it or drafting its capability. Coverage is a derived report, never a committed file.**

- **`Test-type: Human`** is a fifth value alongside ADR-032's `Unit`, `Integration`, `E2E`, and `Performance`, declared the same way on the scenario's first non-blank body line. It takes no `(single-direction)` qualifier and requires no positive/negative pair: `checkACTests` treats it as satisfied by declaration, requiring neither a Go test, a JVM tag, nor a Cucumber tag. Its evidence carrier is procedural, not mechanical: PDR-017 already requires the acceptance brief to list every added or changed criterion verbatim with its scenario, so a criterion newly or materially declaring `Test-type: Human` is already named there under the existing rule — this decision states explicitly that naming it *is* the human class's proof, and the `clue-delta` and `clue-verify` source templates add one line making that reading unambiguous rather than inventing new brief content or new CI cross-referencing. A capability whose criteria are all `Human` or evidence-backed may reach `active`; declaring `Human` is a claim about how a criterion is always verified, reviewed like any other scenario content, not an escape hatch from review.
- **`@draft`** is a second tag-line token alongside `@retired`, e.g. `@AC-050 @draft`. It marks one criterion as not yet proven — whatever `Test-type` it declares, if any, `checkACTests` skips the "has no test" and coverage-pair checks for that AC while it carries the token. Unlike `@retired`, a `@draft` criterion is alive: it may be referenced, worked toward, and once evidence lands the token is removed in that change. A criteria file's own `status` field governs the file and its capability as before ([ADR-025](ADR-025-one-status-lifecycle.md)); `@draft` governs one criterion inside an otherwise-`active` file, so a capability no longer has to sit in `draft` merely because one of its criteria has no evidence yet.
- **Coverage is a derived report, not a registry.** `clue validate --coverage` prints, per capability, one of the book's three states computed from current corpus state: `covered` (every non-retired criterion has satisfied evidence — a coverage pair, a single-direction exemption, or a `Human`/`@draft` exemption), `partial` (some but not all do), `gap` (none do, or the capability has no non-retired criteria). Nothing is written to disk; a stale registry cannot exist because there is no registry.

**Carrier:** `testTypeRe` and the tag-line parser in `checkACTests` (machine); the `--coverage` flag on `runValidate` (machine); the acceptance-brief line in the `clue-delta` and `clue-verify` source templates and their generated copies (agent); CAP-002 and CAP-006 criteria (evidence).

## Consequences

- A capability can be honestly `active` while carrying criteria proven by a human, criteria not yet proven, and criteria proven by tests, without borrowing `@retired` or holding the whole file back.
- The human class's proof lives in the same acceptance-brief line PDR-017 already mandates; no new hosted-state cross-referencing is added, keeping the merge gate advisory-and-honest rather than a second code-review layer.
- Coverage state is always current because it is computed, never edited by hand, and it disappears the moment the underlying evidence does.
