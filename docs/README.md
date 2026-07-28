# The Cliewen corpus

This directory is the **system-of-record**: the permanent, durable truth about the system. Cliewen changes are transient deltas on branches that get **digested** into this corpus at merge — `git log docs/` is the audit trail. This index is the orientation point when a request names no artifact; once an identity is known, `clue context <id>` emits its focused outgoing-link slice and agents read further only when the task discovers another edge. Plain changes under PDR-011 stay outside the corpus.

## How the corpus is wired

Every artifact carries YAML frontmatter with a common core — `id`, `type`, `status`, `links`, `title` — plus small type-specific extensions. **Identity is the ID, the path is only the current address**: tooling discovers artifacts by scanning frontmatter, and external systems reference IDs, never paths. Status lives in frontmatter, never in folder names; status views are generated.

The red thread the linter walks:

```
G-xxx (goal) → P-xxx/M-xxx (plan/milestone) → CH-xxx (change)
  → CAP-xxx (capability) → AC-xxx (acceptance criterion) → test tag
```

Cross-cutting, checked against every proposal: C-xxx (constraints, including verifiable quality bars — see [ADR-027](decisions/ADR-027-quality-scenarios-are-constraints.md)).

Small consumed extensions keep the graph explicit. `ac-prefix` on a criteria.md namespaces its AC IDs (`<PREFIX>-<digits>`, default `AC`; consumer `checkACTests`). `provenance: inferred|verified` marks agent-extracted non-decisions, and an inferred artifact also declares `reversal-cost: low|high`; high-cost inferred meaning blocks an active capability joined to it by one links edge, while decisions carry provenance in `status` and remain visible separately (consumer `checkProvenance`). `supersedes: [ID, …]` names an artifact this one's retirement replaced and requires that ID to be absent (consumer `checkSupersedes`). `reality: contradicted` marks an incident analysis whose links name the capability or criterion that reality disproved (consumer `checkReality` and `clue validate --reality-gaps`).

## Status vocabularies

**The default lifecycle is `draft` → `active`.** It applies to every artifact type — including adopter-defined types ([ADR-026](decisions/ADR-026-adopter-types-default-lifecycle.md)) — except the few below that need a different shape for a stated reason ([ADR-025](decisions/ADR-025-one-status-lifecycle.md)). There is no `retired`: retiring a default-lifecycle artifact means deleting its file in the same change and naming it in a successor's `supersedes:` field, never holding a status a file rests in ([ADR-034](decisions/ADR-034-retirement-is-deletion.md)). This table mirrors the `defaultLifecycle` slice and `statusVocabExceptions` map in `internal/corpus/rules.go` — the consumer that enforces it. Change them together.

| Type | Statuses | Why not the default |
|---|---|---|
| goal | `proposed` → `accepted` | proposed goals ARE the inbox (ADR-002) |
| plan | `draft` → `active` → `completed` | `completed` is immutable, not `retired` (C-008) |
| decision | `inferred` → `verified` | provenance lives in status; human acceptance promotes (ADR-010) |
| log | `active` | one register — rows are its lifecycle (PDR-003) |
| change, tasks | `open` | transient workspace artifacts |
| open-questions | `open` → `resolved` | transient workspace artifacts |

Types on the default: capability, criteria, design, constraint, architecture, analysis, and any type an adopter adds.

## Folders

<!-- clue:index:start -->
- [goals/](goals/README.md) — G-xxx: who wants it, why (the inbox lives here as `status: proposed`)
- [plans/](plans/README.md) — P-xxx: campaign layer; flat, status in frontmatter
- [capabilities/](capabilities/README.md) — CAP-xxx: one folder per capability (README / criteria / design)
- [architecture/](architecture/README.md) — system scope: the whole, the expensive-to-change
- [decisions/](decisions/README.md) — ADR-xxx (architecture) and PDR-xxx (project/process): MADR + provenance, including rejected records
- [constraints/](constraints/README.md) — C-xxx: laws, licenses, policies, and verifiable quality bars you must not break
- [analysis/](analysis/README.md) — spike findings, extraction reports
<!-- clue:index:end -->
