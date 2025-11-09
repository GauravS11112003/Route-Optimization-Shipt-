# 🚨 URGENT: Set Up API Key for Real Roads

## The Problem
Your routes are showing straight lines because the OpenRouteService API is not configured. Without an API key, the system falls back to simple point-to-point connections.

## Quick Fix (5 minutes)

### Step 1: Get FREE API Key

1. **Go to:** https://openrouteservice.org/dev/#/signup
2. **Sign up** (it's free, no credit card needed)
3. **Copy your API key** (it's a long string, ~120 characters)

### Step 2: Set Up the Key

**Option A: Use the Setup Script (Easiest)**
```powershell
cd backend
.\setup-env.ps1
```
When prompted, paste your API key.

**Option B: Manual Setup**
1. Create a file named `.env` in the `backend` folder
2. Add this line:
```
OPENROUTE_API_KEY=your_actual_api_key_here
```
3. Replace `your_actual_api_key_here` with your real key

### Step 3: Restart Backend

**Stop the backend** (press Ctrl+C in the backend terminal)

**Start it again:**
```powershell
cd backend
.\run.ps1
```

**Look for this line:**
```
✓ API Key loaded (120 chars)
```

If you see `⚠ WARNING: OPENROUTE_API_KEY not found`, the .env file wasn't loaded properly.

### Step 4: Test It!

1. **Refresh the browser** (Ctrl+Shift+R)
2. **Load Sample Data**
3. **Make sure "Real Routes ✓" is checked**
4. **Click "Optimize Routes"**

Now the routes should curve and follow actual roads! 🎉

## Verify It's Working

### Backend Console Should Show:
```
✓ API Key loaded: 5b0ce3591b...
🗺️ calculateRouteGeometries called with useRealRoutes: true
  🛣️ Fetching real routes...
    Segment 0 returned 45 points  ← GOOD! Many points
    Segment 1 returned 38 points  ← GOOD! Many points
  Total points for route: 150+   ← GOOD! Lots of points
```

### Frontend Should Show:
- **Solid, curved lines** following roads
- **No warning banner** about fallback
- Browser console: `✅ Using real road geometries!`

### If It Still Shows Straight Lines:

Check backend console for:
```
⚠ OpenRouteService API error - Status: 403
```

This means:
- API key is invalid/expired
- Or you hit the rate limit (2000 requests/day)

**Solution:** Get a fresh API key from OpenRouteService

## What Each Line Type Means

| Line Style | Points | Meaning |
|------------|--------|---------|
| 🔵 **Solid curved lines** | 100-300 | ✅ Real roads (API working) |
| 🔵 **Dashed straight lines** | 5-10 | ❌ Fallback (API not working) |

## Still Having Issues?

### Test the API Directly

1. Open in browser: `http://localhost:8080/api/test-routing`

**Good Response:**
```json
{
  "pointCount": 45,
  "usingFallback": false,
  "apiKeySet": true
}
```

**Bad Response:**
```json
{
  "pointCount": 2,
  "usingFallback": true,
  "apiKeySet": false
}
```

### Common Mistakes

❌ **File named `.env.txt`** instead of `.env`
- Windows hides extensions by default
- Make sure it's exactly `.env` with no extension

❌ **API key has extra spaces or quotes**
- Should be: `OPENROUTE_API_KEY=5b0ce3591b64...`
- NOT: `OPENROUTE_API_KEY="5b0ce3591b64..."`
- NOT: `OPENROUTE_API_KEY= 5b0ce3591b64...` (space after =)

❌ **Didn't restart backend after creating .env**
- Must restart for environment variables to load

❌ **Wrong directory**
- The `.env` file must be in the `backend` folder
- NOT in the root folder

## Rate Limits

Free OpenRouteService plan:
- **2,000 requests per day**
- **40 requests per minute**

With 20 orders and 5 shoppers:
- Each optimization = ~20-25 API calls
- You can optimize ~80-100 times per day

## After Setup

Once working, you'll see:
- Routes that curve and bend
- Routes that follow highways and streets
- Much more realistic-looking paths
- Better distance calculations
- More accurate time estimates

---

**Don't Skip This!** The app works without it, but routes will be straight lines. 
Real routing is the main feature - take 5 minutes to set it up! 🚀

