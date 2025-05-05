# COMPONENTS/ROOM.md

## Room Management (`room/room.go`)

### Purpose:
- Tracks game rooms and their state.
- Manages lifecycle: join, ready, start, end.

### Key Functions:
- `CreateRoom()` – Initializes new room.
- `AddPlayer(player)` – Adds player to room.
- `StartGame()` – Begins session.
- `RemovePlayer()` – Cleanup.

### Interactions:
- Coordinates with player manager for room-player mapping.
- Informs WS handler of room state changes.
