---
id: AN-022
type: analysis
status: active
links: [P-013, M-064, PDR-013, PDR-029, PDR-026, ADR-034, ADR-035, ADR-045, ADR-048, AN-008, C-013]
title: The rules, types, fields, commands, and checks scored against does the core need it, and pattern C determined
---

# AN-022 — The non-prose surface scored, and pattern C determined

**Primary-consumer exception:** anyone reading the Pattern C determination needs the non-prose scoring that precedes it to interpret the verdict, so these rendered documents remain together as one analysis.

## The risks this spike retires

Two, and they are different in kind.

**The scoring.** [PDR-013](../decisions/PDR-013-explicit-core-red-line.md) supplies *does the core need it?* and [PDR-029](../decisions/PDR-029-simplification-tests-by-surface.md) assigns it the surface it governs: rules, artifact types, required fields, commands, checks. The risk is that P-013 closes having applied that test to nothing — the criterion has existed since P-005 and has never been run over a measured population, which is exactly how simplification stalled the previous three times.

**Pattern C.** [AN-008](AN-008-methodology-critiques.md)'s pattern C — the graph only accumulates — has been carried forward by three campaigns in one sentence. The risk is a fourth repetition, or the opposite error: closing it because later work looks like it addressed the same words. [PDR-026](../decisions/PDR-026-campaigns-close-on-re-derived-evidence.md) requires the status to be re-derived from the corpus and the tool rather than read off a table, so this spike re-derives each of pattern C's four claims at head and reports which are answered, which are not, and which turned out to be false as stated.

Nothing is removed by this spike, and no candidate it approves is carried out here.

## Evidence boundary

- **Pinned revision:** `9a632f9` on `main` (`cliewen/cliewen`), the tip this change branched from, plus the working tree of this change, which adds only analysis files and its own workspace. The scored populations are read from `internal/corpus`, `cmd/clue`, and `docs/`; no scored file is edited by this change.
- **Environment:** prepared, not clean. Windows 11, Go 1.26.5, `clue` built from this checkout and stamped `dev`. The `dev` stamp matters for one result: skill-to-binary drift is not detectable from a checkout build, so the `checkSkillVersions` row below is scored on its code rather than on an observed verdict.
- **Confidence classes.** The populations are **observed** — they are enumerations of code and files, and every one is derived by a named command rather than counted by hand. Each *verdict* is an **inference**: it is one reader's judgement about what the core needs, and PDR-013 supplies the criterion but not a procedure. Where a verdict turns on the criterion itself being the wrong instrument, that is reported as a finding rather than smoothed over.
- **Not established.** No candidate here has been carried out, so no verdict has been tested by removal. Two verdicts rest on empirical facts about *this* repository — its inferred-decision population and its use of `supersedes:` — and an adopter repository could differ.

## The commands that produce every figure in this document

Stated once, here, so no figure below is a hand-maintained number:

```sh
# constraints, and their enforcement classes
ls docs/constraints/C-*.md | wc -l
for f in docs/constraints/C-0*.md; do grep -m1 '^enforcement:' "$f"; done | sort | uniq -c

# artifact types the validator recognizes
grep -c '^\t"' internal/corpus/rules.go            # statusVocabExceptions entries
grep -rh '^type: ' docs/ | sort | uniq -c

# checks clue validate runs
grep -c 'issues = append(issues, check' internal/corpus/rules.go

# decisions, their provenance, and prose supersession claims
ls docs/decisions/[AP]DR-*.md | wc -l
grep -l '^status: inferred' docs/decisions/*.md | wc -l
grep -o '^- \[\(ADR\|PDR\)-[0-9]*' docs/decisions/README.md | wc -l
grep 'supersed\|amend' docs/decisions/README.md | grep -o '^- \[\(ADR\|PDR\)-[0-9]*' | wc -l

# machine-visible supersession edges
grep -rln '^supersedes:' docs/ | wc -l
```

## What *does the core need it?* decided, and where it broke

The criterion was applied to five populations. It discriminated cleanly on one of them and under-determined on the other four, and the failure has a single shape worth stating before the tables.

**The core is a set of guarantees; a command, a check, a field, and a type are means.** [ARCH-003](../architecture/core.md) names three elements: the verifiable thread, the human merge boundary, and the deterministic judge. Asked of a *rule that reaches into meaning*, "does the core need it?" is decisive — C-013 and C-012 are the core, C-001 is not, and nobody has to argue about it. Asked of a command, it returns *no* for eleven of thirteen, because the core needs `clue validate` and needs nothing else by name. That is not a list of eleven removals. It is the criterion measuring the wrong thing: `clue id next` is not needed by the core, and it exists because a deleted artifact's identity could otherwise be re-minted, which breaks the criterion-identity guarantee the thread depends on. One step removed from the core, and invisible to the test.

PDR-029 already recorded that the criterion returns the wrong answer on carrier prose, where it passes nearly everything. **This spike finds the mirror-image failure on the tooling surface, where it fails nearly everything**, and the cause is the same: the criterion asks whether the core *needs* a thing, and both surfaces contain things the core does not need and cannot do without.

For the tooling surface the question that does discriminate is: **does removing it move an obligation from a machine back to a human?** That test is checkable, it is not a matter of taste, and it is the one the project actually applied historically — every command in `clue` exists because a ritual was being skipped by people who agreed it mattered. It is proposed rather than adopted, as **Q-04**.

The tables below therefore carry two columns: the criterion's verdict, and what removing the item would cost. Where those disagree, the disagreement is the finding.

## Rules — the constraint register

Seventeen constraints. Enforcement classes as declared: four `machine`, five `partial`, one `agent`, seven `human` — derived by the command above.

| Constraint | Core needs it? | What removal would cost |
|---|---|---|
| C-001 no hard-wrapped Markdown | no | Corpus diffs stop being reviewable line by line, which is how a human verifies meaning at the merge gate. Machine-enforced at zero human cost. |
| C-002 changelog per user-visible change | no | The release body has no reviewed source; [ADR-012](../decisions/ADR-012-release-notes-from-changelog.md)'s failure mode — auto-generated notes — returns. |
| C-003 tasks tick immediately | no | Nothing; the obligation is on a transient artifact the digest deletes. **The strongest candidate in this register.** See below. |
| C-004 never weaken checks | **yes** | The judge stops being deterministic in the only way that matters: a check that may be weakened to pass is not a check. |
| C-005 proposal declares plan item | **yes** | The plan→change edge of the thread is unenforced. |
| C-006 decision records timeless with carrier | no | Decisions accumulate episode prose and stop being readable later; [PDR-019](../decisions/PDR-019-methodology-contract-carriers-move-together.md)'s carrier obligation loses its home. |
| C-007 diagrams inline Mermaid | no | Diagrams become unreviewable binaries. Candidate; cost is small and real. |
| C-008 completed plans immutable | no | The plan index stops being an achievement record; a closed campaign can be rewritten to match a later story. |
| C-009 type-specific frontmatter | **yes** | The judge cannot read the graph. |
| C-010 milestone status vocabulary | **yes** | The judge reads milestone status cells; without the vocabulary they are free text. |
| C-011 decision records typed | no | The human's answers stop being recoverable, which is the merge boundary's memory rather than the boundary itself. |
| C-012 agents never merge own changes | **yes** | It *is* the human merge boundary. |
| C-013 core changes need a decision | **yes** | It *is* the red line. |
| C-014 coverage floor | no | A product quality bar for this repository, admitted deliberately by [ADR-027](../decisions/ADR-027-quality-scenarios-are-constraints.md). Removal costs this repository's own test discipline and costs the method nothing. |
| C-015 onboarding under 30 minutes | no | Same class as C-014, and it is CAP-001's own bar. |
| C-016 index rows state their record | no | Indexes become link lists; navigation degrades. Machine-counted, never failed on. |
| C-017 agentic review loop is bounded | no | The loop stops converging, which costs tokens and reviewer passes rather than a core guarantee. |

**Six of seventeen are needed by the core.** Eleven are not, and none of the eleven should be removed on that ground — which is the under-determination stated above, in its clearest form. Two of the eleven are deliberately admitted repository quality bars (C-014, C-015) that ADR-027 folded into this register on purpose; the rest each hold a defect the project watched happen.

**C-003 is the one candidate this spike would put to a human on the criterion alone.** It obliges an agent to tick a task the moment it completes, on a file the digest deletes, in a workspace that never reaches `main`. Its cost if removed is not zero — a batch-ticked checklist hides which task was in flight when a change stopped, which is what a handoff reads — but that value is about handoff, not about meaning, and [C-011](../constraints/C-011-decision-records-typed.md) plus the workspace already carry handoff. It is raised as **Q-05** rather than scored as waste, because a constraint that governs how an agent works while working is exactly the kind of rule a human should decide to keep or drop.

### A defect the scoring found on the way

**Five of the seventeen constraints name a `source:` that no longer resolves.** C-001, C-002, C-004, C-005, and C-013 name "AGENTS.md rule 3/5/6/7/8". `AGENTS.md` has carried no numbered rule list since CH-132 restructured it; its only numbered items are the three change tiers. C-004 additionally names "clue-verify preamble", and after CH-132 reordered that skill the text in question is no longer a preamble — it sits after the change-tier and review-boundary sections.

M-067 already owns the rule that "a constraint's `source:` must resolve to a live file rather than merely be non-empty". This is its population and its motive: the field is currently checked for non-emptiness only, so five constraints point at prose that does not exist and `clue validate` is green. It is not escalated, because the answer is already recorded; it is reported here so M-067's change has the list.

## Artifact types

Fourteen types the validator recognizes by name — eleven durable, three transient — plus adopter-defined types, which resolve to the default lifecycle by design ([ADR-026](../decisions/ADR-026-unknown-types-are-adopter-extensions.md)).

| Type | Core needs it? | What removal would cost |
|---|---|---|
| goal, plan, capability, criteria | **yes** | They are four of the thread's six links. |
| change | **yes** | The thread's fifth link and the digest boundary's subject. |
| architecture | **yes** | ARCH-003 — the core's own durable statement — is an artifact of this type, and [PDR-031](../decisions/PDR-031-architecture-artifacts-are-traces.md) made the type a valid trace. |
| decision | **yes** | C-013 requires a decision record for a core change; the red line cannot function without the type. |
| constraint | **yes** | C-004, C-012, and C-013 are constraints. The judge's protection and the red line live in this type. |
| analysis | no, then yes | Not needed by the thread. But [ADR-035](../decisions/ADR-035-cost-bounds-inferred-provenance.md)'s edge from reality is carried by an analysis with `reality: contradicted`, so the one feedback edge the core admits has no other home. |
| design | no | Per-capability implementer view. Removal would push design into the capability README or into decisions, and [`guide/corpus.md`](../../guide/corpus.md)'s one-home-per-scope rule exists because a fact with two homes disagrees with itself. Candidate with a stated cost. |
| log | no | [PDR-006](../decisions/PDR-006-decision-records-are-typed.md)'s cheap route. Removing it makes every decision expensive, which is the all-or-nothing ceremony PDR-006 repaired. |
| tasks | no | See C-003. |
| open-questions | no | The stopping mechanism's file. C-011 states the rule; the type gives it a location the judge can see. |
| imported-change | no | [ADR-050](../decisions/ADR-050-in-flight-source-work-is-durable.md)'s record of in-flight source work. Extraction-only; removal loses the migration's only durable record of work that was mid-flight when the source corpus died. |

**Eight of fourteen are needed by the core; one (`analysis`) becomes needed one step removed.** `design` is the only type this spike would put forward as a candidate on the criterion, and its cost — reopening where per-capability design lives — is larger than the saving. No type is recommended for removal.

## Required fields

| Field | Where | Core needs it? | What removal would cost |
|---|---|---|---|
| `id`, `type`, `links` | every artifact | **yes** | The judge cannot read the graph at all. |
| `status` | every artifact | **yes** | Nothing distinguishes a draft promise from an active one, so the evidence contract cannot be applied. |
| `title` | every artifact | no | Index rows and `clue context` output lose their human-readable handle; [ADR-046](../decisions/ADR-046-index-rows-say-what-the-artifact-is-about.md) needs it. |
| `author`, `accepted-by` | decision | **yes** | The merge boundary's signature. `accepted-by` is the only machine-visible record that a human endorsed a decision. |
| `goal` | capability | **yes** | The thread's first edge. |
| `source`, `enforcement` | constraint | no | [ADR-045](../decisions/ADR-045-constraints-name-their-holder.md)'s priced residual disappears, and uniform protection becomes uniform friction — the failure PDR-013 names. |
| `provenance`, `reversal-cost` | extracted non-decision | no | ADR-035's bound on inferred meaning; without them high-cost inferred meaning can back an active capability silently. |
| `reality` | analysis | no | The edge from reality has no carrier. |
| `ac-prefix` | criteria | **yes** | Brownfield criterion identities stop being stable, which breaks the thread's evidence edges. |
| `supersedes` | any | no | **Nothing today.** No artifact in the corpus carries this field. See the pattern C determination below — this is where the two halves of this spike meet. |

## Commands

Thirteen commands, counting `id next` and `id live` separately.

| Command | Core needs it? | Does removing it move work from a machine to a human? |
|---|---|---|
| `validate` | **yes** | It is the judge. |
| `init` | no | Yes — every adopter hand-assembles the taxonomy, hub, skills, and CI caller, which is the acquisition cost the first iteration failed on. |
| `scaffold` | no | Yes — indexes return to being hand-maintained, which is one of the observed failures that produced Cliewen. |
| `context` | no | Yes — [PDR-034](../decisions/PDR-034-the-corpus-is-read-narrowly-by-default.md) obliges narrow reading; without the command the obligation has no instrument. |
| `migrate` | no | Yes — a coordinated upgrade becomes a hand-edited one, and [ADR-039](../decisions/ADR-039-versioned-corpus-migrations.md) exists because that drifted. |
| `id next` / `id live` | no | Yes — and this one protects a core guarantee: [ADR-048](../decisions/ADR-048-corpus-wide-id-ledger.md)'s ledger exists so a deleted artifact's identity is never re-minted, which is criterion-identity immutability. |
| `refs` | no | Partly — a human can open a link. What no human does reliably is open every link in a growing corpus. |
| `parity` | no | Yes — deleting a source corpus without a mechanical parity check is exactly the evidence loss [ADR-049](../decisions/ADR-049-migration-parity-manifest.md) was written for. |
| `carriers` | no | Yes — same shape, for operational carriers ([ADR-051](../decisions/ADR-051-pinned-carrier-inventory.md)). |
| `report` | no | Yes — [ADR-054](../decisions/ADR-054-derived-extraction-report-region.md) makes an extraction report's figures rendered precisely because typed figures went stale. |
| `latest` | no | Yes — nobody remembers to ask whether they are behind, which is [PDR-023](../decisions/PDR-023-tool-notice-and-hub-instruction.md)'s whole reason. |
| `version` | no | No — but it is a single line of output and the only command guaranteed to answer offline forever. |

**One of thirteen is needed by the core. Eleven of thirteen move work back to a human if removed.** `version` is the only command that fails both tests, and it costs nothing. This is the clearest evidence for Q-04: on this surface the criterion produces a near-total *no* and the second question produces a near-total *keep*, and the second question is the one that matches every decision the project actually made.

## Checks

Twenty-two check functions plus the `--forbid-changes` gate.

| Check | Core needs it? |
|---|---|
| `checkCoreFields`, `checkDuplicateIDs`, `checkLinks`, `checkStatusVocab`, `checkTypeFields` | **yes** — the judge's definition of a well-formed graph |
| `checkACTests` | **yes** — the thread's last edge, and the reason `clue` exists |
| `checkProposalPlanItem` | **yes** — C-005, the plan→change edge |
| `checkMilestoneStatus` | **yes** — C-010, read by the plan layer |
| the `--forbid-changes` gate | **yes** — the digest boundary a merge may not cross |
| `checkProvenance` | no — ADR-035's bound on inferred meaning |
| `checkReality` | no — the edge from reality |
| `checkSupersedes` | no — and unexercised; see below |
| `checkConstraints` | no — ADR-045's priced residual |
| `checkFolderReadmes`, `checkIndexes` | no — navigability, and the repair of a watched failure |
| `checkExternalReferences`, `checkForeignPointers` | no — ADR-040's reference honesty |
| `checkFrontmatterHygiene`, `checkProseLayout` | no — C-001 and form |
| `checkSkippedTasks` | no — C-003 |
| `checkSkillVersions` | no — ADR-011's drift rule |
| `checkLedger` | no, then yes — ADR-048, for the identity reason given above |
| `checkImportedChanges` | no — ADR-050 |

**Nine of twenty-three verify something the core needs; the rest verify form, memory, or a bounded lifecycle.** No check is a removal candidate: a check is the cheapest possible location for an obligation, it costs a human nothing, and every one of the fourteen that the core does not need is holding a rule that would otherwise be enforced by memory. `checkSupersedes` is the single exception and it is not a simplification candidate either, for the reason the next section gives.

---

# Pattern C — determined

[AN-008](AN-008-methodology-critiques.md)'s pattern C is one paragraph making five claims. Re-derived at head, they do not have one status. Three are answered, one is false as stated, and one is open — and the open one is narrower and sharper than the paragraph suggests.

**Claim 1 — "Nothing in the corpus bounds, reverses, or consumes state."** Superseded by claims 2–5, which are its instances. Not scored separately.

**Claim 2 — "The born-`inferred` population can only grow, so its counter stops being read." → False as stated at head.** Of 84 decision records, one carries `status: inferred`: PDR-035, written by the change immediately before this one and not yet human-approved. The population is not monotonic, and the counter has not stopped being read — `clue validate`'s OK line reports it on every run and the human has evidently been promoting. For non-decision meaning the bound is mechanical rather than social: [ADR-035](../decisions/ADR-035-cost-bounds-inferred-provenance.md) makes high-cost inferred meaning an activation blocker, so it cannot silently back an active capability. **What remains true is the shape of the risk, not the claim:** nothing *obliges* promotion, and this repository's population is near zero because its one maintainer signs decisions promptly. An adopter with a queue and no such habit would see the growth AN-008 predicted, and no mechanism would report it as anything but a count. So the claim is answered for Cliewen and unanswered as a general property.

**Claim 3 — "Retirement flips a status field but removes nothing from the mandatory read path."** → **Answered.** [ADR-034](../decisions/ADR-034-retirement-is-deletion.md) made retirement deletion, and the answer is provable rather than asserted: `retired` is absent from the default status vocabulary, so no committed file may carry it; `checkSupersedes` reports a `supersedes:` naming an artifact that still exists in the corpus, so a claimed retirement whose file survives fails validation. A retired artifact is not in the corpus, therefore `clue context` cannot walk to it and `clue validate` cannot read it. The mandatory read path shrinks by construction.

**Claim 4 — "No edge returns from reality, so a fully green corpus can describe a wrong product."** → **Answered, within a stated boundary.** ADR-035 added `reality: contradicted` on an analysis, linking the failed capability or criterion alongside the carriers that failed to prevent it; `checkReality` validates it and `clue validate --reality-gaps` derives the affected-capability list. The boundary is deliberate and unchanged: one repository-local edge, no telemetry ingestion, no operations loop.

**Claim 5 — "No supersession edge exists: criteria get tombstones, decisions do not, and when a decision is reversed nothing in the graph answers what was downstream of it."** → **Open, and the mechanism that looks like its answer is not.**

The field and the check both ship. `supersedes:` is parsed, and `checkSupersedes` rejects self-supersession, a superseded artifact that still exists, and an ID claimed by two successors. **No artifact in this corpus carries the field.** It has never been exercised, and that is not neglect — ADR-034 binds `supersedes:` to *deletion*, and a superseded decision is not deleted. PDR-001 is a live `verified` decision whose successor PDR-004 replaced it; PDR-003 is a live `verified` decision superseded by PDR-006. Both files remain, correctly, because a superseded decision is history a later reader must be able to read.

So supersession-by-deletion is answered and supersession-of-a-surviving-record is not. Nine of eighty-four decision index rows claim a supersession or amendment, all of it in prose: in the index row, in the successor's body, sometimes in an amendment blockquote at the head of the superseded record. None of it is an edge. Nothing answers "what was downstream of PDR-001" except reading. And the reverse direction is blocked twice over: `clue context` follows outgoing links only, deliberately, so even a machine-visible edge would not answer the question without a reverse walk that PDR-034's narrow-reading obligation argues against.

**Determination.** Pattern C is not closed, and it is not the four-part family AN-008 described. Three of its claims have been answered by P-007 and P-009, one is empirically false in this repository and unproven in general, and what remains is one finding: **a decision that is superseded but not deleted carries no machine-visible edge, and no command answers what depended on it.** That is worth stating in one sentence in place of the paragraph three campaigns have carried forward.

It is **not** closed here and **not** declined here, for a reason that is itself the finding: closing it would need a mechanism, and adding a mechanism inside a simplification campaign needs the argument PDR-026 asks of any addition. The `supersedes:` field is simultaneously the strongest removal candidate on the whole scored surface — shipped, checked, and used zero times — and the only partial answer pattern C has. Those two facts have to be resolved together, and neither is an agent's call. Escalated as **Q-06**.

The shared-memory obligation is worth restating because this determination touches it. PDR-029 holds that reducing what the corpus can remember is not simplification. Removing `supersedes:` would reduce it — not today, when nothing uses the field, but permanently, by removing the only place a future change could record the edge without new machinery.

---

## Escalations

**Q-04 — the tooling surface needs a second test, and this spike proposes one.** *Class: the criterion under-determines.*

*Does the core need it?* returns *no* for eleven of thirteen commands, fourteen of twenty-three checks, six of fourteen types, and eleven of seventeen constraints, and in almost every case the honest answer is that the item should stay. PDR-029 already recorded the mirror failure on carrier prose and repaired it by assigning that surface its own test. The same repair is available here: **does removing it move an obligation from a machine back to a human?** It is checkable, it is not taste, and it reconstructs every decision the project has actually made — `clue` exists at all because rituals that depend on remembering get skipped.

Cost of adopting it: a third test, and a campaign that has to say which test applies where. Cost of not adopting it: either P-013 reports that most of Cliewen's tooling is unnecessary, which is false and would be quoted, or it reports that the criterion passed everything, which is the fourth deferral reached through a new route. What is required: whether to add the second test for the tooling surface as an amendment to PDR-029, and whether the answer changes what P-013 claims to have measured.

**Q-05 — C-003 (tick tasks immediately) governs a transient artifact.** *Class: is this rule still needed, and is it important enough to bind every change?*

The rule obliges an agent to mark a task `[x]` the moment it completes, `[-]` with a reason when infeasible, on a file the digest deletes and that never reaches `main`. It is `enforcement: partial` — `checkSkippedTasks` verifies the shape of a skipped-task line, not the timing. Retaining it binds every full change to a discipline about a file nobody reads after merge. Removing it costs handoff legibility: a batch-ticked checklist does not say which task was in flight when work stopped, and the workspace is what a second agent reads. What is required: whether the handoff value justifies a constraint, or whether the rule belongs in `clue-delta` as procedure with no register entry.

**Q-06 — `supersedes:` is the strongest removal candidate on the scored surface and pattern C's only partial answer.** *Class: two facts about one field that must be resolved together.*

The field ships, the check ships, and no artifact uses it, because ADR-034 binds it to deletion while the supersession that actually happens in this corpus — nine of eighty-four decisions — is of records that survive. On the removal side: an unused field and its check are exactly what a simplification campaign should find. On the retention side: it is the only place a machine-visible supersession edge could be recorded without new machinery, and PDR-029 forbids reducing what the corpus can remember.

Three answers are available and they are not equivalent. **Widen it** — let a surviving superseded record carry the edge, which closes pattern C's residue and adds an obligation to every superseding change. **Leave it** — keep the field for retirement only, and record that pattern C's residue is declined with the cost stated: prose supersession stays unqueryable and "what was downstream of this decision" stays a reading task. **Remove it** — and accept that a future answer needs new machinery. What is required is the choice, and if it is *widen*, the milestone that builds it; P-013 has no milestone that owns a mechanism, which is itself part of the question.

## Answers

Q-04 through Q-06 were answered by Flemming N. Larsen on 2026-08-08 in conversation, in the change that wrote this document, and recorded under [C-011](../constraints/C-011-decision-records-typed.md).

**Q-04 answered: the second test is adopted, and the standing preference is to remove.** [PDR-037](../decisions/PDR-037-tooling-is-judged-by-what-it-holds.md) amends [PDR-029](../decisions/PDR-029-simplification-tests-by-surface.md): on the tooling surface the operative question is whether removing an item moves an obligation from a machine back to a human, a yes is a keep, and a no is a removal candidate the standing preference removes. The human added the direction that shapes the second half — make Cliewen simpler where something does not contribute meaningfully — so the burden falls on retention rather than removal. [ARCH-003](../architecture/core.md)'s periphery clause now carries the qualification, because it stated the single-criterion version as the standing test.

**Q-05 answered: the timing rule is withdrawn and the reason rule stays.** The human's answer was that Cliewen has not needed the obligation to tick a task the moment it completes. The extrapolation this spike proposed alongside it — that the skipped-task reason would go too, and `checkSkippedTasks` with it — was refused: a task marked infeasible should always say why, whenever it is marked. So [C-003](../constraints/C-003-skipped-tasks-carry-reasons.md) narrows to that rule, its check is untouched, and a green `clue validate` asserts exactly what it asserted before. The proposal that this change would give `supersedes:` its first use is therefore void: nothing is retired here, and the field remains unused.

**Q-06 answered: declined and routed.** [PDR-038](../decisions/PDR-038-supersession-residue-declined.md) keeps `supersedes:` as it is — neither widened to a surviving superseded record nor removed as unused — states the cost of the decline, and routes the widening as a named door for P-013's successor. Pattern C stops being carried forward as AN-008's paragraph: a later campaign names that record and the door.

## Rejected approaches

**Scoring the 84 decisions individually against the criterion.** Attempted and abandoned. Every decision states a rule, so the population looks like the right one — but the rules decisions carry reach their readers through carriers, and AN-021 and AN-018 have now registered every statement of every carrier and traced each to its decision. Scoring the decisions as well would score the same rules twice, and the second scoring would have no carrier to point at. What the decision layer needs is not scoring but the screening AN-018 already demonstrated: DLT-06 was found because a *carrier* instructed something a shipped mechanism had replaced. Only two such cases exist in either register so far (DLT-06 and GD-71), both found from the carrier side.

**Counting artifacts, lines, or commands removed.** Refused by PDR-029 before this spike, and this spike is why the refusal matters: the honest result of the scoring is that almost nothing should be removed, and a campaign measured by removals would have to manufacture some.

**Closing pattern C because three of its five claims are answered.** Tempting, and the temptation is exactly what PDR-026 was written against. Three campaigns carried the paragraph forward unexamined; closing it unexamined is the same failure with the opposite sign.

**Declining pattern C's residue here.** Also refused. A decline must state its cost and be a human's decision under PDR-026, and this residue's cost is entangled with a field this same spike found to be unused — so declining it would silently answer Q-06.

## What this analysis does not establish

The verdicts are one reader's judgement. PDR-013 supplies a criterion and no procedure, which is why the failure mode reported here — under-determination on four of five populations — is a finding about the criterion rather than a defect in the populations.

Two results are facts about this repository and not about the method. The inferred-decision population is near zero because one maintainer signs decisions promptly; the `supersedes:` population is zero because this corpus has retired nothing by deletion since ADR-034 shipped. An adopter could differ on both, and pattern C's claim 2 is answered *here* and unproven *generally* for exactly that reason.

Nothing here tests a removal. Every "what removal would cost" cell is a prediction, and the two candidates that reach a human — C-003 and `supersedes:` — are raised precisely because prediction is where an agent should stop.

The check population is scored from source rather than from behaviour. A `dev` build cannot detect skill-to-binary drift, so `checkSkillVersions` was read rather than exercised.

## Consumer

[P-013](../plans/P-013-simplification.md)'s **M-064**, whose scoring and pattern C determination this is; AN-021 carries the same milestone's carrier register. **M-067** consumes the five stale constraint `source:` values as the population for the rule it already owns. **M-066** consumes the determination: pattern C is routed rather than closed, and the campaign cannot report it closed without Q-06's answer.
