---
id: ADR-062
type: decision
status: verified
links: [G-011, P-019, M-081, ADR-013, ADR-011, ADR-048, ADR-044, CAP-001, CAP-002]
title: A repository declares its Cliewen role as machine state, and an adopter-binding rule names a shipped carrier
binds: adopter
author: agent
accepted-by: Flemming N. Larsen (2026-09-02)
---

# ADR-062 — Repository role is declared machine state

## Context

[ADR-013](ADR-013-ships-generic-vs-repo-local.md) separates the generic shipped surface from the one repo-local layer, but nothing in a repository says which kind of repository it is. Cliewen's own repository and every adopter carry the same `docs/`, `.clue/`, and `.agents/skills/`; the distinction survives only as the incidental presence of `cmd/clue` and as prose in a routing hub that states shared and local rules together. Tooling that needs the distinction has to infer it, an agent reading the hub cannot tell which rules bind, and no check notices when a rule that binds adopter behaviour never reaches the surface adopters receive.

## Decision

**A repository declares its role in `.clue/role.yaml`, and that declaration is machine state rather than a second configuration layer.** The file carries `role: source` or `role: adopter` and nothing else. `clue init` materializes `adopter` without overwriting an existing declaration, and an undeclared repository is an adopter: every repository onboarded before the marker existed is one, and none may be blocked for lacking a file it had no way to write.

**The marker records the role alone.** A repository's Cliewen version already has two canonical carriers — the `version:` stamp each generated `skill.md` holds and the CI caller's `clue-version` input ([ADR-011](ADR-011-version-stamping.md)) — and a third copy would be precisely the drift ADR-013 warns about.

**`clue validate` applies a source-repository-only rule that a decision or constraint binding adopter behaviour names a carrier on the shipped surface**, and never applies that rule to an adopter's corpus. The boundary ADR-013 states in prose becomes checkable in the one repository that can be checked against the templates and canonical skill sources.

This amends ADR-013's rejection of a second config file. That rejection stands for what it addressed: a per-repository layer of rules and forked skills that would split or drift the source of truth. A role marker configures nothing, overrides no rule, and forks no skill; it records a fact the tooling would otherwise guess. It sits beside the identity ledger under `.clue/` for the reason [ADR-048](ADR-048-identity-ledger.md) put the ledger there — derived operational state is not authored corpus prose.

The corpus structure this affects is described in [`docs/architecture/README.md`](../architecture/README.md).

## Rejected: infer the role from the tree

Deriving the role from the presence of `cmd/clue` and `internal/skills/source/` needs no new file, but the inference is wrong in a fork, a monorepo holding Cliewen beside other projects, and any vendored checkout. A rule that changes behaviour by repository kind should not rest on a directory listing that a legitimate layout can invalidate.

## Carrier

The `internal/role` package and its marker, `clue init`'s materialization, the source-only carrier rule in `internal/corpus`, the role-reading instruction in `internal/skills/source/shared/local-conventions.md.tmpl` and the `MIG-012` step in `internal/skills/source/skills/clue-upgrade.md.tmpl`, the scaffolded corpus guidance in `internal/scaffold/templates/docs/README.md`, this repository's own separated routing text, and the adopter-facing changelog entry.
