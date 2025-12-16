# Development Guide

This guide covers development workflows for Stick Toss.

## Prerequisites

- **Go 1.21+**: Download from [golang.org](https://golang.org/dl/)
- **Node.js 20+**: Download from [nodejs.org](https://nodejs.org/)
- **Docker & Docker Compose**: For containerized development (optional)

## Getting Started

1. **Install dependencies**
   ```bash
   make install
   ```

2. **Run backend** (terminal 1)
   ```bash
   make dev-backend
   ```
   Backend runs on http://localhost:8080

3. **Run frontend** (terminal 2)
   ```bash
   make dev-frontend
   ```
   Frontend dev server runs on http://localhost:5173

The frontend dev server proxies API requests to the backend automatically.

## Development Workflow

### Backend Development

The backend is written in Go using the Gin framework. Source code is in `backend/`.

**Project structure:**
- `cmd/server/main.go` - Application entry point
- `internal/api/` - HTTP handlers for routes
- `internal/auth/` - JWT authentication logic
- `internal/db/` - Database connection and configuration
- `internal/models/` - Data models (GORM)
- `internal/teamgen/` - Team generation algorithm

**Running tests:**
```bash
make test-backend
```

**Format code:**
```bash
cd backend && go fmt ./...
```

**Database migrations:**
GORM auto-migrates on startup. Models are defined in `internal/models/models.go`.

### Frontend Development

The frontend is built with Svelte and Vite. Source code is in `frontend/src/`.

**Project structure:**
- `src/App.svelte` - Root component with router
- `src/routes/` - Page components
- `src/lib/` - API client and shared utilities
- `src/components/` - Reusable components (if any)

**Building for production:**
```bash
make build-frontend
```

The build outputs to `frontend/dist/` which the Go backend serves in production.

## Database

### SQLite (Default for Development)

By default, the backend uses SQLite with a file `sticktoss.db` in the project root.

No setup required - it will create the database on first run.

### PostgreSQL (Recommended for Production)

1. **Start PostgreSQL with docker-compose:**
   ```bash
   docker-compose up postgres
   ```

2. **Configure environment:**
   ```bash
   export DB_DRIVER=postgres
   export DATABASE_URL="postgres://sticktoss:sticktoss_dev_password@localhost:5432/sticktoss?sslmode=disable"
   ```

3. **Run backend:**
   ```bash
   make dev-backend
   ```

## Docker Development

Run the entire stack with Docker Compose:

```bash
make docker-compose-up
```

This starts:
- PostgreSQL on port 5432
- Backend + Frontend on port 8080

Access at http://localhost:8080

## Environment Variables

Create a `.env` file in the project root (copy from `.env.example`):

```bash
cp .env.example .env
```

**Available variables:**
- `DB_DRIVER` - `sqlite` or `postgres`
- `DATABASE_URL` - Connection string
- `JWT_SECRET` - Secret for signing JWT tokens
- `PORT` - Server port (default: 8080)
- `GIN_MODE` - `debug` or `release`

## API Testing

You can test the API with curl, Postman, or any HTTP client.

**Example: Sign up**
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**Example: Login**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**Example: Create player (requires auth)**
```bash
curl -X POST http://localhost:8080/api/players \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"John Doe","skill_weight":3}'
```

## Deployment

### Building Production Image

```bash
make docker-build
```

### Running Production Container

```bash
docker run -p 8080:8080 \
  -e DB_DRIVER=postgres \
  -e DATABASE_URL="your-postgres-url" \
  -e JWT_SECRET="your-secret" \
  -e GIN_MODE=release \
  sticktoss:latest
```

### Deploying to Miren

Assuming you have Miren CLI configured:

```bash
miren deploy
```

Make sure to set environment variables in Miren dashboard:
- `DB_DRIVER=postgres`
- `DATABASE_URL=<provided-by-miren>`
- `JWT_SECRET=<generate-random-secret>`
- `GIN_MODE=release`

## Troubleshooting

### Backend won't start
- Check that port 8080 is not in use
- Verify Go is installed: `go version`
- Check database connection string

### Frontend won't start
- Check that port 5173 is not in use
- Verify Node.js is installed: `node --version`
- Delete `node_modules` and reinstall: `cd frontend && rm -rf node_modules && npm install`

### CORS errors
- Make sure frontend is accessing the backend through the Vite proxy (http://localhost:5173, not http://localhost:8080)
- Check CORS configuration in `backend/cmd/server/main.go`

### Database errors
- Delete the SQLite database file and restart: `rm sticktoss.db`
- Check PostgreSQL is running: `docker-compose ps`
- Verify DATABASE_URL is correct

## Code Style

**Go:**
- Use `go fmt` to format code
- Follow standard Go conventions
- Use meaningful variable names

**Svelte/JavaScript:**
- Use 2 spaces for indentation
- Prefer const/let over var
- Use async/await for async code

## Making Changes

1. Create a new branch
2. Make your changes
3. Test locally
4. Format code
5. Commit and push
6. Open a pull request

## Useful Commands

```bash
# See all available make commands
make help

# Clean build artifacts
make clean

# Format all code
make fmt

# Run full stack with Docker
make docker-compose-up

# Stop Docker services
make docker-compose-down
```
