# The skills

Skills are Cliewen's process-knowledge layer. They tell a coding agent what the next right step is while the `clue` binary stays a small deterministic judge.

`clue init` installs the skills under `.agents/skills/` and emits an `AGENTS.md` routing hub. The hub recommends simple or full before loading the corpus: ordinary simple work uses no lifecycle skill, while `clue-upgrade` is the specific simple-work workflow for a human-authorized, coordinated upgrade. A chosen full loop loads the relevant skill and follows a complete workflow without a long prompt.

Assistants disagree about which file they load, so init also emits the bridges that make the hub reachable: the skills are mirrored to `.claude/skills/`, and a `CLAUDE.md` imports `AGENTS.md` because Claude Code reads that name and not this one. The bridge only points — every rule stays in the hub, where every assistant sees it, and the emitted file is yours to extend with genuinely Claude-specific instructions.

## The lifecycle set

| Skill | Use it when | Durable output or hand-off |
|---|---|---|
| `clue-analysis` | Risk or an unknown should be investigated before committing to a design | Verified findings under `/docs/analysis`, then a plan or change |
| `clue-plan` | A goal needs a campaign or an active plan needs a semantic revision | A plan with verifiable milestones, then `clue-delta` |
| `clue-upgrade` | An already-onboarded repository needs to check or take up a newer release | A human-authorized, reviewed coordinated upgrade, or a clear decision to wait |
| `clue-delta` | The user chooses the recommended full loop | A complete branch, digested corpus, and verified pull request |
| `clue-verify` | A full Cliewen pull request is about to become ready | A locally verified and automatically agent-reviewed candidate, then the review hand-off |
| `clue-extract` | An existing repository needs a one-time brownfield conversion | An inferred corpus reviewed through its first change loop |

## Why the skills stay separate

Each skill owns a lifecycle boundary and can be installed independently. Analysis should not load implementation mechanics it does not need, and verification should not improvise the proposal rules after the work is finished. Verification does own the recurring challenge-and-repair hand-off: it delegates review into a clean context where supported, returns findings to the implementing context, and requires a clean pass on the resulting commit.

The files are complete standalone artifacts, but repeated rules are generated from shared canonical sources. This keeps decision routing, simple/full recommendation, repository conventions, and the human acceptance boundary identical across the set without runtime includes.

## Version agreement

Distributed Cliewen skills carry an ownership marker and the same version as the released binary. `clue validate` catches drift among the managed skills and between the skills and a released binary, while ignoring unrelated skills in the shared directory.

This makes process changes reviewable and reproducible: the agent guidance that shaped a branch is a versioned repository artifact, not an invisible service configuration.

## Next

[Make the checks blocking in CI.](./ci-wall)
