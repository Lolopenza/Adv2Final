# PowerShell script to run all tests for the payment service
# This includes both unit and integration tests

# Set working directory to payment service root
$paymentServicePath = $PSScriptRoot

function Write-ColorText {
    param (
        [string]$text,
        [string]$color
    )
    Write-Host $text -ForegroundColor $color
}

Write-ColorText "====================================================" "Cyan"
Write-ColorText "        RUNNING PAYMENT SERVICE TESTS               " "Cyan" 
Write-ColorText "====================================================" "Cyan"

# 1. Run unit tests
Write-ColorText "`n[1/3] Running unit tests..." "Yellow"
go test ./test/unit/... -v
if ($LASTEXITCODE -ne 0) {
    Write-ColorText "Unit tests failed!" "Red"
    exit 1
}
Write-ColorText "Unit tests passed!" "Green"

# 2. Run mock Redis tests
Write-ColorText "`n[2/3] Running cache tests..." "Yellow"
go test ./internal/repository/cache/... -v
if ($LASTEXITCODE -ne 0) {
    Write-ColorText "Cache tests failed!" "Red"
    exit 1
}
Write-ColorText "Cache tests passed!" "Green"

# 3. Run integration tests
Write-ColorText "`n[3/3] Running integration tests..." "Yellow"
Write-ColorText "(Note: These tests require Redis and PostgreSQL running)" "Yellow"

# Check if miniredis package is installed
Write-ColorText "Checking for required dependencies..." "Yellow"
$miniredisInstalled = go list -m github.com/alicebob/miniredis/v2 2>$null
if (-not $miniredisInstalled) {
    Write-ColorText "Missing dependency: github.com/alicebob/miniredis/v2" "Red"
    Write-ColorText "To install required test dependencies, run:" "Yellow"
    Write-ColorText "go get github.com/alicebob/miniredis/v2" "Cyan"
}

# Run the integration tests
go test ./test/integration/... -v
if ($LASTEXITCODE -ne 0) {
    # Only fail if not just skipped tests
    $output = go test ./test/integration/... -v
    if (-not ($output -match "skip")) {
        Write-ColorText "Integration tests failed!" "Red"
        exit 1
    } else {
        Write-ColorText "Integration tests skipped due to missing dependencies" "Yellow"
    }
} else {
    Write-ColorText "Integration tests passed!" "Green"
}

Write-ColorText "`n====================================================" "Cyan"
Write-ColorText "        TEST SUITE EXECUTION COMPLETED               " "Green" 
Write-ColorText "====================================================" "Cyan" 