# COMPONENTS/PLAYER.md

## Player Management (`player/player.go`)

### Purpose:
- Represents and manages each player connection.

### Key Attributes:
- Player ID
- WebSocket connection
- Ready status

### Key Functions:
- `NewPlayer()` – Constructs player instance.
- `SetReady()` – Updates readiness.
- `Disconnect()` – Cleanup on leave.

### Interactions:
- Part of a Room.
- Controlled through WS handler.
