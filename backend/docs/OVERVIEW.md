# OVERVIEW.md

## 🧩 General Overview of the Backend Server

The **Marboota backend** is implemented in **Go (Golang)** and serves as the core real-time server that facilitates gameplay logic, room management, and player synchronization in a multiplayer environment. It leverages **WebSockets** for efficient two-way communication between clients (players) and the server.

### Key Capabilities:
- Persistent WebSocket connections for real-time data exchange
- Dynamic game room creation and player management
- Game session lifecycle management (start, in-progress, end)
- Broadcasting of player and game events
- Lightweight, scalable Go-based architecture

### Key Components:
- `main.go` – Entry point, sets up HTTP and WebSocket server
- `ws/` – WebSocket connection and message dispatcher
- `room/` – Game room lifecycle and coordination
- `player/` – Player instance and state management
- `message/` – Structured message definitions and helpers
- `utils/` – Shared utility functions like ID generation

### High-Level Flow:
1. A client connects via WebSocket.
2. The server assigns the client to a room or creates one.
3. Game room tracks player readiness.
4. Once all players are ready, the game begins.
5. Game events are transmitted as structured messages.
6. On disconnect or game end, resources are cleaned up.
