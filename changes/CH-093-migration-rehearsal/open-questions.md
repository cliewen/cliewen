---
id: CH-093-open-questions
type: open-questions
status: open
links: [CH-093]
title: Open questions for CH-093
---

# Open questions

## OQ-001 — How should the copied CI wall be reconciled before migration?

`clue migrate` can preview the six generated-skill replacements and refuses to write through the `.claude/skills` symlink, but it stops on `.github/workflows/clue.yml` because the file contains copied validation steps rather than a thin upstream caller.

Should the adopter first replace that wall manually with the upstream reusable-workflow caller, or should this rehearsal record a narrower support boundary for repositories whose wall is semantically forked? No automatic migration or target edit is authorized until the boundary is explicit.

## OQ-002 — Which supported environment can establish a target test result?

The target's `gradlew test --no-daemon` run under JDK 26 and Gradle 9.6.1 failed configuration at `:bot-api:typescript:npmTest` because four implicit task dependencies on `npmPack` are rejected; the target documents JDK 17–21 for builds.

Can the maintainer provide a supported JDK 17 or 21 environment and either a compatible Gradle invocation or a target-side resolution for the task dependency before this migration is called operationally viable?

## OQ-003 — Which Kotlin version is authoritative in the adopter instructions?

`AGENTS.md` describes Kotlin 2.2.0 while `gradle/libs.versions.toml` declares Kotlin 2.4.10 at the pinned target head.

Should the instruction be updated to the build source's version, or is the documented version an intentional compatibility constraint? The rehearsal cannot call the target's operating guidance coherent until this is answered.
