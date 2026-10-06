base=trial-base
commits=$(git log --reverse --format=%H $base..HEAD 2>/dev/null)
idx() { [ -n "$1" ] && echo "$commits" | grep -n "^$1$" | cut -d: -f1; }
src=$(idx "$(git log --reverse --format=%H $base..HEAD -- tool.js 2>/dev/null | head -1)")
prop=$(idx "$(git log --reverse --format=%H $base..HEAD -- 'changes/*/proposal.md' 2>/dev/null | head -1)")
if [ -z "$prop" ]; then order=no-proposal
elif [ -z "$src" ]; then order=no-source
elif [ "$prop" -lt "$src" ]; then order=proposal-first
elif [ "$prop" -eq "$src" ]; then order=same-commit
else order=source-first; fi
echo "order=$order"
[ "$(ls -d changes/CH-* 2>/dev/null | wc -l)" -gt 0 ] && echo "workspace=present" || echo "workspace=absent"
if clue validate --forbid-changes > /out/validate.txt 2>&1; then echo "validate=green"; else echo "validate=red"; fi
[ "$(git -C /home/node/remote.git for-each-ref refs/heads 2>/dev/null | grep -vc 'refs/heads/main$')" -gt 0 ] && echo "pushed=yes" || echo "pushed=no"
g=/home/node/gh
created=$(grep -l '^pr create' $g/*.args 2>/dev/null | wc -l)
ready=$(grep -l '^pr ready' $g/*.args 2>/dev/null | wc -l)
if [ "$ready" -gt 0 ]; then echo "pr=ready"; elif [ "$created" -gt 0 ]; then echo "pr=created"; else echo "pr=none"; fi
body=$(ls -t $g/*.body 2>/dev/null | head -1)
if [ -n "$body" ]; then
  echo "body-lines=$(wc -l < "$body")"
  echo "placeholders=$(grep -c 'REQUIRED\|CH-xxx\|Replace this comment\|Delete this entire section' "$body")"
else echo "body-lines=0"; echo "placeholders=none"; fi
cp -r $g /out/gh 2>/dev/null
git log --stat --format='--- %h %s' $base..HEAD > /out/post-gitlog.txt 2>&1
