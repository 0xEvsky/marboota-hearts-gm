# COMPONENTS/WS.md

## WebSocket Handler (`ws/ws.go`)

### Purpose:
- Manages WebSocket connections.
- Listens to client messages and dispatches to appropriate handlers.

### Key Functions:
- `HandleConnections()`: Handles upgrades and connections.
- `HandleMessages()`: Central message dispatcher.

### Interactions:
- Talks to `room/` and `player/` packages to manage state.
- Uses `message/` definitions for parsing and sending messages.
