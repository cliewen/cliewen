---
id: CH-186-tasks
type: tasks
status: open
links: [CH-186]
title: Tasks — the CI wall asks the host, and its setup instructions ship with it
---

# CH-186 — Tasks

- [x] Write `docs/decisions/PDR-060-the-validation-workflow-reports-the-merge-boundary.md` with `binds: adopter`, naming its shipped carriers and PDR-059 as the agent-side counterpart.
- [x] Add AC-197 to `docs/capabilities/CAP-006-collaborative-handoffs/criteria.md`: an enforced boundary is reported as observed, with the unseen bypass list named, and a partial one never is.
- [x] Add AC-198: an unenforced or unanswerable boundary warns in plain terms, or as unknown, and never fails the step.
- [x] Add AC-199 to `docs/capabilities/CAP-001-onboarding/criteria.md`: `clue init` ships a host-neutral checklist beside the caller, which is not a managed carrier.
- [x] Add the probe step to `.github/workflows/clue-validation.yml`, reading GitHub's rules-for-a-branch endpoint with the existing `contents: read` token (AC-197, AC-198).
- [x] Add `internal/scaffold/templates/github/cliewen-wall.md` and point the generated caller at it (AC-199).
- [x] Write AC-197 and AC-198 evidence in `internal/ciguards/mergeboundary_test.go`, running the step's own script with only `curl` stubbed.
- [x] Write AC-199 evidence in `internal/scaffold/wallchecklist_test.go`.
- [ ] Link PDR-060 and the new criteria from the CAP-006 and CAP-001 READMEs, and update each `design.md` where the behaviour is described.
- [ ] Add the `[Unreleased]` CHANGELOG entry naming the workflow report and the checklist file as the adopter-facing surfaces.
- [ ] Update `guide/ci-wall.md` to name the checklist and the workflow's report, using the `humanizer` skill.
- [ ] Mark M-102 done in P-023 with its evidence.
- [ ] Run the Verify Locally block in `CONTRIBUTING.md`, and observe the smoke workflow's real run of the probe on this pull request.
- [ ] Delete the change workspace as the digest.
