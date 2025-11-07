# Shipt Route Optimizer - Backend Startup Script
# This script sets the environment variable and starts the backend

Write-Host "🚀 Starting Shipt Route Optimizer Backend..." -ForegroundColor Green

# Load .env file and set environment variables
if (Test-Path ".env") {
    Get-Content ".env" | ForEach-Object {
        if ($_ -match "^([^#][^=]+)=(.+)$") {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
            if ($key -eq "OPENROUTE_API_KEY") {
                Write-Host "✓ API Key loaded ($($value.Length) chars)" -ForegroundColor Green
            }
        }
    }
} else {
    Write-Host "⚠ Warning: .env file not found" -ForegroundColor Yellow
}

# Start the backend
& "C:\Program Files\Go\bin\go.exe" run cmd/main.go

