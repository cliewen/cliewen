# Operate Cliewen safely

This page is for a repository that has finished the disposable trial and is deciding whether to keep Cliewen. It describes the support boundary that ships today and the recovery paths that preserve your repository's history and review boundary.

## What Cliewen ships and checks

| Surface | Current support |
|---|---|
| `clue` | Versioned release binaries for Windows, macOS, and Linux on amd64 and arm64; source installs are also available through Go |
| Test evidence | Go test names such as `TestAC001_…`, plus Java and Kotlin JUnit `@Tag("AC_001")` references |
| Agent guidance | Five generated Cliewen skills in `.agents/skills/`, mirrored to `.claude/skills/` for Claude Code |
| GitHub CI | `clue init` writes an unarmed GitHub Actions workflow; you vendor the pinned Linux binary and checksum, then require its `validate` check |
| Validation | `clue validate` checks the repository-local corpus, generated indexes, skill ownership/version drift, and active-criterion test references; `--forbid-changes` also rejects an undigested `/changes/` workspace |

Cliewen does not run your tests, synchronize tickets or wikis, does not update installed files in the background, or validate evidence across repositories. It harvests only the Go and JVM conventions above. A different framework needs a supported language profile before its test references can satisfy `clue validate`; do not treat an arbitrary comment or tag as equivalent evidence.

The positive-and-negative evidence pair remains a method and review requirement. The current validator verifies that at least one supported reference exists; it does not count or classify that pair.

## Upgrade one coordinated set

Keep the binary, generated skills, and vendored CI binary on the same release when you upgrade. First make the current repository green and branch through its normal review process. Then choose the release in the [release list](https://github.com/cliewen/cliewen/releases), verify the new platform binary against that release's `SHA256SUMS`, and confirm `clue version` prints the chosen version.

`clue init` is deliberately non-destructive: it skips existing files, so it is not an updater. Replace the five Cliewen-owned skill directories from the selected release's source tree in both `.agents/skills/` and `.claude/skills/`; leave unrelated third-party skills alone. Update `.github/workflows/clue.yml`'s `CLUE_VERSION`, then replace both `.github/tools/clue-<version>-linux-amd64` and `.github/tools/SHA256SUMS` from that same release and verify the checksum before committing. Let the pull request run `clue validate --forbid-changes` before making `validate` required again.

If a released binary reports skill drift, do not edit a version number to silence it. Install the matching released binary or replace the complete generated skill set with the matching version, then run `clue validate`. A checkout build reports `dev` and cannot detect binary-to-skill release drift, so use a released binary for this check.

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
