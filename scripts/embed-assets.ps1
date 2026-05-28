$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$assetRoot = Join-Path $root "server/internal/embedstatic/assets"

$webDist = Join-Path $root "web/dist"
$adminDist = Join-Path $root "admin/dist"
$h5Dist = Join-Path $root "app/dist/build/h5"

New-Item -ItemType Directory -Path $assetRoot -Force | Out-Null

function Copy-AssetDir {
  param(
    [string]$Source,
    [string]$TargetName
  )

  if (!(Test-Path $Source)) {
    Write-Warning "Skip missing source: $Source"
    return
  }

  $target = Join-Path $assetRoot $TargetName
  if (Test-Path $target) {
    Remove-Item -LiteralPath $target -Recurse -Force
  }

  New-Item -ItemType Directory -Path $target -Force | Out-Null
  Copy-Item -Path (Join-Path $Source "*") -Destination $target -Recurse -Force
  Write-Host "Embedded: $Source -> $target"
}

Copy-AssetDir -Source $webDist -TargetName "web"
Copy-AssetDir -Source $adminDist -TargetName "admin"
Copy-AssetDir -Source $h5Dist -TargetName "h5"
