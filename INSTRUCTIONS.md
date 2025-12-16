You will create this web app for me.

The web app should be very easy to use, with a responsive layout that works well on both desktop and mobile devices.

## Technology Stack

### Backend
- **Go** with **Gin** framework
- REST API with JWT authentication
- Database abstraction layer supporting:
  - **SQLite** for local development
  - **PostgreSQL** for production deployment

### Frontend
- **Svelte** with **Vite**
- Single Page Application (SPA)
- Communicates with backend via REST API
- JWT stored in httpOnly cookie or localStorage

### Database Schema
```
users
  - id (primary key)
  - email (unique)
  - password_hash
  - created_at

players
  - id (primary key)
  - user_id (foreign key to users)
  - name
  - skill_weight (1-5)
  - created_at

groups
  - id (primary key)
  - user_id (foreign key to users)
  - name
  - created_at

group_players (junction table for many-to-many)
  - group_id (foreign key to groups)
  - player_id (foreign key to players)
  - primary key (group_id, player_id)
```

**Important**: Players are shared entities across groups. Each player has ONE skill weight that applies universally. When a player's weight is updated, it updates across all groups they belong to. Players can be added to multiple groups.

### Project Structure
```
sticktoss/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── api/         # HTTP handlers
│   │   ├── auth/        # JWT middleware and auth logic
│   │   ├── db/          # Database layer and migrations
│   │   ├── models/      # Data models and structs
│   │   └── teamgen/     # Team balancing algorithm
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── routes/
│   │   └── lib/
│   ├── package.json
│   └── vite.config.js
├── Dockerfile           # Multi-stage build
├── docker-compose.yml   # Local development
└── README.md
```

### Deployment
- Docker container with multi-stage build:
  1. Build Svelte frontend to static files
  2. Build Go binary
  3. Single Go binary serves both API and static frontend files
- Designed for deployment to Miren (Heroku-like platform)
- Uses PostgreSQL in production

## Features

- A simple log in and sign up system to create and manage user accounts.
- Users should be able to create 'groups', which will contain a list of players. This will allow users to manage different sets of players for different games or leagues, to easily switch between them and keep saved lists of players and their weights.
- Users should be able to create players with a name and a skill weight from 1 to 5. Players can be added to multiple groups.
- Within each group, users should be able to add existing players or create new players.
- Users should be able to edit and delete players. Editing a player's weight updates it across all groups.
- Users should be able to remove players from specific groups without deleting the player entirely.
- When a group is selected, users should be able to generate random teams based on the players and their weights. Generated teams are ephemeral (not saved to database).
- The screen which shows the teams should have a small, subtle button for showing or hiding the weight numbers next to each player's name.
- The team generation algorithm should aim to create teams that are as balanced as possible based on the total skill weights of each team.
- The user should have some way to be able to keep players 'together' on a same team if desired, and generate the rest of the teams around those locked players.
- When picking weights, there should be some sort of tooltip or info icon next to the weight selection that, when hovered over or clicked, shows a description of what each weight level means in terms of player skill.

Here's the player level descriptions for the skill weights:

Level 1: Bender
Can't stop without using the boards. Skating is labored and unstable—ankles rolling, choppy strides, limited mobility. Stickhandling is minimal; mostly just tries to get a piece of the puck. Shooting form is rough. Positional awareness is absent—chases the puck everywhere. Struggles with basic rules and gameplay concepts. Probably learning to skate as an adult or just started playing hockey.

Level 2: Pylon
Can skate forward with moderate stability but crossovers are rough and backwards skating is shaky. Stops and starts without total control. Can receive and make simple passes but struggles under pressure. Shot exists but lacks consistency or power. Understands basic positions but gets caught out of place regularly. Easy to play around—opponents skate past like you're a stationary orange cone at practice. Might have been playing for years but the skills just haven't developed. Shows up and tries, which counts for something.

Level 3: Solid
Competent skater in all directions with decent speed and agility. Can execute crossovers and transitions reasonably well. Handles the puck confidently in open ice and makes smart, tape-to-tape passes most of the time. Has a serviceable shot with occasional accuracy. Understands positioning and reads plays adequately. Reliable player who won't hurt the team but won't dominate either. May have played a bit growing up or has several years of dedicated beer league experience.

Level 4: Stud
Strong, fluid skater with good edge work and quick acceleration. Handles the puck well in traffic and can make plays under pressure. Consistently accurate passer who sees the ice well. Has a legitimately dangerous shot—goalies respect it. Solid hockey IQ with good positioning and anticipation. Didn't grow up playing organized hockey but has excellent athleticism and years of experience that shows. Often one of the best players on the ice in most beer league games.

Level 5: Ringer
Effortless, powerful skating with elite edge work and explosive speed. Stickhandles in a phone booth and protects the puck naturally. Makes high-difficulty passes look routine and sees plays developing before they happen. Can pick corners consistently and has a legitimately hard, accurate shot. Reads the game at a different speed than everyone else—always in the right position. Almost certainly played high school hockey, at least. Makes everyone else look slow. The guy who "takes it easy" and still dominates.
