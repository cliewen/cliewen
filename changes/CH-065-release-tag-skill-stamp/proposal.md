---
id: CH-065
type: change
status: open
links: []
title: A release runs the judge stamped as its own tag, so a tag can never ship skills that disagree with it
---

# CH-065 — A release runs the judge stamped as its own tag

## What

`.github/workflows/release.yml` gains one gate, between the release-notes extraction and the build: it runs `clue validate --forbid-changes` from source with the tag's version stamped in, exactly as the published binary will be stamped. Because [`corpus.checkSkillVersions`](../../docs/capabilities/CAP-004-ship/design.md) already fails a stamped binary whose marked skills carry a different version, a tag that disagrees with `internal/skills/source/shared/frontmatter.md.tmpl` and its generated copies now fails the release before any artifact exists, naming both versions. `TestSanity_ReleaseRunsTheJudgeStampedAsTheTag` keeps the step, its stamp, and its position ahead of the build from being lost. The workflow's header comment stops asserting the agreement and points at the gate that proves it, CAP-004's design records the gate, and a decision-log row states why the invariant is enforced in the workflow rather than in the corpus rules.

This change is plan-less: [P-007](../../docs/plans/P-007-core-hardening.md)'s milestones are the core-hardening campaign, and none of them owns the release pipeline. It is a reported defect ([issue #63](https://github.com/cliewen/cliewen/issues/63)) in machinery P-002/M-004 built.

## Why

The release has two version inputs that must agree and no step that compares them. The tag drives the binary (`-X main.version=${VERSION}`); the skills carry a hand-maintained frontmatter stamp. [ADR-011](../../docs/decisions/ADR-011-version-stamping.md) makes their disagreement a failure, but the only thing enforcing it is a human remembering to bump the template before tagging — the exact "ritual that depends on a person remembering" the methodology exists to remove.

The cost of forgetting is not paid here. `go test ./...` proves the committed skills match their canonical sources; it never compares either against the tag. So a mismatched tag publishes a self-inconsistent pair, and the failure surfaces in an adopter's repository as `skill version 0.7.0 != clue 0.7.1 (drift — reinstall the skills or clue)` — a drift error with no local change to explain it, manufactured by our release process and diagnosed by the person least able to fix it.

The rule that detects this already exists and is already shipped; only the invocation was missing. Running the judge stamped as the tag closes the gap without inventing a second, shell-resident copy of the version invariant that could drift from the real one.

## Decision boundary

No `clue` behavior, corpus rule, skill, or acceptance criterion changes: the drift rule is untouched and no new check is added to `clue validate`, so this is neither a core change under [C-013](../../docs/constraints/C-013-core-changes-need-decision.md) nor a criterion change. The one decision — the invariant is enforced at the release, not as a new corpus rule — is cheap and local to reverse and is recorded as a log row. Release-cutting stays manual: nothing here bumps a version, regenerates skills, or tags.
