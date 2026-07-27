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
# append-only contract; TestSanity_InstallScriptsUseTheReleaseAssetContract
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
    # The /releases/latest redirect ends in the tag: no API token, and not
    # subject to api.github.com's unauthenticated rate limit, which a
    # shared address can exhaust.
    #
    # The redirect is followed rather than inspected, because reading an
    # unfollowed 302's Location header differs between Windows PowerShell
    # 5.1 and PowerShell 7 — and 5.1 is what `powershell.exe` starts on
    # every Windows machine. The final URI is exposed under two different
    # names instead, and both are cheap to ask for.
    $head = Invoke-WebRequest -Uri "https://github.com/$repo/releases/latest" -Method Head -UseBasicParsing
    $final = $head.BaseResponse.RequestMessage.RequestUri   # PowerShell 7
    if (-not $final) { $final = $head.BaseResponse.ResponseUri }  # Windows PowerShell 5.1
    if (-not $final) { throw 'Could not determine the latest release; set CLUE_VERSION=<x.y.z> and retry.' }
    $version = ("$final" -split '/')[-1]
}
# Accept 0.7.0 or v0.7.0; the asset names carry bare semver (ADR-011).
$version = $version -replace '^v', ''

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

# Add to the user's PATH, never the machine's, so no elevation is needed.
#
# The raw registry value is read and rewritten deliberately:
# [Environment]::GetEnvironmentVariable expands a REG_EXPAND_SZ value on
# read, so the obvious read-modify-write would silently replace entries
# like %JAVA_HOME%\bin with whatever they happened to resolve to, and
# store the result as a plain string. That quietly breaks a PATH the user
# wrote to track their variables — on a machine this script was only asked
# to install one binary on.
$envKey = 'HKCU:\Environment'
$raw = (Get-Item $envKey).GetValue('Path', '', 'DoNotExpandEnvironmentNames')
$entries = @($raw -split ';' | Where-Object { $_ -ne '' })

if ($entries -notcontains $installDir) {
    $updated = (@($entries) + $installDir) -join ';'
    Set-ItemProperty -Path $envKey -Name Path -Value $updated -Type ExpandString
    Write-Host "Added $installDir to your user PATH."
    Write-Host 'Open a new terminal, then run `clue version` to confirm.'
} else {
    Write-Host 'Run `clue version` to confirm, then `clue init` in a repository.'
}

# The registry write above does not reach this session or any other shell
# already running; update this one so `clue` works immediately here too.
if (($env:Path -split ';') -notcontains $installDir) {
    $env:Path = "$env:Path;$installDir"
}
