---
id: AN-024
type: analysis
status: draft
links: [VIS-001, AN-008]
title: Methodology review — strengthen acceptance judgment and learning from use
---

# AN-024 — Methodology review: strengthen acceptance judgment and learning from use

## Question and consumer

Where should Cliewen focus next to improve as a methodology, given its existing intent graph, evidence conventions, and human acceptance boundary? The maintainer requested an independent perspective on 2026-09-12 and then asked to preserve the findings in the analysis corpus. The consumer is the maintainer's next methodology-prioritization discussion, alongside their own ideas; selected findings can subsequently feed `clue-plan` or `clue-delta`. This record proposes no accepted policy, plan commitment, or change to the core.

## Evidence boundary and approach

This is a qualitative review of the documented methodology at repository commit `1f5fab29d10eea79faa0b6a42926c22b68050520`, conducted in the existing prepared Windows/PowerShell checkout with the installed `clue` binary. It is not a code review, adoption experiment, test of agent compliance, or observation of team outcomes. No performance or statistical claims were measured. The conversation contains the original assessment; this document preserves its reasoning and recommendations as a standalone record.

The review read the guide's explanation of the method, design rationale, change loop, intent, corpus, adoption, operations, prompting, and onboarding examples, together with the draft vision. A focused inspection of the canonical delta and verification skill text checked the treatment of scenario-to-test comparison and Human evidence. When filing this record, the earlier [AN-008 critique](AN-008-methodology-critiques.md) supplied historical context. Its observations describe an earlier revision and are not repeated as current defects.

Evidence links below resolve to the current checkout for navigation; the pinned commit above identifies the text assessed. Findings distinguish documented mechanisms from inferred risks. The recommendations have not been demonstrated to improve outcomes and require human selection before implementation.

The follow-up discussion on 2026-09-12 added the maintainer's proposal for repository-specific guides and automatic preservation of practical learning, together with an expressed interest in pushback on new plans and ideas. The section below attributes that proposal separately from the original review. Its candidate workflow is the reviewer's elaboration, not an accepted requirement or an observed improvement in agent behavior.

The maintainer subsequently named RUP as an inspiration and requested a comparison, then authorized including that comparison in this record. The review consulted the supplied Wikipedia overview and the IBM/Rational sources cited below on 2026-09-12, and compared them with the existing planning and analysis workflows. This is a selective reading of published practices, not a comprehensive RUP conformance assessment or evidence that adopting them improves Cliewen. The proposed adaptations remain reviewer recommendations.

## What already works well

The [design rationale](../../guide/design.md) places preservation of current intent inside the completion of a change. Digesting transient proposals into durable documentation, keeping criterion identities meaning-immutable, and stating the limits of structural validation form a coherent foundation. These are strengths of the documented design, not measured effectiveness claims.

The method already includes a concise acceptance brief, scenario-to-test review, optional use cases, minimal adoption guidance, a simple/full route distinction, and an incident edge for contradictory reality. The improvement question is how effectively these mechanisms support judgment in practice, rather than whether they exist.

## Findings and candidate improvements

### 1. Make human acceptance concrete enough to exercise

**Observed:** The [change loop](../../guide/change-loop.md) puts changed criteria, evidence verdicts, and binding meaning in an acceptance brief. The [design rationale](../../guide/design.md) leaves semantic judgment with the accepting human. The [delta skill source](../../internal/skills/source/skills/clue-delta.md.tmpl) identifies a Human-class criterion's brief line as its proof.

**Inference:** A human merge boundary establishes authority but does not by itself establish sufficient expertise, attention, or evidence. A brief can summarize an assertion without recording the observation needed to support it. Authorization and witnessing evidence are different acts.

**Recommendation:** Develop a worked example in which a human rejects a green candidate because its behavior or evidence is inadequate. Show the actual decision: the behavior change, consequential uncertainty, and what the person must observe or decide. For Human evidence, consider identifying who observed what, under which conditions, and with what result. This is a proposal to clarify evidentiary content, not a claim that existing Human proof is uniformly deficient.

### 2. Challenge the adequacy of criteria as well as their tests

**Observed:** The [verification skill source](../../internal/skills/source/skills/clue-verify.md.tmpl) compares changed scenarios with test setup, action, and assertions. The guide explicitly acknowledges that validation cannot establish whether intent is right.

**Inference:** Criteria, implementation, and tests can share one mistaken assumption. Reviewing in a fresh context helps challenge the implementation conversation, but a reviewer reasoning from the same declared intent can still miss the same omitted user need. Classified positive and negative evidence does not establish that the important counterexamples were selected.

**Recommendation:** For consequential changes, ask: “What plausible implementation would satisfy these criteria and still disappoint or harm the intended user?” For example, a password-reset criterion can establish that an email is sent while omitting account enumeration, expired links, or recovery usability. These are illustrative possibilities, not findings about a Cliewen adopter. Keep the challenge within the existing review and scale further investigation to consequences and uncertainty; no new artifact type is proposed.

### 3. Give post-delivery learning a small, explicit place

**Observed:** The [design rationale](../../guide/design.md) acknowledges that a green corpus can describe the wrong product. Incident findings can link contradictory reality back to failed claims, while the larger production feedback loop remains outside the current scope. The [intent guide](../../guide/intent.md) asks a vision to state success signals and uncertain assumptions.

**Inference:** An incident is only one form of disconfirming evidence. A feature may satisfy its criteria yet fail to produce its intended benefit. A capability such as exporting a report does not establish an outcome such as reducing monthly reconciliation effort.

**Recommendation:** For goals with material uncertainty, consider recording the expected observable outcome, the assumption most likely to fail, and an owner and trigger for revisiting it. Existing goals and analyses may suffice. This would require explicitly considering the standing feedback-loop boundary; it is not permission to reopen it or build telemetry ingestion.

### 4. Test proportionality and remove obligations that do not earn their cost

**Observed:** [Adoption guidance](../../guide/adoption.md) discourages filling unused categories, and the [change loop](../../guide/change-loop.md) routes by accepted-contract impact. The full workflow nevertheless carries numerous precise obligations.

**Inference:** Agent-maintained records still consume human attention at acceptance. Contract impact and operational risk are different dimensions: unchanged-contract work can have severe failure consequences, while a small contract addition can be low risk. The current distinction is useful for documentation obligations but should not be mistaken for a complete basis for verification depth.

**Recommendation:** Exercise routing and acceptance on realistic borderline cases: an undocumented bug, a tiny behavior addition, an uncertain investigation, a compatibility repair, and an urgent production correction. Compare practitioners' reasons and the resulting value to the accepting human. Clarify how failure consequences and uncertainty inform verification without multiplying routes. Assess each proposed rule by the observed failure it addresses and whether it can replace an existing obligation. These are evaluation proposals, not measured deficiencies or a decision to widen simple work.

### 5. Establish the benefit through sustained adopter use

**Observed:** The [operations guide](../../guide/operations.md) explicitly distinguishes foreign-repository trials from adoptions. The [draft vision](../vision.md) identifies scale and long-term maintenance as uncertainties. Earlier adopter evidence exists, including the history discussed in AN-008; this review did not independently reassess those studies.

**Inference:** Dogfooding a methodology that describes its own validator is a favorable environment for explicit rules and traceable criteria. It cannot alone establish sustained value for other product domains. The foreign-soil trials cited in the guide do not establish longitudinal adoption outcomes either.

**Recommendation:** Follow willing adopters through ordinary changes, a mistaken requirement, an incident, and a handover. Observe whether consequential misunderstandings are detected before merge, human acceptance time, a newcomer's ability to recover behavioral rationale, false documentation surviving compliant changes, and repeatedly bypassed rules with their reasons. Define the population and measurement method before drawing quantitative conclusions. Use these as learning measures, not compliance targets; a useful result may justify deleting a rule.

### 6. Use one consistent explanation of intent and delivery

**Observed:** [Methodology](../../guide/methodology.md) illustrates a single thread through goal, plan, change, and capability. [Intent](../../guide/intent.md) explicitly describes intent and delivery as separate threads meeting at the goal.

**Inference:** The two presentations can leave a reader unsure whether delivery artifacts belong to the semantic hierarchy. The two-thread explanation makes the distinction between shipping work and establishing product meaning clearer.

**Recommendation:** Reconcile the explanations with the accepted core before editing carriers. This is a conceptual clarification candidate; this analysis does not decide whether its eventual correction is editorial or changes accepted meaning.

## Maintainer proposal: preserve practical learning in repository guides

**Origin and intent:** In the follow-up conversation, the maintainer proposed `/docs/guides` for non-obvious, repository-specific procedures so humans and agents need not guess how to perform them. The broader intent is for an agent that discovers a smarter way to work to preserve that learning automatically in `/docs`. The maintainer also highlighted the value of pushback on new plans and ideas. This is evidence of interest in a direction, not approval of the particular design below.

**Existing basis:** [ADR-026](../decisions/ADR-026-adopter-types-default-lifecycle.md) permits adopter-defined artifact types and explicitly mentions runbooks. The [cross-cutting design overview](../design/README.md) already describes updating the durable document a reader will need during a change. Neither establishes the proposed practical-learning workflow. This proposal concerns knowledge gained while doing repository work; the post-delivery product-outcome feedback in finding 3 remains a separate question.

### Candidate homes for what is learned

A guide would explain how to accomplish a recurring task whose successful execution is not obvious from the existing documentation or tooling. Possible examples include reproducing a difficult failure, preparing a realistic test environment, adding an integration, and recovering from a failed migration. These are illustrative uses, not instructions to populate a new directory.

| Home | Question it answers |
|---|---|
| Analysis | What did we discover, and how certain are we? |
| Guide | How do I accomplish this task in this repository? |
| Design | How does the system work? |
| Decision | Why did we choose this approach? |
| Script or executable workflow | Can the repeatable procedure be executed reliably? |

The reviewer recommends updating an existing home before creating another. If a clearer error message, fewer setup steps, or a script removes the need to remember a procedure, consider that improvement before adding a guide. A new guide earns its place when useful knowledge remains that a future task will consume. Whether guides need a dedicated artifact type, identity convention, or template remains open.

### Candidate learning loop

1. Notice reusable learning during ordinary work: an unexpected prerequisite, repeated failed approach, human correction, or demonstrably better procedure.
2. Establish its scope and evidence. A procedure that worked once in one prepared environment supports a narrower claim than a generally supported procedure. Preserve prerequisites, applicability, and uncertainty instead of promoting a workaround into a universal rule.
3. Update the appropriate existing document, or propose a guide when a recurring task lacks a suitable explanation. Automate stable mechanical steps when that improves reliability and fits the authorized work.
4. Include the documentation update in the current change without requiring a separate request to document the discovery. State in the handoff what reusable knowledge was preserved; apply the existing route and acceptance rules to its actual meaning.
5. Retrieve that guidance on the next relevant task, compare it with current conditions, and correct or retire it when experience shows it is wrong, incomplete, obsolete, or unnecessarily difficult.

This is learning through shared, versioned repository knowledge, not a claim about model training or private agent memory. Automatic preparation of updates can fit the current change workflow; automatic acceptance of new policy does not follow from the proposal. A discovery that changes a promise or standing obligation still requires the appropriate route and decision. Recording the proposal here changes none of those boundaries.

### Discovery, correction, and evidence-based pushback

Saving information alone does not establish a useful learning loop. A guide needs an explicit task trigger and links from the capability, contributor instructions, or other place where the relevant work begins. Agents should discover applicable guidance without reading every guide for every task. The mechanism for doing so is a design question to resolve before adopting a general documentation obligation.

Prior learning could also make plan review more useful. An agent could identify an assumption in a proposed plan, cite a previous investigation or applicable guide that challenges it, explain the consequence, and offer an alternative. For example: a plan assumes local tests reproduce production conditions, while an earlier investigation established a specific mismatch. That evidence warrants revisiting the assumption; it does not justify reflexively rejecting the plan or treating old guidance as permanently authoritative.

The main inferred risks are accumulation of low-value documents, a one-time workaround becoming policy, duplicated instructions drifting apart, and obsolete guidance reinforcing a past mistake. Scope-qualified evidence, one appropriate home, task-based discovery, and correction on reuse are proposed mitigations whose usefulness needs a trial.

### Candidate first trial and open choices

Consider a small trial of this proposed obligation: when work reveals reusable knowledge that would materially change how a later task is performed, preserve it in the appropriate durable home and make it discoverable. Observe whether a subsequent human or agent finds it, can follow it under its stated prerequisites, and corrects it when conditions differ. Assess saved rediscovery effort and misleading or unused guidance rather than counting documents as progress. No trial is authorized or scheduled by this analysis.

Before implementation, decide which discoveries qualify, how tasks find applicable guides, when an observation is sufficiently supported for procedural guidance, how obsolete knowledge is removed, and whether existing artifact conventions suffice. The choice between human-readable guidance and executable automation should follow the task and evidence rather than require one new workflow mechanism for every discovery.

## RUP comparison: candidate improvements

### Sources and comparison boundary

The maintainer's starting point was the [Rational Unified Process overview](https://en.wikipedia.org/w/index.php?title=Rational_unified_process&oldid=1368248900). The substantive comparisons use [Rational's small-project guidance](https://public.dhe.ibm.com/software/rational/web/whitepapers/2003/tp183.pdf), [IBM's RUP description for System z, chapter 2](https://www.redbooks.ibm.com/redbooks/pdfs/sg247362.pdf), [Rational's process-maturity paper](https://public.dhe.ibm.com/software/rational/web/whitepapers/2003/rup_tp178.pdf), and [IBM's description of scenarios and tool mentors](https://www.ibm.com/docs/en/rational-soft-arch/9.7.0?topic=overview-learning-resources-design-management). These historical sources establish the practices being considered, not their effectiveness in contemporary agent-assisted development. The System z publication supplies a general RUP account within a platform-specific book; no platform or tooling choice is proposed for Cliewen.

### 1. Let risk reduction shape the plan

**RUP practice:** Risk mitigation drives iterative work, and the appropriate amount of process depends on the project's risks. [Rational's small-project guidance](https://public.dhe.ibm.com/software/rational/web/whitepapers/2003/tp183.pdf)

**Already in Cliewen:** The [analysis workflow](../../.agents/skills/clue-analysis/references/analysis-workflow.md) opens with retiring the biggest risk first. The [planning skill](../../internal/skills/source/skills/clue-plan.md.tmpl) defines milestones with verifiable exits and controls changes to plan promises; the change loop reassesses plan health.

**Candidate improvement:** Make the challenge to a proposed approach explicit: which assumption would most seriously undermine it, what is the cheapest credible test, and which result would cause a stop, simplification, or different approach? A first milestone can retire uncertainty even when it produces no feature. For the guides proposal, first testing whether another agent can discover and reuse one guide would challenge the central assumption before building a broader mechanism. This strengthens the existing risk-first practice rather than asserting it is absent.

### 2. Exercise architecture through a difficult scenario

**RUP practice:** Elaboration uses architectural prototypes and significant scenarios to assess whether the proposed architecture can support the system before most construction. [IBM's RUP description, chapter 2](https://www.redbooks.ibm.com/redbooks/pdfs/sg247362.pdf)

**Already in Cliewen:** Architecture and design overviews, decisions, constraints, and analysis findings provide homes for architectural reasoning and evidence. Their existence alone does not demonstrate a proposed architecture under realistic conditions.

**Candidate improvement:** When architectural uncertainty is material, demonstrate a small executable path across the consequential boundaries, selected for the uncertainty it exposes. An integration example might exercise one real request, dependency unavailability, and recovery before implementing the surrounding features. Use the existing analysis and milestone evidence mechanisms; do not require a prototype for every change. This is a proposed application of the current structures, not a finding that adopters never exercise architecture.

### 3. Assess learning at meaningful milestones

**RUP practice:** Iterations span engineering disciplines and typically yield an executable but incomplete system, creating repeated assessment opportunities. [Rational's process-maturity paper](https://public.dhe.ibm.com/software/rational/web/whitepapers/2003/rup_tp178.pdf)

**Already in Cliewen:** Plan-health checks respond when work resumes or evidence challenges the campaign, while the digest preserves durable meaning.

**Candidate improvement:** At a meaningful milestone, ask what was expected, what was observed, and what should change because of the difference. Preserve only useful learning, in the existing home appropriate to it. A result may revise a plan, correct a guide, or support removing an obligation; when there is no useful new information, no new record is needed. This adds a deliberate opportunity to notice learning to the maintainer's proposed capture-and-reuse loop. It is the reviewer's adaptation, not a claim that RUP prescribes automatic repository updates.

### 4. Assess readiness for actual use where the goal requires it

**RUP practice:** Transition includes deployment, user and maintainer preparation, feedback, and adjustments needed for successful use. [IBM's RUP description, Transition phase](https://www.redbooks.ibm.com/redbooks/pdfs/sg247362.pdf)

**Already in Cliewen:** Human merge accepts a repository change; release is adopter-owned. This is a deliberate scope boundary. Relevant capabilities and plan milestones can already state outcomes beyond producing code.

**Candidate improvement:** For a goal that depends on adoption or operation, consider a readiness criterion such as an unfamiliar maintainer performing an upgrade and recovering from its documented failure case. A guide then supports an observable outcome rather than satisfying readiness merely by existing. Preserve the distinction between change acceptance and readiness for use. Any general obligation would require a deliberate methodology decision; this suggestion does not prescribe an adopter release process or reopen the production-feedback boundary.

### 5. Connect practical guidance to the task that needs it

**RUP-related precedent:** Tool mentors provide guidance for applying practices with particular tools, and IBM describes scenarios combining tasks toward a larger outcome. [IBM's scenarios and tool-mentor guidance](https://www.ibm.com/docs/en/rational-soft-arch/9.7.0?topic=overview-learning-resources-design-management)

**Already in Cliewen:** Skills carry process instructions and the corpus supplies bounded task context. The maintainer's guides proposal would add repository-specific practical knowledge, as described above.

**Candidate improvement:** Trial guidance with a task trigger, prerequisites, procedure, expected result, and relevant recovery steps. Connect it to where the task starts and correct it on reuse. This precedent supports task-oriented guidance; automatic preservation and correction through agent work are the proposed Cliewen adaptation. The discovery mechanism and evidentiary boundaries remain the open choices already recorded in the guides section.

### What this comparison does not recommend importing

The reviewer recommends selective adoption of practices that address demonstrated needs. RUP's own [small-project guidance](https://public.dhe.ibm.com/software/rational/web/whitepapers/2003/tp183.pdf) supports tailoring process to risk. Mandatory phase names, a larger role taxonomy, document inventories, and mandatory use cases are not proposed here. Cliewen's optional use cases, proportional routes, and existing artifact homes should be considered before adding structure. These are recommendations for the prioritization discussion, not formal rejection of future alternatives.

## Alternatives considered and limits

The review considered whether the immediate need was another artifact type, broader tooling, or additional universal gates. It did not recommend those as the first investment because the documented method already provides homes for the relevant reasoning, and added obligations would themselves require evidence of value. These are deferred possibilities, not consequential alternatives formally rejected on the project's behalf.

A claim that the human boundary is empty, minimal adoption is absent, or reality has no return edge would repeat historical criticisms after their documented remedies. This review instead asks whether the existing mechanisms work well enough in practice. It cannot answer that empirical question from a guide read.

## Suggested next discussion

The original review emphasizes acceptance judgment and evidence from sustained use. The maintainer's proposal adds reusable practical learning, and the RUP comparison suggests strengthening risk-driven planning and milestone learning assessments. Together they suggest a candidate sequence for discussion:

1. Challenge a proposed plan by identifying its consequential assumption, a credible test, and a result that would change the approach. Include a worked example where a convincing proposal or green candidate is redirected.
2. Trial one guide's discovery and reuse, then assess what the attempt taught before expanding the mechanism. Preserve useful corrections through the ordinary change workflow.
3. Follow subsequent work, preferably with willing adopters, to assess acceptance effort, missed assumptions, and whether preserved knowledge is actually used. Apply executable architecture demonstrations and readiness-for-use criteria where the trial or project exposes those specific risks.

This sequence is a reviewer recommendation for human prioritization, not an accepted plan or authorization to run the trials. No plan, decision, criterion, or shipped skill is changed by preserving this analysis, and no `carried-by` destination is asserted before a finding is actually incorporated into durable product meaning.
