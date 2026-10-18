---
id: P-025
type: plan
status: active
links: [G-009, G-019, VIS-001]
title: Cliewen commands are safe to re-run in a living repository
---

# P-025 — Cliewen commands are safe to re-run in a living repository

An adopter who keeps working in a repository after `clue init` keeps needing two things from the tool: regenerated README indexes and fresh identities. Today the first is only available through `clue init`, which also recreates files the adopter deliberately removed, and the second reserves an identity the moment it is asked for. Observed in the Tank Royale repository on 2026-10-18 while a maintainer worked on a change: running `clue init` only to refresh `docs/decisions/README.md` recreated `.clue/acceptance.yaml` (whose removal had been chosen to select local acceptance), `docs/vision.md`, six evidence-producer templates under `.clue/evidence/` and `.github/cliewen-wall.md`, twice in one session. Two identities were reserved only to look at the next number and were then abandoned in the ledger.

The observations come from one session and from reading the skill text. Nothing here has been checked against the tool's source, so each milestone starts by confirming what the command actually does.

## Milestones

| ID | Milestone (exit criterion) | Status | Evidence |
|---|---|---|---|
| M-110 | **Refreshing README indexes cannot resurrect a deliberate removal.** Either a command regenerates the indexes and nothing else, or `clue init` leaves a removed policy or scaffold file removed once the repository has adopted Cliewen. Disposable-Git evidence removes `.clue/acceptance.yaml`, the vision bootstrap and the evidence templates, refreshes the indexes, and finds all three still absent while an index row for a new artifact appears. The existing non-destructive behaviour of a first `clue init` is unchanged. | `done` | Already satisfied by `clue scaffold`, which regenerates the indexes and materializes nothing (CAP-005). In a disposable Git repository, after `.clue/acceptance.yaml`, the evidence templates and `.github/cliewen-wall.md` were removed, `clue scaffold` left all three absent and refreshed the index, while `clue init` recreated them. The gap was discoverability: the change-loop skill and `clue help` now name `clue scaffold` for an index refresh. Direct route. |
| M-111 | **An identity can be previewed without reserving it.** A documented option (for example a dry run) prints the next identity for a prefix and leaves the ledger byte-for-byte unchanged, including when the prefix is a milestone or a workspace artifact. Focused evidence covers the preview and confirms that the ordinary form still reserves. The option works whichever order the flags and the prefix are written in, or rejects a misplaced flag instead of misreading it. | `done` | `clue id next --dry-run <prefix>` prints the next identity and leaves the ledger byte-for-byte unchanged, for a milestone prefix as for any other; the ordinary form still reserves; a flag after the prefix is rejected. Focused tests in `cmd/clue/main_test.go`. Direct route, chosen over the recommended tracked route. |
| M-112 | **A change workspace's identity formats and the plan-less declaration are discoverable from the tool.** Validator messages and the change-loop skill name the identity prefix and allowed status for the tasks and open-questions files, and say how a change declares itself plan-less, without citing an internal rule an adopter cannot read. A new workspace written from the documentation alone passes `clue validate` on the first run. | `todo` | Not started. |
| M-113 | **The default pull request template does not mislead an outside contributor.** The template `clue init` installs as a repository's default is either a short contributor-facing description or clearly addressed to maintainers, and an outside contributor opening a pull request that is not a tracked change is not handed the full-loop checklist. Evidence opens a pull request in a disposable adopter repository without a tracked change and finds no unanswerable maintainer-only sections, while a tracked change still reaches its acceptance-brief gate. | `todo` | Not started. Observed: a first-time contributor's registration pull request received the maintainer checklist. |
| M-114 | **The tracked change loop publishes only when the maintainer has asked it to.** The skill states, before it branches or opens anything, whether it will push or open a hosted pull request, takes that from the acceptance policy and an explicit instruction, and otherwise keeps the work on a local branch. Evidence runs the loop under a policy that selects local acceptance and finds no push or pull request created without an instruction, and under pull-request acceptance finds the existing behaviour unchanged. | `todo` | Not started. Observed: two draft pull requests were opened for a maintainer who did not want any. |

## Challenge

The assumption most likely to undermine this plan is that the surprise is in the tool and not in how it was used here: a documented command for indexes or a dry run may already exist and simply not be listed in `clue --help` or the skills. The cheapest test is to read the command source and its documentation before building anything, and a milestone that finds the capability already present becomes a documentation fix. A different failure would be an implementation that passes its checks but stops `clue init` from repairing a repository whose scaffold really is missing; the evidence therefore keeps the first-run and repair paths separate from the deliberate-removal path.

## Boundaries

No change to the acceptance policy, the change loop, the identity allocation protocol or the meaning of any criterion is approved by this plan. Recreating files on request stays possible.
