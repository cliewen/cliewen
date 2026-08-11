---
id: PDR-041
type: decision
status: inferred
links: [G-007, C-002, ADR-012, ADR-013]
title: A change owes a release note when it changes what an adopter receives
author: agent
accepted-by: []
---

# PDR-041 — Release-note scope is the shipped surface

## Context and problem statement

`CHANGELOG.md` is published verbatim as the GitHub release body ([ADR-012](ADR-012-release-notes-from-changelog.md)), so its audience is whoever installs `clue` or adopts the skills. The rule saying which changes owe an entry listed the triggers as shipped behavior, a capability, a contract, a command, or a user workflow. Every term in that list is also true of something this repository holds locally and no adopter ever sees: it has its own capabilities, its own contributor contract, its own commands, and its own workflows.

[ADR-013](ADR-013-ships-generic-vs-repo-local.md) already draws the line the rule needs — what ships to adopters is generic, and `AGENTS.md` is the repo-local layer — but the scope rule did not name it. Applying the rule therefore meant reaching past it to that decision and re-deriving the boundary, and two readings of the same change could reasonably reach opposite answers. The failure is asymmetric: an entry wrongly written describes something the reader does not have, and it is published before anyone notices, while an entry wrongly omitted is caught when the release is written.

## Decision outcome

**A change owes a release note when it changes what an adopter receives; every term in the scope list names that shipped surface.**

- **The test.** Does the change alter the behaviour of `clue`, the text of a generated skill, or an artifact `clue init` or `clue scaffold` materializes into an adopter repository? If none of the three, no entry is owed.
- **This repository's own layer is outside the scope**, including its corpus, its contributor guidance, its local conventions, and its CI — even when what changed is genuinely a capability, a contract, a command, or a workflow of this repository. That a change is not plain, and carries full Cliewen bookkeeping, says nothing about whether it owes an entry.
- **A change spanning both writes the entry for the adopter-visible part alone**, in the reader's terms, and says nothing about the repo-local half.
- **The residual stays judgment.** The test names the surface; whether a given edit changes what that surface *means* to a user is still a reading, and the release gates continue to check only that a version section exists.

**Carrier:** [C-002](../constraints/C-002-changelog-per-user-visible-change.md) states the rule and the test; the release-notes convention in `AGENTS.md` and the digest sentence in `CONTRIBUTING.md` state the same rule for the agent and the contributor. The generated skills are deliberately not carriers: they say only that a repository's own conventions may require a user-facing changelog entry, and naming this repository's scope test in shipped text is what [ADR-013](ADR-013-ships-generic-vs-repo-local.md) forbids.

### Rejected: replace the list with "is it user-visible?"

That question is the ambiguity, not a resolution of it — "user" is exactly the word that reads as either an adopter or a contributor, and the list at least names candidate triggers.

### Rejected: derive the obligation mechanically from changed paths

A path list asserting "these files ship" goes stale silently as the repository grows, and it answers the wrong question: a comment in shipped code changes a shipped file without changing anything an adopter receives, while a corpus edit can change a skill's generated text. The judge reads state rather than transitions ([ADR-044](ADR-044-judge-reads-state-not-transitions.md)), and this obligation is about what a change did.

### Rejected: an entry for every Cliewen change

It would remove the judgment and the ambiguity with it, at the cost of release bodies that read as an internal commit log — the outcome ADR-012 exists to prevent.
