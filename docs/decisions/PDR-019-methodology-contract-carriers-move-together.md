---
id: PDR-019
type: decision
status: verified
links: [G-001, P-008, AN-011, PDR-013, ADR-032, C-006, C-013]
title: Methodology contract changes update every live carrier in the same change
author: agent
accepted-by: Flemming N. Larsen (2026-08-02, conversation)
---

# PDR-019 — Methodology contract changes update every live carrier in the same change

## Context and problem statement

Cliewen can mechanically validate its corpus, generated-skill identity, acceptance-evidence references, and guide build while live instructions still disagree about the methodology contract. A contract change can reach its implementation and canonical decision yet leave public guidance, skills, templates, architecture, metadata, or implementation explanations stating the old rule. How does a methodology change keep the contract coherent without pretending the current judge can derive semantic carrier relationships?

## Decision outcome

**A methodology contract change inventories every live carrier that states the affected contract and updates that entire inventory in the same change.**

A live carrier is current material that instructs, constrains, describes, emits, or advertises the contract: permanent corpus truth, agent skills and their canonical sources, generated and scaffolded copies, public and contributor guidance, templates, implementation explanations, tests, CLI text, and distribution metadata. Pinned analyses, completed plans, changelog history, and decision context describing an earlier state remain history; a current reader-facing note may point from an older decision to its refinement without rewriting the old episode as though the new rule already existed.

The implementing change records its inventory in durable intent or evidence, repairs every inventoried carrier together, regenerates derived copies from canonical sources, and adds focused guards for stable claims that a machine can recognize. A change does not split one shared contract across merges when doing so would leave an intermediate accepted `main` that states the contract two ways.

The general inventory obligation remains agent-enforced. The guards introduced for one repair bind only the stable claims they name; they do not prove that an arbitrary future contract's carrier set is complete. A mechanical carrier registry needs evidence for its data model and is a later decision, not something inferred from this rule.

**Refinement of PDR-013:** the protected verifiable thread ends in acceptance evidence, not exclusively executable evidence: goal → plan → change → capability → acceptance criterion → acceptance evidence. For `Unit`, `Integration`, `E2E`, and `Performance`, that final edge is supported Go, JVM, or Cucumber evidence classified by proof type and direction, with `(single-direction)` as the explicit exception. For `Human`, the pull request acceptance brief is the evidence carrier and no code test is invented. A per-criterion `@draft` marks one promise as not yet proven without forcing its active criteria file or capability back to draft. Legacy criteria without a declared proof type retain the one-supported-reference contract. `clue validate` validates these declarations and references; normal test runners execute the tests.

This refinement changes what the protected thread connects and therefore takes effect only through the human merge boundary required by C-013. PDR-013's other core elements and red-line rule remain unchanged.

**Carrier:** C-006 and the shared decision-record guidance generated into every lifecycle skill carry the same-change inventory rule for agents; the repository and scaffold decisions READMEs carry it for humans and adopters; ARCH-003, C-013, and both AGENTS.md routing hubs carry the refined protected core; acceptance criteria and content guards hold the repaired evidence-model claims that are mechanically stable.

### Rejected: update only canonical or executable carriers

Canonical generation prevents byte drift among declared copies, and executable tests hold implemented behavior, but neither reaches independent public prose, architecture, metadata, or current design explanations. A green canonical subset can therefore continue advertising a contradictory contract.

### Rejected: require a machine-readable carrier registry now

The observed contradiction proves the need for a complete repair, not which fields would identify every future semantic claim or carrier. Designing a registry from one evidence-model incident would turn an agent-enforced obligation into speculative schema.

### Rejected: keep the core thread executable-only

Human-class criteria would remain accepted by the judge and lifecycle workflow while excluded by the protected definition those mechanisms serve. Calling the acceptance brief merely review context does not solve the contradiction: ADR-033 already makes that brief the criterion's proof.
