---
id: PDR-060
type: decision
status: verified
links: [P-023, M-102, AN-025, CAP-006, CAP-001, PDR-059, PDR-021, ADR-038, ADR-044, G-017]
title: The validation workflow reports what the host enforces on every run, and never fails because of it
binds: adopter
author: agent
accepted-by: Flemming N. Larsen (2026-09-12, conversation)
---

# PDR-060 — The validation workflow reports the merge boundary on every run

## Context

[PDR-059](PDR-059-the-agent-asks-whether-the-merge-boundary-is-enforced.md) has the agent ask the host once per repository, before the first ready mark. Protection that was right then can be removed later, and a repository whose work never passes through an agent's ready step is never asked at all. The reusable validation workflow ([ADR-038](ADR-038-upstream-validation-workflow.md)) runs at the host on every pull request and every push to the branch it defends, so it is the one continuous observer Cliewen has. `clue validate` cannot be that observer, because it reads repository files and the setting is not in any of them ([ADR-044](ADR-044-judge-reads-state-not-transitions.md)).

## Decision

**Every run of the reusable validation workflow asks GitHub which rules apply to the branch it defends — the pull request's base, or the pushed branch — and reports the answer as an annotation. The answer never changes the job's result.** The boundary it looks for is [PDR-021](PDR-021-supported-merge-commit-history.md)'s: deletion and force pushes blocked, pull requests required, merge commits the only merge method, conversations resolved, and the `validate` check required. When all of that is present the run records a notice; when any part is missing it records a warning naming each missing part in plain terms; when GitHub cannot be asked, or does not answer with a rule list, it records a warning that the state is unknown.

Three limits are part of the decision.

**It reports only what the workflow's token can see, and names what it cannot.** GitHub's rules-for-a-branch endpoint reflects active rulesets and is readable with the `contents: read` permission the caller already grants. Classic branch protection and a ruleset's bypass list need administration access. Asking every adopter to grant an upstream workflow administration rights, to feed a report that changes nothing, is the wrong trade. So the notice says the bypass list went unseen, and the warning says classic branch protection went unseen. A pass the workflow did not observe is never reported, for the same reason PDR-059 gives.

**It never fails the job.** The configuration is the maintainer's decision, and a repository may choose less than the full boundary on purpose. A required check that turned red on a deliberate configuration would be bypassed or removed, and a removed check reports nothing.

**The instructions travel with the report.** `clue init` writes a host-neutral checklist to `.github/cliewen-wall.md`, from `internal/scaffold/templates/github/cliewen-wall.md`, and both the generated caller and every warning point there and to the full guide — the adopter reading the warning is the adopter who has not read the guide. The checklist is not a managed carrier: an adopter may edit it to fit their host, and a migration demanding the file would block every existing adopter's upgrade, while the warning's guide link already reaches them.

[CAP-006](../capabilities/CAP-006-collaborative-handoffs/README.md) holds the report's criteria and [CAP-001](../capabilities/CAP-001-onboarding/README.md) the checklist's.

## Rejected: fail the check when the boundary is missing

It would make the rule enforce itself, and it would override the maintainer's configuration from inside their own required check — the veto PDR-059 already rejected for the agent. It would also turn every run in this repository red, since this repository deliberately allows direct pushes for simple work. A rule that fails everywhere it is not followed is switched off, and then it stops producing the evidence [G-017](../goals/G-017-the-method-is-evidenced-and-simplified.md) asks for.
