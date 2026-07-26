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
# contract; TestSanity_InstallScriptsUseTheReleaseAssetContract holds this
# file to them.
#
# Everything is wrapped in main() and called on the last line, so a
# connection dropped mid-transfer cannot execute a half-read script.
set -eu

REPO="cliewen/cliewen"

die() { printf 'install.sh: %s\n' "$1" >&2; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"; }

detect_os() {
  case "$(uname -s)" in
    Darwin) echo darwin ;;
    Linux)  echo linux ;;
    *) die "unsupported operating system: $(uname -s). Install from source with: go install github.com/${REPO}/cmd/clue@latest" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64 | amd64)  echo amd64 ;;
    arm64 | aarch64) echo arm64 ;;
    *) die "unsupported architecture: $(uname -m). Install from source with: go install github.com/${REPO}/cmd/clue@latest" ;;
  esac
}

# The /releases/latest redirect ends in the tag, so no API token and no
# JSON parsing are needed — and it is not subject to api.github.com's
# unauthenticated rate limit, which a shared address can exhaust.
latest_tag() {
  if [ "$downloader" = curl ]; then
    curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest" \
      | sed 's|.*/||'
  else
    wget -qS --spider --max-redirect=0 "https://github.com/${REPO}/releases/latest" 2>&1 \
      | sed -n 's/.*[Ll]ocation: *//p' | tail -n 1 | tr -d '\r' | sed 's|.*/||'
  fi
}

main() {
  need uname
  need mkdir
  need install
  need mktemp

  if command -v curl >/dev/null 2>&1; then
    downloader=curl
    fetch() { curl -fsSL "$1" -o "$2"; }
  elif command -v wget >/dev/null 2>&1; then
    downloader=wget
    fetch() { wget -qO "$2" "$1"; }
  else
    die "neither curl nor wget is available"
  fi

  os=$(detect_os)
  arch=$(detect_arch)
  install_dir="${CLUE_INSTALL:-$HOME/.local/bin}"

  version="${CLUE_VERSION:-}"
  if [ -z "$version" ]; then
    tag=$(latest_tag) || die "could not reach github.com; set CLUE_VERSION=<x.y.z> and retry"
    [ -n "$tag" ] || die "could not determine the latest release; set CLUE_VERSION=<x.y.z> and retry"
    version="$tag"
  fi
  # Accept 0.7.0 or v0.7.0; the asset names carry bare semver (ADR-011).
  version="${version#v}"

  asset="clue-${version}-${os}-${arch}"
  base="https://github.com/${REPO}/releases/download/v${version}"

  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT
  trap 'rm -rf "$tmp"; exit 130' INT TERM

  printf 'Downloading %s\n' "$asset"
  fetch "${base}/${asset}" "${tmp}/${asset}" || die "download failed: ${base}/${asset}"
  fetch "${base}/SHA256SUMS" "${tmp}/SHA256SUMS" || die "download failed: ${base}/SHA256SUMS"

  # Verify before installing. --ignore-missing lets one checksum file cover
  # every published asset while we hold only one of them; it still fails
  # when nothing matched, so an unlisted asset cannot pass unchecked.
  printf 'Verifying checksum\n'
  if command -v sha256sum >/dev/null 2>&1; then
    (cd "$tmp" && sha256sum -c --ignore-missing SHA256SUMS >/dev/null) \
      || die "checksum verification failed for ${asset} — nothing was installed"
  elif command -v shasum >/dev/null 2>&1; then
    expected=$(grep " ${asset}\$" "${tmp}/SHA256SUMS" | awk '{print $1}')
    [ -n "$expected" ] || die "${asset} has no line in SHA256SUMS"
    actual=$(shasum -a 256 "${tmp}/${asset}" | awk '{print $1}')
    [ "$expected" = "$actual" ] || die "checksum verification failed for ${asset} — nothing was installed"
  else
    die "neither sha256sum nor shasum is available; refusing to install an unverified binary"
  fi

  mkdir -p "$install_dir"
  install -m 0755 "${tmp}/${asset}" "${install_dir}/clue"
  printf 'Installed clue %s to %s/clue\n' "$version" "$install_dir"

  case ":${PATH}:" in
    *":${install_dir}:"*)
      printf '\nRun `clue version` to confirm, then `clue init` in a repository.\n'
      ;;
    *)
      printf '\n%s is not on your PATH. Add it to your shell profile:\n\n    export PATH="%s:$PATH"\n\nThen open a new terminal and run `clue version`.\n' "$install_dir" "$install_dir"
      ;;
  esac
}

main "$@"
