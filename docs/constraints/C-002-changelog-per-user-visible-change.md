---
id: C-002
type: constraint
status: active
links: [ADR-012]
title: Every release-relevant user-visible change adds a changelog entry
source: AGENTS.md (Repository conventions — release notes), ADR-012
enforcement: partial
---

# C-002 — Every release-relevant user-visible change adds a changelog entry

A Cliewen change that affects shipped behavior, a capability, a contract, a command, or a user workflow adds its entry to the `[Unreleased]` section of `CHANGELOG.md` in the digest — written for users, never a PR title or commit subject. A plain editorial change under [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) is not release history and adds no entry; prose that changes normative instructions or user workflow is not plain. The release gates fail a release whose version section is missing, so the rule has a machine at release time; whether *this* change owed an entry, and whether the entry it wrote is written for a user, is judgment.

**Checked by:** the release gates (`.github/scripts/release-gates.sh`, run on the release pull request and again on the merge that tags): a release whose version section is missing from `CHANGELOG.md` cannot be cut.

**Residual:** whether *this* change was release-relevant, and whether its entry is written for a user rather than restated from a commit subject. Both are meaning, and the first is a question about a transition that [ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md) keeps out of the judge. The cost is real: a user-visible change can merge with no entry and nothing notices until the release is being written, when whoever writes it no longer remembers what the change meant to a user — which is exactly the entry this rule exists to get.
