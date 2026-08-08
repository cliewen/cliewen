---
id: AN-018
type: analysis
status: active
links: [P-013, PDR-029, PDR-013, PDR-028, PDR-023, ADR-021, ADR-048, C-011, C-013, G-001]
title: Statement register for the six shipped skills and the routing hub
---

# AN-018 — Statement register for the six shipped skills and the routing hub

## The risk this spike retires

Every previous attempt to simplify Cliewen's methodology carriers argued case by case and stalled, because no one had written down what each statement is for. [PDR-029](../decisions/PDR-029-simplification-tests-by-surface.md) supplies the test; this spike supplies the material it is applied to. The risk is that M-063 begins trimming prose without a register, and simplification stalls a fourth time — or worse, succeeds at removing words while removing rules nobody noticed were load-bearing.

Nothing here is removed, reworded, or reordered. No trace, citation, or marker was added to any carrier.

## Evidence boundary

- **Pinned revision:** the carriers were read at `e9a7d07` on `main` (`cliewen/cliewen`). Every locator below is to that revision.
- **What was read:** `AGENTS.md`, the six generated skills under `.agents/skills/*/skill.md`, and — because [PDR-029](../decisions/PDR-029-simplification-tests-by-surface.md) makes tracing, checkability, and order properties of the *authored* fragments — the generator sources under `internal/skills/source/` that render them. The rendered files and their sources were compared; no divergence was found beyond template expansion.
- **What was not read:** the public guide, `CONTRIBUTING.md`, CLI text, and the scaffolded adopter copies. M-064 carries those. `internal/scaffold/templates/AGENTS.md` was opened once, for the single coverage question in HUB-59, and not registered.
- **Environment:** prepared, not clean. Windows 11, Go 1.26.5, a `clue` binary stamped `dev`, skills stamped `0.13.0`. Nothing in this spike depends on the toolchain: every result is a reading of committed text, reproducible by opening the same files at the same revision.
- **Confidence classes.** Class, duplication, checkability, and order are **observed** — they are properties of the text and a second reader can confirm or refute each one by opening the file. A trace is observed when the named artifact contains the rule; the judgement that it is the *narrowest* such artifact is an **inference**, because it rests on having searched the corpus rather than on having read all of it. Every `NONE` verdict was reached by grepping the live corpus for the rule's distinctive wording and its synonyms, and a `NONE` is therefore a claim that no artifact was *found*, not a proof that none exists. Where a rule traces only to a frozen analysis or to an architecture file, that is stated rather than counted as a trace.
- **Not established.** This is one reading by one agent. The register is a first pass, and the segmentation rule below exists precisely so a second pass can disagree in a way that is legible.

## What counts as a single statement

The milestone requires this to be precise enough that an independent second pass, which did not write the definition, segments the same prose the same way. The rule is applied in order; each step consumes the output of the one above it.

1. **Take the carrier's structural units in file order.** A heading, a paragraph, a list item (including each nested item), a table row, a checklist item, a fenced block, and the frontmatter block are each a unit. The generated-file HTML comment is a unit.
2. **Split each unit at sentence boundaries** — a `.`, `?`, or `!` followed by whitespace and a capital letter — ignoring boundaries inside a code span, an abbreviation, or an ellipsis.
3. **Re-split by obligation.** A sentence carrying more than one independent obligation, joined by a semicolon or by `, and` / `, but` / `, never` / `, nor` where each side has its own obligation verb, becomes one statement per obligation. These carry a letter suffix (`HUB-37a`, `HUB-37b`).
4. **Do not split a qualified obligation.** One obligation with subordinate conditions, however long the sentence, is one statement: the condition is part of the rule and cannot be evaluated apart from it.
5. **Re-join rationale.** A sentence that only supplies the reason for the obligation immediately before it — opening with *because*, *so that*, *which is why*, *which is where* — joins that obligation. Rationale is never its own statement.
6. **A step's leading bold label** (`**Branch:**`) belongs to its step.

Applying steps 1–6 to the seven carriers yields the identifiers used below: `HUB` for `AGENTS.md`, `ANL`/`PLN`/`UPG`/`DLT`/`VFY`/`EXT` for the six skill bodies, and `F-CT`/`F-DR`/`F-LC`/`F-DW`/`F-RB`/`F-FM` for the six shared fragments. A fragment is registered once and inherited by every skill that renders it.

## What each column means

**Class.** *Rule* when the statement states something an agent or human must, must not, or may do, or a condition a corpus state must satisfy — including a definition that decides whether another rule applies. *Connective* otherwise: headings, orientation, mechanism description, worked examples, and pointers to another section of the same file. A connective statement carries no obligation and is not scored against the three conditions.

**Trace.** The narrowest **live** corpus artifact that *states* the rule. `NONE` means no such artifact was found. Per PDR-029, derivability is not a trace.

The accepted types are decision, constraint, goal, acceptance criterion, and — from [PDR-031](../decisions/PDR-031-architecture-artifacts-are-traces.md), which this register's own findings produced — architecture, under the same restriction: it traces when the architecture file *states* the rule, never when the rule is merely derivable from it. The first pass ran before PDR-031 existed and recorded architecture separately; those rows are now ordinary traces.

The frozen Foundation Document ([AN-001](AN-001-foundation-v0.4.md)) is **not** an accepted trace and did not become one. It is an analysis, it is never edited, and its own banner says the corpus wins where the two disagree — an artifact that cannot be amended cannot carry a live rule. Eight statements rested on it and three more traced to nothing; all eleven are now traced to [PDR-030](../decisions/PDR-030-analysis-is-a-bounded-spike.md) and [PDR-006](../decisions/PDR-006-decision-records-are-typed.md), which were written for that purpose.

**Dup.** Whether a live carrier already states the same rule in the same reading path. Reading paths are counted, not files: `AGENTS.md` → `clue-delta` → `clue-verify` is one path an ordinary full change actually walks, so a rule stated in all three is stated three times to one reader. [ADR-021](../decisions/ADR-021-generated-standalone-skills.md)'s deliberate file-level repetition of a shared fragment across skills is **not** scored — a fragment is registered once.

**Chk.** Whether a reader can determine that the statement has been satisfied. `part` marks the failure PDR-029 names: several independent conditions offered as one obligation, which nobody can honestly tick when four hold and the fifth does not. The count of conditions is given.

**Ord.** `!` marks a statement that binds absolutely but is read after the procedure it constrains.

---

## Register — `AGENTS.md` (HUB)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| HUB-01 | `# Agent routing hub` | connective | — | — | — | |
| HUB-02 | "This repo runs **Cliewen**…" | connective | — | — | — | |
| HUB-03 | "Before your first tool call … run `clue latest --quiet`" | rule | PDR-023 | no | yes | |
| HUB-04 | "That is the whole instruction; the rest is why." | connective | — | — | — | |
| HUB-05 | "It prints one line … asking again is free." | connective | — | — | — | |
| HUB-06 | "If it says you are behind, route to `clue-upgrade`" | rule | ADR-043 | no | yes | |
| HUB-07 | "the human decides … nothing upgrades without that answer" | rule | ADR-043 | **UPG-06** | yes | |
| HUB-08 | "If `clue` reports `latest` as an unknown command, that *is* the answer" | rule | PDR-023 | no | yes | |
| HUB-09 | "The ordinary workflow commands … covers a session that runs none." | connective | — | — | — | |
| HUB-10 | "This is the only reason to run it unprompted" | rule | PDR-023 | no | part (prohibition with no enumerated set) | |
| HUB-11 | "it reaches the network, so it never belongs in a validation verdict or a required check" | rule | ADR-042 | no | yes | |
| HUB-12 | "Before loading the corpus, classify the requested work." | rule | PDR-011 | **F-CT-01** | yes | |
| HUB-13 | "Three rules set the tier … take the first rule that matches." | rule | PDR-002 | **F-CT-02** | yes | |
| HUB-14 | Tier 1, plain — definition and protected surfaces | rule | PDR-011 | **F-CT-03** | part (surface inventory, 10 items) | |
| HUB-15 | "Use an ordinary branch from `main` … no Cliewen bookkeeping." | rule | PDR-011 | **F-CT-03** | yes | |
| HUB-16 | "Plain changes do not consume the … slot and never build on unmerged work." | rule | PDR-011 | **F-RB-01** | yes | |
| HUB-17 | Tier 2, light — definition and typical cases | rule | PDR-002 | **F-CT-04** | yes | |
| HUB-18 | "Use a Cliewen branch and ready PR whose description is the proposal" | rule | PDR-002 | **F-CT-04**, **HUB-30** | yes | |
| HUB-19 | Tier 3, full — "Product behavior changes are full even when…" | rule | PDR-018 | **F-CT-05** | yes | |
| HUB-20 | "**Uncertainty escalates:** … take the higher one." | rule | PDR-002 | **F-CT-06** | yes | |
| HUB-21 | "**Discovery escalates immediately:** …" | rule | PDR-002 | **F-CT-06** | yes | |
| HUB-22 | "read `docs/README.md` only when the request does not name or resolve to an artifact" | rule | PDR-034 | no | yes | |
| HUB-23 | "run `clue context` directly and read its outgoing-link slice" | rule | PDR-034 | no | yes | |
| HUB-24 | "Read beyond that slice only when the task or a discovered edge requires it." | rule | PDR-034 | no | part (reader's judgement, no test) | |
| HUB-25 | "The `/docs` corpus remains the system-of-record and working memory." | rule | ARCH-003 | no | part | |
| HUB-26 | `## The rules that bind every change` | connective | — | — | — | |
| HUB-27 | 1 — "Everything that mutates `main` goes through branch + PR." | rule | C-012 | **F-RB-01** | yes | |
| HUB-28 | "the branch is the proposal; transient files live in `/changes/…` and are deleted in the digest commit" | rule | PDR-002 | **DLT-13**, **VFY-14** | yes | |
| HUB-29 | "`main` never contains `/changes/`." | rule | PDR-002 | **VFY-14** | yes | |
| HUB-30 | "A **light** change skips the workspace: the PR description is the proposal." | rule | PDR-002 | **HUB-18**, **F-CT-04** | yes | |
| HUB-31 | "Every change branches from the current tip of `main`" | rule | C-012 | **F-RB-01** | yes | |
| HUB-32 | "one Cliewen change is in flight per initiating author" | rule | PDR-007 | **F-RB-01** | yes | |
| HUB-33 | "**agents never merge their own PRs or push to `main`**" | rule | C-012 | **F-RB-03** | yes | |
| HUB-34 | "Full Cliewen changes use the supported merge-commit mode; squash and rebase are unsupported" | rule | PDR-021 | **F-RB-02** | yes | |
| HUB-35 | "Reviewing or helping update an existing PR does not mint another change" | rule | PDR-016 | **F-RB-01** | yes | |
| HUB-36 | 2 — "Ready means the hosted PR contains the reviewed and verified state." | rule | PDR-016 | **F-RB-07** | yes | |
| HUB-37a | "Every review of an existing PR names its hosted head" | rule | PDR-016 | **F-RB-05**, **VFY-19** | yes | |
| HUB-37b | "actionable findings become unresolved hosted review conversations where supported" | rule | PDR-016 | **F-RB-05**, **VFY-19** | yes | |
| HUB-37c | "a clean result applies only to its named commit" | rule | PDR-016 | **F-RB-05** | yes | |
| HUB-38 | "Any agent that edits becomes the updater for that turn: fetch … and only then resolve satisfied findings." | rule | PDR-016 | **F-RB-06**, **F-RB-07** | **part (6 conditions, one sentence)** | |
| HUB-39 | "A changed head or non-fast-forward rejection requires reconciliation and renewed verification" | rule | PDR-016 | **F-RB-06** | yes | |
| HUB-40 | "newer accepted `main` is merged … without rewriting its history" | rule | PDR-016 | **F-RB-01** | yes | |
| HUB-41 | "a merged or closed PR stops with local work reported as unpublished" | rule | PDR-016 | **F-RB-06** | yes | |
| HUB-42 | "Before publishing … automatically run `clue-verify`, require a clean worktree, and complete this exact hosted-head handoff" | rule | PDR-012 | **F-RB-07**, **DLT-16** | part (3 conditions) | |
| HUB-43 | "A human-requested local stopping point is preserved work, but … not mergeable." | rule | C-012 | **F-RB-07** | yes | |
| HUB-44 | 3 — "Every Cliewen proposal declares which plan item it serves … or plan-less." | rule | C-005 | **DLT-05**, **VFY-06** | yes | |
| HUB-45 | "The merge digest updates plan bookkeeping … sets it `completed` … designates a successor" | rule | log 2026-08-02 | **PLN-08**, **VFY-07** | yes | |
| HUB-46 | 4 — "Open questions are artifacts. When blocked … write it to `open-questions.md` and stop" | rule | C-011 | **DLT-07** | yes | |
| HUB-47 | 5 — "Machines enforce form; humans verify meaning." | connective | — | — | — | |
| HUB-48 | "Never weaken a test or a lint rule to make a build pass — surface the conflict instead." | rule | C-004 | **VFY-01d**, **F-RB-04** | yes | |
| HUB-49 | 6 — "Markdown prose is never hard-wrapped." + line-break rule | rule | C-001 | no | yes | |
| HUB-50 | 7 — "Release notes are written for users, in `CHANGELOG.md`" | rule | C-002 | no | yes | |
| HUB-51 | "A Cliewen change that affects shipped behavior … adds its entry to `[Unreleased]`" | rule | C-002 | **F-LC-01** | yes | |
| HUB-52 | "Plain editorial changes add no release note." | rule | C-002 | no | yes | |
| HUB-53 | "Cutting a release renames that section … the workflow publishes it verbatim and fails without it." | rule | ADR-012 | no | yes | |
| HUB-54 | "Auto-generated changelogs, PR lists, and @mentions never appear on a release." | rule | ADR-012 | no | yes | |
| HUB-55 | "update the locally installed `clue` from this checkout with `go install ./cmd/clue`" | rule | log 2026-07-31 | no | yes | |
| HUB-56 | 8 — "The core is behind a red line." + the three-element definition | rule | ARCH-003 | no | yes | |
| HUB-57 | "A change that alters what any of these means … requires an explicit decision record and human acceptance." | rule | C-013 | no | yes | |
| HUB-58 | "Periphery never constrains the core." | rule | PDR-013 | no | part | |
| HUB-59 | `## Skills` table, 5 rows | rule (routing) | ARCH-002 | no | yes | |
| HUB-60 | "The skill files are generated artifacts … never edit `.agents/skills/` directly" | rule | ADR-021 | no | yes | |

**Coverage gap in HUB-59.** The table lists five skills. Six ship. `clue-extract` is absent from this repository's hub while the scaffolded adopter hub (`internal/scaffold/templates/AGENTS.md`) lists all six. This may be deliberate — Cliewen has no source corpus to extract — but [ADR-043](../decisions/ADR-043-upgrade-skill-is-a-managed-carrier.md) states a managed set of six, and the hub is the routing surface. Escalated as **Q-07**.

---

## Register — `clue-analysis` body (ANL)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| ANL-01 | `# clue-analysis` | connective | — | — | — | |
| ANL-02 | "Use when a change has unclear risks or unknowns — **before** planning or implementing." | rule | PDR-030 | no | yes | |
| ANL-03 | "Spiral's core: retire the biggest risk first." | rule | PDR-030 | no | part | |
| ANL-04 | 1 — "Name the risk or unknown in one sentence." | rule | PDR-030 | no | yes | |
| ANL-05 | "If you cannot, that is the first finding." | rule | PDR-030 | no | yes | |
| ANL-06 | 2 — "Cite outward precisely … full address … `clue:` identity form." | rule | ADR-040 | **VFY-05** | yes | |
| ANL-07 | "A bare forge number means 'this repository' and is wrong the moment it does not" | rule | ADR-040 | **VFY-05** | yes | |
| ANL-08 | "Findings pinned to a revision keep the address they observed" | rule | ADR-040 | no | yes | |
| ANL-09 | 3 — "Establish the evidence boundary before investigating: pin … record … distinguish …" | rule | log 2026-07-18 | no | **part (3 conditions)** | |
| ANL-10 | "Repository activity is evidence of activity, not maintainer intent" | rule | log 2026-07-18 | no | yes | |
| ANL-11 | "Classify every verification result as either a clean disposable or a prepared environment." | rule | log 2026-07-29 | no | yes | |
| ANL-12 | "A clean result supports onboarding reproducibility only when it has no local prerequisites" | rule | log 2026-07-29 | no | yes | |
| ANL-13 | "A prepared result names its prerequisites and establishes only what that environment demonstrated." | rule | log 2026-07-29 | no | yes | |
| ANL-14 | 4 — "Before treating a statistical claim as evidence, name the versioned corpus and population, …" | rule | log 2026-07-29 | no | **part (7 conditions)** | |
| ANL-15 | "Do not turn an environment-sensitive quality claim into a deterministic acceptance criterion." | rule | AC-055 (CAP-003) | no | yes | |
| ANL-16 | "When assessing adoption, name the governance or process changes … not neutral scaffolding." | rule | log 2026-07-29 | no | yes | |
| ANL-17 | 5 — "Run a **spike**: a throwaway investigation…" | rule (definition) | PDR-030 | no | yes | |
| ANL-18 | "Spikes are disposable; their findings are not." | rule | PDR-030 | no | yes | |
| ANL-19 | 6 — "End every spike with a findings document in `/docs/analysis/`" | rule | PDR-030 | no | yes | |
| ANL-20 | "(`AN-xxx-slug.md`, frontmatter: `id`, `type`, `status`, `links`, `title`)" | rule | C-009 | **VFY-02** | yes | |
| ANL-21 | "Include what was tried, what was rejected, and why" | rule | PDR-030 | no | yes | |
| ANL-22 | "If the finding is an incident where the corpus was green but reality proved a claim wrong, add `reality: contradicted` and link every failed capability…" | rule | ADR-035 | no | part (2 conditions) | |
| ANL-23 | "This records the edge from reality; it does not ingest production telemetry or open the operations loop." | rule (boundary) | ADR-035 | no | yes | |
| ANL-24 | 7 — "Route any outcome that constitutes a decision under **Decision records** below." | connective (pointer) | — | — | — | |
| ANL-25 | "A rejected alternative that is itself a decision gets a rejected decision record, not only a paragraph." | rule | PDR-006 | no | yes | |
| ANL-26 | 8 — "Feed findings to `clue-plan` or `clue-delta`." | rule | PDR-030 | no | yes | |
| ANL-27 | "Analysis with no consumer is doc-slop; do not write it." | rule | PDR-030 | no | part | |

**`clue-analysis` was the weakest carrier by tracing, and the split was systematic rather than random.** Its *evidence-discipline* rules (ANL-09 … ANL-16) traced cleanly — seven of them to two decision-log rows written in CH-025 and CH-080, and ANL-15 to AC-055. A fourth review pass found that [AC-055](../capabilities/CAP-003-extract/criteria.md) also states ANL-11 through ANL-14 and ANL-16 word for word, so those five rows have two live artifacts stating them and the register names the log rows. Which is the *narrowest* when a criterion and a decision both state a rule is not something PDR-029 settles; it is left as an observation rather than re-traced here. Its *workflow spine* — what a spike is, that one ends in a findings document, that a findings document records rejected options, that analysis needs a consumer — did not trace at all: most of those statements rested on [AN-001](AN-001-foundation-v0.4.md), the frozen Foundation Document, and the remainder on nothing found. That is the spine of the skill governing how Cliewen writes its shared memory, in a carrier whose evidence-discipline half traces cleanly.

**Answered (Q-01).** They were real rules nobody had recorded, which is the case PDR-029 says is repaired by writing the missing decision rather than deleting the sentence. All but one are a single rule, now stated by [PDR-030](../decisions/PDR-030-analysis-is-a-bounded-spike.md); ANL-25 is about decision records rather than about analysis and is stated by [PDR-006](../decisions/PDR-006-decision-records-are-typed.md), amended in the same change. The rows above carry the resulting traces.

---

## Register — `clue-plan` body (PLN)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| PLN-01 | `# clue-plan` | connective | — | — | — | |
| PLN-02 | "Use when creating a plan or changing what a plan promises." | rule | PDR-008 | no | yes | |
| PLN-03 | 1 — "Create or revise a plan through `clue-delta`; a plan mutation is itself a branch and PR." | rule | C-012 | **F-RB-01** | yes | |
| PLN-04 | "The digest is the plan file in `/docs/plans/`." | rule | PDR-008 | no | yes | |
| PLN-05 | 2 — "Keep plans as flat `P-xxx-slug.md` files with status in frontmatter (`draft` → `active` → `completed`)." | rule | ADR-025 | no | yes | |
| PLN-06 | "Milestones (`M-xxx`) are rows in the plan's milestone table, each with a verifiable exit criterion." | rule | C-010 | no | yes | |
| PLN-07 | 3 — "**Semantic:** Direction, scope, milestone addition/removal … requires human acceptance and a decision record." | rule | PDR-008 | no | yes | |
| PLN-08 | "Agents may propose; only humans accept." | rule | C-012 | **F-RB-03** | yes | |
| PLN-09 | "The default vehicle is a dedicated plan change and PR." | rule | PDR-008 | no | yes | |
| PLN-10 | "A revision discovered during implementation may ride with that implementing change only when …" | rule | PDR-008 | no | **part (4 conditions)** | |
| PLN-11 | "**Bookkeeping:** Marking a milestone done belongs in the implementing change's merge digest, never a separate PR." | rule | log 2026-08-02 | **HUB-45** | yes | |
| PLN-12 | "Closing the plan is the same bookkeeping … sets it `completed`, in that digest." | rule | log 2026-08-02 | **HUB-45**, **VFY-07** | yes | |
| PLN-13 | "A campaign is over the moment its last milestone is evidenced, so leaving it `active` publishes an index claiming work is in flight that is not." | rule + rationale | log 2026-08-02 | no | yes | |
| PLN-14 | "Designate the successor plan there too when one is decided; not having decided one never holds the closure open." | rule | log 2026-08-02 | **HUB-45** | yes | |
| PLN-15 | "Every milestone's evidence must be in the table before that digest lands, because the closed plan is immutable afterwards." | rule | C-008 | no | yes | |
| PLN-16 | 4 — "Treat `status: completed` as immutable and never delete a completed plan." | rule | C-008 | **PLN-15** | yes | |
| PLN-17 | "Before freezing it, distill its durable lessons and rejected paths into decision records." | rule | **NONE — withdrawn, removed by M-063** | no | part | |

`clue-plan` is among the smallest carriers and the only skill with no order defect, but it is not as clean as a first reading suggests. Two of its statements trace to nothing (PLN-06, PLN-17) and two fail checkability (PLN-10, PLN-17). Its duplications split three ways, and which layer owns the repair differs with each: PLN-11, PLN-12 and PLN-14 restate the routing hub; PLN-03 and PLN-08 restate the shared review-boundary fragment; and PLN-16 restates PLN-15 *within the same file*. M-063 should not treat this carrier as finished. In-carrier duplication is not peculiar to it either — rows across the register name a duplicate in their own carrier, most of them in `clue-verify`.

---

## Register — `clue-upgrade` body (UPG)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| UPG-01 | `# clue-upgrade` | connective | — | — | — | |
| UPG-02 | "Use when a repository already uses Cliewen and the human wants to … bring it up to a newer release." | rule | ADR-043 | no | yes | |
| UPG-03 | "This is a route into a reviewed repository change, never a background update or authority to merge." | rule | ADR-043 | **F-RB-03** | yes | |
| UPG-04 | 1 — "Run `clue latest`." + what it determines | rule | ADR-042 | **HUB-03** (different trigger) | yes | |
| UPG-05 | "Do not reproduce an installation command here: one distributed skill cannot know the user's platform." | rule | ADR-043 | no | yes | |
| UPG-06 | "If the release list cannot be reached, explain that Cliewen cannot tell and stop; do not call the repository current." | rule | ADR-042 | no | yes | |
| UPG-07 | 2 — "read that release's notes, including its `### Migration` section" | rule | ADR-039 | no | yes | |
| UPG-08 | "Identify the coordinated corpus, generated-skill, and CI-caller changes before proposing any repository write." | rule | ADR-039 | no | yes | |
| UPG-09 | 3 — "Ask the human whether to upgrade now or later." | rule | ADR-043 | **HUB-07** | yes | |
| UPG-10 | "Do nothing to the repository until they explicitly choose now." | rule | ADR-043 | **HUB-07** | yes | |
| UPG-11 | "A later answer is complete: report the available release and stop without creating a branch…" | rule | ADR-043 | no | yes | |
| UPG-12 | 4 — "make the repository green and create a branch through its normal review process" | rule | C-012 | **F-RB-01** | yes | |
| UPG-13 | "Move the binary and repository together: preview `clue migrate`, resolve every finding and notice … apply only the complete, preflighted plan with the required reversal-cost choice." | rule | ADR-039 | **F-LC-02** | **part (4 conditions)** | |
| UPG-14 | "Keep the managed skills, the thin caller, and any repository corpus obligations on the chosen release together." | rule | ADR-038 | no | yes | |
| UPG-15 | 5 — "Verify the upgraded repository, commit the complete candidate, run its required checks, and open a ready pull request." | rule | PDR-012 | **F-RB-07** | part (4 conditions) | |
| UPG-16 | "Never merge it: the repository's human merge boundary accepts the upgrade." | rule | C-012 | **F-RB-03** | yes | **!** |
| UPG-17 | "The shared methodology and review rules remain below." | connective | — | — | — | |

**Order in `clue-upgrade`.** UPG-16 is the last numbered step, and the review boundary that states the same prohibition is the last section of the file. The prohibition is genuinely inviolable and it is read after everything it constrains. UPG-03 partially rescues this by stating "never … authority to merge" in the opening paragraph, which is why UPG-16 is flagged and UPG-03 is not.

---

## Register — `clue-delta` body (DLT)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| DLT-01 | `# clue-delta` | connective | — | — | — | |
| DLT-02 | "Use for every Cliewen change: features, fixes, docs, and plans whose meaning belongs in the corpus." | rule | PDR-011 | no | yes | |
| DLT-03 | "Plain changes are classified by AGENTS.md before the corpus is loaded and do not invoke this skill." | rule | PDR-011 | **HUB-12**, **F-CT-01** | yes | |
| DLT-04 | "Apply the **Change scope and tiers**, **Decision records**, … below throughout the loop." | connective (pointer) | — | — | — | |
| DLT-05 | 1 — "Follow the review boundary and name the branch `ch-xxx-slug`." | rule | PDR-002 | **F-RB** (pointer) | yes | |
| DLT-06 | "Take the next free CH number by searching Git history and `/changes/` for the highest used number." | rule | ADR-009 — **superseded by ADR-048, whose mechanism has shipped** | no | yes | |
| DLT-07 | 2 — "For a full change, create `/changes/CH-xxx-slug/` and commit it before implementation" | rule | PDR-002 | **HUB-28** | yes | |
| DLT-08 | "`proposal.md` states what and why; its frontmatter `links` names the real plan item … or declares plan-less." | rule | C-005 | **HUB-44**, **VFY-06** | yes | |
| DLT-09 | "`tasks.md` is an ordered `- [ ]` checklist with dependencies first and at most one nested level." | rule | C-003 | no | yes | |
| DLT-10 | "Mark `[x]` the moment a task completes, never in batch at the end." | rule | C-003 | no | yes | |
| DLT-11 | "Mark an infeasible task `[-]` with its reason on the same line." | rule | C-003 | no | yes | |
| DLT-12 | "A behavior-changing task names the acceptance-criterion IDs it serves; if none exists, add the criterion before implementation." | rule | G-001 | no | yes | |
| DLT-13 | "Tests trace to criteria, never transient tasks." | rule | ADR-005 | **VFY-08** | yes | |
| DLT-14 | "`open-questions.md` records blocking questions. When one appears, write it and stop; the human answer becomes a decision record." | rule | C-011 | **HUB-46** | yes | |
| DLT-15 | "A human may opt into a spec-first pause after Propose. Record the pause in `tasks.md` and stop …" | rule | PDR-017 | no | yes | |
| DLT-16 | 3 — "Update the permanent corpus. Capabilities own README, criteria, and design files." | rule | ADR-025 | no | yes | |
| DLT-17 | "Write criteria as Gherkin tagged with their canonical `<PREFIX>-<digits>[lowercase-suffix]` identity" | rule | ADR-009, ADR-037 | **VFY-08** | yes | |
| DLT-18 | "every new or materially revised criterion declares `Test-type: …` and gets focused positive and negative evidence in that class (or `(single-direction)`)" | rule | ADR-032 | **VFY-08** | part (3 conditions) | |
| DLT-19 | "`Human` needs no code evidence — the acceptance brief's criteria line is its proof; … never as a placeholder for a test not yet written." | rule | ADR-033 | **VFY-08** | yes | |
| DLT-20 | "A criterion genuinely not yet proven carries `@draft` … exempting only that criterion" | rule | ADR-033 | **VFY-08** | yes | |
| DLT-21 | "Every test declares exactly one purpose: the criterion ID, `Unit`, `Sanity`, or `Arch` …" | rule | ADR-006 | no | yes | |
| DLT-22 | "On the JVM, all three evidence parts attach to the same Java or Kotlin executable …" | rule | ADR-036 | **VFY-08** | part (3 parts, one sentence) | |
| DLT-23 | "When a criterion's meaning changes, retire it with `@retired`, keep the tombstone, mint a new ID, and remove or retag its tests." | rule | ADR-007 | no | part (4 conditions) | |
| DLT-24 | 4 — "After every task is `[x]` or `[-]` with a reason, update permanent `/docs`, regenerate README indexes, apply repository-local digest conventions, record decisions, and update plan bookkeeping." | rule | C-003 | **VFY-07**, **VFY-13** | **part (5 conditions)** | |
| DLT-25 | "Retiring a non-criterion artifact means deleting its file in this same digest … naming the dead ID in a `supersedes:` field" | rule | ADR-034 | no | part (2 conditions) | |
| DLT-26 | "Delete the change workspace." | rule | PDR-002 | **HUB-28**, **VFY-14** | yes | |
| DLT-27 | "The digest is never a task in `tasks.md`; deletion is the digest, so a self-referential digest task cannot be completed honestly." | rule | C-003 | no | yes | |
| DLT-28 | 5 — "Run `clue-verify`, including its automatic agentic review loop … then open the PR under the review boundary." | rule | PDR-012 | **HUB-42**, **F-RB-03** | yes | |
| DLT-29 | "For a full change, fill the acceptance brief … with the plan item and whether it remains wanted, every added or changed criterion and its scenario-resolution verdict … and what merge binds or supersedes" | rule | PDR-017 | no | **part (4 conditions)** | |
| DLT-30 | "keep it to one screen and never leave template placeholders" | rule | PDR-017 | no | yes | |
| DLT-31 | "Never ask the human to initiate the review." | rule | PDR-012 | **VFY-16** | yes | |
| DLT-32 | "Merging accepts the change; decision provenance follows **Decision records** below." | rule | PDR-004 | **F-DR-02** | yes | |
| DLT-33 | "Keep deltas small: Git merges text, not meaning." | rule | **NONE — withdrawn, removed by M-063** | no | part | |

**DLT-06 is the register's one demonstrated defect, and it was found by running the tool rather than by reading.** [ADR-048](../decisions/ADR-048-corpus-wide-id-ledger.md) supersedes ADR-009's "the corpus is the registry" clause for *every* native ID prefix, `CH` included, and replaces scan-and-max allocation with a persisted ledger — precisely to stop a deleted artifact's ID being re-minted once no scan can see it. ADR-048's closing paragraph defers its own implementation to P-011's M-052; that paragraph is now history. `.clue/id-ledger.yaml` is live in this repository, `clue id next` and `clue id live` are shipped subcommands, and `clue validate` rejects an artifact absent from the ledger. **This change followed DLT-06 as written, derived its CH number by grepping Git history, and produced a corpus `clue validate` rejected.** The skill routes an agent into a failing validate and names no command. Escalated as **Q-02**; the register's own method — reading — would not have caught it, which is a limit stated again under *What this analysis does not establish*.

---

## Register — `clue-verify` body (VFY)

The pre-verification checklist is registered per checkbox; multi-condition checkboxes carry letter suffixes only where the conditions are independently checkable.

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| VFY-01a | Preamble — "Run this verification and review workflow before opening or updating any Cliewen PR." | rule | PDR-012 | **HUB-42** | yes | |
| VFY-01b | "Complete the local checks and agentic review loop before publishing; complete the hosted-head check immediately after publishing." | rule | PDR-016 | **F-RB-07** | yes | |
| VFY-01c | "Plain changes use only checks relevant to their changed surface and do not invoke this skill." | rule | PDR-011 | **HUB-15**, **F-CT-03** | yes | |
| VFY-01d | "Never fix a failure by weakening the check." | rule | C-004 | **HUB-48**, **F-RB-04** | yes | |
| VFY-02 | "The change uses the correct workspace under **Change scope and tiers** below." | rule | PDR-002 | **HUB-13**, **F-CT** | yes | |
| VFY-03 | "Every artifact touched has frontmatter `id`, `type`, `status`, `links`, `title`, plus decision `author`/`accepted-by`, constraint `source`/`enforcement`, capability `goal`, and any other type-specific fields." | rule | C-009 | no | **part (open-ended "any other")** | |
| VFY-04 | "Every `links` entry resolves to an existing ID." | rule | C-009 | no | yes | |
| VFY-05 | "The command name and the citation scheme are written in a code span wherever prose names them" | rule | ADR-040 | no | yes | |
| VFY-06 | "Every reference pointing outside this repository names what it points at … A bare forge number fails." | rule | ADR-040 | **ANL-06**, **ANL-07** | yes | |
| VFY-07 | "The proposal names a real plan item or explicitly declares the change plan-less." | rule | C-005 | **HUB-44**, **DLT-08** | yes | |
| VFY-08 | "Plan bookkeeping reflects the merge, and no completed plan changed. A change completing a plan's last milestone closes that plan `completed` in this same digest." | rule | log 2026-08-02, C-008 | **HUB-45**, **PLN-12** | part (3 conditions) | |
| VFY-09 | "Every active acceptance criterion satisfies its evidence contract: its identity uses the canonical grammar, a declared machine proof type has supported evidence classified by type and direction …" | rule | ADR-032, ADR-033, ADR-036, ADR-037 | **DLT-17…22** | **part (6 conditions in one checkbox)** | |
| VFY-10 | "Every `/docs/**` folder has a README; index blocks list every sibling artifact and no deleted file." | rule | C-016 | no | part (2 conditions) | |
| VFY-11 | "The change was assessed against every constraint (including verifiable quality bars)." | rule | ADR-027, ADR-045 | no | part (unbounded set) | |
| VFY-12 | "Repository-local conventions satisfy the contract below." | connective (pointer) | — | — | — | |
| VFY-13 | "Diagrams use the clearest renderable form: prefer Mermaid, use ASCII art where clearer, retain SVG where neither is adequate." | rule | C-007 | no | yes | |
| VFY-14 | "The full-change workspace is absent after digest; `main` never contains `/changes/`." | rule | PDR-002 | **HUB-28**, **HUB-29**, **DLT-26** | yes | |
| VFY-15 | "Every decision satisfies **Decision records** below, including routing, timeless content, provenance, objections, and pending approval signatures." | rule | C-011, C-006 | **F-DR** | part (5 conditions) | |
| VFY-16 | "The current commit received a pass with no blocking findings … every blocking repair after an earlier clean pass triggered a new pass …" | rule | PDR-012 | **VFY-24** | part (3 conditions) | |
| VFY-17 | "The final handoff identifies the review mode … and the reviewed commit." | rule | PDR-012 | **VFY-25** | yes | |
| VFY-18 | "Every review of an existing PR names its hosted head; actionable findings are unresolved hosted conversations where supported…" | rule | PDR-016 | **HUB-37a/b**, **F-RB-05** | yes | |
| VFY-19 | "Every intended edit, including each review fix, is committed and `git status --porcelain` is empty." | rule | C-012 | **F-RB-07** | yes | |
| VFY-20 | "`git merge-base HEAD origin/main` equals `origin/main` after fetching; no other change workspace is visible on this branch." | rule | C-012 | **F-RB-01** | part (2 conditions) | |
| VFY-21 | "After publishing, the current branch is the ready hosted PR's head branch, its head SHA equals local `HEAD` …" | rule | PDR-016 | **HUB-38**, **F-RB-07** | part (3 conditions) | |
| VFY-22 | "The branch and hosted PR satisfy the **Review boundary** below." | connective (pointer) | — | — | — | |
| VFY-23 | `## Agentic review loop` — "Run this loop automatically; never ask the human to clear context or initiate a separate review." | rule | PDR-012 | **DLT-31** | yes | |
| VFY-24 | 1 — "Finish every intended edit and commit the complete candidate. Determine the current commit and its base … then run the applicable local checks." | rule | PDR-012 | **F-RB-07** | part (3 conditions) | |
| VFY-25 | 2 — "If the host supports context-isolated delegation, start a new read-only reviewer without the implementation conversation." | rule | PDR-012 | no | yes | |
| VFY-26 | "Give it only the repository, branch, base, and declared intent…" | rule | PDR-012 | no | yes | |
| VFY-27 | "If isolated delegation is unavailable, perform an explicit adversarial pass … and label it `in-context fallback`; never describe that fallback as independent review." | rule | PDR-012 | no | yes | |
| VFY-28 | 3 — "Instruct the reviewer to inspect the complete base diff, durable corpus, tests, and constraints." | rule | PDR-012 | no | yes | |
| VFY-29 | "It returns no edits, only actionable findings about correctness, intent mismatch, regressions, security, missing evidence, or unjustified complexity." | rule | PDR-012 | no | yes | |
| VFY-30 | "Each finding includes severity, location, evidence, the operative requirement …, the concrete consequence, and a remediation." | rule | PDR-016 | no | **part (6 conditions)** | |
| VFY-31 | "Apply authoritative decisions and the repository's explicit lifecycle rules before treating nearby wording as contradictory…" | rule | PDR-012 | no | part | |
| VFY-32 | "Historical descriptions, optional activity, alternative readings, and lifecycle-correct state are not actionable defects by themselves." | rule | PDR-012 | no | yes | |
| VFY-33 | "Exclude stylistic preference, optional improvement, and speculative scope expansion." | rule | PDR-012 | no | yes | |
| VFY-34 | "Every finding is classified **blocking** or **advisory** and the reviewer's verdict is about the blocking set alone" + both definitions | rule | PDR-012 | no | yes | |
| VFY-35 | "When reviewing an existing hosted PR, bind the result to its observed head and ensure actionable findings become unresolved hosted review conversations" | rule | PDR-016 | **VFY-18**, **F-RB-05** | yes | |
| VFY-36 | 4 — "For every added or changed acceptance criterion, compare each scenario against its referenced tests' setup, action, and assertions." | rule | PDR-017 | no | yes | |
| VFY-37 | "Record an advisory verdict for the acceptance brief: `verifies`, `verifies-something-adjacent`, or `undetermined`." | rule | PDR-017 | no | yes | |
| VFY-38 | "This scenario-resolution result is not an actionable finding and does not gate `clue validate`." | rule | PDR-017 | no | yes | |
| VFY-39 | "A `Human`-class criterion has no test to compare — name it in the brief instead." | rule | ADR-033 | **DLT-19** | yes | |
| VFY-40 | 5 — "Resolve every blocking finding in the implementing context." | rule | PDR-012 | no | yes | |
| VFY-41 | "A finding that requires a new decision or changed intent becomes an open question and stops the change." | rule | C-011 | **HUB-46**, **DLT-14** | yes | |
| VFY-42 | "Otherwise the implementing context becomes the updater for that turn, follows the **Review boundary**, commits the repairs, and reruns applicable local checks." | rule | PDR-016 | **HUB-38**, **F-RB-06** | part (3 conditions) | |
| VFY-43 | "Advisory findings are carried in the verification evidence and repaired at the author's discretion … they never gate publication." | rule | PDR-012 | no | yes | |
| VFY-44 | 6 — "Start a new review pass after every blocking repair; a previous clean result applies only to the commit it reviewed." | rule | PDR-012 | **HUB-37c**, **F-RB-05** | yes | |
| VFY-45 | "Scope that pass to the diff since the reviewed commit plus the carriers those files declare, not to the whole change again." | rule | PDR-012 | no | yes | |
| VFY-46 | "Repairing an advisory finding does not invalidate a clean result, so the loop terminates." | rule | PDR-012 | no | yes | |
| VFY-47 | "Continue until the current commit receives a pass with no blocking findings. Do not publish with unresolved blocking findings or without such a pass." | rule | PDR-012 | **VFY-16** | yes | |
| VFY-48 | 7 — "Report the final review mode and reviewed commit with the verification evidence." | rule | PDR-012 | **VFY-17** | yes | |
| VFY-49 | "Context isolation reduces implementation anchoring but is not a substitute for human judgment or permission to merge." | rule | C-012 | **F-RB-03** | yes | **!** |

**`clue-verify` is where the checkability defect concentrates.** Most of its rows state more than one independent condition, and three state five or more (VFY-09, VFY-15, VFY-30 — the last a review-loop step rather than a checklist box). VFY-09 is the extreme case: one checkbox covering the identity grammar, a declared machine proof type with evidence classified by type and direction, JVM per-executable attachment, the Human-class brief line, the `@draft` exemption, and the legacy one-reference rule. No reader can honestly tick it when all but one of those hold, and a checklist item that cannot be ticked honestly is a checklist item that gets ticked dishonestly.

---

## Register — `clue-extract` body (EXT)

`clue-extract` is the largest carrier and the only one whose rules are almost entirely about a one-time operation. Its target contract is a numbered list whose items carry anywhere from one obligation to eleven; the register lists them per item with each item's obligation count in the checkability column, because splitting them into single obligations would produce a row set M-063 must consolidate back into a checkable form anyway. **That consolidation is the point** — the register's finding here is the shape, not the enumeration.

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| EXT-01 | `# clue-extract` | connective | — | — | — | |
| EXT-02 | "Brownfield adoption: transform an existing repository's specification corpus into a Cliewen `/docs` corpus." | rule | ADR-008 | no | yes | |
| EXT-03 | "Use once per adopted repository; the extraction is that repository's first `clue-delta` loop." | rule | ADR-008 | no | yes | |
| EXT-04 | "Apply the **Rehearsal before mutation**, **Decision records**, … below." | connective (pointer) | — | — | — | |
| EXT-05 | Rehearsal — "begin with a mandatory report-only pass … do not change the target source corpus, `/docs`, tests, routing, or hosted state" | rule | PDR-020 | no | part (5 protected surfaces) | |
| EXT-06 | "The rehearsal report inventories source formats and entry points, proposed artifact mappings, preserved and minted IDs, confidence and reversal cost, test-purpose work, instruction conflicts, planned deletions, and named plan doors." | rule | PDR-020 | no | **part (8 conditions)** | |
| EXT-07 | "An unresolved conflict becomes an `open-questions.md` entry and stops before mutation." | rule | PDR-024 | **HUB-46** | yes | |
| EXT-08 | "The rehearsal also writes a pinned source manifest … one proof-class, direction, and evidence-location row for every classified reference …" | rule | ADR-049, ADR-053 | no | **part (7 conditions)** | |
| EXT-09 | "The rehearsal also writes a pinned carrier inventory … `deleted-paths` … one row per operational carrier …" | rule | ADR-051 | no | **part (6 conditions)** | |
| EXT-10 | "Only explicit human direction begins the … mutate phase." | rule | PDR-020 | **F-RB-03** | yes | |
| EXT-11 | "That phase digests the rehearsal into the durable extraction report … the ready PR deletes both the transient workspace and the parallel source corpus." | rule | PDR-020 | **EXT-22** | yes | |
| EXT-12 | 1 — "The full taxonomy exists" (folders, indexed READMEs, capability files, extract-don't-invent, empty-but-indexed, row states its record) | rule | ADR-025, C-016, ADR-041, ADR-046 | **VFY-10** | **part (6 conditions)** | |
| EXT-13 | 2 — "Everything extracted is born inferred and cost-routed" | rule | ADR-010, ADR-035 | no | **part (5 conditions)** | |
| EXT-14 | 3 — "Existing criterion IDs survive; a criterion with none is minted deterministically" | rule | ADR-009, ADR-037 | no | **part (7 conditions)** | |
| EXT-15 | 4 — "Every test keeps or gains exactly one purpose without losing source evidence" | rule | ADR-006, ADR-036, PDR-024 | no | **part (11 conditions)** | |
| EXT-16 | 5 — "Evidence status is explicit at the narrowest honest level" | rule | ADR-032, ADR-033 | **DLT-18…20** | **part (6 conditions)** | |
| EXT-17 | 6 — "`clue validate` is green before the ready PR opens." | rule | C-012 | **VFY** (whole checklist) | yes | |
| EXT-18 | 7 — "`clue parity` is clean before the source corpus is deleted" | rule | ADR-049, ADR-053 | no | part (5 conditions) | |
| EXT-19 | 8 — "`clue carriers` is clean before the source corpus is deleted" | rule | ADR-051 | no | part (4 conditions) | |
| EXT-20 | 9 — "The source corpus dies in the same PR without losing its links" | rule | ADR-050, PDR-028 | no | **part (9 conditions)** | |
| EXT-21 | 10 — "Routing is rewritten and reconciled" | rule | ADR-013, PDR-022 | no | part (3 conditions) | |
| EXT-22 | 11 — "An extraction report lands in `/docs/analysis`, and its figures are rendered rather than typed" | rule | ADR-054, PDR-028 | no | **part (6 conditions)** | |
| EXT-23 | 12 — "Unsolved adoption items become named plan doors: Never leave a silent gap." | rule | PDR-024 | no | yes | |
| EXT-24 | 13 — "Every converted file carries exactly one frontmatter block" (incl. BOM handling) | rule | C-009 | no | part (3 conditions) | |
| EXT-25 | 14 — "The committed extraction receives a clean agentic review." | rule | PDR-012 | **VFY-01a**, **DLT-28** | yes | |
| EXT-26 | `## Source mappings` — "Per-source mappings live in this skill's `mappings/` folder" | rule | ADR-008 | no | yes | |
| EXT-27 | "The target contract above governs every extraction; a mapping only describes one source format." | rule | ADR-008 | no | yes | |
| EXT-28 | "A new source format adds a mapping file, never another skill." | rule | ADR-008 | no | yes | |
| EXT-29 | "If no mapping exists, writing one is the extraction PR's first task." | rule | ADR-008 | no | yes | |
| EXT-30 | `## Boundaries` — "Never invent unstated requirements, renumber or rename IDs, leave the source corpus alive, promote your own output to `verified`, or change test code beyond …" | rule | ADR-007, ADR-010, PDR-004 | **EXT-12**, **EXT-14**, **EXT-20** | part (5 prohibitions) | **!** |

**Order in `clue-extract`.** EXT-30 is the file's `## Boundaries` section and it is the last authored section before the shared fragments. It states five absolute prohibitions — including "never promote your own output to `verified`", which is the human merge boundary applied to extraction — after all fourteen contract items and the whole rehearsal procedure they constrain. This is the clearest order defect in the set: a section literally titled *Boundaries* placed after everything it bounds.

---

## Register — shared fragments

A fragment is registered once. The **Paths** column names which reading paths render it, because that is what turns one authored statement into repeated reading.

### `frontmatter` (F-FM) — rendered into all six

| ID | Locator | Class | Trace | Dup | Chk |
|---|---|---|---|---|---|
| F-FM-01 | `cliewen-skill: true` / `version:` block | rule (declaration) | ADR-022, ADR-011 | no | yes |
| F-FM-02 | "Generated from Cliewen's canonical skill sources; edit those sources, not this file." | rule | ADR-021 | **HUB-60** | yes |

### `change-tiers` (F-CT) — rendered into `clue-delta`, `clue-upgrade`, `clue-verify`

| ID | Locator | Class | Trace | Dup | Chk |
|---|---|---|---|---|---|
| F-CT-01 | "Classify scope before using the Cliewen loop." | rule | PDR-011 | **HUB-12**, **DLT-03** | yes |
| F-CT-02 | "Three rules decide the tier … take the first rule that matches." | rule | PDR-002 | **HUB-13** | yes |
| F-CT-03 | Tier 1, plain — definition, protected surfaces, and what plain work does | rule | PDR-011 | **HUB-14**, **HUB-15**, **VFY-01c** | part (9-item surface list) |
| F-CT-04 | Tier 2, light | rule | PDR-002 | **HUB-17**, **HUB-18**, **HUB-30** | yes |
| F-CT-05 | Tier 3, full | rule | PDR-018 | **HUB-19** | yes |
| F-CT-06 | Two guards — uncertainty escalates, discovery escalates immediately | rule | PDR-002 | **HUB-20**, **HUB-21** | yes |

**The whole tier section is stated twice to every reader of an ordinary full change**, once in `AGENTS.md` and once in `clue-delta` — and a third time for anyone who then opens `clue-verify`. This is duplication *within* a reading path, which ADR-021's standalone-skills reasoning does not cover: ADR-021 justifies the same fragment appearing in several *files*, not the hub and the skill both stating the rule to one reader in one session.

### `decision-records` (F-DR) — rendered into all six

| ID | Locator | Class | Trace | Dup | Chk |
|---|---|---|---|---|---|
| F-DR-01 | "Route every decision by reversal cost … log row … ADR … PDR." | rule | C-011, PDR-006, PDR-003 | **HUB-46** | yes |
| F-DR-02 | "A decision adopting a well-established practice cites it by name and records only the local why." | rule | PDR-006 | no | yes |
| F-DR-03 | "Agent-authored decisions start `status: inferred` and `author: agent`. Merging makes them binding without changing that status." | rule | PDR-004 | no | yes |
| F-DR-04 | "Only explicit human approval promotes a decision to `verified`; record every approver …, use the first approval date, and cite the venue." | rule | PDR-004 | no | part (3 conditions) |
| F-DR-05 | "An explicit objection keeps the decision `inferred` and becomes an open question." | rule | PDR-004 | **HUB-46** | yes |
| F-DR-06 | "`accepted-by:` records only approval given under Cliewen's merge boundary, never acceptance a source record already carried." | rule | ADR-029 | no | yes |
| F-DR-07 | "A record converted from a format with its own acceptance history preserves that history as body prose and keeps `accepted-by: []`." | rule | ADR-029 | no | yes |
| F-DR-08 | "Every decision record is timeless: state what is decided and only the enduring context and rationale." | rule | C-006 | no | part |
| F-DR-09 | "Keep triggering incidents, chronology, conversations, implementation details, and review history in findings, the workspace, the PR, and Git history." | rule | C-006 | no | yes |
| F-DR-10 | "A decision that changes a methodology contract inventories every live carrier … and updates that complete inventory in the same change." | rule | PDR-019, C-006 | no | part (unbounded set — C-006 says so itself) |
| F-DR-11 | "Live carriers include current corpus truth, canonical and generated skills, templates, public or contributor guidance, implementation explanations, CLI text, and distribution metadata." | rule (definition) | PDR-019 | no | yes |
| F-DR-12 | "Historical analyses, completed plans, and changelog entries remain pinned history." | rule | PDR-019 | no | yes |
| F-DR-13 | "Add focused guards for stable repaired claims, but do not present those anchors as proof that an arbitrary future carrier inventory is complete." | rule | C-006 | no | yes |

### `local-conventions` (F-LC) — rendered into `clue-delta`, `clue-extract`, `clue-upgrade`, `clue-verify`

| ID | Locator | Class | Trace | Dup | Chk |
|---|---|---|---|---|---|
| F-LC-01 | "apply the repository-local conventions declared in AGENTS.md, including digest requirements such as a user-facing changelog entry" | rule | ADR-013, C-002 | **HUB-51** | yes |
| F-LC-02 | "When a release adds or narrows a corpus obligation, preview and apply the supported `clue migrate` migrations before validating." | rule | ADR-039 | **UPG-13** | yes |
| F-LC-03 | "`clue init` remains a non-destructive materializer, not an updater." | rule | ADR-019 | no | yes |
| F-LC-04 | "Plain changes follow only the repository conventions that apply to their changed surface." | rule | PDR-011 | **HUB-15**, **F-CT-03** | yes |
| F-LC-05 | "Local conventions extend the methodology and never override it." | rule | ADR-013 | no | yes |
| F-LC-06 | "If AGENTS.md conflicts with a skill, record the conflict in `open-questions.md` and stop for a human decision; never choose silently." | rule | ADR-013 | **HUB-46** | yes |

### `durable-work` (F-DW) — rendered into `clue-delta`, `clue-upgrade`, `clue-verify`

| ID | Locator | Class | Trace | Dup | Chk |
|---|---|---|---|---|---|
| F-DW-01 | "An agent's private memory is never where work lives." | rule | PDR-017 | no | yes |
| F-DW-02 | "Anything needed to implement, continue, review, or hand off work belongs in a corpus artifact, the change workspace, or the pull request." | rule | PDR-017 | no | yes |
| F-DW-03 | "A durable record never states a figure a command computes — an artifact count, a coverage percentage, a reported population size. Name the command instead." | rule | **NONE** | no | yes |
| F-DW-04 | "A number written into prose becomes a hand-maintained obligation that goes stale … so the finding regenerates instead of converging." | connective (rationale for F-DW-03) | — | — | — |
| F-DW-05 | "Measurements that are the point of a record … are stated with what produced them and when." | rule | **NONE** | no | yes |

**F-DW-03 and F-DW-05 trace to nothing**, and they are not minor: F-DW-03 is a prohibition that shapes every durable artifact this project writes, including this one. The nearest live artifact is [ADR-054](../decisions/ADR-054-derived-extraction-report-region.md), which makes an *extraction report's* figures a rendered region — a much narrower rule about one mechanism in one document type, not a general prohibition on stating computed figures in prose. Escalated as **Q-03**.

### `review-boundary` (F-RB) — rendered into `clue-delta`, `clue-extract`, `clue-upgrade`, `clue-verify`

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| F-RB-01 | ¶1 — branches from `main`; one change in flight per author; parallel independent authors; plain changes exempt; review does not mint a change; stacking is a blocking open question; rebase before publication; merge `main` in after publication | rule | PDR-007, C-012 | **HUB-16**, **HUB-27**, **HUB-31**, **HUB-32**, **HUB-35**, **HUB-40** | **part (8 obligations, one paragraph)** | |
| F-RB-02 | ¶2 — merge-commit mode; configure the host; rebasing unpublished branches allowed; unsupported forges are outside the path | rule | PDR-021 | **HUB-34** | part (4 conditions) | |
| F-RB-03 | ¶3 — "Open the PR ready … never as a draft"; PR is authorization not duplicate review; unfinished work stays on the branch; **"An agent never merges its own PR, creates a local merge commit into `main`, or pushes to `main`."** | rule | C-012, PDR-007 | **HUB-33**, **UPG-16**, **VFY-49** | part (4 conditions) | **!** |
| F-RB-04 | ¶4 — hosted CI is displayed not enforced; branch protection makes the check a precondition; "Never weaken the workflow or required-check policy to make a change pass." | rule | PDR-027, C-004 | **HUB-48**, **VFY-01d** | yes | |
| F-RB-05 | ¶5 — review bound to observed head; clean result binds that commit; findings are durable PR state; disclose an enforcement gap rather than claiming equivalence | rule | PDR-016 | **HUB-37a/b/c**, **VFY-18** | part (4 conditions) | |
| F-RB-06 | ¶6 — the updater handoff: fetch, record head, recheck, push fast-forward only, reconcile on rejection, stop on a merged PR | rule | PDR-016 | **HUB-38…41**, **VFY-42** | **part (6 conditions)** | |
| F-RB-07 | ¶7 — "Ready means the hosted PR contains the exact locally reviewed and verified state" + the six-step publication sequence + the local-stopping-point clause | rule | PDR-016, C-012 | **HUB-36**, **HUB-38**, **HUB-42**, **HUB-43**, **VFY-19…21** | **part (7 conditions)** | |
| F-RB-08 | ¶8 — an agent stops before initiating another change; review fixes stay on the same branch; a follow-up exists only when a human scopes it | rule | PDR-007 | **HUB-32** | part (3 conditions) | |
| F-RB-09 | ¶9 — "After a human reports the merge, orient before starting anything else: describe the plan's next unfinished step … or say the plan has nothing left and ask what comes next." | rule | **NONE** | no | yes | |

**`review-boundary` carries the single most important prohibition in the method and states it in the third paragraph of the last section of four skills.** F-RB-03's final sentence — an agent never merges its own PR — is the human merge boundary, one of the three core elements behind [ARCH-003](../architecture/core.md)'s red line. In `clue-delta` it is roughly 85% of the way through the file; in `clue-verify` and `clue-extract` it is similar. The hub states it early (HUB-33), which is the only reason an agent reading in the intended order meets it before acting. An agent that opens a skill directly does not.

**F-RB-09 traces to nothing.** The orient-after-merge instruction is real practice and is stated in no decision, constraint, goal, or criterion. Escalated as **Q-04**.

---

## Populations

Counted over the seven carriers at `e9a7d07`, fragments counted once:

These are this spike's own results. **They are the only figures in this document**, and they are not hand-maintained: the script below recomputes every one of them from the register rows, and a reader who disagrees with a number runs it rather than trusting the prose. Everything outside this table describes shape and names statements by ID; where a count would once have been written into a sentence, the sentence now points here.

```awk
# awk -f recount.awk AN-018-skill-statement-register.md
BEGIN{FS="|"}
/^\| (HUB|ANL|PLN|UPG|DLT|VFY|EXT|F-[A-Z][A-Z])-/ && NF>=8 {
  total++
  class=$4; trace=$5; dup=$6; chk=$7
  gsub(/[* ]/,"",class); gsub(/[* ]/,"",trace); gsub(/[* ]/,"",dup); gsub(/[* ]/,"",chk)
  if (class ~ /^connective/) { connective++; next }
  rule++
  if (dup != "no") duplicated++
  if (chk ~ /^part/) uncheckable++
  if (trace ~ /^NONE/) untraced++; else if (trace ~ /^ARCH/) architecture++
  if (NF>=9 && $8 ~ /!/) misordered++
}
END{
  printf "total %d  connective %d  rule-bearing %d
", total, connective, rule
  printf "traceable %d  architecture %d  untraced %d
", rule-untraced-architecture, architecture, untraced
  printf "duplicated %d  uncheckable %d  misordered %d
", duplicated, uncheckable, misordered
}
```

| | Statements |
|---|---|
| Total registered | 279 |
| Connective | 19 |
| Rule-bearing | 260 |
| — traceable to a live decision, constraint, goal, or criterion | 252 |
| — traceable to an architecture artifact ([PDR-031](../decisions/PDR-031-architecture-artifacts-are-traces.md)) | 3 |
| — **traceable to nothing found** | 5 |
| Rule-bearing statements duplicated in at least one reading path | 124 |
| Rule-bearing statements failing checkability (`part`) | 63 |
| Statements binding absolutely but read after what they constrain | 4 |

The trace rows are stated **after** the answers below, and after an adversarial review corrected traces the first pass got wrong in both directions. Rerunning the script against an earlier revision reproduces what it said before; the numbers are not restated here, because a historical figure typed into prose is exactly what this document stopped doing.

Every statement that still traces to nothing has an answer, and the answers split two ways. **Three await their carrier:** F-DW-03 and F-DW-05 trace to the constraint Q-03's answer mints in M-067, F-RB-09 to Q-04's, and neither constraint can be written before the prose it points at exists. **Two are withdrawn:** PLN-17 and DLT-33 were asked what they were for, the answer was *nothing that survives the question*, and M-063 removes them. No rule-bearing statement in these carriers is left both live and unaccounted for.

Two rows deserve their qualification rather than their number. **Duplication** counts rule-bearing statements with at least one duplicate on the `AGENTS.md` → `clue-delta` → `clue-verify` path or the shorter paths through `clue-upgrade` and `clue-extract`; a different path set yields a different total, and the path set is stated above rather than assumed. **Checkability** counts every statement offering more than one independent condition as one obligation. That threshold is two, which is deliberately strict: a two-condition statement is usually fine to read and only awkward to tick, so M-063 should treat that row as a ranked list rather than a defect total.

The shape matters more than any of these numbers. **Tracing is not Cliewen's largest problem — nearly every rule-bearing statement traces, and the ones that do not are spread across the routing hub, two skills, and two shared fragments rather than concentrated in a rotten carrier.** Duplication and checkability are the problem, and those *are* concentrated: the tier rules and the review boundary carry most of the duplication, and `clue-verify` and `clue-extract` carry the largest share of the checkability failures between them, though not a majority. The trace column is the one this pass got wrong twice — first by resting rules on a frozen document, then by claiming traces the named artifacts do not carry — so it is the column a second reader should distrust most.

## Order assessment

These statements bind absolutely and are read after the procedure they constrain:

| Statement | Carrier | Position |
|---|---|---|
| F-RB-03 ("An agent never merges its own PR…") | `clue-delta`, `clue-verify`, `clue-extract`, `clue-upgrade` | last section, third paragraph |
| UPG-16 ("Never merge it") | `clue-upgrade` | last numbered step |
| VFY-49 ("not a substitute for … permission to merge") | `clue-verify` | last line of the review loop |
| EXT-30 (`## Boundaries`, a list of prohibitions) | `clue-extract` | after every contract item |

The shared-fragment order is one authored decision with one repair: the fragments render in the order the skill template lists them, and every skill lists `review-boundary` last. `AGENTS.md` is the counter-example that shows the repair works — it states the merge prohibition as its first binding rule, and nothing about that placement made the file harder to read.

## Compatible overlap candidates for M-063

These are pairs that cover one situation without pulling a reader in different directions. They are consolidation work, not escalations.

1. **The tier section** (HUB-12…21 against F-CT-01…06). One rule set, two statements of it, on one reading path. The hub's version names this repository's protected surfaces; the fragment's names generic ones. Consolidating means deciding which layer owns the surface inventory — likely the hub, under [ADR-013](../decisions/ADR-013-ships-generic-vs-repo-local.md)'s generic/repo-local split.
2. **The review boundary** (HUB-27…43 against F-RB-01…08). The hub restates most of the fragment in compressed form, and that compression is where HUB-38's unreadably conditional sentence came from.
3. **The workspace lifecycle** (HUB-28, HUB-29, DLT-07, DLT-26, VFY-14). One rule — the workspace is created on the branch, deleted in the digest, and absent from `main` — stated once per listed row.
4. **The plan-item declaration** (HUB-44, DLT-08, VFY-07). C-005, restated once per listed row.
5. **Plan closing** (HUB-45, PLN-11, PLN-12, PLN-13, PLN-14, VFY-08). One decision-log row, restated once per listed row — and that log row is itself the record of a change (CH-102) that existed only because the carriers had stated the rule too generically. Consolidating here is repairing the repair.
6. **Never weaken a check** (HUB-48, VFY-01d, F-RB-04). C-004, restated once per listed row, all on one path.
7. **Open questions stop the work** (HUB-46, DLT-14, VFY-41, EXT-07, F-LC-06, F-DR-05). C-011's stopping rule, restated once per listed row and each scoped slightly differently — the widest reuse of one rule in the register, tied with candidate 5.
8. **The criterion evidence contract** (DLT-17…22 against VFY-09 against EXT-16). Stated as procedure in `clue-delta`, as one checkbox in `clue-verify`, and as a contract item in `clue-extract`. This is the pair with the highest consolidation value and the highest risk, because the three statements are not word-identical and the differences may be meaningful.

## Escalations

Q-01 through Q-08. Each names the statement, what it traces to or fails to, what removing or retaining it costs, and what judgment is required. Q-01 through Q-07 were raised by the first pass. Q-08 reached the list two ways — some of its statements surfaced while re-deriving the populations against the first seven answers, and the rest came from an adversarial review of this register. All eight are answered below.

**Q-01 — `clue-analysis`'s workflow spine traces to nothing the method accepts.** ANL-02…05, ANL-17…19, ANL-21, ANL-25…27. *Class: traces to nothing.*

**Q-02 — `clue-delta`'s CH allocator (DLT-06) instructs a method the shipped tool has replaced.** ADR-048 superseded it and its mechanism is live; this change followed the skill and failed `clue validate` as a result. *Class: the decision has outlived its reason.*

**Q-03 — the no-computed-figures rule (F-DW-03, F-DW-05) traces to nothing.** *Class: traces to nothing.*

**Q-04 — the orient-after-merge instruction (F-RB-09) traces to nothing.** *Class: traces to nothing.*

**Q-05 — does an architecture artifact count as a trace?** HUB-25 and HUB-56 trace to [ARCH-003](../architecture/core.md) and HUB-59 to [ARCH-002](../architecture/skills.md), which PDR-029's four types exclude. *Class: definitional; blocks the register's own verdict on the statements above.*

**Q-06 — plain changes and unmerged work (HUB-16 against F-RB-01).** The hub says plain changes *never* build on unmerged work; the fragment permits building on unmerged work with explicit human authorization and does not exempt plain changes. The first pass registered this as an undecided conflict, having mis-traced HUB-16 to PDR-007; review found that [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) states the hub's rule verbatim. It is therefore a carrier that drifted from a verified decision rather than two obligations with no adjudicator, which makes the repair cheaper and its direction already settled.

The answer to this question ranged wider than the question. Asked about the plain-change clause, the human also asked that a suggestion raised mid-change be recorded rather than remembered, and that work depending on an unmerged change be recorded in the repository until the merge happens. Those are two further rules, not readings of this one: the first became [PDR-032](../decisions/PDR-032-mid-change-suggestions-are-triaged.md), and the second is [AN-013](AN-013-distributed-work-and-evidence-boundaries.md)'s F1, routed to [P-013](../plans/P-013-simplification.md)'s M-065. They are recorded here because they arrived here, not because Q-06 asked for them. *Class: recorded as a conflict, resolved as carrier drift.*

**Q-07 — `clue-extract` is missing from the hub's skill table (HUB-59).** The row HUB-59 traces to, [ARCH-002](../architecture/skills.md), carries a section headed "The six skills and how they complement each other" and tables all six. The hub omits `clue-extract`. *Class: coverage, raised here because the fix is a rule-bearing row.*

**Q-08 — statements this register failed to account for.** They failed three different ways, and the third was found only when the answers were being written.

Two were recorded as untraceable and never escalated, because the first pass grouped escalations by carrier and each stood alone in its file: PLN-17 (distil a completed plan's lessons before freezing it) and DLT-33 ("Keep deltas small: Git merges text, not meaning").

Four were recorded as *traced* to artifacts that do not state their rule, and an adversarial review caught them: HUB-22, HUB-23 and HUB-24 — the read protocol — were traced to AC-053, whose scenario is about what `clue context` emits and says nothing about when to read `docs/README.md` or where to stop reading; PLN-06 was traced to C-010, which stated only the milestone status vocabulary.

One was the opposite error. DLT-15, the pause after Propose, was recorded as tracing to nothing — and [PDR-017](../decisions/PDR-017-merge-gate-has-content.md) had decided it all along, in one line inside a record about the acceptance brief. A review caught that too, while checking the record written to fill the supposed gap. *Class: traces to nothing, plus one false negative.*

**Answered.** Each of the seven was put to the human as *what is this rule for?*, and the answers split three ways.

*Two are withdrawn.* PLN-17 asked every campaign to distil its lessons before freezing, and eight consecutive campaigns closed without doing it at no cost — the lessons that mattered had already been recorded as decisions when they were learned. DLT-33 is true advice that binds nothing: no reader can determine whether a delta was small enough. Both are removed from their skills by M-063, each on a decision-log row stating what the removal gives up.

*One traced all along and is now extended.* DLT-15 traces to [PDR-017](../decisions/PDR-017-merge-gate-has-content.md), which decided the pause and settled that it is opt-in and recorded in tasks. What PDR-017 does not settle is what the pause is *for*, and read without that it is a convenience — which is how it has been used. Its purpose is that proposing and implementing are different work, and the boundary between them is where a change can be split, redirected, or handed to a different agent, which is why the proposal is committed first. [PDR-033](../decisions/PDR-033-planning-and-implementation-are-separate-steps.md) amends PDR-017 with that purpose and with what the skill never said: at the pause the agent reports briefly and asks two questions — whether implementation begins, and whether the branch is pushed. Pushing is what makes handoff possible and what ends the branch's freedom to be rebased, so it is the human's call rather than a default.

*Four were real rules with no record.* The read protocol (HUB-22…24) is corpus-reading economy: start at the narrowest point that answers the question, widen on a discovered edge. [CAP-007](../capabilities/CAP-007-focused-context/README.md) built the mechanism and nothing ever obliged anyone to use it, so the expensive default — read defensively, decide afterwards what mattered — never looked wrong. [PDR-034](../decisions/PDR-034-the-corpus-is-read-narrowly-by-default.md) states the obligation and adds no mechanism. PLN-06's verifiable exit criterion now lives in [C-010](../constraints/C-010-milestone-status-vocabulary.md) beside the status vocabulary, which moves to `enforcement: partial`: the judge reads the status cell and will never judge whether a promise can be checked, and the constraint says so rather than overclaiming.

## Answers

Q-01 through Q-07 were answered by Flemming N. Larsen on 2026-08-08 in conversation, and recorded under [C-011](../constraints/C-011-decision-records-typed.md) in the same change that wrote this register: [PDR-030](../decisions/PDR-030-analysis-is-a-bounded-spike.md) (Q-01), [PDR-006](../decisions/PDR-006-decision-records-are-typed.md)'s amendment (Q-01's second half), [PDR-031](../decisions/PDR-031-architecture-artifacts-are-traces.md) (Q-05), [PDR-032](../decisions/PDR-032-mid-change-suggestions-are-triaged.md) (Q-06's suggestion half), and rows in [`log.md`](../decisions/log.md) dated 2026-08-08 for Q-02's routing, Q-03, Q-04, Q-06's plain-change clause, and Q-07.

Four answers reached past the questions, and each landed on machinery the corpus already had. Q-04 generalised to *a reference names what it is about*, which is [ADR-046](../decisions/ADR-046-index-rows-say-what-the-artifact-is-about.md)'s index-row rule applied to agent-facing prose. Q-06's suggestion half reuses [ADR-002](../decisions/ADR-002-inbox-is-proposed-goals.md)'s inbox. Q-03 asked for a check that constraints are applied, which lands on VFY-11 — a defect this register had already found by a different route. Q-07 asked for a guard against forgetting, answered in two layers because the hub is Cliewen's file here and the adopter's file there.

**Q-06's second half is [AN-013](AN-013-distributed-work-and-evidence-boundaries.md)'s F1.** The request that unmerged dependent work be recorded durably in the repository is the finding P-009 declined by removing M-043, and [P-013](../plans/P-013-simplification.md)'s M-065 exists to give it a determination. What M-065 lacked was a human judgement that the mechanism is worth building; that judgement has now been given, so M-065 builds rather than repeating "still open". This register does not design the record — M-065's own change does.

**Q-08 is answered.** M-063 removes PLN-17 and DLT-33, and rewrites DLT-15 to carry PDR-033's report-and-ask step. Nothing in this register is left open.

## Recommendation on the register's durable form

The milestone asks what durable form the register should take, and names the trap: a map kept beside the prose it describes, with no mechanism, is the second stored representation [PDR-028](../decisions/PDR-028-derived-report-is-not-a-committed-registry.md) refused. It goes stale on the first skill edit, and nothing reports it.

**Recommendation: keep it as this pinned analysis artifact, and add no mechanism now.** Three reasons, in order of weight.

**The register's job ends when M-063 lands.** It exists to let M-063 trim on merit rather than by argument. A register that survives as a live index has to answer a question nobody has asked yet — *which rule does this sentence carry, today* — and building the machinery for that question before M-063 has demonstrated it is worth answering is exactly the cost this campaign is supposed to be reducing.

**A `clue`-enforced register would change what a green `clue validate` asserts, and that is core-adjacent under [C-013](../constraints/C-013-core-changes-need-decision.md).** Today green means the corpus is well-formed. A register check would make green also mean *every rule-bearing sentence in every carrier names a live artifact* — a judgement about prose, not form, and one that requires the judge to segment English. PDR-029 already forbids putting the marker in the carrier, so the binding would have to be by quotation or hash, and both fail the same way: an editorial rewording that changes no rule invalidates the row and demands a restatement, which teaches authors to avoid editing carriers. That is a worse outcome than a stale map.

**PDR-028 refused the analogous thing for a stronger case.** A per-criterion extraction registry had a mechanical binding available — IDs — and was still refused as a second stored representation that ages. This register's binding would be to prose, which is weaker.

**What a pinned register costs, stated plainly.** A reviewer reading this file in a year meets a map of a carrier set that has moved. They cannot tell which rows still hold without redoing the walk. That cost is real and it is the reason this recommendation is not free — the mitigation is that the file states its pinned revision in its own evidence boundary, so a later reader knows immediately that they are reading history and how to re-derive it.

**If the human wants the register bound anyway**, the cheapest honest form is a re-runnable spike rather than a committed registry: M-063's evidence already re-runs this walk, and a plan milestone requiring a re-run whenever a carrier changes gets most of the value with none of the second-representation cost. That would be a plan revision with its own decision record, which this milestone recommends and does not build.

**A fourth reason arrived after the recommendation was written, and it is the strongest of them.** The first draft of this document stated its figures in prose in roughly thirty places. The review loop then ran four passes; three of them found nothing but arithmetic — prose saying one total where the rows said another — and each repair wrote new prose carrying new numbers, so the next pass found more. That is the exact failure the shared `durable-work` fragment predicts of a hand-maintained figure, and this register both documented the rule and broke it. The repair was to delete every derived number from the prose and leave the populations table with the script that recomputes it, which is why this file now names statements by ID and describes shape in words.

The lesson generalises past this document. **A register whose numbers are authored is a register that manufactures review findings indefinitely**, and giving that register a mechanism would multiply the surface rather than settle it. A representation that must be kept consistent by hand is not made safe by being checked more often; it is made expensive. The measurement belongs in one place, derived, and the prose around it belongs in prose. That argues against a mechanism more strongly than any of the three reasons above, and it was learned by paying for it.

## Rejected approaches

**Segmenting by sentence alone.** Tried first and abandoned within the hub. It splits qualified obligations from their conditions — HUB-14's protected-surface list becomes a statement with no rule, and the rule it qualifies becomes uncheckable — and it treats every rationale clause as a rule-bearing statement that traces to nothing, which would have buried the untraceable population in noise.

**Registering the rendered skills rather than the sources.** It would have counted each shared fragment once per skill that renders it, and reported a duplication figure that ADR-021 already explains and accepts. The reading-path rule exists because the honest quantity is how many times one *reader* meets one rule, and that is neither the file count nor one.

**Scoring the carriers by word count.** Refused by PDR-029 before this spike began, and the register confirms why: `clue-verify`'s worst defect is VFY-09, whose repair makes the file longer.

**Deciding the escalations.** Most have obvious-looking answers — write the missing decisions for Q-01, Q-03, and Q-04; delete DLT-06 for Q-02; add the row for Q-07. PDR-029 forbids exactly this, and its reason applies here: a statement that traces to nothing may be a real rule nobody recorded, and an agent that resolves it silently has made a methodology decision without a human.

## What this analysis does not establish

One reading of one revision by one agent. The register has not been independently re-segmented, so the segmentation rule's central claim — that a second pass produces the same statements — is untested. The trace column's `NONE` verdicts rest on searches for distinctive wording; a rule stated in a live artifact using entirely different words would have been missed, and the failure mode is asymmetric, producing false `NONE`s rather than false traces. Duplication is counted over the reading paths named above and no others. Nothing here measures whether an agent that reads these carriers actually obeys them, which is the question [PDR-023](../decisions/PDR-023-tool-notice-and-hub-instruction.md) had to answer with a live acceptance test rather than by inspection.

**The method has one demonstrated blind spot.** Reading a carrier against the corpus finds statements that trace to nothing and statements that state one rule twice. It does not find a statement that traces cleanly to a decision whose *mechanism has since shipped and changed the answer* — DLT-06 reads as correct and traces to a real ADR, and only executing it exposed the defect. Everything in this register is a reading; a second pass that ran each instruction rather than reading it would probably find more of these, and no claim is made here that DLT-06 is the only one.

**The trace column was wrong in every direction available to it, and the false-trace direction is the dangerous one.** The first pass rested part of `clue-analysis`'s spine on a frozen document, and it claimed traces the named artifacts do not carry: to an acceptance criterion about a different subject, to a constraint stating only a status vocabulary, to a decision that does not mention the case at all, and to an architecture ID that does not exist. A false `NONE` costs an unnecessary escalation. A false trace costs the opposite and worse — it retires a question nobody will ask again, and M-063 would have defended those statements on evidence that is not there. The asymmetry stated above, that this method produces false `NONE`s rather than false traces, was itself wrong. It was also wrong the third way, once: DLT-15 was recorded as tracing to nothing while PDR-017 had decided it, in one line inside a record about something else — a false `NONE` produced by searching for the rule's wording rather than for the rule. Every entry in this column is one reader's judgement that an artifact states a rule, and the only thing that caught any of these was a second reader opening the artifacts.

## Consumer

[P-013](../plans/P-013-simplification.md)'s **M-063**, which trims and reorders against this register and its answers, and **M-067**, which carries the mechanisms the answers asked for. M-063 is unblocked: every escalation this register raised is answered, and the trim may proceed in full. [M-065](../plans/P-013-simplification.md) consumes Q-06's second half.
