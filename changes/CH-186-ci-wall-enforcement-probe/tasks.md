---
id: CH-186-tasks
type: tasks
status: open
links: [CH-186]
title: Tasks — the CI wall asks the host, and its setup instructions ship with it
---

# CH-186 — Tasks

- [ ] Write `docs/decisions/PDR-060-the-ci-workflow-asks-the-host-continuously.md` with `binds: adopter`, citing `.github/workflows/clue-validation.yml` as its carrier, and naming PDR-059 as the agent-side counterpart this extends to the continuous CI position.
- [ ] Add a probe step to `.github/workflows/clue-validation.yml` that queries the host for the calling branch's protection state (`gh api`/`gh ruleset check` on GitHub) and emits a `::warning` annotation naming what is missing when protection is absent or inadequate, or when the token cannot see the base repository's settings (a fork-originated pull request); never fails the job on a missing or unknown boundary.
- [ ] Add AC-197 to `docs/capabilities/CAP-006-collaborative-handoffs/criteria.md`: an adequately protected branch produces no warning annotation and the `validate` job stays exactly as green as an unrelated pass would leave it.
- [ ] Add AC-198: a missing, inadequate, or unobservable (unknown) boundary produces a plain-language warning annotation naming what is absent or that the workflow could not ask, while the job still succeeds when nothing else fails it.
- [ ] Exercise both paths for real in `clue-validation-smoke.yml` per ADR-038's dogfooding requirement (a shell typo here would reach every adopter's CI having never run), on a disposable ruleset state or via `workflow_dispatch` inputs, and record the observed run as AC-197/AC-198's Integration evidence.
- [ ] Add AC-199 to `docs/capabilities/CAP-001-onboarding/criteria.md`: `clue init` and `clue scaffold` materialize a host-agnostic wall-setup checklist beside the generated `.github/workflows/clue.yml`, distilled from `guide/ci-wall.md`'s "Other forges" contract, pointing to the full guide for the GitHub-specific walkthrough.
- [ ] Add the checklist template under `internal/scaffold/templates/github/` (or the existing workflow-templates directory) and wire it into the same materializer path as `clue.yml`.
- [ ] Write the evidence for AC-199 in `internal/scaffold`'s existing test suite, declaring its purpose and test type.
- [ ] Link PDR-060 and the three new criteria from `docs/capabilities/CAP-006-collaborative-handoffs/README.md` and `docs/capabilities/CAP-001-onboarding/README.md`; update each `design.md` where the new behaviour is described.
- [ ] Add the `[Unreleased]` CHANGELOG entry naming the CI-workflow probe and the shipped checklist file as the adopter-facing surfaces.
- [ ] Assess documentation impact: `guide/ci-wall.md` gains a short pointer to the shipped checklist and a note that the workflow now warns when it detects an unenforced boundary; state in the handoff what changed or why nothing did.
- [ ] Delete the change workspace as the digest.
