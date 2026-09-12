---
id: PDR-059
type: decision
status: inferred
links: [P-023, M-102, AN-025, CAP-006, G-017, PDR-021, PDR-007, ADR-044, ADR-062, C-012]
title: The agent finds out whether the merge boundary is enforced and asks the human, rather than enforcing or assuming
binds: adopter
author: agent
accepted-by: pending human verification (proposed 2026-09-12)
---

# PDR-059 — The agent asks whether the merge boundary is enforced

## Context

Cliewen requires that a full change be accepted through a protected default branch ([PDR-021](PDR-021-supported-merge-commit-history.md)), and the shipped review-boundary reference tells every adopter that a Git host unable to enforce that boundary is outside the supported full-change adoption path. Nothing has ever checked whether an adopter's boundary is enforced.

Nothing could. Branch protection is a setting at the host — GitHub, GitLab, Bitbucket — that no clone carries, that appears in no file, and that `git log` never shows changing. [ADR-044](ADR-044-judge-reads-state-not-transitions.md) keeps `clue validate` reading repository files and off the network, correctly: a validator that called out to a host would answer differently depending on who was logged in, and would stop working offline. So the one tool positioned to judge the corpus is structurally unable to see this.

[AN-025](../analysis/AN-025-unenforced-integration-boundary.md) recorded the consequence in the repository that makes the rule: Cliewen's own default branch had no protection of any kind for the project's entire life, found by accident during unrelated work.

## Decision

**Before marking a pull request ready, the agent finds out whether the branch it is about to rely on actually enforces anything, and stops to ask the human when it does not.** The agent is at the host, with credentials, at the moment the gate starts to matter. It can answer the question no file-reading tool can.

Three limits are part of the decision, not caveats on it.

**It asks; it does not enforce.** A missing boundary does not block the ready mark. Refusing would hand the agent a veto over the maintainer's own configuration choice, and a repository may have decided deliberately — Cliewen's own has, keeping direct pushes available for simple work. The agent asks once, accepts the answer, and states in the readiness handoff that the boundary is unenforced together with the reason given. That record is the point as much as the question is: [G-017](../goals/G-017-the-method-is-evidenced-and-simplified.md) asks for bypassed obligations *with their reasons*, and a decline is exactly that evidence.

**It never changes a setting unasked.** Repository configuration is the human's decision. The agent proposes the exact commands for the host and applies them only on explicit authorization in that exchange. An agent that quietly reconfigured a repository would be doing the thing this finding is about — acting where nobody can see.

**It never reports a pass it did not observe.** Hosts express protection differently, so a check written for one does not work on another. Where Cliewen has no way to ask, the honest report is that it does not know. Silence is not a pass; treating an unanswerable question as a green light is how the original gap stayed invisible.

The rule is carried by `internal/skills/source/shared/review-boundary.md.tmpl`, which generates the review-boundary reference in `clue-delta`, `clue-verify`, `clue-extract`, and `clue-upgrade`, and the scaffolded copies an adopter receives from `clue init`. [CAP-006](../capabilities/CAP-006-collaborative-handoffs/README.md) holds its criteria.

## Scope

This closes the observation gap and nothing else. Whether the underlying requirement is right is a separate question, and AN-025 sets out the case without pre-deciding it: the no-bypass requirement is well justified, while merge-commit-only and the sentence declaring a host outside the supported adoption path are weaker than the obligation that failed here. P-023/M-095 owns that examination.

Two further parts of M-102 remain: a probe in the reusable CI workflow, which can report continuously rather than once per candidate, and getting `guide/ci-wall.md`'s configuration instructions to adopters who never visit this repository's website. The guide is written entirely in GitHub's UI vocabulary, so that is a larger job than it appears.

## Rejected: refuse to mark the pull request ready

Honest to PDR-021 as written, and rejected because it inverts who decides. The maintainer owns the repository's configuration; an agent that withheld a candidate until the configuration matched its expectation would be overriding that. It would also have blocked every change this repository ever made, which is a strong signal that the rule would be bypassed rather than followed — and a bypassed rule teaches nothing, while a recorded decline teaches exactly what G-017 wants to learn.

## Rejected: let `clue validate` report it

`clue validate` could only ever report that a setting which never lives in the repository is not in the repository. That is noise, and it would mean either networking the judge against ADR-044 or asserting a conclusion from no evidence.

## Known gap this does not fix

Neither PDR-021 nor [C-012](../constraints/C-012-agents-never-merge-own-changes.md) carries a `binds:` field, although both reach adopters through this same carrier. [ADR-062](ADR-062-repository-role-is-declared-machine-state.md) exists so an adopter-binding record declares itself and names its carrier, which is why the boundary check never caught the mismatch AN-025 found. Correcting those two records is out of scope here; it is named so the next reader does not rediscover it.
