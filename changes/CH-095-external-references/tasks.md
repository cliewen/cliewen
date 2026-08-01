---
id: CH-095-tasks
type: tasks
status: open
links: [CH-095]
title: Tasks for CH-095
---

# Tasks

- [ ] Record the ADR: the qualified-reference notation, the judge's offline form rule, the separated resolver and its four outcomes, the foreign-evidence pointer, and the rejections that stay closed
- [ ] Revise M-044's exit criterion to the broader boundary this change implements, citing the ADR as its backing decision
- [ ] Add the criteria this change is judged by, with positive and negative evidence for each: the form rule, the resolver's classification, and the foreign-evidence pointer
- [ ] Implement the judge's form rule: a bare forge reference fails, a qualified forge reference and a full URL pass, a foreign corpus ID is qualified by repository, and code fences and headings are never mistaken for references
- [ ] Repair every unqualified reference in this corpus, naming the repository each one actually meant
- [ ] Implement the resolver command: preview by default, explicit write, four outcomes (reachable, redirected, gone, unreachable), a redirect offered as a rewrite, and an unreachable target reported as unknown rather than invalid
- [ ] Implement the foreign-evidence pointer: repository, pinned revision, and identifier, treated as named but locally unproven, never as coverage and never as an imported verdict
- [ ] Move every live carrier together: capability README, criteria and design, canonical and generated skills, scaffold templates, public and contributor guidance, implementation explanations, and `[Unreleased]`
- [ ] Run the complete change verification: build, vet, the full test suite, `clue validate`, and the guide build
