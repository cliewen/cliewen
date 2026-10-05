set -e
mkdir -p /home/node/work && cd /home/node/work
git init -q
printf '# Demo\n\nA small CSV tool. Teh quick start is below.\n' > README.md
clue init . >/dev/null
git add -A
git commit -qm init
