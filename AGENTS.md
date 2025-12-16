# Stick Toss - Agent Context

This document provides context for AI agents working on the Stick Toss codebase.

## ⚠️ Documentation Rules

**IMPORTANT**: When adding or changing features, you MUST update relevant documentation:

- **AGENTS.md** (this file): Update architecture notes, add new gotchas, document design decisions
- **README.md**: Update if adding new features, changing setup/deployment, or modifying API endpoints
- **DEVELOPMENT.md**: Update if changing development workflow or adding new tools
- **QUICKSTART.md**: Update if changing the getting-started experience

Don't leave documentation stale. Future work depends on accurate docs.

## What This Project Is

Stick Toss is a web app for generating balanced teams for pickup hockey games. Users create players with skill ratings (1-5), organize them into groups, and generate fair teams based on skill distribution. It solves the problem of unbalanced pickup games by using a weighted team generation algorithm while maintaining the randomness of traditional "stick toss" team selection.

## Tech Stack

- **Backend**: Go 1.21+ with Gin framework
- **Frontend**: Svelte with Vite (SPA)
- **Database**: SQLite (dev) / PostgreSQL (production)
- **Auth**: JWT tokens (stored in localStorage)
- **Deployment**: Docker multi-stage build

## Architecture Overview

### Request Flow
1. User interacts with Svelte SPA (served as static files)
2. Frontend calls REST API at `/api/*`
3. Auth middleware validates JWT on protected routes
4. Handlers process requests using GORM models
5. Responses returned as JSON

### Key Design Decisions

**Single Binary Deployment**: The Go server serves both the API and static frontend files. In production, everything runs in one container.

**Player Sharing**: Players are NOT scoped to groups. A player entity is created once and can be added to multiple groups. This was an explicit design choice - the same player (with one skill weight) can participate in different groups.

**Ephemeral Teams**: Generated teams are NOT saved to the database. They're computed on-demand and returned to the frontend. This keeps the schema simple and matches the use case (teams change every game).

**Skill Weight is Global**: When you update a player's skill weight, it updates across all groups they're in. This is intentional - a player's skill level is an inherent property, not group-specific.

## Database Schema

```
users
  - id (PK)
  - email (unique)
  - password_hash
  - created_at, updated_at

players
  - id (PK)
  - user_id (FK -> users)
  - name
  - skill_weight (1-5, with CHECK constraint)
  - created_at, updated_at

groups
  - id (PK)
  - user_id (FK -> users)
  - name
  - created_at, updated_at

group_players (junction table)
  - group_id (FK -> groups)
  - player_id (FK -> players)
  - PK (group_id, player_id)
```

**Important Relationships:**
- User owns many Players and Groups
- Players and Groups have many-to-many relationship via group_players
- Deleting a player removes them from all groups
- Deleting a group removes group_players entries but keeps players

## Team Generation Algorithm

Location: `backend/internal/teamgen/teamgen.go`

**How it works:**
1. If locked players are specified, assign them to teams first
2. Collect remaining unassigned players
3. Shuffle remaining players (for randomness)
4. Sort by skill weight descending (for balance)
5. Greedily assign each player to the team with lowest total weight

**Player Locking**: Users can select multiple players and "lock" them together. The algorithm assigns these locked groups to teams first, then balances the remaining players around them. Locked groups are passed as `[][]uint` (array of player ID arrays).

**Minimum Requirements**: Need at least `numTeams` players to generate. The algorithm validates this before processing.

## Skill Level Definitions

These are hardcoded in `frontend/src/lib/store.js` and displayed in tooltips:

- **Level 1 (Bender)**: Learning to skate, very beginner
- **Level 2 (Pylon)**: Can skate but skills haven't developed
- **Level 3 (Solid)**: Competent beer league player
- **Level 4 (Stud)**: Strong player, one of the best on ice
- **Level 5 (Ringer)**: Played high school+, dominates

These descriptions are important context - they're humorous but specific to hockey culture.

## Code Organization

### Backend (`backend/internal/`)

- **`api/`**: HTTP handlers for each resource (auth, players, groups)
  - Each handler is a struct with methods for CRUD operations
  - Handlers validate JWT via middleware and extract user_id from context

- **`auth/`**: JWT generation, validation, bcrypt password hashing
  - `auth.go`: Core auth functions
  - `middleware.go`: Gin middleware that validates Bearer tokens

- **`db/`**: Database connection and configuration
  - Supports both SQLite and PostgreSQL via GORM
  - Config read from env vars (DB_DRIVER, DATABASE_URL)

- **`models/`**: GORM models matching database schema
  - `Migrate()` function runs on startup
  - Foreign key relationships defined with GORM tags

- **`teamgen/`**: Team generation algorithm
  - Single function: `GenerateBalancedTeams()`
  - Returns `[]Team` with player assignments and total weights

### Frontend (`frontend/src/`)

- **`routes/`**: Page components (Login, Signup, Dashboard, Group)
  - `Dashboard.svelte`: List players and groups, CRUD operations
  - `Group.svelte`: Group detail view with team generation

- **`lib/`**: Shared utilities
  - `api.js`: API client functions for all endpoints
  - `store.js`: Svelte stores and skill level definitions

## Important Gotchas

### Docker Volume Mounts
The `docker-compose.yml` has volume mounts commented out. If uncommented, they override built files with empty local directories. Keep them commented unless doing active development with local changes.

### CORS Configuration
The backend allows `localhost:5173` (Vite dev server) and `localhost:3000`. If you add new frontend dev ports, update `backend/cmd/server/main.go`.

### Frontend Routing
The SPA uses `svelte-routing`. The backend's NoRoute handler serves `index.html` for all non-API routes, enabling client-side routing to work.

### SQLite vs PostgreSQL
- SQLite requires CGO, so we build with `CGO_ENABLED=1` in Docker
- The Dockerfile installs `gcc`, `musl-dev`, `sqlite-dev` for this reason
- PostgreSQL doesn't need CGO but requires a separate server

## Common Tasks

### Adding a New API Endpoint
1. Add handler method to appropriate handler in `backend/internal/api/`
2. Register route in `backend/cmd/server/main.go`
3. Add function to `frontend/src/lib/api.js`
4. Use in Svelte component

### Modifying the Database Schema
1. Update model structs in `backend/internal/models/models.go`
2. Restart server - GORM auto-migrates
3. Update API handlers if needed
4. Update frontend types/interfaces if needed

### Adding New Player Properties
Be aware that players are shared across groups. Any new property should be "inherent" to the player (like skill_weight), not group-specific. If you need group-specific data, add it to the junction table or create a new table.

## Development Workflow

**Local Development:**
- Run backend: `cd backend && go run cmd/server/main.go`
- Run frontend: `cd frontend && npm run dev`
- Access at `http://localhost:5173` (Vite proxies API requests)

**Docker Development:**
- `docker-compose up --build`
- Access at `http://localhost:8080`
- Uses PostgreSQL instead of SQLite

**Production Build:**
- `docker build -t sticktoss:latest .`
- Multi-stage build: frontend → Go → Alpine
- Single binary serves everything

## Testing Notes

There are no automated tests yet. Manual testing covers:
- User signup/login/logout
- Player CRUD operations
- Group CRUD operations
- Adding/removing players from groups
- Team generation with various player counts
- Player locking functionality
- Weight visibility toggle

## Future Considerations

If extending this app, keep in mind:
- Players being shared across groups is a core design - don't break this
- Team generation is synchronous - could be async for large groups
- No pagination yet - could be needed for users with many players/groups
- Auth tokens stored in localStorage - could use httpOnly cookies for better security
- No password reset flow yet
- No team history or statistics tracking
