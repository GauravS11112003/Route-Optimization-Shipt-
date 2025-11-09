# Setup script for OpenRouteService API key
# This creates a .env file for real route functionality

Write-Host "🛣️  Route Optimizer - Environment Setup" -ForegroundColor Cyan
Write-Host ""

# Check if .env already exists
if (Test-Path ".env") {
    Write-Host "⚠️  .env file already exists!" -ForegroundColor Yellow
    $overwrite = Read-Host "Do you want to overwrite it? (y/N)"
    if ($overwrite -ne "y" -and $overwrite -ne "Y") {
        Write-Host "Setup cancelled. Your existing .env file was not modified." -ForegroundColor Gray
        exit
    }
}

Write-Host ""
Write-Host "To enable real driving routes, you need an OpenRouteService API key." -ForegroundColor White
Write-Host ""
Write-Host "Steps to get your FREE API key:" -ForegroundColor Cyan
Write-Host "  1. Visit: https://openrouteservice.org/dev/#/signup" -ForegroundColor Gray
Write-Host "  2. Sign up for a free account" -ForegroundColor Gray
Write-Host "  3. Copy your API key" -ForegroundColor Gray
Write-Host ""

$apiKey = Read-Host "Enter your OpenRouteService API key (or press Enter to skip)"

if ([string]::IsNullOrWhiteSpace($apiKey)) {
    Write-Host ""
    Write-Host "ℹ️  No API key provided. Creating .env file with placeholder..." -ForegroundColor Yellow
    $apiKey = "your_api_key_here"
    $needsKey = $true
} else {
    Write-Host ""
    Write-Host "✓ API key received!" -ForegroundColor Green
    $needsKey = $false
}

# Create .env file
$envContent = @"
# OpenRouteService API Key
# Get a free API key at: https://openrouteservice.org/dev/#/signup
# This is required for real driving routes. Without it, the app will use straight-line fallback routes.
OPENROUTE_API_KEY=$apiKey
"@

try {
    $envContent | Out-File -FilePath ".env" -Encoding UTF8 -NoNewline
    Write-Host "✓ .env file created successfully!" -ForegroundColor Green
    
    if ($needsKey) {
        Write-Host ""
        Write-Host "⚠️  Don't forget to:" -ForegroundColor Yellow
        Write-Host "  1. Get your API key from https://openrouteservice.org" -ForegroundColor Gray
        Write-Host "  2. Edit the .env file and replace 'your_api_key_here' with your actual key" -ForegroundColor Gray
        Write-Host "  3. Restart the backend server" -ForegroundColor Gray
    } else {
        Write-Host ""
        Write-Host "🚀 All set! Start the backend with: .\run.ps1" -ForegroundColor Green
    }
    
} catch {
    Write-Host "❌ Error creating .env file: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "For detailed instructions, see: ROUTING_SETUP.md" -ForegroundColor Gray
Write-Host ""

