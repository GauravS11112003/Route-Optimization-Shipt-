# 🛣️ Real Route Setup Guide

## The Issue: Straight Lines vs. Real Roads

If you're seeing straight lines between points instead of routes that follow actual roads, it's because the application needs to be configured to use a routing service API.

## Quick Fix

### Option 1: Enable Real Routes (Already Done! ✅)
The "Real Routes" toggle now defaults to **ON**. When you click "Optimize Routes", it will attempt to fetch real driving routes.

### Option 2: Set Up OpenRouteService API (Recommended)

To get real road-following routes, you need an API key:

#### Step 1: Get a Free API Key
1. Go to [OpenRouteService Sign Up](https://openrouteservice.org/dev/#/signup)
2. Create a free account
3. Copy your API key

#### Step 2: Create `.env` File
In the `backend` folder, create a file named `.env` with this content:

```env
OPENROUTE_API_KEY=your_actual_api_key_here
```

Replace `your_actual_api_key_here` with the key you got from OpenRouteService.

#### Step 3: Restart the Backend
If the backend is already running:
1. Stop it (Ctrl+C in the terminal)
2. Restart it using:
   ```powershell
   cd backend
   .\run.ps1
   ```

You should see:
```
✓ API Key loaded (120 chars)
```

#### Step 4: Test It
1. Load sample data
2. Make sure "Real Routes" is checked (✓)
3. Click "Optimize Routes"
4. Routes should now follow actual roads!

## How It Works

### With API Key ✅
- Routes follow actual roads
- Accurate driving distances and times
- Multiple waypoints per route
- Solid lines on the map

### Without API Key ⚠️
- Straight lines between points
- Estimated distances (as the crow flies)
- Less accurate time predictions
- Dashed lines on the map

## Troubleshooting

### Still seeing straight lines?
1. Check the backend console for:
   - `⚠ No API Key found in environment`
   - `⚠ OpenRouteService API error`

2. Verify your `.env` file:
   - File is in the `backend` folder
   - File is named exactly `.env` (not `.env.txt`)
   - API key has no extra spaces or quotes

3. Make sure you restarted the backend after creating `.env`

### API Rate Limits
The free OpenRouteService plan allows:
- 2,000 requests per day
- 40 requests per minute

For the sample data (5 shoppers, 20 orders), one optimization uses ~20-30 API calls.

## Visual Indicators

The map shows you which mode you're using:
- **Solid lines** = Real routes from API ✅
- **Dashed lines** = Straight-line fallback ⚠️

## Alternative: Use Without API Key

The app works fine without an API key! It just uses:
- Haversine distance formula for calculations
- Straight lines for visualization
- Estimated driving times based on distance

This is useful for:
- Quick prototyping
- Testing the algorithm
- When exact routes aren't critical

---

**Need help?** Check that:
1. "Real Routes" checkbox is enabled
2. `.env` file exists in `backend` folder with valid API key
3. Backend was restarted after creating `.env`
4. Backend console shows "✓ API Key loaded"

