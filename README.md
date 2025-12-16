![Stick Toss Logo](images/logo.png)

# Stick Toss

A weighted team randomizer for pickup pond hockey games. Because sometimes the teams that come out of the literal stick pile are hilariously unbalanced, and nobody wants to play a 12-2 snoozefest when you've only got the rink for an hour. Stick Toss lets you assign skill weights to players and generates random-but-fair teams, so you can keep the spirit of the stick pile while actually having competitive games. Perfect for beer league organizers, pond hockey groups, or anyone tired of watching one team cycle the puck for 45 minutes while the other team practices their defensive zone faceoffs.

## Features

- **Player Management**: Create and manage players with skill weights (1-5, from Bender to Ringer)
- **Groups**: Organize players into different groups for different games or leagues
- **Smart Team Generation**: Algorithm balances teams based on total skill weights
- **Player Locking**: Lock specific players together on the same team
- **Mobile Responsive**: Works great on phones for rink-side team generation
- **User Accounts**: Each user manages their own players and groups

## Tech Stack

- **Backend**: Go with Gin framework
- **Frontend**: Svelte with Vite
- **Database**: SQLite (dev) / PostgreSQL (production)
- **Auth**: JWT tokens
- **Deployment**: Docker

## Quick Start

### Run with Docker Compose (Recommended)

The easiest way to get started:

```bash
docker-compose up --build
```

Access the app at **http://localhost:8080**

This runs the full stack with PostgreSQL. To stop: `docker-compose down`

### Deploy to Miren

TODO

## Documentation

**[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)** - Development guide:
- Configuration and environment variables
- Local development setup (Go + Node.js)
- Building production images
- Running tests
- Detailed troubleshooting
- API testing examples

**[docs/API.md](docs/API.md)** - Complete REST API reference

## Project Structure

```
sticktoss/
├── backend/
│   ├── cmd/server/          # Main application entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   ├── auth/            # Authentication & JWT
│   │   ├── db/              # Database connection
│   │   ├── models/          # Data models
│   │   └── teamgen/         # Team generation algorithm
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── lib/             # API client & stores
│   │   ├── routes/          # Svelte pages
│   │   └── App.svelte
│   └── package.json
├── docs/
│   ├── DEVELOPMENT.md       # Development guide
│   └── API.md               # API documentation
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## License

MIT

## Contributing

Pull requests welcome! You're here, aren't ya?
