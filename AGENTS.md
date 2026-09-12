# Agent routing hub

This repository dogfoods Cliewen and declares `role: source` in `.clue/role.yaml`. The shared methodology below also binds adopters; source-repository conventions apply only here. Before applying a rule that differs by repository kind, read the role marker rather than infer it from the checkout. An adopter-binding rule belongs on a shipped carrier under `internal/skills/source/` or `internal/scaffold/templates/`; `clue validate` enforces that boundary ([ADR-062](docs/decisions/ADR-062-repository-role-is-declared-machine-state.md)).

## Always first

Before your first tool call, including for a question or review, run `clue latest --quiet`. Route a non-empty result, or an unknown `latest` command, to [`clue-upgrade`](.agents/skills/clue-upgrade/skill.md); the human alone decides whether to upgrade. Run this network check unprompted only for that purpose, never as a validation verdict or required check. Ordinary `clue` workflow commands also report an available update.

After that check, once per fresh agent context, run `clue next --all` before substantive work. If the opening request leaves direction open, briefly state the repository position, read the leading candidate with `clue context`, and recommend an option without starting it. If existing work materially affects a concrete request, mention it; otherwise do not add a routine status preamble. Draft milestones and proposed goals are choices for human review, never authorization to begin.

## Before editing: route the work

Inspect the smallest relevant context and tell the user `Recommended route: simple` or `Recommended route: full`, why, and what discovery would change that recommendation.

| Route | Use when | Required work |
|---|---|---|
| Simple | The accepted contract is unchanged: named-consumer analysis, a correction that restores an unchanged criterion, regression evidence, in-contract configuration, refactoring, maintenance, or editorial work. | No CH identity, workspace, plan declaration, digest, acceptance brief, or mandatory agentic review; run checks relevant to changed surfaces. |
| Full | Acceptance-criterion, capability, decision, policy, plan-promise, methodology, or uncovered-behavior meaning changes. Uncertainty is full. | Read and use [`clue-delta`](.agents/skills/clue-delta/skill.md)'s full loop with `/changes/CH-xxx-slug/`. |

Paths and diff size may warn but do not decide the route. Reassess after semantic discovery and against the complete diff; if simple work becomes full, pause and recommend full. If the user chooses simple anyway, keep the repository truthful and add `Cliewen-Route: simple`, `Cliewen-Recommendation: full`, and `Cliewen-Override: user chose simple; <concise risk>` to the final authored commit.

A route does not authorize a push. Push directly to an integration branch only with explicit user authorization and repository permission; humans may integrate as the repository permits, and local conventions may be stricter. Release is not a Cliewen route.

## Work from durable context

`clue next --all` reports open changes to resume, active milestones, draft milestones, and proposed goals in decreasing order of authority. For a full change with a known identity, run `clue context <id>` and read its bounded slice. Otherwise read [`docs/README.md`](docs/README.md), choose the closest artifact, then run `clue context`; `/docs` is the system of record and working memory.

Assess documentation impact before closing every change. Keep `docs/architecture/README.md` current for system structure, `docs/design/README.md` current for cross-cutting behavior, and capability `design.md` for local detail; add or update only information that answers a reader question without duplicating existing material. Draft missing overviews from evidence, ask the user when a material boundary or intent is unclear, use Mermaid when it improves review, retain SVG when it does not, and state in the change or pull-request handoff what durable documentation changed or why none was needed.

For intent work, use [`clue-plan`](.agents/skills/clue-plan/skill.md) or [`clue-extract`](.agents/skills/clue-extract/skill.md): `VIS-001` → goal → optional `UC-xxx` → capability → criterion → evidence states product meaning, while goal → plan → milestone → change → accepted merge states delivery. Links point down; vision and use cases are optional and change only when durable meaning does. A missing vision is elicited (greenfield) or inferred as cited `status: draft`, `provenance: inferred` meaning until a human confirms it; a full change's acceptance brief states its vision or that none exists.

Markdown prose is never hard-wrapped: one line per paragraph and list item; line breaks are structural only ([C-001](docs/constraints/C-001-no-hard-wrapped-markdown.md)).

The core is behind a red line: the verifiable thread (goal → plan → change → capability → criterion → acceptance evidence, including classified executable references and genuine Human proof in the acceptance brief), the full-loop human acceptance boundary, and `clue validate` as deterministic judge ([ARCH-003](docs/architecture/core.md)). Changing what any core element means changes the accepted contract: recommend full and record an explicit decision; the user may choose simple with the route-override trailers above ([C-013](docs/constraints/C-013-core-changes-need-decision.md)). Periphery never constrains the core.

## Read the matching skill

| Skill | Read when |
|---|---|
| [`clue-analysis`](.agents/skills/clue-analysis/skill.md) | Risks or unknowns need a spike that ends in findings documentation. |
| [`clue-plan`](.agents/skills/clue-plan/skill.md) | Creating or revising a plan. |
| [`clue-upgrade`](.agents/skills/clue-upgrade/skill.md) | Checking for or, with human approval, performing a coordinated upgrade; it is simple because the release's contract was accepted upstream ([PDR-043](docs/decisions/PDR-043-upgrade-routes-as-simple-work.md)). |
| [`clue-delta`](.agents/skills/clue-delta/skill.md) | Running a full change from proposal through implementation, digest, verification, and human-controlled merge. |
| [`clue-extract`](.agents/skills/clue-extract/skill.md) | Adopting Cliewen in a brownfield repository by turning one existing corpus into durable Cliewen records. |
| [`clue-verify`](.agents/skills/clue-verify/skill.md) | Verifying a full change and running its bounded agentic review before readiness. |

## Source-repository conventions

These rules govern Cliewen's source repository only; they never become an adopter requirement ([ADR-013](docs/decisions/ADR-013-ships-generic-vs-repo-local.md)).

| When touching or doing | Required action |
|---|---|
| Verifying a change | Read and run the applicable block under **Verify Locally** in [`CONTRIBUTING.md`](CONTRIBUTING.md) verbatim. Its whitespace check compares the branch with its base, unlike `git diff --check`, which can miss an already-committed defect. |
| Integrating agent-authored work | Use a pull request and human merge. The repository's administrative version cut is a local simple-work specialization with its exact release surface and focused checks, not an adopter route. |
| `internal/skills/source/**` | Edit the canonical source, run `go generate ./internal/skills`, and never edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly. Before integration inspect affected generated entry points and references as an adopter receives them; remove exact duplication, simplify repeated state branches without losing outcomes, preserve protected obligations/examples/exceptions, avoid unrelated rewrites, and state in the handoff what improved or that it was already clear and DRY. |
| Prose in `/guide` | Use the installed `humanizer` skill. If unavailable, ask the user to install it in the coding agent's local user directory. Write for software practitioners in plain language, explaining a technical term at first use or linking to a short glossary definition. |
| A shipped surface | Add a user-facing `[Unreleased]` entry to [`CHANGELOG.md`](CHANGELOG.md) when behavior of `clue`, generated-skill text, or an artifact materialized by `clue init` or `clue scaffold` changes. Name the adopter-facing surface, not this repository's equivalent; local corpus, contributor guidance, and conventions alone, and simple editorial corrections, owe none. A mixed change notes only its shipped portion; `.github/pull_request_template.md` and `.github/workflows/clue-validation.yml` are shipped, while `ci.yml` is local. A full change writes its owed entry in the digest; simple work writes it before integration. |
| Cutting a release | Rename `[Unreleased]` to the version; the workflow publishes it verbatim and rejects an empty section. Exclude generated changelogs, PR lists, and `@mentions`; before release verification only, update the locally installed binary with `go install ./cmd/clue`. |
