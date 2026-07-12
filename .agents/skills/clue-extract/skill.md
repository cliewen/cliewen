# clue-extract

Brownfield adoption: transform an existing repository's spec corpus into a Cliewen `/docs` corpus (ADR-008). Use once per adopted repo; the run is that repo's first change loop (`clue-delta` applies — branch, proposal, digest, PR).

## Target contract (source-independent)

The extraction PR is complete only when all of these hold:

1. **The full taxonomy exists**: `/docs` with goals, plans, capabilities (README / criteria / design per folder), decisions, constraints, quality, analysis, architecture — every folder with an indexed README. Extract meaning, don't invent it: a folder with nothing real to hold stays empty-but-indexed.
2. **Everything extracted is born `provenance: inferred`** (ADR-010). Decisions instead use `status: inferred` with `author: agent`. Human review of the PR promotes to `verified` — file-by-file or in bulk; never pre-promote.
3. **Existing AC IDs survive** (ADR-009): declare each capability's namespace with `ac-prefix:` in its criteria.md and keep the source's IDs verbatim. Renumbering is forbidden — IDs are meaning-immutable (ADR-007) and existing test tags must keep resolving.
4. **Every test keeps or gains exactly one purpose** (ADR-006): tests already tagged with an AC ID are done; tests without one get `Unit`, `Sanity` or `Arch` per their actual intent. Where a JVM test framework is present, install an ArchUnit (or equivalent) rule enforcing one purpose tag per test — clue only harvests at file level (ADR-009).
5. **`clue validate` is green** against the repo root before the PR opens. The extracted corpus is judged by exactly the same rules as a greenfield one.
6. **The source corpus dies in the same PR**: spec trees, parallel registries and source-format skills are deleted — git history is their archive. Two systems of record is zero systems of record.
7. **Routing is rewritten**: the repo's AGENTS.md points agents at `/docs/README.md` and the `clue-*` skills; the Cliewen skills (`clue-analysis`, `clue-plan`, `clue-delta`, `clue-verify`, `clue-extract`) are installed under `.agents/skills/`.
8. **An extraction report lands in `/docs/analysis`** (AN-xxx): what was found, what mapped where, what was dropped and why — analysis must leave corpses.
9. **Unsolved adoption items become named doors in the repo's plan** (e.g. clue binary distribution for CI), not silent gaps.

## Source mapping: OpenSpec

Layout: `openspec/config.yaml`, synced truth in `openspec/specs/<capability>/spec.md`, pending work in `openspec/changes/<name>/` (proposal.md, design.md, tasks.md, spec deltas), applied work in `openspec/changes/archive/`.

| OpenSpec | Cliewen |
|---|---|
| `specs/<cap>/spec.md` | `docs/capabilities/CAP-xxx-<cap>/` — one capability per spec file |
| Spec header + purpose prose | capability `README.md` (what/why, `goal:` link) |
| `### Requirement:` + SHALL + `#### Scenario: <name> [ID]` | Gherkin scenarios in `criteria.md`, tag line `@<ID>` — keep the ID, whatever bracket/backtick notation the source used |
| `Test-type:` line per scenario | plain body text inside the Gherkin scenario (required-types enforcement is a door, ADR-006) |
| Scenario ID prefix (e.g. `MG`) | `ac-prefix:` in that criteria.md frontmatter; a delta spanning several prefixes splits into one capability per prefix |
| Pending change (`changes/<name>/`) | a milestone in the repo's plan **plus** a `status: draft` capability holding its criteria (draft = exempt from the test contract until implemented) and its design decisions in `design.md`; its tasks.md dies — `clue-delta` regenerates tasks when implementation starts |
| `changes/archive/…` | git history only — no corpus artifact |
| Nygard/MADR ADRs in `docs/decisions` | `ADR-xxx` born `status: inferred`, `author: agent`; original acceptance dates preserved in the body |
| Architecture docs | `docs/architecture/` artifacts (`status: draft` until reviewed) or capability `design.md` where they are capability-local |
| AC registry / scenario templates (`test/…`) | deleted — the corpus is the registry; next free ID per prefix is max + 1 over declared ACs |
| Project README purpose statements | `G-xxx` goal(s), `status: accepted` (the repo's existence is the acceptance) |
| Coverage/quality gates in build config | `QS-xxx` quality scenarios referencing the enforcing tool |
| OpenSpec workflow skills (`openspec-*`) | deleted with the source corpus |
| JUnit `@Tag("XX_NNN")` | untouched — clue normalizes underscores to hyphens at harvest (ADR-009) |

Watch for: the same logical ID written three ways (`[MG-010]`, `` `PG-001` ``, `MG_010`); `## ADDED/MODIFIED Requirements` delta headers in pending changes (apply the delta meaning, don't copy the header); scenario WHEN/THEN bullets mapping to Gherkin When/Then/And.

## What this skill never does

Invent requirements the source doesn't state; renumber or rename IDs; leave the source corpus alive "for reference"; promote its own output to `verified`; touch test code beyond adding missing purpose tags and the purpose-enforcement rule.
