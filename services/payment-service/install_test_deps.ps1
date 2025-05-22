# PowerShell script to install test dependencies
# This script installs the required dependencies for running tests

function Write-ColorText {
    param (
        [string]$text,
        [string]$color
    )
    Write-Host $text -ForegroundColor $color
}

Write-ColorText "====================================================" "Cyan"
Write-ColorText "        INSTALLING TEST DEPENDENCIES                " "Cyan" 
Write-ColorText "====================================================" "Cyan"

# Install miniredis for integration tests
Write-ColorText "`nInstalling miniredis for Redis mocking in tests..." "Yellow"
go get github.com/alicebob/miniredis/v2

# Install testify for assertions and mocking
Write-ColorText "`nInstalling testify for test assertions and mocks..." "Yellow"
go get github.com/stretchr/testify

# Install go-sqlmock for database mocking
Write-ColorText "`nInstalling go-sqlmock for database mocking..." "Yellow"
go get github.com/DATA-DOG/go-sqlmock

Write-ColorText "`nAll test dependencies installed!" "Green"
Write-ColorText "You can now run tests with: ./run_tests.ps1" "Green"
Write-ColorText "====================================================" "Cyan" 