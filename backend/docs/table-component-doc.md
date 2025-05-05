# Table Component Documentation

## Table Component

The Table component manages the game state, player positions, and coordinates turn handling within a game instance.

### Implementation Files
- `table.go` (primary definition and methods)

### Structure

```go
type TableState int

const (
    TableWaiting TableState = iota
    TableTrumping
    TablePlaying
)

type Table struct {
    instance   *Instance
    players    [4]*Player
    state      TableState
    turn       int
    turnOffset int
    trump      Trump
}

type Trump struct {
    highestCall   int
    highestCaller *Player
    callers       []*Player
}
```

### Key Fields

- `instance`: Reference to the containing Instance
- `players`: Fixed-size array of Player objects (4 seats)
- `state`: Current table state (Waiting, Trumping, Playing)
- `turn`: Current turn index (0-3)
- `turnOffset`: Starting turn offset
- `trump`: Trump phase state information

### Table States

1. **TableWaiting**: Table is waiting for players to join and get ready
2. **TableTrumping**: Trump calling phase is in progress
3. **TablePlaying**: Regular gameplay is in progress

### Primary Responsibilities

1. **Player Management**:
   - Manage player seating
   - Track player positions
   - Handle player ready status

2. **Game Flow Control**:
   - Transition between game phases
   - Deal cards to players
   - Manage turn order

3. **Game Rules**:
   - Enforce trump calling rules
   - Validate player actions
   - Track game progress and scoring

### Key Methods

- `newTable()`: Creates a new Table
- `seatPlayer(c *Client, s int)`: Seats a client at a specific position
- `unseatPlayer(c *Client)`: Removes a client from a seat
- `isEveryoneReady()`: Checks if all players are ready
- `trumpStart()`: Initiates the trump calling phase

### Thread Safety

The Table doesn't have its own mutex, but is protected by the Instance's mutex when accessed. This ensures thread-safe operations on the table state.

### Interaction with Other Components

- **Instance**: The Table belongs to an Instance
- **Players**: The Table manages Player objects
- **Cards**: The Table deals Card objects to Players

### Game Flow

#### Initialization
1. Table is created with the Instance
2. Four Player slots are initialized
3. Table starts in TableWaiting state

#### Player Seating
1. Clients request to sit at specific positions
2. Table assigns clients to Player objects
3. Table broadcasts seating information

#### Game Start
1. Players mark themselves as ready
2. Table checks if all players are ready
3. When all ready, Table transitions to TableTrumping state
4. Cards are dealt to all players

#### Trump Phase
1. Table determines starting player
2. Players make trump calls in turn
3. Trump calls continue until one caller remains
4. Table transitions to TablePlaying state

### Card Management

The Table creates and deals cards to players:
1. Create a standard 52-card deck
2. Shuffle the deck randomly
3. Deal 13 cards to each player
4. Sort cards in each player's hand

### Example Flow

```
Table created in waiting state
  │
  ▼
Players join and sit at the table
  │
  ▼
Players mark themselves as ready
  │
  ▼
Table checks everyone is ready
  │
  ▼
Table deals cards and starts trump phase
  │
  ▼
Players make trump calls in turn
  │
  ▼
Trump is determined
  │
  ▼
Table transitions to playing phase
```

### Error Handling

- Invalid seating requests return errors (seat taken, invalid seat)
- Game state validation prevents out-of-sequence actions
- Turn order is enforced for player actions
