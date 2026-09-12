---
id: AN-025
type: analysis
status: active
links: [G-017, P-023, PDR-021, PDR-007, C-012, ADR-044, ADR-062]
title: Cliewen's strictest requirement was unmet here for the project's whole life and nothing noticed
---

# AN-025 — Cliewen's strictest requirement was unmet here for the project's whole life and nothing noticed

## What happened

On 2026-09-12, while enabling Git-coordinated identity allocation, the maintainer and the agent went looking for how to protect the new allocator branch. The repository turned out to have no protection of any kind: `gh api repos/cliewen/cliewen/rulesets` returned an empty list, and `gh api repos/cliewen/cliewen/branches/main/protection` returned *"Branch protection has been disabled on this repository."* Anyone — or any agent holding the maintainer's credentials — could have force-pushed `main` or deleted it outright.

This was not a recent regression. Nothing in the repository has ever checked, and the discovery was incidental to unrelated work.

Two rulesets now exist, each blocking deletion and force-push: `protect-main-history` on the default branch and `protect-id-allocator` on `clue/id-allocator`. Both were the maintainer's decision, taken with the trade named: the default branch still has no required check and no pull-request requirement, because the maintainer wants direct pushes for simple work.

## Why it matters beyond this repository

[PDR-021](../decisions/PDR-021-supported-merge-commit-history.md) is `status: verified`, accepted by the maintainer on 2026-08-02. It says the protected default branch "must allow merge commits and disable squash and rebase-and-merge, while retaining required validation, pull requests, resolved conversations, no-bypass, deletion, and force-push protections." None of that was configured. The decision described a repository that did not exist.

The same rule reaches every adopter. `internal/skills/source/shared/review-boundary.md.tmpl` generates the review-boundary reference in four shipped skills, and it ends: "A forge that cannot enforce the merge-commit boundary is outside the supported full-change adoption path." A "forge" is the Git host: GitHub, GitLab, Bitbucket, or whatever the team uses. So Cliewen tells adopters that what their host can enforce decides whether they may use the method at all.

So the strictest claim Cliewen makes about integration was, in the repository that makes it, unmet and unnoticed for the project's entire history.

## The three-way mismatch

The severity of the rule, the reachability of the instructions, and the enforcement available are inverted relative to each other.

| | Where it lives | Reaches an adopter? |
|---|---|---|
| The requirement | `review-boundary.md.tmpl` → four scaffolded skills | Yes, on every install |
| The workflow that can fail a pull request | `internal/scaffold/templates/github/workflows/clue.yml` | Yes, materialized by `clue init` |
| How to actually configure it — which settings, why the bypass list must be empty, the probe that proves merges are blocked | [`guide/ci-wall.md`](../../guide/ci-wall.md) | **No.** The guide is this repository's website, not a shipped artifact |

The adopter receives the obligation and the workflow, and not the instructions. The most useful page is the least reachable one.

Nothing checks. `clue validate` runs sixteen checks in `internal/corpus/rules.go`, and every one of them reads files in the repository. None looks at the Git host's settings, because [ADR-044](../decisions/ADR-044-judge-reads-state-not-transitions.md) deliberately keeps `clue validate` offline. That decision is right and this analysis does not propose changing it: a validator that called out to a host would give different answers depending on who was logged in, and would stop working on a plane.

But the consequence has never been stated. **Branch protection is not part of the repository.** It is a setting you make where the repository is hosted — in GitHub's, GitLab's, or Bitbucket's own settings screens — and it stays there. Cloning the repository does not bring it along, it appears in no file, and `git log` never shows it changing. So a tool that only reads the repository cannot see it, however hard it looks, and Cliewen has nothing else looking.

An adopter can therefore run `clue init`, get a green `clue validate`, open pull requests, and have nothing enforced at all — rules anyone can bypass, no required check, or no protection whatsoever — with nothing anywhere telling them. Opening a pull request only *shows* whether the checks passed; it is branch protection, at the host, that stops a red one from being merged. `guide/change-loop.md` makes that distinction, and the adopter does not receive it.

The closest thing to a check is `cmd/clue/main_test.go:1085`, which asserts that the words "branch protection" appear in four documents. It verifies that Cliewen says the thing, not that anything is true.

## What could observe it

Only one candidate is real.

**The CI workflow can ask the host.** The workflow `clue init` installs (`clue.yml`) already runs *at* the host, with credentials. It is in a position to ask "is this branch actually protected?" and to say so plainly when the answer is no. This changes nothing about `clue validate`, keeps ADR-044 intact, and puts the question where the answer lives.

The honest objection is that each host answers differently: GitHub, GitLab, and Bitbucket all express protection their own way, so a check written for one does not work on another. That is a real limit, not a reason to skip it. Where Cliewen has no way to ask, it should say it does not know — never report a pass it did not verify.

**`clue validate` cannot do this.** Considered and rejected. It would only ever be able to report that a setting which never lives in the repository is not in the repository, which tells the reader nothing.

**A one-time adoption checklist** is weaker than it sounds. The failure here was not that someone skipped a step at setup; it is that nothing re-checked for the entire life of the project, and settings drift.

## What this does not settle

This analysis records an unmet obligation and the gap that hid it. It does not decide whether the obligation is the right one. Three claims in the shipped text deserve separate examination, and two of them are weaker than the one that failed here:

1. **No-bypass and a required check.** Well justified, and more so now than when the convention was written for humans alone: an agent holds the maintainer's credentials, so a rule the maintainer can bypass is one the agent bypasses silently.
2. **Merge-commit-only, squash unsupported.** Doubtful. `docs/README.md` states that `/docs` is the system of record and Git history is the archive; if the digest has landed, the reviewed chain remaining reachable from `main` is provenance convenience rather than the acceptance boundary. Squash is a widespread convention, and asking a team to change it repository-wide needs a stronger reason than this one.
3. **"Outside the supported full-change adoption path."** Costly. It turns a configuration preference into an eligibility test, and rules out teams whose Git host — or their pricing tier on it — cannot enforce the setting.

Neither PDR-021 nor [C-012](../constraints/C-012-agents-never-merge-own-changes.md) carries a `binds:` field, yet the rule reaches adopters through a shipped carrier. [ADR-062](../decisions/ADR-062-repository-role-is-declared-machine-state.md) exists so that an adopter-binding record declares itself and names its carrier. This one binds adopters without saying so, which is why the boundary check never caught the mismatch.

## Why this belongs to the campaign

[P-023](../plans/P-023-challenge-plans-and-retain-what-work-teaches.md) took on two questions: challenge an assumption before committing to it, and measure the method by what use reveals rather than by what it produces. This is a worked example of both, found by accident rather than by any mechanism the method provides.

For [G-017](../goals/G-017-the-method-is-evidenced-and-simplified.md), it is the first entry the baseline observation should hold, and it is exactly the class AN-024 named as most informative: an obligation that went unmet, with the reason. The reason here is not disagreement or cost. Nobody knew.

For [G-014](../goals/G-014-riskiest-assumption-is-challenged.md), the assumption worth challenging is stated plainly enough to test: *stating a requirement in a shipped carrier is sufficient to make adopters meet it.* Four documents said "branch protection" and the repository had none. That assumption is false here, and the cheapest test of whether it is false elsewhere is to ask the two existing adopters what their default branch actually enforces.
