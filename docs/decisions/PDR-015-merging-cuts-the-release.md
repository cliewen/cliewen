---
id: PDR-015
type: decision
status: inferred
links: [CAP-004, ADR-011, ADR-012, C-012]
title: Merging the release pull request cuts the release; the tag is derived, not performed
author: agent
accepted-by: []
---

# PDR-015 — Merging cuts the release

## Context and problem statement

Cutting a release took two human acts. First a reviewed pull request raised the skills' version stamp and closed the `CHANGELOG.md` section; then, separately, somebody created and pushed the matching tag, which is what actually publishes.

The second act carries no judgement. By the time the first one has merged, the version is already decided and written down in the repository, in the one file the release's own drift gate reads. Tagging cannot be done differently, only forgotten — and it was: 0.8.0's bump sat on `main` while `https://cliewen.dev/install.sh` kept installing 0.7.0, with nothing anywhere reporting a problem, because from every workflow's point of view nothing had happened.

A step that is mechanical, unskippable and silent when omitted is a step a machine should perform.

## Decision outcome

**Merging a pull request that raises the version stamp cuts the release. `tag-on-merge.yml` reads the version the merge brought in, tags the merge commit, and starts the release workflow.**

- **The merge is the authorization.** [C-012](../constraints/C-012-agents-never-merge-own-changes.md) puts a human at the merge gate, and this decision does not move it. What changes is that the human's approval now reaches all the way to the published release instead of stopping one step short. Nothing publishes that a human did not merge.
- **One version source.** The tag is read from `internal/skills/source/shared/frontmatter.md.tmpl`, the repository's single bump site and already the input [ADR-011](ADR-011-version-stamping.md)'s drift gate judges the tag against. Deriving it from anywhere else would create a second source that could disagree with the gate.
- **Quiet on every ordinary merge.** The workflow runs on every push to `main` and exits when the version's tag already exists. Merges that do not raise the version cost one green no-op run.
- **Notes are required before the tag, not after.** A raised stamp with no `## [X.Y.Z]` section fails here. The release workflow would fail too ([ADR-012](ADR-012-release-notes-from-changelog.md)), but a version number cannot be un-spent, so the cheaper failure is the one that happens before the tag exists.
- **The release is started explicitly.** A tag pushed with `GITHUB_TOKEN` does not trigger workflows — GitHub suppresses it to prevent loops — so the workflow dispatches `release.yml` against the new tag, `workflow_dispatch` being the documented exception to that suppression. That path already existed and is already guarded to tag refs, so the run is identical to one a manually pushed tag would have produced, and no personal access token is introduced.

**Carrier:** `.github/workflows/tag-on-merge.yml` and `TestSanity_TagOnMergeDerivesTheVersionAndStartsTheRelease`, which holds the version source, the dispatch, the permission it needs, the already-tagged no-op and the notes requirement.

### Rejected: keep tagging manual

The status quo, and defensible on the grounds that publishing should feel deliberate. It is rejected because the deliberation had already happened at the merge, and because its failure mode is silence: a forgotten tag produces no error, no red run, and no signal — only users installing a version that is not the current one, which is how this was found.

### Rejected: a "cut a release" button that also writes the bump

A dispatched workflow could raise the stamp, regenerate the skills, close the changelog section and open the pull request. It is genuinely convenient and remains available later. It is not this decision because it needs a stored credential — a pull request opened with `GITHUB_TOKEN` does not run CI, so the release change would arrive at the merge gate with no checks — and because the bump is the one part of the ritual with judgement in it: which number, and whether the notes read well.

### Rejected: tagging from the release pull request's own workflow run

Tagging inside the pull request's checks would fire before merge, tagging a commit that might never reach `main`.

## Consequences

Raising the version stamp in a merged pull request now publishes a release. That is the intent, and it makes the stamp a more dangerous line to edit casually than it was: an unintended bump reaching `main` spends a version number. The bump lives in a reviewed change for that reason, and the notes requirement means an accidental one usually fails before the tag exists rather than after.
