# Operate Cliewen safely

This page is for a repository that has finished the disposable trial and is deciding whether to keep Cliewen. It describes the support boundary that ships today and the recovery paths that preserve your repository's history and review boundary.

## What Cliewen ships and checks

| Surface | Current support |
|---|---|
| `clue` | Versioned release binaries for Windows, macOS, and Linux on amd64 and arm64; source installs are also available through Go |
| Test evidence | Classified Go test names such as `TestAC001_UnitPositive_…` or normalized `TestSNAPSQS001_UnitPositive_…`; Java and Kotlin executables whose own JUnit method annotations carry `@Tag("AC_001") @Tag("UNIT") @Tag("POSITIVE")` or `@Tag("SNAP_SQS_001")`, or whose stable name is `testAC001_UnitPositive_…` or `testSNAPSQS001_UnitPositive_…`; and Cucumber scenario tags such as `@AC-001 @e2e @positive` or `@SNAP-SQS-001 @unit @positive` |
| Agent guidance | Five generated Cliewen skills in `.agents/skills/`, mirrored to `.claude/skills/` for Claude Code |
| GitHub CI | `clue init` writes a thin caller for Cliewen's upstream reusable workflow; choose the runner and binary source, then require its exact `validate` check |
| Full-change merge | Human-controlled merge commits only; disable squash and rebase-and-merge so the proposal, implementation, digest, and durable corpus commits remain reachable from `main` |
| Validation | `clue validate` checks the repository-local corpus, generated indexes, skill ownership/version drift, and active-criterion evidence declarations and references; `--forbid-changes` also rejects an undigested `/changes/` workspace |
| Focused context | `clue context <id>` prints the named artifact and its transitive outgoing-link dependencies; criterion and milestone IDs resolve to their declaring artifact |

Cliewen does not run your tests, does not synchronize tickets or wikis, does not update installed files in the background, and does not validate evidence across repositories. It harvests only the supported Go, JVM, and Cucumber conventions above. On the JVM, all three parts belong to the same executable; class tags and metadata split across methods do not count, and ambiguous or unsupported source syntax is diagnosed instead of guessed. A different framework needs the stable JVM named-executable form or another supported profile before its references can satisfy `clue validate`; do not treat an arbitrary comment or tag as equivalent evidence.

A new or revised machine-proven criterion declares `Test-type: Unit`, `Integration`, `E2E`, or `Performance`, and the validator requires supported evidence classified with that type in positive and negative directions; `(single-direction)` is the explicit narrow exception. An unannotated legacy criterion retains its one-supported-reference rule. `Test-type: Human` uses the pull request acceptance brief rather than code evidence, and `@draft` exempts only the individual criterion that is not yet proven.

## Preserve the full-change archive

For a full Cliewen change, configure the protected default branch to allow the hosting provider's **merge commit** mode and disable **squash and merge** and **rebase and merge**. The merge commit keeps the exact proposal, implementation, digest, and durable corpus commits reachable from `main`; squash and rebase-and-merge can produce the same final tree while discarding or rewriting the reviewed branch chain. A local rebase before first publication is allowed, but the human acceptance mode remains a merge commit.

Because the merge shape is set per branch or per repository rather than per pull request, a default branch restricted this way restricts plain changes into it too; route work that needs another merge shape to a branch the rule does not target. Before adoption, run the [CI wall's branch-protection probe](./ci-wall#_5-prove-failure-blocks-merge) and verify both the protected branch's allowed merge methods and the repository-wide pull-request settings. If the forge cannot enforce this boundary, it is outside Cliewen's supported full-change adoption path rather than an equivalent configuration.

## Upgrade one coordinated set

Keep the binary, generated skills, and CI caller on the same release when you upgrade. First make the current repository green and branch through its normal review process. Then choose the release in the [release list](https://github.com/cliewen/cliewen/releases), verify the new platform binary against that release's `SHA256SUMS`, and confirm `clue version` prints the chosen version.

`clue init` is deliberately non-destructive: it skips existing files, so it is not an updater. Replace the five Cliewen-owned skill directories from the selected release's source tree in both `.agents/skills/` and `.claude/skills/`; leave unrelated third-party skills alone. Update `.github/workflows/clue.yml`'s upstream `uses` reference and `clue-version` together, then either replace `.github/tools/clue-<version>-linux-amd64` and `.github/tools/SHA256SUMS` from that same release or keep `clue-source: release` and let the upstream workflow download them. If your caller selects a self-hosted runner, keep its runner-label JSON and writable `clue-install-directory` unchanged unless the repository policy changes. Keep the existing required `validate` check in place throughout the upgrade; the stable job name lets the same required check verify this pull request. Make it required only when you arm the wall for the first time.

If a released binary reports skill drift, do not edit a version number to silence it. Install the matching released binary or replace the complete generated skill set with the matching version, then run `clue validate`. A checkout build reports `dev` and cannot detect binary-to-skill release drift, so use a released binary for this check.

Re-running the install script moves the binary and nothing else. In a repository already using Cliewen that produces exactly the drift report above, because the skills are committed repository files no installer can update: the machine moved ahead of the repository. This is the check working, not a broken upgrade. Resolve it by completing the coordinated set — skills, the caller reference, and either the vendored assets or its release-source version — in a normal reviewed change, or by pinning the release the repository still carries with `clue-version=<x.y.z>` if you are not ready to upgrade it yet.

## Recover without bypassing the evidence

| Situation | Safe response |
|---|---|
| `clue init` reports a skipped file | Read the skipped path and keep its local prose. Add the missing scaffolded structure or marked index block deliberately; rerunning `init` will not replace it for you. |
| `clue validate` fails | Run it from the repository root, read the named file and rule, and repair the corpus or evidence on a branch. Run your normal test runner too: `clue` does not execute tests. Do not delete a rule, use an untracked copy, or bypass CI to obtain green. |
| CI rejects a transient workspace | Finish the change's tasks, digest its durable documentation into `docs/`, remove `/changes/CH-xxx-*`, and rerun `clue validate --forbid-changes`. |
| Extraction is still unmerged | Close the pull request and delete its branch. The target repository was not changed by the extraction workflow. |
| An accepted extraction or adoption must be undone | Create a new reviewed change that restores the previous repository state or retires the incorrect intent, then validate it. Do not force-push or leave two competing systems of record. |
| You no longer want the local tool | Remove the installed `clue` binary from the directory you added to your user `PATH`. This does not alter existing repositories; remove or retire their Cliewen files only through their normal reviewed change process. |

## Evidence from other repositories

Cliewen's [hyperfine](https://github.com/sharkdp/hyperfine) and [es-toolkit](https://github.com/toss/es-toolkit) work were read-and-apply foreign-soil trials, not adoptions. They produced findings without branches, pull requests, or tracked changes in those repositories. They are useful evidence about methodology boundaries, not proof that those projects use or endorse Cliewen.

When your repository's ownership, test evidence, or merge boundary cannot meet these conditions, [keep the existing lightweight notes and tests](./adoption#when-cliewen-is-a-poor-fit) instead of forcing an adoption.

## Next

[Use the adoption guide to choose the smallest practice your repository needs.](./adoption)
