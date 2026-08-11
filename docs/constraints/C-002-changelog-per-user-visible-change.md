---
id: C-002
type: constraint
status: active
links: [ADR-012, PDR-041]
title: Every release-relevant user-visible change adds a changelog entry
source: AGENTS.md (Repository conventions — release notes), ADR-012, PDR-041
enforcement: partial
---

# C-002 — Every release-relevant user-visible change adds a changelog entry

A Cliewen change that affects shipped behavior, a capability, a contract, a command, or a user workflow adds its entry to the `[Unreleased]` section of `CHANGELOG.md` in the digest — written for users, never a PR title or commit subject. Every term in that list names the surface an adopter receives, never this repository's own equivalent ([PDR-041](../decisions/PDR-041-release-note-scope-is-the-shipped-surface.md)): the test is whether the change alters the behaviour of `clue`, the text of a generated skill, or an artifact `clue init` or `clue scaffold` materializes into an adopter repository. A change confined to this repository's own corpus, contributor guidance, and local conventions owes no entry, even when what changed is genuinely a capability, a contract, a command, or a workflow here; a change spanning both writes the entry for the adopter-visible part alone. That list of categories is a shortcut, and the test is the rule: where the two disagree the test governs. `.github/` holds both kinds — `ci.yml` is this repository's own, while `pull_request_template.md` is held byte-identical to what `clue init` writes and `clue-validation.yml` is the reusable workflow an adopter's caller references, so a change to either of those owes an entry. A plain editorial change under [PDR-011](../decisions/PDR-011-plain-changes-bypass-cliewen.md) is not release history and adds no entry; prose that changes normative instructions or user workflow is not plain. The release gates fail a release whose version section is missing, so the rule has a machine at release time; whether *this* change owed an entry, and whether the entry it wrote is written for a user, is judgment.

**Checked by:** the release gates (`.github/scripts/release-gates.sh`, run on the release pull request and again on the merge that tags): a release whose version section is missing from `CHANGELOG.md` cannot be cut.

**Residual:** whether *this* change was release-relevant, and whether its entry is written for a user rather than restated from a commit subject. Both are meaning, and the first is a question about a transition that [ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md) keeps out of the judge. The cost is real: a user-visible change can merge with no entry and nothing notices until the release is being written, when whoever writes it no longer remembers what the change meant to a user — which is exactly the entry this rule exists to get.
