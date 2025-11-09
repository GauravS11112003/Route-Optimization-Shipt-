# 🔍 Debugging Real Routes Issue

## Problem
Routes are not following roads even when the API is being called.

## Debugging Steps

I've added comprehensive logging to both frontend and backend. Follow these steps to identify the issue:

### 1. Restart Backend
```powershell
cd backend
.\run.ps1
```

### 2. Restart Frontend (if not auto-reloaded)
```bash
cd frontend
npm run dev
```

### 3. Open Browser Console
- Open the app at http://localhost:5173
- Open Browser DevTools (F12)
- Go to Console tab

### 4. Test Real Routes
1. Click "Load Sample Data"
2. Make sure "Real Routes ✓" checkbox is enabled
3. Click "Optimize Routes"

### 5. Check Frontend Console Logs

You should see logs like:
```
🔍 API Response: {optimization: {...}, analytics: {...}}
📍 Route Geometries: [{shopperId: "S1", points: [[...], ...]}, ...]
✅ Use Real Routes: true
🗺️ Setting geometries: 5 routes
📊 First route has 150 points
🗺️ MapView - routeGeometries: 5
🗺️ MapView - assignments: 5
Route 0: 150 points for shopper S1
Route 1: 200 points for shopper S2
...
```

### 6. Check Backend Console Logs

You should see logs like:
```
🗺️  calculateRouteGeometries called with useRealRoutes: true
📊 Assignments count: 5
  Route 0 for shopper S1 has 5 waypoints
    🛣️  Fetching real routes...
      Segment 0 returned 45 points
      Segment 1 returned 38 points
      Segment 2 returned 42 points
      Segment 3 returned 25 points
    Total points for route: 150
...
```

## Common Issues & Solutions

### Issue 1: routeGeometries is empty or has only a few points

**Check backend logs for:**
```
⚠️ Segment X failed, using fallback
⚠ No API Key found in environment
⚠ OpenRouteService API error
```

**Solution**: Set up API key
```powershell
cd backend
.\setup-env.ps1
# Enter your OpenRouteService API key
```

### Issue 2: routeGeometries has many points but lines are still straight

**Check if points are in correct format:**
Frontend console should show:
```javascript
points: [[33.5186, -86.8104], [33.5187, -86.8103], ...]
        [  lat  ,   lng  ]     [  lat  ,   lng  ]
```

If you see `[[lng, lat]]` instead, coordinates are swapped!

**Solution**: Already fixed in `routing.go` lines 123-127

### Issue 3: useRealRoutes is false even though checkbox is checked

**Check frontend console:**
```
✅ Use Real Routes: false  ← Should be true!
```

**Solution**: The state isn't being passed correctly. Check App.jsx line 52.

### Issue 4: API key is set but routes still fail

**Check backend logs for HTTP status:**
```
⚠ OpenRouteService API error - Status: 403
   Response body: {"error": "Invalid API key"}
```

**Solutions:**
- Verify API key is correct (120 characters)
- Check if you hit rate limits (2000 requests/day)
- Try the test endpoint: http://localhost:8080/api/test-routing

### Issue 5: Only 2 points per route segment

**Backend logs show:**
```
Segment 0 returned 2 points  ← Should be 40-60+ points for real roads
```

This means the API is falling back to straight lines.

**Check:**
1. API key validity
2. Network connectivity
3. OpenRouteService service status

## Expected Behavior

### With Real Routes ENABLED and API Key SET:
- Backend: "Segment X returned 40-60 points" (varies by distance)
- Frontend: Each route has 150-300+ points total
- Map: Routes follow curves of roads (solid lines)

### With Real Routes DISABLED or NO API Key:
- Backend: "Using straight lines" or "failed, using fallback"
- Frontend: Each route has only 5-10 points (waypoints only)
- Map: Straight dashed lines between markers

## Quick Test Commands

### Test Backend Routing Endpoint
```powershell
# Open in browser or curl
http://localhost:8080/api/test-routing
```

Should return:
```json
{
  "pointCount": 45,
  "distance": 1.2,
  "duration": 3.5,
  "usingFallback": false,
  "apiKeySet": true
}
```

If `usingFallback: true`, API is not working.

### Check Environment Variable (Backend Terminal)
```powershell
$env:OPENROUTE_API_KEY
```

Should output your 120-character key. If blank, .env wasn't loaded.

## Manual Verification

### 1. Check if geometry data exists
In browser console after optimizing:
```javascript
// Paste this in console
console.log('Geometries:', JSON.stringify(window.lastGeometries))
```

### 2. Check coordinate format
```javascript
// Should be [lat, lng] not [lng, lat]
console.log(window.lastGeometries[0].points[0])
// Expected: [33.5186, -86.8104]  ← Lat first
```

### 3. Check Polyline component
Look at the map, if you see:
- **Dashed lines**: Using fallback (2-5 points per route)
- **Solid lines**: Real routes (many points per route)

## Advanced Debugging

### Add breakpoint in MapView
Add this in `MapView.jsx` after line 68:
```javascript
if (routeGeometries.length > 0) {
    debugger; // Pause here
    console.log('Geometry points:', routeGeometries[0].points);
}
```

### Verify Leaflet is receiving correct data
```javascript
// In browser console
document.querySelectorAll('.leaflet-interactive').forEach(el => {
    console.log('Path element:', el.getAttribute('d'));
});
```

If you see only "M33.5 -86.8 L33.6 -86.7" (few points), it's straight lines.
If you see hundreds of coordinates, real routes are rendering.

## Resolution Checklist

- [ ] Backend console shows "✓ API Key loaded"
- [ ] Frontend console shows "✅ Use Real Routes: true"
- [ ] Backend shows "🛣️ Fetching real routes..."
- [ ] Segments return 40+ points each
- [ ] Frontend receives geometries with 150+ points per route
- [ ] MapView receives non-empty routeGeometries array
- [ ] Routes on map are solid lines (not dashed)
- [ ] Routes curve and follow roads visually

## If All Else Fails

### Nuclear Option: Start Fresh

1. Stop backend (Ctrl+C)
2. Delete `backend/.env`
3. Run setup:
```powershell
cd backend
.\setup-env.ps1
```
4. Enter fresh API key from https://openrouteservice.org
5. Restart backend: `.\run.ps1`
6. Hard refresh browser (Ctrl+Shift+R)
7. Try optimization again

## Contact Info

If issue persists after all debugging:
1. Copy all console logs (frontend + backend)
2. Check backend terminal for any errors
3. Screenshot the map showing straight lines
4. Note which debugging steps showed unexpected results

---

**Next Steps**: Run the app, check logs, and report what you see in each section above.

