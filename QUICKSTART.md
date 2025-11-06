# 🚀 Quick Start Guide

Get the Shipt Route Optimizer running in 2 minutes!

## Prerequisites Check

```bash
# Check Go version (need 1.21+)
go version

# Check Node version (need 18+)
node --version
```

## Step 1: Start Backend (Terminal 1)

```bash
cd backend
go mod download
go run cmd/main.go
```

✅ You should see: `🚀 Shipt Route Optimizer Backend starting on :8080`

## Step 2: Start Frontend (Terminal 2)

```bash
cd frontend
npm install
npm run dev
```

✅ Browser should open at: `http://localhost:5173`

## Step 3: Use the App

1. Click **"Load Sample Data"**
2. Click **"Optimize Routes"**
3. See the magic happen! ✨

## Troubleshooting

### Port 8080 already in use?
```bash
# Windows
netstat -ano | findstr :8080

# Kill the process using that port
taskkill /PID <PID> /F
```

### Backend connection error?
- Make sure backend is running on port 8080
- Check firewall settings
- Try accessing http://localhost:8080/api/health directly

### npm install fails?
```bash
# Clear cache and retry
npm cache clean --force
npm install
```

---

**Need help?** Check the full [README.md](README.md) for detailed documentation.

