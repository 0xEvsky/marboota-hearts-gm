# Player Component Documentation

## Player Component

The Player component represents a game participant with cards, state, and position at the table. Players are associated with Clients when users sit at the game table.

### Implementation Files
- `table.go` (defined within Table implementation)

### Structure

```go
type PlayerState int

const (
    PlayerUnavailable PlayerState = iota
    PlayerWaiting
    PlayerReady
    PlayerTrumping
    PlayerPlaying
)

type Team int

const (
    TeamA Team = iota
    TeamB
)

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

### Key Fields

- `client`: Reference to the associated Client
- `state`: Current player state (Unavailable, Waiting, Ready, Trumping, Playing)
- `hand`: Array of Card objects in the player's hand
- `seat`: Seat position (1-4)
- `team`: Team assignment (A or B)
- `score`: Player's score
- `partner`: Reference to partner Player
- `isTurn`: Flag indicating if it's this player's turn

### Player States

1. **PlayerUnavailable**: Seat is empty or player not active
2. **PlayerWaiting**: Player is seated but not ready
3. **PlayerReady**: Player is ready to start the game
4. **PlayerTrumping**: Player is in trump calling phase
5. **PlayerPlaying**: Player is in regular gameplay phase

### Primary Responsibilities

The Player object:
1. Holds game-related data for a participant
2. Maintains a link between the game state and the client
3. Tracks individual player state and progress
4. Contains the player's cards and score

### Team Organization

Players are organized into two teams:
- TeamA: Players in seats 1 and 3
- TeamB: Players in seats 2 and 4

Partners sit across from each other:
- Seat 1's partner is Seat 3
- Seat 2's partner is Seat 4

### Player Lifecycle

1. **Creation**: Player objects are created when the Table is initialized
2. **Activation**: Player becomes active when a Client sits at the position
3. **Ready**: Player moves to ready state when the Client indicates readiness
4. **Active Play**: Player participates in trump calling and regular gameplay
5. **Deactivation**: Player becomes inactive when the Client leaves the seat

### Card Management

- Players receive cards during the deal phase
- Cards are sorted by suit and value in the player's hand
- Players play cards from their hand during their turn

### Turn Handling

- The `isTurn` flag indicates when it's a player's turn to act
- Players can only make game actions during their turn
- Turn advances according to game rules (sequential or based on tricks)

### Example Flow

```
Player position created with Table
  │
  ▼
Client sits at position
  │
  ▼
Player becomes associated with Client
  │
  ▼
Player marks ready
  │
  ▼
Game starts and cards are dealt
  │
  ▼
Player participates in trump phase
  │
  ▼
Player participates in regular gameplay
```

### Interaction with Other Components

- **Table**: The Player belongs to a Table
- **Client**: The Player is linked to a Client
- **Cards**: The Player holds Card objects
- **Other Players**: The Player interacts with other Players through gameplay

### State Transitions

1. **Unavailable → Waiting**: When a client sits at the position
2. **Waiting → Ready**: When the client indicates readiness
3. **Ready → Trumping**: When the game starts
4. **Trumping → Playing**: When the trump phase ends
5. **Any State → Unavailable**: When the client leaves the seat

### Error Handling

- Players can only take actions appropriate to their current state
- Turn validation prevents out-of-turn actions
- Client disconnections are handled by marking the player as unavailable
