---
id: CH-185
type: change
status: open
links: [P-023, M-102, AN-025, CAP-006, G-017, PDR-021, ADR-044]
title: The agent finds out whether the merge boundary is enforced, and asks the human when it is not
---

# CH-185 — The agent finds out whether the merge boundary is enforced, and asks the human when it is not

[AN-025](../../docs/analysis/AN-025-unenforced-integration-boundary.md) found that this repository had no branch protection of any kind for the project's entire life, while four shipped skill references tell every adopter that a Git host unable to enforce the merge-commit boundary is outside the supported adoption path. The requirement shipped, the instructions did not, and nothing checked. Nothing could have: branch protection is a setting at the host — GitHub, GitLab, Bitbucket — that no clone carries, that appears in no file, and that `git log` never shows changing. `clue validate` reads files, so it can never see it.

The agent can. It is already at the host, with credentials, at the moment the gate starts to matter. [P-023](../../docs/plans/P-023-challenge-plans-and-retain-what-work-teaches.md)/M-102 owns closing that observation gap, and this change delivers its first and cheapest part: before an agent marks a pull request ready, it finds out whether the branch it is about to rely on enforces anything, and stops to ask the human when it does not.

Three things shape what this rule may do, and all three came from working the problem on this repository today.

**It asks; it does not enforce.** A rule that refused to mark the pull request ready would hand the agent a veto over the maintainer's own configuration choice, and it would have blocked every change this repository ever made. The maintainer here has deliberately chosen no pull-request requirement so that simple work can be pushed directly. A decline is a legitimate answer, so the agent records it and proceeds. That record is also what [G-017](../../docs/goals/G-017-the-method-is-evidenced-and-simplified.md)'s baseline observation asks for: a bypassed obligation together with the reason it was bypassed.

**It never configures anything unasked.** Repository settings are the human's decision, and an agent that quietly changed them would be doing the thing this whole finding is about — acting where nobody can see. It proposes the exact commands and applies them only on an explicit yes in that exchange.

**It never reports a pass it did not observe.** Each host expresses protection differently, so a check written for one does not work on another. Where Cliewen cannot ask, the honest answer is that it does not know. Silence is not a pass — treating an unanswerable question as a green light is how the original gap stayed invisible.

The rule lands in `internal/skills/source/shared/review-boundary.md.tmpl`, so it binds every adopter. [PDR-059](../../docs/decisions/PDR-059-the-agent-asks-whether-the-merge-boundary-is-enforced.md) records it with `binds: adopter` and names that carrier, and [CAP-006](../../docs/capabilities/CAP-006-collaborative-handoffs/README.md) gains four criteria — protection adequate, protection missing, human declines, and the host cannot be asked.

This change closes the observation gap only. Whether the underlying requirement is the right one — particularly merge-commit-only and the sentence declaring a forge outside the supported adoption path — is M-095's subject, and AN-025 sets out the case there without pre-deciding it. The remaining parts of M-102, a probe in the reusable CI workflow and getting `guide/ci-wall.md`'s configuration instructions to adopters who never visit this repository's website, are separate changes.
