# Remove the coordinated identity allocation the tracked-brief fixture starts with, leaving the seeded
# ledger: allocation is local again, as in a repository whose maintainer has not run `clue id coordinate`.
# A variant that finds nothing to remove fails the run, so a fixture change cannot make it a silent no-op.
# The commit message does not name the variant, so the agent under test cannot read it from the history.
[ -f .clue/id-coordination.yaml ] || { echo "variant removed nothing: no .clue/id-coordination.yaml in the fixture" >&2; exit 1; }
unlink .clue/id-coordination.yaml
git add -A
git commit -qm "Settings update"
git push -q origin HEAD:main 2>/dev/null || git push -q origin HEAD
# Coordination left an allocator branch on the remote; a repository that never coordinated has none, and an agent reads it as another contributor's claim.
git push -q origin --delete clue/id-allocator 2>/dev/null || { echo "variant removed nothing: no clue/id-allocator on the remote" >&2; exit 1; }
