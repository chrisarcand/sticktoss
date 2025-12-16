# API Documentation

This document describes the REST API endpoints for Stick Toss.

## Base URL

All API endpoints are prefixed with `/api`.

## Authentication

Protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

Tokens are obtained by logging in or signing up.

## Endpoints

### Authentication

#### Sign Up
```
POST /api/auth/signup
```

Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2025-01-15T10:00:00Z"
  }
}
```

#### Login
```
POST /api/auth/login
```

Authenticate and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2025-01-15T10:00:00Z"
  }
}
```

#### Get Current User
```
GET /api/auth/me
```

Get the authenticated user's information.

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "created_at": "2025-01-15T10:00:00Z"
}
```

### Players

All player endpoints require authentication.

#### List Players
```
GET /api/players
```

Get all players for the authenticated user.

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "John Doe",
    "skill_weight": 4,
    "created_at": "2025-01-15T10:00:00Z"
  }
]
```

#### Get Player
```
GET /api/players/:id
```

Get a specific player by ID.

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "John Doe",
  "skill_weight": 4,
  "created_at": "2025-01-15T10:00:00Z"
}
```

#### Create Player
```
POST /api/players
```

Create a new player.

**Request Body:**
```json
{
  "name": "John Doe",
  "skill_weight": 4
}
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "John Doe",
  "skill_weight": 4,
  "created_at": "2025-01-15T10:00:00Z"
}
```

#### Update Player
```
PUT /api/players/:id
```

Update a player's information. Updates apply across all groups.

**Request Body:**
```json
{
  "name": "John Smith",
  "skill_weight": 5
}
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "John Smith",
  "skill_weight": 5,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T11:00:00Z"
}
```

#### Delete Player
```
DELETE /api/players/:id
```

Delete a player. Removes them from all groups.

**Response:**
```json
{
  "message": "player deleted"
}
```

### Groups

All group endpoints require authentication.

#### List Groups
```
GET /api/groups
```

Get all groups for the authenticated user.

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "Tuesday Night Hockey",
    "created_at": "2025-01-15T10:00:00Z"
  }
]
```

#### Get Group
```
GET /api/groups/:id
```

Get a specific group with its players.

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "Tuesday Night Hockey",
  "created_at": "2025-01-15T10:00:00Z",
  "players": [
    {
      "id": 1,
      "name": "John Doe",
      "skill_weight": 4
    }
  ]
}
```

#### Create Group
```
POST /api/groups
```

Create a new group.

**Request Body:**
```json
{
  "name": "Tuesday Night Hockey"
}
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "Tuesday Night Hockey",
  "created_at": "2025-01-15T10:00:00Z"
}
```

#### Update Group
```
PUT /api/groups/:id
```

Update a group's name.

**Request Body:**
```json
{
  "name": "Wednesday Night Hockey"
}
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "Wednesday Night Hockey",
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T11:00:00Z"
}
```

#### Delete Group
```
DELETE /api/groups/:id
```

Delete a group. Players remain in the database.

**Response:**
```json
{
  "message": "group deleted"
}
```

#### Add Player to Group
```
POST /api/groups/:id/players
```

Add an existing player to a group.

**Request Body:**
```json
{
  "player_id": 1
}
```

**Response:**
```json
{
  "message": "player added to group"
}
```

#### Remove Player from Group
```
DELETE /api/groups/:id/players/:player_id
```

Remove a player from a group. The player remains in the database.

**Response:**
```json
{
  "message": "player removed from group"
}
```

#### Generate Teams
```
POST /api/groups/:id/generate-teams
```

Generate balanced teams from the group's players.

**Request Body:**
```json
{
  "num_teams": 2,
  "locked_players": [
    [1, 2],
    [3, 4]
  ]
}
```

- `num_teams`: Number of teams to create (minimum 2)
- `locked_players`: (Optional) Array of player ID arrays. Players in the same inner array will be placed on the same team.

**Response:**
```json
{
  "teams": [
    {
      "number": 1,
      "players": [
        {
          "id": 1,
          "name": "John Doe",
          "skill_weight": 4
        }
      ],
      "total_weight": 12
    },
    {
      "number": 2,
      "players": [
        {
          "id": 2,
          "name": "Jane Smith",
          "skill_weight": 3
        }
      ],
      "total_weight": 11
    }
  ]
}
```

## Error Responses

All endpoints may return error responses:

**401 Unauthorized:**
```json
{
  "error": "authorization header required"
}
```

**404 Not Found:**
```json
{
  "error": "player not found"
}
```

**400 Bad Request:**
```json
{
  "error": "invalid request body"
}
```

**500 Internal Server Error:**
```json
{
  "error": "internal server error"
}
```
