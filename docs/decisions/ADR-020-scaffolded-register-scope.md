---
id: ADR-020
type: decision
status: inferred
links: [ADR-013, ADR-017, CAP-001]
title: The scaffolded register seeds only conventions without a versioned carrier
author: agent
accepted-by:
---

# ADR-020 — Scope of the scaffolded constraint register

## Context and problem statement

[ADR-017](ADR-017-conventions-are-constraints.md) makes the constraints folder the register of every prose-only convention. `clue init` scaffolds that register into adopting repos — which raises the scope question there: the methodology arrives with rules of its own, some carried by the versioned skills, some declared only by the generated AGENTS.md. Registering everything duplicates the skills' content into unversioned artifacts; registering nothing leaves prose-only rules uninventoried from the first commit.

The scoping first entered the decision log (row of 2026-07-17) as a template choice judged cheap to reverse. This record promotes that row, citing its date, per the register's promotion practice: the choice turned out to constrain future changes — it fixes what `clue init` may ever seed and how the register doctrine applies in every scaffolded repo, so reversing it is a doctrine change across adopters, not a local template edit.

## Considered options

1. **Seed a constraint for every methodology rule** — a complete inventory, but each skill-carried rule gains a second, unversioned carrier that drifts the moment a skill upgrade lands.
2. **Seed nothing** — no duplication, but a rule only the generated AGENTS.md declares is prose-only from the first commit, which is exactly the failure the register exists to prevent.
3. **Seed exactly the methodology conventions no versioned skill carries.**

## Decision outcome

**Option 3.** The register a scaffolded repo starts with holds the methodology conventions no versioned skill carries; the skills are their own carrier and version, so their rules are never duplicated into the register. The repo's own conventions register as they are declared, per ADR-017. Consequences:

- Seeded constraints are **self-sourced** (`source:` names the Cliewen methodology as scaffolded by `clue init`), so they stand whether the repo keeps the generated AGENTS.md or its own — the register entry is the rule's authoritative declaration either way; the generated AGENTS.md mirrors it as readable prose.
- The markdown hard-wrap prohibition is such a convention and seeds as the register's first entry.
- Seeded constraints live the ordinary register lifecycle from there (ADR-017): an `agent`-enforced seed states its promotion trigger, promotion flips `enforcement:` to `machine`, and the entry stays in the register as inventory. The seeding boundary governs what `init` emits, never what a repo later registers or promotes.
