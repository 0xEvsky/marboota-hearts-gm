# Client Component Documentation

## Client Component

The Client component represents a connected user and manages WebSocket communication between the server and the user's client application.

### Implementation Files
- `client.go` (primary definition and methods)

### Structure

```go
type ClientState int

const (
    ClientUnavailable ClientState = iota
    ClientIdle
    ClientSeated
)

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

### Key Fields

- `mu`: Mutex for thread safety
- `conn`: WebSocket connection
- `isAuthed`: Authentication status
- `instance`: Reference to the Instance the client is in
- `id`: Unique client identifier
- `name`: Username
- `iconUrl`: URL to the user's icon/avatar
- `state`: Current client state (Unavailable, Idle, Seated)
- `player`: Reference to the Player object if seated
- `requestId`: Current request ID for response matching

### Primary Responsibilities

1. **Communication**:
   - Send messages to the user
   - Receive messages from the user
   - Track request IDs for response matching

2. **State Management**:
   - Track authentication status
   - Manage client state (idle, seated)
   - Maintain connection to Instance and Player

3. **User Identity**:
   - Store user identifier
   - Store username and icon

### Client States

1. **ClientUnavailable**: Initial state, client is not ready for use
2. **ClientIdle**: Client is authenticated but not seated at a table
3. **ClientSeated**: Client is seated at a table and associated with a Player

### Key Methods

- `newClient(conn *websocket.Conn)`: Creates a new Client object
- `writeJson(msg map[string]string)`: Sends a JSON message to the client
- `writeError(msg string)`: Sends an error message to the client
- `writeOk()`: Sends a success response to the client
- `broadcastToMates(msg map[string]string)`: Broadcasts a message to other clients in the same instance

### Thread Safety

The Client uses a mutex to protect concurrent write operations to the WebSocket connection. This ensures that multiple goroutines don't try to write to the connection simultaneously.

### Interaction with Other Components

- **Server**: The Server creates and manages Client objects
- **Instance**: The Client belongs to an Instance
- **Player**: The Client may be associated with a Player
- **ActionHandler**: Action handlers operate on Client objects

### Example Flow

```
Client connects
  │
  ▼
Client authenticates
  │
  ▼
Client joins an Instance
  │
  ▼
Client sits at a table (becomes associated with a Player)
  │
  ▼
Client sends game actions
  │
  ▼
Client disconnects
```

### Communication Protocol

Clients communicate with the server using JSON messages. Each message must have an `ACTION` key that specifies the type of action being performed or event being communicated.

#### Message Structure

```json
{
    "ACTION": "[action type]",
    "REQUESTID": "[optional request id]",
    "[additional fields]": "[values]"
}
```

#### Response Structure

```json
{
    "ACTION": "OK" | "ERROR",
    "REQUESTID": "[echoed request id]",
    "MESSAGE": "[error message if ACTION is ERROR]"
}
```

### Error Handling

- WebSocket write errors are logged
- Invalid actions return error responses to the client
- Client disconnections trigger cleanup procedures

### Authentication Flow

1. Client connects via WebSocket
2. Client sends AUTH message with instance ID, user ID, username, and icon URL
3. Server validates the authentication data
4. Server associates the client with an instance
5. Server notifies other clients in the instance

### Key Interactions

- Clients interact with the Server through WebSocket
- Clients interact with other Clients through broadcasting
- Clients interact with the game through Player objects
