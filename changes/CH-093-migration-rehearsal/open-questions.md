---
id: CH-093-open-questions
type: open-questions
status: open
links: [CH-093]
title: Open questions for CH-093
---

# Open questions

## OQ-001 — How should the copied CI wall be reconciled before migration?

`clue migrate` can preview the six generated-carrier replacements and refuses to write through the `.claude/skills` symlink, but it stops on `.github/workflows/clue.yml` because the file contains copied validation steps rather than a thin upstream caller.

Should the adopter first replace that wall manually with the upstream reusable-workflow caller, or should this rehearsal record a narrower support boundary for repositories whose wall is semantically forked? No automatic migration or target edit is authorized until the boundary is explicit.

## OQ-002 — Which supported environment can establish a target test result?

The target's `gradlew test --no-daemon` run under JDK 26 and Gradle 9.6.1 failed configuration at `:bot-api:typescript:npmTest` because four implicit task dependencies on `npmPack` are rejected. This is a task-configuration defect rather than a JDK rejection, so a supported JDK alone will not produce a test result.

The supported JDK is also not stated consistently by the target: `CONTRIBUTING.md` and `DEVELOPMENT.md` require JDK 17+, while `docs/architecture/report/architectural-health-report.md` narrows the build range to JDK 17–21.

Can the maintainer resolve the `npmPack` task dependency, and name which JDK statement is authoritative, before this migration is called operationally viable?

## OQ-003 — How does the rehearsal's evidence survive the digest?

The rehearsal report and these questions live in `/changes/CH-093-migration-rehearsal/`, which the digest deletes because `main` never contains `/changes/`. PDR-020 lands a durable report in `/docs/analysis` only in the mutate phase, and this change does not authorize mutation. As it stands, the unmet M-042 conjunct, the three named limits of this target, the narrower support boundary, and OQ-001 and OQ-002 would leave no trace on `main` beyond one plan evidence cell.

Should the digest promote this report to a durable `/docs/analysis` artifact — it already carries analysis frontmatter and sits naturally beside AN-013 and AN-014 — and carry OQ-001 and OQ-002 forward as named P-009 doors or recorded decisions? The workspace must not be deleted until that landing is decided.

A yes answer widens this change beyond what it proposed. `proposal.md` states that this repository's durable `/docs` corpus remains unchanged and that the rehearsal report is transient evidence, and PDR-020 reserves the `/docs/analysis` landing for the mutate phase. So the answer must also direct the amendment: `proposal.md` records the added durable artifact in the same change, and the digest states that the landing happened under human direction rather than under PDR-020's mutate phase, which this change still does not authorize.

## OQ-004 — Is the narrower support boundary the honest one?

One M-042 exit conjunct is unmet: the proposed thin-caller wall running in the target CI shape. M-042's own exit criterion admits this ending — the phase closes when the rehearsal demonstrates the migration is truthful and operationally viable "or records a narrower honest support boundary" — so taking that branch is compliance with the milestone as written. P-009's mutation rules let the digest set the milestone's status and evidence fields, which is all the narrower branch needs; no plan revision or decision record is required to record it.

What is not a mechanical question is whether this boundary is the honest one. The rehearsal was retrospective against an already-adopted target, so it inventoried stable-ID preservation and the pre-mutation checkpoint without exercising either, and the target's own test half is a configuration failure rather than a suite verdict. Those are limits of this target, not unmet conjuncts, but they bear directly on whether one rehearsal against one converted repository is enough to close the migration-readiness phase.

Does the human accept this boundary as satisfying M-042 and closing the phase, or does the phase need a second rehearsal against an unconverted target that carries stable source IDs? Until this is answered, M-042 stays `todo`; on a yes, the digest sets it `done` with the boundary and OQ-001 as its evidence.
