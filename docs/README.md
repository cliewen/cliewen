# The Cliewen corpus

This directory is the **system-of-record**: the permanent, durable truth about the system. Full Cliewen changes are transient deltas on branches that get **digested** into this corpus at a human-controlled merge commit; reachable Git history is the full provenance archive, and `git log docs/` remains the durable-corpus audit trail. This index is the orientation point when a request names no artifact; once an identity is known, `clue context <id>` emits its focused outgoing-link slice and agents read further only when the task discovers another edge. Plain changes under PDR-011 stay outside the corpus.

## How the corpus is wired

Every artifact carries YAML frontmatter with a common core — `id`, `type`, `status`, `links`, `title` — plus small type-specific extensions. **Identity is the ID, the path is only the current address**: tooling discovers artifacts by scanning frontmatter, and external systems reference IDs, never paths. Status lives in frontmatter, never in folder names; status views are generated.

The red thread the linter walks:

```
G-xxx (goal) → P-xxx/M-xxx (plan/milestone) → CH-xxx (change)
  → CAP-xxx (capability) → AC-xxx (acceptance criterion) → acceptance evidence
    → classified Go/JVM/Cucumber test reference, or Human acceptance brief
```

Cross-cutting, checked against every proposal: C-xxx (constraints, including verifiable quality bars — see [ADR-027](decisions/ADR-027-quality-scenarios-are-constraints.md)).

When a released `clue` adds or narrows a corpus obligation, preview `clue migrate` and apply its complete plan only after resolving any semantic or local-edit findings. `clue init` remains a non-destructive materializer, not an updater.

Small consumed extensions keep the graph explicit. `ac-prefix` on a criteria.md namespaces its AC IDs (`<PREFIX>-<digits>[lowercase-suffix]`, where `<PREFIX>` may contain uppercase hyphen-separated segments; default `AC`; consumer `checkACTests`). Declarations and links use the exact canonical spelling; supported evidence carriers document their limited underscore or hyphen-removal aliases. `provenance: inferred|verified` marks agent-extracted non-decisions, and an inferred artifact also declares `reversal-cost: low|high`; high-cost inferred meaning blocks an active capability joined to it by one links edge, while decisions carry provenance in `status` and remain visible separately (consumer `checkProvenance`). `supersedes: [ID, …]` names an artifact this one's retirement replaced and requires that ID to be absent (consumer `checkSupersedes`). `reality: contradicted` marks an incident analysis whose links name the capability or criterion that reality disproved (consumer `checkReality` and `clue validate --reality-gaps`). `carried-by: [ID, …]` on an analysis names the durable artifacts its findings reached, which is what lets a spike be reported as spent once every plan it serves is complete and no live decision or constraint still cites it (consumer `checkCarriedBy` and `corpus.SpentAnalyses`; [PDR-052](decisions/PDR-052-a-spent-analysis-is-reported-not-retained.md) and [PDR-053](decisions/PDR-053-cited-evidence-stays-readable.md)). `binds: adopter|repo` on a decision or constraint names whose behaviour it governs; in this repository alone, an adopter-binding record must cite a carrier under `internal/skills/source/` or `internal/scaffold/templates/` (consumer `checkBoundary`; [ADR-062](decisions/ADR-062-repository-role-is-declared-machine-state.md)).

**This repository is Cliewen's source, and `.clue/role.yaml` says so.** An adopter's marker reads `role: adopter`, and a repository without one is an adopter. The distinction is not cosmetic: rules about releases, generated skills, and the shipped surface exist here and reach no adopter, so a rule stated only in this corpus never binds anyone else. That is what `binds: adopter` and its carrier check are for.

## Status vocabularies

**The default lifecycle is `draft` → `active`.** It applies to every artifact type — including adopter-defined types ([ADR-026](decisions/ADR-026-adopter-types-default-lifecycle.md)) — except the few below that need a different shape for a stated reason ([ADR-025](decisions/ADR-025-one-status-lifecycle.md)). There is no `retired`: retiring a default-lifecycle artifact means deleting its file in the same change and naming it in a successor's `supersedes:` field, never holding a status a file rests in ([ADR-034](decisions/ADR-034-retirement-is-deletion.md)). This table mirrors the `defaultLifecycle` slice and `statusVocabExceptions` map in `internal/corpus/rules.go` — the consumer that enforces it. Change them together.

| Type | Statuses | Why not the default |
|---|---|---|
| goal | `proposed` → `accepted` | proposed goals ARE the inbox (ADR-002) |
| plan | `draft` → `active` → `completed` | `completed` is immutable, not `retired` (C-008) |
| decision | `inferred` → `verified` | provenance lives in status; human acceptance promotes (ADR-010) |
| change, tasks | `open` | transient workspace artifacts |
| open-questions | `open` → `resolved` | transient workspace artifacts |
| imported-change | `in-progress` → `complete` | durable, never `retired` — the record survives its extracted source (ADR-050) |

Types on the default: capability, criteria, design, constraint, architecture, analysis, and any type an adopter adds.

## Folders

<!-- clue:index:start -->
- [goals/](goals/README.md) — G-xxx: who wants it, why (the inbox lives here as `status: proposed`)
- [plans/](plans/README.md) — P-xxx: campaign layer; flat, status in frontmatter
- [capabilities/](capabilities/README.md) — CAP-xxx: one folder per capability (README / criteria / design)
- [architecture/](architecture/README.md) — system scope: the whole, the expensive-to-change
- [design/](design/README.md) — cross-cutting behaviour: flows, interactions, and shared patterns
- [decisions/](decisions/README.md) — future-shaping choices by subject: ADR-xxx architecture, PDR-xxx project/process, IDR-xxx implementation
- [constraints/](constraints/README.md) — C-xxx: laws, licenses, policies, and verifiable quality bars you must not break
- [analysis/](analysis/README.md) — spike findings, extraction reports
- [imported-changes/](imported-changes/README.md) — IC-xxx: durable records of in-flight source work brownfield extraction preserves
- [use-cases/](use-cases/README.md)
- [VIS-001 — Cliewen — durable intent that a machine can check before a human accepts it](vision.md) · `draft` — **What this is.** A methodology and one small command-line judge that keep a repository's durable intent — what is wanted, what the system can do, what proves it — in the repository itself, wired…
<!-- clue:index:end -->
