# Instance Component Documentation

## Instance Component

The Instance component represents a game room that groups clients together and contains a table for gameplay.

### Implementation Files
- `instance.go` (primary definition and methods)

### Structure

```go
type Instance struct {
    mu      sync.Mutex
    id      string
    clients map[string]*Client // key is userid
    table   Table
}
```

### Key Fields

- `mu`: Mutex for thread safety
- `id`: Unique instance identifier
- `clients`: Map of user IDs to Client objects
- `table`: Game table that manages the game state

### Primary Responsibilities

1. **Client Management**:
   - Group clients together in a game room
   - Track active clients
   - Manage client joining and leaving

2. **Game Management**:
   - Host a Table for gameplay
   - Maintain game state
   - Coordinate communication between clients and table

3. **Broadcasting**:
   - Send messages to all clients in the instance

### Key Methods

- `newInstance(c *Client, id string)`: Creates a new Instance
- `joinInstance(c *Client, id string)`: Adds a client to an existing Instance or creates a new one
- `Broadcast(msg map[string]string)`: Sends a message to all clients in the Instance

### Thread Safety

The Instance uses a mutex to protect concurrent access to its client map and table state. This ensures that operations like adding/removing clients are thread-safe.

### Interaction with Other Components

- **Server**: The Server creates and manages Instance objects
- **Clients**: The Instance groups Client objects together
- **Table**: The Instance contains a Table for gameplay

### Lifecycle

1. Creation: Instance is created when a client joins a non-existent instance
2. Active: Clients join and leave, game play occurs
3. Termination: Instance is removed when the last client leaves

### Example Flow

```
Client requests to join an instance
  │
  ▼
Server creates/finds the Instance
  │
  ▼
Instance adds the Client to its map
  │
  ▼
Client interacts with the Table
  │
  ▼
Game progresses in the Instance
  │
  ▼
All Clients leave
  │
  ▼
Server removes the empty Instance
```

### Broadcasting

The Instance provides broadcasting capability to send messages to all authenticated clients in the instance. This is used for:

- Game state updates
- Player actions
- Client join/leave notifications

### Management Functions

When clients join an instance:
1. The client is added to the instance's client map
2. The client is informed of other clients already in the instance (catch-up)
3. Other clients are informed of the new client

When clients leave an instance:
1. The client is removed from the instance's client map
2. Other clients are informed of the client's departure
3. If the instance becomes empty, it is removed from the server

### Instance Creation

Instances are created on-demand when a client attempts to join an instance that doesn't exist. The creation process:

1. Create a new Instance object
2. Create a new Table object
3. Associate the Table with the Instance
4. Add the client to the Instance
5. Add the Instance to the Server's instances map

### Error Handling

- Client disconnections are handled by removing the client from the instance
- Empty instances are automatically cleaned up by the server
- Instance-level errors are reported to the clients
