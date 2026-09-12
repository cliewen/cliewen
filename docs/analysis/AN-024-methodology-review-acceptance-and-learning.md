---
id: AN-024
type: analysis
status: active
links: [VIS-001, AN-008]
carried-by: [G-014, G-015, G-016, G-017, P-023]
title: Methodology review — strengthen acceptance judgment and learning from use
---

# AN-024 — Methodology review: strengthen acceptance judgment and learning from use

## Purpose and evidence boundary

This analysis supports the maintainer's next methodology-prioritization discussion. It combines the independent review requested on 2026-09-12, the maintainer's proposal for repository-specific guides and automatic preservation of practical learning, and a selective RUP comparison requested by the maintainer. The four priorities and candidate plans below are reviewer recommendations for consideration, not accepted policy, scheduled trials, or authorization to change the methodology.

The review examined the guide, draft vision, and relevant planning, analysis, delta, and verification instructions at repository revision `1f5fab29d10eea79faa0b6a42926c22b68050520`. It used the prepared Windows/PowerShell checkout with the installed `clue` binary. No code review, adoption experiment, or statistical measurement was performed. [AN-008](AN-008-methodology-critiques.md) supplies historical context; its earlier defects are not presumed to remain. Repository links navigate current files; the revision identifies the assessed text.

The RUP comparison began with the supplied [Wikipedia overview](https://en.wikipedia.org/w/index.php?title=Rational_unified_process&oldid=1368248900) and consulted the IBM/Rational sources cited below on 2026-09-12. They establish historical practices, not demonstrated effectiveness in agent-assisted development. Proposed adaptations and inferred risks still need testing.

Cliewen's [existing foundation](../../guide/design.md) is strong: current intent is preserved within the change, criterion meanings remain stable, and structural validation has explicit limits. Acceptance briefs, scenario review, optional use cases, proportional routes, and contradictory-reality records already exist. The question is how well they support judgment and learning in practice.

## Prioritized recommendations

### 1. Challenge ideas and plans before commitment

**Basis and gap:** The [analysis workflow](../../.agents/skills/clue-analysis/references/analysis-workflow.md) already prioritizes risk, and the [planning skill](../../internal/skills/source/skills/clue-plan.md.tmpl) requires verifiable milestones. However, criteria, implementation, and tests can share one mistaken assumption. Scenario-to-test alignment does not establish that the scenario represents the right outcome.

RUP offers two relevant practices: iteration and process effort driven by risk, and executable architectural demonstrations before substantial construction. See [Rational's small-project guidance](https://public.dhe.ibm.com/software/rational/web/whitepapers/2003/tp183.pdf) and [IBM's RUP description, chapter 2](https://www.redbooks.ibm.com/redbooks/pdfs/sg247362.pdf).

**Recommendation:** Before a consequential commitment, identify the assumption most likely to undermine the plan, a credible alternative, the cheapest useful test, and the result that would cause a stop or revision. Ask what implementation could satisfy the criteria while still failing the intended user. Cite applicable repository experience when challenging assumptions; old guidance is evidence to reassess, not permanent authority.

Where architecture is materially uncertain, exercise a small path across the consequential boundaries, including failure and recovery. A milestone may produce decision-changing evidence without delivering a feature. Scale investigation to consequences and uncertainty rather than requiring prototypes universally.

**Why first and what success means:** Avoiding misguided implementation has immediate value. Demonstrate both a weak proposal redirected by evidence and a sound proposal proceeding without unnecessary investigation.

### 2. Preserve and reuse practical learning

**Basis and gap:** The maintainer proposed `/docs/guides` so humans and agents need not rediscover non-obvious repository procedures. [ADR-026](../decisions/ADR-026-adopter-types-default-lifecycle.md) already permits extensions such as runbooks, and the [design overview](../design/README.md) describes maintaining useful documentation during changes. Neither establishes the proposed capture-and-reuse behavior.

RUP-related [tool mentors](https://www.ibm.com/docs/en/rational-soft-arch/9.7.0?topic=overview-learning-resources-design-management) provide task-specific guidance. Its [iterative approach](https://public.dhe.ibm.com/software/rational/web/whitepapers/2003/rup_tp178.pdf) also provides assessment opportunities. Automatic repository updates are the proposed Cliewen adaptation.

**Recommendation:** Notice a reusable discovery, establish its evidence and applicability, update its appropriate home within the current change, and make it discoverable on the next relevant task. Correct or retire it when reuse exposes a problem. At meaningful milestones, compare expected and observed results; preserve useful learning without requiring a record when nothing useful emerged.

| Home | Purpose |
|---|---|
| Analysis | Findings and uncertainty |
| Guide | How to accomplish a recurring repository task |
| Design | How the system works |
| Decision | Why a consequential approach was chosen |
| Script or executable workflow | Reliable execution of repeatable steps |

Prefer correcting an existing home. Consider removing confusing steps or automating them before documenting a workaround. A guide should identify its task trigger, prerequisites, procedure, expected result, and relevant recovery. Link it from where work begins so agents need not read every guide. One success in a prepared environment supports only a scoped claim.

**Why second and what success means:** Reusable knowledge compounds the benefit of better planning. A fresh agent must find and successfully use earlier guidance, then correct it when conditions change. Writing documents alone does not pass this trial. Automatic preparation remains within ordinary change authority; new policy still follows the existing acceptance boundary.

### 3. Make acceptance an informed decision about usable behavior

**Basis and gap:** The [change loop](../../guide/change-loop.md) provides an acceptance brief; the [verification skill](../../internal/skills/source/skills/clue-verify.md.tmpl) compares scenarios with test assertions. The [delta skill](../../internal/skills/source/skills/clue-delta.md.tmpl) treats a Human-class criterion's brief line as its proof. These structures establish responsibilities, but cannot alone establish adequate expertise, attention, or an actual observation. Authorization and witnessing evidence are different acts.

RUP's Transition activities include deployment, preparation, and user feedback. [IBM's description](https://www.redbooks.ibm.com/redbooks/pdfs/sg247362.pdf) supplies a useful readiness perspective; it does not justify prescribing an adopter's release process.

**Recommendation:** Make the brief state changed behavior, supporting observations, and consequential uncertainty requiring judgment. For Human evidence, identify who observed what, under which conditions, and the result. When a goal depends on adoption or operation, demonstrate the relevant task: an unfamiliar maintainer performing an upgrade and recovering from a documented failure, for example. A guide's existence does not establish readiness.

**Why third and what success means:** Better planning still needs an informed final decision. A human should be able to explain acceptance or rejection from the brief and evidence, including rejecting an inadequate green candidate without repeating the entire code review.

### 4. Evaluate the method's value and simplify it

**Basis and gap:** [Adoption guidance](../../guide/adoption.md) already discourages unnecessary artifacts. The [operations guide](../../guide/operations.md) distinguishes foreign-soil trials from adoptions, and the [draft vision](../vision.md) acknowledges uncertainty about scale and sustained maintenance. Earlier adopter studies were not independently reassessed here. Dogfooding alone cannot establish general effectiveness.

**Recommendation:** Follow real work through planning, acceptance, and subsequent use, preferably with willing adopters. Observe consequential misunderstandings, acceptance effort, rediscovery, misleading documentation, and bypassed obligations with their reasons. Include handover, mistaken requirements, incidents, and borderline routing cases. Define populations, conditions, and sampling before quantitative conclusions; measure learning rather than document counts or green checks.

Contract impact and operational risk differ: unchanged-contract work can have severe consequences. Evaluate verification depth accordingly, and assess each obligation by the failure it prevents and whether it can replace existing ceremony.

Also distinguish practical task learning from product outcomes. A working report export does not establish reduced reconciliation effort. For uncertain goals, consider an observable benefit, its key assumption, and an owner and trigger for revisiting it through existing goals and analyses.

**Why fourth and what success means:** Evidence should justify keeping, revising, or removing particular practices. Start baseline observations alongside priority 1; evaluation runs throughout delivery rather than waiting until the other priorities finish.

## Potential plans and sequencing

These are candidate campaign boundaries, with no plan or milestone identities allocated.

| Candidate plan | Priority | Bounded result |
|---|---|---|
| Challenge plans and retain what work teaches | 1–2 | Demonstrate constructive pushback and one complete discovery, reuse, and correction loop before expanding the mechanism. |
| Make human acceptance concrete | 3 | Exercise acceptance and rejection examples, including Human observations and practical readiness where relevant. |
| Establish value through sustained use | 4 | Begin baseline observation with the first campaign, continue across subsequent work, and publish evidence-supported simplifications. |

The first candidate is the recommended initial delivery. Its guides proposal provides a useful test subject for its own plan challenge. Subsequent scope should follow the trials' findings, with human selection before implementing changes to accepted meaning.

## Open choices and deferrals

Resolve guide eligibility, evidence sufficiency, discovery, retirement, and any need for a new type or template before adopting a general learning obligation. Guard against low-value accumulation, duplicated instructions, and workarounds becoming policy. Preserve proportional routes and the distinction between structural validation and semantic judgment.

One accompanying clarification remains: [methodology](../../guide/methodology.md) presents a single thread through delivery artifacts, while [intent](../../guide/intent.md) separates intent and delivery. Reconcile these explanations with the accepted core; this analysis does not pre-decide whether that correction changes meaning.

Defer broad tooling, mandatory RUP phases or artifact inventories, and general production-feedback infrastructure until evidence warrants them. Product-outcome learning must respect the existing feedback boundary, and release remains adopter-owned. No alternatives are formally rejected here. Record `carried-by` only when findings actually reach durable product meaning.
