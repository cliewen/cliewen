---
id: CAP-012-design
type: design
status: active
links: [CAP-012, PDR-067, ADR-044]
title: Local acceptance records and integration
---

# Local acceptance

## Opt-in and preparation

The accepted base and candidate both carry `.clue/acceptance.yaml` with `mode: local` and an explicit `branch`, such as `main`. The adopter declares the same choice in repository instructions. Absent opt-in means PR acceptance. `role: source` refuses local acceptance. A candidate cannot authorize its own change of mode. Identity allocation retains its existing ledger and coordination behavior, independent of acceptance.

`clue accept <candidate-sha> --base <base-sha> --brief <file> --check` performs preflight from the clean integration checkout at the base. Both IDs are full lowercase commit hashes. The brief is Markdown with frontmatter `type: acceptance-brief`, `title`, `change`, `candidate`, `base`, and `reviewed`; the reviewed revision equals the candidate. Its required nonempty sections are Intent, Criteria and evidence, Binding decisions, Verification, and Review. Use the shipped `.clue/acceptance/brief.md` template and store the prepared copy under the worktree-specific location returned by `git rev-parse --git-path clue/acceptance/CH-xxx.md`. The command rejects comments/placeholders and structural omissions but cannot judge whether prose is true or adequate.

## State and history boundary

The command reads committed blobs into private temporary snapshots without checkout filters, hooks, or export attributes. It verifies opt-in in the base, candidate ancestry, a matching proposal introduced in the candidate's reachable history after the base, and the candidate's corpus with the same validation and parity checks as `clue validate --forbid-changes`. The first version refuses tracked symlinks and submodules rather than following external files or omitting proof. `--check` changes no refs, index or tracked files and never calls exporters or test runners. Repository verification and review are recorded declarations, not command-authenticated executions.

The human runs the same command without `--check` on an interactive terminal. It displays the entire brief and trust boundary and requires `accept CH-xxx`. After confirmation it repeats preflight and refuses changed inputs. It creates a commit from the candidate tree with ordered parents base and candidate and the complete brief in its message. A prepared Git ref transaction locks the integration branch and compares its old value before `read-tree` updates the index and working tree; the transaction then advances the branch. Each Git command owns an isolated empty hooks directory for its full lifetime and disables executable filesystem monitoring. There is no merge conflict resolution, Git-hook execution, or push. The committer is the configured Git identity, which does not establish human presence.

## Recovery and limits

Cancellation and ordinary preflight failures leave integration state untouched. A moved base must be merged into the change branch before repeating verification/review and generating the new brief. A changed candidate invalidates the prior brief and review. Git commit objects keep proposal history and the acceptance brief reachable after branch deletion. Find the brief with `git show --no-patch --format=%B <acceptance-commit>` and the candidate at its second parent.

Process interruption can leave the candidate tree staged before the acceptance ref advances. Inspect `git status`, `git log -1`, and the expected base/candidate IDs. If HEAD is the acceptance merge and the tree is clean, acceptance completed; do not repeat it. If HEAD is still the base, preserve unrelated edits first. When the index and working tree are exactly the recorded candidate tree with no intervening edits, `git read-tree -m -u <candidate-sha> <base-sha>` restores the base without a force reset; then repeat preflight. If state differs, stop for the human to recover it. A failed ref commit attempts that same non-forcing rollback and reports the outcome. No prompt or local Git lock can prevent an actor with equivalent operating-system access from bypassing the method.

Human evidence records actual observations. The command can enforce structure and revision binding but cannot establish that the observer is human, that tests passed, or that the outcome is wanted. A maintainer usability trial remains human evidence rather than something a disposable-Git test can supply.
