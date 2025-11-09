# Shipt Route Optimizer - Backend Startup Script
# This script sets the environment variable and starts the backend

Write-Host "Starting Shipt Route Optimizer Backend..." -ForegroundColor Green
Write-Host ""

# Load .env file and set environment variables
if (Test-Path ".env") {
    Write-Host "Loading .env file..." -ForegroundColor Cyan
    Get-Content ".env" | ForEach-Object {
        if ($_ -match "^([^#][^=]+)=(.+)$") {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
            if ($key -eq "OPENROUTE_API_KEY") {
                $keyPreview = $value.Substring(0, [Math]::Min(20, $value.Length))
                Write-Host "API Key loaded: $keyPreview... ($($value.Length) chars)" -ForegroundColor Green
            }
        }
    }
    Write-Host ""
} else {
    Write-Host "WARNING: .env file not found" -ForegroundColor Yellow
    Write-Host ""
}

# Start the backend
Write-Host "Starting backend on port 8080..." -ForegroundColor Cyan
& "C:\Program Files\Go\bin\go.exe" run cmd/main.go

