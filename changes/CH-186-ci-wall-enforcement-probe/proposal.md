---
id: CH-186
type: change
status: open
links: [P-023, M-102, AN-025, CAP-006, G-017, PDR-059, PDR-021, ADR-044, ADR-038]
title: The reusable CI workflow asks the host whether it is enforced, and the wall-setup instructions reach an adopter who never visits the website
---

# CH-186 — The CI wall asks the host, and its setup instructions ship with it

[PDR-059](../../docs/decisions/PDR-059-the-agent-asks-whether-the-merge-boundary-is-enforced.md) closed one part of P-023/M-102: before an agent marks a pull request ready, it now asks the host whether the branch it relies on enforces anything. That decision named two parts of M-102 still open, because they are a different mechanism and a different carrier from the agent-side check. This change delivers both.

**The reusable workflow can ask continuously, where the agent only asks once.** `clue-validation.yml` already runs at the host, with credentials, on every pull request and every push to the protected branch. It is in the same position PDR-059 put the agent in, but it runs far more often — every candidate, not just the one an agent happens to mark ready — so it is the cheaper place to catch protection that regresses after the one-time check passed. The same three limits PDR-059 established apply unchanged: it asks, it never enforces (a missing boundary does not fail the check — that would hand the workflow a veto over the maintainer's own configuration, the same objection PDR-059 rejected for the agent); it never reconfigures anything; and where it cannot ask, it reports unknown rather than a pass, and it names the settings its token cannot read — classic branch protection and a ruleset's bypass list — rather than implying an answer for them. The probe runs as an advisory warning annotation on the `validate` job, never a failing step, so a repository that has deliberately chosen no protection (this one, for direct-push simple work) keeps a green required check.

**`guide/ci-wall.md`'s configuration instructions never reach an adopter who does not visit this repository's website.** AN-025 found the requirement and the workflow both shipped, and the instructions for actually configuring the host did not — the least reachable page is the most useful one. The guide is written entirely in GitHub UI vocabulary (menu paths, screen names), which is right for a human reading prose but not for a file that ships beside `clue.yml`. This change materializes a short, host-agnostic checklist — the same contract `guide/ci-wall.md`'s "Other forges" section already states in host-neutral terms — as `.github/cliewen-wall.md`, which `clue init` writes next to the generated `.github/workflows/clue.yml`, and points to the full guide for the GitHub-specific walkthrough. The checklist is the reachable minimum; the guide remains the authoritative detail.

This change closes M-102's observation gap; it does not revisit whether the underlying requirement is right, which remains M-095's subject.
