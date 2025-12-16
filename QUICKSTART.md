# Quick Start Guide

Your Stick Toss app is ready to go! Here are the fastest ways to get started:

## Option 1: Docker (Recommended - Easiest)

```bash
# Start the app with Docker Compose (includes database)
docker-compose up

# Access at: http://localhost:8080
```

That's it! The app is running with PostgreSQL.

To stop:
```bash
docker-compose down
```

## Option 2: Simple Docker Run

If you just want to run the built image with SQLite:

```bash
docker run -p 8080:8080 \
  -e DB_DRIVER=sqlite \
  -e DATABASE_URL=/app/sticktoss.db \
  sticktoss:latest
```

Access at: http://localhost:8080

## Option 3: Local Development (No Docker)

### First Time Setup

1. **Install dependencies**
   ```bash
   # Backend
   cd backend
   go mod download
   cd ..

   # Frontend
   cd frontend
   npm install
   cd ..
   ```

### Running

Open two terminal windows:

**Terminal 1 - Backend:**
```bash
cd backend
go run cmd/server/main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

Access at: http://localhost:5173 (frontend dev server with hot reload)

## What to Do First

1. **Sign up** - Create an account at the signup page
2. **Add players** - Create some players with skill weights (1-5)
3. **Create a group** - Make a group like "Tuesday Night Hockey"
4. **Add players to the group** - Select which players are playing tonight
5. **Generate teams** - Click the generate button and get balanced teams!

## Features to Try

- **Lock players together**: Select multiple players and click "Lock Selected Players Together" to keep them on the same team
- **Show/hide weights**: Toggle the weight display on generated teams
- **Multiple groups**: Create different groups for different game nights or leagues
- **Edit players**: Update player names or skill levels anytime (updates across all groups)

## Environment Variables

Create a `.env` file if needed:

```bash
DB_DRIVER=sqlite
DATABASE_URL=sticktoss.db
JWT_SECRET=your-secret-key-here
PORT=8080
GIN_MODE=debug
```

## Troubleshooting

**Port already in use?**
- Change the port: `PORT=3000 go run cmd/server/main.go`
- Or stop whatever's using 8080: `lsof -ti:8080 | xargs kill`

**Frontend won't connect to backend?**
- Make sure backend is running on port 8080
- Check the Vite proxy config in `frontend/vite.config.js`

**Database errors?**
- Delete the database file: `rm sticktoss.db`
- Restart the server (it will create a fresh database)

## Next Steps

- Check out `README.md` for full documentation
- See `DEVELOPMENT.md` for detailed dev info
- Read `INSTRUCTIONS.md` for the original requirements

Enjoy making your hockey games more balanced!
