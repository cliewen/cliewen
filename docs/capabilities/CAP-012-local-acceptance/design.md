---
id: CAP-012-design
type: design
status: active
links: [CAP-012, PDR-067, ADR-044]
title: Local acceptance records and integration
---

# Local acceptance

## Opt-in and preparation

The accepted base and candidate must resolve to the same local policy. Without `.clue/acceptance.yaml`, that policy is local acceptance on `main`; an explicit local/main file is equivalent. A file selects `mode: local` or `mode: pr` and an explicit valid `branch`. Init materializes local by default regardless of existing Cliewen state; `--acceptance=pr` actively selects PR. The shared strict reader is used by init, migration and acceptance. Init resolves conflicts before writes and preserves existing explicit policy. Migration validates policy but never creates or rewrites it. MIG-021 and historical PR fallback are removed by PDR-069. Repository role does not select acceptance policy: source and adopter repositories share the same local/PR behavior under PDR-070. This source repository explicitly selects local on main. A candidate cannot override accepted PR policy to accept itself locally. Identity allocation retains its ledger and coordination behavior, independent of acceptance.

`clue accept <candidate-sha> --base <base-sha> --brief <file> --check` performs preflight from the clean integration checkout at the base. Both IDs are full lowercase commit hashes. The brief is Markdown with frontmatter `type: acceptance-brief`, `title`, `change`, `candidate`, `base`, and `reviewed`; the reviewed revision equals the candidate. Its required nonempty sections are Intent, Criteria and evidence, Binding decisions, Verification, and Review. Use the shipped `.clue/acceptance/brief.md` template and store the prepared copy under the worktree-specific location returned by `git rev-parse --git-path clue/acceptance/CH-xxx.md`. The command rejects comments/placeholders and structural omissions but cannot judge whether prose is true or adequate.

## State and history boundary

The command reads committed blobs into private temporary snapshots without checkout filters, hooks, or export attributes. It verifies declared local policy in the base, candidate ancestry, a matching proposal introduced in the candidate's reachable history after the base, and the candidate's corpus with the same validation and parity checks as `clue validate --forbid-changes`. The snapshot resolves file, directory and chained links in the committed revision graph and copies regular targets without OS links. It rejects escaping or missing targets, metadata paths, cycles, excessive expansion and submodules; acceptance retains the original tree and link modes. `--check` changes no refs, index or tracked files and never calls exporters or test runners. Repository verification and review are recorded declarations, not command-authenticated executions.

The human decides acceptance. They can run the same command without `--check` on an interactive terminal, or approve an exact bound record for delegated execution with `--approval <file>`. The strict YAML record names decision, candidate, base, brief-sha256, approved-by, approval-source, statement and recorded-at (RFC3339). Its actual statement and source are procedural provenance, not authenticated identity. It displays the entire brief and trust boundary and requires `accept CH-xxx`. After interactive confirmation or approval-record verification it repeats preflight and refuses changed inputs. The recorded decision must match the exact candidate, base and complete brief bytes; the approval record is retained with the brief. It creates a commit from the candidate tree with ordered parents base and candidate and the complete brief in its message. A prepared Git ref transaction locks the integration branch and compares its old value before `read-tree` updates the index and working tree; the transaction then advances the branch. Each Git command owns an isolated empty hooks directory for its full lifetime and disables executable filesystem monitoring. There is no merge conflict resolution, Git-hook execution, or push. The committer is the configured Git identity, which does not establish human presence.

## Recovery and limits

Cancellation and ordinary preflight failures leave integration state untouched. A moved base must be merged into the change branch before repeating verification/review and generating the new brief. A changed candidate invalidates the prior brief and review. Git commit objects keep proposal history and the acceptance brief reachable after branch deletion. Find the brief with `git show --no-patch --format=%B <acceptance-commit>` and the candidate at its second parent.

Process interruption can leave the candidate tree staged before the acceptance ref advances. Inspect `git status`, `git log -1`, and the expected base/candidate IDs. If HEAD is the acceptance merge and the tree is clean, acceptance completed; do not repeat it. If HEAD is still the base, preserve unrelated edits first. When the index and working tree are exactly the recorded candidate tree with no intervening edits, `git read-tree -m -u <candidate-sha> <base-sha>` restores the base without a force reset; then repeat preflight. If state differs, stop for the human to recover it. A failed ref commit attempts that same non-forcing rollback and reports the outcome. No prompt or local Git lock can prevent an actor with equivalent operating-system access from bypassing the method.

Human evidence records actual observations. The command can enforce structure and revision binding but cannot establish that the observer is human, that tests passed, or that the outcome is wanted. A maintainer usability trial remains human evidence rather than something a disposable-Git test can supply.
