# The corpus

This directory is the **system-of-record**: the permanent, durable truth about the system. Cliewen changes are transient deltas on branches that get **digested** into this corpus at merge — `git log docs/` is the audit trail. Entry point for humans and agents alike; agents treat this tree as working memory when a change affects product or methodology meaning. Plain changes classified by AGENTS.md stay outside the corpus.

## How the corpus is wired

Every artifact carries YAML frontmatter with a common core — `id`, `type`, `status`, `links`, `title` — plus small type-specific extensions. **Identity is the ID, the path is only the current address**: tooling discovers artifacts by scanning frontmatter, and external systems reference IDs, never paths.

The red thread `clue validate` walks:

```
G-xxx (goal) → P-xxx/M-xxx (plan/milestone) → CH-xxx (change)
  → CAP-xxx (capability) → AC-xxx (acceptance criterion) → test tag
```

Cross-cutting, checked against every proposal: C-xxx (constraints), QS-xxx (quality scenarios).

## What lives where — and when a change updates it

Each folder below holds one kind of record. A Cliewen change (the `clue-delta` loop) updates, in the same branch and PR as the code, every record its work touches:

- **Goals** (`goals/`) — who wants the system and why. A new wish enters here as `status: proposed`; a change rarely touches goals.
- **Plans** (`plans/`) — campaigns with verifiable milestones. Every Cliewen change names the plan item it serves (or declares itself plan-less); the digest updates milestone bookkeeping.
- **Capabilities** (`capabilities/`) — one folder per capability: `README.md` (what and why), `criteria.md` (acceptance criteria as Gherkin, each tied to tests), `design.md` (how it works). **Design is documented per capability** — a change that alters a capability's behavior updates its criteria and design in the same PR.
- **Architecture** (`architecture/`) — the shape of the whole: the expensive-to-change. Updated when a change alters the system's structure or public surface, not for local detail.
- **Decisions** (`decisions/`) — why things are the way they are. An **ADR** records an architectural decision, a **PDR** a decision about how the project works; both are expensive to reverse. Cheap-and-local-to-reverse decisions are one-line rows in `log.md`. Every decision made during a change is recorded in its digest.
- **Constraints** (`constraints/`) — external rules the system must not break: laws, licenses, policies. Each names its `source` and how it is `enforcement`-checked. Updated only when the outside world imposes something.
- **Quality** (`quality/`) — quality scenarios: verifiable non-functional requirements. Updated when a change introduces or moves a quality bar.
- **Analysis** (`analysis/`) — findings from spikes and extractions. Historical records: written once, never rewritten.

## Status vocabularies

| Type | Statuses |
|---|---|
| goal | `proposed` → `accepted` → `retired` |
| plan | `draft` → `active` → `completed` (completed ⇒ immutable) |
| capability, criteria, design | `draft` → `active` → `retired` |
| decision | `inferred` → `verified` (human acceptance promotes) |
| log | `active` |
| constraint | `active` → `retired` |
| quality | `draft` → `active` → `retired` |
| analysis | `verified` |
| architecture | `draft` → `verified` |

## Folders

<!-- clue:index:start -->
- [goals/](goals/README.md) — G-xxx: who wants it, why
- [plans/](plans/README.md) — P-xxx: campaigns and milestones
- [capabilities/](capabilities/README.md) — CAP-xxx: one folder per capability (README / criteria / design)
- [architecture/](architecture/README.md) — the whole, the expensive-to-change
- [decisions/](decisions/README.md) — ADR-xxx, PDR-xxx and the decision log
- [constraints/](constraints/README.md) — C-xxx: laws, licenses, policies you must not break
- [quality/](quality/README.md) — QS-xxx: quality scenarios (verifiable NFRs)
- [analysis/](analysis/README.md) — spike findings, extraction reports
<!-- clue:index:end -->
