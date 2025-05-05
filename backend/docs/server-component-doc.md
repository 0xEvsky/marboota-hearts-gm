# Server Component Documentation

## Server Component

The Server is the top-level component that manages WebSocket connections and game instances.

### Implementation Files
- `main.go` (primary definition and HTTP handlers)
- `messageHandler.go` (message routing)

### Structure

```go
type Server struct {
    mu        sync.Mutex
    conns     map[*websocket.Conn]*Client
    instances map[string]*Instance // key is instanceid
}
```

### Key Fields

- `mu`: Mutex for thread safety
- `conns`: Map of WebSocket connections to Client objects
- `instances`: Map of instance IDs to Instance objects

### Primary Responsibilities

1. **Connection Management**:
   - Accept WebSocket connections
   - Create Client objects for each connection
   - Track active connections
   - Handle disconnections

2. **Instance Management**:
   - Maintain active game instances
   - Create new instances when needed
   - Clean up empty instances

3. **Message Routing**:
   - Read messages from WebSocket connections
   - Route messages to appropriate handlers
   - Send responses back to clients

### Key Methods

- `wsHandler(w http.ResponseWriter, r *http.Request)`: Handles WebSocket connection requests
- `authHandler(w http.ResponseWriter, r *http.Request)`: Handles authentication requests (WIP)
- `read(ws *websocket.Conn)`: Reads messages from WebSocket connection
- `newServer()`: Creates a new Server instance

### Initialization and Lifecycle

1. Server is initialized during application startup
2. HTTP endpoints are registered
3. WebSocket connections are accepted and managed
4. Server runs until application termination

### Thread Safety

The Server uses a mutex to protect concurrent access to its maps of connections and instances. This ensures that operations like adding/removing clients or instances are thread-safe.

### Interaction with Other Components

- **Clients**: The Server creates and manages Client objects
- **Instances**: The Server creates and tracks Instance objects
- **MessageHandler**: The Server delegates message processing to the message handler

### Error Handling

- Connection errors trigger cleanup of the related Client and Instance
- Empty instances are automatically removed
- Message parsing errors are logged and reported to the client

### Example Flow

```
Server startup
  │
  ▼
Client connects via WebSocket
  │
  ▼
Server creates Client object
  │
  ▼
Client sends AUTH message
  │
  ▼
Server associates Client with Instance
  │
  ▼
Server reads messages and routes to handlers
  │
  ▼
Client disconnects
  │
  ▼
Server cleans up Client and Instance
```
