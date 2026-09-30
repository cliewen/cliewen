---
id: CH-194
type: change
status: open
links: [CAP-002, CAP-003, G-001]
title: Framework-neutral acceptance evidence from repository-owned exporters
---

# CH-194 — Framework-neutral acceptance evidence

This change is plan-less. It implements the human-directed design from the September 30 conversation: Gherkin states behavior, an agent establishes tests and repeatable export in each repository's actual frameworks, and Cliewen checks one versioned evidence format for every framework instead of maintaining framework parsers in its judge.

## Contract

All producers export to `.clue/evidence.yaml`. Each executable has one AC identity, proof type, and direction, with native metadata preferred, executable-bound custom metadata next, and a stable name or title as the fallback. Ordinary proximity comments and container AC identities give no credit. Multiple producers share one manifest without losing their executable identities. Input scopes and normalized content fingerprints make stale, added, and deleted inputs visible from a clean checkout. Cliewen never invokes exporters or runners; repository-owned exporters own discovery and attribution, which their tests and review establish.

Existing Go/JVM/Cucumber extraction moves out of the judge. Repository examples cover their conventions plus Python unittest, pytest, JUnit, NUnit, Vitest, and Gatling. The migration reports the need to establish export, and this source repository adopts its own export in the same change. Legacy proof, Human proof, draft, retired references, and direction requirements retain their meanings.

## Challenge

The riskiest assumption is that repeatable repository-owned export can stay complete and cheaper than a growing central parser catalogue. Keeping built-in profiles with an external extension is the credible alternative, but retains two evidence paths and framework maintenance in the judge. The cheapest useful test is an exporter fixture path for Go, Python, Vitest, and Gatling, repeated byte-for-byte and then exercised with added, moved, and deleted inputs. A need for manual AC mappings or undetected source changes stops removal of the built-in harvesters and requires revising the approach.

A formally valid manifest could still misattribute a container tag or name a Gatling scenario with no relevant assertions. Exporter attribution regressions and the human acceptance brief must expose that residual; input hashes establish freshness, not the meaning or correctness of tests.

## Verification and carriers

Add classified Gherkin criteria and focused Go evidence for generic import, multiple producers, freshness, diagnostic propagation, metadata attribution, and migration. Inventory current corpus, canonical skills, generated skills, scaffold templates, guide, contributor guidance, CLI and distribution metadata; repair all live claims in this change while leaving historical analyses and completed plans pinned. Use the declared draft VIS-001; this change does not promote the vision.

Run the CONTRIBUTING.md gates, guide build, example checks, and clue-verify review on the committed candidate. Publish a draft PR before implementation, digest the workspace before readiness, and leave acceptance and release to the human boundary.

## Live evidence-contract carriers

The inventory covers the current judge/import/coverage/parity code, exporter examples and tests, capability criteria/design, architecture/design overviews, constraints and current decision outcomes, guide, contributor and CLI text, canonical lifecycle skills and their generated entrypoints/references, scaffold documentation and PR template. Historical analyses, completed plans and changelog history remain pinned.

- `docs/README.md`
- `docs/capabilities/CAP-002-validate/criteria.md`
- `docs/capabilities/CAP-003-extract/criteria.md`
- `docs/capabilities/CAP-003-extract/design.md`
- `docs/decisions/ADR-005-test-reference-convention.md`
- `docs/decisions/ADR-006-test-purpose-taxonomy.md`
- `docs/decisions/ADR-032-classified-ac-evidence.md`
- `docs/decisions/ADR-036-jvm-evidence-per-executable.md`
- `docs/decisions/ADR-071-framework-neutral-evidence-export.md`
- `docs/decisions/PDR-019-methodology-contract-carriers-move-together.md`
- `guide/change-loop.md`
- `guide/design.md`
- `guide/getting-started.md`
- `guide/methodology.md`
- `guide/operations.md`
- `guide/what-you-can-do.md`
- `internal/scaffold/templates/docs/README.md`
- `internal/scaffold/templates/docs/capabilities/README.md`
- `internal/skills/source/skills/clue-verify.md.tmpl`

Additional carriers: `AGENTS.md`, `CONTRIBUTING.md`, `docs/architecture/core.md`, `docs/architecture/README.md`, `docs/design/README.md`, `internal/scaffold/templates/github/pull_request_template.md`, `.github/pull_request_template.md`, `internal/scaffold/scaffold.go`, `internal/migrate/migrate.go`, `cmd/clue/main.go`, and their semantic guards. Generated copies are repaired by `go generate ./internal/skills`, never by hand. Source export and its focused CI regeneration are local; adopter workflow changes remain opt-in.
