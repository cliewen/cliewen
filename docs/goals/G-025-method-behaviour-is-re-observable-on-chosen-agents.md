---
id: G-025
type: goal
status: accepted
links: [G-017, VIS-001]
title: What the method makes an agent do can be re-observed, on chosen agents and models, whenever the maintainer asks
---

# G-025 — What the method makes an agent do can be re-observed, on chosen agents and models, whenever the maintainer asks

**Who wants it:** the maintainer, who has to decide which of the method's obligations to keep, and anyone weighing whether to adopt it (2026-10-05, from [AN-028](../analysis/AN-028-scenario-trials-on-different-agents.md)).

**Why:** [G-017](G-017-the-method-is-evidenced-and-simplified.md) asks for evidence from real use, and today that evidence is hand-run and one-off. A skill changes and nothing says whether agents still route, stop, and ask the way the method intends. Every observation so far was made on the agent that wrote the method, so nothing says whether it holds elsewhere. The vision names an agent orienting from a bounded read and a reviewer trusting a merge because of what the corpus shows; neither is established for an agent other than the one that built the method. G-017 also asks for each obligation to be assessed by the failure it prevents and for the method to be simplified against that evidence, and nothing yet lets an obligation be taken away for a trial and the difference observed.

**Success looks like:**

- A scenario drawn from a goal or use case is defined once, runs on an agent and model the maintainer chooses, and starts only when the maintainer starts it.
- Each run records the conditions it needs to be read against, and it starts from a clean configuration, so the repository's own instructions are what the agent had.
- What a machine can check is checked the same way for every agent, and what only a person can judge is put in front of that person rather than scored.
- Results from two runs of the same scenario, or of the same scenario on two agents, can be compared.
- The same scenarios can be run with one of the method's obligations removed from the agent's copy of the instructions, so what the method needs and what it only adds can be told apart. An obligation is a candidate for removal only when a scenario that provokes the failure it prevents ran without it and the failure did not appear; a scenario that never provokes the failure shows nothing about the obligation.
- The harness is repository tooling. It does not ship to adopters and does not make the method depend on any vendor's agent.
