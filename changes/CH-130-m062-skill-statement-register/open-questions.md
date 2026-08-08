---
id: CH-130-questions
type: open-questions
status: open
links: [CH-130, P-013, PDR-029, C-011]
title: Open questions from the M-062 statement register
---

# CH-130 — open questions

Seven questions from [AN-018](../../docs/analysis/AN-018-skill-statement-register.md)'s register. [PDR-029](../../docs/decisions/PDR-029-simplification-tests-by-surface.md) forbids an agent settling any of them: silent deletion and silent retention are the same failure. Each carries a recommendation, which is a proposal for the human to accept, reject, or replace — not a decision. Answers become decision records under [C-011](../../docs/constraints/C-011-decision-records-typed.md), and M-063 applies them.

Q-02, Q-05, and Q-06 block M-063 directly. Q-01, Q-03, Q-04, and Q-07 can be answered later without stalling the trim, but each leaves a statement M-063 must not touch until it is answered.

---

## Q-01 — `clue-analysis`'s workflow spine traces to nothing the method accepts

**Class:** a rule-bearing statement that traces to nothing.

**Statements:** ANL-02 ("Use when a change has unclear risks or unknowns — before planning or implementing"), ANL-03 ("retire the biggest risk first"), ANL-04 ("Name the risk or unknown in one sentence"), ANL-05 ("If you cannot, that is the first finding"), ANL-17 (what a spike is), ANL-18 ("Spikes are disposable; their findings are not"), ANL-19 ("End every spike with a findings document in `/docs/analysis/`"), ANL-21 ("Include what was tried, what was rejected, and why"), ANL-25 ("A rejected alternative that is itself a decision gets a rejected decision record"), ANL-26 ("Feed findings to `clue-plan` or `clue-delta`"), ANL-27 ("Analysis with no consumer is doc-slop").

**What they trace to.** Eight of the eleven trace to [AN-001](../../docs/analysis/AN-001-foundation-v0.4.md), the frozen Foundation Document, which is an analysis rather than one of PDR-029's four accepted types and whose own banner states that the corpus wins where the two disagree. Three trace to nothing found at all. The rest of the skill — its evidence-discipline rules — traces cleanly to two decision-log rows.

**What retaining them costs.** Eleven statements in the skill that governs how Cliewen writes its shared memory would remain unanswerable to the adversarial question PDR-029 says a reviewer asks: *why do you have this statement?* The answer would be "the founding document said so", which is the answer PDR-029 explicitly rules out.

**What removing them costs.** The whole of `clue-analysis`'s spine. The skill would state how to bound evidence and nothing about when to run a spike, what a spike is, that it must end in a findings document, or that findings need a consumer. `/docs/analysis/` holds seventeen findings documents produced under exactly these rules.

**Judgment required.** These read as real rules nobody ever recorded — the repair PDR-029 names for that case is writing the missing decision, not deleting the sentence. But that is eleven statements, and writing one decision to bless all eleven at once would be a decision written to make a register green rather than because someone decided something.

**Recommendation.** Write one PDR — *analysis is a bounded spike with a named consumer* — covering ANL-02…05, ANL-17…19, ANL-21, ANL-26, ANL-27, because those ten are one coherent rule about the lifecycle of a spike and have been practised without exception across seventeen findings documents. Handle ANL-25 separately: "a rejected alternative that is itself a decision gets a rejected decision record" is a rule about decision records, not about analysis, and it belongs with [PDR-006](../../docs/decisions/PDR-006-decision-records-are-typed.md) as an amendment rather than in a new record. The reason for one PDR rather than ten is that the ten statements are not independent — removing any one of them makes the others incoherent — and a decision record that states a rule its author cannot decompose is more honest than ten records that pretend it decomposes.

---

## Q-02 — `clue-delta` instructs an allocator the shipped tool has replaced

**Class:** a statement whose decision has outlived its reason. **Blocks M-063.** This is the register's one demonstrated defect rather than a judgement call, and it was found by running the tool rather than by reading.

**Statement:** DLT-06 — "Take the next free CH number by searching Git history and `/changes/` for the highest used number."

**What it traces to.** [ADR-009](../../docs/decisions/ADR-009-ac-id-namespaces.md)'s "the corpus is the registry" clause, which [ADR-048](../../docs/decisions/ADR-048-corpus-wide-id-ledger.md) supersedes for *every* native ID prefix including `CH`. ADR-048 replaces scan-and-max with a persisted ledger precisely because [ADR-034](../../docs/decisions/ADR-034-retirement-is-deletion.md) makes retirement a file deletion, after which no scan can see the dead ID and nothing stops it being re-minted.

**The ledger is built and live.** ADR-048's closing paragraph says the ledger, `clue id next`, the `checkLedger` rule, and the migrate backfill are follow-up work under P-011's M-052. That paragraph is now history: `.clue/id-ledger.yaml` exists in this repository with a `CH` counter, `clue id next <prefix>` and `clue id live <id>` are shipped subcommands, and `clue validate` rejects an artifact whose ID is absent from the ledger. The skill was not updated when the mechanism landed.

**What retaining it costs.** Measured, not argued: this change followed DLT-06 as written, derived `CH-130` by grepping Git history, and produced a corpus `clue validate` rejected — twice, once for `CH-130` and once for `AN-018`. The skill's instruction does not merely trace to a superseded clause; it routes an agent into a failing validate and gives no hint that a command exists to do the job. An adopter following it can also re-mint a deleted artifact's number, which is the correctness gap ADR-048 was written to close.

**What changing it costs.** DLT-06 is in `clue-delta`, and the same scan-and-max idea appears in `clue-extract`'s minting rules (EXT-14), which are about a different namespace and a different situation. A repair that only fixes `clue-delta` leaves the two carriers disagreeing about how IDs are allocated.

**Judgment required.** Whether DLT-06 is replaced by an instruction to run `clue id next`, and whether ADR-048's superseded implementation paragraph is corrected in the same change — and whether that reaches EXT-14 or stops at `clue-delta`.

**Recommendation.** Replace DLT-06 with `clue id next CH` and treat this as a methodology-carrier repair under [PDR-019](../../docs/decisions/PDR-019-methodology-contract-carriers-move-together.md) rather than a wording fix: ADR-048's "does not itself implement" paragraph is a live carrier now stating something false, and the skill, the scaffolded copies, and that paragraph move together. Leave EXT-14 alone — brownfield minting preserves source IDs and mints into an `ac-prefix:` namespace from source order, which is a genuinely different rule from allocating the next native `CH`, and folding them together would be the kind of consolidation PDR-029 warns produces one ambiguous rule from two clear ones. This is not a log row: it changes what four shipped skills instruct.

---

## Q-03 — the no-computed-figures rule traces to nothing

**Class:** a rule-bearing statement that traces to nothing.

**Statements:** F-DW-03 ("A durable record never states a figure a command computes — an artifact count, a coverage percentage, a reported population size. Name the command instead") and F-DW-05 ("Measurements that are the point of a record … are stated with what produced them and when"), in the shared `durable-work` fragment rendered into `clue-delta`, `clue-upgrade`, and `clue-verify`.

**What they trace to.** Nothing found. The nearest live artifact is [ADR-054](../../docs/decisions/ADR-054-derived-extraction-report-region.md), which makes an *extraction report's* figures a rendered region bounded by markers — a much narrower rule about one mechanism in one document type. [PDR-028](../../docs/decisions/PDR-028-derived-report-is-not-a-committed-registry.md) is adjacent for the same reason and equally narrow. Neither states a general prohibition on stating computed figures in durable prose.

**What retaining them costs.** A prohibition binding every durable artifact this project writes rests on nothing a reader can open. It is also the most frequently *applied* rule in the set — it shapes how every plan, decision, and analysis states its evidence, including AN-018 itself, which had to justify its own population table against it.

**What removing them costs.** Prose reacquires hand-maintained numbers. The rule's own rationale (F-DW-04) states the failure mode it prevents: a figure written into prose goes stale on the next change, every later reviewer re-derives it, and repairing one writes new prose carrying new numbers, so the finding regenerates instead of converging.

**Judgment required.** Whether this is a general rule that was never recorded, or an over-generalisation of ADR-054's narrow mechanism that spread into a shared fragment because it sounded right.

**Recommendation.** Write it as a constraint rather than a decision. It behaves like C-001 and C-006 — a prose convention with a human enforcement class and a stated residual — and [ADR-017](../../docs/decisions/ADR-017-conventions-are-constraints.md) already says prose conventions register as constraint artifacts. A constraint also gets it into the register that [ADR-045](../../docs/decisions/ADR-045-register-names-the-machine.md) requires to name the machine that holds it, which for this rule is honestly "none — human judgement", and saying so is worth more than the rule's current silence.

---

## Q-04 — the orient-after-merge instruction traces to nothing

**Class:** a rule-bearing statement that traces to nothing.

**Statement:** F-RB-09 — "After a human reports the merge, orient before starting anything else: describe the plan's next unfinished step in plain language and ask whether to start it, or say that the plan has nothing left and ask what comes next."

**What it traces to.** Nothing found. It is the last statement of the `review-boundary` fragment, rendered into four skills.

**What retaining it costs.** Little in practice, and one thing in principle: it is the only statement in the review boundary that is about conversational behaviour rather than about the boundary, and it is stated in the fragment named after the boundary.

**What removing it costs.** The behaviour it produces is the reason this very question is being asked in a session that began with "what's next?". Removing it costs the handoff that keeps a campaign moving between sessions.

**Judgment required.** Whether this is a rule at all, or a helpful habit that drifted into a normative carrier — and if it is a rule, whether the review boundary is where it belongs.

**Recommendation.** Record it as a decision-log row and move it out of `review-boundary` into `durable-work`, whose subject — what survives a change of agent — is what this rule is actually about. The move is M-063's work; the log row is what authorises it. A log row rather than a PDR because reverting it restores one sentence to one fragment.

---

## Q-05 — does an architecture artifact count as a trace?

**Class:** definitional. **Blocks M-063**, and blocks this register's own verdict on three statements.

**Statements:** HUB-25 ("The `/docs` corpus remains the system-of-record and working memory"), HUB-56 (the core's three-element definition), HUB-59 (the skills routing table). All three trace to `ARCH-003` or `ARCH-004`.

**The problem.** PDR-029 names four artifact types as valid traces: decision, constraint, goal, acceptance criterion. Architecture is not among them. Yet `AGENTS.md` rule 8 cites `ARCH-003` by ID as the source of the core definition, and [ARCH-003](../../docs/architecture/core.md) is the durable statement of what the core *is* — the thing [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) protects. Reading PDR-029 literally makes the routing hub's most load-bearing rule untraceable.

**What treating architecture as a valid trace costs.** It widens the accepted set from four types to five, and architecture files are more narrative than decisions — the risk is that "traces to an architecture file" becomes the loophole that lets a statement pass on prose it can be argued from, which is the derivability failure PDR-029 exists to prevent.

**What refusing it costs.** Three statements in the hub are reported untraceable when the artifact stating them is the one the constraint protecting them points at. It would also mean M-063 must either delete the core definition from the hub or write a decision record restating what ARCH-003 already says.

**Judgment required.** Whether PDR-029's four-type list was a considered boundary or an omission.

**Recommendation.** Amend PDR-029 to accept an architecture artifact as a trace under the same restriction the other types already carry: it traces when the architecture file *states* the rule, and not when the rule is merely derivable from it. ARCH-003 states the core's three elements in numbered form and states the red line in one sentence; that is a statement, not an intent to be argued from. The restriction is what keeps this from being the loophole — and it is the same restriction that already stops G-001 from validating the entire skill set at once.

---

## Q-06 — plain changes and unmerged work pull in different directions

**Class:** a pair of obligations over one situation. **Blocks M-063.**

**Statements:** HUB-16 — "Plain changes do not consume the one-Cliewen-change-in-flight slot and never build on unmerged work." Against F-RB-01 — "If work must build on an unmerged change, record a blocking open question and stop unless the human explicitly authorizes it."

**The divergence.** The hub states an absolute prohibition scoped to plain changes. The fragment states a general rule with an explicit human-authorisation escape and does not exempt plain changes from it. A reader who meets both can reasonably conclude either that a human may authorise a plain change to build on unmerged work, or that plain changes are the one case where no authorisation is available. The fragment is the more authoritative carrier by placement; the hub is the more specific by scope.

**What each reading costs.** Reading the hub as absolute means a human cannot authorise something they are permitted to authorise for every other tier — an odd asymmetry, since plain changes are the *cheapest* class. Reading the fragment's escape as reaching plain changes means the hub's "never" is not a never, which is the ambiguity PDR-029 warns about: a rule obeyed inconsistently while appearing to be obeyed.

**Judgment required.** Which one is the rule. This is exactly PDR-029's conflict case, where an agent may not preserve the union or pick a winner.

**Recommendation.** The hub's absolute reading is the one to keep, and the fragment should be narrowed to say so. A plain change carries no CH identity and no proposal, so there is nowhere to record the authorisation the fragment's escape depends on — the escape is not merely undesirable for plain changes, it is unrecordable, and an authorisation that leaves no artifact is the thing [PDR-007](../../docs/decisions/PDR-007-review-boundary.md) built the escape to avoid. Recording that as a decision-log row and letting M-063 add the exemption clause to F-RB-01 resolves the pair without weakening either rule.

---

## Q-07 — `clue-extract` is missing from the hub's skill table

**Class:** coverage. Raised as a question rather than fixed, because the fix is a rule-bearing routing row.

**Statement:** HUB-59, the `## Skills` table in `AGENTS.md`. It lists five skills; six ship. `clue-extract` is absent. The scaffolded adopter hub (`internal/scaffold/templates/AGENTS.md`) lists all six.

**What leaving it out costs.** [ADR-043](../../docs/decisions/ADR-043-upgrade-skill-is-a-managed-carrier.md) states a managed set of six, and the hub is this repository's routing surface. An agent working here that needs `clue-extract` — most plausibly to change it — is routed to it by nothing. It also makes this repository's hub and the hub it ships diverge in what they claim the skill set is.

**What adding it costs.** One row, and a slightly less honest table: Cliewen has no source corpus to extract, so the row would route to a skill nobody here will run for its stated purpose.

**Judgment required.** Whether the hub's table is *the skill set* or *the skills this repository uses*. Both readings are defensible and the file does not say which it is.

**Recommendation.** Add the row. The table's heading is "Skills" and its column is "When", so it already reads as the set rather than as a usage list; and this repository is the reference implementation, where a divergence from the scaffolded hub reads to an adopter as a statement about the skill rather than about Cliewen's own corpus. The "When" cell can say what is true — brownfield adoption of an existing corpus — which is honest whether or not anyone here runs it. Cheap and local to reverse, so a decision-log row carries it.
