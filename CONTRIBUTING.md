---
type: contributor-guide
title: Contributing to Cliewen
---

# Contributing to Cliewen

Thank you for helping improve Cliewen. Participation is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Choose the Right Route

- Suspected security vulnerability: follow the private reporting route in [SECURITY.md](SECURITY.md). Never disclose it in a public issue or pull request.
- Private conduct concern: use the [private conduct-reporting address](mailto:flemming&#46;n&#46;larsen&#43;cliewen-conduct&#64;gmail&#46;com) with the subject `[Cliewen Conduct] <private report>`. Never open a public conduct issue.
- Reproducible defect: open the structured bug form.
- Desired outcome or unmet need: open the proposed-goal form. A goal issue records demand for consideration; it does not add the goal to Cliewen's accepted plan.
- Work that leaves the accepted contract unchanged: use the direct route below.

## Before Starting a Change

Before editing, inspect the smallest relevant context and state `Recommended route: direct` or `Recommended route: tracked`, why, and what discovery would change that recommendation. Direct work leaves the accepted contract unchanged: observational analysis with a named consumer, a defect correction restoring an unchanged criterion, regression evidence for an unchanged criterion, in-contract configuration, refactoring, maintenance, and editorial work. Tracked work changes acceptance-criterion, capability, decision, policy, plan-promise, methodology, or uncovered-behavior meaning. Paths and diff size may warn but never decide meaning; uncertainty makes tracked the honest recommendation.

Reassess on semantic discovery and against the complete diff before integration. If direct work grows into tracked work, pause and recommend tracked. If the user explicitly declines, proceed as direct without making code, tests, or durable documentation untruthful and add the three PDR-042 trailers to the final authored commit: `Cliewen-Route: direct`, `Cliewen-Recommendation: tracked`, and `Cliewen-Override: user chose direct; <concise risk>`.

Direct work uses no CH number, plan declaration, workspace, digest, acceptance brief, or mandatory agentic review; it runs checks relevant to the changed surfaces. This repository selects local acceptance on main. Tracked work uses a private acceptance brief and human-run `clue accept`; direct work uses a human-controlled local merge after relevant checks. Agents never accept their own work. Its administrative `clue` version cut is a local direct-route specialization limited to the release files and relevant release checks. It is not a release process Cliewen imposes on adopters.

For a tracked change, search existing issues, pull requests, and the system-of-record under [`docs/`](docs/README.md). Every tracked change serves an existing plan item or explicitly declares itself plan-less. A contributor may initiate one tracked change at a time; direct work, reviewing, and helping update an existing pull request consume no tracked-change slot. Every tracked-change branch starts from accepted `main` and never from unmerged work unless a human explicitly authorizes a genuine dependency and its workspace and acceptance brief record that authorization.

After classification, load the smallest durable slice that can govern the work: run `clue context <id>` when the request names or resolves to an artifact; otherwise orient at `docs/README.md`, select the closest plan, capability, criterion, or decision, and run `clue context` from there. Follow additional edges only when the task discovers them.

Allocate the next `CH-xxx` identifier with `clue id next CH`, which reserves it in the repository's identity ledger, and mark it live with `clue id live CH-xxx` once the proposal exists. Never derive the number by reading git history: an identifier that a deleted artifact once used would be re-minted, and `clue validate` rejects an artifact missing from the ledger. Then create a descriptive Cliewen branch such as `ch-031-short-slug`.

## Run the Tracked Route

Before implementation, add `/changes/<CH-xxx-slug>/proposal.md`, `tasks.md`, and `open-questions.md`, and commit that proposal by itself on a branch from accepted `main`. Under this repository's local policy, follow the generated `clue-delta` review boundary for publication and handoff; a PR is optional publication, not the acceptance mechanism. Record unresolved decisions in `open-questions.md` and stop until a human answer can be captured as a typed decision.

A bug fix that restores an unchanged acceptance criterion may be direct work even though runtime behavior changes. If the accepted criterion itself must change, or the fix introduces behavior it does not cover, recommend the tracked route.

## Implement and Digest

Keep the change focused on its proposal, and give a task you mark `[-]` as infeasible its reason on the same line. Update permanent capability, acceptance-criteria, decision, constraint, architecture, and plan artifacts when their meaning changes.

Commit each working session that changed anything, and push the work branch as authorized for collaboration; neither action claims readiness or acceptance. Local acceptance does not require a PR. The review-boundary sections about draft PRs, hosted conversations and protected admission apply only when policy selects PR.

New or revised machine-proven acceptance criteria declare `Test-type: Unit`, `Integration`, `E2E`, or `Performance` and require framework-neutral exported executable evidence classified by that type and positive/negative direction, unless the criterion records `(single-direction)`. Repository-owned exporters attribute native or custom metadata to one executable, using a stable name/title fallback where needed; ordinary comments and container tags cannot complete the triple. Regenerate `.clue/evidence.yaml` when its inputs change. A genuine `Test-type: Human` criterion is named in the acceptance brief as its proof and needs no code test; `@draft` exempts only one genuinely not-yet-proven criterion; an unannotated legacy criterion retains its one-supported-reference contract. `clue validate` validates declarations and references but does not execute tests. Never weaken a test, lint rule, or quality gate to make a build pass. If a Cliewen-owned skill changes, edit `internal/skills/source/` and run `go generate ./internal/skills`; do not edit `.agents/skills/` or `internal/scaffold/templates/skills/` directly.

Before review, digest a tracked change into the permanent corpus, update its plan bookkeeping, and remove its `/changes/` workspace. Add a [`CHANGELOG.md`](CHANGELOG.md) entry when the change alters what an adopter receives — `clue`'s behaviour, a generated skill's text, or an artifact `clue init` or `clue scaffold` materializes into an adopter repository. A change confined to this repository's own corpus, contributor guidance, and local conventions adds no release note even when it is a tracked Cliewen change, and neither does a simple editorial correction; a change spanning both writes the entry for the adopter-visible part alone. Those categories are a shortcut and the test above is the rule: `.github/pull_request_template.md` is held byte-identical to what `clue init` writes and `.github/workflows/clue-validation.yml` is the reusable workflow an adopter's caller references, so a change to either owes an entry even though `.github/` reads as this repository's own CI ([C-002](docs/constraints/C-002-changelog-per-user-visible-change.md)). The final tree proposed for merge must not contain transient change files.

The generated `.github/workflows/clue.yml` is a thin caller for Cliewen's upstream reusable validation workflow. Keep runner labels, binary source, and writable install-directory choices in that caller; do not copy validation steps or action references into it. A reusable-workflow reference update is the reviewed path for importing upstream scope, warning, acceptance-brief, and digest-gate fixes.

## Verify Locally

For direct work, run only checks relevant to its changed surfaces. A guide-Markdown-only edit runs the whitespace check below and `npm run guide:build`; an analysis-only corpus change also runs `clue validate`.

For a Cliewen change, commit the complete candidate, then run the repository's full mechanical gates against that commit:

```text
go build ./...
go test ./... -coverprofile coverage.out
go tool cover -func coverage.out
go run ./cmd/clue validate --forbid-changes
git diff --check $(git merge-base HEAD origin/main) HEAD
```

On Windows, when Bitdefender blocks the temporary `clue.exe` produced by `go run`, use this block instead. It builds the same judge at the stable, ignored checkout path `clue.exe` and serializes package builds; it runs the full test suite and retains every gate:

```text
go build ./...
go test -p=1 ./... -coverprofile coverage.out
go tool cover -func coverage.out
go build -o ./clue.exe ./cmd/clue
./clue.exe validate --forbid-changes
git diff --check $(git merge-base HEAD origin/main) HEAD
```

The stable build was observed to run with Bitdefender enabled on Windows on 2026-10-09. This changes no antivirus settings. If that exact executable is also blocked, inspect Bitdefender's notification and use its supported exception interface for that executable; do not exclude the entire temporary directory or the Go toolchain. Rebuild after source changes so validation uses the candidate's code. This is a local verification build, not an installed-binary upgrade or release.

Total Go statement coverage must remain at least 80%. `clue-verify` then automatically reviews that same commit before the tracked candidate is handed to the human for local acceptance. A coding-agent host with context-isolated delegation starts a fresh read-only reviewer; other hosts disclose an in-context fallback. The loop owns its classification regardless of the reviewer brief: a blocking finding is actionable and enters the hosted repair lifecycle; an advisory is a non-actionable observation for the readiness gate and stays in the verification handoff. Counts and arithmetic disagreements are advisory, while a wrong, missing, or reused identity remains blocking, and the reviewer spends no pass re-deriving figures. A blocking finding returns to the implementing context, is committed, checked against that commit, and reviewed again — scoped to what changed and the carriers it declares — until the current commit receives a pass with no blocking findings. An advisory repair may ride before a pass already required by a blocking repair; an advisory first reported by a pass with no blocking findings stays in the handoff for a later change so the handed-off candidate remains the exact reviewed commit without making the advisory a merge gate. The loop runs up to the maximum number of passes [C-017](docs/constraints/C-017-agentic-review-loop-is-bounded.md) states for this repository, and a further pass runs only after a pass with a blocking finding. Reaching the maximum with blocking findings outstanding stops the loop, reports them to the maintainer, and asks whether to run more; it never permits handing off a candidate as ready; the work branch shows the unfinished state. The final verification evidence identifies the review mode, reviewed commit, number of passes run, and advisory findings left open.

## Writing Multi-Line PR and Commit Text on Windows

**Trigger:** Authoring a `gh pr create`/`gh pr edit`/`gh issue create` body, or a multi-line git commit message, from a PowerShell shell.

**Prerequisites:** PowerShell is the shell in use (not Bash or another POSIX shell), and the text contains backtick-quoted inline code (for example `` `npm run build` ``) or is meant to contain literal newlines.

**Procedure:** Write the body or message to a file with a file-writing tool, then pass it with `--body-file <path>` (`gh pr create`, `gh pr edit`, `gh issue create`) or `git commit -F <path>`. Prefer running the `gh`/`git` command itself through a POSIX-shell tool when one is available, since it has no backtick-escape hazard. If PowerShell must build the string in place, use a single-quoted here-string (`@'...'@`); never a double-quoted string or `@"..."@` for text containing backticks or intended newlines.

**Expected result:** The published body or message renders with intact headers, unbroken inline code, and correct line breaks.

**Recovery:** If a body or message was already published, edit it with the file-based form above (`gh pr edit --body-file <path>`, or `git commit --amend -F <path>` for an unpushed commit) and verify by re-reading the rendered result.

**Why:** PowerShell's escape character is backtick, not backslash. A backtick followed by `n` — as inside backtick-quoted inline code — collapses to a real newline and drops the `n`; backtick followed by `f` becomes a form-feed. Meanwhile a literal `\n` typed with intent to mean "newline" is not an escape in a PowerShell double-quoted string, so it survives as literal text instead of becoming a line break. This combination corrupted this repository's [PR #49](https://github.com/cliewen/cliewen/pull/49) description at first publication — `npm` lost its leading letter to a swallowed newline and a backtick-quoted commit hash turned into a form-feed character — and it was corrected shortly after using the file-based procedure above, so the live body now reads clean; the original corrupted text is still visible in the PR's content-edit history, for example via `gh api graphql -f query='{ repository(owner: "cliewen", name: "cliewen") { pullRequest(number: 49) { userContentEdits(first: 10) { nodes { editedAt diff } } } } }'`. Observed on Windows PowerShell with `gh`; not yet confirmed under other shells or hosts.

## Propose for Human Acceptance

For direct work, provide a summary, exact candidate and base commits, relevant verification and remaining uncertainty. The human integrates the branch locally; a CH identity, tracked brief and automatic agentic review are not required.

For tracked work, follow `clue-delta`'s local review boundary. Commit the complete candidate, pass mechanical checks and automatic agentic review, and fill the private acceptance brief for the exact candidate, base and reviewed revisions. Run `clue accept <candidate-sha> --base <base-sha> --brief <file> --check` on the clean main checkout at the base. Hand the human the brief and the same command without `--check`; only the human runs the accepting command and confirms. It creates the history-preserving merge but never pushes. The human publishes accepted main through the repository's authorized process.

An edit or advanced base requires renewed verification, review and a new brief. Cancellation or failed preflight integrates nothing. Agents never run accepting `clue accept`, merge their own work into main, or push main without explicit authorization. Repository hosting permissions still apply; local policy does not change branch protection. A transition from PR to local is accepted under the old PR policy.

The administrative release cut remains direct work without CH identity, workspace, digest, tracked acceptance brief or agentic review. Its focused release gates run before a human-controlled local merge; publication starts when accepted main is pushed. Existing release PRs can finish under their earlier handoff, but future releases do not require PR creation.

Cliewen does not currently require a Contributor License Agreement or Developer Certificate of Origin sign-off.
