---
id: ADR-071
type: decision
status: inferred
links: [CAP-002, CAP-003, ADR-005, ADR-032, ADR-036, ADR-044, PDR-019]
title: The judge consumes framework-neutral evidence exported by the repository
author: agent
accepted-by: []
binds: adopter
---

# ADR-071 — Framework-neutral evidence export

## Context and problem statement

Acceptance criteria describe behavior in Gherkin, while executable proof belongs to each repository's chosen test frameworks. A parser catalogue inside the judge makes adoption depend on upstream knowledge of those frameworks.

## Decision outcome

All executable evidence reaches the judge through the versioned `.clue/evidence.yaml` contract. Repository-owned, repeatable exporters attribute each reference to one executable and aggregate every declared producer before replacing the manifest. Native tags, markers, properties, or categories are preferred; executable-bound custom annotations, attributes, or decorators come next; stable names or titles are the fallback. Ordinary proximity comments and container-level AC identities supply no proof. A producer identity and source-qualified executable identity distinguish references across frameworks and suites.

The judge checks canonical identities, proof classification, references, exporter diagnostics, and fingerprints of the complete declared input scopes. It runs neither exporters nor test runners and interprets no framework syntax. Exporters own discovery and attribution; their tests and human review establish those claims. A current fingerprint proves freshness, not test execution, complete discovery outside the declared scopes, or behavioral meaning. Existing Human, draft, retired, legacy, and direction contracts remain unchanged.

This amends ADR-005 and ADR-036 by replacing their per-language harvester boundary with repository-owned export, while retaining native metadata and per-executable attribution. It extends ADR-032's carrier list. The [architecture overview](../architecture/README.md) describes the producer/judge boundary; the [design overview](../design/README.md) reaches the common contract. The shipped carrier is `internal/skills/source/shared/evidence-workflow.md.tmpl`, used by the lifecycle skills, together with `internal/scaffold/templates/evidence/README.md` and its exporter examples.
