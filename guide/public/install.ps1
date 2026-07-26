# Install clue, the Cliewen corpus judge, on Windows.
#
#   irm https://cliewen.dev/install.ps1 | iex
#
# The binary is verified against the release's SHA256SUMS before it is
# installed; a mismatch aborts without writing anything. Nothing requires
# an elevated shell — the default target is a directory you own, and the
# PATH entry is added to the user environment, not the machine.
#
# Options (environment variables):
#   CLUE_VERSION   release to install, e.g. 0.7.0   (default: latest)
#   CLUE_INSTALL   directory to install into        (default: %LOCALAPPDATA%\Programs\clue)
#
# This script downloads the same `clue-<version>-<os>-<arch>.exe` asset the
# release publishes for every other channel (ADR-030). Those names are an
# append-only contract; TestSanity_InstallScriptUsesTheReleaseAssetContract
# holds this file to them.

$ErrorActionPreference = 'Stop'

$repo = 'cliewen/cliewen'
$installDir = if ($env:CLUE_INSTALL) { $env:CLUE_INSTALL } else { Join-Path $env:LOCALAPPDATA 'Programs\clue' }

$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { 'amd64' }
    'ARM64' { 'arm64' }
    default { throw "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE). Install from source with: go install github.com/$repo/cmd/clue@latest" }
}

$version = $env:CLUE_VERSION
if (-not $version) {
    $latest = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" -Headers @{ 'User-Agent' = 'clue-install' }
    if (-not $latest.tag_name) { throw 'Could not determine the latest release; set CLUE_VERSION=<x.y.z> and retry.' }
    $version = $latest.tag_name -replace '^v', ''
}

$asset = "clue-$version-windows-$arch.exe"
$base = "https://github.com/$repo/releases/download/v$version"
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) ([System.IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $tmp | Out-Null

try {
    Write-Host "Downloading $asset"
    Invoke-WebRequest -Uri "$base/$asset" -OutFile (Join-Path $tmp $asset)
    Invoke-WebRequest -Uri "$base/SHA256SUMS" -OutFile (Join-Path $tmp 'SHA256SUMS')

    # Verify before installing.
    Write-Host 'Verifying checksum'
    $line = Get-Content (Join-Path $tmp 'SHA256SUMS') | Where-Object { $_ -match "\s$([regex]::Escape($asset))$" } | Select-Object -First 1
    if (-not $line) { throw "$asset has no line in SHA256SUMS" }
    $expected = ($line -split '\s+')[0]
    $actual = (Get-FileHash (Join-Path $tmp $asset) -Algorithm SHA256).Hash
    if ($expected -ine $actual) { throw "Checksum verification failed for $asset - nothing was installed" }

    if (-not (Test-Path $installDir)) { New-Item -ItemType Directory -Path $installDir -Force | Out-Null }
    Copy-Item (Join-Path $tmp $asset) (Join-Path $installDir 'clue.exe') -Force
    Write-Host "Installed clue $version to $installDir\clue.exe"
} finally {
    Remove-Item $tmp -Recurse -Force -ErrorAction SilentlyContinue
}

# Persist the PATH entry for the user, not the machine, so no elevation is
# needed. The current session is updated separately: a change made here
# does not reach an already-running shell.
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$installDir", 'User')
    Write-Host "Added $installDir to your user PATH."
    Write-Host 'Open a new terminal, then run `clue version` to confirm.'
} else {
    Write-Host 'Run `clue version` to confirm, then `clue init` in a repository.'
}
