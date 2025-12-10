# CleanGo Installation Script for Windows (PowerShell)
# Version 1.0

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║         CleanGo Installation Script v1.0              ║" -ForegroundColor Blue
Write-Host "╚═══════════════════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

# Check if Go is installed
$goCommand = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCommand) {
    Write-Host "❌ Go is not installed. Please install Go 1.20 or higher." -ForegroundColor Red
    Write-Host "Visit: https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

$goVersion = (go version) -replace 'go version go', '' -replace ' .*', ''
Write-Host "✓ Go version: $goVersion" -ForegroundColor Green

# Get GOPATH
$GOPATH = go env GOPATH
$GOBIN = Join-Path $GOPATH "bin"

Write-Host "📦 Installing cleango..." -ForegroundColor Blue
Write-Host ""

# Install using go install
try {
    go install github.com/YeridStick/cleango/cmd/cleango@latest

    Write-Host ""
    Write-Host "✅ CleanGo installed successfully!" -ForegroundColor Green
    Write-Host ""

    # Check if GOBIN is in PATH
    $envPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not ($envPath -like "*$GOBIN*")) {
        Write-Host "⚠️  Warning: $GOBIN is not in your PATH" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Adding $GOBIN to your PATH..." -ForegroundColor Yellow

        # Add to PATH for current session
        $env:Path += ";$GOBIN"

        # Add to PATH permanently for user
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$envPath;$GOBIN",
            "User"
        )

        Write-Host "✓ Added $GOBIN to PATH" -ForegroundColor Green
        Write-Host "⚠️  Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
        Write-Host ""
    }

    # Verify installation
    $cleangoCommand = Get-Command cleango -ErrorAction SilentlyContinue
    if ($cleangoCommand) {
        $cleangoVersion = (cleango --version 2>&1) -as [string]
        Write-Host "✓ cleango is now available" -ForegroundColor Green
        Write-Host "✓ Location: $($cleangoCommand.Source)" -ForegroundColor Green
        if ($cleangoVersion) {
            Write-Host "✓ Version: $cleangoVersion" -ForegroundColor Green
        }
    } else {
        Write-Host "⚠️  cleango command not found. Please restart your terminal." -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Blue
    Write-Host "🚀 Quick Start:" -ForegroundColor Green
    Write-Host ""
    Write-Host "  1. Create a new project:" -ForegroundColor Yellow
    Write-Host "     cleango new my-api" -ForegroundColor Green
    Write-Host ""
    Write-Host "  2. Navigate to your project:" -ForegroundColor Yellow
    Write-Host "     cd my-api" -ForegroundColor Green
    Write-Host ""
    Write-Host "  3. Run the application:" -ForegroundColor Yellow
    Write-Host "     go run ./cmd/api" -ForegroundColor Green
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Blue
    Write-Host ""
    Write-Host "For more information, visit: https://github.com/YeridStick/cleango" -ForegroundColor Blue
    Write-Host ""

} catch {
    Write-Host ""
    Write-Host "❌ Installation failed: $_" -ForegroundColor Red
    exit 1
}
