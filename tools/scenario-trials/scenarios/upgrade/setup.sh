set -e
mkdir -p /home/node/work && cd /home/node/work
git init -q
printf '# Demo\n\nA small tool.\n' > README.md
clue init . >/dev/null
git add -A
git commit -qm init
# A newer release exists: `clue latest` reports one, every other command is the real clue.
mkdir -p /home/node/bin
cat > /home/node/bin/clue <<'SHIM'
#!/bin/sh
if [ "$1" = "latest" ]; then
  case " $* " in
    *" --quiet "*) echo "v99.0.0" ;;
    *) echo "A newer Cliewen release is available: v99.0.0 (installed: dev). See https://github.com/cliewen/cliewen/releases for the installation route for this machine." ;;
  esac
  exit 0
fi
exec /usr/local/bin/clue "$@"
SHIM
chmod +x /home/node/bin/clue
