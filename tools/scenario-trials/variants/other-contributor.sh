# Local allocation, as in local-allocation, and another clone's change branch already pushed to the bare
# remote, so the agent under test can find another contributor with plain Git and no forge.
# A variant that finds nothing to remove fails the run, so a fixture change cannot make it a silent no-op.
[ -f .clue/id-coordination.yaml ] || { echo "variant removed nothing: no .clue/id-coordination.yaml in the fixture" >&2; exit 1; }
unlink .clue/id-coordination.yaml
git add -A
git commit -qm "Settings update"
git push -q origin HEAD:main 2>/dev/null || git push -q origin HEAD
home=$(git rev-parse --abbrev-ref HEAD)
git checkout -q -b ch-150-colleague-change
printf 'A colleague is working on a separate change.\n' > colleague-notes.txt
git add colleague-notes.txt
git commit -qm "Start a colleague's change"
git push -q origin ch-150-colleague-change
git checkout -q "$home"
git branch -q -D ch-150-colleague-change
git fetch -q origin
