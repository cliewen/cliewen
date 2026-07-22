# The corpus

The `/docs` tree is Cliewen's permanent working memory. Agents read it before acting, people review it with the implementation, and Git records every accepted mutation.

## The taxonomy

| Folder | Artifact | Question it answers |
|---|---|---|
| `goals/` | `G-xxx` | Who needs an outcome, and why? |
| `plans/` | `P-xxx` with `M-xxx` milestones | What bounded campaign moves a goal forward? |
| `capabilities/` | `CAP-xxx` with criteria and design | What can the system do, how is it verified, and how is it built? |
| `architecture/` | `ARCH-xxx` | What describes the whole system or an expensive-to-change boundary? |
| `decisions/` | `ADR-xxx`, `PDR-xxx`, and a decision log | Why is the system or project shaped this way? |
| `constraints/` | `C-xxx` | What rule must every relevant change obey? |
| `quality/` | `QS-xxx` | What measurable non-functional behavior must hold? |
| `analysis/` | `AN-xxx` | What did a time-boxed investigation find? |

Each folder has a README that explains its type and contains a generated index of the artifacts beside it.

## Identity is not location

Every artifact begins with YAML frontmatter:

```yaml
---
id: CAP-002
type: capability
status: active
links: [G-001]
title: clue validate
goal: G-001
---
```

The ID is the identity; the path is only its current address. `clue` scans frontmatter, checks IDs and status vocabularies, resolves every `links` entry, and verifies that generated indexes match the files on disk.

This makes refactoring the corpus safe. A file can move without becoming a different decision or capability, while duplicate IDs and broken references fail loudly.

## One home per scope

System-wide and expensive-to-change design belongs under `architecture/`. Per-capability design lives beside the capability. Decisions explain durable choices but do not become substitute design documents. Findings record what an investigation observed but do not silently become accepted intent.

The separation is intentionally strict: a fact with two homes will eventually disagree with itself.

## Choose the right decision record

Start with the cost of reversing the decision, then ask what it changes:

| Decision | Record |
|---|---|
| Cheap and local to reverse | A dated row in `docs/decisions/log.md` |
| Expensive to reverse and about software architecture or the corpus format, such as frontmatter fields or extraction mappings | An ADR, or Architectural Decision Record |
| Expensive to reverse and about project workflow or process | A PDR, or Project/Process Decision Record |

The log is for choices that do not need a full argument preserved. ADRs and PDRs use the same decision template because expensive choices need context, alternatives, and rationale; the different names tell readers whether the decision shapes the system or the way the project works.

## See a living corpus

Cliewen dogfoods the methodology. Browse its [corpus entry point](https://github.com/cliewen/cliewen/blob/main/docs/README.md), [active public campaign](https://github.com/cliewen/cliewen/blob/main/docs/plans/P-003-goes-public.md), or [validator capability](https://github.com/cliewen/cliewen/tree/main/docs/capabilities/CAP-002-validate) to see real artifacts rather than a toy example.

## Next

[Follow one proposal through the change loop.](./change-loop)
