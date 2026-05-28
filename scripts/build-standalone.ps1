param(
  [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$dist = Join-Path $root $OutputDir
New-Item -ItemType Directory -Path $dist -Force | Out-Null

Write-Host "Building web frontend..."
npm --prefix (Join-Path $root "web") run build

Write-Host "Building admin frontend (base=/admin/)..."
npm --prefix (Join-Path $root "admin") run build -- --base=/admin/

Write-Host "Building app H5 frontend for embedded mode..."
npm --prefix (Join-Path $root "app") run build:h5:embed

Write-Host "Embedding frontend assets into Go server..."
pwsh -File (Join-Path $root "scripts/embed-assets.ps1")

Push-Location (Join-Path $root "server")
try {
  Write-Host "Building lyshop.exe..."
  $env:GOOS = "windows"
  $env:GOARCH = "amd64"
  go build -o (Join-Path $dist "lyshop.exe") .

  Write-Host "Building lyshop-gui.exe..."
  go build -ldflags "-H=windowsgui" -o (Join-Path $dist "lyshop-gui.exe") ./cmd/lyshop-gui
}
finally {
  Remove-Item Env:GOOS -ErrorAction SilentlyContinue
  Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
  Pop-Location
}

Write-Host "Build complete. Artifacts: $dist"
