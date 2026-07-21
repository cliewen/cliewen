# The skills

Skills are Cliewen's process-knowledge layer. They tell a coding agent what the next right step is while the `clue` binary stays a small deterministic judge.

`clue init` installs the skills under `.agents/skills/` and emits an `AGENTS.md` routing hub. The routing file classifies plain work before loading the corpus; a plain change uses no Cliewen skill. For Cliewen work, the agent loads the relevant skill and follows a complete workflow without requiring a long prompt from the user.

## The lifecycle set

| Skill | Use it when | Durable output or hand-off |
|---|---|---|
| `clue-analysis` | Risk or an unknown should be investigated before committing to a design | Verified findings under `/docs/analysis`, then a plan or change |
| `clue-plan` | A goal needs a campaign or an active plan needs a semantic revision | A plan with verifiable milestones, then `clue-delta` |
| `clue-delta` | A light or full Cliewen change will mutate `main` | A complete branch, digested corpus, and verified pull request |
| `clue-verify` | A Cliewen pull request is about to open or update | An automatically agent-reviewed candidate, the pre-merge checklist, and review hand-off |
| `clue-extract` | An existing repository needs a one-time brownfield conversion | An inferred corpus reviewed through its first change loop |

## Why the skills stay separate

Each skill owns a lifecycle boundary and can be installed independently. Analysis should not load implementation mechanics it does not need, and verification should not improvise the proposal rules after the work is finished. Verification does own the recurring challenge-and-repair hand-off: it delegates review into a clean context where supported, returns findings to the implementing context, and requires a clean pass on the resulting commit.

The files are complete standalone artifacts, but repeated rules are generated from shared canonical sources. This keeps decision routing, change tiers, repository conventions, and the human review boundary identical across the set without creating runtime includes.

## Version agreement

Distributed Cliewen skills carry an ownership marker and the same version as the released binary. `clue validate` catches drift among the managed skills and between the skills and a released binary, while ignoring unrelated skills in the shared directory.

This makes process changes reviewable and reproducible: the agent guidance that shaped a branch is a versioned repository artifact, not an invisible service configuration.

Ready to practice the loop? Return to [Get started](./getting-started).
