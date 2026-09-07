---
id: ADR-067
type: decision
status: inferred
links: [G-012, P-021, M-090, CAP-001, CAP-002, CAP-009, C-023, ADR-044, PDR-011]
title: A corpus without a vision stays valid, and no migration writes one
binds: adopter
author: agent
accepted-by: []
---

# ADR-067 — A corpus without a vision stays valid

## Context

Every repository already running Cliewen lacks a vision, because there was nothing to lack. Making the artifact required would turn those corpora red for an obligation their authors could not have met, and would attach the failure to work that has nothing to do with product meaning.

The opposite failure is as real. A concept that is never surfaced is a concept nobody adopts, and a repository that quietly has no statement of direction is indistinguishable from one that decided it does not need one.

## Decision

**A missing vision is never a validation failure, and a missing use case never is either.** `clue validate` judges the form of what exists and requires nothing to exist. It has no coverage rule for use cases and reports no coverage figure for them, because a completeness metric over an optional artifact rewards creating artifacts nobody needs.

**Three surfaces make the absence visible without blocking anyone.** `clue init` writes a marked bootstrap into a *new* repository, which validation rejects until it is replaced — the same treatment the architecture and design overviews already get, and the reason a repository that starts today starts with a direction. `clue migrate` reports a missing vision as a notice and writes no vision content. `clue validate --intent` states what exists whenever anyone asks.

**Migration never drafts a vision.** A repository proves what a system does and cannot prove why anyone wanted it; a migration that produced one would be inventing the single thing in the corpus that has no evidence base. Structure that carries no meaning is a different matter and is written — the optional use-case folder and its index row are inert scaffolding.

**Full work that changes meaning discloses the vision it proceeds under**, as a required line in the acceptance brief ([C-023](../constraints/C-023-full-work-discloses-its-vision.md)). This is what separates *not yet established* from *accidentally omitted*: a repository that has decided it does not need one says so, once, where a human reads it, rather than being asked again by a tool that cannot tell the two apart.

## Rejected: have `clue migrate` write the same bootstrap `clue init` writes

It is one line of code and it fails the case it exists to serve. The bootstrap is rejected by validation, so a migration that wrote one would turn a green adopter corpus red in the same run that was supposed to move it forward — for a file the adopter did not ask for and cannot fill in mechanically.

## Carrier

`MIG-014` in `internal/migrate`, the `--intent` report, the bootstrap in `internal/scaffold/templates/docs/vision.md`, C-023 and the brief line in `internal/scaffold/templates/github/pull_request_template.md`, the migration walk in `internal/skills/source/skills/clue-upgrade.md.tmpl`, and the adopter-facing changelog entry.
