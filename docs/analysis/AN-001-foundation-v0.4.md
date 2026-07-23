---
id: AN-001
type: analysis
status: active
links: [G-001]
title: Cliewen Foundation Document (v0.4, 2026-07-11)
---

> **Historical record — frozen.** This document seeded the corpus and is never edited (beyond this banner). The living truth is `/docs` itself: method improvements land as ADR-backed changes to the permanent corpus through the change loop, never as edits here. Where this document and the corpus disagree, the corpus wins — that is the inversion it proposed.

# Cliewen — Foundation Document (v0.4, updated 2026-07-11)

> Cliewen (Old English *cliewen*, "ball of thread" — the word that became *clue*). A methodology and toolchain for agent-driven software development where the documentation corpus is the system-of-record and the thread from goal to test is mechanically enforced. CLI binary: `clue`.
>
> **Short description (org / domain / README first line — keep identical):** *"A verifiable thread from goal to test — docs as the system of record for agent-driven development."*

**How to read this document (note to Claude Code):** this is a contract, not inspiration. Implement against it; do not redesign the taxonomy or the methodology. Open questions are marked as such — everything else is decided.

---

## 1. The Core Claim

**SDD documents the change; Cliewen documents the system.**

Every SDD framework (Spec-Kit, OpenSpec, Kiro) treats the *change* as the unit: spec → plan → tasks → implement → next spec. Cliewen inverts this: the *permanent documentation corpus* (`/docs`) is the durable truth about the system, and changes are transient deltas that are **digested into** the corpus at merge. Provenance does not stop at "this feature matched its spec" — it answers "why does the whole system look like this, and what chain of goals, decisions and acceptance criteria led here."

**The red thread is core, not profile.** The goal → scenario → AC → test chain is the methodology's load-bearing claim; capabilities are its middle anchor and are therefore core. Granularity is flexible, the anchor is mandatory: all software has behavior with acceptance criteria — a CLI's capabilities are its commands, a library's its public API surfaces, a platform's its services. A version of Cliewen without capabilities is a document management system, not Cliewen.

## 2. The Founding Argument

**RUP and the spiral model were right about *what*, and failed on *who*.** Humans cannot be driven by process — they skip ceremony wherever the fence is lowest, and enforcement (reviews, managers) cost more than the discipline was worth. Agents are no more disciplined — unsupervised they delete failing tests and declare victory — but agent cheating can be **enforced away mechanically**, human cheating cannot. Hence the division of labor:

| Actor | Role | One-liner |
|---|---|---|
| **Skills** | Process knowledge | Tell the agent what the next right step is |
| **CLI (`clue`)** | Deterministic judge | Tells everyone whether it was done right |
| **CI** | The wall | Refuses to proceed if not (same binary as local) |
| **Human** | Decision-maker | Settles what machines cannot check: meaning |

Corollary (Goodhart guard): **machines enforce form, humans verify meaning.** The linter checks that AC-042 has a test; only PR review checks the test means anything. PR review is the one gate an agent cannot game.

Second half of the RUP autopsy: ceremony failed not only because it was expensive to produce but because nobody *consumed* the artifacts (shelf-ware). In Cliewen, `/docs` is the agent's working memory on every subsequent change. **An artifact nothing reads gets removed from the taxonomy.** This is the guard against AI doc-slop.

## 3. Git Is the Engine

No archive mechanism (OpenSpec's biggest machinery, replaced by convention):

- **The branch IS the proposal.** Transient artifacts live in `/changes/` on the branch only.
- **The PR IS the review gate.** Human verifies meaning here.
- **The merge commit IS the acceptance.** Before merge, the delta must be digested: permanent docs updated, ADRs written, transient files deleted.
- **Git history IS the provenance archive.** `git log docs/` is the audit trail; signed commits + required reviews give *verified* intent provenance for free.
- **CI gate: `main` must never contain `/changes/`.**

Provenance must be **repo-native, not forge-native** (never rely on PR descriptions — that makes the method GitHub-dependent).

Known open edges (acknowledged, not solved): git merges text, not meaning — two clean-merging proposals can still contradict semantically. Mitigation: small deltas + post-rebase consistency check. Long-lived branches rot — engine should refuse to proceed when N commits behind main on `/docs`.

## 4. Identity Is Frontmatter, Paths Are Addresses

**The ID (`CAP-012`) is the identity — eternal, path-independent. The file path is only the current address.** Resolution of the earlier tension ("paths are not identity, yet rules bind to folders"):

- `clue` discovers artifacts by **scanning for frontmatter**, not by path. The graph is built from `id` + `links`, never from directory position.
- The folder layout (§10) is **convention**: where `scaffold` puts things and where humans look. Layout linting is a *default-on profile rule*, not a kernel requirement. `validate` parses frontmatter first, checks layout second.
- **External systems reference IDs, never paths.** A ticket says "CH-041", not a GitHub URL. Internal links a linter can fix on a move; URLs in Jira tickets, Slack archives and compliance reports rot forever at the first `git mv`. This is the hardest argument for frontmatter identity.
- **Path stability remains a virtue.** Moves being *survivable* does not make them free (`git log --follow` is fragile; wild-world URLs die regardless). Convention: folders are stable addresses you don't move without reason; frontmatter is the insurance that makes it survivable when the reason arises. Insurance, not invitation.
- Future CLI affordance (post-baseline): `clue locate CAP-012` — ID to current path. The reason generated views beat hand-written links.

**Principle: never encode mutable state in file paths — only identity.** Status lives in frontmatter, never in folder names (no TODO/in-progress/done directories). Folder structure is *storage*; status views are *generated* (README indexes, `clue status`).

## 5. Three Artifact Lifetime Classes

1. **Permanent** — `/docs`. The system-of-record. Lives forever, updated at every merge.
2. **Transient** — `/changes/<id>/` on a branch. Proposal, tasks, open questions. Dies at merge (digested into permanent docs).
3. **Campaign** — plans. Live on `main` (parallel deltas must branch off them), mutate continuously, and are **frozen when the goal is reached** — never deleted (deleting breaks the G→P→CH chain; auditors ask about finished things). A completed plan gets `status: completed` and becomes **immutable** (lint rule: any change to a completed plan fails `validate`). Rot only afflicts living documents claiming to reflect the present; a frozen plan claims only to reflect the past — perfectly. Before freezing, the plan's *lessons* (direction changes, rejected paths) must be distilled into ADRs: the plan is the ledger, the ADRs are the wisdom. `/docs/plans/` thereby doubles as the project's achievement overview — the generated README index groups by status (TODO / active / completed).

Plan freshness is enforced by the **merge hook**: every proposal declares which plan item it serves (or explicitly declares itself plan-less — no fake plan items), and the merge updates plan status in the same commit. A plan can by construction never be more than one merge stale: it is an *account*, not a promise.

**The plan is itself a change.** Everything that mutates `main` goes through branch + PR — including creating a plan (a change whose digest is the plan file itself). Two kinds of plan mutation, with different rules:

- **Semantic change** (direction, scope, new/removed milestones — anything that alters what the plan *promises*): requires its own change/PR with human acceptance, backed by an ADR. Plan adjustments ARE decisions.
- **Bookkeeping** (milestone M-003 marked done): happens as part of the *feature change's* merge digest — never a separate PR, or process overhead kills the method.

Rule of thumb: *changing what the plan promises is a change; recording what happened is part of the digest.* Lintable: status fields may mutate in any merge; everything else in a plan file only in changes declaring themselves plan revisions.

## 6. The Red Thread (ID Scheme)

Machine-checkable chain, every link lintable:

```
G-xxx (goal: who wants it, why)
  └─ P-xxx / M-xxx (plan / milestone)
       └─ CH-xxx (change / delta)
            └─ CAP-xxx (capability)
                 └─ AC-xxx (acceptance criterion, Gherkin)
                      └─ test tag (positive + negative pair)
```

Cross-cutting, checked against every proposal:
- **C-xxx** constraints (laws, regulations, licenses, policies — things you *must not break*, distinct from requirements). This is the compliance gold: "show that every change was assessed against GDPR constraint C-007."
- **QS-xxx** quality scenarios (NFRs made verifiable: "under X load, respond within Y ms").

Constraints carry two extra frontmatter fields:
- **`source:`** — where the constraint comes from. The field exists from day one (three words), but **external catalog mechanics are deferred**: v1 uses local constraints only. External catalogs (e.g. dot-principles) plug in later via this field — door defined, door closed, same maneuver as the V3 operations door.
- **`enforcement:`** — `machine` | `agent` | `human`. The field exists from day one; **the elaboration baseline implements only the `machine` class.** (`agent`: a skill instructs the agent to assess; human verifies in PR. `human`: PR review only. Boundaries are soft — an agent check run in CI becomes machine-ish — which is exactly why only `machine` ships first.)

**dot-principles is orthogonal to Cliewen — explicitly.** Cliewen enforces *process and traceability*; dot-principles guides *craft quality*. A repo can run either, both, or neither. Coupling them would make 373 principles an adoption barrier and repeat the RUP mistake: all-or-nothing bundles die. The `/dot-*` commands live on as independent skills beside the `clue-*` skills.

## 7. Frontmatter: The Machine-Readable Graph

Every markdown artifact in the corpus carries YAML frontmatter. Without it, `clue validate` would parse free prose for IDs and relations — fragile. With it, the corpus is a machine-readable graph: **frontmatter is for the CLI, the markdown body is for humans.** (Plain-Text-as-Code applied to the methodology itself.)

Design rules:

- **Common core + type-specific extension.** Every artifact shares: `id`, `type`, `status`, `links` (the relations the linter traverses). Each type adds its own: ADRs get the MADR fields plus `author: agent|human`, `accepted-by:`, and provenance status `inferred|verified`; capabilities get a `goal:` reference; constraints get `source:` and `enforcement:` (see §6).
- **Every field must have a consumer.** A field neither `clue` nor a skill reads gets removed — otherwise agents fill 20-field headers with plausible values nobody checks: doc-slop in metadata form. Start with five fields; extend when the linter needs something, never the other way around. Standardize schemas for only the artifact types the baseline actually exercises; more types crystallize through use, not upfront.
- **Keep it short.** GitHub renders frontmatter as a table at the top of the file (not hidden). Acceptable — the table is informative — but it is one more reason for brevity.

Exact schema per exercised artifact type is settled in the elaboration baseline (cheap, on paper).

## 8. Two-Tier Provenance

Every reconstructed or agent-made claim carries a status:

- **`inferred`** — agent-reconstructed (brownfield extraction, agent decisions). Not yet truth.
- **`verified`** — human has accepted. Promotion inferred → verified is the human act that makes provenance auditable.

Applies to: historical decisions extracted from brownfield code, agent framework choices (`author: agent`, `accepted-by: human`), and plan adjustments (agents may *propose* plan changes; humans accept — plan adjustments ARE decisions).

**Open questions are artifacts**: the agent writes them, blocks, and human answers are converted into ADRs. Every human decision joins the provenance chain. Analysis must leave corpses: every spike ends in a findings document; every rejected alternative becomes a rejected-ADR. Discarded options are half of "why does the system look like this."

## 9. Brownfield: Extraction First

The gap every competitor ignores (all are greenfield-biased). Flow: extract the artifact graph from existing code (all artifacts born `inferred`), human promotes to `verified` incrementally. Historical decisions are accepted *as historical*.

**Dogfood targets — two guinea pigs, two distinct roles, do not mix:**
- **The Cliewen repo itself**: greenfield dogfood — conventions from commit one. "Does it work?" → `git log docs/` in the engine's own repo.
- **model2diagram**: brownfield extraction test case (exit criterion 3). Right size (smaller than Tank Royale, real code unlike c64-ctx), and the meta-point writes itself: the tool that generates diagrams becomes the first system whose corpus contains them.

## 10. The Lifecycle (classic, deliberately)

Idea → rough plan sketch → **analysis** (clear risks and unknowns first — spiral's core; prototypes/spikes with findings docs) → real plan with verifiable milestones → per milestone: change loop (branch → implement → digest → merge) → plan adjusted as reality bites (human-accepted). Iterative by construction: each branch is one turn of the spiral, git is the spiral.

Terminology (precision matters in a book about RUP): **XP spikes** are the throwaway investigations in the analysis phase, ending as findings docs in `/docs/analysis`. The initial 2–3 week build is not a spike — it is an **elaboration baseline** in RUP's sense: a running skeleton that proves the architecture and is built upon, not thrown away. Risk-driven elaboration: the exit criteria are Boehm's "retire the biggest risks first" as running code.

Post-merge door (defined now, built later/V3): production findings enter the chain as a new goal or a new constraint.

Positioning sentence: *this is not a new development method — it is ordinary good software engineering (analysis, design, implementation, test, iteration) where agents do the work and the documentation never falls behind.* Everyone else sells a paradigm shift; Cliewen sells what senior engineers already believe, finally enforced.

### SDLC mapping (positioning, not core — the colleagues' question)

Cliewen IS an SDLC — with enforced artifacts:

| SDLC phase | Cliewen home |
|---|---|
| Requirements / planning | goals, plans, analysis phase (spikes with findings) |
| Design | architecture.md (system scope) + per-capability design.md |
| Implementation | the change loop (branch → implement → digest → merge) |
| Testing | AC↔test contract in CI — not a phase but a gate in every iteration |
| Deployment / operations | **deliberately out until V3** — only the door is defined |

Be honest about the gap when asked: release management, versioning and operations are not covered; Cliewen owns the stretch from intention to verified merge. One-liner for veterans: *RUP's disciplines without RUP's ceremony cost.*

### Ticket systems (integration pattern, not core)

Plan = epic, change = story, tasks = subtasks. For solo and small teams, Cliewen replaces the ticket system — with the crucial difference that these "tickets" cannot drift from the code: same repo, same gate. But ticket systems also do things Cliewen deliberately does not: intake from non-developers, cross-team coordination, SLAs, management reporting, non-code work. **Do not build that** (the Jira-reinvention trap). Enterprise positioning: *Cliewen is the system-of-record for engineering truth; the ticket system is a coordination tool — and the arrow points inward:* tickets reference CH-IDs (never paths, per §4), never the reverse.

Open gap: the **inbox** — where does an idea/bug report live *before* it is a change? Candidates: GitHub Issues as non-binding forecourt, or goals with `status: proposed`. Decide early in the baseline.

---

## 11. Directory Structure

### Permanent corpus (on `main`, forever)

```
/docs
  README.md                    # map of the corpus; entry point for agents & humans
  /goals                       # G-xxx: one file per goal — who wants it, why
    README.md
  /plans                       # P-xxx: campaign layer — FLAT (no status folders);
    README.md                  #   status in frontmatter; index groups by status
  /capabilities                # one folder per capability
    README.md
    /CAP-012-user-auth
      README.md                # what + why, goal reference (context, free prose)
      criteria.md              # ONLY Gherkin with AC-IDs — the lintable contract
      design.md                # how — ALL its diagrams inline (Mermaid)
  /architecture                # system scope: architecture.md, system context,
    README.md                  #   cross-cutting sequence diagrams (inline Mermaid)
  /decisions                   # ADR-xxx: MADR + status (inferred|verified),
    README.md                  #   author (agent|human), incl. rejected ADRs
  /constraints                 # C-xxx: laws, regulations, licenses, policies;
    README.md                  #   source: + enforcement: machine|agent|human
  /quality                     # QS-xxx: quality scenarios (verifiable NFRs)
    README.md
  /analysis                    # spike findings, extraction reports
    README.md
```

Layout is convention with default-on linting (see §4): identity lives in frontmatter; this layout is where `scaffold` puts things and humans look.

**One primary consumer per file** (capability-folder invariant, applies broadly): README.md is for humans/agents who must *understand*; criteria.md is for the linter/tester that must *verify* (a pure-Gherkin file is trivially parseable; Gherkin mixed into prose breaks on every rephrasing); design.md is for the implementer who must *build*. A fourth file that cannot answer "who reads it?" does not get in.

**Architecture vs design — different scopes, not a hierarchy.** Architecture is the whole (system context, cross-cutting structure, the expensive-to-change). Design is the part (how one capability is built) and is **colocated with its capability**. No top-level `/design` tree competing with `/architecture`; one home per scope, or the trees drift apart.

**README.md in every `/docs/**` folder, no exceptions** (GitHub renders README.md when browsing a directory; index.md does not). Hand-written index files are the first artifact type to rot, so indexes are *generated*: human-written prose at the top (the folder's purpose), generated index below between markers — `clue scaffold` updates them, `clue validate` fails on drift (file missing from index, or index referencing a deleted file).

**All documents must render on GitHub as-is.** Consequence: diagrams are Mermaid in fenced code blocks *inside* design.md / architecture.md — never standalone `.mmd` files (GitHub does not render those). If a diagram is too large to read when rendered, that is a design problem, not a rendering problem: split it, or lift it to architecture.md. **Diagrams humans cannot read are not usable and do not belong in the corpus.**

### Transient workspace (branch only — CI forbids on `main`)

```
/changes
  /<CH-xxx-slug>
    proposal.md          # what & why; references G/P/M items (or "plan-less")
    tasks.md             # task breakdown (kept from OpenSpec — it works)
    open-questions.md    # blocking questions; answers become ADRs
```

### Method machinery (on `main`)

Rule of thumb: `.agents/` = what agents read, `/.cliewen` = what the CLI reads, `/docs` = what the system *is*. Anything that can be convention instead of configuration stays out of `/.cliewen` — the emptier it is, the closer to the convention-over-configuration promise.

```
/.agents
  /skills
    /clue-analysis/skill.md    # book's existing .agents/ hub convention;
    /clue-plan/skill.md        #   clue- prefix namespaces against user skills
    /clue-delta/skill.md
    /clue-verify/skill.md
/.cliewen                      # CLI-private ONLY — keep small
  /templates                   # scaffolding sources for `clue scaffold`
AGENTS.md                      # routing hub (existing pattern from the book)
.github/workflows/             # CI: clue validate + tests + AC traceability gate
```

---

## 12. Getting Started Is a Capability

**Onboarding is CAP-001 in the Cliewen repo itself:** *a new user goes from install to first green `clue validate` in under 30 minutes.* The method's first enforced requirement is its own accessibility — Gherkin-testable and usable in marketing. (Lesson from Spec-Kit's 88k stars: instant usability beats conceptual superiority.)

The layered guide — keep the layers separate (OpenSpec's guide was hard to penetrate because it mixed them):

1. **Command (seconds):** `clue init` scaffolds the /docs taxonomy, README indexes, skills in `.agents/`, CI workflow — the whole convention materialized in one call. The most important part of the guide is not a document.
2. **Quickstart (5 minutes):** one page — install, `clue init`, run your first change loop, watch `validate` go green.
3. **Skills (learned during use).**
4. **Book (the why — depth, secondary).**

---

## 13. Project Overview & Elaboration Baseline

**Naming architecture (settled and executed):**
- Methodology/brand: **Cliewen** · CLI binary: **`clue`** · package: `cliewen`
- Domain: **cliewen.dev** — *purchased* ✓
- GitHub: **github.com/cliewen** org — *reserved* ✓ (verify spelling: must be c-l-i-e-w-e-n; rename today if not — external links to a misspelling never rot away)
- Short description everywhere (org, domain, README line one — identical): *"A verifiable thread from goal to test — docs as the system of record for agent-driven development."*

**The CLI is the judge and only the judge.** Deterministic, stateless, no AI, no orchestration, no TUI luxuries. Commands: `init`, `validate`, `scaffold`, `status`, later `extract` and `locate`. `scaffold` also (re)generates the README index blocks; `validate` fails on index drift; `clue status` renders the kanban-style view (status lives in frontmatter, never in paths). Everything that can be a skill/convention is pushed out of the CLI. Goal: a CLI so boring it is finished. First ADR: implementation language (Go vs Rust; criteria: single-binary distribution, maintainer fluency, agent fluency). `author: human`.

**Elaboration baseline (2–3 weeks, hard exit criteria):**
1. The loop closes once, end-to-end: goal in → proposal branch → implement → permanent /docs updated → transient files gone → clean merge.
2. The thread is machine-checkable: linter fails the build when an AC lacks a test or a test lacks an AC (reuse ADR-015 pattern).
3. Brownfield extraction works on **model2diagram**.

Not a throwaway spike: a running skeleton built upon afterwards (§10 terminology). Dogfooding is the proof — the Cliewen repo runs Cliewen conventions from commit one.

**Explicitly OUT of the baseline:** multi-agent orchestration, semantic consistency checking of shared docs (parked as known edge), plan milestone mechanics beyond one markdown file + one linter rule (plan-reference check), production feedback loop (V3 door only), dot-principles / external catalog integration (door defined via constraint `source:`, nothing more), `enforcement:` classes beyond `machine`, `clue locate`.

**Post-baseline refactoring note (not current architecture):** a kernel/profile/tooling/integrations layering (small artifact protocol + opinionated engineering profile on top) is plausible and useful as the *book's* disposition — but abstractions are extracted from multiple working instances (rule of three), not designed from zero. The two guinea pigs provide the data points; extract the kernel *after* they run. Do not build the layering now.

**Relationship to the books:** V1 (*Intent Engineering for Coding Agents*) ships now, unchanged — its market test must run unpolluted. A dedicated Cliewen book comes later, when V2 has matured through real use; it will harvest its own comments and feedback. Build-log articles during the baseline are decision-free bridge material: they become chapters in whichever book structure wins, and they market along the way. Whether the Cliewen book *supplements* or *supersedes* V1 is an open question (`author: human`, post-baseline) — the engine likely obsoletes the OpenSpec chapters.

**Market context (July 2026):** the territory is empty but converging — Tessl pivoted spec→context, Microsoft publicly argues "the agent decided is not acceptable to auditors," arXiv formalized SDD-as-gates. Nobody has inverted change-centric → corpus-centric, and nobody owns brownfield + verified provenance. The window is short. Speed over polish.
