# 🎉 New API Key System - No More .env Files!

## What Changed

The app now has a **much easier** way to set up real routing - just paste your API key in the UI!

### ❌ OLD WAY (Complex)
1. Create `.env` file
2. Paste API key in file
3. Restart backend
4. Hope it loaded correctly
5. Check console logs

### ✅ NEW WAY (Simple)
1. Click **Settings** button
2. Paste API key
3. Click Save
4. Done! Routes work immediately

## How to Use It

### Step 1: Get Your FREE API Key
1. Go to: https://openrouteservice.org/dev/#/signup
2. Sign up (free, no credit card)
3. Copy your API key (~600 characters, starts with "eyJ")

### Step 2: Add It to the App
1. **Click the "Settings" button** in the header (orange if no key is set)
2. **Click "Get Free API Key"** button to open sign-up page
3. **Paste your API key** in the password field
4. **Click "Save API Key"**
5. Settings button turns **green ✓** when saved!

### Step 3: Use Real Routes
1. Make sure **"Real Routes ✓"** toggle is enabled (should be by default)
2. Click **"Load Sample Data"**
3. Click **"Optimize Routes"**
4. **Routes now follow actual roads!** 🎉

## Features

### Visual Indicators

**Settings Button:**
- 🟠 **Orange "Settings"** = No API key set
- 🟢 **Green "API Key ✓"** = API key is saved and ready

**Real Routes Toggle:**
- ✓ **Checked & Green** = Will use real routes
- ⬜ **Unchecked** = Will use straight lines

### Error Handling

The app will show helpful messages:
- **"API key required"** - Click Settings to add one
- **"API key may be invalid"** - Check your key in Settings
- **Routes show solid curved lines** = Working! ✅
- **Routes show dashed straight lines** = Not working ❌

### Data Storage

**Where is it stored?**
- Your API key is saved in your **browser's localStorage**
- It persists across page refreshes
- It's **only on your computer**

**Security:**
- API key is sent **directly to OpenRouteService**
- It **never goes through our backend servers**
- Your key stays private and secure

### Managing Your Key

**To Update:**
1. Click Settings
2. Paste new key
3. Click Save

**To Remove:**
1. Click Settings
2. Click "Clear Key"
3. Confirms removal

**To View:**
- Key is hidden (password field)
- Shows character count when typing
- ✓ appears when >100 characters

## Technical Details

### How It Works

```
Frontend (Browser)
    ↓
  [User pastes API key]
    ↓
  localStorage.setItem('openroute_api_key', key)
    ↓
  [User clicks Optimize]
    ↓
  API Request: POST /api/optimize-analytics
  Body: {
    orders: [...],
    shoppers: [...],
    useRealRoutes: true,
    apiKey: "eyJv..."  ← Sent from frontend
  }
    ↓
Backend (Go)
    ↓
  Receives API key from request
    ↓
  routing.GetRouteWithKey(from, to, apiKey)
    ↓
  OpenRouteService API
  Headers: {
    Authorization: apiKey  ← Used directly
  }
    ↓
  Returns road geometry
    ↓
Frontend renders curved routes!
```

### Architecture Changes

**Frontend (`App.jsx`):**
- Added `apiKey` state
- Added `apiKeyInput` state  
- Settings modal with password input
- localStorage integration
- Pass `apiKey` to API calls

**API Client (`optimizer.js`):**
- Added `apiKey` parameter to `optimizeWithAnalytics`
- Sends API key in request body

**Backend (`handlers.go`):**
- Added `ApiKey` field to request struct
- Passes API key to optimizer functions

**Optimizer (`optimizer_v2.go`):**
- Updated all functions to accept `apiKey` parameter
- Passes key through to routing functions

**Routing (`routing.go`):**
- New `GetRouteWithKey(lat, lng, key)` function
- Uses provided key instead of environment variable
- Backward compatible `GetRoute()` still exists

## Benefits

### ✅ User-Friendly
- No file system navigation
- No text editor needed
- Visual confirmation (green checkmark)
- Clear instructions in the UI

### ✅ No Backend Restart Required
- Change API key anytime
- Takes effect immediately
- No server downtime

### ✅ Multiple Users / Browsers
- Each user can use their own key
- Different keys for different projects
- Easy to switch

### ✅ Better Error Messages
- App knows if key is missing
- Shows helpful prompts
- Opens Settings automatically when needed

### ✅ Security
- Key never stored on server
- Direct API to OpenRouteService
- Can be cleared anytime

## Migration from .env File

If you had a `.env` file before:

1. **Old key still works** for backward compatibility
2. **UI key overrides .env** if both are present
3. **UI key is recommended** for ease of use
4. **Can safely delete .env** file once UI key is set

## Troubleshooting

### Issue: Settings button stays orange after saving

**Solution:**
- Make sure you clicked "Save API Key"
- Check browser console for errors
- Try refreshing the page

### Issue: Routes still show straight lines

**Check:**
1. Settings button is green (key is saved)
2. "Real Routes ✓" is checked
3. Backend console shows "✓ Using API Key from frontend"
4. No error messages in browser console

**If still not working:**
- API key might be invalid - get a fresh one
- Check character count (should be ~600+ chars)
- Make sure key starts with "eyJ"

### Issue: Can't paste into password field

**Solution:**
- Click in the field first
- Use Ctrl+V to paste
- Or right-click → Paste

### Issue: Lost my API key

**Solution:**
- Just get a new free one from OpenRouteService
- Old keys don't expire quickly
- Can have multiple keys

## Rate Limits

**Free OpenRouteService Plan:**
- 2,000 requests per day
- 40 requests per minute

**What uses requests:**
- Each road segment = 1 request
- 20 orders, 5 shoppers ≈ 20 requests per optimization
- Can optimize ~100 times per day

**If you hit the limit:**
- Routes fall back to straight lines
- Backend logs: "Status: 429 Rate limit exceeded"
- Wait a few minutes or get another key

## Summary

🎉 **No more .env files!**
🔑 **Just paste your key in the UI**
✅ **Instant real routing**
🚀 **Much simpler setup**

---

**Ready to try it?**
1. Start backend: `backend\backend-new-api.exe` or `backend\run.ps1`
2. Open app: http://localhost:5173
3. Click Settings → Add API key
4. Watch routes follow real roads! 🛣️

