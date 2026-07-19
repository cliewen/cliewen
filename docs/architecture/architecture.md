---
id: ARCH-001
type: architecture
status: verified
links: [G-001]
title: System architecture — actors, lifetime classes, the frontmatter graph
---

# Architecture

## The four actors

Process cannot drive humans, and agents cheat unless mechanically prevented. The division of labor:

| Actor | Role | One-liner |
|---|---|---|
| **Skills** (`.agents/skills/clue-*`) | Process knowledge | Tell the agent what the next right step is |
| **CLI (`clue`)** | Deterministic judge | Tells everyone whether it was done right; also materializes the starting point (`clue init`, [ADR-018](../decisions/ADR-018-init-templates-embedded.md)) and regenerates the index blocks (`clue scaffold`, [ADR-019](../decisions/ADR-019-init-regenerates-indexes.md)) |
| **CI** | The wall | Refuses to proceed if not (same binary as local) |
| **Human** | Decision-maker | Settles what machines cannot check: meaning |

Goodhart guard: **machines enforce form, humans verify meaning.** The linter checks that AC-042 has a test; only PR review checks the test means anything.

```mermaid
flowchart TD
    A[Agent] --> Q{Cliewen<br/>relevant?}
    Q -->|no: plain| P[Ordinary branch<br/>focused checks]
    P --> H{Human PR review:<br/>change OK?}
    Q -->|yes| S[Skills<br/>process knowledge]
    S -->|guide| B[Branch = proposal<br/>light or full CH-xxx]
    B -->|clue validate| C{CLI: form OK?}
    C -->|no| A
    C -->|yes| H
    H -->|no| A
    H -->|merge = acceptance| D[/docs corpus<br/>system-of-record/]
    H -->|merge plain change| R[Repository state<br/>outside corpus]
    D -->|working memory| A
    W[CI wall<br/>same binary] -.enforces.- C
```

## Three artifact lifetime classes

1. **Permanent** — `/docs`. Lives forever, updated by every Cliewen merge.
2. **Transient** — `/changes/<CH-xxx>/` on a branch only. Dies at merge, digested into permanent docs. CI gate: `main` never contains `/changes/`.
3. **Campaign** — `/docs/plans`. Live on `main`, mutate continuously (bookkeeping in digests, semantic changes via ADR-backed revisions), frozen immutable at `status: completed` — never deleted.

Git is the engine: for a Cliewen change the branch is the proposal, the PR is the review gate, the merge commit is the acceptance, and `git log docs/` is the provenance archive. A plain change stays outside the artifact graph but retains the branch, PR, and human merge boundary ([PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md)). Repo-native, never forge-native.

## The frontmatter graph

Every artifact carries YAML frontmatter (`id`, `type`, `status`, `links`, `title` + small type-specific extensions). `clue` discovers artifacts by scanning frontmatter, never by path: **the ID is the identity, the path is only the current address.** Every field must have a consumer — a field neither `clue` nor a skill reads gets removed.

```mermaid
flowchart TD
    G[G-xxx goal] --> P[P-xxx / M-xxx plan]
    P --> CH[CH-xxx change]
    CH --> CAP[CAP-xxx capability]
    CAP --> AC[AC-xxx acceptance criterion]
    AC --> T[test tag<br/>positive + negative]
    C[C-xxx constraints] -. checked against every proposal .-> CH
    QS[QS-xxx quality scenarios] -. checked against every proposal .-> CH
```

## Deliberately out (doors defined, doors closed)

Deployment/operations (V3 door: production findings enter as new goals or constraints); external constraint catalogs (plug in via `source:`); kernel/profile layering (extracted after multiple working instances, not designed from zero). The `enforcement:` classes beyond `machine` door was opened by ADR-017: `agent`-enforced constraints are the lintable promotion backlog.
