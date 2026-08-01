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
