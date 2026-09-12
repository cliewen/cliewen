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

## Alternatives considered and limits

The review considered whether the immediate need was another artifact type, broader tooling, or additional universal gates. It did not recommend those as the first investment because the documented method already provides homes for the relevant reasoning, and added obligations would themselves require evidence of value. These are deferred possibilities, not consequential alternatives formally rejected on the project's behalf.

A claim that the human boundary is empty, minimal adoption is absent, or reality has no return edge would repeat historical criticisms after their documented remedies. This review instead asks whether the existing mechanisms work well enough in practice. It cannot answer that empirical question from a guide read.

## Suggested next discussion

The reviewer's preferred next investment is a worked example of rejecting a convincing green change, paired with a sustained adopter trial that examines acceptance effort and missed assumptions. Compare that suggestion with the maintainer's ideas before selecting work. No plan, decision, criterion, or shipped skill is changed by preserving this analysis, and no `carried-by` destination is asserted before a finding is actually incorporated into durable product meaning.
