$ErrorActionPreference = "Stop"

$App = "entrox"
$DefaultInstallBaseUrl = "https://entrox.996icu.wiki"
$InstallBaseUrl = if ($env:ENTROX_INSTALL_BASE_URL) { $env:ENTROX_INSTALL_BASE_URL.TrimEnd("/") } else { $DefaultInstallBaseUrl }
$DownloadBaseUrl = if ($env:ENTROX_DOWNLOAD_BASE_URL) { $env:ENTROX_DOWNLOAD_BASE_URL.TrimEnd("/") } else { "$InstallBaseUrl/downloads/entrox-dev" }
$InstallDir = if ($env:ENTROX_INSTALL_DIR) { $env:ENTROX_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\Entrox\bin" }
$Version = if ($env:VERSION) { $env:VERSION.TrimStart("v") } else { "" }

if (-not [Environment]::Is64BitOperatingSystem) {
  throw "Entrox Windows installer currently supports Windows x64 only."
}

if ($Version) {
  if ($Version -like "0.0.0-ci.*") {
    $DownloadBaseUrl = "$DownloadBaseUrl/$Version"
  } else {
    $DownloadBaseUrl = "https://github.com/hexonal/entrox/releases/download/v$Version"
  }
  $DisplayVersion = "v$Version"
} else {
  $ManifestUrl = "$DownloadBaseUrl/latest.json"
  Write-Host "Resolving latest Entrox release"
  $Manifest = Invoke-RestMethod -UseBasicParsing -Uri $ManifestUrl
  $DisplayVersion = $Manifest.version
}

$Asset = "entrox-cli-windows-x64.zip"
if ($Version) {
  $Url = "$DownloadBaseUrl/$Asset"
  $ExpectedHash = ""
} else {
  $AssetInfo = $Manifest.assets | Where-Object { $_.name -eq $Asset } | Select-Object -First 1
  if (-not $AssetInfo -or -not $AssetInfo.url) {
    throw "Latest Entrox manifest did not include $Asset."
  }
  $Url = $AssetInfo.url
  $ExpectedHash = $AssetInfo.sha256
}
$TempDir = Join-Path ([IO.Path]::GetTempPath()) ("entrox_install_" + [Guid]::NewGuid().ToString("N"))
$Archive = Join-Path $TempDir $Asset

New-Item -ItemType Directory -Force -Path $TempDir, $InstallDir | Out-Null

try {
  Write-Host "Installing Entrox $DisplayVersion for windows-x64"
  Write-Host "Downloading $Asset"
  Invoke-WebRequest -UseBasicParsing -Uri $Url -OutFile $Archive
  if ($ExpectedHash) {
    $ActualHash = (Get-FileHash -Algorithm SHA256 -Path $Archive).Hash.ToLowerInvariant()
    if ($ActualHash -ne $ExpectedHash.ToLowerInvariant()) {
      throw "SHA-256 verification failed for $Asset. Expected $ExpectedHash, actual $ActualHash."
    }
  }
  Expand-Archive -Path $Archive -DestinationPath $TempDir -Force

  $Candidates = @(
    (Join-Path $TempDir "bin\entrox.exe"),
    (Join-Path $TempDir "bin\entrox"),
    (Join-Path $TempDir "entrox.exe"),
    (Join-Path $TempDir "entrox")
  )
  $Binary = $Candidates | Where-Object { Test-Path $_ } | Select-Object -First 1
  if (-not $Binary) {
    throw "Downloaded archive did not contain the Entrox binary."
  }

  $Target = Join-Path $InstallDir "entrox.exe"
  Copy-Item -Path $Binary -Destination $Target -Force

  if (-not $env:NO_MODIFY_PATH) {
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $PathParts = @()
    if ($UserPath) {
      $PathParts = $UserPath -split ";" | Where-Object { $_ }
    }
    if ($PathParts -notcontains $InstallDir) {
      $NewPath = if ($UserPath) { "$UserPath;$InstallDir" } else { $InstallDir }
      [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
      Write-Host "Added Entrox to the user PATH. Open a new PowerShell window before running entrox."
    }
  }

  Write-Host ""
  Write-Host "Entrox installed to $Target"
  Write-Host "Run: entrox login"
} finally {
  Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
}
