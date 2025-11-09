# 🔍 Diagnose Routing Issue - Straight Lines Problem

## The Problem
Routes are showing straight lines even after adding API key.

## Diagnostic Steps

### Step 1: Stop Current Backend
Press `Ctrl+C` in the backend terminal window

### Step 2: Start Debug Backend
```powershell
cd "C:\Route Optimizer\backend"
.\backend-debug.exe
```

### Step 3: Check Browser Console
1. Open browser DevTools (F12)
2. Go to Console tab
3. Click "Load Sample Data"
4. Add your API key in Settings (if not already done)
5. Click "Optimize Routes"

### Step 4: Look for These Logs

#### In Browser Console - Look for:
```javascript
📊 First route has X points
```

**If X is 5-10:** ❌ Fallback mode (straight lines)
**If X is 100+:** ✅ Real routes are working

Also check:
```javascript
📍 First geometry sample: {
  pointsCount: X,
  firstPoint: [lat, lng],
  allPoints: [[lat1, lng1], [lat2, lng2], ...]
}
```

#### In Backend Console - Look for:
```
📡 Received API key from frontend: X chars
```

**If X is 0:** ❌ API key not being sent
**If X is 600+:** ✅ API key is being sent

Then look for:
```
✓ Using API Key from frontend: eyJv... (X chars)
✓ OpenRouteService API success - Status: 200
🗺️ Coordinate count: X
```

**If Status is 403:** ❌ Invalid API key
**If Status is 429:** ❌ Rate limit exceeded
**If Coordinate count is 2:** ❌ Fallback mode
**If Coordinate count is 40+:** ✅ API is working

## Common Issues & Solutions

### Issue 1: API Key Not Being Sent

**Backend shows:**
```
📡 Received API key from frontend: 0 chars
⚠️ No API Key provided from frontend
```

**Solution:**
1. Click Settings button
2. Make sure API key is pasted
3. Click "Save API Key"
4. Button should turn green
5. Try optimizing again

### Issue 2: Invalid API Key

**Backend shows:**
```
⚠ OpenRouteService API error - Status: 403
Response body: {"error": "Invalid API key"}
```

**Solution:**
1. Go to https://openrouteservice.org/dev/#/signup
2. Get a fresh API key
3. Click Settings → Clear Key
4. Paste new key → Save
5. Try again

### Issue 3: Rate Limited

**Backend shows:**
```
⚠ OpenRouteService API error - Status: 429
Response body: {"error": "Rate limit exceeded"}
```

**Solution:**
- Wait 1-2 minutes
- Free tier: 40 requests/minute, 2000/day
- Or get another API key

### Issue 4: Network/Firewall Issue

**Backend shows:**
```
⚠ OpenRouteService HTTP request failed: dial tcp...
```

**Solution:**
- Check internet connection
- Check firewall settings
- Try from different network

### Issue 5: Coordinates Swapped

**Browser console shows points but map shows straight lines**

Check if coordinates are like:
```
firstPoint: [33.5186, -86.8104]  ← Should be lat first
```

If showing `[-86.8104, 33.5186]` (lng first), coordinates are swapped!

## Test API Key Manually

### Open PowerShell and run:
```powershell
$headers = @{
    "Authorization" = "YOUR_API_KEY_HERE"
    "Content-Type" = "application/json"
}

$body = @{
    coordinates = @(
        @(-86.8104, 33.5186),
        @(-86.8050, 33.5250)
    )
} | ConvertTo-Json

Invoke-RestMethod -Uri "https://api.openrouteservice.org/v2/directions/driving-car" -Method POST -Headers $headers -Body $body
```

**Good response:**
```json
{
  "routes": [{
    "summary": { "distance": 1234.5, "duration": 123.4 },
    "geometry": { "coordinates": [[...], [...], ...] }
  }]
}
```

**Bad response:**
```json
{
  "error": "Invalid API key"
}
```

## Expected vs Actual

### Expected Behavior:

**Backend:**
```
📡 Received API key from frontend: 656 chars
✓ Using API Key from frontend: eyJvcmci... (656 chars)
✓ OpenRouteService API success - Status: 200
📦 Response sample: {"routes":[{"summary":{"distance":1234.5...
📏 Route distance: 1234.5 meters
⏱️ Route duration: 123.4 seconds
🗺️ Coordinate count: 45
   First coord (raw): [-86.8104 33.5186]
   Converted point 0: [lng: -86.8104, lat: 33.5186] -> RoutePoint{Lat: 33.5186, Lng: -86.8104}
✅ Final geometry has 45 points
    Total points for route: 150+
```

**Frontend:**
```
📊 First route has 150 points
✅ Using real road geometries!
📍 First geometry sample: {
  pointsCount: 150,
  firstPoint: [33.5186, -86.8104],
  allPoints: [[33.5186, -86.8104], [33.5187, -86.8103], ...]
}
Route 0: 150 points for shopper S1
  First point: [33.5186, -86.8104]
  Last point: [33.5250, -86.8050]
```

**Map:**
- Solid curved lines following roads ✅

### Actual (Current Issue):

**Likely Backend:**
```
📡 Received API key from frontend: 656 chars
✓ Using API Key from frontend: eyJvcmci... (656 chars)
⚠ OpenRouteService API error - Status: 403 or 429
⚠️ Segment 0 failed, using fallback
    Segment 0 returned 2 points  ← Only 2!
    Total points for route: 10
```

**Likely Frontend:**
```
📊 First route has 5 points
⚠️ Real routes unavailable - API key may be invalid
Route 0: 5 points for shopper S1
```

**Map:**
- Dashed straight lines ❌

## Next Steps

1. **Run the debug backend** (`backend-debug.exe`)
2. **Open browser console** (F12)
3. **Optimize routes**
4. **Copy all console output** from both backend and frontend
5. **Share the logs** to identify the exact issue

## Quick Checklist

- [ ] Backend debug version is running
- [ ] Browser console is open
- [ ] Settings button shows green "API Key ✓"
- [ ] "Real Routes ✓" is checked
- [ ] Clicked "Optimize Routes"
- [ ] Checked backend console for errors
- [ ] Checked browser console for point counts
- [ ] API key is valid and not rate-limited

---

**Most likely cause:** API key is invalid or expired. Get a fresh one from OpenRouteService!

