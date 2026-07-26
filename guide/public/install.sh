#!/bin/sh
# Install clue, the Cliewen corpus judge, on macOS or Linux.
#
#   curl -fsSL https://cliewen.dev/install.sh | sh
#
# The binary is verified against the release's SHA256SUMS before it is
# installed; a mismatch aborts without writing anything. Nothing runs with
# elevated privileges — the default target is a directory you own.
#
# Options (environment variables):
#   CLUE_VERSION   release to install, e.g. 0.7.0        (default: latest)
#   CLUE_INSTALL   directory to install into              (default: ~/.local/bin)
#
# This script downloads the same `clue-<version>-<os>-<arch>` asset an
# adopter's CI wall installs (ADR-030). Those names are an append-only
# contract; TestSanity_InstallScriptUsesTheReleaseAssetContract holds this
# file to them.
set -eu

REPO="cliewen/cliewen"
INSTALL_DIR="${CLUE_INSTALL:-$HOME/.local/bin}"

die() { printf 'install.sh: %s\n' "$1" >&2; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"; }

need uname
need mkdir
need install

if command -v curl >/dev/null 2>&1; then
  fetch() { curl -fsSL "$1" -o "$2"; }
  fetch_stdout() { curl -fsSL "$1"; }
elif command -v wget >/dev/null 2>&1; then
  fetch() { wget -qO "$2" "$1"; }
  fetch_stdout() { wget -qO- "$1"; }
else
  die "neither curl nor wget is available"
fi

case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux)  os=linux ;;
  *)      die "unsupported operating system: $(uname -s). Install from source with: go install github.com/${REPO}/cmd/clue@latest" ;;
esac

case "$(uname -m)" in
  x86_64 | amd64)  arch=amd64 ;;
  arm64 | aarch64) arch=arm64 ;;
  *)               die "unsupported architecture: $(uname -m). Install from source with: go install github.com/${REPO}/cmd/clue@latest" ;;
esac

version="${CLUE_VERSION:-}"
if [ -z "$version" ]; then
  # The redirect target of /releases/latest ends in the tag; no API token
  # and no JSON parser needed.
  tag=$(fetch_stdout "https://api.github.com/repos/${REPO}/releases/latest" \
    | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
    | head -n 1) || true
  [ -n "${tag:-}" ] || die "could not determine the latest release; set CLUE_VERSION=<x.y.z> and retry"
  version="${tag#v}"
fi

asset="clue-${version}-${os}-${arch}"
base="https://github.com/${REPO}/releases/download/v${version}"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT INT TERM

printf 'Downloading %s\n' "$asset"
fetch "${base}/${asset}" "${tmp}/${asset}" || die "download failed: ${base}/${asset}"
fetch "${base}/SHA256SUMS" "${tmp}/SHA256SUMS" || die "download failed: ${base}/SHA256SUMS"

# Verify before installing. --ignore-missing lets one checksum file cover
# every published asset while we hold only one of them.
printf 'Verifying checksum\n'
if command -v sha256sum >/dev/null 2>&1; then
  (cd "$tmp" && sha256sum -c --ignore-missing SHA256SUMS >/dev/null) || die "checksum verification failed for ${asset} — nothing was installed"
elif command -v shasum >/dev/null 2>&1; then
  expected=$(grep " ${asset}\$" "${tmp}/SHA256SUMS" | awk '{print $1}')
  [ -n "$expected" ] || die "${asset} has no line in SHA256SUMS"
  actual=$(shasum -a 256 "${tmp}/${asset}" | awk '{print $1}')
  [ "$expected" = "$actual" ] || die "checksum verification failed for ${asset} — nothing was installed"
else
  die "neither sha256sum nor shasum is available; refusing to install an unverified binary"
fi

mkdir -p "$INSTALL_DIR"
install -m 0755 "${tmp}/${asset}" "${INSTALL_DIR}/clue"
printf 'Installed clue %s to %s/clue\n' "$version" "$INSTALL_DIR"

case ":${PATH}:" in
  *":${INSTALL_DIR}:"*)
    printf '\nRun `clue version` to confirm, then `clue init` in a repository.\n'
    ;;
  *)
    printf '\n%s is not on your PATH. Add it to your shell profile:\n\n    export PATH="%s:$PATH"\n\nThen open a new terminal and run `clue version`.\n' "$INSTALL_DIR" "$INSTALL_DIR"
    ;;
esac
