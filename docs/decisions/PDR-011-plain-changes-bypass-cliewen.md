---
id: PDR-011
type: decision
status: verified
links: [PDR-002, PDR-007, C-002, C-005, C-012, AN-006]
title: Plain changes stay outside Cliewen while retaining human merge
author: agent
accepted-by: Flemming N. Larsen (2026-07-20, planning conversation)
---

# PDR-011 — Plain changes stay outside Cliewen

## Context and problem statement

PDR-002 removed the transient workspace from light changes but kept CH identity, proposal metadata, plan declaration, Cliewen verification, and full-repository checks. AN-006 shows that this still charges a guide-only editorial change for a methodology that produces no relevant evidence. Where should Cliewen stop?

## Decision outcome

**Cliewen owns light and full changes, not every repository edit. A third route, the plain change, stays outside Cliewen.**

A change is plain only when it has no effect on product behavior, intent, executable evidence, decisions, plans, policy, or methodology. `/docs`, `/changes`, product code, tests, configuration, build and release machinery, security and governance policy, AGENTS.md rules, skills, and lint rules are protected surfaces and cannot use the plain route. Changes to commands, contracts, user workflow, or normative instructions are not editorial. Uncertainty fails closed: classify the work as light or full before implementation.

Classify the requested work before loading the corpus. A plain change uses an ordinary branch rooted at accepted `main`, runs only checks relevant to the changed surface, opens a ready PR, and is merged by a human. It has no CH identity, plan or plan-less declaration, proposal artifact, corpus read, Cliewen skill, `clue validate` requirement, Cliewen verification checklist, plan bookkeeping, or Cliewen-mandated changelog entry.

Plain PRs do not consume the one-Cliewen-change-in-flight slot. They may proceed independently but never build on an unmerged branch. Agents still never push to `main`, merge their own PRs, or manufacture a local merge as acceptance.

PDR-002 continues to define the light/full split after a change enters Cliewen. Its global CH sequence covers those two tiers only. Editorial copy is not release history; a plain prose correction needs no changelog entry. If prose changes shipped behavior, a contract, a command, or a user workflow, it is not plain and the repository's release-note rule applies.

Repositories choose focused verification for their own surfaces. The generated corpus wall treats only Markdown outside protected corpus, policy, configuration, and methodology paths as eligible to pass through; every other diff fails closed. Cliewen's repository is narrower: only a guide-Markdown-only diff receives focused whitespace and production-guide checks. Mixed, protected, empty, unreadable, or unknown input uses the full CI path.

**Carrier:** the pre-corpus fast path in generated and repository routing hubs; the canonical change-tier and review-boundary skill sources and their generated standalone skills; the contributor and PR guidance; and repository-specific focused CI selection. The CLI and corpus format do not change.

### Rejected: direct commits to main

The expensive work was Cliewen metadata and irrelevant verification, not the ordinary protected-repository review boundary. Direct pushes make an agent's self-classification unreviewable.

### Rejected: keeping CH identity for provenance

A CH number says Cliewen participated. Assigning one when no corpus, plan, decision, capability, or evidence is involved creates bookkeeping without traceability.

### Rejected: path patterns as the definition

Paths are useful for fail-closed automation but cannot decide meaning. A one-line guide edit can change a command or policy, while a large copy edit may preserve meaning. Protected paths can disqualify a plain change; an allowed path cannot prove one.

### Rejected: full verification after skipping metadata

Running unrelated builds and corpus checks preserves most of the time and resource cost while producing no evidence about the edited surface.
