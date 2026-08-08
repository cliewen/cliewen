---
id: AN-021
type: analysis
status: active
links: [P-013, M-064, PDR-029, PDR-031, PDR-019, PDR-035, AN-018, C-011, C-017]
title: Statement register for the public guide, contributor guidance, and CLI text
---

# AN-021 — Statement register for the public guide, contributor guidance, and CLI text

## The risk this spike retires

[AN-018](AN-018-skill-statement-register.md) registered the seven carriers an agent reads. It deliberately left three carriers unread: the public guide, `CONTRIBUTING.md`, and the CLI's own text. Those are what a *human* reads — an evaluator deciding whether to adopt Cliewen, a contributor deciding how to work, a user reading `clue --help` — and [PDR-029](../decisions/PDR-029-simplification-tests-by-surface.md) puts them under the same test as the skills.

The risk is asymmetric with AN-018's. There, the danger was trimming prose without knowing what each statement was for. Here the danger is that a rule was repaired in the skills and left standing in the guide. A carrier nobody registered is a carrier whose contract changes silently, and a methodology whose public guide states a rule its skills have withdrawn is arguing against itself in the one place a reader is most likely to check.

Nothing here is removed, reworded, or reordered. No trace, citation, or marker was added to any carrier.

## Evidence boundary

- **Pinned revision:** the carriers were read at `9a632f9` on `main` (`cliewen/cliewen`) — the tip this change branched from. No carrier is edited by the change that writes this register, so every locator below is to that revision and to the working tree alike.
- **What was read:** the twelve files of the public guide (`guide/*.md`, including `index.md`), `CONTRIBUTING.md`, and the CLI's user-facing text — the `usage` constant in `cmd/clue/main.go`, the post-`init` hint in `cmd/clue/init.go`, and the diagnostics under `internal/` that instruct a reader rather than merely reporting a state.
- **What was not read:** the corpus itself, the skills and routing hub (AN-018 holds those), `README.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, the issue and pull-request templates, and the guide's VitePress configuration except for its sidebar, which was read to establish the reading paths below and is not registered. The scaffolded adopter hub belongs to AN-018's surface and is not re-registered here.
- **Environment:** prepared, not clean. Windows 11, Go 1.26.5, a `clue` binary built from this checkout and stamped `dev`. Nothing in the register depends on the toolchain: every result is a reading of committed text. The one derived figure that is not a reading — the inferred-decision population quoted in the pattern C determination — belongs to [AN-022](AN-022-remaining-surface-scored.md), not here.
- **Confidence classes.** Class, duplication, checkability, and order are **observed**: a second reader confirms or refutes each by opening the file. A trace is observed when the named artifact contains the rule; the judgement that it is the *narrowest* such artifact is an **inference**. AN-018 learned that the trace column is the one to distrust, in both directions, and this register reduces that exposure rather than eliminating it — see *Where the traces came from* below.
- **Not established.** One reading by one agent, not independently re-segmented.

## Where the traces came from

Most rule-bearing statements in these carriers state a rule AN-018 already traced, because the guide explains to a human what the skills instruct an agent to do. Where that is so, **this register reuses AN-018's trace and names the skill statement that carries the same rule**, rather than re-deriving the trace from the corpus. That is deliberate and it cuts both ways: it keeps the two registers consistent and it inherits AN-018's errors. A trace this register derived fresh — for a rule that exists only in these carriers, such as installation, CI configuration, or the plugin boundary — is marked with a dagger (†) so a second reader knows which column entries have no prior pass behind them.

## What counts as a single statement

[AN-018's segmentation rule](AN-018-skill-statement-register.md) applies unchanged, and is not restated here. Two extensions were needed for these carriers, both stated so an independent second pass can apply them the same way:

1. **A maximal run of adjacent connective statements inside one section is registered as one row.** Its ID spans the range (`GD-06…GD-14`) and its class cell carries the run's length (`connective ×9`). The guide is majority argument and explanation; registering each sentence of it separately would have produced a document whose bulk was rows carrying no obligation. A run is broken by any rule-bearing statement, by a heading, and by a section boundary, so no rule hides inside one.
2. **CLI text segments by output unit.** A usage line, a command's description paragraph, a flag's description, and one emitted diagnostic are each one statement. A diagnostic's format placeholders are part of it.

## Reading paths

Duplication is counted per reading path, as PDR-029 requires. These carriers have four, and they are not the agent's:

- **P1 — the adopter's path.** The guide's own sidebar section *Start here*, in its declared order: `index` → `what-is-cliewen` → `getting-started` → `plugin` → `adoption` → `ci-wall` → `operations`.
- **P2 — the evaluator's path.** The sidebar section *How the method works*: `design` → `methodology` → `corpus` → `change-loop` → `skills`. Each page's closing *Next* link joins the two sections, so a reader following *Next* from `what-is-cliewen` enters P2 and returns to P1 from `skills`; a rule stated once on each path is met twice by that reader and is scored as duplicated only where both statements sit on one path.
- **P3 — the contributor's path.** `CONTRIBUTING.md` → `AGENTS.md` → the skill the hub routes to. This path crosses into AN-018's surface, so a `CONTRIBUTING.md` statement that restates a hub or skill statement **is** scored as duplicated, and the duplicate's AN-018 ID is named.
- **P4 — the CLI.** `clue --help` and the diagnostics one command emits.

**A guide statement and a skill statement are not on one path** unless P3 joins them. A human reading the public guide is not the agent reading `clue-delta`, and scoring that repetition would produce a large figure describing nothing — the same error AN-018 refused when it declined to count a shared fragment once per skill. This is the register's most consequential judgement and a second reader who rejects it gets a different duplication population; the rule is stated here so the disagreement is legible rather than buried.

## What each column means

As in AN-018: **Class** (*rule* when the statement says what a reader must, must not, or may do, or states a condition a corpus state must satisfy; *connective* otherwise), **Trace** (the narrowest live artifact that *states* the rule, `NONE` when none was found, † when derived by this pass), **Dup** (a live carrier stating the same rule on the same reading path), **Chk** (whether a reader can determine satisfaction; `part` marks several independent conditions offered as one obligation), **Ord** (`!` marks a statement that binds absolutely but is read after the procedure it constrains).

One class judgement recurs often enough to state once. **A description of what Cliewen or `clue` does, addressed to someone deciding whether to adopt it, carries no obligation and is connective.** "The `clue` CLI checks structure, links, and traceability without executing tests" tells a reader what to expect, not what to do. The same sentence in a skill would be a rule about what the agent may rely on; in the guide it is mechanism description, which AN-018's class definition already places outside the three conditions. Without this rule most of the guide would be scored as uncheckable rules, which would be a measurement artefact rather than a finding.

---

## Register — `guide/index.md` (GI)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GI-01 | frontmatter hero block — name, text, tagline, two actions | connective ×1 | — | — | — | |
| GI-02…GI-04 | feature cards 1–3 — audience, the thread's shape, the three roles | connective ×3 | — | — | — | |
| GI-05a | feature 4 — "Agents prepare the corpus and verified proposal; humans keep control of intent and merge." | rule | C-012 | no | yes | |
| GI-05b | feature 4 — "Small work that changes no meaning stays outside the full loop." | rule | PDR-011 | no | yes | |
| GI-06…GI-07 | `## Next` and its link | connective ×2 | — | — | — | |

---

## Register — `guide/what-is-cliewen.md` (GW)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GW-01…GW-02 | definition and the *cliewen* etymology | connective ×2 | — | — | — | |
| GW-03 | "the durable documentation describes the system as it exists, not a pile of past change requests" | rule | ARCH-003 | no | part (no test a reader can apply) | |
| GW-04 | "A goal leads to a capability, a capability owns acceptance criteria, and each active criterion reaches its declared acceptance evidence." | rule | ARCH-003 | no | yes | |
| GW-05 | "Machine-proven criteria use supported, classified test references; genuine Human-class criteria use the pull request acceptance brief." | rule | ADR-032, ADR-033 | no | yes | |
| GW-06…GW-12 | "The `clue` command checks that this thread is intact"; `## Why another workflow?`; the acceptance-bottleneck argument; "Cliewen separates mechanical checks from human judgment" | connective ×7 | — | — | — | |
| GW-13 | "The corpus under `/docs` is the system of record." | rule | ARCH-003 | no | part | |
| GW-14 | "A branch is a proposal, and a pull request is the authorization boundary: the agent may publish a candidate but cannot accept it into `main`." | rule | C-012 | no | yes | |
| GW-15 | "A full change keeps its working delta in `/changes/CH-xxx-*`; the digest deletes that workspace before merge." | rule | PDR-002 | no | yes | |
| GW-16 | "The `clue` CLI checks structure, links, and acceptance-evidence traceability without executing tests." | connective ×1 | — | — | — | |
| GW-17 | "A human controls acceptance by merging; this safeguard does not require repeating a code review already completed locally." | rule | C-012 | **GW-14** | yes | |
| GW-18 | "The pull request is also where hosted CI becomes enforceable when the repository requires its status check and protects `main`." | rule | PDR-027 | no | yes | |
| GW-19 | "A pull request without a required check and branch protection only displays CI; the combination is what prevents an agent from silently skipping the gate." | rule | PDR-027 | **GW-18** | yes | |
| GW-20…GW-27 | `## Born from Intent Engineering…`; the book's provenance; the OpenSpec comparison and the no-archive-step argument | connective ×8 | — | — | — | |
| GW-28 | "every Cliewen change is required to leave it true" | rule | ARCH-003 | **GW-03** | part | |
| GW-29 | "The pull request authorizes that merge but is not the system of record, so squash and rebase-and-merge are outside the full-change support boundary." | rule | PDR-021 | no | yes | |
| GW-30 | "A repository already using the book's extended OpenSpec format can be adopted with its IDs and test traceability intact" | connective ×1 | — | — | — | |
| GW-31 | "A decision an agent records during a change is born `inferred`; merging the pull request makes it binding, and a later explicit human approval promotes it to `verified`." | rule | PDR-004 | no | yes | |
| GW-32 | "Shipping never blocks on an approval ritual, and decisions no human has yet endorsed stay visible separately." | connective ×1 | — | — | — | |
| GW-33 | "Extracted non-decision meaning is classified by reversal cost: cheap inferred findings may remain deferred, while expensive inferred meaning cannot silently support an active capability." | rule | ADR-035 | no | yes | |
| GW-34…GW-36 | the two-failures summary; `## What Cliewen is not`; the not-an-issue-tracker and not-a-test-runner boundaries | connective ×3 | — | — | — | |
| GW-37 | "Canonical criterion IDs use `<PREFIX>-<digits>[lowercase-suffix]`, so brownfield identities such as `SNAP-SQS-001` and `ADP-045b` remain stable" | rule | ADR-037 | no | yes | |
| GW-38 | "Go/JVM named forms remove prefix hyphens and literal JVM/Cucumber tags may use underscores as documented aliases." | rule | ADR-037 | no | yes | |
| GW-39 | "A new or revised machine-proven criterion declares its proof type and needs classified positive and negative evidence through … unless it explicitly records `(single-direction)`." | rule | ADR-032 | no | **part (4 conditions)** | |
| GW-40 | "JVM metadata split across methods or inherited from a class receives no evidence credit." | rule | ADR-036 | no | yes | |
| GW-41 | "An unannotated legacy criterion keeps the one-supported-reference rule." | rule | ADR-032 | no | yes | |
| GW-42 | "A genuine `Test-type: Human` criterion is proven by its acceptance-brief line without fake code evidence, while `@draft` exempts only one not-yet-proven criterion inside an otherwise active file." | rule | ADR-033 | no | part (2 conditions) | |
| GW-43…GW-44 | `## Next` and its link | connective ×2 | — | — | — | |

**GW-37…GW-42 are the first of six statements of the criterion-evidence contract in these carriers.** The others are GM-11…GM-16, GL-27…GL-31, GG-38…GG-43, GA-08…GA-10, GO-11…GO-14, and CTB-16…CTB-20. Four of them — GW, GG, GA, GO — sit on P1, which means an adopter following the guide's own declared order reads the same contract four times, in four different compressions, none of them identical. That is the largest single duplication in this register and the one whose consolidation carries the most risk, because the four wordings differ and the differences may be meaningful. It is candidate 1 below.

---

## Register — `guide/design.md` (GD)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GD-01…GD-04 | `# The design of Cliewen`; what the page is for | connective ×4 | — | — | — | |
| GD-05…GD-14 | `## The problem: acceptance is the new bottleneck` and its argument | connective ×10 | — | — | — | |
| GD-15…GD-18 | `## Learned, not invented`; the first iteration and the book | connective ×4 | — | — | — | |
| GD-19…GD-27 | the nine observed failures of the first iteration | connective ×9 | — | — | — | |
| GD-28…GD-33 | the lesson, the two repairs, and the list of rules each repair produced | connective ×6 | — | — | — | |
| GD-34 | "keeping the durable documentation true is the workflow's job, performed inside every change — never a chore a human occasionally remembers to assign" | rule | ARCH-003 | no | part | |
| GD-35…GD-36 | `## The core: three elements and a red line`; the kernel image | connective ×2 | — | — | — | |
| GD-37 | **The verifiable thread** — the chain and "Every durable claim about the system traces to its declared proof" | rule | ARCH-003 | no | yes | |
| GD-38 | "'the agent decided' is never the answer" | rule | ARCH-003 | **GD-37** | yes | |
| GD-39 | **The human merge boundary** — "An agent never merges its own change. The merge is the act of acceptance, performed by a human, every time." | rule | C-012 | no | yes | |
| GD-40 | "Agents do the work; humans own the truth." | connective ×1 | — | — | — | |
| GD-41 | **The deterministic judge** — what `clue validate` checks and classifies | rule | ARCH-003 | no | part (6 checked shapes) | |
| GD-42 | "does not execute tests or inspect the pull request acceptance brief" | rule (boundary) | ADR-033 | no | yes | |
| GD-43 | "Remove any one element and the other two stop meaning anything" | connective ×1 | — | — | — | |
| GD-44 | "a change that alters what the thread connects, what a merge accepts, or what a green validate asserts always requires an explicit decision record and human acceptance" | rule | C-013 | no | yes | |
| GD-45 | "It never rides silently inside another change." | rule | C-013 | **GD-44** | yes | |
| GD-46 | "Everything outside the core … is periphery, changed at ordinary cost." | rule | PDR-013 | no | part | |
| GD-47 | `## The principles, and why each one` | connective ×1 | — | — | — | |
| GD-48 | "Machines enforce form; humans verify meaning." | rule | ARCH-003 | no | part | |
| GD-49 | the judge's checked evidence shapes, enumerated again | rule | ADR-032, ADR-033 | **GD-41** | **part (5 conditions)** | |
| GD-50 | "It does not check that a test establishes the right thing or that the brief contains the promised Human proof" | rule | ADR-033 | **GD-42** | yes | |
| GD-51…GD-53 | why the line is drawn explicitly | connective ×3 | — | — | — | |
| GD-54 | "The documentation is the spec, and the digest keeps it true." | rule | ARCH-003 | **GD-34** | part | |
| GD-55 | "durable truth lives in `/docs`, describing the system as it exists" | rule | ARCH-003 | **GD-54** | part | |
| GD-56 | "A full change works in a transient `/changes` workspace, and before merge that workspace is *digested* … the workspace deleted." | rule | PDR-002 | no | yes | |
| GD-57…GD-59 | two-systems-of-record rationale; what the digest does not do | connective ×3 | — | — | — | |
| GD-60 | "Decisions are routed by reversal cost." + log row, ADR, PDR | rule | PDR-006 | no | yes | |
| GD-61…GD-62 | proportional-ceremony rationale | connective ×2 | — | — | — | |
| GD-63 | "Agent decisions are born `inferred`; merge binds, approval signs." | rule | PDR-004 | no | yes | |
| GD-64 | "the CLI reports those decisions as their own signature backlog" | connective ×1 | — | — | — | |
| GD-65 | "Extracted non-decision meaning declares whether reversing it is cheap or expensive; … expensive inferred meaning cannot sit in an active capability's immediate graph slice." | rule | ADR-035 | no | yes | |
| GD-66 | "Reality gets one repository-local edge back." + the incident analysis marks the contradiction and links the failed claim | rule | ADR-035 | no | yes | |
| GD-67 | "It does not ingest telemetry or operate production; that larger feedback loop remains deliberately outside the current system." | rule (boundary) | ADR-035 | no | yes | |
| GD-68 | "Acceptance criteria are meaning-immutable." + retire the ID, mint a new one | rule | ADR-007 | no | yes | |
| GD-69 | "traceability that permits redefinition is not traceability" | connective ×1 | — | — | — | |
| GD-70 | "The work is challenged before a human sees it." + a reviewer in a fresh context | rule | PDR-012 | no | yes | |
| GD-71 | "Findings get fixed, and any fix invalidates the pass; the loop ends only on a clean review." | rule | PDR-012 — **contradicted by PDR-035 and C-017** | no | yes | |
| GD-72…GD-73 | loop-engineering rationale; what still needs the human | connective ×2 | — | — | — | |
| GD-74 | "**Small deltas, because Git merges text and not meaning.**" + the two-large-changes rationale | rule | **NONE — the same rule was withdrawn from `clue-delta` as DLT-33** | no | part | |
| GD-75 | "every change branches from accepted `main` and stays small enough that a human can actually hold its meaning at the merge gate" | rule | PDR-007 | no | part | |
| GD-76 | "each initiating author takes one Cliewen change to its pull request before starting the next" | rule | PDR-007 | no | yes | |
| GD-77 | "when a sibling change merges first, the candidate incorporates the new tip and re-runs its checks before proceeding" | rule | PDR-016 | no | yes | |
| GD-78 | "Before publication that may be a rebase; an existing PR merges current `main` into its branch so publication never needs a force push." | rule | PDR-016 | **GD-77** | yes | |
| GD-79…GD-80 | team parallelism; the constraint is a deliberate cost | connective ×2 | — | — | — | |
| GD-81 | "Process is tiered by depth into meaning." + the three tiers | rule | PDR-011, PDR-002 | no | **part (3 tiers in one statement)** | |
| GD-82 | "Two guards hold above the rules: an unclear tier escalates to the higher one, and the moment a decision or meaning change appears mid-work, the change moves to the full loop." | rule | PDR-002 | no | yes | |
| GD-83 | "keeping that price proportionate to ordinary work is treated as a live design obligation, not a settled matter" | rule | PDR-013 | no | part | |
| GD-84 | the cost paragraph — expect more tokens, more time, and more changes | connective ×1 | — | — | — | |
| GD-85…GD-86 | `## How adopters extend it`; "The core deliberately does not enumerate what a corpus may contain." | rule | ARCH-003 | no | yes | |
| GD-87 | "adopter-defined artifact types validate against the same form rules as everything else … without needing Cliewen's permission to exist" | rule | ADR-026 | no | yes | |
| GD-88 | "Your rules enter as constraints, each naming its source and whether a machine, an agent, or a human enforces it" | rule | ADR-045 | no | yes | |
| GD-89 | "Repository-local conventions extend the methodology in `AGENTS.md`; they may add to the rules but never override them, and a conflict stops the change for a human decision instead of being resolved silently." | rule | ADR-013 | no | yes | |
| GD-90 | "What you cannot do is redefine the core and still call it Cliewen" | rule | C-013 | **GD-44** | yes | |
| GD-91…GD-92 | `## What Cliewen does not solve`; why limits are stated | connective ×2 | — | — | — | |
| GD-93 | "**A green corpus can still describe the wrong product.**" + the thread stops at merge | connective ×1 | — | — | — | |
| GD-94 | "Production feedback is a deliberately closed door today: findings from operation re-enter the corpus as new goals or constraints through the ordinary loop, not through an automated pipeline." | rule (boundary) | ADR-035 | **GD-67** | yes | |
| GD-95…GD-96 | "**The judge checks form, not semantic alignment.**" + what it proves and cannot prove | rule | ADR-033 | **GD-41**, **GD-50** | part (3 conditions) | |
| GD-97 | "**Evidence harvesting has edges.**" + the supported conventions and the conservative JVM scanner | rule | ADR-036 | no | **part (7 conditions)** | |
| GD-98 | "A project rule whose verification is inherently human … lives as a constraint with `enforcement: human`" | rule | ADR-045 | **GD-88** | yes | |
| GD-99 | "a genuine Human-class acceptance criterion declares `Test-type: Human` and uses its acceptance-brief line as proof" | rule | ADR-033 | **GD-49** | yes | |
| GD-100 | "A criterion that is merely not proven yet uses `@draft` on its own tag line without drafting the entire capability." | rule | ADR-033 | **GD-49** | yes | |
| GD-101 | "**The human's half of the merge gate is written, but it is not automated.**" + the brief's five contents | rule | PDR-017 | no | **part (5 contents)** | |
| GD-102 | "its one-screen cap pushes oversized changes to split rather than hide meaning" | rule | PDR-017 | no | yes | |
| GD-103…GD-105 | "**It does not remove humans, and it does not manage projects.**"; ordinary Git discipline; where human attention goes | connective ×3 | — | — | — | |
| GD-106…GD-107 | `## Next` and its link | connective ×2 | — | — | — | |

**GD-71 and GD-74 are the finding this register was worth writing for.** Both are rules the method has already changed, still stated in the page a reviewer reads to evaluate the argument.

GD-71 says a fix invalidates the pass and the loop ends only on a clean review. [PDR-035](../decisions/PDR-035-bounded-agentic-review-loop.md) and [C-017](../constraints/C-017-agentic-review-loop-is-bounded.md), landed by CH-131 one change before this one, decided the opposite in two respects: an advisory repair does not invalidate a clean pass, and the loop has a bounded ordinary budget rather than an open-ended one. CH-131 updated `guide/change-loop.md`, `CONTRIBUTING.md`, both pull-request templates, and every skill carrier. It did not update `guide/design.md`, which states the same rule as a *principle*.

GD-74 states "small deltas" as a named principle of the design. The same rule was DLT-33 in AN-018, was asked what it was for, and was withdrawn as true advice that binds nothing — no reader can determine whether a delta was small enough. CH-132 removed it from `clue-delta`. It remains here, with a rationale paragraph, as one of nine principles.

Neither is an implementation defect in CH-131 or CH-132: both are [PDR-019](../decisions/PDR-019-methodology-contract-carriers-move-together.md) carrier-inventory misses, and PDR-019's own carrier list names "public or contributor guidance" explicitly. They are escalated together as **Q-01**, because the two cases want different answers — one rule is superseded and its statement must change, the other is withdrawn and its statement must go, and whether a *design rationale* for a withdrawn rule survives its rule is a question about what this page is for.

---

## Register — `guide/methodology.md` (GM)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GM-01…GM-03 | `# The verifiable thread`; the graph sentence; the Mermaid diagram | connective ×3 | — | — | — | |
| GM-04 | `## Goal` — "A goal states who wants an outcome and why." | rule | ARCH-003 | no | yes | |
| GM-05 | "Proposed goals form the inbox; accepting a goal says it is real, not that it must be built immediately." | rule | ADR-002 | no | yes | |
| GM-06 | `## Plan` — "A plan is a finite campaign serving a goal. Its milestones have explicit exit criteria and evidence." | rule | C-010 | no | yes | |
| GM-07 | "Completed plans are frozen rather than rewritten, so the plan index also records what the project has achieved." | rule | C-008 | no | yes | |
| GM-08 | `## Change` — "Cliewen does not own every repository edit." + the three tiers and first-match rule | rule | PDR-011, PDR-002 | no | **part (3 tiers)** | |
| GM-09 | "A full change uses a transient workspace under `/changes/CH-xxx-*` … deleted during the digest … A light change skips that workspace and its ready pull-request description becomes the proposal, but the branch and human merge boundary remain." | rule | PDR-002 | **GM-08** | part (3 conditions) | |
| GM-10 | "Two guards hold above the rules." + both guards | rule | PDR-002 | **GM-08** | yes | |
| GM-11 | "Product behavior stays full even when an existing criterion already states the behavior" | rule | PDR-018 | no | yes | |
| GM-12 | "`clue context <id>` resolves an artifact, criterion, or milestone identity and prints the transitive outgoing-link slice that governs it." | connective ×1 | — | — | — | |
| GM-13 | `## Capability and acceptance criterion` — "A capability owns three views: a plain-language explanation, Gherkin acceptance criteria, and implementer-facing design." | rule | ADR-025 | no | yes | |
| GM-14 | "Criterion IDs use the exact canonical grammar `<PREFIX>-<digits>[lowercase-suffix]` … only supported evidence-carrier aliases normalize their syntax." | rule | ADR-037 | no | yes | |
| GM-15 | "A new or revised machine-proven criterion declares `Test-type: …` and has supported Go, JVM, or Cucumber evidence classified by that type and positive/negative direction; JVM evidence carries all three parts on one supported executable." | rule | ADR-032, ADR-036 | no | **part (4 conditions)** | |
| GM-16 | "A genuinely single-direction scenario says so explicitly." | rule | ADR-032 | no | yes | |
| GM-17 | "`Test-type: Human` routes proof to the pull request acceptance brief and needs no code reference." | rule | ADR-033 | no | yes | |
| GM-18 | "A genuinely not-yet-proven criterion carries `@draft` on its own tag line without drafting its active siblings or capability, while an unannotated legacy criterion retains the one-supported-reference rule." | rule | ADR-033, ADR-032 | no | part (2 conditions) | |
| GM-19 | "`clue validate` checks these declarations and references but does not execute tests." | rule (boundary) | ADR-033 | no | yes | |
| GM-20 | "If a criterion's meaning changes, the old ID is retired as a tombstone and a new one is minted." | rule | ADR-007 | no | yes | |
| GM-21…GM-22 | "That immutability matters. A test tagged `AC-042` should always mean the same promise, even years later." | connective ×2 | — | — | — | |
| GM-23 | `## Constraints` — "Constraints are rules a Cliewen change must not break: a law, license, policy, project convention, or a verifiable quality bar" | rule | ADR-027 | no | yes | |
| GM-24 | "Each one names its source and whether a machine, agent, or human enforces it" | rule | ADR-045 | no | yes | |
| GM-25 | "and every Cliewen proposal is assessed against all of them" | rule | ADR-045 | **VFY-11** (P3 only) | **part (unbounded set — the same defect AN-018 recorded as VFY-11)** | |
| GM-26 | `## Four actors, one boundary` — "Skills carry process knowledge, `clue` is the deterministic judge, protected CI is the wall, and the human controls acceptance." | connective ×1 | — | — | — | |
| GM-27 | "A full-change PR begins with an acceptance brief that puts the remaining semantic questions … in front of the human." | rule | PDR-017 | no | part (3 questions) | |
| GM-28 | "the human does not have to repeat a locally completed code review, but the agent can never perform the merge that accepts its own work" | rule | C-012 | no | yes | |
| GM-29 | "CI becomes a wall only when its PR check is required and branch protection blocks integration without it." | rule | PDR-027 | no | yes | |
| GM-30 | "The wall enforces admission to merge; it is not acceptance evidence — that remains the criterion's classified executable reference or its Human-class acceptance-brief entry." | rule | PDR-027 | no | yes | |
| GM-31…GM-32 | `## Next` and its link | connective ×2 | — | — | — | |

**GM-25 is VFY-11 again, on a different surface.** AN-018 found `clue-verify`'s "assessed against every constraint" item uncheckable because the set is unbounded, and M-067 owns its repair. The guide states the same unbounded obligation to the human. Whatever repairs VFY-11 must reach GM-25, and M-067's scope names only the skill.

---

## Register — `guide/corpus.md` (GC)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GC-01 | `# The corpus` — "The `/docs` tree is Cliewen's permanent working memory." | rule | ARCH-003 | no | part | |
| GC-02 | "After classifying a task, agents read the focused slice from `clue context <id>` and follow more edges only when needed" | rule | PDR-034 | **HUB-22…24** (P3 only) | yes | |
| GC-03 | "people review the durable artifacts with the implementation, and Git records every accepted mutation" | connective ×1 | — | — | — | |
| GC-04…GC-05 | `## The taxonomy` and its seven-row table | connective ×2 | — | — | — | |
| GC-06 | "Each folder has a README that explains its type and contains a generated index of the artifacts beside it." | rule | C-016 | no | yes | |
| GC-07 | `## Identity is not location` — "Every artifact begins with YAML frontmatter" + the example block | rule | C-009 | no | yes | |
| GC-08…GC-12 | "The ID is the identity; the path is only its current address"; what `clue` scans; why refactoring is safe; what `clue context` prints and does not follow | connective ×5 | — | — | — | |
| GC-13 | `## One home per scope` — "System-wide and expensive-to-change design belongs under `architecture/`. Per-capability design lives beside the capability." | rule | ADR-025 | no | yes | |
| GC-14 | "Decisions explain durable choices but do not become substitute design documents." | rule | C-006 | no | yes | |
| GC-15 | "Findings record what an investigation observed but do not silently become accepted intent." | rule | PDR-030 | no | yes | |
| GC-16 | "The separation is intentionally strict: a fact with two homes will eventually disagree with itself." | connective ×1 | — | — | — | |
| GC-17 | `## Choose the right decision record` — "Start with the cost of reversing the decision, then ask what it changes" + the three-row routing table | rule | PDR-006 | no | yes | |
| GC-18…GC-19 | what the log is for; why ADRs and PDRs share a template | connective ×2 | — | — | — | |
| GC-20 | `## See a living corpus` — the three links into this repository's own corpus | connective ×1 | — | — | — | |
| GC-21…GC-22 | `## Next` and its link | connective ×2 | — | — | — | |

**GC-20 points at a completed campaign as the live example.** The link is to `P-007-core-hardening.md`, described as the "active campaign"; P-007 closed on 2026-07-28 and five campaigns have opened since. This is not a rule-bearing defect — the row is connective — but it is a stale claim in the guide's only pointer to a working corpus, and it belongs in the follow-on change rather than in an escalation, so it is recorded here and listed under *Editorial defects with no rule attached* below.

---

## Register — `guide/change-loop.md` (GL)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GL-01 | `# The change loop` | connective ×1 | — | — | — | |
| GL-02 | "The change loop applies when work belongs in Cliewen." | rule | PDR-011 | no | yes | |
| GL-03 | "Before loading the corpus, classify the request: if nothing about meaning changes … it is plain, and it uses an ordinary branch, relevant checks, a ready pull request, and human merge." | rule | PDR-011 | no | part (9-item surface list) | |
| GL-04 | "That plain route has no CH number, proposal metadata, corpus work, Cliewen verification, plan bookkeeping, or changelog entry." | rule | PDR-011 | **GL-03** | part (6 exemptions) | |
| GL-05 | "The light tier fits when meaning is only touched and no decision, acceptance meaning, plan semantics, or methodology carrier is affected; everything else is full." | rule | PDR-002 | no | yes | |
| GL-06 | "When the tier is unclear, take the higher one, and move to the full loop the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work." | rule | PDR-002 | no | yes | |
| GL-07 | "After classification, start from the smallest durable context that governs the task." | rule | PDR-034 | no | yes | |
| GL-08 | "`clue context <id>` prints the named artifact and the transitive closure of its outgoing links" | connective ×1 | — | — | — | |
| GL-09 | "If the request gives no usable ID, orient at `docs/README.md`, choose the closest artifact, and run the command from there." | rule | PDR-034 | no | yes | |
| GL-10 | "Shared goals have many reverse dependents, so `context` deliberately follows declared outgoing dependencies only" | connective ×1 | — | — | — | |
| GL-11…GL-15 | `## One real change, end to end`; the narrative of [`cliewen/cliewen` pull request 2](https://github.com/cliewen/cliewen/pull/2); the eight-row stage table; "That same shape applies to an ordinary product request" | connective ×5 | — | — | — | |
| GL-16 | "An existing criterion does not make a behavior change light." | rule | PDR-018 | **GM-11** | yes | |
| GL-17 | "The implementation changes executable evidence and may reveal that criterion, test boundary, and product reality disagree, so behavior remains a full reviewed delta." | rule | PDR-018 | **GL-16** | yes | |
| GL-18 | "The first adopter-history measurement found real workspace cost—144 transient lines for 76 durable corpus additions across its two full semantic changes—but no behavior-under-existing-criterion example" | connective ×1 | — | — | — | |
| GL-19 | `## 1. Branch` — "Create `ch-xxx-your-slug` from the current tip of `main`." | rule | C-012 | no | yes | |
| GL-20 | "One initiating author takes one initiated Cliewen change to its pull request before starting another, and a change never starts from unaccepted work." | rule | PDR-007 | no | yes | |
| GL-21 | "Plain changes, reviews, and help updating an existing pull request do not consume another initiated-change slot." | rule | PDR-011, PDR-016 | no | yes | |
| GL-22 | `## 2. Propose` — "A full change commits `/changes/CH-xxx-your-slug/proposal.md` before implementation." | rule | PDR-002 | **GM-09** | yes | |
| GL-23 | "The proposal says what will change, why it matters, which plan item it serves or that it is plan-less, and where the decision boundary lies." | rule | C-005 | no | part (4 contents) | |
| GL-24 | "they can opt into a spec-first pause. Record the pause in tasks and stop until they direct the work to continue; it is not the default loop." | rule | PDR-017, PDR-033 — **states the pause without PDR-033's report-and-ask** | no | yes | |
| GL-25 | "`tasks.md` is an ordered checklist with dependencies first." | rule | C-003 | no | yes | |
| GL-26 | "Tick a task the moment it completes." | rule | C-003 | no | yes | |
| GL-27 | "If a blocking decision appears, write it to `open-questions.md` and stop; the answer becomes a typed decision record rather than disappearing into chat." | rule | C-011 | no | yes | |
| GL-28 | `## 3. Implement` — "Change the permanent corpus and implementation together." | rule | ARCH-003 | **GD-34** | part | |
| GL-29 | "Behavior-changing work names the acceptance criteria it serves." | rule | G-001 | no | yes | |
| GL-30 | "Canonical criterion IDs use `<PREFIX>-<digits>[lowercase-suffix]` …; carrier aliases are limited to documented underscore tags and hyphen-free Go/JVM named prefixes." | rule | ADR-037 | **GM-14** | yes | |
| GL-31 | "A new or revised machine-proven criterion declares `Test-type: …` and gets supported Go, JVM, or Cucumber evidence classified by that type and positive/negative direction, unless it records `(single-direction)`; JVM evidence attaches all three parts to the same supported executable." | rule | ADR-032, ADR-036 | **GM-15**, **GM-16** | **part (4 conditions)** | |
| GL-32 | "A genuine `Test-type: Human` criterion is named in the acceptance brief as its proof; a not-yet-proven criterion carries `@draft` on its own tag line; an unannotated legacy criterion retains one supported reference." | rule | ADR-033, ADR-032 | **GM-17**, **GM-18** | part (3 conditions) | |
| GL-33 | "Every test declares one purpose: an AC ID, unit, sanity, or architecture." | rule | ADR-006 | no | yes | |
| GL-34 | "Never weaken a test or lint rule to make the build pass. A failing check is evidence about the change." | rule | C-004 | no | yes | |
| GL-35 | `## 4. Digest` — "update durable documentation, decisions, indexes, plan bookkeeping, and release notes for shipped behavior or workflow changes" | rule | C-003, C-002 | no | **part (5 conditions)** | |
| GL-36 | "Then delete the `/changes` workspace." | rule | PDR-002 | **GM-09** | yes | |
| GL-37 | "When a change completes a campaign's last milestone, the same digest sets that plan `completed`" | rule | log 2026-08-02 | no | yes | |
| GL-38 | "A successor plan is designated in that digest when one is decided; not having decided one is no reason to keep the finished plan open." | rule | log 2026-08-02 | **GL-37** | yes | |
| GL-39 | "Deletion is the digest … `main` never contains `/changes`." | rule | PDR-002 | **GL-36** | yes | |
| GL-40 | `## 5. Verify and review` — "Commit the complete candidate, run the repository tests and `clue validate --forbid-changes` against that commit, then run `clue-verify` on the same commit." | rule | PDR-012 | no | part (3 conditions) | |
| GL-41 | "a host with context-isolated delegation starts a fresh read-only reviewer with the declared intent but without the implementation conversation; other hosts disclose an in-context fallback" | rule | PDR-012 | no | yes | |
| GL-42 | "The reviewer returns correctness, regression, security, evidence, intent, or unjustified-complexity findings without editing." | rule | PDR-012 | no | yes | |
| GL-43 | "Every finding identifies the operative requirement or declared intent that is violated and its concrete consequence" | rule | PDR-016 | no | yes | |
| GL-44 | "authoritative decisions and explicit lifecycle rules govern before alternative readings become findings" | rule | PDR-012 | no | yes | |
| GL-45 | "Human-controlled merge does not imply duplicate human code review, and a release cut uses its versioned changelog section instead of `[Unreleased]`." | rule | C-012, ADR-012 | no | **part (two unrelated obligations in one sentence)** | |
| GL-46 | "The loop owns its classification regardless of the reviewer brief: a blocking finding is actionable and enters the hosted repair lifecycle; an advisory is a non-actionable observation for the publication gate" | rule | C-017 | no | yes | |
| GL-47 | "Counts and arithmetic disagreements are advisory, while a wrong, missing, or reused identity remains blocking, and the reviewer spends no pass re-deriving figures." | rule | C-017 | no | yes | |
| GL-48 | "The implementing context fixes blocking findings, commits the repaired candidate, reruns checks against that commit, and starts a new review pass on the same commit." | rule | PDR-012 | no | part (4 conditions) | |
| GL-49 | "An advisory repair may ride before a pass already required by a blocking repair; an advisory first reported by a pass with no blocking findings stays in the handoff for a later change" | rule | C-017 | no | part (2 conditions) | |
| GL-50 | "Three passes are the ordinary budget, and a fourth or later pass runs only after an immediately preceding pass with a blocking finding." | rule | C-017 | no | yes | |
| GL-51 | "The current commit needs a pass with no blocking findings before it is locally ready." | rule | PDR-012 | **GL-48** | yes | |
| GL-52 | "Fetch the latest `main`; if another change merges before the branch is first published, rebase and repeat review and verification." | rule | PDR-016 | no | yes | |
| GL-53 | "Once a PR exists, merge newer accepted `main` into its branch with a normal push instead of rewriting hosted history, then repeat both checks." | rule | PDR-016 | **GD-78** | yes | |
| GL-54 | "Separate authors keep separate branches from accepted `main`; collaboration is scoped to one pull request." | rule | PDR-016 | **GL-20** | yes | |
| GL-55 | "A review of an existing PR names the hosted head it inspected, and actionable findings live as unresolved hosted review conversations where the forge supports them." | rule | PDR-016 | no | yes | |
| GL-56 | "Any agent asked to fix one becomes the updater for that turn: it fetches and records that head, commits and reviews the repair, pushes without force, confirms the hosted head is the reviewed commit, and only then resolves the finding." | rule | PDR-016 | no | **part (5 conditions)** | |
| GL-57 | "If another updater moved the head, normal Git non-fast-forward protection forces reconciliation and a new verification pass." | rule | PDR-016 | no | yes | |
| GL-58 | "If findings cannot be published as resolvable conversations, the agent says the PR is not merge-ready and names the enforcement gap" | rule | PDR-016 | no | yes | |
| GL-59 | "no forge can detect an edit or intention that remained solely in a private worktree" | connective ×1 | — | — | — | |
| GL-60 | `## 6. Open the review gate` — "The pull request is an authorization and protected-integration gate, not a demand for duplicate human code review." | rule | C-012 | **GL-45** | yes | |
| GL-61 | "A solo developer may already have accepted the local candidate; the PR still prevents the agent that prepared it from accepting its own work." | rule | C-012 | **GL-60** | yes | |
| GL-62 | "The agent may publish the branch, but it never merges the pull request or pushes to `main`." | rule | C-012 | **GL-60**, **GL-61** | yes | **!** |
| GL-63 | "the human-controlled merge commit is the acceptance act; configure the forge to disable squash and rebase-and-merge so the original proposal, implementation, and digest commits remain reachable from `main`" | rule | PDR-021 | no | yes | |
| GL-64 | "A local rebase before first publication is allowed, but it is not the acceptance mode." | rule | PDR-021 | **GL-63** | yes | |
| GL-65 | "the PR starts with an acceptance brief. It asks whether the plan item is still wanted, puts the added or changed criteria and their scenarios in front of the human, and names what merge binds." | rule | PDR-017 | **GM-27** | part (3 contents) | |
| GL-66 | "The review loop adds an advisory verdict for each changed criterion — whether its referenced tests verify the scenario, something adjacent, or leave it undetermined." | rule | PDR-017 | no | yes | |
| GL-67 | "That is evidence for human judgment, not a semantic claim by `clue validate`: a green build and a fluent agent do not establish that the outcome is right." | rule | PDR-017 | no | yes | |
| GL-68 | "Enforcement requires the CI workflow to run on the PR, its result to be a required status check, and branch protection to block merge until that check passes." | rule | PDR-027 | **GM-29** | part (3 conditions) | |
| GL-69…GL-70 | "Local verification remains fast evidence; protected hosted CI is the safeguard"; the pointer to the CI wall guide | connective ×2 | — | — | — | |
| GL-71 | "Workflow and protection changes must never weaken the gate merely to make a change pass." | rule | C-004 | **GL-34** | yes | |
| GL-72 | "Open a ready pull request only after local review and verification pass." | rule | PDR-012 | **GL-51** | yes | |
| GL-73 | "Report the review mode, reviewed commit, number of passes run, and advisory findings left open, then confirm that the hosted head branch and SHA equal the clean, locally reviewed branch and `HEAD` before reporting it ready." | rule | C-017, PDR-016 | no | **part (5 conditions)** | |
| GL-74 | "A requested local branch or commit stopping point preserves work, but it is incomplete and not mergeable." | rule | C-012 | no | yes | |
| GL-75 | "Review fixes are committed, locally verified, and agent-reviewed again on the same branch." | rule | PDR-016 | **GL-48** | yes | |
| GL-76 | "Once the current commit has a clean pass and the worktree is clean, push it to the existing pull request and repeat the hosted-head check before reporting it ready again." | rule | PDR-016 | **GL-73** | part (3 conditions) | |
| GL-77 | "After the human merges, orient on the next unfinished plan milestone and ask before beginning it." | rule | **NONE — F-RB-09's rule; its constraint is owed by M-067** | no | yes | |
| GL-78…GL-79 | `## Next` and its link | connective ×2 | — | — | — | |

**GL-62's order defect is real but narrower than the skills'.** The merge prohibition appears in section 6 of six, after the whole loop. On P2 the reader met it at GD-39, early in `design.md`, and on P1 at GW-14; only a reader who opens `change-loop.md` directly meets the procedure first. That is the same shape AN-018 recorded for the skills, with the same mitigation — an earlier carrier states it — which makes it a weaker instance of one defect rather than a separate one.

**GL-24 is the one carrier edit M-064 names by hand.** [PDR-033](../decisions/PDR-033-planning-and-implementation-are-separate-steps.md) settled that at the pause the agent reports briefly and asks two questions — whether implementation begins, and whether the branch is pushed. `clue-delta` now carries that after CH-132; this page still describes the pause as a review of shape and says nothing about what happens at it. The rule is not missing from the corpus, only from this carrier, so it is a repair rather than an escalation.

**GL-77 is F-RB-09 restated to a human.** AN-018 found the orient-after-merge rule tracing to nothing and M-067 owes it a constraint. That constraint must reach this row too; M-067's scope names the shared fragment and `durable-work`, not the guide.

---

## Register — `guide/skills.md` (GS)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GS-01…GS-03 | `# The skills`; what skills are for; what `clue init` installs | connective ×3 | — | — | — | |
| GS-04 | "The routing file classifies plain work before loading the corpus; a plain change uses no Cliewen skill." | rule | PDR-011 | **GL-02**, **GL-03** | yes | |
| GS-05…GS-07 | how an agent loads a skill; why init also emits the bridges | connective ×3 | — | — | — | |
| GS-08 | "The bridge only points — every rule stays in the hub, where every assistant sees it, and the emitted file is yours to extend with genuinely Claude-specific instructions." | rule | PDR-022 | no | yes | |
| GS-09 | `## The lifecycle set` — the six-row skill table | rule (routing) | ARCH-002 | no | yes | |
| GS-10…GS-11 | `## Why the skills stay separate`; each skill owns a lifecycle boundary | connective ×2 | — | — | — | |
| GS-12 | "Verification does own the recurring challenge-and-repair hand-off: it delegates review into a clean context where supported, returns findings to the implementing context, and requires a clean pass on the resulting commit." | rule | PDR-012 | **GL-41**, **GL-51** | yes | |
| GS-13 | "The files are complete standalone artifacts, but repeated rules are generated from shared canonical sources." | rule | ADR-021 | no | yes | |
| GS-14 | "This keeps decision routing, change tiers, repository conventions, and the human review boundary identical across the set without creating runtime includes." | connective ×1 | — | — | — | |
| GS-15 | `## Version agreement` — "Distributed Cliewen skills carry an ownership marker and the same version as the released binary." | rule | ADR-022, ADR-011 | no | yes | |
| GS-16…GS-17 | what `clue validate` catches; why versioned guidance is reviewable | connective ×2 | — | — | — | |
| GS-18…GS-19 | `## Next` and its link | connective ×2 | — | — | — | |

**GS-09 lists all six skills.** AN-018's Q-07 found `clue-extract` missing from this repository's routing hub, and M-067 owns adding the row and its generator guard. The public guide has had the complete set all along, which is corroboration that the omission is a hub defect rather than a deliberate scope choice.

---

## Register — `guide/getting-started.md` (GG)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GG-01…GG-02 | `# Get started`; what the path does and that it touches no existing project | connective ×2 | — | — | — | |
| GG-03 | `## Prerequisites` — Git is required | rule | CAP-001 † | no | yes | |
| GG-04 | the install script's own requirements per platform | rule | CAP-001 † | no | yes | |
| GG-05 | permission to add one directory to the user `PATH`; no administrator rights | rule | CAP-001 † | no | yes | |
| GG-06 | the Go toolchain is optional | rule | CAP-001 † | no | yes | |
| GG-07 | an authenticated `gh` is recommended later for the pull-request loop | rule | CAP-001 † | no | yes | |
| GG-08 | "Node.js and npm are needed only to build this guide or contribute to Cliewen itself." | connective ×1 | — | — | — | |
| GG-09 | `## 1. Install clue` — the three install commands | rule | ADR-030 † | no | yes | |
| GG-10 | the pointer to the Claude Code plugin route | connective ×1 | — | — | — | |
| GG-11 | "Then open a new terminal and run `clue version` — a `PATH` change does not reach an already-running shell." | rule | ADR-030 † | no | yes | |
| GG-12…GG-18 | the printed version example; the profile note; what the script detects, verifies, and installs; `CLUE_INSTALL` and `CLUE_VERSION`; read-it-first; the Go route's `dev` stamp | connective ×7 | — | — | — | |
| GG-19 | "Upgrading later means re-running the same command. That moves the binary only — for a repository already using Cliewen it is half an upgrade; preview and apply the coordinated corpus and carrier migration with `clue migrate`" | rule | ADR-039 | no | yes | |
| GG-20 | "To find out whether there is anything to upgrade to, run `clue latest`." | rule | ADR-042 | no | yes | |
| GG-21 | what `clue latest` reports when you are behind | connective ×1 | — | — | — | |
| GG-22 | "It installs nothing, writes nothing in your repository, and exits 0 even when it cannot reach the network — so it is safe to run from a script or a coding agent's session start" | rule (boundary) | ADR-042 | no | yes | |
| GG-23 | `### Download a binary instead` — the manual route and the release-asset table | connective ×1 | — | — | — | |
| GG-24 | step 1 — "Verify the downloaded binary's SHA-256 matches its line in `SHA256SUMS`." | rule | ADR-030 † | no | yes | |
| GG-25 | the per-system checksum-command table | connective ×1 | — | — | — | |
| GG-26 | step 2 — rename to `clue.exe` or `clue`, and make it executable on macOS and Linux | rule | ADR-030 † | no | yes | |
| GG-27 | step 3 — move it into a directory on the user `PATH` | rule | ADR-030 † | no | yes | |
| GG-28 | step 4 — open a new terminal and run `clue version` | rule | ADR-030 † | **GG-11** | yes | |
| GG-29 | "The macOS binaries are unsigned and not notarized … First confirm the checksum matches, try `clue version` once, then open **System Settings → Privacy & Security**" | rule | ADR-030 † | no | part (3 steps) | |
| GG-30 | "The install script avoids this: a download made outside the browser carries no quarantine attribute." | connective ×1 | — | — | — | |
| GG-31 | `## 2. Initialize a disposable repository` — "Create an empty directory instead of experimenting in an existing project" and the five commands | rule | CAP-001 † | no | yes | |
| GG-32…GG-34 | the reported output; "The exact count can grow in a future release. The important result is the final `OK`."; the top-level tree | connective ×3 | — | — | — | |
| GG-35 | how the behind-notice reaches you without asking | connective ×1 | — | — | — | |
| GG-36 | "Never from `clue validate`, never when `CI` carries a value, and `CLUE_NO_UPDATE_NOTIFIER` turns the unasked notice off — standard output and the exit code never change." | rule (boundary) | ADR-042, PDR-023 | no | part (3 conditions) | |
| GG-37 | "Cliewen emits no configuration for any assistant to make that happen — the hub is the file they all read, and what your tools run is your business." | rule | PDR-023 | no | yes | |
| GG-38 | "`clue init` copies defaults but does not take ownership of your repository. You and your agent own the corpus prose and repository-specific instructions." | rule | ADR-013 | no | yes | |
| GG-39 | "`clue scaffold` and repeated `clue init` regenerate only the marked README index blocks; existing files are otherwise skipped, never replaced." | rule | ADR-019 | no | yes | |
| GG-40 | "The copied skills and workflow are versioned repository files, not background-managed services." | connective ×1 | — | — | — | |
| GG-41 | `## 3. See clue catch a broken thread` — add a goal and capability, keep the criterion `draft`, create these three files | rule | CAP-001 † | no | yes | |
| GG-42 | "Regenerate the two taxonomy indexes and validate" and its two commands | rule | CAP-001 † | no | yes | |
| GG-43…GG-44 | why the draft criterion is green; what flipping it to `active` prints | connective ×2 | — | — | — | |
| GG-45 | "That is the product's job: an active machine-proven promise cannot silently lose its acceptance evidence." | rule | ARCH-003 | no | yes | |
| GG-46 | "This deliberately small example is an unannotated legacy criterion, so one supported reference would satisfy it." | connective ×1 | — | — | — | |
| GG-47 | "To return the demo to green, set the whole criteria file back to `draft`." | rule | CAP-001 † | no | yes | |
| GG-48 | "In a real new or revised criterion, declare `Test-type: …` and add classified positive and negative evidence …; on the JVM, the AC identity, type, and direction belong to the same executable." | rule | ADR-032, ADR-036 | **GW-39**, **GW-40** | **part (4 conditions)** | |
| GG-49 | "Use `(single-direction)` only when one direction is honest." | rule | ADR-032 | **GW-39** | yes | |
| GG-50 | "If only that criterion is not ready, put `@draft` on its tag line instead of drafting proven siblings or the capability." | rule | ADR-033 | **GW-42** | yes | |
| GG-51 | "For a criterion whose proof is inherently human, declare `Test-type: Human`; naming it in the pull request acceptance brief is its proof, and no code test is invented." | rule | ADR-033 | **GW-42** | yes | |
| GG-52 | "`clue validate` recognizes the Human declaration and waives code evidence, but it cannot check that the brief supplies the proof — the pull request workflow and human merge gate do that." | rule (boundary) | ADR-033 | **GW-05** | yes | |
| GG-53 | "The judge also checks classified pairs, single-direction declarations, and `@draft`; it does not run tests, so your normal test runner remains responsible for whether executable evidence passes." | rule (boundary) | ADR-033 | **GG-52** | yes | |
| GG-54 | `## 4. Remove the experiment or continue` — leave the directory and delete it to undo the trial | rule | CAP-001 † | no | yes | |
| GG-55…GG-56 | removing the binary is not required; where to go next | connective ×2 | — | — | — | |
| GG-57 | "Before real work lands, configure the protected default branch for human-controlled merge commits only, choose and arm the thin CI caller, and run its disposable probe" | rule | PDR-021, PDR-027, ADR-038 | no | part (3 conditions) | |
| GG-58…GG-59 | `## Next` and its link | connective ×2 | — | — | — | |

**GG's prerequisite and install rows are the largest block of fresh traces in this register, and they trace to a capability rather than a decision.** The onboarding path is a *capability of the product* under [CAP-001](../capabilities/CAP-001-onboarding/README.md), and this page is the surface its criteria are about. That is a valid trace under PDR-029's four types and it is worth naming, because it is the one place in these carriers where the guide is not restating a rule stated elsewhere: a statement here that drifted from CAP-001's criteria would be a product defect, not a carrier defect. It also means these rows are checkable by a mechanism the others are not — `clue validate` will not judge the prose, but CAP-001's criteria have evidence.

---

## Register — `guide/adoption.md` (GA)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GA-01…GA-05 | `# Greenfield and brownfield`; the first-move difference; `## Who keeps the documentation current?`; tell the agent the outcome; you should not mirror changes by hand | connective ×5 | — | — | — | |
| GA-06 | "The agent reads `AGENTS.md`, loads the relevant Cliewen skill, and updates the implementation and durable corpus together on the same branch." | rule | ARCH-003 | **GL-28** | yes | |
| GA-07 | what `clue validate` judges | connective ×1 | — | — | — | |
| GA-08 | "A human still reviews whether the documentation and implementation say the right thing." | rule | C-012 | **GW-17** | yes | |
| GA-09…GA-10 | agent-maintained, not background synchronization; `clue` watches nothing and invents nothing | connective ×2 | — | — | — | |
| GA-11 | "The change loop requires local validation before a Cliewen pull request is ready for review." | rule | PDR-012 | no | yes | |
| GA-12 | "Once the generated CI caller is armed and its upstream validation job is a required check, broken traceability blocks merge." | rule | PDR-027 | **GW-18**, **GW-19** | yes | |
| GA-13 | "Plain changes keep the same required job but do not invoke the corpus validator." | rule | PDR-011 | no | yes | |
| GA-14 | `## Start with the minimum` — "Do not fill every corpus folder because the scaffold created it." | rule | ADR-025 | no | yes | |
| GA-15 | "A useful first thread needs four things" and the four-row table | rule | ARCH-003 | **GW-04** | yes | |
| GA-16 | "The criterion carries a canonical stable ID such as `AC-001`, `SNAP-SQS-001`, or `ADP-045b`" | rule | ADR-037 | **GW-37** | yes | |
| GA-17 | "A new or revised machine-proven criterion declares `Test-type: …`; focused positive and negative evidence both reference its ID and declared type through supported Go test names, per-executable JVM tags, or Cucumber scenario tags." | rule | ADR-032 | **GW-39**, **GG-48** | **part (4 conditions)** | |
| GA-18 | "A JVM executable carries all three evidence parts itself; class tags and unrelated methods cannot supply missing parts." | rule | ADR-036 | **GW-40**, **GG-48** | yes | |
| GA-19 | "A genuinely one-direction scenario says `(single-direction)`." | rule | ADR-032 | **GG-49** | yes | |
| GA-20 | "`Test-type: Human` instead uses the pull request acceptance brief as its proof and needs no code reference." | rule | ADR-033 | **GG-51** | yes | |
| GA-21 | "If one criterion is not ready, add `@draft` to that criterion's tag line while its proven siblings and active capability remain active." | rule | ADR-033 | **GG-50** | yes | |
| GA-22 | what `clue validate` classifies, counts, recognizes, and preserves | connective ×1 | — | — | — | |
| GA-23 | "It cannot check that the acceptance brief supplies Human proof; the pull request workflow and human merge gate do that." | rule (boundary) | ADR-033 | **GG-52** | yes | |
| GA-24 | "It validates executable evidence references but does not run the tests" | rule (boundary) | ADR-033 | **GG-53** | yes | |
| GA-25 | `## Add the wider corpus when it earns its keep` and its five-row table | rule | ADR-025 | no | yes | |
| GA-26 | "Leave unused categories empty. Cliewen is supposed to expose necessary reasoning, not reward document volume." | rule | ADR-025 | **GA-14** | yes | |
| GA-27 | `## When Cliewen is a poor fit` — "Do not adopt Cliewen when the repository cannot own both the intent and its acceptance evidence." | rule | **NONE** | no | yes | |
| GA-28 | the five poor-fit conditions | rule | **NONE** | no | part (5 conditions) | |
| GA-29 | "In those cases, use the project's existing lightweight notes and tests instead of creating a corpus nobody will maintain." | rule | **NONE** | **GA-27** | yes | |
| GA-30…GA-32 | `## Prompts that get useful work started`; you need not speak Cliewen's language; the greenfield prompt block | connective ×3 | — | — | — | |
| GA-33 | "The agent should establish the goal, make uncertainty visible, and propose the smallest verifiable plan before implementation." | rule | PDR-008 | no | part (3 conditions) | |
| GA-34…GA-35 | `### Make a routine change` and its prompt block | connective ×2 | — | — | — | |
| GA-36 | "The agent follows the change loop and leaves the merge decision to a human." | rule | C-012 | **GA-08** | yes | |
| GA-37…GA-38 | `### Adopt one existing repository` and its prompt block | connective ×2 | — | — | — | |
| GA-39 | "Use `clue-extract` once when the repository already contains specifications, decision notes, tagged tests, or other durable intent" | rule | ADR-008 | no | yes | |
| GA-40 | "Extraction is a meaning-level conversion, not a file copy." | connective ×1 | — | — | — | |
| GA-41 | "After its full change is proposed, the agent first writes a report-only rehearsal in that change's `/changes/` workspace." | rule | PDR-020 | no | yes | |
| GA-42 | "It inventories the source, proposed mappings, ID preservation or minting, uncertainty, test-purpose work, instruction conflicts, planned deletions, and plan doors without changing the target corpus, tests, routing, or hosted state." | rule | PDR-020 | no | **part (8 contents plus 4 protected surfaces)** | |
| GA-43 | "An unresolved conflict stops as an open question." | rule | PDR-024 | no | yes | |
| GA-44 | "Only explicit human direction starts the same change's mutation phase; that phase digests the rehearsal into the durable extraction report and eventually removes the old parallel specification corpus in the ready pull request." | rule | PDR-020 | no | part (3 conditions) | |
| GA-45 | "The report's criterion counts and mapping table are not typed: they live in one region rendered by `clue report` … and `clue validate` re-renders that region" | rule | ADR-054 | no | yes | |
| GA-46 | "That report is the readable summary, not a committed per-criterion registry: to inspect one criterion's mapping, follow the report's manifest reference and read the pinned manifest" | rule | PDR-028 | no | yes | |
| GA-47 | "Cliewen deliberately does not store a second document rendering every criterion, so this costs you some navigation and saves you a duplicate representation of the same mapping." | rule | PDR-028 | **GA-46** | yes | |
| GA-48 | "Every extracted artifact begins inferred: non-decision artifacts use `provenance: inferred` plus `reversal-cost: low|high` …, while decisions use `status: inferred` and `author: agent`" | rule | ADR-010, ADR-035 | no | part (2 conditions) | |
| GA-49 | "Human review promotes only the meaning it verifies; an active capability cannot depend on high-cost inferred meaning in its immediate graph slice, while low-cost findings may remain deferred." | rule | ADR-035 | **GW-33** | yes | |
| GA-50 | "Two extraction mappings ship today." + the OpenSpec mapping's clauses | rule | ADR-008, ADR-037 | no | part (5 conditions) | |
| GA-51 | the MADR mapping — the filename prefix survives as the ID, every record is born `inferred`, and source acceptance stays body prose rather than `accepted-by:` | rule | ADR-029 | no | part (3 conditions) | |
| GA-52 | "If the source format has no extraction mapping yet, writing that mapping is the first extraction task." | rule | ADR-008 | no | yes | |
| GA-53…GA-54 | `## When the system spans repositories, wikis, and tickets`; when a discovery pass helps; the prompt block | connective ×2 | — | — | — | |
| GA-55 | "Wiki pages and tickets can be evidence, preferably through revision-pinned links or stable exports. They do not become a second system of record, and Cliewen does not live-sync them after adoption." | rule | ADR-040 | no | yes | |
| GA-56 | the four-bullet repository boundary — one repository per extraction, evidence discovered only inside the validated repository, separate corpora, and no cross-repository claim | rule | ADR-044 | no | part (4 conditions) | |
| GA-57…GA-58 | `## Next` and its link | connective ×2 | — | — | — | |

**GA-27…GA-29 are the only rule-bearing statements in this register that trace to nothing and are not already owed a carrier by M-067.** The five poor-fit conditions are the method's own statement of where it does not apply, and `operations.md` links to them as its fallback. Nothing in the corpus states them. Escalated as **Q-02**.

---

## Register — `guide/ci-wall.md` (GB)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GB-01…GB-02 | `# Make CI enforce Cliewen`; what `clue init` gives you | connective ×2 | — | — | — | |
| GB-03 | "The caller starts unarmed when it uses the default vendored source, so the job warns and skips corpus validation until the pinned Linux release binary and its checksum file are committed under `.github/tools/`." | rule | ADR-038 | no | yes | |
| GB-04 | "There are three separate jobs here" and the three-item list | connective ×1 | — | — | — | |
| GB-05 | "CI without branch protection is a dashboard. Branch protection without the validator cannot see a broken Cliewen thread. You need both." | rule | PDR-027 | **GW-19** | yes | |
| GB-06 | `## 1. Choose the caller inputs` — one `uses:` reference at an immutable reference, plus four inputs | rule | ADR-038 | no | part (2 conditions) | |
| GB-07 | "Use the exact generated version; do not substitute `latest`." | rule | ADR-038 | no | yes | |
| GB-08 | what the default caller inputs are and what they stage | connective ×1 | — | — | — | |
| GB-09 | "A repository that needs a self-hosted/no-root runner changes only the caller's runner-label JSON and writable install directory." | rule | ADR-038 | no | yes | |
| GB-10 | "A repository that downloads the release instead of committing `.github/tools/` changes only `clue-source: release`" | rule | ADR-038 | no | yes | |
| GB-11 | "Prefer an install directory outside the checkout." | rule | ADR-038 | no | yes | |
| GB-12 | "If your policy requires a path inside the workspace, add it to `.gitignore`." | rule | ADR-038 | no | yes | |
| GB-13 | why same-release verification catches corruption rather than establishing trust | connective ×1 | — | — | — | |
| GB-14 | "The reusable workflow owns its action references and all validation steps. Do not copy its checkout, scope, warning, acceptance-brief, or `clue validate` steps into the caller." | rule | ADR-038 | no | yes | |
| GB-15 | "Updating the one upstream reference is the reviewed upgrade that imports those fixes while retaining the caller's local choices." | rule | ADR-038 | **GB-14** | yes | |
| GB-16 | what reference a release binary and a source build emit | connective ×1 | — | — | — | |
| GB-17 | "Both forms are immutable; protect release tags from force updates, and replace a tag with the exact source SHA when your hosting policy requires SHA-only references." | rule | ADR-038 | no | part (2 conditions) | |
| GB-18 | `## 2. Arm the pinned judge` — "The examples below use `0.7.0`; replace it with the `clue-version` in your caller." | rule | ADR-038 | no | yes | |
| GB-19 | "Create the tools directory, then download the Linux amd64 binary and the release checksum file" and its command group | rule | ADR-038 | no | yes | |
| GB-20 | "The runner is Linux amd64 even when you develop on Windows or macOS." | connective ×1 | — | — | — | |
| GB-21 | "Verify the vendored file before committing it" and the per-system check table | rule | ADR-030 | **GG-24** | yes | |
| GB-22 | "Commit both files with the generated caller" and its commands | rule | ADR-038 | no | yes | |
| GB-23 | "Do not edit `clue-version` without replacing both vendored files from the matching release." | rule | ADR-038 | no | yes | |
| GB-24 | the upstream workflow re-verifies on every run; the staging directory defaults to `RUNNER_TEMP` | connective ×1 | — | — | — | |
| GB-25 | `## 3. Know what armed means` and its four-state table | connective ×1 | — | — | — | |
| GB-26 | "The `--forbid-changes` flag is the digest boundary. A pull request with a transient `/changes/CH-xxx-*` workspace is unfinished, even if ordinary validation passes." | rule | PDR-002 | **GW-15** | yes | |
| GB-27 | "The hosted check turns red until the change is digested into the permanent corpus and the workspace is removed." | connective ×1 | — | — | — | |
| GB-28 | `## 4. Require the check on GitHub` — "Push the caller and let its pull request run once." | rule | PDR-027 | no | yes | |
| GB-29 | "The caller and reusable job are both named `validate`; select the exact check GitHub displays" | rule | ADR-038 | no | yes | |
| GB-30 | where the ruleset UI lives and GitHub's own reference documentation | connective ×1 | — | — | — | |
| GB-31 | "Two separate surfaces control which merge methods are available, and Cliewen needs both pointed the same way." | rule | PDR-021 | no | yes | |
| GB-32 | "allow **Create a merge commit**, and disable **Squash and merge** and **Rebase and merge**" | rule | PDR-021 | **GW-29** | yes | |
| GB-33 | "Configure one active ruleset" and its twelve-row setting table | rule | PDR-027 | no | **part (12 settings as one obligation)** | |
| GB-34 | "An empty bypass list matters. A rule that the normal maintainer or automation can silently bypass is not the merge boundary Cliewen assumes." | rule | PDR-027 | **GB-33** | yes | |
| GB-35…GB-36 | how GitHub combines overlapping rulesets; which plans offer rulesets | connective ×2 | — | — | — | |
| GB-37 | "If the Rulesets menu is unavailable but classic branch protection is offered, configure the same default-branch requirements there" | rule | PDR-027 | **GB-33** | part (5 conditions) | |
| GB-38 | "Classic branch protection carries no merge-method control of its own, so the repository-wide options are the only lever — enable merge commits and disable squash and rebase-and-merge there." | rule | PDR-021 | **GB-32** | yes | |
| GB-39 | "If the hosting plan offers neither enforcement surface, the workflow can report failures and agents can warn about unresolved findings, but neither can block integration." | rule | PDR-027 | no | yes | |
| GB-40 | "After saving, inspect the effective default-branch rules" and what you should see | rule | PDR-027 | no | part (5 expectations) | |
| GB-41 | "Align them too, so nobody is offered a button that the default branch will reject" and the `gh api` command | rule | PDR-021 | **GB-31** | yes | |
| GB-42 | "Expect `merge: true`, `squash: false`, and `rebase: false`." | rule | PDR-021 | **GB-32** | yes | |
| GB-43 | "Do not remove an existing stronger requirement merely to match this minimum." | rule | C-004 | no | yes | |
| GB-44 | "If the default branch still permits squash or rebase-and-merge, the repository is not ready for a full Cliewen change." | rule | PDR-021 | **GB-32** | yes | |
| GB-45 | `## 5. Prove failure blocks merge` — "Do this once in a disposable branch." | rule | PDR-027 | no | yes | |
| GB-46 | what the probe creates and why normal validation stays green | connective ×1 | — | — | — | |
| GB-47 | "Run it only after the merge-method checks above pass; the final pull request must be accepted with a merge commit." | rule | PDR-021 | **GB-32** | part (the referent of "the final pull request" is not stated, and the probe is closed unmerged) | |
| GB-48…GB-51 | the probe's four steps — branch, the `CH-999` proposal, the two local checks and their expected outputs, commit/push/create/watch | rule ×4 | PDR-027 | no | yes | |
| GB-52 | "The `validate` check must fail and GitHub must show the pull request as blocked. If the check is red but the merge button is available, the workflow works but the ruleset does not enforce it yet." | rule | PDR-027 | no | yes | |
| GB-53 | "Close the probe without merging it" and its commands | rule | PDR-027 | no | yes | |
| GB-54 | `## Other forges` — "Other forges are outside this workflow contract unless they can provide an equivalent reusable-workflow boundary and the same immutable reference, checksum, stable-check, and protected-merge guarantees." | rule | ADR-038 | no | part (5 guarantees) | |
| GB-55 | "copy the contract rather than GitHub's labels" and its seven bullets | rule | PDR-027, ADR-038 | **GB-33** | **part (7 conditions)** | |
| GB-56 | "If your forge or hosting plan cannot enforce those conditions, local Cliewen validation still catches mistakes, but CI is evidence rather than a wall." | rule | PDR-027 | **GB-39** | yes | |
| GB-57…GB-58 | `## Next` and its link | connective ×2 | — | — | — | |

---

## Register — `guide/operations.md` (GO)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GO-01…GO-03 | `# Operate Cliewen safely`; who the page is for; what it describes | connective ×3 | — | — | — | |
| GO-04 | `## What Cliewen ships and checks` — the `clue` row: versioned release binaries per platform and a Go source route | connective ×1 | — | — | — | |
| GO-05 | the Test evidence row — the three supported conventions and their exact forms | rule | ADR-005, ADR-036 | no | part (3 conventions) | |
| GO-06 | the Agent guidance row — six generated skills in `.agents/skills/`, mirrored for Claude Code | rule | ADR-021, PDR-022 | no | yes | |
| GO-07 | the Session-start discovery row — which commands print the notice, the four never-cases, and that no vendor configuration is emitted | rule | ADR-042, PDR-023 | **GG-36**, **GG-37** | **part (6 conditions)** | |
| GO-08 | the GitHub CI row — the thin caller, the runner and binary-source choice, and requiring the exact `validate` check | rule | ADR-038 | **GB-06**, **GB-29** | yes | |
| GO-09 | the Full-change merge row — human-controlled merge commits only | rule | PDR-021 | **GB-32** | yes | |
| GO-10 | the Validation row — what `clue validate` checks, what `--forbid-changes` adds, and that an outward reference must name what it points at | rule | ADR-040 | no | part (4 conditions) | |
| GO-11 | the Focused context row | connective ×1 | — | — | — | |
| GO-12 | the Release check row — reports and installs nothing, never a required check, calm when it cannot tell | rule | ADR-042 | **GG-22** | part (3 conditions) | |
| GO-13 | the External addresses row — `clue refs` classifications, only *gone* is an error, `--apply` leaves completed plans alone, keep it out of branch protection | rule | ADR-040 | no | **part (4 conditions)** | |
| GO-14 | "Cliewen does not run your tests, does not synchronize tickets or wikis, does not update installed files in the background, and does not validate evidence across repositories." | rule (boundary) | ADR-044 | **GA-56** | yes | |
| GO-15 | "On the JVM, all three parts belong to the same executable; class tags and metadata split across methods do not count, and ambiguous or unsupported source syntax is diagnosed instead of guessed." | rule | ADR-036 | **GA-18**, **GO-05** | yes | |
| GO-16 | "A different framework needs the stable JVM named-executable form or another supported profile before its references can satisfy `clue validate`; do not treat an arbitrary comment or tag as equivalent evidence." | rule | ADR-005 | no | yes | |
| GO-17 | "A new or revised machine-proven criterion declares `Test-type: …`, and the validator requires supported evidence classified with that type in positive and negative directions; `(single-direction)` is the explicit narrow exception." | rule | ADR-032 | **GA-17**, **GA-19** | part (3 conditions) | |
| GO-18 | "An unannotated legacy criterion retains its one-supported-reference rule." | rule | ADR-032 | **GW-41** | yes | |
| GO-19 | "`Test-type: Human` uses the pull request acceptance brief rather than code evidence, and `@draft` exempts only the individual criterion that is not yet proven." | rule | ADR-033 | **GA-20**, **GA-21** | yes | |
| GO-20 | `## Preserve the full-change archive` — allow the merge-commit mode and disable squash and rebase | rule | PDR-021 | **GB-32** | yes | |
| GO-21 | "A local rebase before first publication is allowed, but the human acceptance mode remains a merge commit." | rule | PDR-021 | no | yes | |
| GO-22 | "Because the merge shape is set per branch or per repository rather than per pull request, a default branch restricted this way restricts plain changes into it too; route work that needs another merge shape to a branch the rule does not target." | rule | PDR-021 | no | yes | |
| GO-23 | "Before adoption, run the CI wall's branch-protection probe and verify both the protected branch's allowed merge methods and the repository-wide pull-request settings." | rule | PDR-027 | **GB-45**, **GB-41** | yes | |
| GO-24 | "If the forge cannot enforce this boundary, it is outside Cliewen's supported full-change adoption path rather than an equivalent configuration." | rule | PDR-021 | **GB-54** | yes | |
| GO-25 | `## Upgrade one coordinated set` — "Ask first whether there is anything to upgrade to." | rule | ADR-042 | **GG-20** | yes | |
| GO-26…GO-27 | what `clue version` can and cannot tell you; what `clue latest` prints when you are behind | connective ×2 | — | — | — | |
| GO-28 | "It never installs anything and never writes in your repository: the upgrade stays a reviewed change." | rule | ADR-042 | **GG-22** | yes | |
| GO-29…GO-31 | offline behaviour and `--quiet`; the day-long cache and that asking yourself always reaches the list; that neither you nor your agent must remember to ask | connective ×3 | — | — | — | |
| GO-32 | "it is bounded the same way: one line, never on standard output, never changing an exit code, and cached for a day" | rule | ADR-042 | **GO-07** | part (4 conditions) | |
| GO-33 | the hour-long memory of a non-answer | connective ×1 | — | — | — | |
| GO-34 | "It is deliberately narrow about where it appears. Never from `clue validate` …, never from `clue version` …, never when `CI` holds a value …, and never when `CLUE_NO_UPDATE_NOTIFIER` is present at all" | rule | ADR-042 | **GO-07**, **GG-36** | part (4 conditions) | |
| GO-35 | "standard output and the exit code are byte-identical with the notice and without it, on every command" | rule | ADR-042 | **GO-32** | yes | |
| GO-36 | "If you capture standard error too and want it clean, set `CLUE_NO_UPDATE_NOTIFIER`." | connective ×1 | — | — | — | |
| GO-37 | "The `AGENTS.md` that `clue init` materializes also tells the agent to run the quiet check before its first tool call and to route a non-empty answer to `clue-upgrade`." | rule | PDR-023 | **HUB-03**, **HUB-06** | yes | |
| GO-38 | "Cliewen ships no hook and no settings file for any assistant: the method belongs to no tool, and configuring yours is your choice, in your file." | rule | PDR-023 | **GG-37** | yes | |
| GO-39 | "If you would rather have it run deterministically at session start, wire `clue latest --quiet` into your assistant's own mechanism" | connective ×1 | — | — | — | |
| GO-40 | "If you see `unknown command \"latest\"`, that is the answer — that binary is behind." | rule | PDR-023 | **HUB-08** | yes | |
| GO-41 | "`clue migrate` reports a hub that never asks — run `clue init` if you have no hub at all, or add the line yourself to the one you wrote. It repairs neither, because your hub is your prose." | rule | ADR-013, PDR-023 | no | yes | |
| GO-42 | "Keep the binary, generated skills, and CI caller on the same release when you upgrade." | rule | ADR-038 | no | yes | |
| GO-43 | "First make the current repository green and branch through its normal review process." | rule | C-012 | no | yes | |
| GO-44 | "Then choose the release …, verify the new platform binary against that release's `SHA256SUMS`, and confirm `clue version` prints the chosen version." | rule | ADR-030 | **GG-24** | part (3 conditions) | |
| GO-45 | "`clue init` is deliberately non-destructive: it skips existing files, so it is not an updater." | rule | ADR-019 | **GG-39** | yes | |
| GO-46 | "Preview the coordinated upgrade with `clue migrate`; when the plan reports only deterministic work, apply it with `clue migrate --apply --reversal-cost=low|high`" | rule | ADR-039 | **GG-19** | yes | |
| GO-47 | what migration updates, what it adds only under a supported preceding release, and what it leaves unchanged | connective ×2 | — | — | — | |
| GO-48 | "A missing semantic choice, unsupported old release, copied workflow, or locally edited generated file is reported without partial writes; resolve it in the repository's reviewed change and rerun the preview." | rule | ADR-039, ADR-052 | no | part (4 conditions) | |
| GO-49 | "The preview also lists every external reference that names no repository … and repairs none of them … Those are yours to resolve, and `clue validate` fails until you do." | rule | ADR-040 | **GO-10** | yes | |
| GO-50 | "Keep the existing required `validate` check in place throughout the upgrade … Make it required only when you arm the wall for the first time." | rule | PDR-027 | no | part (2 conditions) | |
| GO-51 | "If a released binary reports skill drift, do not edit a version number to silence it." | rule | C-004, ADR-011 | no | yes | |
| GO-52 | "The message names both ways out: move forward … or stay where you are … Either way, run `clue validate` afterwards." | rule | ADR-011 | no | part (2 routes) | |
| GO-53 | "A checkout build reports `dev` and cannot detect binary-to-skill release drift, so use a released binary for this check." | rule | ADR-011 | no | yes | |
| GO-54 | "Re-running the install script moves the binary and nothing else." + why that produces the drift report | rule | ADR-011 | **GG-19**, **GO-51** | yes | |
| GO-55 | "This is the check working, not a broken upgrade. Resolve it by completing the coordinated set … or by pinning the release the repository still carries" | rule | ADR-011 | **GO-52** | yes | |
| GO-56 | `## Recover without bypassing the evidence` and its six-row recovery table | rule | C-004, ADR-034 | no | **part (6 situations as one obligation)** | |
| GO-57 | `## Evidence from other repositories` — "Cliewen's hyperfine and es-toolkit work were read-and-apply foreign-soil trials, not adoptions." | rule | PDR-005 | no | yes | |
| GO-58 | "They are useful evidence about methodology boundaries, not proof that those projects use or endorse Cliewen." | rule | PDR-005 | **GO-57** | yes | |
| GO-59 | "When your repository's ownership, test evidence, or merge boundary cannot meet these conditions, keep the existing lightweight notes and tests instead of forcing an adoption." | rule | **NONE — GA-27's rule** | **GA-29** | yes | |
| GO-60…GO-61 | `## Next` and its link | connective ×2 | — | — | — | |

**GO-07 and GO-34 name five notice-printing commands; the shipped set is seven.** `notifierCommands` in `cmd/clue/main.go` holds `init`, `scaffold`, `context`, `migrate`, `id`, `refs`, and `report`; the CLI's own usage text lists the same seven. This page and `AGENTS.md` (HUB-09) both name only `context`, `migrate`, `refs`, `init`, and `scaffold`. No judgement is required — the code and the usage text agree, and the two prose carriers understate them — so this is a repair rather than an escalation. It is recorded here because the repair crosses milestones: `AGENTS.md` belongs to AN-018's surface and M-063 has closed.

---

## Register — `guide/plugin.md` (GP)

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| GP-01…GP-02 | `# Install from Claude Code`; the two plugin commands | connective ×2 | — | — | — | |
| GP-03 | "The plugin adds one skill, `/cliewen:setup`. It … confirms the binary reports a release version, and then **asks** before running `clue init`." | rule | ADR-031 | no | yes | |
| GP-04 | "Nothing is written into your repository until you say so." | rule | ADR-031 | **GP-03** | yes | |
| GP-05…GP-06 | "That is the entire plugin."; what the rest of the page is about | connective ×2 | — | — | — | |
| GP-07 | `## What the plugin does not install` — "It does not ship … the six managed skills that run the Cliewen loop. Those arrive when you run `clue init` in a repository, and only then." | rule | ADR-031 | no | yes | |
| GP-08 | "Cliewen's skills are committed repository files … each stamped with the version of the binary that wrote it, and `clue validate` fails if the binary and those skills ever disagree." | rule | ADR-011 | no | yes | |
| GP-09…GP-11 | why a green validate then means something; how a plugin's components live outside the drift check; what bundling the six would produce | connective ×3 | — | — | — | |
| GP-12 | "`clue init` is the only supported way to put Cliewen skills in a repository" | rule | ADR-031 | **GP-07** | yes | |
| GP-13 | "The same reasoning applies to copying a Cliewen skill into your personal `~/.claude/skills/` directory by hand. It will appear to work, and it will be invisible to the check that exists to catch it going stale." | rule | ADR-031 | **GP-12** | yes | |
| GP-14…GP-15 | `## Which route to use`; "neither is more official than the other" | connective ×2 | — | — | — | |
| GP-16 | the three-row route table, including "Setting up CI … the pinned, checksum-verified binary the generated workflow installs — never a plugin" | rule | ADR-038 | no | yes | |
| GP-17 | "Upgrading `clue` moves the binary only … preview and apply the coordinated repository migration with `clue migrate`, and use the drift report as the check that the pair still needs work." | rule | ADR-039 | **GG-19** | yes | |
| GP-18…GP-19 | `## Next` and its link | connective ×2 | — | — | — | |

---

## Register — `CONTRIBUTING.md` (CTB)

This carrier is on P3, so a statement restating the routing hub or a skill **is** scored as duplicated and the duplicate's AN-018 ID is named.

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| CTB-01 | `# Contributing to Cliewen` and its thanks | connective ×1 | — | — | — | |
| CTB-02 | "Participation is governed by the Code of Conduct." | rule | PDR-010 | no | yes | |
| CTB-03 | `## Choose the Right Route` — "Suspected security vulnerability: follow the private reporting route in SECURITY.md. Never disclose it in a public issue or pull request." | rule | PDR-010 | no | yes | |
| CTB-04 | "Private conduct concern: use the private conduct-reporting address … Never open a public conduct issue." | rule | PDR-010 | no | yes | |
| CTB-05 | "Reproducible defect: open the structured bug form." | rule | PDR-010 | no | yes | |
| CTB-06 | "Desired outcome or unmet need: open the proposed-goal form. A goal issue records demand for consideration; it does not add the goal to Cliewen's accepted plan." | rule | ADR-002, PDR-010 | no | yes | |
| CTB-07 | "Small editorial correction with no effect on behavior, intent, evidence, decisions, plans, policy, or methodology: use the plain-change route below." | rule | PDR-011 | **HUB-14** | yes | |
| CTB-08 | `## Before Starting a Change` — "Classify the work before loading the corpus. Three rules set the tier — plain, then light, then full — and you take the first that matches" | rule | PDR-011, PDR-002 | **HUB-12**, **HUB-13**, **F-CT-01**, **F-CT-02** | yes | |
| CTB-09 | "A change is plain when nothing about meaning changes: it has no effect on product behavior, intent, acceptance evidence, decisions, plans, policy, or methodology." | rule | PDR-011 | **HUB-14**, **F-CT-03** | part (7 surfaces) | |
| CTB-10 | "Acceptance evidence includes executable evidence and genuine `Test-type: Human` proof carried by the pull-request acceptance brief." | rule | ADR-033 | no | yes | |
| CTB-11 | "Protected surfaces are never plain: `/docs`, `/changes`, code, tests, configuration, build and release machinery, security and governance policy, `AGENTS.md`, skills, and lint rules." | rule | PDR-011 | **HUB-14**, **F-CT-03** | part (10-item list) | |
| CTB-12 | "Changes to commands, contracts, user workflow, or normative instructions are not editorial." | rule | PDR-011 | **HUB-14** | yes | |
| CTB-13 | "Two guards hold above the rules, from this first classification onward." and both guards | rule | PDR-002 | **HUB-20**, **HUB-21**, **F-CT-06** | yes | |
| CTB-14 | "A plain change starts from the current tip of `main`, uses an ordinary branch, runs checks relevant to the changed surface, and opens a ready pull request." | rule | PDR-011 | **HUB-15**, **F-CT-03** | part (4 conditions) | |
| CTB-15 | "It needs no CH number, plan declaration, proposal, corpus read, Cliewen skill, full verification checklist, plan bookkeeping, or changelog entry." | rule | PDR-011 | **HUB-15**, **GL-04** | part (8 exemptions) | |
| CTB-16 | "For every other change, search existing issues, pull requests, and the system-of-record under `docs/`." | rule | PDR-034 | no | yes | |
| CTB-17 | "Every Cliewen change serves an existing item under `docs/plans/` or explicitly declares itself plan-less." | rule | C-005 | **HUB-44**, **DLT-08** | yes | |
| CTB-18 | "A contributor may initiate one Cliewen change at a time; plain changes, reviewing, and helping update an existing pull request do not consume another initiated-change slot." | rule | PDR-007 | **HUB-32**, **HUB-35**, **F-RB-01** | yes | |
| CTB-19 | "Every branch starts from accepted `main` and never from unmerged work." | rule | C-012 | **HUB-31**, **F-RB-01** | yes | |
| CTB-20 | "load the smallest durable slice that can govern the work: run `clue context <id>` when the request names or resolves to an artifact; otherwise orient at `docs/README.md`, select the closest plan, capability, criterion, or decision, and run `clue context` from there." | rule | PDR-034 | **HUB-22**, **HUB-23** | part (2 branches) | |
| CTB-21 | "Follow additional edges only when the task discovers them." | rule | PDR-034 | **HUB-24** | part | |
| CTB-22 | "Use the next free `CH-xxx` identifier visible in git history and any active `/changes/` workspace" | rule | ADR-009 — **superseded by ADR-048, whose ledger and `clue id next` have shipped** | **DLT-06 (removed by M-063)** | yes | |
| CTB-23 | `## Choose the Change Tier` — "A change is light when meaning is touched but not changed: it makes no decision, changes no acceptance-criterion or capability meaning, makes no semantic plan mutation, and touches no methodology carrier" | rule | PDR-002 | **HUB-17**, **F-CT-04** | part (4 conditions) | |
| CTB-24 | "A light change has no `/changes/` workspace; its pull-request description is the proposal and states what, why, and the plan item or plan-less declaration." | rule | PDR-002 | **HUB-30**, **F-CT-04** | yes | |
| CTB-25 | "Every other change is full. Before implementation, add `/changes/<CH-xxx-slug>/proposal.md`, `tasks.md`, and `open-questions.md`, then commit that proposal by itself." | rule | PDR-002 | **HUB-28**, **DLT-07** | yes | |
| CTB-26 | "Record unresolved decisions in `open-questions.md` and stop until a human answer can be captured as a typed decision." | rule | C-011 | **HUB-46**, **DLT-14** | yes | |
| CTB-27 | "A product behavior change remains full when an existing criterion already describes the intended behavior." | rule | PDR-018 | **HUB-19**, **F-CT-05** | yes | |
| CTB-28 | `## Implement and Digest` — "Keep the change focused on its proposal and tick each task immediately when it is complete." | rule | C-003 | **DLT-10** | yes | |
| CTB-29 | "Update permanent capability, acceptance-criteria, decision, constraint, architecture, and plan artifacts when their meaning changes." | rule | ADR-025 | **DLT-16** | yes | |
| CTB-30 | the criterion-evidence contract — declared type, classified pair, `(single-direction)`, the JVM triple on one executable, Human proof in the brief, `@draft`, and the legacy one-reference rule | rule | ADR-032, ADR-033, ADR-036 | **VFY-09**, **DLT-17…DLT-22** | **part (7 conditions)** | |
| CTB-31 | "`clue validate` validates declarations and references but does not execute tests." | rule (boundary) | ADR-033 | no | yes | |
| CTB-32 | "Never weaken a test, lint rule, or quality gate to make a build pass." | rule | C-004 | **HUB-48**, **VFY-01d** | yes | |
| CTB-33 | "If a Cliewen-owned skill changes, edit `internal/skills/source/` and run `go generate ./internal/skills`; do not edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly." | rule | ADR-021 | **HUB-60** | yes | |
| CTB-34 | "Before review, digest a full change into the permanent corpus, update its plan bookkeeping and release-relevant `CHANGELOG.md` entry where applicable, and remove its `/changes/` workspace." | rule | PDR-002, C-002 | **DLT-24**, **DLT-26** | part (3 conditions) | |
| CTB-35 | "Plain editorial changes add no release note." | rule | C-002 | **HUB-52** | yes | |
| CTB-36 | "The final tree proposed for merge must not contain transient change files." | rule | PDR-002 | **HUB-29**, **VFY-14** | yes | |
| CTB-37 | "Keep runner labels, binary source, and writable install-directory choices in that caller; do not copy validation steps or action references into it." | rule | ADR-038 | no | yes | |
| CTB-38 | "A reusable-workflow reference update is the reviewed path for importing upstream scope, warning, acceptance-brief, and digest-gate fixes." | rule | ADR-038 | **CTB-37** | yes | |
| CTB-39 | `## Verify Locally` — "For a plain change, run only checks relevant to its changed surface. A guide-Markdown-only edit runs `git diff --check` and `npm run guide:build`." | rule | PDR-011 | **CTB-14**, **VFY-01c** | yes | |
| CTB-40 | "For a Cliewen change, commit the complete candidate, then run the repository's full mechanical gates against that commit" and the five commands | rule | PDR-012 | **VFY-24** | yes | |
| CTB-41 | "Total Go statement coverage must remain at least 80%." | rule | C-014 | no | yes | |
| CTB-42 | "`clue-verify` then automatically reviews that same commit before publication." | rule | PDR-012 | **VFY-01a**, **DLT-28** | yes | |
| CTB-43 | "A coding-agent host with context-isolated delegation starts a fresh read-only reviewer; other hosts disclose an in-context fallback." | rule | PDR-012 | **VFY-25**, **VFY-27** | yes | |
| CTB-44 | "The loop owns its classification regardless of the reviewer brief: a blocking finding is actionable and enters the hosted repair lifecycle; an advisory is a non-actionable observation" | rule | C-017 | **VFY-34**, **GL-46** | yes | |
| CTB-45 | "Counts and arithmetic disagreements are advisory, while a wrong, missing, or reused identity remains blocking, and the reviewer spends no pass re-deriving figures." | rule | C-017 | **GL-47** | yes | |
| CTB-46 | "A blocking finding returns to the implementing context, is committed, checked against that commit, and reviewed again — scoped to what changed and the carriers it declares — until the current commit receives a pass with no blocking findings." | rule | PDR-012 | **VFY-40**, **VFY-44**, **VFY-45** | part (4 conditions) | |
| CTB-47 | "An advisory repair may ride before a pass already required by a blocking repair; an advisory first reported by a pass with no blocking findings stays in the handoff for a later change" | rule | C-017 | **GL-49** | part (2 conditions) | |
| CTB-48 | "Three passes are the ordinary budget, and a fourth or later pass runs only after an immediately preceding pass with a blocking finding." | rule | C-017 | **GL-50** | yes | |
| CTB-49 | "The final verification evidence identifies the review mode, reviewed commit, number of passes run, and advisory findings left open." | rule | C-017, PDR-012 | **VFY-17**, **GL-73** | part (4 conditions) | |
| CTB-50 | `## Open the Pull Request` — "For a plain change, complete only the pull-request summary and relevant verification, then open the pull request after the applicable checks pass." | rule | PDR-011 | **CTB-14** | yes | |
| CTB-51 | "For a Cliewen change, also complete the template's proposal, traceability, and Cliewen checklist, and open the pull request only after the applicable checks and automatic agentic review pass." | rule | PDR-012 | **VFY-01b** | yes | |
| CTB-52 | "Keep review fixes on the same branch and pull request; for a Cliewen change, each substantive fix invalidates the earlier clean pass." | rule | PDR-012 | **F-RB-08**, **VFY-44** | **part — "substantive" is not the classification C-017 uses, and an advisory repair is substantive without invalidating** | |
| CTB-53 | "The branch and pull request are a proposal; a human maintainer performs the human-controlled merge commit that accepts a full Cliewen change." | rule | PDR-021, C-012 | **HUB-33**, **HUB-34**, **F-RB-02** | yes | |
| CTB-54 | "Configure the protected default branch to allow merge commits and disable squash and rebase-and-merge, so the proposal, implementation, digest, and durable corpus history remain reachable from `main`." | rule | PDR-021 | **HUB-34**, **F-RB-02** | yes | |
| CTB-55 | "Agents must never merge their own pull requests, create local merge commits into `main`, or push directly to `main`." | rule | C-012 | **HUB-33**, **F-RB-03** | yes | **!** |
| CTB-56 | "Cliewen does not currently require a Contributor License Agreement or Developer Certificate of Origin sign-off." | connective ×1 | — | — | — | |

**CTB-22 is DLT-06 alive on a second carrier.** AN-018's Q-02 found `clue-delta` instructing an agent to derive the next CH number by searching Git history, an instruction [ADR-048](../decisions/ADR-048-corpus-wide-id-ledger.md) superseded and whose replacement — the persisted ledger and `clue id next` — has shipped; following it produces a corpus `clue validate` rejects. M-063 repaired the skill. `CONTRIBUTING.md` gives a human contributor the same superseded instruction, and a contributor who follows it hits the same failure. The answer is already recorded, so this is a repair, not an escalation — but it is the clearest illustration of why this register was needed: one defect, found once, repaired on one of its two carriers.

**CTB-52 disagrees with C-017 on the same page that states C-017.** CTB-45 through CTB-49 carry the bounded loop and the intrinsic severity model; CTB-52, four paragraphs later, says every *substantive* fix invalidates the earlier clean pass. Under [C-017](../constraints/C-017-agentic-review-loop-is-bounded.md) an advisory repair is substantive and does not invalidate. Two obligations over one situation that pull a reader in different directions is PDR-029's definition of a conflict, and PDR-029 forbids an agent resolving one. Escalated as **Q-03** — though the direction looks settled, because C-017 is a constraint and CTB-52 is prose that CH-131 did not reach.

**CTB-55 is the merge prohibition in the final paragraph of the contributor's first carrier.** On P3 the contributor reads this file before the hub, so the earliest statement of the prohibition on that path is its last section, after every procedural instruction. The hub states it early (HUB-33), which rescues the *agent*; nothing rescues the human reading in the declared order.

---

## Register — CLI text (CLI)

Segmented by output unit. The command descriptions in `usage` are predominantly mechanism description addressed to someone deciding what to run, and are connective under the class rule stated above; the rows below are the units that carry an obligation, plus the connective runs they sit in.

| ID | Locator | Class | Trace | Dup | Chk | Ord |
|---|---|---|---|---|---|---|
| CLI-01…CLI-02 | the banner line; the `Usage:` synopsis block | connective ×2 | — | — | — | |
| CLI-03 | `init` — "Idempotent: existing files are never overwritten (they are reported and skipped)" | rule | ADR-019 | no | yes | |
| CLI-04 | `init` — a target directory that is a symlink is left untouched and reported as linked | rule | ADR-018 † | no | yes | |
| CLI-05…CLI-07 | `scaffold`, `context`, and `id next` descriptions | connective ×3 | — | — | — | |
| CLI-08 | `scaffold` — "prose outside the clue:index markers is never touched" | rule | ADR-019 | **CLI-03** | yes | |
| CLI-09 | `migrate` — "Existing prose and locally modified generated files are never overwritten." | rule | ADR-039 | no | yes | |
| CLI-10 | `id live` — "Refuses an ID that is not reserved." | rule | ADR-048 | no | yes | |
| CLI-11 | `refs` — "Only \"gone\" is an error; an outage elsewhere never condemns a corpus." | rule | ADR-040 | no | yes | |
| CLI-12 | `refs` — "A clue: identity is never followed." | rule | ADR-040 | no | yes | |
| CLI-13 | `refs` — "Never make this a required check: another host's uptime must not gate a merge." | rule | ADR-040 | no | yes | |
| CLI-14 | `parity` — "The report is derived and never edited by hand; a clean run is the only passing result." | rule | ADR-049 | no | yes | |
| CLI-15 | `report` — "A report states its criterion counts and mapping table only there, and validate re-renders the region and fails on any difference, so a typed figure is not a supported way to write one." | rule | ADR-054 | no | yes | |
| CLI-16 | `carriers` — "The report is derived and never edited by hand; a clean run is the only passing result." | rule | ADR-051 | **CLI-14** | yes | |
| CLI-17 | `carriers` — "A blocked entry has no target and is not reconciled — its presence is the record of a known gap." | rule | ADR-051 | no | yes | |
| CLI-18 | `validate` — the scanned graph and its checked properties | connective ×1 | — | — | — | |
| CLI-19 | `validate --forbid-changes` — "fail when /changes contains files — the digest-before-merge gate used by CI" | rule | PDR-002 | no | yes | |
| CLI-20 | `validate --coverage` / `--reality-gaps` / `--index-rows` — each "never a committed registry" | rule | PDR-028 | no | yes | |
| CLI-21 | `latest` — "Reaches the network, so it is never part of a validation verdict and must never be a required check." | rule | ADR-042 | **CLI-13** | yes | |
| CLI-22 | `latest` — "It writes no file in the repository and never replaces the binary." | rule | ADR-042 | no | yes | |
| CLI-23 | `latest --quiet` and `--timeout`; `version` | connective ×2 | — | — | — | |
| CLI-24 | "Release notice: init, scaffold, context, migrate, id, refs, and report print one line to standard error … Never from validate or version, never when CI carries a value, and never when CLUE_NO_UPDATE_NOTIFIER is set at all, the empty string included." | rule | ADR-042, PDR-023 | no | **part (4 conditions)** | |
| CLI-25 | "Exit codes: 0 corpus valid · 1 issues found · 2 usage error" | connective ×1 | — | — | — | |
| CLI-26 | `init`'s closing hint — "next: run `clue validate` — green on a fresh scaffold; then read docs/README.md" | connective ×1 | — | — | — | |
| CLI-27 | validate — "transient workspace present — digest before merge (main must never contain /changes)" | rule | PDR-002 | **CLI-19** | yes | |
| CLI-28 | `id next` — "identity ledger is missing; run `clue migrate --apply` first" | rule | ADR-048 | no | yes | |
| CLI-29 | validate — "derived region disagrees with … regenerate it with clue report, never by hand" | rule | ADR-054 | **CLI-15** | yes | |
| CLI-30 | migrate — a bare outward reference is listed and "a bare reference must never be repaired mechanically" | rule | ADR-040 | no | yes | |
| CLI-31 | migrate — the "resolve it by hand before resuming" family: ambiguous frontmatter, a provenance field on a decision, a `reversal-cost` that is neither low nor high | rule ×3 | ADR-039 | no | yes | |
| CLI-32 | migrate — "managed carrier differs from every supported generated release; local edits are never overwritten" | rule | ADR-039 | **CLI-09** | yes | |
| CLI-33 | migrate — "exists but never imports AGENTS.md … add a line containing just `@AGENTS.md` — migration does not edit a file you wrote" | rule | PDR-022 | no | yes | |
| CLI-34 | migrate — "never asks whether this repository is behind … add a line telling the agent to run `clue latest --quiet` when it starts — migration does not edit your hub" | rule | PDR-023 | no | yes | |
| CLI-35 | migrate — "corpus artifact is behind a symlink; resolve the repository-owned path before migrating" | rule | ADR-039 | no | yes | |

**The CLI is the cleanest carrier in this register, and its shape says why.** Almost every obligation it states is a prohibition on the tool itself — never overwrite, never rewrite, never gate a merge, never repair mechanically — and each traces to the decision that drew that boundary. The one duplication inside the help text (CLI-14 against CLI-16) is two commands stating the same derived-report rule in identical words, which is a candidate for one sentence stated once at the level both share.

**CLI-24 names seven notice-printing commands and is correct.** `AGENTS.md` and `guide/operations.md` name five. The CLI is the carrier that agrees with the code.

---

## Populations

Counted over the fourteen carriers at `9a632f9`. These are this spike's own results, they are **the only figures in this document**, and they are not hand-maintained: the script recomputes every one from the register rows above, and a reader who disagrees with a number runs it rather than trusting the prose. A row registering a run contributes its run length, which is the trailing number in its class cell.

```awk
# awk -f recount.awk AN-021-remaining-carrier-register.md
BEGIN{FS="|"}
/^\| (GI|GW|GD|GM|GC|GL|GS|GG|GA|GB|GO|GP|CTB|CLI)-/ && NF>=8 {
  class=$4; trace=$5; dup=$6; chk=$7; ord=$8
  gsub(/[* ]/,"",class); gsub(/[* ]/,"",trace); gsub(/[* ]/,"",dup); gsub(/[* ]/,"",chk); gsub(/[* ]/,"",ord)
  n=1
  if (match(class,/[0-9]+$/)) n=substr(class,RSTART,RLENGTH)+0
  total+=n
  if (class ~ /^connective/) { connective+=n; next }
  rule+=n
  if (dup != "no") duplicated+=n
  if (chk ~ /^part/) uncheckable+=n
  if (trace ~ /^NONE/) untraced+=n
  if (index(trace,"\342\200\240")>0) fresh+=n
  if (ord ~ /!/) misordered+=n
}
END{
  printf "total %d  connective %d  rule-bearing %d\n", total, connective, rule
  printf "traced %d  untraced %d  fresh-trace %d\n", rule-untraced, untraced, fresh
  printf "duplicated %d  uncheckable %d  misordered %d\n", duplicated, uncheckable, misordered
}
```

| | Statements |
|---|---|
| Total registered | 659 |
| Connective | 246 |
| Rule-bearing | 413 |
| — traceable to a live artifact | 407 |
| — of those, traced by this pass rather than inherited from AN-018 (†) | 18 |
| — **traceable to nothing found** | 6 |
| Rule-bearing statements duplicated in at least one reading path | 167 |
| Rule-bearing statements failing checkability (`part`) | 89 |
| Statements binding absolutely but read after what they constrain | 2 |

Three qualifications belong with these rows rather than after them.

**The connective share is the reason these carriers cannot be scored like skills.** More than a third of what these files say carries no obligation, and that is correct rather than waste: `design.md` exists to make an argument, and an argument that has been stripped to its rules is no longer an argument. PDR-029 already places connective statements outside the three conditions; the population is reported so nobody later mistakes the guide's size for its rule count.

**Duplication is the dominant defect here, as it was in the skills, but its shape is different.** In the skills the repetition was hub-against-fragment: one rule authored twice for one agent. Here it is a rule restated in four pages an adopter is told to read in order, in four different compressions. That is worse in one respect and better in another — worse because the wordings are not identical and a reader cannot tell which is authoritative, better because consolidating them requires no generator change and breaks no standalone-artifact guarantee.

**The six untraced statements are not six problems.** GL-77 and GO-59 are rules already owed a carrier — GL-77 is the orient-after-merge rule M-067 owes a constraint, GO-59 restates GA-27. GD-74 is a *withdrawn* rule still standing. GA-27, GA-28, and GA-29 are one rule with no record anywhere. So the untraced population is two live questions, one of which is already someone else's milestone.

**GD-71 is not in the untraced count and is the most serious single row in this register.** It traces cleanly — to PDR-012 — and is still wrong, because PDR-035 amended PDR-012 and the carrier kept the old reading. AN-018 named this exact blind spot: a statement that traces to a real decision whose meaning has since changed reads as correct and is not caught by tracing. That register found the blind spot by *running* an instruction (DLT-06). This one found it by registering a carrier a later change had edited around, which is a second, cheaper way to catch the same class — and it argues that whatever a change's carrier inventory claims, a register of the carriers it did *not* touch is worth running afterwards.

## Order assessment

Two statements bind absolutely and are read after the procedure they constrain:

| Statement | Carrier | Position |
|---|---|---|
| GL-62 ("The agent may publish the branch, but it never merges the pull request or pushes to `main`") | `guide/change-loop.md` | section 6 of 6 |
| CTB-55 ("Agents must never merge their own pull requests, create local merge commits into `main`, or push directly to `main`") | `CONTRIBUTING.md` | final paragraph |

Both are the human merge boundary, and both sit last. The guide as a whole is rescued by its reading paths — an evaluator meets the prohibition at GD-39 near the top of `design.md`, and an adopter at GW-14 — so the defect bites only a reader who opens one page directly, which is the ordinary case for a page reached from a search engine. `CONTRIBUTING.md` has no such rescue: it is the first file on the contributor's path, and its last paragraph is the first place the boundary is stated. Every other carrier here states its prohibitions where they belong, and the ten pages with no order flag are not an absence of evidence — they carry few absolute prohibitions at all, because the guide explains and the skills instruct.

## Editorial defects with no rule attached

Recorded so the follow-on change can carry them and so they are not mistaken for escalations. None requires a decision.

- **GC-20** links to `P-007-core-hardening.md` as this repository's "active campaign". P-007 completed on 2026-07-28; five campaigns have opened since, and the live one is P-013.
- **GG-12, GG-28** use `clue 0.7.0` as the printed-version example, and **GB-18…GB-23** use `0.7.0` throughout the arming procedure. GB-18 tells the reader to substitute their own version, which makes the CI-wall instance sound; the getting-started examples have no such note.
- **GL-18** states a derived measurement — transient and durable line counts from the first adopter-history measurement — without naming the analysis that produced it or when. AN-010 is the source. Under the durable-work rule a measurement may be stated when it is the point of the record and is stated with what produced it; this one names neither, on a page whose point is the loop rather than the measurement.
- **GB-47** says "the final pull request must be accepted with a merge commit" inside a procedure whose pull request is closed unmerged three paragraphs later. The referent is the reader's first real change, not the probe, and nothing on the page says so.

## Consolidation candidates for the follow-on change

Compatible overlap, not conflict. Each is one rule stated more than once on one reading path.

1. **The criterion-evidence contract** — GW-37…GW-42, GG-48…GG-53, GA-16…GA-24, GO-15…GO-19 on P1, and GM-14…GM-19 with GL-30…GL-32 on P2. Four statements of one contract on the adopter's path and two on the evaluator's, in six different compressions. Highest value and highest risk in this register: the wordings differ, and the differences may be meaningful. Consolidating means deciding which page owns the contract and reducing the others to a pointer — and the honest candidate for owner is `methodology.md`, because it is the page whose subject *is* the thread.
2. **The merge-shape configuration** — GW-29, GB-31…GB-32, GB-38, GB-41…GB-44, GO-20…GO-24, CTB-54. The same "allow merge commits, disable squash and rebase" instruction in three pages and the contributor guide, with the two-surfaces distinction restated three times inside `ci-wall.md` alone.
3. **The release notice's boundaries** — GG-36, GO-07, GO-32, GO-34, GO-35, CLI-24. One rule about where a notice may appear, stated twice inside `operations.md` and once more in `getting-started.md`. This is also where the five-versus-seven command discrepancy lives, so consolidating and repairing are the same edit.
4. **The plain route's exemptions** — GL-03, GL-04, GS-04, CTB-09, CTB-11, CTB-14, CTB-15, CTB-39, CTB-50. The contributor guide states the plain tier's definition, protected surfaces, route, and exemptions four times in four sections.
5. **The bounded review loop** — GL-46…GL-50 against CTB-44…CTB-49. Landed by CH-131 in both carriers in nearly identical words. Both are on different paths, so this is not scored as in-path duplication; it is listed because the two copies are what let CTB-52 contradict one of them without the contradiction being visible.
6. **The derived-report rule in the CLI** — CLI-14 against CLI-16, and CLI-15 against CLI-29. Identical sentences under two commands and one diagnostic.

## Escalations

Three, each naming the statement, what it traces to or fails to, what removing or retaining it costs, and the judgement required. Everything else this register found is either a repair whose answer is already recorded, or a consolidation candidate.

**Q-01 — `guide/design.md` states two rules the method has changed (GD-71, GD-74).** *Class: the decision behind the statement has outlived it, twice over, in two different ways.*

GD-71 says a fix invalidates the review pass and the loop ends only on a clean review. C-017 and PDR-035 say an advisory repair does not invalidate a clean pass, and that the ordinary budget is three passes. GD-74 states "small deltas" as a named principle with a rationale paragraph; the rule itself was withdrawn as uncheckable and removed from `clue-delta`.

Retaining either publishes a rule the method does not hold, in the page written to be read adversarially. Removing GD-71's wording is mechanical — the corpus already states the replacement. GD-74 is not mechanical, and that is the judgement: the *rule* is withdrawn, but the *reason* — two large parallel changes can be textually mergeable and semantically incompatible, and no tool notices — is a true and load-bearing part of why the one-change-in-flight rule and the reconcile-and-recheck rule exist. Those rules survive and trace to PDR-007 and PDR-016. So the question is whether `design.md` may keep the rationale as an explanation of rules that remain, having lost the unenforceable rule it was written to justify, or whether a withdrawn rule takes its argument with it. What is being asked: is a design page allowed to argue for a boundary the method no longer states as a rule of its own?

**Q-02 — the poor-fit conditions trace to nothing (GA-27, GA-28, GA-29, GO-59).** *Class: traces to nothing.*

Five conditions under which the method should not be adopted — work outside Git and a human merge boundary, no reliable tests or stable CI check, disposable or vendored code, evidence spread across repositories, a team that will not let agents maintain the corpus or review meaning before merge. `operations.md` links to them as its fallback and the sentence appears in adoption conversations. Nothing in the corpus states them; the nearest live artifacts are ARCH-003's extension clause, which says what may not be redefined rather than where the method does not apply, and ADR-044's repository boundary, which states one of the five.

Removing them costs the method its only published statement of its own limits, on a page an evaluator reads before adopting — and PDR-029's own argument is that a methodology that overclaims trains readers to stop reading its claims. Retaining them untraced leaves five load-bearing conditions that no change is obliged to keep true, and four of the five are about the *adopter's* situation rather than Cliewen's design, which is unusual for anything the corpus records. What human judgement is required: are these a real rule the corpus should state — most naturally as a constraint or as prose in ARCH-003's periphery section — or are they editorial advice that should say so, and stop being written in the imperative?

**Q-03 — `CONTRIBUTING.md` contradicts C-017 on the same page that states it (CTB-52 against CTB-44…CTB-49).** *Class: two obligations over one situation.*

CTB-52 says every substantive fix invalidates the earlier clean pass. Under C-017 an advisory repair is substantive and does not invalidate; that is the whole point of the bound, because a loop where any repair reopens the loop does not converge. A contributor following CTB-52 runs passes C-017 says are unnecessary; a contributor following CTB-48 publishes a candidate CTB-52 calls unreviewed.

The generated `clue-verify` carrier states the same clause and is rescued by its next sentence: after "every substantive edit invalidates it" it adds that an advisory from a clean pass stays in the handoff *rather than editing the clean commit*, so no advisory repair ever reaches the invalidating case. `CONTRIBUTING.md` carries the clause without that sentence. That is the repair's shape as well as the defect's.

The direction looks settled — a constraint outranks contributor prose, and CH-131 wrote the constraint and edited this file without reconciling this paragraph — so the cost of resolving it toward C-017 looks like zero. It is escalated anyway, and the reason is PDR-029's: an agent that resolves a conflict silently has made a methodology decision without a human, and the cheapness of a repair is not evidence that it is the right one. What is required is a one-line confirmation that C-017 governs and CTB-52 is redrafted to its language, or a statement of why the contributor's rule should be stricter than the agent's.

## Answers

Q-01 through Q-03 were answered by Flemming N. Larsen on 2026-08-08 in conversation, in the change that wrote this register, and recorded under [C-011](../constraints/C-011-decision-records-typed.md).

**Q-01 answered: the argument survives its rule.** `guide/design.md`'s review-loop wording is corrected to the bounded loop, and "small deltas" loses its bold heading and becomes the rationale under the branching and one-change-per-author rules it actually explains. What failed the test was stating an uncheckable rule as a named principle, not the observation behind it. Recorded as a [`log.md`](../decisions/log.md) row dated 2026-08-08.

**Q-02 answered: the poor-fit conditions move into the corpus as architecture, not as a constraint.** They are conditions on the adopting repository rather than rules a change can violate, so they now live in [ARCH-003](../architecture/core.md) beside the periphery and extension clauses, where [PDR-031](../decisions/PDR-031-architecture-artifacts-are-traces.md) makes them a valid trace. `guide/adoption.md` and `guide/operations.md` keep their wording and now trace there. Recorded as a `log.md` row.

**Q-03 answered: the constraint governs.** `CONTRIBUTING.md`'s sentence is redrafted to say that a blocking repair invalidates the clean pass and starts a new one while an advisory from a clean pass stays in the handoff — the sentence the generated skill already had. The answer then ranged wider than the question: asked what the bounded loop actually says, the human changed it. The budget is five passes rather than three, and a fifth pass that still returns blocking findings reports them and asks whether to continue instead of earning further passes on the loop's own authority. That is [PDR-036](../decisions/PDR-036-review-loop-budget-and-human-checkpoint.md), and it moved every carrier of the pass rule in the same change.

**Nothing in this register is left open.** The repairs that needed no judgement — the superseded CH-allocation instruction in `CONTRIBUTING.md`, the five-versus-seven notice-command lists in `AGENTS.md` and `guide/operations.md`, and the withdrawn tick-timing rule — were carried out in the same change. The editorial defects listed above and the six consolidation candidates are not: they are the follow-on change's work, and none of them is a rule.

## Rejected approaches

**Registering every sentence of the guide.** Started and abandoned inside `design.md`. Its nine observed failures of the first iteration, its problem statement, and its cost paragraph are argument: registered per sentence they produce dozens of rows reading `connective | — | — | —`, which buries the rule-bearing rows they surround and makes the document unreadable by the person it is written for. The run form keeps the population exact — a run states its length and the script reads it — while leaving the rows a reader must actually check.

**Scoring guide-against-skill repetition as duplication.** Refused for the reason stated under *Reading paths*: it would have produced a very large figure describing nobody's reading. The guide explains to a human what the skills instruct an agent to do, and that division is deliberate. A second reader who thinks the adopter and the agent are one reader will get a different duplication population, which is why the path set is written down rather than assumed.

**Deriving fresh traces for every row.** Refused for the opposite reason to the obvious one: the guide restates rules AN-018 already traced, and re-deriving them from the corpus would have produced a second, independent, differently-wrong trace column for the same rules. Reusing AN-018's traces makes the two registers consistent and inherits its errors, which is the honest trade and is marked in the column so a reader knows which entries had a prior pass behind them.

**Deciding the escalations.** All three have plausible-looking answers — delete GD-71's wording, write a constraint for the poor-fit conditions, redraft CTB-52. PDR-029 forbids exactly that, and Q-01 is the case that proves the rule: the mechanical half of it is obvious and the half that matters, whether a rationale outlives its rule, is not something an agent should settle.

## What this analysis does not establish

One reading of one revision by one agent, not independently re-segmented, so the segmentation rule's central claim — that a second pass yields the same statements — remains untested here as it was in AN-018.

The trace column is weaker than AN-018's in one direction and stronger in another. Stronger, because most entries are inherited from a column that has already survived an adversarial review. Weaker, because that inheritance is not independent evidence: if AN-018 mis-traced a rule, this register mis-traces it again and the agreement between the two registers will look like corroboration. The eighteen fresh traces are the rows with no prior pass behind them, and they are concentrated in `getting-started.md`, where they point at a capability rather than a decision.

Duplication is counted over the four reading paths named above and no others. A reader who arrives at one page from a search engine walks no path at all, and for that reader every rule is stated once — which is an argument for consolidation being safe, not for the count being wrong.

**Nothing here measures whether a human who reads these carriers does what they say.** That is the same limit AN-018 recorded, and it is sharper on this surface: the guide's audience has no obligation to follow it, so a rule stated four times in four compressions may be read once, partially, on a phone.

**The class rule that a description is connective was applied by one reader and it moves the population.** Under a stricter reading — where "the `clue` CLI checks structure, links, and traceability" obliges a reader to rely on nothing more — the rule-bearing count rises and the uncheckable count rises with it, because most such descriptions state several properties at once. The rule is stated in *What each column means* so a second reader can apply the other reading and see exactly which rows move.

## Consumer

[P-013](../plans/P-013-simplification.md)'s **M-064**, which owns this surface; [AN-022](AN-022-remaining-surface-scored.md) carries the same milestone's scoring of the rules, artifact types, required fields, commands, and checks. Q-01 through Q-03 were answered and their repairs carried out in the same change that wrote this register, because a decision that changes a methodology contract moves its carriers with it. What is left for the follow-on change is the six consolidation candidates and the four editorial defects, none of which is a rule. **M-067** consumes GL-77, whose constraint it already owes.

