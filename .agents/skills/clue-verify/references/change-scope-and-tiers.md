## Change scope and tiers

Classify scope before using the Cliewen loop. Three rules decide the tier, by how deeply the change reaches into meaning; take the first rule that matches.

1. **Plain — nothing about meaning changes.** No product behavior, intent, evidence, decision, plan, policy, or methodology changes. Protected product, corpus, test, configuration, build/release, governance/security, agent-rule, skill, and lint surfaces are never plain; neither are commands, contracts, workflows, or normative instructions. Plain work stays outside this skill: branch from `main`, run relevant checks, open a ready PR for human merge, and use no CH identity or Cliewen bookkeeping.
2. **Light — meaning is touched but not changed.** No decision, acceptance-criterion or capability meaning change, semantic plan mutation, or methodology carrier. Typical: protected-surface clarity, dependency bumps, pure refactors, and CI plumbing. Use a Cliewen branch and ready PR whose description names the plan item or plan-less scope, but no transient workspace.
3. **Full — everything else.** Product behavior changes are full even when an existing criterion already states the behavior. Use the whole loop with `/changes/CH-xxx-slug/`.

Two guards hold above the rules. **Uncertainty escalates:** when the tier is unclear, take the higher one. **Discovery escalates immediately:** the moment a decision, an open question, a meaning change, or a methodology-carrier edit appears during work, move to the full loop before continuing.
