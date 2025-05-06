# Marboota Backend Documentation

## Overview

Marboota is a card game server written in Go that manages game instances, player connections, and game state through WebSockets. This document explains the server architecture, components, and flow to help new contributors understand the system.

## System Architecture

The server follows a hierarchical structure:

```
Server
├── Instances (Game Rooms)
│   ├── Clients (Connected Users)
│   └── Table (Game State)
│       ├── Players (Game Participants)
│       └── Trump (Game Phase State)
```

### Key Components

1. **Server**: The top-level component that manages all WebSocket connections and game instances
2. **Instance**: A game room that groups clients together and contains a table
3. **Client**: A connected user with a WebSocket connection
4. **Table**: The game state manager with player positions and game phases
5. **Player**: A game participant with cards, state, and position
6. **Card**: Game cards with suit and value

## Connection Flow

1. Client connects via WebSocket (`/ws` endpoint)
2. Server creates a new `Client` object
3. Client authenticates via the `AUTH` action
4. Server associates the client with an `Instance` (or creates a new one)
5. Game setup begins when players sit and get ready

## Game Flow

1. **Waiting Phase**: Players join, sit, and mark themselves as ready
2. **Trump Phase**: When all players are ready, the game deals cards and enters the trump calling phase
3. **Playing Phase**: After trump is determined, regular gameplay begins

## Component Details

### Server (`main.go`, `messageHandler.go`)

- Handles WebSocket connections and upgrades
- Routes messages to appropriate handlers
- Maintains maps of active connections and game instances

```go
type Server struct {
    mu        sync.Mutex
    conns     map[*websocket.Conn]*Client
    instances map[string]*Instance
}
```

### Client (`client.go`)

- Represents a connected user
- Manages WebSocket communication
- Maintains user info (ID, name, icon URL)
- Tracks client state (idle, seated)

```go
type Client struct {
    mu        sync.Mutex
    conn      *websocket.Conn
    isAuthed  bool
    instance  *Instance
    id        string
    name      string
    iconUrl   string
    state     ClientState
    player    *Player
    requestId string
}
```

### Instance (`instance.go`)

- Represents a game room
- Groups clients together
- Contains a table for gameplay

```go
type Instance struct {
    mu      sync.Mutex
    id      string
    clients map[string]*Client
    table   Table
}
```

### Table (`table.go`)

- Manages game state and phases
- Holds player positions
- Coordinates turn handling

```go
type Table struct {
    instance   *Instance
    players    [4]*Player
    state      TableState
    turn       int
    turnOffset int
    trump      Trump
}
```

### Player (`table.go`)

- Represents a game participant
- Holds player's hand of cards
- Tracks player state and score

```go
type Player struct {
    client  *Client
    state   PlayerState
    hand    []Card
    seat    int
    team    Team
    score   int
    partner *Player
    isTurn  bool
}
```

### Card (`card.go`)

- Represents playing cards
- Has suit, value, and display name

```go
type Card struct {
    name  string
    suit  suit
    value int
}
```

## Message Protocol

The server uses a JSON-based protocol over WebSockets as described in `PROTOCOL.md`. All values are sent as strings, and each message must include an `ACTION` field.

### Client to Server Actions

- `AUTH`: Authenticate with the server
- `SIT`: Request to sit at a specific seat
- `UNSIT`: Leave a seat
- `READY`: Mark as ready to play
- `UNREADY`: Mark as not ready to play
- `TRUMPCALL`: Make a trump call or pass

### Server to Client Events

- `JOIN`: Notification of a client joining
- `LEAVE`: Notification of a client leaving
- `SIT`: Notification of a client sitting
- `UNSIT`: Notification of a client unseating
- `READY`: Notification of a client being ready
- `UNREADY`: Notification of a client being unready
- `DEAL`: Notification of cards dealt to the player
- `TRUMPSTART`: Notification that the trump phase has started
- `YOURTRUMPCALL`: Notification that it's the player's turn to call trump
- `TRUMPCALL`: Notification that another player made a trump call

## Action Handlers (`actionHandler.go`)

The server implements several action handlers:

- `authClient`: Processes authentication requests
- `seatClient`: Handles seating requests
- `unseatClient`: Handles unseating requests
- `setReady`: Sets a player as ready
- `unsetReady`: Sets a player as not ready
- `advanceTrump`: Processes trump calls

## Error Handling

The server responds with error messages for invalid actions:
- Authentication errors
- Seating errors (seat taken, invalid seat)
- Game state errors (trying to sit when game already started)
- Turn errors (acting out of turn)

## System Flow Examples

### User Connection and Authentication

1. User connects to WebSocket endpoint
2. Server assigns a new client object
3. Client sends AUTH message with ID, name, and icon
4. Server validates and associates client with an instance
5. Server broadcasts JOIN message to other clients in the instance

### Game Setup

1. Client sends SIT request with seat number
2. Server validates and seats the client
3. Server broadcasts SIT message to other clients
4. Client sends READY message
5. Server validates and marks client as ready
6. When all players are ready, server starts the trump phase

### Trump Phase

1. Server deals cards to all players
2. Server broadcasts TRUMPSTART
3. Server notifies the first player with YOURTRUMPCALL
4. Player makes a trump call or passes
5. Server processes the call and advances to the next player
6. When trump is determined, the playing phase begins

## Contributing

To contribute to the server:

1. Understand the component structure and message flow
2. Identify which component needs modification
3. Ensure thread safety by using mutex locks when needed
4. Follow the existing protocol for message handling
5. Test your changes with the client implementation

## Running Locally

- Install Go v1.24
- Run `go run .` for development
- Run `go build .` to compile a binary
